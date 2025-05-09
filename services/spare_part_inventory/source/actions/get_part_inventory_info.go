package actions

import (
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/services/spare_part_inventory/source/dto"
	"net/http"

	"go.cerex.io/transcendflow/header_parser"
	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (v *Actions) GetPartInventoryInfo(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
	partInfoPayload := dto.PartInfoRequest{}
	if err := ctx.ShouldBindBodyWith(&partInfoPayload, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)
	var condition = "object_info->>'$.sparePartNumber' = '" + partInfoPayload.PartNumber + "'"

	errMaster, listOfRecords := v.Repository.GetResourceWithCondition(targetTable, condition)
	if errMaster != nil {
		v.Logger.Error("error getting records", logging.String("error", errMaster.Error()))
		error_util.SendResourceNotFound(ctx)
		return
	}

	ctx.JSON(http.StatusOK, listOfRecords)
}
