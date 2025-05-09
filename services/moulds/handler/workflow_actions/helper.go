package workflow_actions

import (
	"cx-micro-flake/pkg/util"
	"strconv"
)

func GetTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}
