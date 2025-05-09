package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ScheduleRequest struct {
	StartDate          string      `json:"startDate"`
	EndDate            string      `json:"endDate"`
	ScheduledQty       interface{} `json:"scheduledQty"`
	CompletedQty       interface{} `json:"completedQty"`
	RejectedQty        interface{} `json:"rejectedQty"`
	Name               string      `json:"name"`
	EventSourceId      int         `json:"eventSourceId"`
	EnableCustomCavity bool        `json:"enableCustomCavity"`
	MouldId            int         `json:"mouldId"`
	CustomCavity       int         `json:"customCavity"`
	MouldUp            string      `json:"mouldUp"`
	MouldDown          string      `json:"mouldDown"`
	IsRecoverySchedule bool        `json:"isRecoverySchedule"`
	RecoveryScheduleId int         `json:"recoveryScheduleId"`
}

func initEventWithDates(productionOrder string, scheduleRequest ScheduleRequest, machineId, orderStatusId int) datatypes.JSON {
	var createEvent = make(map[string]interface{})
	createEvent["name"] = scheduleRequest.Name
	createEvent["productionOrder"] = productionOrder
	createEvent["eventType"] = "production_schedule"
	createEvent["eventStatus"] = orderStatusId
	createEvent["scheduledQty"] = scheduleRequest.ScheduledQty
	createEvent["draggable"] = true
	createEvent["startDate"] = scheduleRequest.StartDate
	createEvent["endDate"] = scheduleRequest.EndDate
	createEvent["eventSourceId"] = scheduleRequest.EventSourceId
	createEvent["machineId"] = machineId
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0
	createEvent["completedQty"] = 0
	createEvent["rejectedQty"] = 0
	createEvent["isUpdate"] = true
	createEvent["canForceStop"] = false
	createEvent["isRecoverySchedule"] = scheduleRequest.IsRecoverySchedule
	createEvent["recoveryScheduleId"] = scheduleRequest.RecoveryScheduleId

	if scheduleRequest.IsRecoverySchedule {
		createEvent["iconCls"] = "b-fa b-fa-bullseye"
	} else {
		createEvent["iconCls"] = "b-fa fa-anchor"
	}
	createEvent["priorityLevel"] = 1
	createEvent["plannedManpower"] = 0
	createEvent["remarks"] = ""
	rawEvent, _ := json.Marshal(createEvent)
	return rawEvent
}

func initEvent(productionOrder string, dayOffset int, dailyRate interface{}, eventSourceId int, machineId int, orderStatusId int) datatypes.JSON {
	var createEvent = make(map[string]interface{})
	createEvent["name"] = productionOrder + " - batch " + strconv.Itoa(dayOffset)
	createEvent["productionOrder"] = productionOrder
	createEvent["iconCls"] = "b-fa fa-anchor"
	createEvent["eventType"] = "production_schedule"
	createEvent["eventStatus"] = orderStatusId
	createEvent["scheduledQty"] = dailyRate
	createEvent["draggable"] = true
	startDate, endDate := getStartDateAndEndDate(dayOffset)
	createEvent["startDate"] = startDate
	createEvent["endDate"] = endDate
	createEvent["eventSourceId"] = eventSourceId
	createEvent["machineId"] = machineId
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0
	createEvent["completedQty"] = 0
	createEvent["rejectedQty"] = 0
	createEvent["isAbortEnabled"] = true
	createEvent["isUpdate"] = true
	createEvent["canForceStop"] = false
	createEvent["priorityLevel"] = 1
	createEvent["plannedManpower"] = 0
	createEvent["remarks"] = ""
	rawEvent, _ := json.Marshal(createEvent)
	return rawEvent
}

func initSplitEvent(startDate string, endDate string, productionOrder string, dayOffset int, dailyRate interface{}, eventSourceId int, machineId int, orderStatusId int) map[string]interface{} {

	var createEvent = make(map[string]interface{})
	createEvent["name"] = productionOrder + " - batch " + strconv.Itoa(dayOffset)
	createEvent["productionOrder"] = productionOrder
	createEvent["iconCls"] = "b-fa fa-anchor"
	createEvent["eventType"] = "production_schedule"
	createEvent["eventStatus"] = orderStatusId
	createEvent["scheduledQty"] = dailyRate
	createEvent["draggable"] = true
	createEvent["startDate"] = util.ConvertSingaporeTimeToUTC(startDate)
	createEvent["endDate"] = util.ConvertSingaporeTimeToUTC(endDate)
	createEvent["eventSourceId"] = eventSourceId
	createEvent["machineId"] = machineId
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0
	createEvent["completedQty"] = 0
	createEvent["rejectedQty"] = 0
	createEvent["isAbortEnabled"] = true
	createEvent["isUpdate"] = true
	createEvent["canForceStop"] = false
	createEvent["priorityLevel"] = 1
	createEvent["plannedManpower"] = 0
	createEvent["remarks"] = ""
	return createEvent
}

func initToolingSplitEvent(startDate string, endDate string, productionOrder string, dayOffset int, eventSourceId int, machineId int, orderStatusId int, toolingPartId int) map[string]interface{} {

	var createEvent = make(map[string]interface{})
	createEvent["name"] = productionOrder + " - batch " + strconv.Itoa(dayOffset)
	createEvent["productionOrder"] = productionOrder
	createEvent["iconCls"] = "b-fa fa-anchor"
	createEvent["eventType"] = "tooling_schedule"
	createEvent["eventStatus"] = orderStatusId
	createEvent["draggable"] = true
	createEvent["startDate"] = util.ConvertSingaporeTimeToUTC(startDate)
	createEvent["endDate"] = util.ConvertSingaporeTimeToUTC(endDate)
	createEvent["eventSourceId"] = eventSourceId
	createEvent["machineId"] = machineId
	createEvent["objectStatus"] = "Active"
	createEvent["partId"] = toolingPartId
	createEvent["percentDone"] = 0
	createEvent["completedQty"] = 0
	createEvent["rejectedQty"] = 0
	createEvent["isAbortEnabled"] = true
	createEvent["canSetupTime"] = true
	createEvent["isUpdate"] = true
	createEvent["canForceStop"] = false
	createEvent["priorityLevel"] = 1
	createEvent["plannedManpower"] = 0
	createEvent["remarks"] = ""
	return createEvent
}

func initToolingEventWithDates(productionOrder string, scheduleRequest map[string]interface{}, orderStatusId int) datatypes.JSON {
	var createEvent = make(map[string]interface{})
	createEvent["name"] = scheduleRequest["name"]
	createEvent["productionOrder"] = productionOrder
	createEvent["iconCls"] = "b-fa fa-anchor"
	createEvent["eventType"] = "tooling_schedule"
	createEvent["eventStatus"] = orderStatusId
	createEvent["draggable"] = true
	createEvent["startDate"] = scheduleRequest["startDate"]
	createEvent["endDate"] = scheduleRequest["endDate"]
	createEvent["eventSourceId"] = scheduleRequest["eventSourceId"]
	createEvent["machineId"] = scheduleRequest["machineId"]
	createEvent["partId"] = scheduleRequest["partId"]
	createEvent["objectStatus"] = "Active"
	createEvent["percentDone"] = 0
	createEvent["completedQty"] = 0
	createEvent["rejectedQty"] = 0
	createEvent["isUpdate"] = true
	createEvent["canForceStop"] = false
	createEvent["canSetupTime"] = true
	createEvent["priorityLevel"] = 1
	createEvent["plannedManpower"] = 0
	createEvent["remarks"] = ""
	rawEvent, _ := json.Marshal(createEvent)
	return rawEvent
}

func SplitDailySchedule(iOrderQty int,
	iDailyRate int,
	fullAllocation int,
	remainder int,
	interval int,
	startDateTimeStr string,
	endTime string,
	rangeObject map[string]interface{},
	productionOrder string,
	eventSourceId int, machineId int, orderStatusId int, noOfExistingEvent int, isEndTimeSpecified bool) ([]map[string]interface{}, error, int) {

	startDateTime, _ := time.Parse(TimeLayout, startDateTimeStr)
	rangeType := util.InterfaceToString(rangeObject["type"])

	prefreneceLevel := ScheduleStatusPreferenceThree

	scheduledDailyRate := 0
	listOfScheduledEvents := make([]map[string]interface{}, 0)

	var timeList []string
	var endHours int
	var endMinutes int

	if !isEndTimeSpecified {
		timeList = strings.Split(endTime, ":")
		endHours, _ = strconv.Atoi(timeList[0])
		endMinutes, _ = strconv.Atoi(timeList[1])

	}

	if rangeType == "endDate" {
		var endDateTimeStr string
		var endDateTime time.Time
		if isEndTimeSpecified {
			endHoursStr := strconv.Itoa(startDateTime.Hour())
			endMinutesStr := strconv.Itoa(startDateTime.Minute())

			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endHoursStr + ":" + endMinutesStr + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		} else {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endTime + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		}

		index := 0

		for endDateTime.After(startDateTime) {
			if scheduledDailyRate > iOrderQty {
				return listOfScheduledEvents, getError("Schedule quantity exceeds the order quantity"), prefreneceLevel
			}
			index += 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				//This is indicator for remainder
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					eventEndDateTime = startDateTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}
				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+index, remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				eventEndDateTime = startDateTime.AddDate(0, 0, 1)
			} else {
				eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+index, iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)

			startDateTime = startDateTime.AddDate(0, 0, interval)

		}

		if scheduledDailyRate < iOrderQty {
			prefreneceLevel = ScheduleStatusPreferenceTwo
		}

	} else if rangeType == "numbered" {
		//Repeat this event 10 times
		noOfOccurrence := util.InterfaceToInt(rangeObject["numberOfOccurrences"])
		index := 0

		for index < noOfOccurrence {
			if scheduledDailyRate > iOrderQty {
				return listOfScheduledEvents, getError("Schedule quantity exceeds the order quantity"), prefreneceLevel
			}

			index += 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				//This is indicator for remainder
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					eventEndDateTime = startDateTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}

				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				eventEndDateTime = startDateTime.AddDate(0, 0, 1)
			} else {
				eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)

			startDateTime = startDateTime.AddDate(0, 0, interval)

		}

		if scheduledDailyRate < iOrderQty {
			prefreneceLevel = ScheduleStatusPreferenceTwo
		}

	} else if rangeType == "noEnd" {
		// 56000 , 5000
		index := 0
		for scheduledDailyRate < iOrderQty {
			index += 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)
				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					endTime := time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), startDateTime.Hour(), startDateTime.Minute(), 0, 0, time.UTC)
					eventEndDateTime = endTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}

				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				endTime := time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), startDateTime.Hour(), startDateTime.Minute(), 0, 0, time.UTC)
				eventEndDateTime = endTime.AddDate(0, 0, 1)

			} else {
				eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
			startDateTime = startDateTime.AddDate(0, 0, 1*interval)

		}

	} else {
		return listOfScheduledEvents, getError("invalid range type"), prefreneceLevel
	}

	return listOfScheduledEvents, nil, prefreneceLevel

}

func SplitDailyAssemblySchedule(iOrderQty int,
	iDailyRate int,
	fullAllocation int,
	remainder int,
	interval int,
	startDateTimeStr string,
	endTime string,
	rangeObject map[string]interface{},
	productionOrder string,
	eventSourceId int, machineId int, orderStatusId int, noOfExistingEvent int, isEndTimeSpecified bool, shiftPeriod int) ([]map[string]interface{}, error, int) {

	startDateTime, _ := time.Parse(TimeLayout, startDateTimeStr)
	rangeType := util.InterfaceToString(rangeObject["type"])

	prefreneceLevel := ScheduleStatusPreferenceThree

	scheduledDailyRate := 0
	listOfScheduledEvents := make([]map[string]interface{}, 0)

	var timeList []string
	var endHours int
	var endMinutes int

	if !isEndTimeSpecified {
		timeList = strings.Split(endTime, ":")
		endHours, _ = strconv.Atoi(timeList[0])
		endMinutes, _ = strconv.Atoi(timeList[1])

	}

	if rangeType == "endDate" {
		var endDateTimeStr string
		var endDateTime time.Time
		if isEndTimeSpecified {
			endHoursStr := strconv.Itoa(startDateTime.Hour())
			endMinutesStr := strconv.Itoa(startDateTime.Minute())

			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endHoursStr + ":" + endMinutesStr + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		} else {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endTime + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		}

		index := 0

		for endDateTime.After(startDateTime) {

			if scheduledDailyRate > iOrderQty {
				return listOfScheduledEvents, getError("Schedule quantity exceeds the order quantity"), prefreneceLevel
			}
			index += 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				//This is indicator for remainder
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					eventEndDateTime = startDateTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}
				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+index, remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				eventEndDateTime = startDateTime.AddDate(0, 0, 1)
			} else {
				// eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				eventEndDateTime = startDateTime.Add(time.Duration(shiftPeriod) * time.Hour)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+index, iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)

			startDateTime = startDateTime.AddDate(0, 0, interval)

		}

		if scheduledDailyRate < iOrderQty {
			prefreneceLevel = ScheduleStatusPreferenceTwo
		}

	} else if rangeType == "numbered" {
		//Repeat this event 10 times
		noOfOccurrence := util.InterfaceToInt(rangeObject["numberOfOccurrences"])
		index := 0

		for index < noOfOccurrence {
			if scheduledDailyRate > iOrderQty {
				return listOfScheduledEvents, getError("Schedule quantity exceeds the order quantity"), prefreneceLevel
			}

			index += 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				//This is indicator for remainder
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					eventEndDateTime = startDateTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}

				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				eventEndDateTime = startDateTime.AddDate(0, 0, 1)
			} else {
				eventEndDateTime = startDateTime.Add(time.Duration(shiftPeriod) * time.Hour)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)

			startDateTime = startDateTime.AddDate(0, 0, interval)

		}

		if scheduledDailyRate < iOrderQty {
			prefreneceLevel = ScheduleStatusPreferenceTwo
		}

	} else if rangeType == "noEnd" {
		// 56000 , 5000
		index := 0
		remainingDays := fullAllocation
		for remainingDays > 0 {

			index += 1
			remainingDays -= 1
			scheduledDailyRate += iDailyRate

			if (iOrderQty - scheduledDailyRate) < 0 {
				remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)
				var eventEndDateTime time.Time
				if isEndTimeSpecified {
					endTime := time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), startDateTime.Hour(), startDateTime.Minute(), 0, 0, time.UTC)
					eventEndDateTime = endTime.AddDate(0, 0, 1)
				} else {
					eventEndDateTime = time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), endHours, endMinutes, 0, 0, time.UTC)
				}

				eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), remainderQty, eventSourceId, machineId, orderStatusId)
				listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
				break
			}

			var eventEndDateTime time.Time
			if isEndTimeSpecified {
				endTime := time.Date(startDateTime.Year(), startDateTime.Month(), startDateTime.Day(), startDateTime.Hour(), startDateTime.Minute(), 0, 0, time.UTC)
				eventEndDateTime = endTime.AddDate(0, 0, 1)

			} else {
				eventEndDateTime = startDateTime.Add(time.Duration(shiftPeriod) * time.Hour)
			}

			eventObject := initSplitEvent(startDateTime.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, (noOfExistingEvent + index), iDailyRate, eventSourceId, machineId, orderStatusId)
			listOfScheduledEvents = append(listOfScheduledEvents, eventObject)
			startDateTime = startDateTime.AddDate(0, 0, 1*interval)

		}

	} else {
		return listOfScheduledEvents, getError("invalid range type"), prefreneceLevel
	}

	return listOfScheduledEvents, nil, prefreneceLevel

}

func SplitWeeklySchedule(iOrderQty int,
	iDailyRate int,
	fullAllocation int,
	remainder int,
	interval int,
	startDateTimeStr string,
	endTime string,
	rangeObject map[string]interface{},
	patternObject map[string]interface{},
	productionOrder string,
	eventSourceId int, machineId int, orderStatusId int, noOfExistingEvent int, isEndTimeSpecified bool) ([]map[string]interface{}, error, int) {

	startDateTime, _ := time.Parse(TimeLayout, startDateTimeStr)
	rangeType := util.InterfaceToString(rangeObject["type"])
	productionDaysWeek := util.InterfaceToStringArray(patternObject["daysOfWeek"])

	eventScheduled := make([]map[string]interface{}, 0)

	preferenceLevel := ScheduleStatusPreferenceThree

	scheduledDailyRate := 0
	var timeList []string
	var endHours int
	var endMins int

	if !isEndTimeSpecified {
		timeList = strings.Split(endTime, ":")
		endHours, _ = strconv.Atoi(timeList[0])
		endMins, _ = strconv.Atoi(timeList[1])

	}

	if rangeType == "endDate" {
		var endDateTimeStr string
		var endDateTime time.Time
		if isEndTimeSpecified {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T23:59:00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		} else {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endTime + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		}

		incre := 0
		remainderBreak := false

		datePointer := startDateTime
		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)
		for endDateTime.After(datePointer) {
			if remainderBreak {
				break
			}
			for dateWeekPointer.After(datePointer) {
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate
					if (iOrderQty - scheduledDailyRate) < 0 {
						//This is indicator for remainder
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)
						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}

		if scheduledDailyRate < iOrderQty {
			preferenceLevel = ScheduleStatusPreferenceTwo
		}
	} else if rangeType == "numbered" {
		noOfOccurrence := util.InterfaceToInt(rangeObject["numberOfOccurrences"])

		datePointer := startDateTime

		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)
		incre := 0
		remainderBreak := false

		for incre < noOfOccurrence {
			for dateWeekPointer.After(datePointer) && incre < noOfOccurrence {
				if remainderBreak {
					break
				}
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate
					if (iOrderQty - scheduledDailyRate) < 0 {
						//This is indicator for remainder
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}
		if scheduledDailyRate < iOrderQty {
			preferenceLevel = ScheduleStatusPreferenceTwo
		}
	} else if rangeType == "noEnd" {
		datePointer := startDateTime

		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)

		remainderBreak := false

		incre := 0

		for scheduledDailyRate < iOrderQty {
			if remainderBreak {
				break
			}

			for dateWeekPointer.After(datePointer) && (incre <= fullAllocation) && (scheduledDailyRate < iOrderQty) {
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate

					if (iOrderQty - scheduledDailyRate) < 0 {
						//Reminder quantity
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}

	} else {
		fmt.Println("Non supported range type")
	}

	return eventScheduled, nil, preferenceLevel
}
func SplitWeeklyAssemblySchedule(iOrderQty int,
	iDailyRate int,
	fullAllocation int,
	remainder int,
	interval int,
	startDateTimeStr string,
	endTime string,
	rangeObject map[string]interface{},
	patternObject map[string]interface{},
	productionOrder string,
	eventSourceId int, machineId int, orderStatusId int, noOfExistingEvent int, isEndTimeSpecified bool, shiftPeriod int) ([]map[string]interface{}, error, int) {

	startDateTime, _ := time.Parse(TimeLayout, startDateTimeStr)
	rangeType := util.InterfaceToString(rangeObject["type"])
	productionDaysWeek := util.InterfaceToStringArray(patternObject["daysOfWeek"])

	eventScheduled := make([]map[string]interface{}, 0)

	preferenceLevel := ScheduleStatusPreferenceThree

	scheduledDailyRate := 0
	var timeList []string
	var endHours int
	var endMins int

	if !isEndTimeSpecified {
		timeList = strings.Split(endTime, ":")
		endHours, _ = strconv.Atoi(timeList[0])
		endMins, _ = strconv.Atoi(timeList[1])

	}

	if rangeType == "endDate" {
		var endDateTimeStr string
		var endDateTime time.Time
		if isEndTimeSpecified {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T23:59:00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		} else {
			endDateTimeStr = util.InterfaceToString(rangeObject["endDate"]) + "T" + endTime + ":00.000Z"
			endDateTime, _ = time.Parse(TimeLayout, endDateTimeStr)
		}

		incre := 0
		remainderBreak := false

		datePointer := startDateTime
		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)
		for endDateTime.After(datePointer) {
			if remainderBreak {
				break
			}
			for dateWeekPointer.After(datePointer) {
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate
					if (iOrderQty - scheduledDailyRate) < 0 {
						//This is indicator for remainder
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)
						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						eventEndDateTime = datePointer.Add(time.Duration(shiftPeriod) * time.Hour)
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}

		if scheduledDailyRate < iOrderQty {
			preferenceLevel = ScheduleStatusPreferenceTwo
		}
	} else if rangeType == "numbered" {
		noOfOccurrence := util.InterfaceToInt(rangeObject["numberOfOccurrences"])

		datePointer := startDateTime

		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)
		incre := 0
		remainderBreak := false

		for incre < noOfOccurrence {
			for dateWeekPointer.After(datePointer) && incre < noOfOccurrence {
				if remainderBreak {
					break
				}
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate
					if (iOrderQty - scheduledDailyRate) < 0 {
						//This is indicator for remainder
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						eventEndDateTime = datePointer.Add(time.Duration(shiftPeriod) * time.Hour)
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}
		if scheduledDailyRate < iOrderQty {
			preferenceLevel = ScheduleStatusPreferenceTwo
		}
	} else if rangeType == "noEnd" {
		datePointer := startDateTime

		weekEnding := 7 - int(startDateTime.Weekday())
		dateWeekPointer := startDateTime.AddDate(0, 0, weekEnding)

		remainderBreak := false

		incre := 0
		remainingDays := fullAllocation
		for remainingDays > 0 {
			if remainderBreak {
				break
			}
			remainingDays -= 1

			for dateWeekPointer.After(datePointer) && (incre <= fullAllocation) && (scheduledDailyRate < iOrderQty) {
				if contains(productionDaysWeek, datePointer.Weekday().String()) {
					if scheduledDailyRate > iOrderQty {
						return eventScheduled, getError("Schedule quantity exceeds the order quantity"), preferenceLevel
					}
					incre += 1
					scheduledDailyRate += iDailyRate

					if (iOrderQty - scheduledDailyRate) < 0 {
						//Reminder quantity
						remainderQty := iOrderQty - (scheduledDailyRate - iDailyRate)

						var eventEndDateTime time.Time
						if isEndTimeSpecified {
							eventEndDateTime = datePointer.AddDate(0, 0, 1)
						} else {
							eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						}

						eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, remainderQty, eventSourceId, machineId, orderStatusId)
						eventScheduled = append(eventScheduled, eventObject)
						remainderBreak = true
						break
					}

					var eventEndDateTime time.Time
					if isEndTimeSpecified {
						eventEndDateTime = datePointer.AddDate(0, 0, 1)
					} else {
						// eventEndDateTime = time.Date(datePointer.Year(), datePointer.Month(), datePointer.Day(), endHours, endMins, 0, datePointer.Nanosecond(), datePointer.Location())
						eventEndDateTime = datePointer.Add(time.Duration(shiftPeriod) * time.Hour)
					}

					eventObject := initSplitEvent(datePointer.Format(TimeLayout), eventEndDateTime.Format(TimeLayout), productionOrder, noOfExistingEvent+incre, iDailyRate, eventSourceId, machineId, orderStatusId)
					eventScheduled = append(eventScheduled, eventObject)
				}
				datePointer = datePointer.AddDate(0, 0, 1)
			}
			dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
			datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
		}

	} else {
		fmt.Println("Non supported range type")
	}

	return eventScheduled, nil, preferenceLevel
}

func getNumberOfTotalScheduledEvents(productionOrderId int, dbConnection *gorm.DB) int {
	conditionString := " JSON_EXTRACT(object_info, \"$.eventSourceId\") = '" + strconv.Itoa(productionOrderId) + "'"
	listOfObjects, err := GetConditionalObjects(dbConnection, ScheduledOrderEventTable, conditionString)
	if err != nil {
		return 0
	} else {
		return len(*listOfObjects)
	}
}

func (v *ProductionOrderService) SplitSchedule(projectId string, dbConnection *gorm.DB, object map[string]interface{}, orderStatusId int, requestObject map[string]interface{}, scheduledOrderTableName string) (*response.DetailedError, int, int) {
	iOrderQty := util.InterfaceToInt(object["orderQty"])
	balanceQty := util.InterfaceToInt(object["balance"])
	iDailyRate := util.InterfaceToInt(object["dailyRate"])
	machineId := util.InterfaceToInt(object["machineId"])
	remainder := math.Mod(float64(iOrderQty), float64(iDailyRate))
	productionOrder := util.InterfaceToString(object["prodOrder"])
	eventSourceId := util.InterfaceToInt(object["id"])
	remainingScheduleQty := util.InterfaceToInt(object["remainingScheduledQty"])

	preferenceLevel := ScheduleStatusPreferenceThree

	v.BaseService.Logger.Info("order details", zap.Any("order quantity:", iOrderQty), zap.Any("daily rate:", iDailyRate), zap.Any("machine id", machineId),
		zap.Any("balance quantity", balanceQty), zap.Any("remainder", remainder))
	// if balanceQty > 0 {
	// 	po.BaseService.Logger.Info("Balance is more than 0")
	// 	return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
	// }

	patternObject := requestObject["pattern"].(map[string]interface{})
	rangeObject := requestObject["range"].(map[string]interface{})
	interval := util.InterfaceToInt(patternObject["interval"])

	var startTime string
	isBasedOnShift := util.InterfaceToBool(requestObject["isBasedOnShift"])
	var scheduleArray []map[string]interface{}

	if scheduledOrderTableName == AssemblyScheduledOrderEventTable && isBasedOnShift {
		shiftTemplates := requestObject["shiftTemplates"]
		shiftTemplatesList := util.InterfaceToIntArray(shiftTemplates)
		labourManagementService := common.GetService("labour_management_module").ServiceInterface.(common.LabourManagementInterface)
		labourManagementInfo := labourManagementService.GetLabourManagementShiftTemplate(shiftTemplatesList)
		fullAllocationDays := iOrderQty / (iDailyRate * len(labourManagementInfo))

		for _, template := range labourManagementInfo {

			var scheduleArrayList []map[string]interface{}
			templateInfo := make(map[string]interface{})
			json.Unmarshal(template.ObjectInfo, &templateInfo)
			var shiftPeriod int

			startTime = util.GetCurrentDate() + "T" + util.InterfaceToString(templateInfo["shiftStartTime"]) + ":00.000Z"

			// if end time isn't specified
			isEndTimeSpecified := util.InterfaceToBool(requestObject["is24HoursSplitting"])
			var endTime string
			var endDateTime string
			if isEndTimeSpecified {
				endTime = util.GetCurrentDate() + "T" + util.InterfaceToString(templateInfo["shiftStartTime"]) + ":00.000Z"
				stringToDate := util.ConvertStringToDateTime(endTime)
				oneDayIncrement := stringToDate.DateTime.AddDate(0, 0, 1)
				// endTime = oneDayIncrement.Format("2006-01-02T15:04:05.000Z")

				endHour := strconv.Itoa(oneDayIncrement.Hour())
				endMin := strconv.Itoa(oneDayIncrement.Minute())

				endTime = endHour + ":" + endMin
			} else {

				parsedTime, _ := time.Parse("2006-01-02T15:04:05.000Z", startTime)
				shiftPeriod = util.InterfaceToInt(templateInfo["shiftPeriod"])
				newTime := parsedTime.Add(time.Duration(shiftPeriod) * time.Hour)
				endDateTime = newTime.Format("2006-01-02T15:04:05.000Z")
				endTime = newTime.Format("15:04")

			}

			if !util.CompareDateTime(startTime, endDateTime) {
				return getDetailedError(InvalidTimeRange, TimeRangeErrorDescription), preferenceLevel, 0
			}
			startDateTime := util.InterfaceToString(rangeObject["startDate"]) + "T" + util.InterfaceToString(templateInfo["shiftStartTime"]) + ":00.000Z"

			var err error

			noOfExistingEvents := getNumberOfTotalScheduledEvents(eventSourceId, dbConnection)

			if util.InterfaceToString(patternObject["type"]) == "weekly" {
				if remainingScheduleQty != iOrderQty {
					remainder = math.Mod(float64(remainingScheduleQty), float64(iDailyRate))
					scheduleArrayList, err, preferenceLevel = SplitWeeklyAssemblySchedule(remainingScheduleQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, endTime, rangeObject, patternObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified, shiftPeriod)
				} else {
					remainder = math.Mod(float64(iOrderQty), float64(iDailyRate))
					scheduleArrayList, err, preferenceLevel = SplitWeeklyAssemblySchedule(iOrderQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, endTime, rangeObject, patternObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified, shiftPeriod)
				}

				v.BaseService.Logger.Info("weekly scheduled array", zap.Any("scheduleArray", scheduleArrayList))
				if err != nil {
					v.BaseService.Logger.Info("Error in split weekly schedule")
					return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
				}

			} else if util.InterfaceToString(patternObject["type"]) == "daily" {
				if remainingScheduleQty != iOrderQty {
					remainder = math.Mod(float64(remainingScheduleQty), float64(iDailyRate))
					scheduleArrayList, err, preferenceLevel = SplitDailyAssemblySchedule(remainingScheduleQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, endTime, rangeObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified, shiftPeriod)
				} else {
					remainder = math.Mod(float64(iOrderQty), float64(iDailyRate))
					scheduleArrayList, err, preferenceLevel = SplitDailyAssemblySchedule(iOrderQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, endTime, rangeObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified, shiftPeriod)
				}

				v.BaseService.Logger.Info("daily scheduled array", zap.Any("scheduleArray", scheduleArrayList))
				if err != nil {
					v.BaseService.Logger.Info("Error in split daily schedule")
					return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
				}
			} else {
				return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
			}
			for _, item := range scheduleArrayList {
				scheduleArray = append(scheduleArray, item)
			}
		}
	} else {
		fullAllocationDays := iOrderQty / iDailyRate

		startTime = util.GetCurrentDate() + "T" + util.InterfaceToString(requestObject["startTime"]) + ":00.000Z"

		// if end time isn't specified
		isEndTimeSpecified := util.InterfaceToBool(requestObject["is24HoursSplitting"])
		var endTime string
		if isEndTimeSpecified {
			endTime = util.GetCurrentDate() + "T" + util.InterfaceToString(requestObject["startTime"]) + ":00.000Z"
			stringToDate := util.ConvertStringToDateTime(endTime)
			oneDayIncrement := stringToDate.DateTime.AddDate(0, 0, 1)
			endTime = oneDayIncrement.Format("2006-01-02T15:04:05.000Z")

			endHour := strconv.Itoa(oneDayIncrement.Hour())
			endMin := strconv.Itoa(oneDayIncrement.Minute())

			requestObject["endTime"] = endHour + ":" + endMin
		} else {
			endTime = util.GetCurrentDate() + "T" + util.InterfaceToString(requestObject["endTime"]) + ":00.000Z"
		}

		if !util.CompareDateTime(startTime, endTime) {
			return getDetailedError(InvalidTimeRange, TimeRangeErrorDescription), preferenceLevel, 0
		}
		startDateTime := util.InterfaceToString(rangeObject["startDate"]) + "T" + util.InterfaceToString(requestObject["startTime"]) + ":00.000Z"

		var err error

		noOfExistingEvents := getNumberOfTotalScheduledEvents(eventSourceId, dbConnection)

		if util.InterfaceToString(patternObject["type"]) == "weekly" {
			if remainingScheduleQty != iOrderQty {
				remainder = math.Mod(float64(remainingScheduleQty), float64(iDailyRate))
				scheduleArray, err, preferenceLevel = SplitWeeklySchedule(remainingScheduleQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, util.InterfaceToString(requestObject["endTime"]), rangeObject, patternObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified)
			} else {
				remainder = math.Mod(float64(iOrderQty), float64(iDailyRate))
				scheduleArray, err, preferenceLevel = SplitWeeklySchedule(iOrderQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, util.InterfaceToString(requestObject["endTime"]), rangeObject, patternObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified)
			}

			v.BaseService.Logger.Info("weekly scheduled array", zap.Any("scheduleArray", scheduleArray))
			if err != nil {
				v.BaseService.Logger.Info("Error in split weekly schedule")
				return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
			}

		} else if util.InterfaceToString(patternObject["type"]) == "daily" {
			if remainingScheduleQty != iOrderQty {
				remainder = math.Mod(float64(remainingScheduleQty), float64(iDailyRate))
				scheduleArray, err, preferenceLevel = SplitDailySchedule(remainingScheduleQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, util.InterfaceToString(requestObject["endTime"]), rangeObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified)
			} else {
				remainder = math.Mod(float64(iOrderQty), float64(iDailyRate))
				scheduleArray, err, preferenceLevel = SplitDailySchedule(iOrderQty, iDailyRate, fullAllocationDays, int(remainder), interval, startDateTime, util.InterfaceToString(requestObject["endTime"]), rangeObject, productionOrder, eventSourceId, machineId, orderStatusId, noOfExistingEvents, isEndTimeSpecified)
			}

			v.BaseService.Logger.Info("daily scheduled array", zap.Any("scheduleArray", scheduleArray))
			if err != nil {
				v.BaseService.Logger.Info("Error in split daily schedule")
				return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
			}
		} else {
			return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero), preferenceLevel, 0
		}
	}

	// time validation

	scheduledQty := 0
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

	for _, eventObject := range scheduleArray {
		v.BaseService.Logger.Info("creating event", zap.Any("event", eventObject))

		if scheduledOrderTableName == AssemblyScheduledOrderEventTable {
			// get default manpower
			var defaultManpower = machineService.GetAssemblyMachineDefaultManpower(machineId)
			v.BaseService.Logger.Info("assigning default manpower", zap.Int("default_manpower", defaultManpower))
			eventObject["plannedManpower"] = defaultManpower
		}
		jsonObject, _ := json.Marshal(eventObject)
		scheduledQty += util.InterfaceToInt(eventObject["scheduledQty"])
		generalObject := component.GeneralObject{ObjectInfo: jsonObject}
		err, recordId := Create(dbConnection, scheduledOrderTableName, generalObject)
		if err != nil {
			v.BaseService.Logger.Error("error creating schedule", zap.Any("error", err.Error()))
			continue
		} else {
			// now create the history with event status
			if scheduledOrderTableName == AssemblyScheduledOrderEventTable {
				if eventStatusId, ok := eventObject["eventStatus"]; ok {
					v.BaseService.Logger.Info("creating shift line history", zap.Int("event_id", recordId))
					err := v.ViewManager.CreateLabourManagementShiftLinesHistory(recordId, util.InterfaceToInt(eventStatusId))
					if err != nil {
						v.BaseService.Logger.Error("error creating shift lines history", zap.Error(err))
					}
				}
			}
		}
		v.createSystemNotification(projectId, util.InterfaceToString(eventObject["name"]), util.InterfaceToString(eventObject["name"])+" is created", recordId)
	}

	return nil, preferenceLevel, scheduledQty
}

func (v *ProductionOrderService) CreateSchedule(dbConnection *gorm.DB, object map[string]interface{}, orderStatusId int) *response.DetailedError {
	iOrderQty := util.InterfaceToInt(object["orderQty"])
	balanceQty := util.InterfaceToInt(object["balance"])
	iDailyRate := util.InterfaceToInt(object["dailyRate"])
	machineId := util.InterfaceToInt(object["machineId"])
	remainder := math.Mod(float64(iOrderQty), float64(iDailyRate))
	productionOrder := util.InterfaceToString(object["prodOrder"])
	eventSourceId := util.InterfaceToInt(object["id"])

	v.BaseService.Logger.Info("order details", zap.Any("order qty", iOrderQty), zap.Any("daily rate:", iDailyRate), zap.Any("machine id", machineId))
	if balanceQty > 0 {
		return getDetailedError(PartiallyScheduled, BalanceQuantityGreaterThanZero)
	}
	fullAllocationDays := iOrderQty / iDailyRate
	fmt.Println("fullAllocationDays:", fullAllocationDays)
	for i := 0; i < fullAllocationDays; i++ {
		// create a event, and assigned to machine
		eventObject := initEvent(productionOrder, i, iDailyRate, eventSourceId, machineId, orderStatusId)
		v.BaseService.Logger.Info("creating event", zap.Any("event", eventObject))
		generalObject := component.GeneralObject{ObjectInfo: eventObject}
		err, _ := Create(dbConnection, ScheduledOrderEventTable, generalObject)
		if err != nil {
			v.BaseService.Logger.Error("error creating schedule", zap.Any("error", err.Error()))
		}
	}
	if remainder != 0 {

		eventObject := initEvent(productionOrder, fullAllocationDays+1, remainder, eventSourceId, machineId, orderStatusId)
		v.BaseService.Logger.Info("creating event", zap.Any("event", eventObject))
		generalObject := component.GeneralObject{ObjectInfo: eventObject}
		err, _ := Create(dbConnection, ScheduledOrderEventTable, generalObject)
		if err != nil {
			v.BaseService.Logger.Error("error creating schedule", zap.Any("error", err.Error()))
		}
		// now do the remaining

	}
	return nil
}

func (v *ProductionOrderService) handleUpdateSchedule(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	scheduleUpdateRequest := ScheduleRequest{}
	if err := ctx.ShouldBindBodyWith(&scheduleUpdateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !scheduleValidation(scheduleUpdateRequest.StartDate, scheduleUpdateRequest.EndDate) {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(InvalidScheduleEventStatusError).Error(),
				Description: "This schedule has invalid schedule start date and end date",
			})
		return
	}

	err, scheduledEventObject := Get(dbConnection, ScheduledOrderEventTable, recordId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(ObjectNotFound).Error(),
				Description: "Invalid event, this event is no longer available in the system, are you sending request out of system environment?",
			})
		return
	}
	if !common.ValidateObjectStatus(scheduledEventObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

	//preferenceOrderStatusId3 := po.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	currentEventStatus := scheduledEventInfo.EventStatus
	currentEventSourceId := scheduledEventInfo.EventSourceId
	currentEventName := scheduledEventInfo.Name
	completedQuantity := util.InterfaceToInt(scheduleUpdateRequest.CompletedQty)
	rejectedQuantity := util.InterfaceToInt(scheduleUpdateRequest.RejectedQty)
	scheduledQuantity := util.InterfaceToInt(scheduleUpdateRequest.ScheduledQty)
	existingScheduledQty := scheduledEventInfo.ScheduledQty
	if completedQuantity > 0 {
		// we are trying to update the complete quantity , so check the order status is proudction stopped
		orderStatusIdPreferenceLevel6 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSix)

		if currentEventStatus != orderStatusIdPreferenceLevel6 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "You are trying to update the completed quantity of the order  which the order is not stopped or started, please stop or process the order from HMI and proceed",
				})
			return
		}
	}
	v.BaseService.Logger.Info("received event source id used :", zap.Any("event_source_id", currentEventSourceId))
	productionOrderInfo := ProductionOrderInfo{}
	err, generalObject := Get(dbConnection, ProductionOrderMasterTable, currentEventSourceId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Production Order").Error(),
				Description: "Relationship missing between scheduled orders and parent order",
			})

		return
	}
	json.Unmarshal(generalObject.ObjectInfo, &productionOrderInfo)
	if scheduledQuantity > 0 && existingScheduledQty != scheduledQuantity {
		orderStatusIdPreferenceLevel4 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
		if currentEventStatus == orderStatusIdPreferenceLevel4 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "This schedule is already confirmed, please make it un-confirmed to amend it",
				})

			return
		}
		var deltaVariation int
		if existingScheduledQty > scheduledQuantity {
			deltaVariation = existingScheduledQty - scheduledQuantity
			productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty + deltaVariation
		} else {
			deltaVariation = scheduledQuantity - existingScheduledQty
			productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty - deltaVariation
		}

		if productionOrderInfo.RemainingScheduledQty < deltaVariation {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError("Invalid Scheduled Quantity").Error(),
					Description: "Remaining scheduled quantity is lower than requested one, check the quantity again",
				})
			return
		}
		// we reached the remaining scheduled quantity 0, so this means, we already scheduled
		orderStatusIdPreferenceLevel3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
		orderStatusIdPreferenceLevel2 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceTwo)
		if productionOrderInfo.RemainingScheduledQty == 0 {
			productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel3
		} else {
			productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel2
		}
		scheduledEventInfo.ScheduledQty = scheduledQuantity

		productionOrderUpdate := make(map[string]interface{})
		rawProductionOrderInfo, _ := json.Marshal(productionOrderInfo)
		productionOrderUpdate["object_info"] = rawProductionOrderInfo
		err = Update(dbConnection, ProductionOrderMasterTable, currentEventSourceId, productionOrderUpdate)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
			return
		}
		scheduledEventInfo.EventStatus = orderStatusIdPreferenceLevel3
	}

	if completedQuantity > 0 {
		scheduledEventInfo.CompletedQty = completedQuantity
	}
	if rejectedQuantity > 0 {
		scheduledEventInfo.RejectedQty = scheduledEventInfo.RejectedQty + rejectedQuantity

		//Can't add rejected quantity for completed schedule
		orderStatusIdPreferenceLevel7 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)
		if currentEventStatus == orderStatusIdPreferenceLevel7 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "You are trying to update the completed schedule",
				})
			return
		}
	}

	if scheduleUpdateRequest.Name != "" {
		scheduledEventInfo.Name = scheduleUpdateRequest.Name
	}
	if scheduleUpdateRequest.StartDate != "" {
		scheduledEventInfo.StartDate = util.ConvertSingaporeTimeToUTC(scheduleUpdateRequest.StartDate)
	}
	if scheduleUpdateRequest.EndDate != "" {
		scheduledEventInfo.EndDate = util.ConvertSingaporeTimeToUTC(scheduleUpdateRequest.EndDate)
	}

	if scheduleUpdateRequest.MouldId != 0 {
		scheduledEventInfo.MouldId = scheduleUpdateRequest.MouldId
		scheduledEventInfo.EnableCustomCavity = scheduleUpdateRequest.EnableCustomCavity
		scheduledEventInfo.CustomCavity = scheduleUpdateRequest.CustomCavity
	}

	if scheduleUpdateRequest.MouldUp != "" {
		scheduledEventInfo.MouldUp = scheduleUpdateRequest.MouldUp
	}

	if scheduleUpdateRequest.MouldDown != "" {
		scheduledEventInfo.MouldDown = scheduleUpdateRequest.MouldDown
	}

	userId := common.GetUserId(ctx)

	preferenceOrderStatusId3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	if scheduledEventInfo.EventStatus > preferenceOrderStatusId3 {
		fmt.Println("Before: ", scheduleUpdateRequest.MouldId)
		if scheduleUpdateRequest.MouldUp != "" || scheduleUpdateRequest.MouldDown != "" || scheduleUpdateRequest.MouldId != 0 {
			fmt.Println(scheduleUpdateRequest.MouldId)
			err = Update(dbConnection, ScheduledOrderEventTable, recordId, scheduledEventInfo.DatabaseSerialize(userId))
			if err != nil {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
				return
			}
		} else {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSchedulePosition), InvalidScheduleStatus, "Schedule order passed the update stage")
			return
		}

	}

	err = Update(dbConnection, ScheduledOrderEventTable, recordId, scheduledEventInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
		return
	}

	notificationHeader := currentEventName + " is updated"
	notificationDescription := "The scheduler event is modified, fields are updated. Check the timeline for further details"
	v.createSystemNotification(projectId, notificationHeader, notificationDescription, recordId)
}

func (v *ProductionOrderService) handleUpdateToolingSchedule(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	scheduleUpdateRequest := make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&scheduleUpdateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err, scheduledEventObject := Get(dbConnection, ToolingScheduledOrderEventTable, recordId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(ObjectNotFound).Error(),
				Description: "Invalid event, this event is no longer available in the system, are you sending request out of system environment?",
			})
		return
	}

	if !scheduleValidation(util.InterfaceToString(scheduleUpdateRequest["startDate"]), util.InterfaceToString(scheduleUpdateRequest["endDate"])) {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(InvalidScheduleEventStatusError).Error(),
				Description: "This schedule has invalid schedule start date and end date",
			})
		return
	}

	updatingData := make(map[string]interface{})

	scheduledOrderEvent := ToolingScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledEventInfo := scheduledOrderEvent.getToolingScheduledOrderEventInfo()

	if val, ok := scheduleUpdateRequest["setupTime"]; ok {
		setupTime := util.InterfaceToString(val)
		scheduledEventInfo.SetupTime = setupTime
		scheduledEventInfo.CanSetupTime = false
	} else {
		scheduledEventInfo.Name = util.InterfaceToString(scheduleUpdateRequest["name"])

		requestedStartDate := util.InterfaceToString(scheduleUpdateRequest["startDate"])

		if scheduledEventInfo.StartDate != requestedStartDate {
			_, partOrder := Get(dbConnection, ToolingPartMasterTable, util.InterfaceToInt(scheduleUpdateRequest["partId"]))
			partOrderInfo := make(map[string]interface{})
			json.Unmarshal(partOrder.ObjectInfo, &partOrderInfo)

			day := util.InterfaceToInt(partOrderInfo["day"])
			hour := util.InterfaceToInt(partOrderInfo["hour"])
			minute := util.InterfaceToInt(partOrderInfo["minute"])

			durationSec := day*86400 + hour*3600 + minute*60

			startTime, _ := time.Parse(TimeLayout, requestedStartDate)
			endTime := startTime.Add(time.Second * time.Duration(durationSec))

			endTimeStr := endTime.Format(TimeLayout)

			scheduledEventInfo.StartDate = util.ConvertSingaporeTimeToUTC(requestedStartDate)
			scheduledEventInfo.EndDate = util.ConvertSingaporeTimeToUTC(endTimeStr)
		}
		scheduledEventInfo.EventSourceId = util.InterfaceToInt(scheduleUpdateRequest["eventSourceId"])
		scheduledEventInfo.PartId = util.InterfaceToInt(scheduleUpdateRequest["partId"])
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(scheduledEventInfo.Serialize(), ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, recordId, updatingData)
}

func (v *ProductionOrderService) handleUpdateAssemblySchedule(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	scheduleUpdateRequest := ScheduleRequest{}
	if err := ctx.ShouldBindBodyWith(&scheduleUpdateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !scheduleValidation(scheduleUpdateRequest.StartDate, scheduleUpdateRequest.EndDate) {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(InvalidScheduleEventStatusError).Error(),
				Description: "This schedule has invalid schedule start date and end date",
			})
		return
	}

	err, scheduledEventObject := Get(dbConnection, AssemblyScheduledOrderEventTable, recordId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(ObjectNotFound).Error(),
				Description: "Invalid event, this event is no longer available in the system, are you sending request out of system environment?",
			})
		return
	}
	if !common.ValidateObjectStatus(scheduledEventObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	scheduledOrderEvent := ScheduledOrderEvent{ObjectInfo: scheduledEventObject.ObjectInfo}
	scheduledEventInfo := scheduledOrderEvent.getScheduledOrderEventInfo()

	currentEventStatus := scheduledEventInfo.EventStatus
	currentEventSourceId := scheduledEventInfo.EventSourceId
	currentEventName := scheduledEventInfo.Name
	completedQuantity := util.InterfaceToInt(scheduleUpdateRequest.CompletedQty)
	rejectedQuantity := util.InterfaceToInt(scheduleUpdateRequest.RejectedQty)
	scheduledQuantity := util.InterfaceToInt(scheduleUpdateRequest.ScheduledQty)
	existingScheduledQty := scheduledEventInfo.ScheduledQty
	if completedQuantity > 0 {
		// we are trying to update the complete quantity , so check the order status is proudction stopped
		orderStatusIdPreferenceLevel6 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSix)

		if currentEventStatus != orderStatusIdPreferenceLevel6 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "You are trying to update the completed quantity of the order  which the order is not stopped or started, please stop or process the order from HMI and proceed",
				})
			return
		}
	}
	v.BaseService.Logger.Info("received event source id used :", zap.Any("event_source_id", currentEventSourceId))
	productionOrderInfo := ProductionOrderInfo{}
	err, generalObject := Get(dbConnection, AssemblyProductionOrderTable, currentEventSourceId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Production Order").Error(),
				Description: "Relationship missing between scheduled orders and parent order",
			})

		return
	}
	json.Unmarshal(generalObject.ObjectInfo, &productionOrderInfo)
	if scheduledQuantity > 0 && existingScheduledQty != scheduledQuantity {
		orderStatusIdPreferenceLevel4 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceFour)
		if currentEventStatus == orderStatusIdPreferenceLevel4 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "This schedule is already confirmed, please make it un-confirmed to amend it",
				})

			return
		}
		var deltaVariation int
		if existingScheduledQty > scheduledQuantity {
			deltaVariation = existingScheduledQty - scheduledQuantity
			productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty + deltaVariation
		} else {
			deltaVariation = scheduledQuantity - existingScheduledQty
			productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty - deltaVariation
		}

		if productionOrderInfo.RemainingScheduledQty < deltaVariation {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError("Invalid Scheduled Quantity").Error(),
					Description: "Remaining scheduled quantity is lower than requested one, check the quantity again",
				})
			return
		}
		// we reached the remaining scheduled quantity 0, so this means, we already scheduled
		orderStatusIdPreferenceLevel3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
		orderStatusIdPreferenceLevel2 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceTwo)
		if productionOrderInfo.RemainingScheduledQty == 0 {
			productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel3
		} else {
			productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel2
		}
		scheduledEventInfo.ScheduledQty = scheduledQuantity

		productionOrderUpdate := make(map[string]interface{})
		rawProductionOrderInfo, _ := json.Marshal(productionOrderInfo)
		productionOrderUpdate["object_info"] = rawProductionOrderInfo
		err = Update(dbConnection, AssemblyProductionOrderTable, currentEventSourceId, productionOrderUpdate)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
			return
		}
		scheduledEventInfo.EventStatus = orderStatusIdPreferenceLevel3
	}

	if completedQuantity > 0 {
		scheduledEventInfo.CompletedQty = completedQuantity
	}
	if rejectedQuantity > 0 {
		scheduledEventInfo.RejectedQty = scheduledEventInfo.RejectedQty + rejectedQuantity

		//Can't add rejected quantity for completed schedule
		orderStatusIdPreferenceLevel7 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceSeven)
		if currentEventStatus == orderStatusIdPreferenceLevel7 {
			response.DispatchDetailedError(ctx, FieldValidationFailed,
				&response.DetailedError{
					Header:      getError(InvalidScheduleEventStatusError).Error(),
					Description: "You are trying to update the completed schedule",
				})
			return
		}
	}
	if scheduleUpdateRequest.Name != "" {
		scheduledEventInfo.Name = scheduleUpdateRequest.Name
	}
	if scheduleUpdateRequest.StartDate != "" {
		scheduledEventInfo.StartDate = util.ConvertSingaporeTimeToUTC(scheduleUpdateRequest.StartDate)
	}
	if scheduleUpdateRequest.EndDate != "" {
		scheduledEventInfo.EndDate = util.ConvertSingaporeTimeToUTC(scheduleUpdateRequest.EndDate)
	}
	scheduledEventInfo.MouldId = scheduleUpdateRequest.MouldId
	scheduledEventInfo.EnableCustomCavity = scheduleUpdateRequest.EnableCustomCavity
	scheduledEventInfo.CustomCavity = scheduleUpdateRequest.CustomCavity

	if scheduleUpdateRequest.MouldUp != "" {
		scheduledEventInfo.MouldUp = scheduleUpdateRequest.MouldUp
	}

	if scheduleUpdateRequest.MouldDown != "" {
		scheduledEventInfo.MouldDown = scheduleUpdateRequest.MouldDown
	}

	userId := common.GetUserId(ctx)
	preferenceOrderStatusId3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	if scheduledEventInfo.EventStatus > preferenceOrderStatusId3 {
		fmt.Println("Before: ", scheduleUpdateRequest.MouldId)
		if scheduleUpdateRequest.MouldUp != "" || scheduleUpdateRequest.MouldDown != "" {
			fmt.Println(scheduleUpdateRequest.MouldId)
			err = Update(dbConnection, AssemblyScheduledOrderEventTable, recordId, scheduledEventInfo.DatabaseSerialize(userId))
			if err != nil {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
				return
			}
		} else {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError(InvalidSchedulePosition), InvalidScheduleStatus, "Schedule order passed the update stage")
			return
		}

	}

	err = Update(dbConnection, AssemblyScheduledOrderEventTable, recordId, scheduledEventInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Update Event Information Failed"), ErrorUpdatingObjectInformation, "Updating event information has failed due to ["+err.Error()+"]")
		return
	}

	notificationHeader := currentEventName + " is updated"
	notificationDescription := "The scheduler event is modified, fields are updated. Check the timeline for further details"
	v.createSystemNotification(projectId, notificationHeader, notificationDescription, recordId)
}

func (v *ProductionOrderService) handleCreateSchedule(ctx *gin.Context) {
	scheduleRequest := ScheduleRequest{}
	projectId := util.GetProjectId(ctx)
	if err := ctx.ShouldBindBodyWith(&scheduleRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !scheduleValidation(scheduleRequest.StartDate, scheduleRequest.EndDate) {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(InvalidScheduleEventStatusError).Error(),
				Description: "This schedule has invalid schedule start date and end date",
			})
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	productionOrderId := scheduleRequest.EventSourceId
	err, orderObject := Get(dbConnection, ProductionOrderMasterTable, productionOrderId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Production Order").Error(),
				Description: "Relationship missing between scheduled orders and parent order",
			})

		return
	}
	productionOrderInfo := ProductionOrderInfo{}
	json.Unmarshal(orderObject.ObjectInfo, &productionOrderInfo)
	// we need to call the machine services
	orderStatusIdPreferenceLevel3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	scheduleRequest.StartDate = util.ConvertSingaporeTimeToUTC(scheduleRequest.StartDate)
	scheduleRequest.EndDate = util.ConvertSingaporeTimeToUTC(scheduleRequest.EndDate)

	rawEvent := initEventWithDates(productionOrderInfo.ProdOrder, scheduleRequest, productionOrderInfo.MachineId, orderStatusIdPreferenceLevel3)

	updatingData := make(map[string]interface{})

	if productionOrderInfo.RemainingScheduledQty == 0 {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Exceeding Remaining Scheduled Quantity").Error(),
				Description: "Can not create new schedule as it is already filled, remove one to create new",
			})

		return
	}

	scheduledQty := util.InterfaceToInt(scheduleRequest.ScheduledQty)
	if productionOrderInfo.RemainingScheduledQty < scheduledQty {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Scheduled Quantity").Error(),
				Description: "Remaining scheduled quantity is lower than requested one, check the quantity again, this is mainly that you are scheduling greater than water mark level",
			})
		return
	} else {
		productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty - scheduledQty
	}

	// we reached the remaining scheduled quantity 0, so this means, we already scheduled
	if productionOrderInfo.RemainingScheduledQty == 0 {
		productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel3
	}

	serializedProductionOrder, _ := json.Marshal(productionOrderInfo)
	updatingData["object_info"] = serializedProductionOrder
	Update(dbConnection, ProductionOrderMasterTable, productionOrderId, updatingData)

	scheduledOrderEvent := component.GeneralObject{ObjectInfo: rawEvent}
	_, eventId := Create(dbConnection, ScheduledOrderEventTable, scheduledOrderEvent)

	notificationHeader := scheduleRequest.Name
	notificationDescription := "New order is scheduled with id [" + productionOrderInfo.ProdOrder + "] click view to see more details"
	v.createSystemNotification(projectId, notificationHeader, notificationDescription, eventId)
	// create a system notification

	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "New schedule is successfully created",
	})
}

func (v *ProductionOrderService) handleCreateAssemblySchedule(ctx *gin.Context) {
	scheduleRequest := ScheduleRequest{}
	projectId := util.GetProjectId(ctx)
	if err := ctx.ShouldBindBodyWith(&scheduleRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !scheduleValidation(scheduleRequest.StartDate, scheduleRequest.EndDate) {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError(InvalidScheduleEventStatusError).Error(),
				Description: "This schedule has invalid schedule start date and end date",
			})
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	productionOrderId := scheduleRequest.EventSourceId
	err, orderObject := Get(dbConnection, AssemblyProductionOrderTable, productionOrderId)
	if err != nil {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Production Order").Error(),
				Description: "Relationship missing between scheduled orders and parent order",
			})

		return
	}
	productionOrderInfo := ProductionOrderInfo{}
	json.Unmarshal(orderObject.ObjectInfo, &productionOrderInfo)
	// we need to call the machine services
	orderStatusIdPreferenceLevel3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	scheduleRequest.StartDate = util.ConvertSingaporeTimeToUTC(scheduleRequest.StartDate)
	scheduleRequest.EndDate = util.ConvertSingaporeTimeToUTC(scheduleRequest.EndDate)

	rawEvent := initEventWithDates(productionOrderInfo.ProdOrder, scheduleRequest, productionOrderInfo.MachineId, orderStatusIdPreferenceLevel3)

	updatingData := make(map[string]interface{})

	if productionOrderInfo.RemainingScheduledQty == 0 {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Exceeding Remaining Scheduled Quantity").Error(),
				Description: "Can not create new schedule as it is already filled, remove one to create new",
			})

		return
	}

	scheduledQty := util.InterfaceToInt(scheduleRequest.ScheduledQty)
	if productionOrderInfo.RemainingScheduledQty < scheduledQty {
		response.DispatchDetailedError(ctx, FieldValidationFailed,
			&response.DetailedError{
				Header:      getError("Invalid Scheduled Quantity").Error(),
				Description: "Remaining scheduled quantity is lower than requested one, check the quantity again, this is mainly that you are scheduling greater than water mark level",
			})
		return
	} else {
		productionOrderInfo.RemainingScheduledQty = productionOrderInfo.RemainingScheduledQty - scheduledQty
	}

	// we reached the remaining scheduled quantity 0, so this means, we already scheduled
	if productionOrderInfo.RemainingScheduledQty == 0 {
		productionOrderInfo.OrderStatus = orderStatusIdPreferenceLevel3
	}

	serializedProductionOrder, _ := json.Marshal(productionOrderInfo)
	updatingData["object_info"] = serializedProductionOrder
	Update(dbConnection, AssemblyProductionOrderTable, productionOrderId, updatingData)

	scheduledOrderEvent := component.GeneralObject{ObjectInfo: rawEvent}
	_, eventId := Create(dbConnection, AssemblyScheduledOrderEventTable, scheduledOrderEvent)

	notificationHeader := scheduleRequest.Name
	notificationDescription := "New order is scheduled with id [" + productionOrderInfo.ProdOrder + "] click view to see more details"
	v.createSystemNotification(projectId, notificationHeader, notificationDescription, eventId)
	// create a system notification

	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "New schedule is successfully created",
	})
}

func (v *ProductionOrderService) SplitToolingSchedule(toolingProductionOrder map[string]interface{}, dbConnection *gorm.DB, productionId int, requestObject map[string]interface{}, duration string) []map[string]interface{} {
	listOfSplitOrders := make([]map[string]interface{}, 0)

	rangeObject := requestObject["range"].(map[string]interface{})

	startDateTime := util.InterfaceToString(rangeObject["startDate"]) + "T" + util.InterfaceToString(requestObject["startTime"]) + ":00.000Z"
	//startDateTime := util.InterfaceToString(rangeObject["startDate"])
	patternObject := requestObject["pattern"].(map[string]interface{})
	interval := util.InterfaceToInt(patternObject["interval"])
	typeOdSplit := util.InterfaceToString(patternObject["type"])

	startTime, _ := time.Parse(TimeLayout, startDateTime)

	productionOrder := util.InterfaceToString(toolingProductionOrder["prodOrder"])
	machineId := util.InterfaceToInt(toolingProductionOrder["machineId"])
	orderStatusIdPreferenceLevel3 := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)

	listOfPartNo := util.InterfaceToIntArray(toolingProductionOrder["partNo"])
	var inQuery string
	inQuery = ""
	if len(listOfPartNo) == 0 {
		inQuery = "-1"
	} else {
		for index, partNo := range listOfPartNo {
			if index == len(listOfPartNo)-1 {
				inQuery += strconv.Itoa(partNo)
			} else {
				inQuery += strconv.Itoa(partNo) + ","
			}

		}
	}
	conditionString := " id IN (" + inQuery + ")"

	generalObjects, _ := GetConditionalObjects(dbConnection, ToolingPartMasterTable, conditionString)

	if typeOdSplit == "weekly" {
		dayWeeks := util.InterfaceToStringArray(patternObject["daysOfWeek"])
		listOfSplitOrders = v.WeeklySplitToolingSchedule(generalObjects, startTime, productionOrder, productionId, machineId, orderStatusIdPreferenceLevel3, interval, dayWeeks, duration)
	} else {
		listOfSplitOrders = v.DailySplitToolingSchedule(generalObjects, startTime, productionOrder, productionId, machineId, orderStatusIdPreferenceLevel3, interval, duration)
	}

	return listOfSplitOrders

}

func (v *ProductionOrderService) DailySplitToolingSchedule(generalObjects *[]component.GeneralObject,
	startTime time.Time,
	productionOrder string,
	productionId int,
	machineId int,
	orderStatusIdPreferenceLevel3 int,
	interval int,
	duration string,
) []map[string]interface{} {
	listOfSplitOrders := make([]map[string]interface{}, 0)
	splitDuration := strings.Split(duration, ":")
	durationInSec := util.InterfaceToInt(splitDuration[0])*3600 + util.InterfaceToInt(splitDuration[1])*60
	splitIndex := 0
	for _, partMaster := range *generalObjects {
		toolingPartInfo := make(map[string]interface{})
		json.Unmarshal(partMaster.ObjectInfo, &toolingPartInfo)

		day := util.InterfaceToInt(toolingPartInfo["day"])
		hour := util.InterfaceToInt(toolingPartInfo["hour"])
		minute := util.InterfaceToInt(toolingPartInfo["minute"])

		partDurationSec := day*86400 + hour*3600 + minute*60

		for partDurationSec > 0 {
			var endTimeDuration int
			if partDurationSec > durationInSec {
				endTimeDuration = durationInSec
			} else {
				endTimeDuration = partDurationSec
			}
			endTime := startTime.Add(time.Second * time.Duration(endTimeDuration))

			partDurationSec = partDurationSec - durationInSec

			startTimeStr := startTime.Format(TimeLayout)
			endTimeStr := endTime.Format(TimeLayout)
			splitObject := initToolingSplitEvent(startTimeStr, endTimeStr, productionOrder, splitIndex, productionId, machineId, orderStatusIdPreferenceLevel3, partMaster.Id)

			listOfSplitOrders = append(listOfSplitOrders, splitObject)

			startTime = startTime.AddDate(0, 0, 1*interval)
			splitIndex += 1
		}

	}
	return listOfSplitOrders
}

func (v *ProductionOrderService) WeeklySplitToolingSchedule(generalObjects *[]component.GeneralObject,
	startTime time.Time,
	productionOrder string,
	productionId int,
	machineId int,
	orderStatusIdPreferenceLevel3 int,
	interval int,
	productionDaysWeek []string,
	duration string,
) []map[string]interface{} {
	datePointer := startTime
	listOfSplitOrders := make([]map[string]interface{}, 0)
	weekEnding := 7 - int(startTime.Weekday())
	dateWeekPointer := startTime.AddDate(0, 0, weekEnding)

	partInfoArray := *generalObjects

	splitDuration := strings.Split(duration, ":")
	durationInSec := util.InterfaceToInt(splitDuration[0])*3600 + util.InterfaceToInt(splitDuration[1])*60
	fmt.Println("partInfoArray: ", len(partInfoArray))
	splitIndex := 0
	if generalObjects != nil {
		for _, partMasterObj := range partInfoArray {
			toolingPartInfo := make(map[string]interface{})
			json.Unmarshal(partMasterObj.ObjectInfo, &toolingPartInfo)

			day := util.InterfaceToInt(toolingPartInfo["day"])
			hour := util.InterfaceToInt(toolingPartInfo["hour"])
			minute := util.InterfaceToInt(toolingPartInfo["minute"])

			partDurationSec := day*86400 + hour*3600 + minute*60

			for partDurationSec > 0 {
				for dateWeekPointer.After(datePointer) {
					if contains(productionDaysWeek, datePointer.Weekday().String()) && partDurationSec > 0 {
						var endTimeDuration int
						if partDurationSec > durationInSec {
							endTimeDuration = durationInSec
						} else {
							endTimeDuration = partDurationSec
						}
						endTime := datePointer.Add(time.Second * time.Duration(endTimeDuration))

						partDurationSec = partDurationSec - durationInSec

						startTimeStr := datePointer.Format(TimeLayout)
						endTimeStr := endTime.Format(TimeLayout)

						splitObject := initToolingSplitEvent(startTimeStr, endTimeStr, productionOrder, splitIndex, productionId, machineId, orderStatusIdPreferenceLevel3, partMasterObj.Id)

						listOfSplitOrders = append(listOfSplitOrders, splitObject)
						splitIndex += 1
						startTime = datePointer.AddDate(0, 0, 1)

					}
					datePointer = datePointer.AddDate(0, 0, 1)
				}
				dateWeekPointer = dateWeekPointer.AddDate(0, 0, 7*interval)
				datePointer = datePointer.AddDate(0, 0, 7*(interval-1))
			}

		}

	}

	return listOfSplitOrders
}

func IsValidToolingOrder(toolingRequest map[string]interface{}, generalObjects *[]component.GeneralObject) bool {
	// This function checks total duration is exceed the by part master durations
	// inputs
	//		tooling order master info : toolingRequest
	//		List of tooling part master : generalObjects
	validity := true

	day := util.InterfaceToInt(toolingRequest["day"])
	hour := util.InterfaceToInt(toolingRequest["hour"])
	minute := util.InterfaceToInt(toolingRequest["minute"])

	totalDurationSec := day*86400 + hour*3600 + minute*60

	partTotalDuration := 0
	for _, partMaster := range *generalObjects {
		toolingPartInfo := make(map[string]interface{})
		json.Unmarshal(partMaster.ObjectInfo, &toolingPartInfo)

		dayPart := util.InterfaceToInt(toolingPartInfo["day"])
		hourPart := util.InterfaceToInt(toolingPartInfo["hour"])
		minutePart := util.InterfaceToInt(toolingPartInfo["minute"])

		durationSec := dayPart*86400 + hourPart*3600 + minutePart*60

		partTotalDuration += durationSec

	}

	if totalDurationSec < partTotalDuration {
		return true
	}

	return validity
}

func (v *ProductionOrderService) handleCreateToolingSchedule(ctx *gin.Context) {
	scheduleRequest := make(map[string]interface{})
	projectId := util.GetProjectId(ctx)
	if err := ctx.ShouldBindBodyWith(&scheduleRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	// Based on event Source id we can get tooling order event
	eventSourceId := util.InterfaceToInt(scheduleRequest["eventSourceId"])
	_, toolingProductionOrder := Get(dbConnection, ToolingOrderMasterTable, eventSourceId)
	toolingOrderInfo := make(map[string]interface{})
	json.Unmarshal(toolingProductionOrder.ObjectInfo, &toolingOrderInfo)
	productionOrder := util.InterfaceToString(toolingOrderInfo["prodOrder"])
	partInfoInOrder := util.InterfaceToIntArray(toolingOrderInfo["partNo"])

	// Find existing schedule orders to get part ids
	conditionString := " object_info ->> '$.eventSourceId' = " + strconv.Itoa(eventSourceId) + " AND object_info ->> '$.objectStatus' = 'Active'"
	existingScheduledOrders, _ := GetConditionalObjects(dbConnection, ToolingScheduledOrderEventTable, conditionString)

	scheduleStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceThree)
	partiallyScheduleStatusId := v.getOrderStatusId(dbConnection, ScheduleStatusPreferenceTwo)
	// Check all the part master are scheduled
	if len(partInfoInOrder) == len(*existingScheduledOrders)+1 {
		toolingOrderInfo["orderStatus"] = scheduleStatusId
	} else {
		toolingOrderInfo["orderStatus"] = partiallyScheduleStatusId
	}

	// check the production order status
	// check total duration is exceeded
	requestedPartId := util.InterfaceToInt(scheduleRequest["partId"])
	_, requestedPartMaster := Get(dbConnection, ToolingPartMasterTable, requestedPartId)

	requestedPartMasterInfo := make(map[string]interface{})
	json.Unmarshal(requestedPartMaster.ObjectInfo, &requestedPartMasterInfo)

	day := util.InterfaceToInt(requestedPartMasterInfo["day"])
	hour := util.InterfaceToInt(requestedPartMasterInfo["hour"])
	minute := util.InterfaceToInt(requestedPartMasterInfo["minute"])

	requestedDurationSec := day*86400 + hour*3600 + minute*60

	requestedStartTime := util.ConvertSingaporeTimeToUTC(util.InterfaceToString(scheduleRequest["startDate"]))
	scheduleRequest["startDate"] = requestedStartTime
	startTime, _ := time.Parse(TimeLayout, requestedStartTime)
	endTime := startTime.Add(time.Second * time.Duration(requestedDurationSec))
	endTimeStr := endTime.Format(TimeLayout)

	scheduleRequest["endDate"] = endTimeStr
	scheduleRequest["partId"] = requestedPartId

	toolingScheduleOrder := initToolingEventWithDates(productionOrder, scheduleRequest, scheduleStatusId)

	scheduledOrderEvent := component.GeneralObject{ObjectInfo: toolingScheduleOrder}
	Create(dbConnection, ToolingScheduledOrderEventTable, scheduledOrderEvent)

	finalUpdatingData := make(map[string]interface{})
	serializedProductionOrder, _ := json.Marshal(toolingOrderInfo)
	finalUpdatingData["object_info"] = serializedProductionOrder

	Update(dbConnection, ToolingOrderMasterTable, toolingProductionOrder.Id, finalUpdatingData)

	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "New schedule is successfully created",
	})
}

func isValidEndDate(endDate string, startDate string, startTime string, generalObjects *[]component.GeneralObject) bool {
	isValid := true
	startDateTime := startDate + "T" + startTime + ":00.000Z"
	endDateTime := endDate + "T00:00:00.000Z"

	startDateTimeObj := util.ConvertStringToDateTime(startDateTime)
	endDateTimeObj := util.ConvertStringToDateTime(endDateTime)

	timeDiff := endDateTimeObj.DateTimeEpoch - startDateTimeObj.DateTimeEpoch

	partTotalDuration := 0
	for _, partMaster := range *generalObjects {
		toolingPartInfo := make(map[string]interface{})
		json.Unmarshal(partMaster.ObjectInfo, &toolingPartInfo)

		dayPart := util.InterfaceToInt(toolingPartInfo["day"])
		hourPart := util.InterfaceToInt(toolingPartInfo["hour"])
		minutePart := util.InterfaceToInt(toolingPartInfo["minute"])

		durationSec := dayPart*86400 + hourPart*3600 + minutePart*60

		partTotalDuration += durationSec

	}

	if timeDiff < int64(partTotalDuration) {
		isValid = false
	}

	return isValid
}

func isValidDuration(orderRequest map[string]interface{}) bool {
	valid := true
	durationType := util.InterfaceToString(orderRequest["durationType"])
	if durationType == "D" {
		if util.InterfaceToInt(orderRequest["day"]) <= 0 {
			valid = false
		}
	} else if durationType == "H" {
		if util.InterfaceToInt(orderRequest["hour"]) <= 0 {
			valid = false
		}
	} else if durationType == "M" {
		if util.InterfaceToInt(orderRequest["minute"]) <= 0 {
			valid = false
		}
	} else {
		valid = false
	}
	return valid
}

func reverseSlice(generalObjects *[]component.GeneralObject) *[]component.GeneralObject {
	if generalObjects != nil {
		for i, j := 0, len(*generalObjects)-1; i < j; i, j = i+1, j-1 {
			(*generalObjects)[i], (*generalObjects)[j] = (*generalObjects)[j], (*generalObjects)[i]
		}
	}

	return generalObjects
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
