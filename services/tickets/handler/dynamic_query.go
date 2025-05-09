package handler

import (
	"cx-micro-flake/pkg/util"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

func (is *TicketsService) route(dbConnection *gorm.DB, serviceRequestId, sourceServiceStatusId int) {
	err, serviceStatusObject := Get(dbConnection, TicketsServiceRequestStatusTable, sourceServiceStatusId)
	if err == nil {
		serviceRequestStatus := TicketsServiceRequestStatus{ObjectInfo: serviceStatusObject.ObjectInfo}
		routingConditionList := serviceRequestStatus.getRequestStatusInfo().RoutingCondition
		unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers
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
				listOfServiceRequest, err := GetConditionalObjects(dbConnection, TicketsServiceRequestTable, SQLCondition)
				if err == nil {
					// should be always one, not more than one, otherwise, something is wrong.
					if len(*listOfServiceRequest) == 1 {
						TicketsServiceRequest := TicketsServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
						notificationList := routingCondition.NotificationUserList
						for _, notificationUserId := range notificationList {
							is.emailGenerator(dbConnection, SubmitForExecutionEmailTemplateType, notificationUserId, TicketsServiceMyReviewRequestComponent, TicketsServiceRequest.Id)
						}
					}
				}

			}
		} else {
			SQLConditionRequestId := " id = " + strconv.Itoa(serviceRequestId)
			listOfServiceRequest, err := GetConditionalObjects(dbConnection, TicketsServiceRequestTable, SQLConditionRequestId)
			if err == nil {
				// should be always one, not more than one, otherwise, something is wrong.
				if len(*listOfServiceRequest) == 1 {
					TicketsServiceRequest := TicketsServiceRequest{ObjectInfo: (*listOfServiceRequest)[0].ObjectInfo}
					for _, notificationUserId := range unmatchedRoutingUsers {
						is.emailGenerator(dbConnection, SubmitForExecutionEmailTemplateType, notificationUserId, TicketsServiceMyReviewRequestComponent, TicketsServiceRequest.Id)
					}
				}
			}

			// lets handle the unmatached routing //unmatchedRoutingUsers

		}
	}
}

func (ts *TicketsService) getTargetWorkflowGroup(dbConnection *gorm.DB, requestId int, workflowId int) []int {
	err, serviceStatusObject := Get(dbConnection, TicketsServiceRequestStatusTable, workflowId)
	var listOfRequestIds []int
	if err == nil {
		serviceRequestStatus := TicketsServiceRequestStatus{ObjectInfo: serviceStatusObject.ObjectInfo}
		routingConditionList := serviceRequestStatus.getRequestStatusInfo().RoutingCondition
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

				SQLCondition += ") and  id = " + strconv.Itoa(requestId)
				fmt.Println("SQLCondition: ", SQLCondition)
				listOfRequests, err := GetConditionalObjects(dbConnection, TicketsServiceRequestTable, SQLCondition)
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
	}
	return listOfRequestIds
}
func (ts *TicketsService) getExecutionRequests(dbConnection *gorm.DB, userId int) []int {
	err, serviceStatusObject := Get(dbConnection, TicketsServiceRequestStatusTable, WorkFlowExecutionParty)
	var listOfRequestIds []int
	if err == nil {
		serviceRequestStatus := TicketsServiceRequestStatus{ObjectInfo: serviceStatusObject.ObjectInfo}
		routingConditionList := serviceRequestStatus.getRequestStatusInfo().RoutingCondition
		//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers
		if len(routingConditionList) > 0 {
			for _, routingCondition := range routingConditionList {
				var SQLCondition string
				condition := routingCondition.Query.Condition
				notificationList := routingCondition.NotificationUserList

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
				if util.ArrayContains(notificationList, userId) {
					// now check any records for him
					listOfServiceRequest, err := GetConditionalObjects(dbConnection, TicketsServiceRequestTable, SQLCondition)
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
	return listOfRequestIds
}
