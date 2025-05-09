package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func updateToolingSprint(toolingTaskId int, toolingProjectTask ToolingProjectTask, dbConnection *gorm.DB) {
	toolingTaskInfo := toolingProjectTask.getToolingTaskInfo()
	projectId := toolingTaskInfo.ProjectId

	conditionalString := " JSON_EXTRACT(object_info, \"$.projectId\")=" + strconv.Itoa(projectId) + " "
	springGeneralObject, err := GetConditionalObjects(dbConnection, ToolingProjectSprintTable, conditionalString)

	if err != nil || len(*springGeneralObject) == 0 {
		return
	}
	isBreakSpringLoop := false
	//This method put assigned task id into completed
	//But this has time constraint
	//There must be a scheduler update sprint when sprint is overdued
	for _, toolingSprintObject := range *springGeneralObject {
		if isBreakSpringLoop {
			break
		}

		toolingSprint := ToolingProjectSprint{Id: toolingSprintObject.Id, ObjectInfo: toolingSprintObject.ObjectInfo}
		toolingSprintInfo := toolingSprint.getToolingProjectSprintInfo()

		taskList := toolingSprintInfo.ListOfAssignedTasks
		for _, assignedTask := range taskList {
			if assignedTask == toolingTaskId {
				totalTask := len(toolingSprintInfo.ListOfAssignedTasks)
				forcastTask := totalTask - len(toolingSprintInfo.ListOfCompletedTasks)
				toolingSprintInfo.ListOfCompletedTasks = append(toolingSprintInfo.ListOfCompletedTasks, toolingTaskId)
				// toolingSprintInfo.ListOfAssignedTasks = remove(toolingSprintInfo.ListOfAssignedTasks, index)

				if totalTask != 0 {
					toolingSprintInfo.ActualAmount = float64(len(toolingSprintInfo.ListOfCompletedTasks)) / float64(totalTask)
					toolingSprintInfo.ForcastAmount = float64(forcastTask) / float64(totalTask)
				}

				_ = Update(dbConnection, ToolingProjectSprintTable, toolingSprintObject.Id, toolingSprintInfo.DatabaseSerialize())
				updateCompletedPercentage(dbConnection, projectId, toolingSprint)
				isBreakSpringLoop = true
				break
			}
		}
	}
}

func updateCompletedPercentage(dbConnection *gorm.DB, toolingProjectId int, toolingProjectSprint ToolingProjectSprint) {

	_, toolingProjectGeneralObj := Get(dbConnection, ToolingProjectTable, toolingProjectId)
	toolingProject := ToolingProject{ObjectInfo: toolingProjectGeneralObj.ObjectInfo}
	projectInfo := toolingProject.getToolingProjectInfo()

	toolingSpringInfo := toolingProjectSprint.getToolingProjectSprintInfo()
	startDate := util.ConvertStringToDateTime(toolingSpringInfo.StartDate)

	sprintYear, _, _ := startDate.DateTime.Date()
	dateCondition := time.Date(sprintYear, 1, 1, 0, 0, 0, 0, time.UTC)

	dateString := dateCondition.Format("2006-01-02T15:04:05.000Z")
	conditionalString := " object_info->>'$.startDate' > '" + dateString + "' AND object_info->>'$.projectId' = " + strconv.Itoa(toolingProjectId) + " "
	currentYearSprintGeneralObject, _ := GetConditionalObjects(dbConnection, ToolingProjectSprintTable, conditionalString)

	monthlyActualAmount := calcActualAmountMonthly(currentYearSprintGeneralObject, *projectInfo)

	var actualAmountForYear float64
	for _, element := range monthlyActualAmount {
		actualAmountForYear += element
	}

	monthlyForcastAmount := calcForcastAmountMonthly(currentYearSprintGeneralObject, *projectInfo)

	var actualForcastForYear float64
	for _, element := range monthlyForcastAmount {
		actualForcastForYear += element
	}

	fmt.Println("actualAmountForYear ", actualAmountForYear)

	completedPercentage := (actualAmountForYear / float64(projectInfo.AllocatedBudget)) * float64(100)

	fmt.Println("completedPercentage ", completedPercentage)

	projectInfo.TotalAmountUsed = actualAmountForYear
	projectInfo.TotalForcastAmount = actualForcastForYear
	projectInfo.CompletedPercentage = int(completedPercentage)
	projectInfo.Status = setStatus(projectInfo.CompletedPercentage)

	if projectInfo.Status == ProjectAPPROVAL {
		projectInfo.CanApprove = true
	}

	_ = Update(dbConnection, ToolingProjectTable, toolingProjectId, projectInfo.DatabaseSerialize())

}

func setStatus(completedPercentage int) int {
	status := 1
	switch {
	case completedPercentage <= 2:
		status = ProjectDFM
	case completedPercentage >= 3 && completedPercentage <= 5:
		status = ProjectDESIGN
	case completedPercentage >= 6 && completedPercentage <= 85:
		status = ProjectFABRICATION
	case completedPercentage >= 86 && completedPercentage <= 99:
		status = ProjectAPPROVAL
	case completedPercentage == 100:
		status = ProjectCLOSED
	default:
		status = ProjectCREATED
	}

	return status
}

func calcActualAmountMonthly(currentYearSprintGeneralObject *[]component.GeneralObject, projectInfo ToolingProjectInfo) map[int]float64 {
	monthlyActualAmount := make(map[int]float64)

	for _, sprintObject := range *currentYearSprintGeneralObject {
		toolingSprint := ToolingProjectSprint{Id: sprintObject.Id, ObjectInfo: sprintObject.ObjectInfo}
		toolingSprintInfo := toolingSprint.getToolingProjectSprintInfo()

		startDate := util.ConvertStringToDateTime(toolingSprintInfo.StartDate)
		_, sprintMonth, _ := startDate.DateTime.Date()
		monthInt := int(sprintMonth)

		if sprintMonth == 1 {
			monthlyActualAmount[monthInt] = (toolingSprintInfo.ActualAmount * float64(projectInfo.AllocatedBudget)) - projectInfo.LastYearAmountUsed
			fmt.Println("Sprint mont ", monthInt, " Amount ", monthlyActualAmount[monthInt])
		} else {
			monthlyAmount := (toolingSprintInfo.ActualAmount * float64(projectInfo.AllocatedBudget)) - projectInfo.LastYearAmountUsed
			for i := monthInt - 1; i > 0; i-- {
				monthlyAmount -= monthlyActualAmount[i]
			}
			fmt.Println("Sprint month ", monthInt, " Amount ", monthlyActualAmount[monthInt])
			monthlyActualAmount[monthInt] = monthlyAmount
		}
	}

	return monthlyActualAmount
}

func calcForcastAmountMonthly(currentYearSprintGeneralObject *[]component.GeneralObject, projectInfo ToolingProjectInfo) map[int]float64 {
	monthlyForcastAmount := make(map[int]float64)

	for _, sprintObject := range *currentYearSprintGeneralObject {
		toolingSprint := ToolingProjectSprint{Id: sprintObject.Id, ObjectInfo: sprintObject.ObjectInfo}
		toolingSprintInfo := toolingSprint.getToolingProjectSprintInfo()

		startDate := util.ConvertStringToDateTime(toolingSprintInfo.StartDate)
		_, sprintMonth, _ := startDate.DateTime.Date()
		monthInt := int(sprintMonth)

		if sprintMonth == 1 {
			monthlyForcastAmount[monthInt] = projectInfo.AllocatedBudget * (toolingSprintInfo.ForcastAmount - float64(projectInfo.CompletionBF))
			fmt.Println("Sprint mont ", monthInt, " Amount ", monthlyForcastAmount[monthInt])
		} else {
			monthlyAmount := (toolingSprintInfo.ForcastAmount * float64(projectInfo.AllocatedBudget)) - projectInfo.LastYearAmountUsed
			for i := monthInt - 1; i > 0; i-- {
				monthlyAmount -= monthlyForcastAmount[i]
			}
			fmt.Println("Sprint month ", monthInt, " Amount ", monthlyForcastAmount[monthInt])
			monthlyForcastAmount[monthInt] = monthlyAmount
		}
	}

	return monthlyForcastAmount
}

// func remove(slice []int, s int) []int {
// 	return append(slice[:s], slice[s+1:]...)
// }
