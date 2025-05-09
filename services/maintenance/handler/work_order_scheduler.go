package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

func (v *MaintenanceService) OnTimer() {

	//c := cron.New()
	//c.AddFunc("@every 120s", func() { v.createWorkOrders() })
	//c.Start()
}

func (v *MaintenanceService) createWorkOrders() {
	projectId := "906d0fd569404c59956503985b330132"
	workOrders := getAllWorkOrders(v, projectId)
	dbConnection := v.BaseService.ReferenceDatabase
	now := time.Now()
	var workOrderInfo WorkOrderInfo
	v.BaseService.Logger.Info("starting the work order processing...")
	for _, workOrder := range workOrders {
		err := json.Unmarshal(workOrder.ObjectInfo, &workOrderInfo)
		if err != nil {
			continue
		}
		v.BaseService.Logger.Info("processing the work-order", zap.Any("work_order_name", workOrderInfo.Name))
		if workOrderInfo.IsParent {
			childWorkOrderInfo := workOrderInfo
			nextTimeWorkOrder := cronexpr.MustParse(workOrderInfo.WorkOrderRepetitiveCron).Next(time.Unix(workOrderInfo.LastTimeWorkOrderCreation, 0))
			if nextTimeWorkOrder.Unix() <= now.Unix() {
				childWorkOrderInfo.WorkOrderScheduledStartDate = time.Unix(nextTimeWorkOrder.Unix(), 0).Format(ISOTimeLayout)
				childWorkOrderInfo.WorkOrderScheduledEndDate = time.Unix(nextTimeWorkOrder.Unix(), 0).Format(ISOTimeLayout)
				childWorkOrderInfo.IsParent = false
				childWorkOrderInfo.IsRepetitive = false

				insertWorkOrder(v, projectId, childWorkOrderInfo)

				workOrderInfo.LastTimeWorkOrderCreation = nextTimeWorkOrder.Unix()
				updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
			}
		}

		// We have to send email to supervisor based on frequency
		workOrderStartDate := util.ConvertStringToDateTime(workOrderInfo.RemainderDate)
		workOrderStartDateTs := workOrderStartDate.DateTimeEpoch

		workOrderEndDate := util.ConvertStringToDateTime(workOrderInfo.RemainderEndDate)
		workOrderEndDateTs := workOrderEndDate.DateTimeEpoch

		if now.Unix() >= workOrderStartDateTs && workOrderEndDateTs >= now.Unix() {
			lastTimeEmailSend := util.ConvertStringToDateTime(workOrderInfo.EmailLastSendDate)
			lastTimeEmailSendTs := lastTimeEmailSend.DateTimeEpoch
			fmt.Println("workOrderInfo.RepeatFrequency .Name", workOrderInfo.RepeatFrequency)
			if workOrderInfo.RepeatFrequency == DailyRepeat {
				if workOrderInfo.IsFirstTime {
					fmt.Println("IsFirstTime", workOrderInfo.IsFirstTime)
					for _, userId := range workOrderInfo.Supervisors {
						err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
						if err != nil {
							fmt.Println("errr")
							v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
						}
					}
					workOrderInfo.IsFirstTime = false
					workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
					updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
				} else {
					considerDayInterval := 86400 * workOrderInfo.RepeatInterval
					nextTimeEmailSend := lastTimeEmailSendTs + int64(considerDayInterval)

					if now.Unix() >= nextTimeEmailSend {
						for _, userId := range workOrderInfo.Supervisors {
							err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
							if err != nil {
								v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
							}
						}
						workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
						updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
					}
				}
			} else if workOrderInfo.RepeatFrequency == WeeklyRepeat {
				if workOrderInfo.IsFirstTime {
					for _, userId := range workOrderInfo.Supervisors {
						err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
						if err != nil {
							v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
						}
					}
					workOrderInfo.IsFirstTime = false
					workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
					updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
				} else {
					considerDayInterval := 604800 * workOrderInfo.RepeatInterval
					nextTimeEmailSend := lastTimeEmailSendTs + int64(considerDayInterval)

					if now.Unix() >= nextTimeEmailSend {
						for _, userId := range workOrderInfo.Supervisors {
							err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
							if err != nil {
								v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
							}
						}
						workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
						updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
					}
				}

			} else if workOrderInfo.RepeatFrequency == MonthlyRepeat {
				if workOrderInfo.IsFirstTime {
					for _, userId := range workOrderInfo.Supervisors {
						err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
						if err != nil {
							v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
						}
					}
					workOrderInfo.IsFirstTime = false
					workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
					updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
				} else {
					considerDayInterval := 2629743 * workOrderInfo.RepeatInterval
					nextTimeEmailSend := lastTimeEmailSendTs + int64(considerDayInterval)

					if now.Unix() >= nextTimeEmailSend {
						for _, userId := range workOrderInfo.Supervisors {
							err = v.emailGenerator(dbConnection, EmailNotificationPreventiveOrder, userId, MaintenanceWorkOrderComponent, workOrder.Id)
							if err != nil {
								v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
							}
						}
						workOrderInfo.EmailLastSendDate = util.GetCurrentTime(ISOTimeLayout)
						updateWorkOrder(v, projectId, workOrder.Id, workOrderInfo)
					}
				}
			} else {
				continue
			}
		}
	}
}

func getAllWorkOrders(ms *MaintenanceService, projectId string) []MaintenanceWorkOrder {
	dbConnection := ms.BaseService.ServiceDatabases["906d0fd569404c59956503985b330132"]
	workOrderQuery := "select * from maintenance_work_order"
	var workOrders []MaintenanceWorkOrder
	dbConnection.Raw(workOrderQuery).Scan(&workOrders)

	return workOrders
}

func updateWorkOrder(ms *MaintenanceService, projectId string, id int, workOrderInfo WorkOrderInfo) {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]

	workOrderObjectInfo, err := json.Marshal(workOrderInfo)

	if err != nil {
		return
	}

	dbError := dbConnection.Model(&MaintenanceWorkOrder{}).Where("id = ?", id).Update("object_info", workOrderObjectInfo).Error

	if dbError != nil {
		ms.BaseService.Logger.Error("error update work order info", zap.String("error", err.Error()))
	}
}

func insertWorkOrder(ms *MaintenanceService, projectId string, workOrderInfo WorkOrderInfo) {
	dbConnection := ms.BaseService.ServiceDatabases[projectId]

	workOrderObjectInfo, err := json.Marshal(workOrderInfo)

	if err != nil {
		return
	}

	object := component.GeneralObject{
		ObjectInfo: workOrderObjectInfo,
	}
	dbError, _ := Create(dbConnection, MaintenanceWorkOrderTable, object)

	if dbError != nil {
		ms.BaseService.Logger.Error("error insert work order info", zap.String("error", err.Error()))
	}
}
