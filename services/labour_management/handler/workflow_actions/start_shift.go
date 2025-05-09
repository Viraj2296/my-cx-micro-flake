package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) HandleStartShift(ctx *gin.Context) {
	v.Logger.Info("handle start shift is received")
	recordId, dbConnection, userInfo, shiftMasterInfo := v.getBasicInfo(ctx)
	if shiftMasterInfo.CanShiftStart {
		// update the shift status as active
		shiftMasterInfo.ShiftStatus = const_util.ShiftStatusActive
		shiftMasterInfo.CanShiftStop = true
		shiftMasterInfo.CanShiftStart = false
		shiftMasterInfo.CanCheckIn = true
		shiftMasterInfo.LastUpdatedBy = userInfo.UserId
		shiftMasterInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		var updateObject = make(map[string]interface{})
		updateObject["object_info"] = shiftMasterInfo.Serialised()
		err := database.Update(dbConnection, const_util.LabourManagementShiftMasterTable, recordId, updateObject)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Server Exception",
					Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
				})
			return
		}
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		orderStatusId := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(const_util.ProjectID, const_util.ScheduleStatusPreferenceFive)

		// update all the events as running
		for _, eventId := range shiftMasterInfo.ScheduledOrderEvents {
			var updatingData = make(map[string]interface{})
			updatingData["eventStatus"] = orderStatusId
			serializedObject, _ := json.Marshal(updatingData)
			err = productionOrderInterface.UpdateAssemblyScheduledOrderFields(const_util.ProjectID, eventId, serializedObject)
			if err != nil {
				v.Logger.Error("error updating scheduler order event during start of the shift", zap.Error(err))
			}

			//	Insert new hmi entry
			err = v.createAssemblyHmiEntry(eventId, userInfo.UserId, "started")
			if err != nil {
				v.Logger.Error("error in creating start hmi through labour management", zap.Error(err))
			}
		}

		// update the manpower flag to shift production now
		var shiftProductionCondition = " object_info->>'$.shiftId' = " + strconv.Itoa(recordId)
		err, shiftProdInterfaceObjects := database.GetConditionalObjects(v.Database, const_util.LabourManagementShiftProductionTable, shiftProductionCondition)
		if err == nil {
			for _, shiftProdInterfaceObject := range shiftProdInterfaceObjects {
				shiftProductionInfo := database.GetShiftMasterProductionInfo(shiftProdInterfaceObject.ObjectInfo)
				shiftProductionInfo.CanUpdateManpower = true
				var updateShiftProdObject = make(map[string]interface{})
				updateShiftProdObject["object_info"] = shiftProductionInfo.Serialised()
				err = database.Update(dbConnection, const_util.LabourManagementShiftProductionTable, shiftProdInterfaceObject.Id, updateShiftProdObject)
				if err != nil {
					v.Logger.Error("error updating scheduler production canUpdateManpower flag", zap.Error(err))
				}
			}
		}
		v.Logger.Info("handle user submit is successfully processed", zap.Any("record_id", recordId))
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    0,
			Message: "Your shift has been started successfully",
		})
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Action",
				Description: "This operation is invalid, you can not start the shift at this movement, Please report this to system admin",
			})
	}

}
