package api_access

import (
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PrinterDetails struct {
	Id          int    `json:"id"`
	NetworkIP   string `json:"network_ip"`
	NetworkPort int    `json:"networkPort"`
	Location    int    `json:"location"`
}

func (v *APIService) GetPrinterDetails(ctx *gin.Context) {
	err, objects := database.GetObjects(v.DBConnection, const_util.BatchManagementPrinterTable)
	if err != nil {
		v.Logger.Error("error getting batch management printer details", zap.Error(err))
		response.SendResourceNotFound(ctx)
	}

	var result []PrinterDetails

	for _, object := range objects {
		printerInfo := database.GetPrinterInfo(object.ObjectInfo)
		result = append(result, PrinterDetails{Id: object.Id, NetworkPort: printerInfo.NetworkPort, NetworkIP: printerInfo.NetworkIP, Location: printerInfo.Location})
	}

	v.Logger.Info("get printer details", zap.Any("objects", result))
	ctx.JSON(http.StatusOK, result)
}
