package helper

import (
	"fmt"
	"time"
)

func GenerateShiftID(number int) string {
	return "SHIFT" + fmt.Sprintf("%04d", number)
}

// ValidateShiftTimes  checks if the end date-time is after the start date-time
func ValidateShiftTimes(shiftStartDate, shiftStartTime, shiftEndDate, shiftEndTime string) (bool, error) {
	// Combine date and time strings into a single datetime string
	startDateTimeStr := shiftStartDate + " " + shiftStartTime
	endDateTimeStr := shiftEndDate + " " + shiftEndTime

	// Define the datetime format
	layout := "2006-01-02 15:04"

	// Parse the start and end datetime strings
	startDateTime, err := time.Parse(layout, startDateTimeStr)
	if err != nil {
		return false, fmt.Errorf("error parsing start datetime: %v", err)
	}
	endDateTime, err := time.Parse(layout, endDateTimeStr)
	if err != nil {
		return false, fmt.Errorf("error parsing end datetime: %v", err)
	}

	// Compare the two datetime values
	return endDateTime.After(startDateTime), nil
}

func IsDateLessThanCurrent(dateString string) (bool, error) {
	// Parse the input date string
	shiftStartDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return false, err
	}

	// Get the current date minus one day
	currentDate := time.Now().AddDate(0, 0, -1)

	// Compare the dates
	return shiftStartDate.Before(currentDate), nil
}
