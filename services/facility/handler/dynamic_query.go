package handler

import (
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

func (v *FacilityService) route(dbConnection *gorm.DB, serviceRequestId, sourceServiceStatusId int) {
	err, serviceStatusObject := database.Get(dbConnection, const_util.FacilityServiceRequestStatusTable, sourceServiceStatusId)
	if err == nil {
		requestStatusInfo := database.GetRequestStatusInfo(serviceStatusObject.ObjectInfo)
		routingConditionList := requestStatusInfo.RoutingCondition
		unmatchedRoutingUsers := requestStatusInfo.UnmatchedRoutingUsers
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
				listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLCondition)
				if err == nil {
					// should be always one, not more than one, otherwise, something is wrong.
					if len(*listOfServiceRequest) == 1 {
						FacilityServiceRequest := database.FacilityServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
						notificationList := routingCondition.NotificationUserList
						for _, notificationUserId := range notificationList {
							v.NotificationHandler.EmailGenerator(dbConnection, const_util.ExecutionPartyAssignTemplateType, notificationUserId, const_util.FacilityServiceMyReviewRequestComponent, FacilityServiceRequest.Id)
							v.BaseService.Logger.Info("routing request email", zap.Any("user_id", notificationUserId))
						}
					}
				}

			}
		} else {
			SQLConditionRequestId := " id = " + strconv.Itoa(serviceRequestId)
			v.BaseService.Logger.Info("generated condition", zap.Any("condition", SQLConditionRequestId))
			listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLConditionRequestId)
			if err == nil {
				// should be always one, not more than one, otherwise, something is wrong.
				if len(*listOfServiceRequest) == 1 {
					FacilityServiceRequest := database.FacilityServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
					for _, notificationUserId := range unmatchedRoutingUsers {
						v.NotificationHandler.EmailGenerator(dbConnection, const_util.ExecutionPartyAssignTemplateType, notificationUserId, const_util.FacilityServiceMyReviewRequestComponent, FacilityServiceRequest.Id)
						v.BaseService.Logger.Info("un matched routing request email", zap.Any("user_id", notificationUserId))
					}
				}
			}

			// lets handle the unmatached routing //unmatchedRoutingUsers

		}
	}
}

func (v *FacilityService) getTargetWorkflowGroup(dbConnection *gorm.DB, requestId int, workflowId int) []int {
	err, serviceStatusObject := database.Get(dbConnection, const_util.FacilityServiceRequestStatusTable, workflowId)
	var listOfRequestIds []int
	if err == nil {
		requestStatusInfo := database.GetRequestStatusInfo(serviceStatusObject.ObjectInfo)
		routingConditionList := requestStatusInfo.RoutingCondition
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
				listOfRequests, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLCondition)
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
			listOfRequestIds = requestStatusInfo.UnmatchedRoutingUsers
		}
	}

	return listOfRequestIds
}

func (v *FacilityService) getExecutionRequests(dbConnection *gorm.DB, userId int) []int {
	listOfWorkflows, err := database.GetObjects(dbConnection, const_util.FacilityServiceWorkflowEngineTable)

	listOfRequestIds := make([]int, 0)
	if err == nil {
		for _, serviceStatusObject := range *listOfWorkflows {
			workFlowEngineInfo := database.GetWorkFlowEngineInfo(serviceStatusObject.ObjectInfo)
			entityList := workFlowEngineInfo.Entities
			//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers

			var executionPartEntryIndex int

			switch serviceStatusObject.Id {
			case const_util.HodWorkFlow:
				executionPartEntryIndex = const_util.HodWorkFlowExecutionEntry
			case const_util.ExecutionWorkFlow:
				executionPartEntryIndex = const_util.ExecutionWorkFlowExecutionEntry
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
							listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLCondition)

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

func (v *FacilityService) getEHSManagerRequests(dbConnection *gorm.DB, userId int) []int {
	listOfWorkflows, err := database.GetObjects(dbConnection, const_util.FacilityServiceWorkflowEngineTable)

	listOfRequestIds := make([]int, 0)
	if err == nil {
		for _, serviceStatusObject := range *listOfWorkflows {
			workFlowEngineInfo := database.GetWorkFlowEngineInfo(serviceStatusObject.ObjectInfo)
			entityList := workFlowEngineInfo.Entities
			//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers

			var ehsPartEntryIndex int

			switch serviceStatusObject.Id {
			case const_util.ExecutionWorkFlow:
				ehsPartEntryIndex = const_util.HodWorkFlowExecutionEntry
			}

			if ehsPartEntryIndex > 0 {
				if len(entityList) >= ehsPartEntryIndex {
					routingConditionList := entityList[ehsPartEntryIndex-1].RoutingCondition
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
								listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLCondition)

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

	}

	listOfRequestIds = util.RemoveDuplicateInt(listOfRequestIds)
	return listOfRequestIds
}

func (v *FacilityService) getFacilityManagerRequests(dbConnection *gorm.DB, userId int) []int {
	listOfWorkflows, err := database.GetObjects(dbConnection, const_util.FacilityServiceWorkflowEngineTable)

	listOfRequestIds := make([]int, 0)
	if err == nil {
		for _, serviceStatusObject := range *listOfWorkflows {
			workFlowEngineInfo := database.GetWorkFlowEngineInfo(serviceStatusObject.ObjectInfo)
			entityList := workFlowEngineInfo.Entities
			//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers

			var ehsPartEntryIndex int

			switch serviceStatusObject.Id {
			case const_util.ExecutionWorkFlow:
				ehsPartEntryIndex = const_util.FacilityManager
			}

			if ehsPartEntryIndex > 0 {
				if len(entityList) >= ehsPartEntryIndex {
					routingConditionList := entityList[ehsPartEntryIndex-1].RoutingCondition
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
								listOfServiceRequest, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceRequestTable, SQLCondition)

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
