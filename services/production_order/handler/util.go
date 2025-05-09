package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func (v *ProductionOrderService) getOrderStatusId(dbConnection *gorm.DB, preferenceLevel int) int {
	listOfOrderStatusObjects, err := GetObjects(dbConnection, ProductionOrderStatusTable)
	if err != nil {
		return -1
	}

	for _, orderStatusInterface := range *listOfOrderStatusObjects {
		productionOrderStatus := ProductionOrderStatus{
			Id:         orderStatusInterface.Id,
			ObjectInfo: orderStatusInterface.ObjectInfo,
		}
		if productionOrderStatus.getProductionOrderStatusInfo().Preference == preferenceLevel {
			return productionOrderStatus.Id
		}
	}
	return -1
}

func (v *ProductionOrderService) getColorCode(dbConnection *gorm.DB, orderStatusId int) (error, string) {
	err, orderStatusGeneralObject := Get(dbConnection, ProductionOrderStatusTable, orderStatusId)
	if err != nil {
		return err, ""
	}

	productionOrderStatus := ProductionOrderStatus{
		Id:         orderStatusGeneralObject.Id,
		ObjectInfo: orderStatusGeneralObject.ObjectInfo,
	}
	return nil, productionOrderStatus.getProductionOrderStatusInfo().ColorCode
}

func getStartDateAndEndDate(offsetDays int) (string, string) {
	currentTime := time.Now()
	after := currentTime.AddDate(0, 0, offsetDays)

	year, month, day := after.Date()
	minute := after.Minute()
	hour := after.Hour()
	seconds := after.Second()

	var monthString string
	var dayString string
	var hourString string
	var minString string
	var secondString string
	if int(month) < 11 {
		monthString = "0" + strconv.Itoa(int(month))
	} else {
		monthString = strconv.Itoa(int(month))
	}

	if int(day) < 10 {
		dayString = "0" + strconv.Itoa(int(day))
	} else {
		dayString = strconv.Itoa(int(day))
	}

	if int(hour) < 10 {
		hourString = "0" + strconv.Itoa(int(hour))
	} else {
		hourString = strconv.Itoa(int(hour))
	}
	if int(minute) < 10 {
		minString = "0" + strconv.Itoa(int(minute))
	} else {
		minString = strconv.Itoa(int(minute))
	}
	if int(seconds) < 10 {
		secondString = "0" + strconv.Itoa(int(seconds))
	} else {
		secondString = strconv.Itoa(int(seconds))
	}

	fmt.Println("second string :", secondString)
	startDate := strconv.Itoa(year) + "-" + monthString + "-" + dayString + "T" + hourString + ":" + minString + ":" + "00" + ".000Z"
	fmt.Println("preparing schedule [startDate] :::::::::::::::::", startDate)
	endSchedule := after.Add(20 * time.Hour)
	year, month, day = endSchedule.Date()
	if int(day) < 10 {
		dayString = "0" + strconv.Itoa(int(day))
	} else {
		dayString = strconv.Itoa(int(day))
	}
	if int(month) < 10 {
		monthString = "0" + strconv.Itoa(int(month))
	} else {
		monthString = strconv.Itoa(int(month))
	}
	endHour := endSchedule.Hour()
	if endHour < 10 {
		hourString = "0" + strconv.Itoa(endHour)
	} else {
		hourString = strconv.Itoa(endHour)
	}

	minute = endSchedule.Minute()
	if int(minute) < 10 {
		minString = "0" + strconv.Itoa(int(minute))
	} else {
		minString = strconv.Itoa(int(minute))
	}
	seconds = endSchedule.Second()

	if int(seconds) < 10 {
		secondString = "0" + strconv.Itoa(int(seconds))
	} else {
		secondString = strconv.Itoa(int(seconds))
	}
	endDate := strconv.Itoa(year) + "-" + monthString + "-" + dayString + "T" + hourString + ":" + minString + ":" + "00" + ".000Z"

	return startDate, endDate
}

func (v *ProductionOrderService) createSystemNotification(projectId, header, description string, recordId int) error {
	systemNotification := common.SystemNotification{}
	systemNotification.Name = header
	systemNotification.ColorCode = "#14F44E"
	systemNotification.IconCls = "icon-park-outline:transaction-order"
	systemNotification.RecordId = recordId
	systemNotification.RouteLinkComponent = "timeline"
	systemNotification.Component = "Production Order"
	systemNotification.Description = description
	systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
	rawSystemNotification, _ := json.Marshal(systemNotification)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	err := notificationService.CreateSystemNotification(projectId, rawSystemNotification)
	return err
}
