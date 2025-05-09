package handler

import (
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

func (v *ITService) route(dbConnection *gorm.DB, serviceRequestId, sourceServiceStatusId int) {
	err, serviceStatusObject := database.Get(dbConnection, const_util.ITServiceRequestStatusTable, sourceServiceStatusId)
	if err == nil {
		serviceRequestStatus := database.ITServiceRequestStatus{ObjectInfo: serviceStatusObject.ObjectInfo}
		routingConditionList := serviceRequestStatus.GetRequestStatusInfo().RoutingCondition
		unmatchedRoutingUsers := serviceRequestStatus.GetRequestStatusInfo().UnmatchedRoutingUsers
		if len(routingConditionList) > 0 {
			for _, routingCondition := range routingConditionList {
				var SQLCondition string
				condition := routingCondition.Query.Condition
				for index, rule := range routingCondition.Query.Rules {
					if rule.Field == "category" {
						if len(routingCondition.Query.Rules)-1 == index {
							SQLCondition += " object_info->>'$.category' = " + strconv.Itoa(rule.Value)
						} else {
							SQLCondition += " object_info->>'$.category' = " + strconv.Itoa(rule.Value) + " " + condition
						}
					} else if rule.Field == "sub_category" {
						if len(routingCondition.Query.Rules)-1 == index {
							SQLCondition += " object_info->>'$.subCategory' = " + strconv.Itoa(rule.Value)
						} else {
							SQLCondition += " object_info->>'$.subCategory' = " + strconv.Itoa(rule.Value) + " " + condition
						}
					}

				}
				SQLCondition += " and id = " + strconv.Itoa(serviceRequestId)
				v.BaseService.Logger.Info("generated condition", zap.Any("condition", SQLCondition))
				listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.ITServiceRequestTable, SQLCondition)
				if err == nil {
					// should be always one, not more than one, otherwise, something is wrong.
					if len(*listOfServiceRequest) == 1 {
						itServiceRequest := database.ITServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
						notificationList := routingCondition.NotificationUserList
						for _, notificationUserId := range notificationList {
							v.NotificationHandler.EmailGenerator(dbConnection, const_util.SubmitForExecutionEmailTemplateType, notificationUserId, const_util.ITServiceMyReviewRequestComponent, itServiceRequest.Id)
							v.BaseService.Logger.Info("routing request email", zap.Any("user_id", notificationUserId))
						}
					}
				}

			}
		} else {
			SQLConditionRequestId := " id = " + strconv.Itoa(serviceRequestId)
			v.BaseService.Logger.Info("generated condition", zap.Any("condition", SQLConditionRequestId))
			listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.ITServiceRequestTable, SQLConditionRequestId)
			if err == nil {
				// should be always one, not more than one, otherwise, something is wrong.
				if len(*listOfServiceRequest) == 1 {
					itServiceRequest := database.ITServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
					for _, notificationUserId := range unmatchedRoutingUsers {
						v.NotificationHandler.EmailGenerator(dbConnection, const_util.SubmitForExecutionEmailTemplateType, notificationUserId, const_util.ITServiceMyReviewRequestComponent, itServiceRequest.Id)
						v.BaseService.Logger.Info("un matched routing request email", zap.Any("user_id", notificationUserId))
					}
				}
			}

		}
	}
}

func (v *ITService) getTargetWorkflowGroup(dbConnection *gorm.DB, requestId int, workflowId int) []int {
	err, serviceStatusObject := database.Get(dbConnection, const_util.ITServiceRequestStatusTable, workflowId)
	var listOfRequestIds []int
	if err == nil {
		serviceRequestStatus := database.ITServiceRequestStatus{ObjectInfo: serviceStatusObject.ObjectInfo}
		routingConditionList := serviceRequestStatus.GetRequestStatusInfo().RoutingCondition
		fmt.Println("routingConditionList : ", routingConditionList)
		//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers
		if len(routingConditionList) > 0 {
			for _, routingCondition := range routingConditionList {
				var SQLCondition string
				SQLCondition = "( "
				condition := routingCondition.Query.Condition
				notificationList := routingCondition.NotificationUserList

				for index, rule := range routingCondition.Query.Rules {
					if rule.Field == "category" {
						if len(routingCondition.Query.Rules)-1 == index {
							SQLCondition += " object_info->>'$.categoryId' = " + strconv.Itoa(rule.Value)
						} else {
							SQLCondition += " object_info->>'$.categoryId' = " + strconv.Itoa(rule.Value) + " " + condition
						}
					} else if rule.Field == "sub_category" {
						if len(routingCondition.Query.Rules)-1 == index {
							SQLCondition += " object_info->>'$.subCategory' = " + strconv.Itoa(rule.Value)
						} else {
							SQLCondition += " object_info->>'$.subCategory' = " + strconv.Itoa(rule.Value) + " " + condition
						}
					}

				}

				SQLCondition += ") and  id = " + strconv.Itoa(requestId)
				fmt.Println("SQLCondition: ", SQLCondition)
				listOfRequests, err := database.GetConditionalObjects(dbConnection, const_util.ITServiceRequestTable, SQLCondition)
				if err == nil {
					if len(*listOfRequests) > 0 {
						// yes we found , now check who is tagged
						for _, existingRequestId := range notificationList {
							listOfRequestIds = append(listOfRequestIds, existingRequestId)
						}
					}
					listOfRequestIds = util.RemoveDuplicateInt(listOfRequestIds)
				}

			}
		}
		if len(listOfRequestIds) == 0 {
			listOfRequestIds = serviceRequestStatus.GetRequestStatusInfo().UnmatchedRoutingUsers
		}
	}

	return listOfRequestIds
}

func (v *ITService) getExecutionRequests(dbConnection *gorm.DB, userId int) []int {
	listOfWorkflows, err := database.GetObjects(dbConnection, const_util.ITServiceWorkflowEngineTable)

	listOfRequestIds := make([]int, 0)
	if err == nil {
		for _, serviceStatusObject := range *listOfWorkflows {
			serviceRequestStatus := database.ITServiceWorkflowEngine{ObjectInfo: serviceStatusObject.ObjectInfo}
			entityList := serviceRequestStatus.GetWorkFlowEngineInfo().Entities
			var executionPartEntryIndex int

			switch serviceStatusObject.Id {
			case const_util.HodWorkFlow:
				executionPartEntryIndex = const_util.HodWorkFlowExecutionEntry
			case const_util.ExecutionWorkFlow:
				executionPartEntryIndex = const_util.ExecutionWorkFlowExecutionEntry
			case const_util.SapManagerWorkFlow:
				executionPartEntryIndex = const_util.SapManagerWorkFlowExecutionEntry
			case const_util.ITManagerWorkFlow:
				executionPartEntryIndex = const_util.ITManagerWorkFlowExecutionEntry
			}

			if len(entityList) >= executionPartEntryIndex {
				routingConditionList := entityList[executionPartEntryIndex-1].RoutingCondition
				if len(routingConditionList) > 0 {
					for _, routingCondition := range routingConditionList {
						var SQLCondition string
						condition := routingCondition.Query.Condition
						notificationList := routingCondition.NotificationUserList
						listOfRules := routingCondition.Query.Rules
						lastIndexOfCategory := getLastIdOfCategoryRule(listOfRules)
						for index, rule := range routingCondition.Query.Rules {
							if rule.Field == "category" {
								if lastIndexOfCategory == index {

									SQLCondition += " object_info->>'$.categoryId' = " + strconv.Itoa(rule.Value)
								} else {

									SQLCondition += " object_info->>'$.categoryId' = " + strconv.Itoa(rule.Value) + " " + condition
								}
							}
						}
						if util.ArrayContains(notificationList, userId) {
							// now check any records for him
							listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.ITServiceRequestTable, SQLCondition)
							if err == nil {
								if len(*listOfServiceRequest) > 0 {
									for _, serviceRequestInterface := range *listOfServiceRequest {
										listOfRequestIds = append(listOfRequestIds, serviceRequestInterface.Id)
									}
								}
							}

						}

					}
				}
			}
		}

	}

	listOfRequestIds = util.RemoveDuplicateInt(listOfRequestIds)
	return listOfRequestIds
}

func getLastIdOfCategoryRule(listOfRules []database.Rules) int {
	lastIndex := 0
	for index, ruleCategory := range listOfRules {
		if ruleCategory.Field == "category" {
			lastIndex = index
		}
	}

	return lastIndex
}
