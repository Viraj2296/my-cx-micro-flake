package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}

type Assignment struct {
	Event    int `json:"int"`
	Id       int `json:"id"`
	Resource int `json:"resource"`
}

type OrderScheduledEventUpdateRequest struct {
	EndDate    string      `json:"endDate"`
	StartDate  string      `json:"startDate"`
	Assignment *Assignment `json:"assignment"`
}

func (v *MaintenanceService) updateSchedulerEvent(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	updateScheduleEventRequest := OrderScheduledEventUpdateRequest{}
	if err := ctx.ShouldBindBodyWith(&updateScheduleEventRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	userId := common.GetUserId(ctx)
	currentTime := time.Now().Unix()

	// convert the stat time and end time
	//authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	//startTime := authService.ConvertToUserTimezoneToISO(userId, updateScheduleEventRequest.StartDate)
	//endTime := authService.ConvertToUserTimezoneToISO(userId, updateScheduleEventRequest.EndDate)
	startTime := util.ConvertSingaporeTimeToUTC(updateScheduleEventRequest.StartDate)
	endTime := util.ConvertSingaporeTimeToUTC(updateScheduleEventRequest.EndDate)

	v.BaseService.Logger.Info("date validation", zap.String("start_date", updateScheduleEventRequest.StartDate), zap.String("end_date", updateScheduleEventRequest.EndDate),
		zap.String("start_time", startTime), zap.String("end_time", endTime), zap.Any("current_time", currentTime))

	parsedStartTime, _ := time.Parse(ISOTimeLayout, startTime)
	parsedEndTime, _ := time.Parse(ISOTimeLayout, endTime)
	if (currentTime > parsedStartTime.Unix()) && (currentTime > parsedEndTime.Unix()) {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSourceError), InvalidScheduleStatus, "You are trying to modify the event which is already passed, This modification is rejected")
		return
	}

	if componentName == "maintenance_work_order" {
		maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()

		if updateScheduleEventRequest.Assignment != nil {
			if updateScheduleEventRequest.Assignment.Resource != maintenanceWorkOrder.getWorkOrderInfo().AssetId {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSourceError), InvalidScheduleStatus, "Changing resource which already assigned is not permitted, consider moving along with the resources")
				return
			}
		}

		workOrderInfo.WorkOrderScheduledStartDate = startTime
		workOrderInfo.WorkOrderScheduledEndDate = endTime

		err = Update(dbConnection, targetTable, recordId, workOrderInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

	} else if componentName == "maintenance_work_order_task" {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: eventObject.ObjectInfo}
		workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()
		workOrderTaskInfo.TaskDate = startTime
		workOrderTaskInfo.EstimatedTaskEndDate = endTime

		err = Update(dbConnection, targetTable, recordId, workOrderTaskInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})
}

func getWorkOrderStatus(dbConnection *gorm.DB, workOrderId int, targetTable string) (error, int) {
	err, workOrderGeneralObject := Get(dbConnection, targetTable, workOrderId)
	if err == nil {
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(workOrderGeneralObject.ObjectInfo, &workOrderInfo)
		return nil, util.InterfaceToInt(workOrderInfo["workOrderStatus"])
	}
	return err, -1
}

func (v *MaintenanceService) checkInTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var workOrderTable string
	if targetTable == MaintenanceWorkOrderTaskTable {
		workOrderTable = MaintenanceWorkOrderTable
	} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTaskTable {
		workOrderTable = MouldMaintenanceCorrectiveWorkOrderTable
	} else if targetTable == MouldMaintenancePreventiveWorkOrderTaskTable {
		workOrderTable = MouldMaintenancePreventiveWorkOrderTable
	} else {
		workOrderTable = MaintenanceCorrectiveWorkOrderComponent
	}

	err, taskGeneralObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	workOrderTaskInfo := make(map[string]interface{})
	json.Unmarshal(taskGeneralObject.ObjectInfo, &workOrderTaskInfo)

	// before check in, can we allowed to start the task?
	workOrderId := util.InterfaceToInt(workOrderTaskInfo["workOrderId"])
	err, workOrderStatus := getWorkOrderStatus(dbConnection, workOrderId, workOrderTable)
	fmt.Println("workOrderTable : ", workOrderTable, "workOrderStatus: ", workOrderStatus)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "Internal system error, requested work order not available, please report the error code to system admin",
			})
		return
	}
	if workOrderStatus == WorkOrderScheduled || workOrderStatus == WorkOrderInProgress {
		if util.InterfaceToInt(workOrderTaskInfo["taskStatus"]) == WorkOrderTaskInProgress {
			// already checked in , then why?
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "You have already checked in to this task, please proceed to complete once it is actually done",
				})
			return
		}
		//workOrderTaskInfo.EstimatedTaskEndDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		workOrderTaskInfo["taskStatus"] = WorkOrderTaskInProgress

		workOrderId = util.InterfaceToInt(workOrderTaskInfo["workOrderId"])

		_, workOrderObject := Get(dbConnection, workOrderTable, workOrderId)
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(workOrderObject.ObjectInfo, &workOrderInfo)

		if util.InterfaceToInt(workOrderInfo["workOrderStatus"]) != WorkOrderInProgress {
			workOrderInfo["workOrderStatus"] = WorkOrderInProgress

			workOrderInfo["canComplete"] = false
			workOrderInfo["canUpdate"] = false
			workOrderInfo["canRelease"] = false
			workOrderInfo["canUnRelease"] = false

			updateObject := make(map[string]interface{})
			updateObject["object_info"], _ = json.Marshal(workOrderInfo)

			err = Update(dbConnection, workOrderTable, workOrderId, updateObject)
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
			// now put the machine on maintenace
			if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveMachineToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
				mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
				if util.InterfaceToString(workOrderInfo["workOrderType"]) == WorkOrderTypePreventive {
					mouldInterface.PutToMaintenanceMode(projectId, util.InterfaceToInt(workOrderInfo["assetId"]), userId)
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveAssemblyMachineToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveToolingMachineToMaintenance(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			}

		}

		workOrderTaskInfo["canCheckIn"] = false
		workOrderTaskInfo["canCheckOut"] = true
		workOrderTaskInfo["checkInDate"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

		updateRequest := make(map[string]interface{})
		updateRequest["object_info"], _ = json.Marshal(workOrderTaskInfo)
		err = Update(dbConnection, targetTable, recordId, updateRequest)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
			return
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Message: "You have successfully checked in",
			Code:    0,
		})

	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "Please wait until work order is get scheduled and finalised, Check-In task is not allowed at this moment",
			})
		return
	}

}

func (v *MaintenanceService) completeTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	workOrderTaskInfo := make(map[string]interface{})
	json.Unmarshal(taskGeneralObject.ObjectInfo, &workOrderTaskInfo)
	workOrderTaskInfo["taskStatus"] = WorkOrderDone

	workOrderTaskInfo["checkOutDate"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

	workOrderTaskInfo["canApprove"] = true
	workOrderTaskInfo["canReject"] = true
	workOrderTaskInfo["canCheckOut"] = false
	workOrderTaskInfo["canComplete"] = true

	workOrderTaskInfo["lastUpdatedBy"] = userId

	requestData := make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&requestData, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if componentName == MaintenanceWorkOrderMyCorrectiveTaskComponent || componentName == MyMouldMaintenanceCorrectiveWorkOrderTaskComponent {

		workOrderTaskInfo["faultCode"] = requestData["faultCode"]
		workOrderTaskInfo["cost"] = requestData["cost"]
		workOrderTaskInfo["remarks"] = requestData["remarks"]
		workOrderTaskInfo["labourCost"] = requestData["labourCost"]
		workOrderTaskInfo["canUpdate"] = false

	}

	if componentName == MaintenanceWorkOrderMyTaskComponent || componentName == MyMouldMaintenancePreventiveWorkOrderTaskComponent {
		workOrderTaskInfo["remarks"] = requestData["remarks"]
	}

	updateObject := make(map[string]interface{})
	updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

	err = Update(dbConnection, targetTable, recordId, updateObject)
	fmt.Println("Err: ", err)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	workOrderId := util.InterfaceToInt(workOrderTaskInfo["workOrderId"])
	// when every time task is getting finished, we need to see all the task is completed
	// condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\"))  = \"" + strconv.Itoa(workOrderId) + "\""
	condition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(workOrderId) + "'"
	listOfWorkOrderTasks, err := GetConditionalObjects(dbConnection, targetTable, condition)
	var numberOfCompletedTasks int
	if err == nil {
		for _, workOrderTaskObjectInterface := range *listOfWorkOrderTasks {
			workOrderTask := make(map[string]interface{})
			json.Unmarshal(workOrderTaskObjectInterface.ObjectInfo, &workOrderTask)
			if util.InterfaceToInt(workOrderTask["taskStatus"]) == WorkOrderTaskDone {
				numberOfCompletedTasks = numberOfCompletedTasks + 1
			}
		}
	}
	var workOrderTable string
	if targetTable == MaintenanceWorkOrderCorrectiveTaskComponent {
		workOrderTable = MaintenanceCorrectiveWorkOrderComponent
	} else if targetTable == MaintenanceWorkOrderTaskComponent {
		workOrderTable = MaintenanceWorkOrderTable
	} else if targetTable == MouldMaintenancePreventiveWorkOrderTaskComponent {
		workOrderTable = MouldMaintenancePreventiveWorkOrderTable
	} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTaskComponent {
		workOrderTable = MouldMaintenanceCorrectiveWorkOrderTable
	}
	if numberOfCompletedTasks == len(*listOfWorkOrderTasks) {
		// yes we are done
		_, workOrderObject := Get(dbConnection, workOrderTable, workOrderId)
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(workOrderObject.ObjectInfo, &workOrderInfo)

		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = true
		workOrderInfo["canUpdate"] = false
		workOrderInfo["workOrderActualEndDate"] = util.GetCurrentTime(ISOTimeLayout)

		// calculating the downtime hours of machine /mould corrective maintenance work orders
		if workOrderTable == MaintenanceCorrectiveWorkOrderComponent || workOrderTable == MouldMaintenanceCorrectiveWorkOrderTable {
			// extract actual and scheduled date strings
			actualEndDateStr, _ := workOrderInfo["workOrderActualEndDate"].(string)
			scheduledStartDateStr, _ := workOrderInfo["workOrderScheduledStartDate"].(string)

			// validate the extracted dates
			if actualEndDateStr == "" || scheduledStartDateStr == "" {
				v.BaseService.Logger.Warn("missing date fields for downtime calculation")
				return
			}

			// parse date strings using RFC3339 format
			layout := time.RFC3339
			actualEndDate, err1 := time.Parse(layout, actualEndDateStr)
			scheduledStartDate, err2 := time.Parse(layout, scheduledStartDateStr)

			if err1 == nil && err2 == nil {
				// Calculate duration in hours
				duration := actualEndDate.Sub(scheduledStartDate).Hours()

				// prevent invalid negative durations
				if duration < 0 {
					v.BaseService.Logger.Warn("actual end date is before scheduled start date",
						zap.Float64("duration", duration))
					return
				}

				// Round to 2 decimal places and store it
				workOrderInfo["downTimeHours"] = math.Round(duration*100) / 100
			} else {
				v.BaseService.Logger.Error("error calculating downtime hours", zap.String("actualEndDateError", err1.Error()), zap.String("scheduledStartDateError", err2.Error()))
				return
			}
		}
		updateOrderObject := make(map[string]interface{})
		updateOrderObject["object_info"], _ = json.Marshal(workOrderInfo)
		err = Update(dbConnection, workOrderTable, workOrderId, updateOrderObject)

		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			_ = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		}
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			_ = mouldService.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
		}

		var generalWorkOrderStatus component.GeneralObject
		notificationId := CorrectiveWorkOrderCompletionNotification
		emailMainTable := MaintenanceCorrectiveWorkOrderComponent
		if targetTable == MaintenanceWorkOrderCorrectiveTaskComponent {
			_, generalWorkOrderStatus = Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
		} else if targetTable == MaintenanceWorkOrderTaskComponent {
			notificationId = WorkOrderCompletionNotification
			emailMainTable = MaintenanceWorkOrderComponent
			_, generalWorkOrderStatus = Get(dbConnection, MaintenancePreventiveWorkOrderStatusTable, WorkOrderDone)
		} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTaskComponent {
			_, generalWorkOrderStatus = Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
			notificationId = MouldCorrectiveWorkOrderCompletionNotification
			emailMainTable = MouldMaintenanceCorrectiveWorkOrderComponent
		} else if targetTable == MouldMaintenancePreventiveWorkOrderTaskComponent {
			notificationId = MouldPreventiveWorkOrderCompletionNotification
			emailMainTable = MouldMaintenancePreventiveWorkOrderComponent
			_, generalWorkOrderStatus = Get(dbConnection, MaintenancePreventiveWorkOrderStatusTable, WorkOrderDone)
		}

		// Send notification to configured user in creation status
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, notificationId, workflowUser, emailMainTable, workOrderObject.Id)
				fmt.Println("Errror:", err)
			}
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully completed",
		Code:    0,
	})

}

func (v *MaintenanceService) approveTask(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)

	var workOrderTable string
	if targetTable == MaintenanceWorkOrderTaskTable {
		workOrderTable = MaintenanceWorkOrderTable
	} else if targetTable == MouldMaintenanceCorrectiveWorkOrderTaskTable {
		workOrderTable = MouldMaintenanceCorrectiveWorkOrderTable
	} else if targetTable == MouldMaintenancePreventiveWorkOrderTaskTable {
		workOrderTable = MouldMaintenancePreventiveWorkOrderTable
	} else {
		workOrderTable = MaintenanceCorrectiveWorkOrderComponent
	}
	workOrderTaskInfo := make(map[string]interface{})
	json.Unmarshal(taskGeneralObject.ObjectInfo, &workOrderTaskInfo)
	workOrderId := util.InterfaceToInt(workOrderTaskInfo["workOrderId"])

	_, workOrderObject := Get(dbConnection, workOrderTable, workOrderId)
	workOrderInfo := make(map[string]interface{})
	json.Unmarshal(workOrderObject.ObjectInfo, &workOrderInfo)

	workOrderTaskInfo["taskStatus"] = WorkOrderTaskApproved

	workOrderTaskInfo["canApprove"] = false
	workOrderTaskInfo["canReject"] = false
	workOrderTaskInfo["lastUpdatedBy"] = userId
	updateObject := make(map[string]interface{})
	updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

	err = Update(dbConnection, targetTable, recordId, updateObject)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	// when every time task is getting finished, we need to see all the task is completed
	condition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\"))  = \"" + strconv.Itoa(workOrderId) + "\""
	listOfWorkOrderTasks, err := GetConditionalObjects(dbConnection, targetTable, condition)
	var numberOfApprovedTasks int
	if err == nil {
		for _, workOrderTaskObjectInterface := range *listOfWorkOrderTasks {
			workOrderTask := make(map[string]interface{})
			json.Unmarshal(workOrderTaskObjectInterface.ObjectInfo, &workOrderTask)
			if util.InterfaceToInt(workOrderTask["taskStatus"]) == WorkOrderTaskApproved {
				numberOfApprovedTasks = numberOfApprovedTasks + 1
			}
		}
	}
	if numberOfApprovedTasks == len(*listOfWorkOrderTasks) {
		// yes we are done

		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = true

		updateOrderObject := make(map[string]interface{})
		updateOrderObject["object_info"], _ = json.Marshal(workOrderInfo)
		err = Update(dbConnection, workOrderTable, workOrderId, updateOrderObject)

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "You have successfully approved",
		Code:    0,
	})

}

func (v *MaintenanceService) rejectTask(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	var rejectRemarksFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&rejectRemarksFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	rejectRemarks := util.InterfaceToString(rejectRemarksFields["remark"])

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, taskGeneralObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	workOrderTaskInfo := make(map[string]interface{})
	json.Unmarshal(taskGeneralObject.ObjectInfo, &workOrderTaskInfo)

	workOrderTaskInfo["taskStatus"] = WorkOrderTaskReDo
	workOrderTaskInfo["canApprove"] = false
	workOrderTaskInfo["canReject"] = false
	workOrderTaskInfo["canCheckIn"] = true
	workOrderTaskInfo["canCheckOut"] = false
	existingActionRemarks := ConvertRemarkArray(workOrderTaskInfo["actionRemarks"])
	existingActionRemarks = append(existingActionRemarks, ActionRemarks{
		ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		Status:        "TASK  REJECTED",
		UserId:        userId,
		Remarks:       rejectRemarks,
		ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderTaskInfo["createdAt"])),
	})
	workOrderTaskInfo["actionRemarks"] = existingActionRemarks
	workOrderTaskInfo["lastUpdatedBy"] = userId

	updateObject := make(map[string]interface{})
	updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

	err = Update(dbConnection, targetTable, recordId, updateObject)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating maintenance work task information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Task rejection is successfully updated",
		Code:    0,
	})

}

func (v *MaintenanceService) kanbanMoveTask(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)

	kanbanTaskStatus := KanbanTaskStatusRequest{}
	if err := ctx.ShouldBindBodyWith(&kanbanTaskStatus, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if componentName == MaintenanceWorkOrderTaskComponent {

		// before update, check the work order is completed , then don.t let the user to move
		if err, isTaskIsDone := v.isWorkOrderCompletedFromTaskId(dbConnection, recordId); err == nil {
			if isTaskIsDone {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, this work order is already been done, further modification is not allowed",
					})
				return
			}
		}

		// we need to check whether we have any tasks created before release
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: eventObject.ObjectInfo}
		workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()

		if workOrderTaskInfo.TaskStatus == WorkOrderTaskToDo {
			// it is illegal to move the task when the work order is not schedueld

			if err, isWorkOrderScheduled := v.isWorkOrderScheduledFromTaskId(dbConnection, recordId); err == nil {
				if !isWorkOrderScheduled {
					response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
						&response.DetailedError{
							Header:      getError(common.InvalidObjectStatusError).Error(),
							Description: "Please wait until the work order get scheduled. The scheduler will schedule and inform the date to start the task!!",
						})
					return
				}
			}
			if kanbanTaskStatus.TaskStatus > 1 {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, once the task is done only, it can be moved to REDO stage if necessary.",
					})
				return
			}

			if kanbanTaskStatus.TaskStatus == WorkOrderTaskReDo {
				response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
					&response.DetailedError{
						Header:      getError(common.InvalidObjectStatusError).Error(),
						Description: "Sorry, once the task is done only, it can be moved to REDO stage if necessary.",
					})
				return
			}

		}

		// existing task status id
		workOrderTaskInfo.TaskStatus = kanbanTaskStatus.TaskStatus

		err = Update(dbConnection, targetTable, recordId, workOrderTaskInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})

}
func (v *MaintenanceService) releaseWorkOrder(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}

	userId := common.GetUserId(ctx)
	if componentName == "maintenance_work_order" {
		// we need to check whether we have any tasks created before release
		// taskCondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = " + strconv.Itoa(recordId)
		taskCondition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, taskCondition)
		if len(*listOfTasks) == 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "There is no tasks defined to release the work order, please define corresponding tasks and assigned name",
				})
			return
		}

		for _, workOrderTask := range *listOfTasks {
			maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
			workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()
			workOrderTaskInfo.IsOrderReleased = true
			err = Update(dbConnection, MaintenanceWorkOrderTaskTable, workOrderTask.Id, workOrderTaskInfo.DatabaseSerialize(userId))

			if err != nil {
				v.BaseService.Logger.Error("error updating work order task record", zap.String("error", err.Error()))
			}
		}
		maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
		workOrderInfo.WorkOrderStatus = WorkOrderScheduled
		workOrderInfo.CanUpdate = false
		workOrderInfo.CanRelease = false
		workOrderInfo.CanUnRelease = true

		err = Update(dbConnection, targetTable, recordId, workOrderInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == MouldMaintenanceCorrectiveWorkOrderComponent {
		// we need to check whether we have any tasks created before release
		taskCondition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, taskCondition)
		if len(*listOfTasks) == 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "There is no tasks defined to release the work order, please define corresponding tasks and assigned name",
				})
			return
		}

		for _, workOrderTask := range *listOfTasks {

			workOrderTaskInfo := make(map[string]interface{})
			json.Unmarshal(workOrderTask.ObjectInfo, &workOrderTaskInfo)
			workOrderTaskInfo["isOrderReleased"] = true

			workOrderTaskInfo["lastUpdatedBy"] = userId
			updateObject := make(map[string]interface{})
			updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)
			err = Update(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, workOrderTask.Id, updateObject)

			if err != nil {
				v.BaseService.Logger.Error("error updating work order task record", zap.String("error", err.Error()))
			}
		}

		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)

		workOrderInfo["workOrderStatus"] = WorkOrderScheduled
		workOrderInfo["canUpdate"] = false
		workOrderInfo["canRelease"] = false
		workOrderInfo["canUnRelease"] = true
		workOrderInfo["canComplete"] = true

		workOrderInfo["lastUpdatedBy"] = userId
		updateWorkOrderObject := make(map[string]interface{})
		updateWorkOrderObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateWorkOrderObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == "maintenance_work_order_task" {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: eventObject.ObjectInfo}
		workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()
		workOrderTaskInfo.TaskStatus = 2

		err = Update(dbConnection, targetTable, recordId, workOrderTaskInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == "maintenance_corrective_work_order" {
		// we need to check whether we have any tasks created before release
		// taskCondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = " + strconv.Itoa(recordId)
		taskCondition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, taskCondition)
		if len(*listOfTasks) == 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "There is no tasks defined to release the work order, please define corresponding tasks and assigned name",
				})
			return
		}

		for _, workOrderTask := range *listOfTasks {

			workOrderTaskInfo := make(map[string]interface{})
			json.Unmarshal(workOrderTask.ObjectInfo, &workOrderTaskInfo)
			workOrderTaskInfo["isOrderReleased"] = true

			workOrderTaskInfo["lastUpdatedBy"] = userId
			updateObject := make(map[string]interface{})
			updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)
			err = Update(dbConnection, MaintenanceWorkOrderCorrectiveTaskComponent, workOrderTask.Id, updateObject)

			if err != nil {
				v.BaseService.Logger.Error("error updating work order task record", zap.String("error", err.Error()))
			}
		}

		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = WorkOrderScheduled
		workOrderInfo["canUpdate"] = false
		workOrderInfo["canRelease"] = false
		workOrderInfo["canUnRelease"] = true
		workOrderInfo["canComplete"] = true

		workOrderInfo["lastUpdatedBy"] = userId
		updateWorkOrderObject := make(map[string]interface{})
		updateWorkOrderObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateWorkOrderObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == MouldMaintenancePreventiveWorkOrderComponent {
		// we need to check whether we have any tasks created before release
		taskCondition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, taskCondition)
		if len(*listOfTasks) == 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      getError(common.InvalidObjectStatusError).Error(),
					Description: "There is no tasks defined to release the work order, please define corresponding tasks and assigned name",
				})
			return
		}

		for _, workOrderTask := range *listOfTasks {
			maintenanceWorkOrderTask := MouldMaintenancePreventiveWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
			workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()
			workOrderTaskInfo.IsOrderReleased = true
			err = Update(dbConnection, MaintenanceWorkOrderTaskTable, workOrderTask.Id, workOrderTaskInfo.DatabaseSerialize(userId))

			if err != nil {
				v.BaseService.Logger.Error("error updating work order task record", zap.String("error", err.Error()))
			}
		}
		maintenanceWorkOrder := MouldMaintenancePreventiveWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
		workOrderInfo.WorkOrderStatus = WorkOrderScheduled
		workOrderInfo.CanUpdate = false
		workOrderInfo.CanRelease = false
		workOrderInfo.CanUnRelease = true

		err = Update(dbConnection, targetTable, recordId, workOrderInfo.DatabaseSerialize(userId))
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})

}

func (v *MaintenanceService) holdOrder(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	if componentName == "maintenance_work_order" || componentName == "maintenance_corrective_work_order" || componentName == MouldMaintenanceCorrectiveWorkOrderComponent || componentName == MouldMaintenancePreventiveWorkOrderComponent {
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = 1
		workOrderInfo["lastUpdatedBy"] = userId
		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	} else if componentName == "maintenance_work_order_task" || componentName == "maintenance_work_order_corrective_task" {
		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderTaskInfo)
		workOrderTaskInfo["taskStatus"] = 1
		workOrderTaskInfo["lastUpdatedBy"] = userId
		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}
	}
	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})
}

// completeOrder ShowAccount godoc
// @Summary schedule orders
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/action/schedule [get]
func (v *MaintenanceService) forceCompleteOrder(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	
	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	if componentName == "maintenance_work_order" || componentName == "maintenance_corrective_work_order" {
		//maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = false
		workOrderInfo["canComplete"] = false
		workOrderInfo["canForceStop"] = false
		workOrderInfo["workOrderActualEndDate"] = util.GetCurrentTime(ISOTimeLayout)

		existingActionRemarks := util.AppendToObjectArray(workOrderInfo["actionRemarks"], ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "FORCEFULLY COMPLETED",
			UserId:        userId,
			Remarks:       "FORCEFULLY COMPLETED",
			ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		})
		//existingActionRemarks = append(existingActionRemarks, ActionRemarks{
		//	ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		//	Status:        "FORCEFULLY COMPLETED",
		//	UserId:        userId,
		//	Remarks:       "FORCEFULLY COMPLETED",
		//	ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		//})

		workOrderInfo["actionRemarks"] = existingActionRemarks

		// calculating the downtime hours of machine corrective maintenance work orders
		if componentName == "maintenance_corrective_work_order" {
			// calculating the downtime hours of machine /mould corrective maintenance work orders
			// extract actual and scheduled date strings
			actualEndDateStr, _ := workOrderInfo["workOrderActualEndDate"].(string)
			scheduledStartDateStr, _ := workOrderInfo["workOrderScheduledStartDate"].(string)

			// validate the extracted dates
			if actualEndDateStr == "" || scheduledStartDateStr == "" {
				v.BaseService.Logger.Warn("missing date fields for downtime calculation")
				return
			}

			// parse date strings using RFC3339 format
			layout := time.RFC3339
			actualEndDate, err1 := time.Parse(layout, actualEndDateStr)
			scheduledStartDate, err2 := time.Parse(layout, scheduledStartDateStr)

			if err1 == nil && err2 == nil {
				// Calculate duration in hours
				duration := actualEndDate.Sub(scheduledStartDate).Hours()

				// prevent invalid negative durations
				if duration < 0 {
					v.BaseService.Logger.Warn("actual end date is before scheduled start date",
						zap.Float64("duration", duration))
					return
				}

				// Round to 2 decimal places and store it
				workOrderInfo["downTimeHours"] = math.Round(duration*100) / 100
			} else {
				v.BaseService.Logger.Error("error calculating downtime hours", zap.String("actualEndDateError", err1.Error()), zap.String("scheduledStartDateError", err2.Error()))
				return
			}
		}

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTaskTable string
		if targetTable == MaintenanceWorkOrderTable {
			workOrderTaskTable = MaintenanceWorkOrderTaskTable
		} else {
			workOrderTaskTable = MaintenanceWorkOrderCorrectiveTaskComponent
		}

		// complete work order will complete all the task assigned
		// condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(recordId) + "\""
		condition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, workOrderTaskTable, condition)
		if err == nil {
			for _, workOrderTaskObject := range *listOfTasks {
				workOrderTask := make(map[string]interface{})
				json.Unmarshal(workOrderTaskObject.ObjectInfo, &workOrderTask)
				if util.InterfaceToInt(workOrderTask["taskStatus"]) != WorkOrderTaskDone {
					workOrderTask["taskStatus"] = WorkOrderTaskDone

					updateObject = make(map[string]interface{})
					updateObject["object_info"], _ = json.Marshal(workOrderTask)

					Update(dbConnection, workOrderTaskTable, workOrderTaskObject.Id, updateObject)
				}
			}
		}

		// now put the machine on Active
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
			_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		}

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, CorrectiveWorkOrderCompletionNotification, workflowUser, MaintenanceCorrectiveWorkOrderComponent, eventObject.Id)
				fmt.Println("Errror:", err)
			}
		}

	} else if componentName == "maintenance_work_order_task" || componentName == "maintenance_work_order_corrective_task" {
		//maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: eventObject.ObjectInfo}
		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderTaskInfo)
		workOrderTaskInfo["taskStatus"] = WorkOrderTaskDone
		workOrderTaskInfo["canComplete"] = false

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTable string
		if componentName == MaintenanceWorkOrderTaskTable {
			workOrderTable = MaintenanceWorkOrderTable
		} else {
			workOrderTable = MaintenanceWorkOrderCorrectiveTaskComponent
		}

		if checkAllTaskAreDone(util.InterfaceToInt(workOrderTaskInfo["workOrderId"]), dbConnection, componentName) {
			_, workOrderObject := Get(dbConnection, workOrderTable, util.InterfaceToInt(workOrderTaskInfo["workOrderId"]))
			//maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrderObject.ObjectInfo}
			workOrderInfo := make(map[string]interface{})
			json.Unmarshal(workOrderObject.ObjectInfo, &workOrderInfo)
			workOrderInfo["canComplete"] = true

			workOrderInfo["workOrderStatus"] = WorkOrderDone

			updateObject = make(map[string]interface{})
			updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

			err = Update(dbConnection, workOrderTable, util.InterfaceToInt(workOrderTaskInfo["workOrderId"]), updateObject)

			if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
				mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
				err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			}
		}
	} else if componentName == MouldMaintenancePreventiveWorkOrderComponent || componentName == MouldMaintenanceCorrectiveWorkOrderComponent {
		//maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = false
		workOrderInfo["canComplete"] = false
		workOrderInfo["canForceStop"] = false
		workOrderInfo["workOrderActualEndDate"] = util.GetCurrentTime(ISOTimeLayout)

		existingActionRemarks := util.AppendToObjectArray(workOrderInfo["actionRemarks"], ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "FORCEFULLY COMPLETED",
			UserId:        userId,
			Remarks:       "FORCEFULLY COMPLETED",
			ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		})

		workOrderInfo["actionRemarks"] = existingActionRemarks

		// calculating the downtime hours of mould corrective maintenance work orders
		if componentName == MouldMaintenanceCorrectiveWorkOrderComponent {
			// extract actual and scheduled date strings
			actualEndDateStr, _ := workOrderInfo["workOrderActualEndDate"].(string)
			scheduledStartDateStr, _ := workOrderInfo["workOrderScheduledStartDate"].(string)

			// validate the extracted dates
			if actualEndDateStr == "" || scheduledStartDateStr == "" {
				v.BaseService.Logger.Warn("missing date fields for downtime calculation")
				return
			}

			// parse date strings using RFC3339 format
			layout := time.RFC3339
			actualEndDate, err1 := time.Parse(layout, actualEndDateStr)
			scheduledStartDate, err2 := time.Parse(layout, scheduledStartDateStr)

			if err1 == nil && err2 == nil {
				// Calculate duration in hours
				duration := actualEndDate.Sub(scheduledStartDate).Hours()

				// prevent invalid negative durations
				if duration < 0 {
					v.BaseService.Logger.Warn("actual end date is before scheduled start date",
						zap.Float64("duration", duration))
					return
				}

				// Round to 2 decimal places and store it
				workOrderInfo["downTimeHours"] = math.Round(duration*100) / 100
			} else {
				v.BaseService.Logger.Error("error calculating downtime hours", zap.String("actualEndDateError", err1.Error()), zap.String("scheduledStartDateError", err2.Error()))
				return
			}
		}

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTaskTable string
		if targetTable == MouldMaintenancePreventiveWorkOrderTable {
			workOrderTaskTable = MouldMaintenancePreventiveWorkOrderTaskTable
		} else {
			workOrderTaskTable = MouldMaintenanceCorrectiveWorkOrderTaskTable
		}

		// complete work order will complete all the task assigned
		condition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, workOrderTaskTable, condition)
		if err == nil {
			for _, workOrderTaskObject := range *listOfTasks {
				workOrderTask := make(map[string]interface{})
				json.Unmarshal(workOrderTaskObject.ObjectInfo, &workOrderTask)
				if util.InterfaceToInt(workOrderTask["taskStatus"]) != WorkOrderTaskDone {
					workOrderTask["taskStatus"] = WorkOrderTaskDone

					updateObject = make(map[string]interface{})
					updateObject["object_info"], _ = json.Marshal(workOrderTask)

					Update(dbConnection, workOrderTaskTable, workOrderTaskObject.Id, updateObject)
				}
			}
		}

		// now put the machine on Active
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
			_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		}

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, CorrectiveWorkOrderCompletionNotification, workflowUser, MouldMaintenanceCorrectiveWorkOrderComponent, eventObject.Id)
				fmt.Println("Errror:", err)
			}
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance event is successfully updated",
		Code:    0,
	})
}

// completeOrder ShowAccount godoc
// @Summary schedule orders
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId}/action/schedule [get]
func (v *MaintenanceService) completeOrder(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)
	if componentName == "maintenance_work_order" || componentName == "maintenance_corrective_work_order" {
		//maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: eventObject.ObjectInfo}
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = false
		workOrderInfo["canUpdate"] = false
		workOrderInfo["canForceStop"] = false
		workOrderInfo["workOrderActualEndDate"] = util.GetCurrentTime(ISOTimeLayout)

		existingActionRemarks := util.AppendToObjectArray(workOrderInfo["actionRemarks"], ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "COMPLETED",
			UserId:        userId,
			Remarks:       "COMPLETED",
			ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		})
		//existingActionRemarks = append(existingActionRemarks, ActionRemarks{
		//	ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
		//	Status:        "COMPLETED",
		//	UserId:        userId,
		//	Remarks:       "COMPLETED",
		//	ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		//})

		workOrderInfo["actionRemarks"] = existingActionRemarks

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTaskTable string
		if targetTable == MaintenanceWorkOrderTable {
			workOrderTaskTable = MaintenanceWorkOrderTaskTable
		} else {
			workOrderTaskTable = MaintenanceWorkOrderCorrectiveTaskComponent
		}

		// complete work order will complete all the task assigned
		condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = \"" + strconv.Itoa(recordId) + "\""
		listOfTasks, err := GetConditionalObjects(dbConnection, workOrderTaskTable, condition)
		if err == nil {
			for _, workOrderTaskObject := range *listOfTasks {
				workOrderTask := make(map[string]interface{})
				json.Unmarshal(workOrderTaskObject.ObjectInfo, &workOrderTask)
				if util.InterfaceToInt(workOrderTask["taskStatus"]) != WorkOrderTaskDone {
					workOrderTask["taskStatus"] = WorkOrderTaskDone

					updateObject = make(map[string]interface{})
					updateObject["object_info"], _ = json.Marshal(workOrderTask)

					Update(dbConnection, workOrderTaskTable, workOrderTaskObject.Id, updateObject)
				}
			}
		}

		// now put the machine on Active
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
			_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		}

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, CorrectiveWorkOrderCompletionNotification, workflowUser, MaintenanceCorrectiveWorkOrderComponent, eventObject.Id)
				fmt.Println("Errror:", err)
			}
		}

	} else if componentName == "maintenance_work_order_task" || componentName == "maintenance_work_order_corrective_task" {
		//maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: eventObject.ObjectInfo}
		workOrderTaskInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderTaskInfo)
		workOrderTaskInfo["taskStatus"] = WorkOrderTaskDone
		workOrderTaskInfo["canComplete"] = false

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTable string
		if componentName == MaintenanceWorkOrderTaskTable {
			workOrderTable = MaintenanceWorkOrderTable
		} else {
			workOrderTable = MaintenanceWorkOrderCorrectiveTaskComponent
		}

		if checkAllTaskAreDone(util.InterfaceToInt(workOrderTaskInfo["workOrderId"]), dbConnection, componentName) {
			_, workOrderObject := Get(dbConnection, workOrderTable, util.InterfaceToInt(workOrderTaskInfo["workOrderId"]))
			//maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrderObject.ObjectInfo}
			workOrderInfo := make(map[string]interface{})
			json.Unmarshal(workOrderObject.ObjectInfo, &workOrderInfo)
			workOrderInfo["canComplete"] = true

			workOrderInfo["workOrderStatus"] = WorkOrderDone

			updateObject = make(map[string]interface{})
			updateObject["object_info"], _ = json.Marshal(workOrderTaskInfo)

			err = Update(dbConnection, workOrderTable, util.InterfaceToInt(workOrderTaskInfo["workOrderId"]), updateObject)

			if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
				mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
				err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
				machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
				err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
				if err != nil {
					response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
						&response.DetailedError{
							Header:      getError(common.UpdateResourceFailedError).Error(),
							Description: "Internal system error happened during resource update, please report error code to system admin",
						})
					return
				}
			}
		}
	} else if componentName == MouldMaintenancePreventiveWorkOrderComponent || componentName == MouldMaintenanceCorrectiveWorkOrderComponent {
		workOrderInfo := make(map[string]interface{})
		json.Unmarshal(eventObject.ObjectInfo, &workOrderInfo)
		workOrderInfo["workOrderStatus"] = WorkOrderDone
		workOrderInfo["canComplete"] = false
		workOrderInfo["canUpdate"] = false
		workOrderInfo["canForceStop"] = false
		workOrderInfo["workOrderActualEndDate"] = util.GetCurrentTime(ISOTimeLayout)

		existingActionRemarks := util.AppendToObjectArray(workOrderInfo["actionRemarks"], ActionRemarks{
			ExecutedTime:  util.GetCurrentTime(ISOTimeLayout),
			Status:        "COMPLETED",
			UserId:        userId,
			Remarks:       "COMPLETED",
			ProcessedTime: getTimeDifference(util.InterfaceToString(workOrderInfo["createdAt"])),
		})

		workOrderInfo["actionRemarks"] = existingActionRemarks

		updateObject := make(map[string]interface{})
		updateObject["object_info"], _ = json.Marshal(workOrderInfo)

		err = Update(dbConnection, targetTable, recordId, updateObject)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
			return
		}

		var workOrderTaskTable string
		if targetTable == MouldMaintenancePreventiveWorkOrderTable {
			workOrderTaskTable = MouldMaintenancePreventiveWorkOrderTaskTable
		} else {
			workOrderTaskTable = MouldMaintenanceCorrectiveWorkOrderTaskTable
		}

		// complete work order will complete all the task assigned
		condition := "object_info->>'$.workOrderId' = '" + strconv.Itoa(recordId) + "'"
		listOfTasks, err := GetConditionalObjects(dbConnection, workOrderTaskTable, condition)
		if err == nil {
			for _, workOrderTaskObject := range *listOfTasks {
				workOrderTask := make(map[string]interface{})
				json.Unmarshal(workOrderTaskObject.ObjectInfo, &workOrderTask)
				if util.InterfaceToInt(workOrderTask["taskStatus"]) != WorkOrderTaskDone {
					workOrderTask["taskStatus"] = WorkOrderTaskDone

					updateObject = make(map[string]interface{})
					updateObject["object_info"], _ = json.Marshal(workOrderTask)

					Update(dbConnection, workOrderTaskTable, workOrderTaskObject.Id, updateObject)
				}
			}
		}

		// now put the machine on Active
		if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
			_ = machineService.MoveMachineLiveStatusToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassMoulds {
			mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			err = mouldInterface.PutToActiveMode(projectId, util.InterfaceToInt(workOrderInfo["mouldId"]), userId)
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassAssemblyMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveAssemblyMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		} else if util.InterfaceToString(workOrderInfo["moduleName"]) == AssetClassToolingMachines {
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			err = machineService.MoveToolingMachineToActive(projectId, util.InterfaceToInt(workOrderInfo["assetId"]))
			if err != nil {
				response.DispatchDetailedError(ctx, common.UpdateResourceFailed,
					&response.DetailedError{
						Header:      getError(common.UpdateResourceFailedError).Error(),
						Description: "Internal system error happened during resource update, please report error code to system admin",
					})
				return
			}
		}

		var notificationEmailTemplate int
		if componentName == MouldMaintenancePreventiveWorkOrderComponent {
			notificationEmailTemplate = MouldPreventiveWorkOrderCompletionNotification
		} else {
			notificationEmailTemplate = MouldCorrectiveWorkOrderCompletionNotification
		}

		// Send notification to configured user in creation status
		_, generalWorkOrderStatus := Get(dbConnection, MaintenanceWorkOrderStatusTable, WorkOrderDone)
		workOrderStatus := MaintenanceWorkOrderStatus{ObjectInfo: generalWorkOrderStatus.ObjectInfo}
		listOfWorkflowUsers := workOrderStatus.getMaintenanceWorkOrderStatusInfo().NotificationUserList

		if workOrderStatus.getMaintenanceWorkOrderStatusInfo().IsEmailNotificationEnabled {
			for _, workflowUser := range listOfWorkflowUsers {
				err = v.emailGenerator(dbConnection, notificationEmailTemplate, workflowUser, componentName, eventObject.Id)
				fmt.Println("Errror:", err)
			}
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Maintenance mould event is successfully updated",
		Code:    0,
	})
}

func checkAllTaskAreDone(workOrderId int, dbConnection *gorm.DB, targetTable string) bool {
	isAllDone := false
	taskCondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = " + strconv.Itoa(workOrderId)
	listOfTasks, err := GetConditionalObjects(dbConnection, targetTable, taskCondition)

	if err != nil {
		return isAllDone
	}

	noOfTasks := len(*listOfTasks)
	if noOfTasks == 0 {
		return true
	}

	for _, object := range *listOfTasks {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: object.ObjectInfo}
		workOrderTaskInfo := maintenanceWorkOrderTask.getWorkOrderTaskInfo()

		if workOrderTaskInfo.TaskStatus == WorkOrderTaskDone {
			noOfTasks -= 1
		}
	}

	if noOfTasks == 0 {
		isAllDone = true
	}

	return isAllDone
}

func (v *MaintenanceService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Resource",
				Description: "The resource that you are trying to delete doesn't exist, Please check refresh page and try again",
			})
		return
	}
	if component.IsArchived(generalObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Resource Archived",
				Description: "The resource that you are trying to delete is already archived. This operation is not allowed",
			})
		return
	}
	var dependencyComponents []string
	var dependencyRecords int
	var canDelete = true
	v.checkReference(dbConnection, componentName, recordId, &dependencyComponents, &dependencyRecords)
	if dependencyRecords > 0 {
		var dependencyString string
		canDelete = false
		dependencyComponents = util.RemoveDuplicateString(dependencyComponents)
		dependencyString = " ["
		for index, dependencyComponent := range dependencyComponents {
			if index == len(dependencyComponents)-1 {
				dependencyString += dependencyComponent
			} else {
				dependencyString += dependencyComponent + " ->"
			}
		}
		dependencyString += " ]"
		ctx.JSON(http.StatusOK, response.ValidationResponse{
			Code:      100,
			Message:   "There are dependencies bound to the resource that you are trying to remove. Removing this resource would create the chain removal on following resources " + dependencyString + " in " + strconv.Itoa(dependencyRecords) + " places, Please understand the risk of deleting this resource as all the dependencies would be achieved immediately, and this process is not reversible",
			CanDelete: canDelete,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.ValidationResponse{
		Code:      100,
		Message:   "There are no dependencies bound to the resource that you are trying to remove. So, removing this resource won't affect others resource now, you can proceed !!",
		CanDelete: canDelete,
	})
}

func ConvertRemarkArray(remark interface{}) []ActionRemarks {
	data := remark.([]map[string]interface{})
	finalObject := make([]ActionRemarks, 0)
	for _, action := range data {
		insertedData := ActionRemarks{
			ExecutedTime:  util.InterfaceToString(action["executedTime"]),
			Status:        util.InterfaceToString(action["status"]),
			UserId:        util.InterfaceToInt(action["userId"]),
			Remarks:       util.InterfaceToString(action["remarks"]),
			ProcessedTime: util.InterfaceToString(action["processedTime"]),
		}
		finalObject = append(finalObject, insertedData)
	}
	return finalObject
}

func hasAuthorizationOnTask(userList []common.UserBasicInfo, currentUserId int) bool {

	for _, userId := range userList {
		if currentUserId == userId.UserId {
			return true
		}
	}
	return false
}
