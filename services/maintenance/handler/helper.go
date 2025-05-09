package handler

import (
	"gorm.io/gorm"
	"strconv"
)

func (v *MaintenanceService) isTaskDone(dbConnection *gorm.DB, taskId int) (error, bool) {

	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		if maintenanceWorkOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
			return nil, true
		}
	}
	return err, false
}

func (v *MaintenanceService) isTaskStarted(dbConnection *gorm.DB, taskId int) (error, bool) {

	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		if maintenanceWorkOrderTask.getWorkOrderTaskInfo().TaskStatus > WorkOrderCreated {
			return nil, true
		}
	}
	return err, false
}

func (v *MaintenanceService) isWorkOrderCompletedFromTaskId(dbConnection *gorm.DB, taskId int) (error, bool) {
	var err error
	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		err, workOrder := Get(dbConnection, MaintenanceWorkOrderTable, maintenanceWorkOrderTask.getWorkOrderTaskInfo().WorkOrderId)
		if err == nil {
			maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrder.ObjectInfo}
			workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
			if workOrderInfo.WorkOrderStatus == WorkOrderDone {
				return nil, true
			}
		}
	}

	return err, false

}

func (v *MaintenanceService) isWorkOrderScheduledFromTaskId(dbConnection *gorm.DB, taskId int) (error, bool) {
	var err error
	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		err, workOrder := Get(dbConnection, MaintenanceWorkOrderTable, maintenanceWorkOrderTask.getWorkOrderTaskInfo().WorkOrderId)
		if err == nil {
			maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrder.ObjectInfo}
			workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
			if workOrderInfo.WorkOrderStatus == WorkOrderScheduled {
				return nil, true
			}
		}
	}

	return err, false

}

func (v *MaintenanceService) isWorkOrderCompletedFromWorkOrderId(dbConnection *gorm.DB, workOrderId int) (error, bool) {
	err, workOrder := Get(dbConnection, MaintenanceWorkOrderTable, workOrderId)
	if err == nil {
		maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrder.ObjectInfo}
		workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
		if workOrderInfo.WorkOrderStatus == WorkOrderDone {
			return nil, true
		}
	}

	return err, false

}

func (v *MaintenanceService) completeWorkOrder(userId int, dbConnection *gorm.DB, taskId int) error {
	var err error
	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		workOrderId := maintenanceWorkOrderTask.getWorkOrderTaskInfo().WorkOrderId
		err, workOrder := Get(dbConnection, MaintenanceWorkOrderTable, workOrderId)
		if err == nil {
			maintenanceWorkOrder := MaintenanceWorkOrder{ObjectInfo: workOrder.ObjectInfo}
			workOrderInfo := maintenanceWorkOrder.getWorkOrderInfo()
			workOrderInfo.WorkOrderStatus = WorkOrderDone
			err = Update(dbConnection, MaintenanceWorkOrderTable, workOrderId, workOrderInfo.DatabaseSerialize(userId))
		}
	}
	return err
}

func (v *MaintenanceService) isAllTaskCompleted(dbConnection *gorm.DB, taskId int) (error, bool) {

	err, workOrderTask := Get(dbConnection, MaintenanceWorkOrderTaskTable, taskId)
	if err == nil {
		maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
		workOrderId := maintenanceWorkOrderTask.getWorkOrderTaskInfo().WorkOrderId
		// now load all the
		conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.workOrderId\")) = " + strconv.Itoa(workOrderId)
		listOfWorkOrderTasks, err := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskTable, conditionString)
		var totalWorkOrderTasks = len(*listOfWorkOrderTasks)
		if err == nil {
			for _, workOrderTasksInterface := range *listOfWorkOrderTasks {
				workOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTasksInterface.ObjectInfo}
				if workOrderTask.getWorkOrderTaskInfo().TaskStatus == WorkOrderTaskDone {
					totalWorkOrderTasks = totalWorkOrderTasks - 1
				}
			}
			if totalWorkOrderTasks == 0 {
				return nil, true
			}
		}
	}
	return err, false
}
