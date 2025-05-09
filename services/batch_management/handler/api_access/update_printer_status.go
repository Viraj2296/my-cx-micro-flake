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

type PrinterStatusRequest struct {
	Id               int    `json:"id"`
	ConnectionStatus string `json:"connectionStatus"`
	Message          string `json:"message"`
}

func (v *APIService) UpdatePrinterStatus(ctx *gin.Context) {

	var params PrinterStatusRequest
	if err := ctx.ShouldBindBodyWith(&params, binding.JSON); err != nil {
		err := ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("abort", zap.Error(err))
		}
		return
	}

	err, object := database.Get(v.DBConnection, const_util.BatchManagementPrinterTable, params.Id)
	if err != nil {
		v.Logger.Error("error getting batch management printer details", zap.Error(err))
		response.SendResourceNotFound(ctx)
	}

	printerInfo := database.GetPrinterInfo(object.ObjectInfo)
	printerInfo.Message = params.Message
	printerInfo.ConnectionStatus = params.ConnectionStatus

	var updateObject = make(map[string]interface{})
	updateObject["object_info"] = printerInfo.Serialised()

	updateError := database.Update(v.DBConnection, const_util.BatchManagementPrinterTable, params.Id, updateObject)
	if updateError != nil {
		v.Logger.Error("error updating printer status", zap.Error(updateError))
		error_util.SendResourceUpdateFailed(ctx)
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the resource",
	})
}
