package actions

import (
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/models"
	"github.com/gin-gonic/gin"
	"go.cerex.io/transcendflow/orm"
	"go.uber.org/zap"
	"net/http"
)

type FaultCodes struct {
	Code string `json:"code"`
	Id   int    `json:"id"`
}
type FaultType struct {
	FaultType string
	Id        int
}
type FaultCode struct {
	Code      string
	FaultType int
	Id        int
}
type FaultInfo struct {
	FaultType  string       `json:"faultType"`
	Id         int          `json:"id"`
	FaultCodes []FaultCodes `json:"faultCodes"`
}

func (v *Actions) HandleGetFaultInfo(ctx *gin.Context) {
	// get the list of fault type
	err, i := orm.GetObjects(v.Database, consts.MachineDownTimeFaultTypeTable)
	var faultTypeList []FaultType
	if err == nil {
		for _, faultTypeObjects := range i {
			faultTypeInfo := models.GetMachineDowntimeFaultTypeInfo(faultTypeObjects.ObjectInfo)
			faultTypeList = append(faultTypeList, FaultType{
				FaultType: faultTypeInfo.Name,
				Id:        faultTypeObjects.Id,
			})
		}
	}
	err, i = orm.GetObjects(v.Database, consts.MachineDownTimeFaultCodeTable)
	var faultCodeList []FaultCode
	if err == nil {
		for _, faultCodeObject := range i {
			faultCodeInfo := models.GetMachineDowntimeFaultCodeInfo(faultCodeObject.ObjectInfo)
			faultCodeList = append(faultCodeList, FaultCode{
				FaultType: faultCodeInfo.FaultType,
				Code:      faultCodeInfo.Name,
				Id:        faultCodeObject.Id,
			})
		}
	}

	// Generate the array of FaultInfo
	var faultInfoList []FaultInfo
	for _, faultType := range faultTypeList {
		// Filter FaultCodes that match the FaultType ID
		var faultCodes []FaultCodes
		for _, faultCode := range faultCodeList {
			if faultCode.FaultType == faultType.Id {
				faultCodes = append(faultCodes, FaultCodes{
					Code: faultCode.Code,
					Id:   faultCode.Id,
				})
			}
		}

		// Add a new FaultInfo entry
		faultInfoList = append(faultInfoList, FaultInfo{
			FaultType:  faultType.FaultType,
			Id:         faultType.Id,
			FaultCodes: faultCodes,
		})
	}
	v.Logger.Info("generated fault info list ", zap.Any("fault_info_list", faultInfoList))
	ctx.JSON(http.StatusOK, faultInfoList)
}
