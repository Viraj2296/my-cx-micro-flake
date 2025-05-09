package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/header_parser"
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/dto"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.cerex.io/transcendflow/logging"

	"github.com/gin-gonic/gin"
)

func (v *Actions) RequestedSpareParts(ctx *gin.Context) {
	v.Logger.Info("handle check in received")
	componentName := header_parser.GetComponentName(ctx)
	var targetTable = v.ComponentManager.GetTargetTable(componentName)
	var condition = "object_status = 'Active'"
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	// labourService := common.GetService("labour_management_module").ServiceInterface.(common.LabourManagementInterface)

	err, listOfRecords := v.Repository.GetResourceWithCondition(targetTable, condition)
	if err != nil {
		v.Logger.Error("error getting records", logging.String("error", err.Error()))
		error_util.SendResourceNotFound(ctx)
		return
	}
	var sparePartList = make([]dto.SparePartRepairData, 0)
	for _, repair := range listOfRecords {
		v.Logger.Info("found repair", logging.Any("repair", repair.Id))
		repairObjectData := dto.SparePartRepairData{}

		err, repairObject := models.GetSparePartInventoryRepairRequestInfo(repair.ObjectInfo)
		if err != nil {
			v.Logger.Error("error getting repair object", logging.String("error", err.Error()))
			continue
		}
		err, machineObject := v.Repository.GetResource("assembly_machine_master", repairObject.MachineId)

		var info map[string]interface{}
		_ = json.Unmarshal([]byte(machineObject.ObjectInfo), &info)

		newMachineId, _ := info["newMachineId"].(string)

		if err != nil {
			v.Logger.Error("error getting records", logging.String("error", err.Error()))
			error_util.SendResourceNotFound(ctx)
			return
		}

		basicUserInfo := authService.GetUserInfoById(repair.CreatedBy)
		var partList = make([]dto.PartList, 0)
		for _, part := range repairObject.SpareParts {
			err, objectInterface := v.Repository.GetResource(consts.SparePartInventoryMasterComponent, part.SparePartId)
			if err != nil {
				v.Logger.Error("error getting records", logging.String("error", err.Error()))
			}
			err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(objectInterface.ObjectInfo)
			if err != nil {
				v.Logger.Error("error getting records", logging.String("error", err.Error()))
			}

			partObject := dto.PartList{}
			partObject.Quantity = part.Quantity
			partObject.SparePartId = part.SparePartId
			partObject.SparePartName = inventoryMasterInfo.SparePartNumber
			partList = append(partList, partObject)

		}

		createDate := customDateTime(repair.CreatedAt.String())
		repairObjectData.Id = repair.Id
		repairObjectData.SpareParts = partList
		repairObjectData.JobId = repairObject.JobId
		repairObjectData.RequestStatus = repairObject.RequestStatus
		repairObjectData.IsNeedSparePart = repairObject.IsNeedSparePart
		repairObjectData.Created_by = basicUserInfo.FullName
		repairObjectData.Created_at = createDate
		repairObjectData.MachineName = newMachineId
		sparePartList = append(sparePartList, repairObjectData)

	}

	ctx.JSON(http.StatusOK, sparePartList)
}
func customDateTime(createdAt string) string {
	parts := strings.Fields(createdAt)

	if len(parts) > 3 {
		createdAt = strings.Join(parts[:3], " ")
	}

	layout := "2006-01-02 15:04:05 -0700"
	utcTime, err := time.Parse(layout, createdAt)
	if err != nil {
		panic(err)
	}

	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		panic(err)
	}
	singaporeTime := utcTime.In(loc)
	output := singaporeTime.Format("02 Jan 2006 03:04 PM")

	return output
}
