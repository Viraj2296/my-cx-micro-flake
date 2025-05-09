package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (v *ToolingService) getActiveSprint(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	sprintCondition := "object_info->>'$.projectId'=" + strconv.Itoa(recordId) + " and object_info->>'$.status'= 2"
	activeSprints, _ := GetConditionalObjects(dbConnection, ToolingProjectSprintTable, sprintCondition)
	if len(*activeSprints) == 0 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "No Active Sprints",
				Description: "Sorry, there is no active sprints in the system. Please activate your project sprint",
			})
		return
	}
	if len(*activeSprints) > 1 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Sprints",
				Description: "Sorry, there are more than 1 sprint available in the system, please make sure only one sprint is available to show in board",
			})
		return
	}
	sprintObject := (*(activeSprints))[0].ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(sprintObject, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, recordId, componentName, rawJSONObject)
	ctx.JSON(http.StatusOK, response)
}
