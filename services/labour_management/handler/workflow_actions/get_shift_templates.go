package workflow_actions

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetShiftTemplates  : Getting all the shift templates
func (v *ActionService) GetShiftTemplates(ctx *gin.Context) {
	v.Logger.Info("handle get shift templates")

	err, listOfShiftTemplateInterface := database.GetObjects(v.Database, const_util.LabourManagementShiftTemplateTable)
	if err != nil {
		v.Logger.Error("error getting shift template", zap.String("error", err.Error()))
		response.SendResourceNotFound(ctx)
	}

	var dropDownArray = make([]component.OrderedData, 0)
	var recordInfo = component.RecordInfo{}
	var defaultIndex = 0
	var defaultIndexValue = ""
	for index, shiftTemplateInterface := range listOfShiftTemplateInterface {
		shiftTemplateInfo := database.GetSShiftTemplateInfo(shiftTemplateInterface.ObjectInfo)
		if shiftTemplateInfo.ObjectStatus == component.ObjectStatusActive {
			if index == 0 {
				defaultIndex = index + 1
				defaultIndexValue = shiftTemplateInfo.Name
			}
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:                       index + 1,
				Value:                    shiftTemplateInfo.Name,
				OnValueConditionalFields: nil,
			})
		}

	}
	recordInfo.Value = defaultIndexValue
	recordInfo.Type = "int"
	recordInfo.Index = defaultIndex
	recordInfo.Data = dropDownArray
	recordInfo.IsEdit = false
	ctx.JSON(http.StatusOK, recordInfo)
}
