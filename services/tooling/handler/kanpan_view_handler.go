package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type KanbanResponse struct {
	Status        string      `json:"status"`
	ColorCode     string      `json:"colorCode"`
	Id            int         `json:"id"`
	Cards         interface{} `json:"cards"`
	NumberOfCards int         `json:"numberOfCards"`
}

type KanbanTaskStatus struct {
	Id        int
	ColorCode string
}

func (v *ToolingService) getTasksWorkOrderTaskKanbanView(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	project := ctx.Query("project")
	taskNature := ctx.Query("nature")
	taskInterval := ctx.Query("interval")
	userId := common.GetUserId(ctx)
	var projectRecordId int
	if project != "" {
		projectRecordId = util.InterfaceToInt(project)
	} else {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Project",
				Description: "Sorry, system couldn't able to generate your task board due to invalid project identifier",
			})
		return
	}
	// first get the active sprint
	sprintCondition := "object_info->>'$.projectId'=" + strconv.Itoa(projectRecordId) + " and object_info->>'$.status'= 2"
	activeSprints, _ := GetConditionalObjects(dbConnection, ToolingProjectSprintTable, sprintCondition)
	if len(*activeSprints) > 1 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Sprints",
				Description: "Sorry, there are more than 1 sprint available in the system, please make sure only one sprint is available to show in board",
			})
		return
	}
	if len(*activeSprints) == 0 {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Sprints",
				Description: "Sorry, there are no active sprints available in the system, please make sure only one sprint is available to show in board",
			})
		return
	}
	sprintObject := (*(activeSprints))[0].ObjectInfo
	activeProjectSprint := ToolingProjectSprint{ObjectInfo: sprintObject}
	projectSprintInfo := activeProjectSprint.getToolingProjectSprintInfo()

	var taskStatusCache = make(map[int]*TaskStatusInfo)
	var taskStatusNameCache = make(map[string]*KanbanTaskStatus)
	var taskNameArray []string
	// since we don't have order by condition in where, we should using ID> 0 as an additional input
	listOfTaskStatus, _ := GetConditionalObjects(dbConnection, ToolingProjectTaskStatusTable, " ID > 0 ORDER BY ID ASC")
	for _, workOrderTaskStatusInterface := range *listOfTaskStatus {
		taskStatus := ToolingProjectTaskStatus{ObjectInfo: workOrderTaskStatusInterface.ObjectInfo}
		fmt.Println("task_id", workOrderTaskStatusInterface.Id, "task_name", taskStatus.getToolingProjectTaskStatusInfo().Status)
		taskStatusCache[workOrderTaskStatusInterface.Id] = taskStatus.getToolingProjectTaskStatusInfo()
		taskStatusNameCache[taskStatus.getToolingProjectTaskStatusInfo().Status] = &KanbanTaskStatus{
			Id:        workOrderTaskStatusInterface.Id,
			ColorCode: taskStatus.getToolingProjectTaskStatusInfo().ColorCode,
		}
		taskNameArray = append(taskNameArray, taskStatus.getToolingProjectTaskStatusInfo().Status)
	}
	var isStatusMatched bool

	var kanbanResponse []KanbanResponse
	// this for to send the list of card based on all the task status name, in order, we should maintain the task order status id
	for _, taskStatusName := range taskNameArray {
		isStatusMatched = false
		var cards []interface{}
		taskStatusInfo := taskStatusNameCache[taskStatusName]
		for _, taskId := range projectSprintInfo.ListOfAssignedTasks {
			_, taskObject := Get(dbConnection, ToolingProjectTaskTable, taskId)
			toolingProjectTask := ToolingProjectTask{ObjectInfo: taskObject.ObjectInfo}
			taskInfo := toolingProjectTask.getToolingTaskInfo()
			if taskNature != "" {
				if taskNature == "my_tasks" && taskInterval == "recently_updated" {
					if userId != taskInfo.AssignedUserId {
						continue
					}
					if !isLessThan24Hours(taskInfo.LastUpdatedAt) {
						continue
					}
				} else if taskNature == "my_tasks" {
					if userId != taskInfo.AssignedUserId {
						continue
					}
				} else if taskInterval == "recently_updated" {
					if !isLessThan24Hours(taskInfo.LastUpdatedAt) {
						continue
					}

				}
			}
			taskStatusId := taskInfo.Status
			if taskStatusInfo.Id == taskStatusId {
				assignedUserId := taskInfo.AssignedUserId
				authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
				userInfo := authService.GetUserInfoById(assignedUserId)
				var listOfUserInfo []interface{}
				listOfUserInfo = append(listOfUserInfo, userInfo)
				var taskFields = make(map[string]interface{})
				json.Unmarshal(toolingProjectTask.ObjectInfo, &taskFields)
				projectRecordId := util.InterfaceToInt(taskFields["projectId"])
				name := util.InterfaceToString(taskFields["name"])
				description := util.InterfaceToString(taskFields["description"])
				plannedStartDate := util.InterfaceToString(taskFields["plannedStartDate"])
				taskFields["parent"] = projectRecordId
				taskFields["taskName"] = name
				taskFields["shortDescription"] = description
				taskFields["taskDate"] = plannedStartDate
				cards = append(cards, taskFields)
				taskFields["assignedUsers"] = listOfUserInfo
				taskFields["id"] = taskId
				isStatusMatched = true
			}

		}
		var response KanbanResponse
		if !isStatusMatched {
			isStatusMatched = false
			response = KanbanResponse{
				Status:        taskStatusName,
				NumberOfCards: len(cards),
				Cards:         make([]string, 0),
				ColorCode:     taskStatusNameCache[taskStatusName].ColorCode,
				Id:            taskStatusNameCache[taskStatusName].Id,
			}
		} else {
			response = KanbanResponse{
				Status:        taskStatusName,
				NumberOfCards: len(cards),
				Cards:         cards,
				ColorCode:     taskStatusNameCache[taskStatusName].ColorCode,
				Id:            taskStatusNameCache[taskStatusName].Id,
			}
		}
		kanbanResponse = append(kanbanResponse, response)
	}

	ctx.JSON(http.StatusOK, kanbanResponse)
}

func (v *ToolingService) kanbanMoveTask(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)

	recordId := util.GetRecordId(ctx)
	projectId := ctx.Param("projectId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, eventObject := Get(dbConnection, targetTable, recordId)
	if err != nil {
		response.SendResourceNotFound(ctx)
		return
	}
	userId := common.GetUserId(ctx)

	kanbanTaskStatus := KanbanTaskStatusRequest{}
	if err := ctx.ShouldBindBodyWith(&kanbanTaskStatus, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//// before update, check the work order is completed , then don.t let the user to move
	//if err, isTaskIsDone := ms.isWorkOrderCompletedFromTaskId(dbConnection, recordId); err == nil {
	//	if isTaskIsDone {
	//		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
	//			&response.DetailedError{
	//				Header:      getError(common.InvalidObjectStatusError).Error(),
	//				Description: "Sorry, this work order is already been done, further modification is not allowed",
	//			})
	//		return
	//	}
	//}

	// we need to check whether we have any tasks created before release
	toolingProjectTask := ToolingProjectTask{ObjectInfo: eventObject.ObjectInfo}
	toolingTaskInfo := toolingProjectTask.getToolingTaskInfo()

	// existing task status id
	toolingTaskInfo.Status = kanbanTaskStatus.TaskStatus

	err = Update(dbConnection, targetTable, recordId, toolingTaskInfo.DatabaseSerialize(userId))
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines event information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Message: "Tooling task is successfully updated",
		Code:    0,
	})

}
