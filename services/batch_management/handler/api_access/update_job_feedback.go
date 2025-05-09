package api_access

import (
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
)

type JobFeedbackRequest struct {
	Id      int    `json:"id"`
	JobType string `json:"jobType"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (v *APIService) UpdateJobFeedback(ctx *gin.Context) {
	var params JobFeedbackRequest
	if err := ctx.ShouldBindBodyWith(&params, binding.JSON); err != nil {
		err := ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("abort", zap.Error(err))
		}
		return
	}
	if params.JobType == "raw_material" {
		err, object := database.Get(v.DBConnection, const_util.BatchManagementRawMaterialTable, params.Id)
		if err != nil {
			v.Logger.Error("error getting batch management printer details", zap.Error(err))
			response.SendResourceNotFound(ctx)
		}
		rawMaterialBatchInfo := database.GeRawMaterialBatchInfo(object.ObjectInfo)
		rawMaterialBatchInfo.LabelStatus = params.Status
		rawMaterialBatchInfo.LabelPrintMessage = params.Message
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = rawMaterialBatchInfo.Serialised()

		updateError := database.Update(v.DBConnection, const_util.BatchManagementRawMaterialTable, params.Id, updateObject)
		if updateError != nil {
			v.Logger.Error("error updating raw material job status", zap.Error(updateError))
			error_util.SendResourceUpdateFailed(ctx)
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Successfully updated the resource",
		})
	} else if params.JobType == "mold_batch" {
		err, object := database.Get(v.DBConnection, const_util.BatchManagementMouldTable, params.Id)
		if err != nil {
			v.Logger.Error("error getting batch management printer details", zap.Error(err))
			response.SendResourceNotFound(ctx)
		}
		mouldBatchInfo := database.GeRawMouldBatchInfo(object.ObjectInfo)
		mouldBatchInfo.LabelStatus = params.Status
		mouldBatchInfo.LabelPrintMessage = params.Message
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = mouldBatchInfo.Serialised()

		updateError := database.Update(v.DBConnection, const_util.BatchManagementMouldTable, params.Id, updateObject)
		if updateError != nil {
			v.Logger.Error("error updating mould batch job status", zap.Error(updateError))
			error_util.SendResourceUpdateFailed(ctx)
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Successfully updated the resource",
		})
	}

}
