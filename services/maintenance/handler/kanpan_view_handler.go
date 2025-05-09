package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (v *MaintenanceService) getMaintenanceWorkOrderTaskKanbanView(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfWorkOrderTasks, _ := GetObjects(dbConnection, MaintenanceWorkOrderTaskTable)

	var taskStatusCache = make(map[int]*MaintenanceWorkOrderTaskStatusInfo)
	var taskStatusNameCache = make(map[string]*KanbanTaskStatus)
	var taskNameArray []string
	// since we don't have order by condition in where, we should using ID> 0 as an additional input
	listOfTaskStatus, _ := GetConditionalObjects(dbConnection, MaintenanceWorkOrderTaskStatusTable, " ID > 0 ORDER BY ID ASC")
	for _, workOrderTaskStatusInterface := range *listOfTaskStatus {
		workOrderTaskStatus := MaintenanceWorkOrderTaskStatus{ObjectInfo: workOrderTaskStatusInterface.ObjectInfo}
		fmt.Println("task_id", workOrderTaskStatusInterface.Id, "task_name", workOrderTaskStatus.getMaintenanceWorkOrderTaskStatusInfo().Status)
		taskStatusCache[workOrderTaskStatusInterface.Id] = workOrderTaskStatus.getMaintenanceWorkOrderTaskStatusInfo()
		taskStatusNameCache[workOrderTaskStatus.getMaintenanceWorkOrderTaskStatusInfo().Status] = &KanbanTaskStatus{
			Id:        workOrderTaskStatusInterface.Id,
			ColorCode: workOrderTaskStatus.getMaintenanceWorkOrderTaskStatusInfo().ColorCode,
		}
		taskNameArray = append(taskNameArray, workOrderTaskStatus.getMaintenanceWorkOrderTaskStatusInfo().Status)
	}
	var isStatusMatched bool

	var kanbanResponse []KanbanResponse
	// this for to send the list of card based on all the task status name, in order, we should maintain the task order status id
	for _, taskStatusName := range taskNameArray {
		isStatusMatched = false
		var cards []interface{}
		taskStatusInfo := taskStatusNameCache[taskStatusName]
		for _, workOrderTask := range *listOfWorkOrderTasks {
			maintenanceWorkOrderTask := MaintenanceWorkOrderTask{ObjectInfo: workOrderTask.ObjectInfo}
			taskStatusId := maintenanceWorkOrderTask.getWorkOrderTaskInfo().TaskStatus
			if taskStatusInfo.Id == taskStatusId {
				assignedUserId := maintenanceWorkOrderTask.getWorkOrderTaskInfo().AssignedUserId
				authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
				userInfo := authService.GetUserInfoById(assignedUserId)
				var listOfUserInfo []interface{}
				listOfUserInfo = append(listOfUserInfo, userInfo)
				var taskFields = make(map[string]interface{})
				json.Unmarshal(maintenanceWorkOrderTask.ObjectInfo, &taskFields)
				workOrderId := util.InterfaceToInt(taskFields["workOrderId"])
				workOrderTaskId := util.InterfaceToString(taskFields["workOrderTaskId"])
				taskFields["parent"] = workOrderId
				taskFields["taskReferenceId"] = workOrderTaskId
				cards = append(cards, taskFields)
				taskFields["assignedUsers"] = listOfUserInfo
				taskFields["id"] = workOrderTask.Id
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
