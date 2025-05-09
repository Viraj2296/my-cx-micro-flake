package workflow_actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (v *ActionService) AddModificationCount(ctx *gin.Context) {

	recordId := util.GetRecordId(ctx)
	v.Logger.Info("handing add modification count ", zap.Int("recordId", recordId))
	err, mouldMasterGeneralObject := database.Get(v.Database, const_util.MouldMasterTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEventByMouldId(const_util.ProjectID, recordId)

	if scheduledEventObject != nil {
		if len(*scheduledEventObject) > 0 {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      const_util.GetError(const_util.ErrorCannotModifyData).Error(),
					Description: "Cannot change the status of this mould to 'Qualification' as it is currently being used for a scheduled production order",
				})
			return
		}
	}
	conditionString := " object_info ->> '$.mouldId' = " + strconv.Itoa(recordId)
	orderBy := " object_info ->> '$.createdAt' desc"
	listOfObjects, _ := database.GetConditionalObjectsOrderBy(v.Database, const_util.MouldTestRequestTable, conditionString, orderBy, 1)
	mouldTestId := 0
	if listOfObjects != nil {
		mouldTestId = (*listOfObjects)[0].Id
	}

	userId := common.GetUserId(ctx)
	var objectInfo = make(map[string]interface{})
	json.Unmarshal(mouldMasterGeneralObject.ObjectInfo, &objectInfo)

	if util.InterfaceToBool(objectInfo["canModify"]) {
		var modificationCount = util.InterfaceToInt(objectInfo["modificationCount"])
		objectInfo["modificationCount"] = modificationCount + 1
		objectInfo["mouldStatus"] = const_util.QualificationStatus

		objectInfo["lastUpdatedBy"] = userId
		objectInfo["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

		transactionIn := database.MouldModificationHistoryInfo{
			ModifyCount:        util.InterfaceToInt(objectInfo["modificationCount"]),
			MouldId:            recordId,
			CreatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			UpdatedAt:          util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
			CreatedBy:          userId,
			UpdatedBy:          0,
			MouldTestRequestId: mouldTestId, //assume modification is did for latest mould test request
		}

		var updateDate = make(map[string]interface{})
		var serialisedData, _ = json.Marshal(objectInfo)
		updateDate["object_info"] = serialisedData
		//v.Logger.Info("updating approve mould test request", zap.Any("approve_request", string(mouldMasterInfo.Serialize())))
		err = database.Update(v.Database, const_util.MouldMasterTable, recordId, updateDate)
		if err != nil {
			v.Logger.Error("error updating mould test request table", zap.Error(err))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "Internal system error",
					Description: "Internal system error starting test request, please summit error code to system admin",
				})
			return
		}
		generalObject := component.GeneralObject{
			ObjectInfo: transactionIn.Serialised(),
		}
		err, i := database.Create(v.Database, const_util.MouldModificationHistoryTable, generalObject)
		if err == nil {
			v.Logger.Info("Mould Modify History is successfully created", zap.Any("job", i))
		} else {
			v.Logger.Error("Mould Modify History creation failed", zap.Error(err))
		}

	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Modification count is added",
		Code:    0,
	})

}
