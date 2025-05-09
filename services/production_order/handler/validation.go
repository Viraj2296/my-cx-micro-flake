package handler

import (
	"time"
)

func scheduleValidation(scheduleStartDate string, scheduleEndDate string) bool {
	if scheduleStartDate == "" || scheduleEndDate == "" {
		return true
	}
	layout := "2006-01-02T15:04:05.000Z"
	startDateTime, errStart := time.Parse(layout, scheduleStartDate)
	endDateTime, errEnd := time.Parse(layout, scheduleEndDate)

	if errEnd != nil || errStart != nil {
		return false
	}

	if startDateTime.Unix() > endDateTime.Unix() {
		return false
	}

	return true

}
