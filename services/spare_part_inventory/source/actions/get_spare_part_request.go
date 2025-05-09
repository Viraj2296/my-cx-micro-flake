package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/dto"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.cerex.io/transcendflow/logging"
	"go.uber.org/zap"
)

func (v *Actions) GetSparePartRequest(ctx *gin.Context) {
	v.Logger.Info("transfer out request received")
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)

	jobInfoRequest := dto.SparePartRequest{}
	if err := ctx.ShouldBindBodyWith(&jobInfoRequest, binding.JSON); err != nil {
		err = ctx.AbortWithError(http.StatusBadRequest, err)
		if err != nil {
			v.Logger.Error("error sending response", zap.Error(err))
		}
		return
	}

	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
	err, machineId := machineService.GetAssemblyMachineInfoFromEquipmentId(const_util.ProjectID, jobInfoRequest.EquipmentId)
	if err != nil {
		v.Logger.Error("error getting machine Id response", zap.Error(err))
		return
	}

	var condition = "JSON_CONTAINS(object_info->'$.machineIds', CAST(" + strconv.Itoa(machineId) + " AS JSON))"
	err, listOfSpareParts := v.Repository.GetResourceWithCondition(targetTable, condition)
	if err != nil {
		v.Logger.Error("error getting spare part master records", logging.String("error", err.Error()))
	}

	var sparePartList = make([]dto.SparePartRequestList, 0)

	for _, sparePart := range listOfSpareParts {

		err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(sparePart.ObjectInfo)
		if err != nil {
			v.Logger.Error("error getting records", logging.String("error", err.Error()))
		}

		err, objectInterface := v.Repository.GetResource(consts.SparePartInventoryLocationComponent, inventoryMasterInfo.Location)
		if err != nil {
			v.Logger.Error("error getting records", logging.String("error", err.Error()))
		}

		sparePartInventoryMasterInfo := models.GetSparePartInventoryLocationInfo(objectInterface.ObjectInfo)
		sparePartObject := dto.SparePartRequestList{}
		sparePartObject.Id = sparePart.Id
		sparePartObject.LocationName = sparePartInventoryMasterInfo.Name
		sparePartObject.OnHandQty = inventoryMasterInfo.OnHandQty
		sparePartObject.SparePartNumber = inventoryMasterInfo.SparePartNumber

		sparePartList = append(sparePartList, sparePartObject)

	}

	ctx.JSON(http.StatusOK, sparePartList)
	v.Logger.Info("sending downtime response", zap.Any("response", sparePartList))

}
