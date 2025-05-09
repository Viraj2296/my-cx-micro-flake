package util

import (
	"strconv"
	"time"
)

const (
	ISOTimeLayout       = "2006-01-02T15:04:05.000Z"
	ISODateLayout       = "2006-01-02T00:00:00.000Z"
	TableDateTimeLayout = "Jan 2, 2006, 03:04:05 PM"
	TableDateLayout     = "Jan 2, 2006"
)

type DateTimeInfo struct {
	DateTimeString     string
	DateTime           time.Time
	DateTimeEpochMilli int64
	DateTimeEpoch      int64
	Error              error
}

func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetCurrentTime(format string) string {
	return time.Now().Format(format)
}

func GetCurrentTimeInSingapore(format string) string {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return ""
	}
	return time.Now().In(loc).Format(format)
}

func GetCurrentTimeWithOffSet(duration time.Duration, format string) string {
	return time.Now().Add(duration).Format(format)
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetCurrentTimeOnly Function to get the current time in HH-MM-SS format
func GetCurrentTimeOnly() string {
	return time.Now().Format("15:04:05")
}

// GetZoneCurrentTime download all the zone avaiable from this https://github.com/golang/go/blob/99fa49e4b7f74fb21cc811270fe42c9b7fa99668/lib/time/zoneinfo.zip
func GetZoneCurrentTime(zone string) string {
	// zone : "Asia/Singapore"
	layout := "2006-01-02T15:04:05.000Z"
	localZone, _ := time.LoadLocation(zone)
	currentTimeInLocal := time.Now().In(localZone)
	layoutTime := currentTimeInLocal.Format(layout)
	return layoutTime
}

func GetZoneCurrentTimeYYYYMMDD(zone string) string {
	// zone : "Asia/Singapore"
	layout := "20060102"
	localZone, _ := time.LoadLocation(zone)
	currentTimeInLocal := time.Now().In(localZone)
	layoutTime := currentTimeInLocal.Format(layout)
	return layoutTime
}
func GetZoneCurrentTimeInPMFormat(zone string) string {
	// zone : "Asia/Singapore"
	layout := "2006-01-02 03:04:05 PM"
	localZone, _ := time.LoadLocation(zone)
	currentTimeInLocal := time.Now().In(localZone)
	layoutTime := currentTimeInLocal.Format(layout)
	return layoutTime
}

func ConvertDateToTimeZonCorrectedPrimeNg(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format("01/02/2006")
}

func ConvertUserTimezoneToUTC(userTimeZone string, datetime string) string {
	parsedTime, _ := time.Parse(ISOTimeLayout, datetime)

	userZone, _ := time.LoadLocation(userTimeZone)
	parsedDateTimeLocal := parsedTime.In(userZone)
	_, offSet := parsedDateTimeLocal.Zone()
	offSetHours := int(offSet / 3600)
	duration := time.Duration(int(time.Hour) * -1 * offSetHours)
	dateTimeParsed := parsedTime.Add(duration)
	return dateTimeParsed.Format(ISODateLayout)
}
func ConvertDateTimePrimeNgTable(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(ISOTimeLayout)
}

func ISO2TableDateTimeFormat(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(TableDateTimeLayout)
}

func ISO2TableDateFormat(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(TableDateLayout)
}

func ConvertDatePrimeNgTable(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(ISODateLayout)
}

// ConvertSingaporeTimeToUTC we need to have global function to convert zone to zone conversion
func ConvertSingaporeTimeToUTC(datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	dateTimeParsed = dateTimeParsed.Add(time.Hour * -8)

	layout := "2006-01-02T15:04:05.000Z"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return dateTimeParsed.Format(layout)
}

func ConvertTimeToTimeZonCorrectedPrimeNgTable(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02T15:04:05.000Z"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func ConvertTimeToTimeZonCorrectedFormat(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02 03:04:05 PM"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func ConvertTimeToTimeZonCorrectedPrimeNg(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format("01/02/2006 15:04:05")
}

func ConvertTimeToTimeZonCorrectedBryntum(zone string, datetime string) string {
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02T15:04:05.000Z"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func ConvertPrimeNgDateTimeToUTC(zone string, datetime string) string {
	// Parse the input datetime string in "01/02/2006 15:04:05" format
	dateTimeParsed, err := time.Parse("01/02/2006 15:04:05", datetime)
	if err != nil {
		return ""
	}

	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02T15:04:05.000Z"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func CalculateHoursFromCurrentTime(datetime string) (float64, error) {
	// Parse the input datetime string in the "2006-01-02T15:04:05.000Z" format
	parsedTime, err := time.Parse(time.RFC3339Nano, datetime)
	if err != nil {
		return 0, err
	}

	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Calculate the duration where parsed time is subtracted from current time
	duration := parsedTime.Sub(currentTime)

	// Convert the duration to hours
	hours := duration.Hours()

	return hours, nil
}

func ConvertTimeToTimeZonCorrected(zone string, datetime string) string {
	if datetime == "" {
		return ""
	}
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02T15:04:05"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func ConvertTimeToTimeZonLongCorrected(zone string, datetime string) string {
	if datetime == "" {
		return "-"
	}
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "Mon Jan 2 15:04:05 2006"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}
func ConvertTimeToZeroZone(zone string, datetime string) string {
	if datetime == "" {
		return ""
	}
	dateTimeParsed, _ := time.Parse(time.RFC3339, datetime)
	localZone, _ := time.LoadLocation(zone)
	parsedDateTimeLocal := dateTimeParsed.In(localZone)
	layout := "2006-01-02T15:04:05Z"
	// convert this time mm/dd/yyyy, this is the format front-end can understand
	return parsedDateTimeLocal.Format(layout)
}

func CompareDateTime(src string, dst string) bool {
	layout := "2006-01-02T15:04:05.000Z"
	startDateTime, err := time.Parse(layout, src)
	if err != nil {
		return false
	}
	endDateTime, err := time.Parse(layout, dst)
	if err != nil {
		return false
	}
	if startDateTime.UnixMilli() < endDateTime.UnixMilli() {
		return true
	}

	return false
}

func ConvertStringToDateTime(dateTimeString string) DateTimeInfo {
	layout := "2006-01-02T15:04:05.000Z"
	dateTime, err := time.Parse(layout, dateTimeString)
	if err != nil {
		return DateTimeInfo{Error: err}
	}
	return DateTimeInfo{DateTimeString: dateTimeString, DateTime: dateTime, DateTimeEpochMilli: dateTime.UTC().UnixMilli(), DateTimeEpoch: dateTime.UTC().Unix(), Error: nil}
}

func ConvertStringToDateTimeV2(dateTimeString string) DateTimeInfo {
	layout := "2006-01-02T15:04:05Z"
	dateTime, err := time.Parse(layout, dateTimeString)
	if err != nil {
		return DateTimeInfo{Error: err}
	}
	return DateTimeInfo{DateTimeString: dateTimeString, DateTime: dateTime, DateTimeEpochMilli: dateTime.UTC().UnixMilli(), DateTimeEpoch: dateTime.UTC().Unix(), Error: nil}
}

func ConvertReferenceTimeToString(difference int64) string {
	if difference < 60 {
		// this shoud be seconds
		return strconv.Itoa(int(difference)) + " sec"
	} else if difference > 60 && difference < 3600 {
		approximateDiff := difference / 60
		return strconv.Itoa(int(approximateDiff)) + " min"
	} else if difference > 3600 && difference < 86400 {
		approximateDiff := difference / 3600
		return strconv.Itoa(int(approximateDiff)) + " hour"
	} else if difference > 86400 {
		approximateDiff := difference / 86400
		return strconv.Itoa(int(approximateDiff)) + " hour"
	}
	return " invalid"
}

func IsDateString(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
func HumanReadable12HoursDateTimeFormat(datetime string) string {
	dateTimeParsed, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", datetime)
	return dateTimeParsed.Format("02-01-2006 15:04")
}
func HumanReadableDateFormat(datetime string) string {
	dateTimeParsed, _ := time.Parse("2006-01-02", datetime)
	return dateTimeParsed.Format("02-01-2006")
}
