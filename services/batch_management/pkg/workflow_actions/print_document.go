package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) PrintDocument(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	err, generalObject := database.Get(v.Database, targetTable, recordId)

	if err != nil {
		v.Logger.Error("error getting records", zap.String("error", err.Error()))
		error_util.SendResourceNotFound(ctx)
		return
	}
	var objectFields = make(map[string]interface{})
	err = json.Unmarshal(generalObject.ObjectInfo, &objectFields)
	if err != nil {
		v.Logger.Error("error un-marshal object", zap.Error(err))
		error_util.SendInternalSystemError(ctx)
	}
	var canPrint bool
	if canPrintInterface, ok := objectFields["canPrint"]; ok {
		canPrint = util.InterfaceToBool(canPrintInterface)
	}
	if canPrint {
		// publish the message to broker
		var updatingFields = make(map[string]interface{})
		objectFields["labelStatus"] = const_util.PrintingStatus
		objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		serialisedData, _ := json.Marshal(objectFields)
		updatingFields["object_info"] = serialisedData
		database.Update(v.Database, targetTable, recordId, updatingFields)
	} else {
		response.DispatchDetailedError(ctx, common.OperationNotPermitted,
			&response.DetailedError{
				Header:      "Invalid operation",
				Description: "This action can not be performed",
			})
	}
}

func (v *ActionService) publishJob(jobId int) {
	job, err := GetRawMaterialPrintJob(v.Database, jobId)
	if err != nil {
		v.Logger.Error("error getting job", zap.Error(err))
	}
	v.Logger.Info("publishing the job to broker", zap.Any("job", job))

}
