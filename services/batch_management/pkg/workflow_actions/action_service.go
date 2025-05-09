package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type ActionService struct {
	Logger           *zap.Logger
	Database         *gorm.DB
	ComponentManager *common.ComponentManager
}

func (v *ActionService) Init(logger *zap.Logger) {
	v.Logger = logger
}
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
