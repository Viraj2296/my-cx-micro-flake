package workflow_actions

import (
	"cx-micro-flake/pkg/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetShiftLines to get the shift lines
func (v *ActionService) GetShiftLines(ctx *gin.Context) {
	v.Logger.Info("handle get shift lines")
	var shiftId = util.GetRecordId(ctx)

	var listOfAssemblyLines = v.getLinesForSchedulerEvents(shiftId)
	v.Logger.Info("got the list of lines to shift", zap.Int("shift_id", shiftId), zap.Any("assemblies", listOfAssemblyLines))
	ctx.JSON(http.StatusOK, listOfAssemblyLines)
}
