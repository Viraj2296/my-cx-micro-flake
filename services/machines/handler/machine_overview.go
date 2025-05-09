package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type GenericDB struct {
	Name    string
	Records int
}

type Machinestats struct {
	Id      int
	Name    string
	NextRun string
	Image   string
	Status  string
}

type RunDetails struct {
	Type       string `json:"type"` // e.g., "production", "test", "maintenance"
	StartTime  string `json:"startTime"`
	Data       string `json:"data"`       // example "Order 345324243"
	ResourceID string `json:"resourceId"` // Unique identifier for navigation
	Component  string `json:"component"`  // Front-end component to route (e.g., "ProductionDetails", "TestView", "MaintenancePage")
	ColorType  string `json:"colorType"`
}

type MachineStatsResponse struct {
	Name       string      `json:"name"`
	CurrentRun *RunDetails `json:"currentRun,omitempty"` // Pointer to omit when empty
	NextRun    *RunDetails `json:"nextRun,omitempty"`    // Pointer to omit when empty
	Image      string      `json:"image"`
	Status     string      `json:"status"`
	Color      string      `json:"color"`
}

func getNumberOfRecords(dbConnection *gorm.DB, query string) map[string]interface{} {
	var dataResponse = make(map[string]interface{}, 1)
	var numberOfRecords int
	dbConnection.Raw(query).Scan(&numberOfRecords)
	dataResponse["data"] = numberOfRecords
	return dataResponse
}

func getMachineSummary(dbConnection *gorm.DB, query string) map[string]interface{} {
	var dataResponse = make(map[string]interface{}, 1)
	genericDB := GenericDB{}
	dbConnection.Raw(query).Scan(&genericDB)
	dataResponse["data"] = genericDB.Records
	return dataResponse
}

func getBackgroundColour(dbConnection *gorm.DB, id int, table string) string {
	err, objectInterface := Get(dbConnection, table, id)
	if err != nil {
		return "#49C4ED"
	} else {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
		return util.InterfaceToString(objectFields["colorCode"])
	}

}

func getKPIData(value interface{}, label string) component.OverviewData {
	var arrayResponse []map[string]interface{}
	var numberOfUsersData = make(map[string]interface{}, 0)
	numberOfUsersData["v1"] = value
	arrayResponse = append(arrayResponse, numberOfUsersData)

	return component.OverviewData{
		Value:           arrayResponse,
		IsVisible:       true,
		Label:           label,
		Icon:            "bx:task",
		BackgroundColor: "#49C4ED",
	}

}
func (v *MachineService) getOverview(ctx *gin.Context) {
	var overviewResponse = make(map[string]interface{})
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)

	var err error
	listOfMachines, err := GetObjects(dbConnection, MachineMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	listOfAssemblyMachines, err := GetObjects(dbConnection, AssemblyMachineMasterTable)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	listOfToolingMachines, err := GetObjects(dbConnection, ToolingMachineMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	err, historyListResponse := productionOrderInterface.GetNumberOfOdersByStatus(projectId)
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	newData := getScheduleOderData(historyListResponse, dbConnection)
	mouldingOverviewData := getMachineData(listOfMachines, dbConnection)
	assemblyOverviewData := getMachineData(listOfAssemblyMachines, dbConnection)
	toolingOverviewData := getMachineData(listOfToolingMachines, dbConnection)

	mouldingOverviewData = append(mouldingOverviewData, newData...)

	overviewResponse["mouldOverview"] = [1]component.OverviewResponse{{
		Data:  mouldingOverviewData,
		Label: "Machines Summary",
	}}

	overviewResponse["assemblyOverview"] = [1]component.OverviewResponse{{
		Data:  assemblyOverviewData,
		Label: "Assembly Summary",
	}}

	overviewResponse["toolingOverview"] = [1]component.OverviewResponse{{
		Data:  toolingOverviewData,
		Label: "Tooling Summary",
	}}

	ctx.JSON(http.StatusOK, overviewResponse)

}

func getMachineData(listOfMachines *[]component.GeneralObject, dbConnection *gorm.DB) []component.OverviewData {
	var numberOfMachines = len(*listOfMachines)
	var noOfActive int
	var noOfMaintenance int
	var noOfRepair int
	var noOfInactive int

	for _, machineInterface := range *listOfMachines {
		machineMaster := make(map[string]interface{})
		json.Unmarshal(machineInterface.ObjectInfo, &machineMaster)
		if util.InterfaceToInt(machineMaster["machineStatus"]) == 1 {
			noOfActive = noOfActive + 1
		}
		if util.InterfaceToInt(machineMaster["machineStatus"]) == 2 {
			noOfMaintenance = noOfMaintenance + 1
		}
		if util.InterfaceToInt(machineMaster["machineStatus"]) == 3 {
			noOfRepair = noOfRepair + 1
		}
		if util.InterfaceToInt(machineMaster["machineStatus"]) == 4 {
			noOfInactive = noOfInactive + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)
	totalMachines := getKPIData(numberOfMachines, "Total Machines")
	active := getKPIData(noOfActive, "Active")
	maintenance := getKPIData(noOfMaintenance, "Maintenance")
	repair := getKPIData(noOfRepair, "Repair")
	inactive := getKPIData(noOfInactive, "Inactive")

	overviewData = append(overviewData, totalMachines)
	active.BackgroundColor = getBackgroundColour(dbConnection, 1, MachineStatusTable)
	overviewData = append(overviewData, active)
	maintenance.BackgroundColor = getBackgroundColour(dbConnection, 2, MachineStatusTable)
	overviewData = append(overviewData, maintenance)
	repair.BackgroundColor = getBackgroundColour(dbConnection, 3, MachineStatusTable)
	overviewData = append(overviewData, repair)
	inactive.BackgroundColor = getBackgroundColour(dbConnection, 4, MachineStatusTable)
	overviewData = append(overviewData, inactive)

	return overviewData
}

func getScheduleOderData(listOfMachines *[]component.GeneralObject, dbConnection *gorm.DB) []component.OverviewData {

	var noOfScheduled int
	var noOfConfirmed int

	for _, machineInterface := range *listOfMachines {

		machineMaster := make(map[string]interface{})
		json.Unmarshal(machineInterface.ObjectInfo, &machineMaster)

		if util.InterfaceToInt(machineMaster["eventStatus"]) == 3 {
			noOfScheduled = noOfScheduled + 1
		}
		if util.InterfaceToInt(machineMaster["eventStatus"]) == 4 {
			noOfConfirmed = noOfConfirmed + 1
		}

	}

	overviewData := make([]component.OverviewData, 0)

	scheduled := getKPIData(noOfScheduled, "Scheduled Orders")
	confirmed := getKPIData(noOfConfirmed, "Confirmed Orders")

	scheduled.BackgroundColor = getBackgroundColour(dbConnection, 1, MachineStatusTable)
	overviewData = append(overviewData, scheduled)
	confirmed.BackgroundColor = getBackgroundColour(dbConnection, 2, MachineStatusTable)
	overviewData = append(overviewData, confirmed)

	return overviewData
}

func (v *MachineService) getMachineStatsSummary(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)

	// get the user's department, and get the display enabed machines, and get all these machines only
	//GetDepartmentDisplayEnabledMachines
	var err error
	listOfMachines, err := GetObjects(dbConnection, MachineMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	listOfAssemblyMachines, err := GetObjects(dbConnection, AssemblyMachineMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}

	listOfToolingMachines, err := GetObjects(dbConnection, ToolingMachineMasterTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Server Exception",
				Description: "Your action can not be processed due to internal server error, please report this error code to system admin",
			})
		return
	}
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	mouldStatSummary := getMachineStatSummary(listOfMachines, productionOrderInterface, MachineMasterTable, projectId)
	assemblyStatSummary := getMachineStatSummary(listOfAssemblyMachines, productionOrderInterface, AssemblyMachineMasterTable, projectId)
	toolingStatSummary := getMachineStatSummary(listOfToolingMachines, productionOrderInterface, ToolingMachineMasterTable, projectId)

	machineStatsList := make(map[string]interface{})
	machineStatsList["mouldStatSummary"] = mouldStatSummary
	machineStatsList["assemblyStatSummary"] = assemblyStatSummary
	machineStatsList["toolingStatSummary"] = toolingStatSummary
	ctx.JSON(http.StatusOK, machineStatsList)
}

func getMachineStatSummary(listOfMachines *[]component.GeneralObject, productionOrderInterface common.ProductionOrderInterface, targetTable, projectId string) []interface{} {

	machineStatusMap := make(map[int][]MachineStatsResponse)
	for _, machineInterface := range *listOfMachines {

		var currentOrderStatus string
		var currentStatusColor string

		var currentRun *RunDetails
		var nextRun *RunDetails

		masterInfo := make(map[string]interface{})
		json.Unmarshal(machineInterface.ObjectInfo, &masterInfo)

		switch targetTable {
		case MachineMasterTable:
			err, scheduledEventObject := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineInterface.Id)
			mouldService := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
			maintenanceService := common.GetService("maintenance_module").ServiceInterface.(common.MaintenanceInterface)
			_, mouldName := mouldService.GetMouldNameByMachineId(machineInterface.Id)
			correctiveErr, correctiveName := maintenanceService.GetMaintenanceWorkOrderByMachineId(machineInterface.Id, "maintenance_corrective_work_order")
			_, preventiveName := maintenanceService.GetMaintenanceWorkOrderByMachineId(machineInterface.Id, "maintenance_work_order")

			if err == nil {

				currentScheduledEventInfo := GetScheduledOrderEventInfo(scheduledEventObject.ObjectInfo)
				// currentOrderName = currentScheduledEventInfo.Name
				if currentScheduledEventInfo.EventStatus == ScheduleStatusPreferenceFour {

					currentRun = GetRunDetailsValue(currentScheduledEventInfo.EventType, currentScheduledEventInfo.Name, util.InterfaceToString(currentScheduledEventInfo.EventSourceId), "MachineData", "success")
					currentOrderStatus = productionOrderInterface.OrderStatusId2String(projectId, currentScheduledEventInfo.EventStatus)
					_, orderStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, currentScheduledEventInfo.EventStatus)
					orderStatusInfo := make(map[string]interface{})
					json.Unmarshal(orderStatusObject.ObjectInfo, &orderStatusInfo)
					currentStatusColor = util.InterfaceToString(orderStatusInfo["colorCode"])
					err, nextScheduledEventObject := productionOrderInterface.GetNextToCurrentScheduledEvent(projectId, machineInterface.Id, scheduledEventObject.Id)

					if err == nil {
						nextScheduledEventInfo := GetScheduledOrderEventInfo(nextScheduledEventObject.ObjectInfo)
						// nextRunOrderName = nextScheduledEventInfo.Name
						nextRun = GetRunDetailsValue(currentScheduledEventInfo.EventType, nextScheduledEventInfo.Name, util.InterfaceToString(currentScheduledEventInfo.EventSourceId), "MachineData", "success")
					}
				} else if currentScheduledEventInfo.EventStatus == ScheduleStatusPreferenceFive {

					currentRun = GetRunDetailsValue(currentScheduledEventInfo.EventType, currentScheduledEventInfo.Name, util.InterfaceToString(currentScheduledEventInfo.EventSourceId), "MachineData", "success")
					currentOrderStatus = productionOrderInterface.OrderStatusId2String(projectId, currentScheduledEventInfo.EventStatus)
					_, orderStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, currentScheduledEventInfo.EventStatus)
					orderStatusInfo := make(map[string]interface{})
					json.Unmarshal(orderStatusObject.ObjectInfo, &orderStatusInfo)
					currentStatusColor = util.InterfaceToString(orderStatusInfo["colorCode"])
					err, nextScheduledEventObject := productionOrderInterface.GetNextToCurrentScheduledEvent(projectId, machineInterface.Id, scheduledEventObject.Id)

					if err == nil {
						nextScheduledEventInfo := GetScheduledOrderEventInfo(nextScheduledEventObject.ObjectInfo)
						// nextRunOrderName = nextScheduledEventInfo.Name
						nextRun = GetRunDetailsValue(currentScheduledEventInfo.EventType, nextScheduledEventInfo.Name, util.InterfaceToString(currentScheduledEventInfo.EventSourceId), "MachineData", "success")
					}
				} else if currentScheduledEventInfo.EventStatus == ScheduleStatusPreferenceThree {
					if mouldName != "" {
						currentRun = GetRunDetailsValue("mould_test_request", mouldName, "", "MouldData", "info")
					}
				} else {
					if (correctiveErr != nil || correctiveName == "") && preventiveName != "" {
						currentRun = GetRunDetailsValue("maintenance_preventive_work_order", preventiveName, "", "MouldData", "warn")
					} else if correctiveName != "" {
						currentRun = GetRunDetailsValue("maintenance_corrective_work_order", correctiveName, "", "MouldData", "warn")
					}
				}
			} else {
				if currentRun == nil && nextRun == nil {
					if mouldName == "" {
						if (correctiveErr != nil || correctiveName == "") && preventiveName != "" {
							currentRun = GetRunDetailsValue("maintenance_preventive_work_order", preventiveName, "", "MouldData", "warn")
						} else if correctiveName != "" {
							currentRun = GetRunDetailsValue("maintenance_corrective_work_order", correctiveName, "", "MouldData", "warn")
						}
					} else {
						currentRun = GetRunDetailsValue("mould_test_request", mouldName, "", "MouldData", "info")
					}

				}
			}
		case AssemblyMachineMasterTable:
			err, scheduledEventObject := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineInterface.Id)

			if err == nil {
				currentScheduledEventInfo := make(map[string]interface{})
				json.Unmarshal(scheduledEventObject.ObjectInfo, &currentScheduledEventInfo)

				// currentOrderName = util.InterfaceToString(currentScheduledEventInfo["name"])
				currentRun = &RunDetails{
					Type:       util.InterfaceToString(currentScheduledEventInfo["eventType"]),
					StartTime:  util.InterfaceToString(util.GetCurrentTime("2006-01-02T15:04:05")),
					Data:       util.InterfaceToString(currentScheduledEventInfo["name"]),
					ResourceID: util.InterfaceToString(currentScheduledEventInfo["eventSourceId"]),
					Component:  "AssemblyData",
					ColorType:  "danger",
				}
				currentOrderStatus = productionOrderInterface.OrderStatusId2String(projectId, util.InterfaceToInt(currentScheduledEventInfo["eventStatus"]))
				err, nextScheduledEventObject := productionOrderInterface.GetNextToCurrentAssemblyScheduledEvent(projectId, machineInterface.Id, scheduledEventObject.Id)

				_, orderStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, util.InterfaceToInt(currentScheduledEventInfo["eventStatus"]))
				orderStatusInfo := make(map[string]interface{})
				json.Unmarshal(orderStatusObject.ObjectInfo, &orderStatusInfo)
				currentStatusColor = util.InterfaceToString(orderStatusInfo["colorCode"])
				if err == nil {
					nextScheduledEventInfo := make(map[string]interface{})
					json.Unmarshal(nextScheduledEventObject.ObjectInfo, &nextScheduledEventInfo)
					// nextRunOrderName = util.InterfaceToString(nextScheduledEventInfo["name"])
					nextRun = &RunDetails{
						Type:       util.InterfaceToString(currentScheduledEventInfo["eventType"]),
						StartTime:  util.InterfaceToString(util.GetCurrentTime("2006-01-02T15:04:05")),
						Data:       util.InterfaceToString(nextScheduledEventInfo["name"]),
						ResourceID: util.InterfaceToString(currentScheduledEventInfo["eventSourceId"]),
						Component:  "AssemblyData",
						ColorType:  "danger",
					}
				}
			}
		case ToolingMachineMasterTable:
			err, scheduledEventObject := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineInterface.Id)

			if err == nil {
				currentScheduledEventInfo := make(map[string]interface{})
				json.Unmarshal(scheduledEventObject.ObjectInfo, &currentScheduledEventInfo)

				// currentOrderName = util.InterfaceToString(currentScheduledEventInfo["name"])
				currentRun = &RunDetails{
					Type:       util.InterfaceToString(currentScheduledEventInfo["eventType"]),
					StartTime:  util.InterfaceToString(util.GetCurrentTime("2006-01-02T15:04:05")),
					Data:       util.InterfaceToString(currentScheduledEventInfo["name"]),
					ResourceID: util.InterfaceToString(currentScheduledEventInfo["eventSourceId"]),
					Component:  "ToolingData",
					ColorType:  "contrast",
				}
				currentOrderStatus = productionOrderInterface.OrderStatusId2String(projectId, util.InterfaceToInt(currentScheduledEventInfo["eventStatus"]))
				err, nextScheduledEventObject := productionOrderInterface.GetNextToCurrentToolingScheduledEvent(projectId, machineInterface.Id, scheduledEventObject.Id)

				_, orderStatusObject := productionOrderInterface.GetProductionOrderStatus(projectId, util.InterfaceToInt(currentScheduledEventInfo["eventStatus"]))
				orderStatusInfo := make(map[string]interface{})
				json.Unmarshal(orderStatusObject.ObjectInfo, &orderStatusInfo)
				currentStatusColor = util.InterfaceToString(orderStatusInfo["colorCode"])

				if err == nil {
					nextScheduledEventInfo := make(map[string]interface{})
					json.Unmarshal(nextScheduledEventObject.ObjectInfo, &nextScheduledEventInfo)
					// nextRunOrderName = util.InterfaceToString(nextScheduledEventInfo["name"])
					nextRun = &RunDetails{
						Type:       util.InterfaceToString(currentScheduledEventInfo["eventType"]),
						StartTime:  util.InterfaceToString(util.GetCurrentTime("2006-01-02T15:04:05")),
						Data:       util.InterfaceToString(nextScheduledEventInfo["name"]),
						ResourceID: util.InterfaceToString(currentScheduledEventInfo["eventSourceId"]),
						Component:  "ToolingData",
						ColorType:  "contrast",
					}
				}
			}

		}

		machineStatsResponse := MachineStatsResponse{
			Name:       util.InterfaceToString(masterInfo["newMachineId"]),
			CurrentRun: currentRun,
			NextRun:    nextRun,
			Image:      util.InterfaceToString(masterInfo["machineImage"]),
			Status:     currentOrderStatus,
			Color:      currentStatusColor,
		}
		machineStatusMap[util.InterfaceToInt(masterInfo["machineStatus"])] = append(machineStatusMap[util.InterfaceToInt(masterInfo["machineStatus"])], machineStatsResponse)

	}

	machineStatsList := make([]interface{}, 0)
	for key, object := range machineStatusMap {
		groupByView := component.GroupByView{}
		if key == MachineStatusActive {
			groupByView.GroupByField = "active"
			groupByView.DisplayField = "ACTIVE"
		} else if key == MachineStatusMaintenance {
			groupByView.GroupByField = "maintenance"
			groupByView.DisplayField = "MAINTENANCE"
		} else if key == MachineStatusRepair {
			groupByView.GroupByField = "repair"
			groupByView.DisplayField = "REPAIR"
		} else if key == MachineStatusInactive {
			groupByView.GroupByField = "INACTIVE"
			groupByView.DisplayField = "INACTIVE"
		}

		groupByView.Cards = object

		machineStatsList = append(machineStatsList, groupByView)
	}
	return machineStatsList
}

func GetRunDetailsValue(typeValue string, data string, resourceId string, componentName string, color string) *RunDetails {
	currentRun := &RunDetails{
		Type:       typeValue,
		StartTime:  util.InterfaceToString(util.GetCurrentTime("2006-01-02T15:04:05")),
		Data:       data,
		ResourceID: resourceId,
		Component:  componentName,
		ColorType:  color,
	}
	return currentRun
}

//func (ms *MachineService) getStopTimeSummary(ctx *gin.Context) {
//	projectId := util.GetProjectId(ctx)
//	componentName := ctx.Param("componentName")
//
//	targetTable := ms.ComponentManager.GetTargetTable(componentName)
//	offsetValue := ctx.Query("offset")
//	limitValue := ctx.Query("limit")
//	fields := ctx.Query("fields")
//	values := ctx.Query("values")
//	condition := ctx.Query("condition")
//
//	dbConnection := ms.BaseService.ServiceDatabases[projectId]
//
//	var listOfObjects *[]component.GeneralObject
//	var totalRecords int64
//	var err error
//	statusCondition := " object_info ->> '$.hmiStatus' == 'started' OR object_info ->> '$.hmiStatus' == 'stopped'"
//	totalRecords = CountByCondition(dbConnection, targetTable, statusCondition)
//
//	orderBy := "object_info ->> '$.createdAt' desc"
//	if limitValue == "" {
//		listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), orderBy)
//	} else {
//		limitVal, _ := strconv.Atoi(limitValue)
//		listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), orderBy, limitVal)
//	}
//
//	for _, generalObj := range *listOfObjects {
//		machineHmi := make(map[string]interface{})
//		json.Unmarshal(generalObj.ObjectInfo, &machineHmi)
//
//	}
//
//}

func (v *MachineService) getMachineSummary(ctx *gin.Context) {
	projectId := util.GetProjectId(ctx)
	userId := common.GetUserId(ctx)

	factoryService := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	userInfo := authService.GetUserInfoById(userId)

	departmentList := userInfo.Department

	var dropDownArray []component.OrderedData
	var listOfDepartments []int

	for _, departmentId := range departmentList {
		id := departmentId
		dropdownValue := factoryService.GetDepartmentName(departmentId)
		dropDownArray = append(dropDownArray, component.OrderedData{
			Id:    id,
			Value: dropdownValue,
		})

		listOfDepartments = append(listOfDepartments, id)

		for _, secondaryDepartmentId := range userInfo.SecondaryDepartmentList {
			var isExist bool
			isExist = false
			listOfDepartments = append(listOfDepartments, secondaryDepartmentId)
			for _, dropDownId := range dropDownArray {
				if dropDownId.Id == secondaryDepartmentId {
					isExist = true
				}
			}
			if !isExist {
				dropdownValue = factoryService.GetDepartmentName(secondaryDepartmentId)

				dropDownArray = append(dropDownArray, component.OrderedData{
					Id:    secondaryDepartmentId,
					Value: dropdownValue,
				})
			}

		}
	}

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var err error

	tvDisplayResponse := v.GetDepartmentDisplayEnabledMachines(ProjectID, listOfDepartments, dropDownArray)
	displayResponse := make(map[string]interface{})
	json.Unmarshal(tvDisplayResponse, &displayResponse)

	displayEnabledMachines := displayResponse["displayEnabledMachines"].([]interface{})

	idList := "("
	for _, enabledMachines := range displayEnabledMachines {
		enabledMachine := enabledMachines.(map[string]interface{})
		listOfMachines := enabledMachine["listOfMachines"].([]interface{})

		for _, machine := range listOfMachines {
			machineId := util.InterfaceToInt(machine)
			id := strconv.Itoa(machineId)
			idList += id + ","
		}
	}

	idList = util.TrimSuffix(idList, ",") + ")"
	searchQuery := " id in " + idList

	listOfSubcategoryObjects, _ := GetObjects(dbConnection, MachineSubCategoryTable)
	var cacheData = make(map[int]string, 0)
	for _, groupByObject := range *listOfSubcategoryObjects {
		machineCategory := MachineSubCategory{ObjectInfo: groupByObject.ObjectInfo}
		cacheData[groupByObject.Id] = machineCategory.getSubCategoryInfo().Name
	}

	listOfMachineMasterObjects, err := GetConditionalObjects(dbConnection, MachineMasterTable, searchQuery)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}

	var dd []GroupByCardView
	for _, machineMasterObject := range *listOfMachineMasterObjects {
		machineId := machineMasterObject.Id
		machineDashboardQuery := "select * from machine_statistics where machine_id = " + strconv.Itoa(machineId) + " order by ts desc limit 2"
		//dbConnection := ms.BaseService.ServiceDatabases[projectId]
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		var machineDashboardList []MachineStatistics
		dbConnection.Raw(machineDashboardQuery).Scan(&machineDashboardList)
		_, machineObject := Get(dbConnection, MachineMasterTable, machineId)
		machineMater := MachineMaster{ObjectInfo: machineObject.ObjectInfo}
		var currentMachineStats MachineStatisticsInfo
		var previousMachineStats MachineStatisticsInfo

		subCategoryId := machineMater.getMachineMasterInfo().SubCategory

		if len(machineDashboardList) == 2 {
			_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
			_ = json.Unmarshal(machineDashboardList[1].StatsInfo, &previousMachineStats)
		} else if len(machineDashboardList) == 1 {
			_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
			previousMachineStats = MachineStatisticsInfo{}
		} else {
			currentMachineStats = MachineStatisticsInfo{}
			previousMachineStats = MachineStatisticsInfo{}
		}
		_, productionOrderObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, currentMachineStats.ProductionOrderId, machineId)
		productionOrderInfo := GetProductionOrderInfo(productionOrderObject.ObjectInfo)
		mouldModuleInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
		_, partObject := mouldModuleInterface.GetPartInfo(projectId, productionOrderInfo.PartNumber)
		partInfo := GetPartInfo(partObject.ObjectInfo)

		var hmiStartUserName string
		var hmiStartUserAvatarUrl string

		var hmiStopUserName string
		var hmiStopUserAvatarUrl string

		orderBy := "object_info ->> '$.createdAt' desc"
		hmiConditionString := "object_info ->> '$.hmiStatus'='started' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
		listMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, AssemblyMachineHmiTable, hmiConditionString, orderBy)

		if len(*listMachineHmi) > 0 {
			machineHmiInfo := make(map[string]interface{})
			json.Unmarshal((*listMachineHmi)[0].ObjectInfo, machineHmiInfo)

			hmiUserInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
			if hmiUserInfo.FullName != "" {
				hmiStartUserName = hmiUserInfo.FullName
				hmiStartUserAvatarUrl = hmiUserInfo.AvatarUrl
			}

			// get the stopped by
			hmiLastStopConditionString := "object_info ->> '$.hmiStatus'='stopped' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
			listOfStoppedMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, AssemblyMachineHmiTable, hmiLastStopConditionString, orderBy)
			if len(*listOfStoppedMachineHmi) > 0 {
				machineHmiInfo = make(map[string]interface{})
				json.Unmarshal((*listOfStoppedMachineHmi)[0].ObjectInfo, machineHmiInfo)

				hmiUserInfo = authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
				if hmiUserInfo.FullName != "" {
					hmiStopUserName = hmiUserInfo.FullName
					hmiStopUserAvatarUrl = hmiUserInfo.AvatarUrl
				}
			}

		}

		machineDashboardResponse := getDashboardResult(currentMachineStats, previousMachineStats, machineMater.getMachineMasterInfo(), partInfo, "", hmiStartUserName, hmiStartUserAvatarUrl, hmiStopUserName, hmiStopUserAvatarUrl)
		rawDashboardResponse, _ := json.Marshal(machineDashboardResponse)

		var cardResponse = make(map[string]interface{})

		json.Unmarshal(rawDashboardResponse, &cardResponse)

		subCategoryValue := util.InterfaceToString(cacheData[subCategoryId])

		var isElementFound bool
		isElementFound = false
		for index, mm := range dd {
			if mm.GroupByField == subCategoryValue {
				dd[index].Cards = append(dd[index].Cards, cardResponse)
				isElementFound = true
			}
		}
		if !isElementFound {
			xl := GroupByCardView{}
			xl.GroupByField = subCategoryValue
			xl.Cards = append(xl.Cards, cardResponse)
			dd = append(dd, xl)
		}
	}

	ctx.JSON(http.StatusOK, dd)

}

func (v *MachineService) getMachineDashboardResponse(projectId string, machineId int) datatypes.JSON {

	machineDashboardQuery := "select * from machine_statistics where machine_id = " + strconv.Itoa(machineId) + " order by ts desc limit 2"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	var machineDashboardList []MachineStatistics
	dbConnection.Raw(machineDashboardQuery).Scan(&machineDashboardList)
	_, machineObject := Get(dbConnection, MachineMasterTable, machineId)
	machineMater := MachineMaster{ObjectInfo: machineObject.ObjectInfo}
	var currentMachineStats MachineStatisticsInfo
	var previousMachineStats MachineStatisticsInfo

	if len(machineDashboardList) == 2 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		_ = json.Unmarshal(machineDashboardList[1].StatsInfo, &previousMachineStats)
	} else if len(machineDashboardList) == 1 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		previousMachineStats = MachineStatisticsInfo{}
	} else {
		currentMachineStats = MachineStatisticsInfo{}
		previousMachineStats = MachineStatisticsInfo{}
	}
	_, productionOrderObject := productionOrderInterface.GetMachineProductionOrderInfo(projectId, currentMachineStats.ProductionOrderId, machineId)
	productionOrderInfo := GetProductionOrderInfo(productionOrderObject.ObjectInfo)
	mouldModuleInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, partObject := mouldModuleInterface.GetPartInfo(projectId, productionOrderInfo.PartNumber)
	partInfo := GetPartInfo(partObject.ObjectInfo)

	err, currentScheduleEvent := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineId)

	if err == nil {
		scheduledEventInfo := make(map[string]interface{})
		json.Unmarshal(currentScheduleEvent.ObjectInfo, &scheduledEventInfo)

		orderStatus := util.InterfaceToInt(scheduledEventInfo["eventStatus"])
		scheduleOrderStatus := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if orderStatus <= scheduleOrderStatus && machineMater.getMachineMasterInfo().MachineConnectStatus == machineConnectStatusLive && currentMachineStats.EventId != currentScheduleEvent.Id {
			currentMachineStats = MachineStatisticsInfo{}
			previousMachineStats = MachineStatisticsInfo{}
			currentMachineStats.EventId = currentScheduleEvent.Id
		}
	}

	// _, eventId := getCurrentScheduleDate(timeLineEvents)
	// actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId)
	var maintenanceColorCode string

	if machineMater.getMachineMasterInfo().MachineStatus != 0 {
		_, statusObject := Get(dbConnection, MachineStatusTable, machineMater.getMachineMasterInfo().MachineStatus)
		machineStatus := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, &machineStatus)
		maintenanceColorCode = util.InterfaceToString(machineStatus["colorCode"])
	}

	var hmiStartUserName string
	var hmiStartUserAvatarUrl string
	var hmiStopUserName string
	var hmiStopUserAvatarUrl string
	orderBy := "object_info ->> '$.createdAt' desc"
	hmiConditionString := "object_info ->> '$.hmiStatus'='started' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
	listMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, MachineHMITable, hmiConditionString, orderBy)
	if len(*listMachineHmi) > 0 {
		machineHmiInfo := make(map[string]interface{})
		json.Unmarshal((*listMachineHmi)[0].ObjectInfo, &machineHmiInfo)
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
		if userInfo.FullName != "" {
			hmiStartUserName = userInfo.FullName
			hmiStartUserAvatarUrl = userInfo.AvatarUrl
		}

		// get the last stopped by username, and avatar details.
		hmiLastStoppedByConditionString := "object_info ->> '$.hmiStatus'='stopped' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
		listOfLastStopByMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, MachineHMITable, hmiLastStoppedByConditionString, orderBy)
		if len(*listOfLastStopByMachineHmi) > 0 {
			machineHmiInfo = make(map[string]interface{})
			json.Unmarshal((*listOfLastStopByMachineHmi)[0].ObjectInfo, &machineHmiInfo)
			userInfo = authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
			if userInfo.FullName != "" {
				hmiStopUserName = userInfo.FullName
				hmiStopUserAvatarUrl = userInfo.AvatarUrl
			}
		}
	}

	machineDashboardResponse := getDashboardResult(currentMachineStats, previousMachineStats, machineMater.getMachineMasterInfo(), partInfo, maintenanceColorCode, hmiStartUserName, hmiStartUserAvatarUrl, hmiStopUserName, hmiStopUserAvatarUrl)
	rawDashboardResponse, _ := json.Marshal(machineDashboardResponse)
	return rawDashboardResponse
}

func (v *MachineService) getAssemblyMachineDashboardResponse(projectId string, machineId int) datatypes.JSON {

	machineDashboardQuery := "select * from assembly_machine_statistics where machine_id = " + strconv.Itoa(machineId) + " order by ts desc limit 2"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	var machineDashboardList []AssemblyMachineStatistics
	dbConnection.Raw(machineDashboardQuery).Scan(&machineDashboardList)
	_, machineObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	machineMater := AssemblyMachineMaster{ObjectInfo: machineObject.ObjectInfo}
	var currentMachineStats MachineStatisticsInfo
	var previousMachineStats MachineStatisticsInfo

	if len(machineDashboardList) == 2 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		_ = json.Unmarshal(machineDashboardList[1].StatsInfo, &previousMachineStats)
	} else if len(machineDashboardList) == 1 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		previousMachineStats = MachineStatisticsInfo{}
	} else {
		currentMachineStats = MachineStatisticsInfo{}
		previousMachineStats = MachineStatisticsInfo{}
	}
	_, productionOrderObject := productionOrderInterface.GetAssemblyProductionOrderInfo(projectId, currentMachineStats.ProductionOrderId, machineId)
	productionOrderInfo := GetProductionOrderInfo(productionOrderObject.ObjectInfo)
	manufacturingModuleInterface := common.GetService("manufacturing_module").ServiceInterface.(common.ManufacturingInterface)
	_, partObject := manufacturingModuleInterface.GetAssemblyPartInfo(projectId, productionOrderInfo.PartNumber)
	partInfo := GetPartInfo(partObject.ObjectInfo)

	err, currentScheduleEvent := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineId)

	if err == nil {
		scheduledEventInfo := make(map[string]interface{})
		json.Unmarshal(currentScheduleEvent.ObjectInfo, &scheduledEventInfo)

		orderStatus := util.InterfaceToInt(scheduledEventInfo["eventStatus"])
		scheduleOrderStatus := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if orderStatus <= scheduleOrderStatus && machineMater.getAssemblyMachineMasterInfo().MachineConnectStatus == machineConnectStatusLive && currentMachineStats.EventId != currentScheduleEvent.Id {
			currentMachineStats = MachineStatisticsInfo{}
			previousMachineStats = MachineStatisticsInfo{}
			currentMachineStats.EventId = currentScheduleEvent.Id
		}
	}

	// _, eventId := getCurrentScheduleDate(timeLineEvents)
	// actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId)
	var maintenanceColorCode string

	if machineMater.getAssemblyMachineMasterInfo().MachineStatus != 0 {
		_, statusObject := Get(dbConnection, MachineStatusTable, machineMater.getAssemblyMachineMasterInfo().MachineStatus)
		machineStatus := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, machineStatus)
		maintenanceColorCode = util.InterfaceToString(machineStatus["colorCode"])
	}

	orderBy := "object_info ->> '$.createdAt' desc"
	var hmiStartUserName string
	var hmiStartUserAvatarUrl string
	var hmiStopUserName string
	var hmiStopUserAvatarUrl string
	hmiConditionString := "object_info ->> '$.hmiStatus'='started' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
	listMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, ToolingMachineHmiTable, hmiConditionString, orderBy)
	if len(*listMachineHmi) > 0 {
		machineHmiInfo := make(map[string]interface{})
		json.Unmarshal((*listMachineHmi)[0].ObjectInfo, &machineHmiInfo)
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
		if userInfo.FullName != "" {
			hmiStartUserName = userInfo.FullName
			hmiStartUserAvatarUrl = userInfo.AvatarUrl
		}

		// get the last stopped by username, and avatar details.
		hmiLastStoppedByConditionString := "object_info ->> '$.hmiStatus'='stopped' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
		listOfLastStopByMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, ToolingMachineHmiTable, hmiLastStoppedByConditionString, orderBy)
		if len(*listOfLastStopByMachineHmi) > 0 {
			machineHmiInfo = make(map[string]interface{})
			json.Unmarshal((*listOfLastStopByMachineHmi)[0].ObjectInfo, &machineHmiInfo)
			userInfo = authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
			if userInfo.FullName != "" {
				hmiStopUserName = userInfo.FullName
				hmiStopUserAvatarUrl = userInfo.AvatarUrl
			}
		}
	}
	machineDashboardResponse := getAssemblyDashboardResult(currentMachineStats, previousMachineStats, machineMater.getAssemblyMachineMasterInfo(), partInfo, maintenanceColorCode,
		hmiStartUserName, hmiStartUserAvatarUrl, hmiStopUserName, hmiStopUserAvatarUrl)
	rawDashboardResponse, _ := json.Marshal(machineDashboardResponse)
	return rawDashboardResponse
}

func (v *MachineService) getToolingMachineDashboardResponse(projectId string, machineId int) datatypes.JSON {

	machineDashboardQuery := "select * from tooling_machine_statistics where machine_id = " + strconv.Itoa(machineId) + " order by ts desc limit 2"
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	var machineDashboardList []ToolingMachineStatistics
	dbConnection.Raw(machineDashboardQuery).Scan(&machineDashboardList)
	_, machineObject := Get(dbConnection, ToolingMachineMasterTable, machineId)
	machineMater := ToolingMachineMaster{ObjectInfo: machineObject.ObjectInfo}
	var currentMachineStats MachineStatisticsInfo
	var previousMachineStats MachineStatisticsInfo

	if len(machineDashboardList) == 2 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		_ = json.Unmarshal(machineDashboardList[1].StatsInfo, &previousMachineStats)
	} else if len(machineDashboardList) == 1 {
		_ = json.Unmarshal(machineDashboardList[0].StatsInfo, &currentMachineStats)
		previousMachineStats = MachineStatisticsInfo{}
	} else {
		currentMachineStats = MachineStatisticsInfo{}
		previousMachineStats = MachineStatisticsInfo{}
	}
	_, productionOrderObject := productionOrderInterface.GetToolingScheduledOrderInfo(projectId, currentMachineStats.EventId)
	scheduledOrderInfo := make(map[string]interface{})
	json.Unmarshal(productionOrderObject.ObjectInfo, &scheduledOrderInfo)

	orderId := util.InterfaceToInt(scheduledOrderInfo["eventSourceId"])
	_, bomMaster := productionOrderInterface.GetBomInfo(projectId, orderId)
	bomInfo := make(map[string]interface{})
	json.Unmarshal(bomMaster.ObjectInfo, &bomInfo)
	bomId := util.InterfaceToString(bomInfo["name"])

	partId := util.InterfaceToInt(scheduledOrderInfo["partId"])
	_, partMaster := productionOrderInterface.GetToolingPartById(projectId, partId)

	partInfo := make(map[string]interface{})
	json.Unmarshal(partMaster.ObjectInfo, &partInfo)

	err, currentScheduleEvent := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineId)
	if err == nil {
		scheduledEventInfo := make(map[string]interface{})
		json.Unmarshal(currentScheduleEvent.ObjectInfo, &scheduledEventInfo)

		orderStatus := util.InterfaceToInt(scheduledEventInfo["eventStatus"])
		scheduleOrderStatus := productionOrderInterface.GetOrderStatusIdFromPreferenceLevel(projectId, ScheduleStatusPreferenceThree)

		if orderStatus <= scheduleOrderStatus && machineMater.getToolingMachineMasterInfo().MachineConnectStatus == machineConnectStatusLive && currentMachineStats.EventId != currentScheduleEvent.Id {
			currentMachineStats = MachineStatisticsInfo{}
			previousMachineStats = MachineStatisticsInfo{}
			currentMachineStats.EventId = currentScheduleEvent.Id
		}
	}

	// _, eventId := getCurrentScheduleDate(timeLineEvents)
	// actualStartTime, actualEndTime := getActualStartStopTime(dbConnection, machineId, eventId)

	var maintenanceColorCode string

	if machineMater.getToolingMachineMasterInfo().MachineStatus != 0 {
		_, statusObject := Get(dbConnection, MachineStatusTable, machineMater.getToolingMachineMasterInfo().MachineStatus)
		machineStatus := make(map[string]interface{})
		json.Unmarshal(statusObject.ObjectInfo, machineStatus)
		maintenanceColorCode = util.InterfaceToString(machineStatus["colorCode"])
	}
	orderBy := "object_info ->> '$.createdAt' desc"
	var hmiStartUserName string
	var hmiStartUserAvatarUrl string
	var hmiStopUserName string
	var hmiStopUserAvatarUrl string
	hmiConditionString := "object_info ->> '$.hmiStatus'='started' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
	listMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, ToolingMachineHmiTable, hmiConditionString, orderBy)
	if len(*listMachineHmi) > 0 {
		machineHmiInfo := make(map[string]interface{})
		json.Unmarshal((*listMachineHmi)[0].ObjectInfo, &machineHmiInfo)
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		userInfo := authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
		if userInfo.FullName != "" {
			hmiStartUserName = userInfo.FullName
			hmiStartUserAvatarUrl = userInfo.AvatarUrl
		}

		// get the last stopped by username, and avatar details.
		hmiLastStoppedByConditionString := "object_info ->> '$.hmiStatus'='stopped' and object_info ->> '$.eventId'=" + strconv.Itoa(currentMachineStats.EventId)
		listOfLastStopByMachineHmi, _ := GetConditionalObjectsOrderBy(dbConnection, ToolingMachineHmiTable, hmiLastStoppedByConditionString, orderBy)
		if len(*listOfLastStopByMachineHmi) > 0 {
			machineHmiInfo = make(map[string]interface{})
			json.Unmarshal((*listOfLastStopByMachineHmi)[0].ObjectInfo, &machineHmiInfo)
			userInfo = authService.GetUserInfoById(util.InterfaceToInt(machineHmiInfo["createdBy"]))
			if userInfo.FullName != "" {
				hmiStopUserName = userInfo.FullName
				hmiStopUserAvatarUrl = userInfo.AvatarUrl
			}
		}
	}

	machineDashboardResponse := getToolingDashboardResult(currentMachineStats, previousMachineStats, machineMater.getToolingMachineMasterInfo(), partInfo, bomId, maintenanceColorCode, hmiStartUserName, hmiStartUserAvatarUrl, hmiStopUserName, hmiStopUserAvatarUrl)
	rawDashboardResponse, _ := json.Marshal(machineDashboardResponse)
	return rawDashboardResponse
}

func (v *MachineService) getIntialMachineDashboardResponse(projectId string, userId int) datatypes.JSON {

	factoryService := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	userInfo := authService.GetUserInfoById(userId)
	//disabledMenuIds = append(disabledMenuIds, "it_my_request")
	//disabledMenuIds = append(disabledMenuIds, "tooling_project")
	recordInfo := component.RecordInfo{}
	var dropDownArray []component.OrderedData
	index := 0
	var listOfDepartments []int
	var dropdownValue string

	for _, departmentId := range userInfo.Department {
		id := departmentId
		dropdownValue := factoryService.GetDepartmentName(departmentId)
		dropDownArray = append(dropDownArray, component.OrderedData{
			Id:    id,
			Value: dropdownValue,
		})
		listOfDepartments = append(listOfDepartments, id)

		if index == 0 {
			recordInfo.Index = id
			recordInfo.Value = dropdownValue
		}
		index = index + 1
	}

	for _, secondaryDepartmentId := range userInfo.SecondaryDepartmentList {
		var isExist bool
		isExist = false
		listOfDepartments = append(listOfDepartments, secondaryDepartmentId)
		for _, dropDownId := range dropDownArray {
			if dropDownId.Id == secondaryDepartmentId {
				isExist = true
			}
		}
		if !isExist {
			dropdownValue = factoryService.GetDepartmentName(secondaryDepartmentId)

			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    secondaryDepartmentId,
				Value: dropdownValue,
			})
		}

	}

	tvDisplayResponse := v.GetDepartmentDisplayEnabledMachines(ProjectID, listOfDepartments, dropDownArray)
	return tvDisplayResponse
}

func getPercentageDiff(latestValue float64, previousValue float64) float64 {
	if latestValue != 0 {
		return ((latestValue - previousValue) / latestValue) * 100
	}
	return 0
}

func getDashboardResult(machineDashboard MachineStatisticsInfo,
	previousDashboardResult MachineStatisticsInfo,
	machineMasterInfo *MachineMasterInfo,
	partInfo *PartInfo,
	maintenanceColorCode string,
	hmiStartUserName string,
	hmiStartUserAvatarUrl string, hmiStopUserName string,
	hmiStopUserAvatarUrl string) MachineDashboardResponse {

	decimalCorrectedOEE := float64(machineDashboard.Oee) / float64(100)
	decimalCorrectedPerformance := float64(machineDashboard.Performance) / float64(100)
	decimalCorrectedAvailability := float64(machineDashboard.Availability) / float64(100)
	decimalCorrectedQuality := float64(machineDashboard.Quality) / float64(100)

	previousDecimalCorrectedOEE := float64(previousDashboardResult.Oee) / float64(100)
	previousDecimalCorrectedPerformance := float64(previousDashboardResult.Performance) / float64(100)
	previousDecimalCorrectedAvailability := float64(previousDashboardResult.Availability) / float64(100)
	previousDecimalCorrectedQuality := float64(previousDashboardResult.Quality) / float64(100)

	var currentStat string
	underMaintenance := false
	if machineMasterInfo.MachineConnectStatus == machineConnectStatusMaintenance {
		currentStat = "Maintenance"
		underMaintenance = true
	} else {
		if machineMasterInfo.MachineConnectStatus == machineConnectStatusLive {
			currentStat = "Running"
			maintenanceColorCode = "#32602e"
		} else {
			currentStat = "Stopped"
			maintenanceColorCode = "#E60E18"
		}

		if machineDashboard.CurrentStatus == "" {
			currentStat = "Pending"
			maintenanceColorCode = "#A31041"
		}
	}

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, eventObject := productionOrderInterface.GetScheduledOrderInfo(ProjectID, machineDashboard.EventId)

	mouldInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
	_, brandName := mouldInterface.GetBrandName(ProjectID, machineMasterInfo.Brand)

	schedulerEventInfo := GetScheduledOrderEventInfo(eventObject.ObjectInfo)
	oeeDiff := getPercentageDiff(decimalCorrectedOEE, previousDecimalCorrectedOEE)
	availabilityDiff := getPercentageDiff(decimalCorrectedAvailability, previousDecimalCorrectedAvailability)
	performanceDiff := getPercentageDiff(decimalCorrectedPerformance, previousDecimalCorrectedPerformance)
	qualityDiff := getPercentageDiff(decimalCorrectedQuality, previousDecimalCorrectedQuality)
	completedPercentageDiff := getPercentageDiff(float64(machineDashboard.CompletedPercentage), float64(previousDashboardResult.CompletedPercentage))

	machineDashboardResponse := MachineDashboardResponse{
		Actual:                     int(machineDashboard.Actual),
		Brand:                      getDefaultStringValue(brandName),
		MachineImage:               getDefaultStringValue(machineMasterInfo.MachineImage),
		Model:                      getDefaultStringValue(machineMasterInfo.Model),
		NewMachineId:               getDefaultStringValue(machineMasterInfo.NewMachineId),
		CurrentStatus:              currentStat,
		ProductionOrder:            getDefaultStringValue(schedulerEventInfo.Name),
		PartName:                   getDefaultStringValue(partInfo.PartNumber),
		PartDescription:            getDefaultStringValue(partInfo.Description),
		ScheduleStartTime:          util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(schedulerEventInfo.StartDate)),
		ScheduleEndTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(schedulerEventInfo.EndDate)),
		ActualStartTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualStartTime)),
		ActualEndTime:              util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualEndTime)),
		EstimatedEndTime:           util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.EstimatedEndTime)),
		Oee:                        math.Round(decimalCorrectedOEE*100) / 100,
		Availability:               int(math.RoundToEven(decimalCorrectedAvailability)),
		Performance:                int(math.RoundToEven(decimalCorrectedPerformance)),
		Quality:                    int(math.RoundToEven(decimalCorrectedQuality)),
		OeeDiff:                    math.Round(oeeDiff*100) / 100,
		AvailabilityDiff:           math.Round(availabilityDiff*100) / 100,
		PerformanceDiff:            math.Round(performanceDiff*100) / 100,
		QualityDiff:                math.Round(qualityDiff*100) / 100,
		CompletedPercentageDiff:    completedPercentageDiff,
		PlannedQuality:             machineDashboard.PlannedQuality,
		DailyPlannedQuality:        machineDashboard.DailyPlannedQty,
		Completed:                  machineDashboard.Completed,
		Rejects:                    machineDashboard.Rejects,
		CompletedPercentage:        math.Round(float64(machineDashboard.CompletedPercentage)*100) / 100,
		OverallRejectedQty:         machineDashboard.OverallRejectedQty,
		OverallCompletedPercentage: math.Round(float64(machineDashboard.OverallCompletedPercentage)*100) / 100,
		Remark:                     machineDashboard.Remark,
		ColorCode:                  maintenanceColorCode,
		IsUnderMaintenance:         underMaintenance,
		HmiStartUserName:           hmiStartUserName,
		HmiStartUserAvatarUrl:      hmiStartUserAvatarUrl,
		HmiStopUserName:            hmiStopUserName,
		HmiStopUserAvatarUrl:       hmiStopUserAvatarUrl,
	}

	return machineDashboardResponse
}

func getAssemblyDashboardResult(machineDashboard MachineStatisticsInfo,
	previousDashboardResult MachineStatisticsInfo,
	machineMasterInfo *AssemblyMachineMasterInfo,
	partInfo *PartInfo,
	maintenanceColorCode string,
	hmiStartUsername string,
	hmiStartUserAvatarUrl string,
	hmiStopUserName string, hmiStopUserAvatarUrl string) MachineDashboardResponse {

	decimalCorrectedOEE := float64(machineDashboard.Oee) / float64(100)
	decimalCorrectedPerformance := float64(machineDashboard.Performance) / float64(100)
	decimalCorrectedAvailability := float64(machineDashboard.Availability) / float64(100)
	decimalCorrectedQuality := float64(machineDashboard.Quality) / float64(100)

	previousDecimalCorrectedOEE := float64(previousDashboardResult.Oee) / float64(100)
	previousDecimalCorrectedPerformance := float64(previousDashboardResult.Performance) / float64(100)
	previousDecimalCorrectedAvailability := float64(previousDashboardResult.Availability) / float64(100)
	previousDecimalCorrectedQuality := float64(previousDashboardResult.Quality) / float64(100)

	var currentStat string

	underMaintenance := false
	if machineMasterInfo.MachineConnectStatus == machineConnectStatusMaintenance {
		currentStat = "Maintenance"
		underMaintenance = true
	} else {
		if machineMasterInfo.MachineConnectStatus == machineConnectStatusLive {
			currentStat = "Running"
			maintenanceColorCode = "#32602e"
		} else {
			currentStat = "Stopped"
			maintenanceColorCode = "#E60E18"
		}

		if machineDashboard.CurrentStatus == "" {
			currentStat = "Pending"
			maintenanceColorCode = "#A31041"
		}
	}

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, eventObject := productionOrderInterface.GetAssemblyScheduledOrderInfo(ProjectID, machineDashboard.EventId)
	schedulerEventInfo := GetScheduledOrderEventInfo(eventObject.ObjectInfo)
	oeeDiff := getPercentageDiff(decimalCorrectedOEE, previousDecimalCorrectedOEE)
	availabilityDiff := getPercentageDiff(decimalCorrectedAvailability, previousDecimalCorrectedAvailability)
	performanceDiff := getPercentageDiff(decimalCorrectedPerformance, previousDecimalCorrectedPerformance)
	qualityDiff := getPercentageDiff(decimalCorrectedQuality, previousDecimalCorrectedQuality)
	completedPercentageDiff := getPercentageDiff(float64(machineDashboard.CompletedPercentage), float64(previousDashboardResult.CompletedPercentage))

	machineDashboardResponse := MachineDashboardResponse{
		Actual:                     int(machineDashboard.Actual),
		Brand:                      getDefaultStringValue(""),
		MachineImage:               getDefaultStringValue(machineMasterInfo.MachineImage),
		Model:                      getDefaultStringValue(""),
		NewMachineId:               getDefaultStringValue(machineMasterInfo.NewMachineId),
		CurrentStatus:              currentStat,
		ProductionOrder:            getDefaultStringValue(schedulerEventInfo.Name),
		PartName:                   getDefaultStringValue(partInfo.PartNumber),
		PartDescription:            getDefaultStringValue(partInfo.Description),
		ScheduleStartTime:          util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(schedulerEventInfo.StartDate)),
		ScheduleEndTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(schedulerEventInfo.EndDate)),
		ActualStartTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualStartTime)),
		ActualEndTime:              util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualEndTime)),
		EstimatedEndTime:           util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.EstimatedEndTime)),
		Oee:                        math.Round(decimalCorrectedOEE*100) / 100,
		Availability:               int(math.RoundToEven(decimalCorrectedAvailability)),
		Performance:                int(math.RoundToEven(decimalCorrectedPerformance)),
		Quality:                    int(math.RoundToEven(decimalCorrectedQuality)),
		OeeDiff:                    math.Round(oeeDiff*100) / 100,
		AvailabilityDiff:           math.Round(availabilityDiff*100) / 100,
		PerformanceDiff:            math.Round(performanceDiff*100) / 100,
		QualityDiff:                math.Round(qualityDiff*100) / 100,
		CompletedPercentageDiff:    completedPercentageDiff,
		PlannedQuality:             machineDashboard.PlannedQuality,
		DailyPlannedQuality:        machineDashboard.DailyPlannedQty,
		Completed:                  machineDashboard.Completed,
		Rejects:                    machineDashboard.Rejects,
		CompletedPercentage:        math.Round(float64(machineDashboard.CompletedPercentage)*100) / 100,
		OverallRejectedQty:         machineDashboard.OverallRejectedQty,
		OverallCompletedPercentage: math.Round(float64(machineDashboard.OverallCompletedPercentage)*100) / 100,
		Remark:                     machineDashboard.Remark,
		IsUnderMaintenance:         underMaintenance,
		ColorCode:                  maintenanceColorCode,
		HmiStartUserAvatarUrl:      hmiStartUserAvatarUrl,
		HmiStartUserName:           hmiStartUsername,
		HmiStopUserAvatarUrl:       hmiStopUserAvatarUrl,
		HmiStopUserName:            hmiStopUserName,
	}

	return machineDashboardResponse
}

func getToolingDashboardResult(machineDashboard MachineStatisticsInfo,
	previousDashboardResult MachineStatisticsInfo,
	machineMasterInfo *ToolingMachineMasterInfo,
	partInfo map[string]interface{}, bomId string, maintenanceColorCode string,
	hmiStartUserName string,
	hmiStartUserAvatarUrl string, hmiStopUserName string,
	hmiStopUserAvatarUrl string) MachineDashboardResponse {

	decimalCorrectedOEE := float64(machineDashboard.Oee) / float64(100)
	decimalCorrectedPerformance := float64(machineDashboard.Performance) / float64(100)
	decimalCorrectedAvailability := float64(machineDashboard.Availability) / float64(100)
	decimalCorrectedQuality := float64(machineDashboard.Quality) / float64(100)

	previousDecimalCorrectedOEE := float64(previousDashboardResult.Oee) / float64(100)
	previousDecimalCorrectedPerformance := float64(previousDashboardResult.Performance) / float64(100)
	previousDecimalCorrectedAvailability := float64(previousDashboardResult.Availability) / float64(100)
	previousDecimalCorrectedQuality := float64(previousDashboardResult.Quality) / float64(100)

	var currentStat string

	underMaintenance := false
	if machineMasterInfo.MachineConnectStatus == machineConnectStatusMaintenance {
		currentStat = "Maintenance"
		underMaintenance = true
	} else {
		if machineMasterInfo.MachineConnectStatus == machineConnectStatusLive {
			currentStat = "Running"
			maintenanceColorCode = "#32602e"
		} else {
			currentStat = "Stopped"
			maintenanceColorCode = "#E60E18"
		}

		if machineDashboard.CurrentStatus == "" {
			currentStat = "Pending"
			maintenanceColorCode = "#A31041"
		}
	}

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
	_, eventObject := productionOrderInterface.GetToolingScheduledOrderInfo(ProjectID, machineDashboard.EventId)

	schedulerEventInfo := make(map[string]interface{})
	json.Unmarshal(eventObject.ObjectInfo, &schedulerEventInfo)

	oeeDiff := getPercentageDiff(decimalCorrectedOEE, previousDecimalCorrectedOEE)
	availabilityDiff := getPercentageDiff(decimalCorrectedAvailability, previousDecimalCorrectedAvailability)
	performanceDiff := getPercentageDiff(decimalCorrectedPerformance, previousDecimalCorrectedPerformance)
	qualityDiff := getPercentageDiff(decimalCorrectedQuality, previousDecimalCorrectedQuality)
	completedPercentageDiff := getPercentageDiff(float64(machineDashboard.CompletedPercentage), float64(previousDashboardResult.CompletedPercentage))

	stringTimeFormat := strconv.Itoa(machineDashboard.PlannedQuality) + "s"
	plannedQualitySecond, _ := time.ParseDuration(stringTimeFormat)
	plannedQualityHrs := fmt.Sprintf("%.2f", plannedQualitySecond.Hours()) + " hours"

	stringTimeFormatQty := strconv.Itoa(machineDashboard.DailyPlannedQty) + "s"
	dailyPlannedQualitySecond, _ := time.ParseDuration(stringTimeFormatQty)
	dailyPlannedQualityHrs := fmt.Sprintf("%.2f", dailyPlannedQualitySecond.Hours()) + " hours"

	stringTimeFormatActual := strconv.Itoa(machineDashboard.Actual) + "s"
	actualInSeconds, _ := time.ParseDuration(stringTimeFormatActual)
	actualHrs := fmt.Sprintf("%.2f", actualInSeconds.Hours()) + " hours"

	stringTimeFormatCompleted := strconv.Itoa(machineDashboard.Completed) + "s"
	completedInSeconds, _ := time.ParseDuration(stringTimeFormatCompleted)
	completedHrs := fmt.Sprintf("%.2f", completedInSeconds.Hours()) + " hours"

	machineDashboardResponse := MachineDashboardResponse{
		ActualHrs:                  actualHrs,
		Brand:                      getDefaultStringValue(""),
		MachineImage:               getDefaultStringValue(machineMasterInfo.MachineImage),
		Model:                      getDefaultStringValue(""),
		NewMachineId:               getDefaultStringValue(machineMasterInfo.NewMachineId),
		CurrentStatus:              currentStat,
		ProductionOrder:            getDefaultStringValue(util.InterfaceToString(schedulerEventInfo["name"])),
		PartName:                   getDefaultStringValue(util.InterfaceToString(partInfo["name"])),
		PartDescription:            getDefaultStringValue(util.InterfaceToString(partInfo["description"])),
		ScheduleStartTime:          util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(util.InterfaceToString(schedulerEventInfo["startDate"]))),
		ScheduleEndTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(util.InterfaceToString(schedulerEventInfo["endDate"]))),
		ActualStartTime:            util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualStartTime)),
		ActualEndTime:              util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.ActualEndTime)),
		EstimatedEndTime:           util.ConvertTimeToZeroZone("Asia/Singapore", getDefaultDateValue(machineDashboard.EstimatedEndTime)),
		Oee:                        math.Round(decimalCorrectedOEE*100) / 100,
		Availability:               int(math.RoundToEven(decimalCorrectedAvailability)),
		Performance:                int(math.RoundToEven(decimalCorrectedPerformance)),
		Quality:                    int(math.RoundToEven(decimalCorrectedQuality)),
		OeeDiff:                    math.Round(oeeDiff*100) / 100,
		AvailabilityDiff:           math.Round(availabilityDiff*100) / 100,
		PerformanceDiff:            math.Round(performanceDiff*100) / 100,
		QualityDiff:                math.Round(qualityDiff*100) / 100,
		CompletedPercentageDiff:    completedPercentageDiff,
		PlannedQualityHrs:          plannedQualityHrs,
		DailyPlannedQualityHrs:     dailyPlannedQualityHrs,
		CompletedHrs:               completedHrs,
		Rejects:                    machineDashboard.Rejects,
		CompletedPercentage:        math.Round(float64(machineDashboard.CompletedPercentage)*100) / 100,
		OverallRejectedQty:         machineDashboard.OverallRejectedQty,
		OverallCompletedPercentage: math.Round(float64(machineDashboard.OverallCompletedPercentage)*100) / 100,
		RejectsHrs:                 "N/A",
		OverallRejectedHrs:         util.InterfaceToInt(partInfo["totalQuantity"]),
		BomId:                      bomId,
		Remark:                     machineDashboard.Remark,
		IsUnderMaintenance:         underMaintenance,
		ColorCode:                  maintenanceColorCode,
		HmiStartUserName:           hmiStartUserName,
		HmiStartUserAvatarUrl:      hmiStartUserAvatarUrl,
		HmiStopUserName:            hmiStopUserName,
		HmiStopUserAvatarUrl:       hmiStopUserAvatarUrl,
	}

	return machineDashboardResponse
}

func getDefaultStringValue(dashboardValue string) string {
	if dashboardValue == "" {
		return "-"
	}
	return dashboardValue
}

func getDefaultDateValue(dashboardValue string) string {
	if dashboardValue == "" {
		return "2022-06-19T10:00:00.000Z"
	}
	return dashboardValue
}

func getRejectList(hmiList *[]component.GeneralObject, eventId int) []map[string]interface{} {
	rejectList := make([]map[string]interface{}, 0)
	// var hmiInfo MachineHMIInfo

	for index := len(*hmiList) - 1; index >= 0; index-- {
		hmiInfo := MachineHMIInfo{}
		hmi := (*hmiList)[index]
		err := json.Unmarshal(hmi.ObjectInfo, &hmiInfo)

		if err != nil {
			continue
		}

		if hmiInfo.EventId != eventId {
			continue
		}
		//Remove non rejected quantity insertion
		if hmiInfo.RejectedQuantity == nil {
			continue
		}
		rejectObject := make(map[string]interface{})
		rejectObject["rejectedQty"] = *hmiInfo.RejectedQuantity
		rejectObject["created"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", hmiInfo.CreatedAt)

		rejectList = append(rejectList, rejectObject)
	}
	return rejectList
}

func getSetupTimeList(hmiList *[]component.GeneralObject) []map[string]interface{} {
	setUpTimeList := make([]map[string]interface{}, 0)
	// var hmiInfo MachineHMIInfo

	for _, hmiInfoObj := range *hmiList {
		hmiInfo := MachineHMIInfo{}
		err := json.Unmarshal(hmiInfoObj.ObjectInfo, &hmiInfo)

		if err != nil {
			continue
		}

		//Remove non rejected quantity insertion
		if hmiInfo.SetupTime == nil {
			continue
		}
		rejectObject := make(map[string]interface{})
		rejectObject["setupTime"] = *hmiInfo.SetupTime
		rejectObject["created"] = util.ConvertTimeToTimeZonCorrectedFormat("Asia/Singapore", hmiInfo.CreatedAt)

		setUpTimeList = append(setUpTimeList, rejectObject)
	}
	return setUpTimeList
}

func (v *MachineService) getStopList(eventId int, dbConnection *gorm.DB, targetTable string) []map[string]interface{} {
	stopList := make([]map[string]interface{}, 0)
	//select All hmi based on the event id
	hmiQuery := "select * from " + targetTable + " where JSON_EXTRACT(object_info, \"$.eventId\") = " + strconv.Itoa(eventId) + " order by JSON_EXTRACT(object_info, \"$.createdAt\") asc"
	var machineHmi []MachineHMI

	dbConnection.Raw(hmiQuery).Scan(&machineHmi)
	fmt.Println("hmiQuery: ", hmiQuery)
	if len(machineHmi) == 0 {
		return stopList
	}

	stopReasonIds := "("

	//Extract the stop reason id from the hmi
	for _, hmi := range machineHmi {
		hmiInfo := MachineHMIInfo{}
		err := json.Unmarshal(hmi.ObjectInfo, &hmiInfo)
		if err != nil {
			v.BaseService.Logger.Error("error in unmarshall machine hmi", zap.String("error", err.Error()))
			continue
		}
		if hmiInfo.ReasonId == 0 {
			continue
		}
		stopReasonIds = stopReasonIds + strconv.Itoa(hmiInfo.ReasonId) + ","
	}
	stopReasonIds = util.TrimSuffix(stopReasonIds, ",")
	stopReasonIds = stopReasonIds + ")"

	if stopReasonIds == "()" {
		return stopList
	}
	//map reason id with reason list
	stopReasonListQuery := "select id as id, JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.name\")) as name from machine_hmi_stop_reason where id in " + stopReasonIds

	fmt.Println("stop list query :", stopReasonListQuery)
	var stopReasonIdNameMap []StopReasons
	dbConnection.Raw(stopReasonListQuery).Scan(&stopReasonIdNameMap)
	fmt.Println("stopReasonIdNameMap : ", stopReasonIdNameMap)
	//create stoptime, start time and reason list

	statusFlag := "stopped"
	for index, machineHMI := range machineHmi {

		if machineHMI.getMachineHMIInfo().ReasonId != 0 {
			if machineHMI.getMachineHMIInfo().HMIStatus == "stopped" {
				stopReasonMap := make(map[string]interface{})
				stopReasonMap["startTime"] = util.ConvertTimeToTimeZonCorrectedFormat("Asia/Singapore", machineHMI.getMachineHMIInfo().CreatedAt)
				stopReasonMap["endTime"] = ""
				stopReasonMap["reason"] = getStopReasonById(stopReasonIdNameMap, machineHMI.getMachineHMIInfo().ReasonId)

				for _, startHmi := range machineHmi[index:] {
					if (startHmi.getMachineHMIInfo().HMIStatus != statusFlag) && (startHmi.getMachineHMIInfo().HMIStatus != "") {
						stopReasonMap["endTime"] = util.ConvertTimeToTimeZonCorrected("Asia/Singapore", startHmi.getMachineHMIInfo().CreatedAt)
						break
					}
				}

				if stopReasonMap["endTime"] == "" {
					stopReasonMap["endTime"] = util.GetZoneCurrentTimeInPMFormat("Asia/Singapore")
				}
				stopReasonMap["id"] = machineHMI.Id
				stopReasonMap["reasonId"] = machineHMI.getMachineHMIInfo().ReasonId
				stopReasonMap["remark"] = machineHMI.getMachineHMIInfo().Remark
				stopList = append(stopList, stopReasonMap)
			}
		}

	}
	//statusFlag = stop
	//Start to loop through hmiList
	//	check isReasonId
	//		check hmi.HMIStatus == statusFlag
	//			create map[start, readson]
	//			statusFlag = hmi.HMIStatus
	//			loop through hmiList[index]
	//				check secondHmi == statsFlag
	//					update map[stop]
	//					break
	//			check ma[stop] is empty
	//				assign currenttime
	return stopList
}

func getStopReasonById(stopReasonIdNameMap []StopReasons, id int) string {
	for _, reasons := range stopReasonIdNameMap {
		if reasons.Id == id {
			return reasons.Name
		}
	}
	return ""
}

func getUnassignedHMIResponse(machineId int, machineMasterInfo map[string]interface{},
	operatorInfo common.UserBasicInfo,
	hmiStopReasons component.RecordInfo,
	scheduledEventInfo *ScheduledOrderEventInfo, partInfo *PartInfo, eventStatusInfo *ProductionOrderStatusInfo, actualStartTime string, actualEndTime string, connectStatusColorCode string, machineConnectStatus string, locationList component.RecordInfo) datatypes.JSON {

	var stopList = make([]map[string]interface{}, 0)
	var rejectList = make([]map[string]interface{}, 0)

	hmiInfoResponse := HMIInfoResponse{
		MachineImage:            component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineImage"]), "text"),
		MachineId:               component.GetRecordIntInfo(machineId, "number"),
		MachineName:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["newMachineId"]), "text"),
		Status:                  component.GetRecordInfo(machineConnectStatus, "text"),
		HMIStatus:               component.GetRecordInfo("stopped", "text"),
		OperatorName:            component.GetRecordInfo(operatorInfo.FullName, "text"),
		OperatorAvatarUrl:       component.GetRecordInfo(operatorInfo.AvatarUrl, "text"),
		ProductionOrderId:       component.GetRecordInfo("", "text"),
		PartId:                  component.GetRecordInfo("", "text"),
		Description:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineDescription"]), "text"),
		StartTime:               component.GetRecordInfo("", "text"),
		EndTime:                 component.GetRecordInfo("", "text"),
		EventId:                 component.GetRecordIntInfo(0, "number"),
		StopReason:              hmiStopReasons,
		RejectQuantity:          component.GetRecordIntInfo(0, "number"),
		ActualStartTime:         component.GetRecordInfo(actualStartTime, "date"),
		ActualEndTime:           component.GetRecordInfo(actualEndTime, "date"),
		EventName:               component.GetRecordInfo("", "text"),
		PartDescription:         component.GetRecordInfo("", "text"),
		StopList:                stopList,
		StartUserUrl:            component.GetRecordInfo("", "text"),
		EndUserUrl:              component.GetRecordInfo("", "text"),
		ScheduledQty:            component.GetRecordIntInfo(0, "number"),
		CompletedQty:            component.GetRecordIntInfo(0, "number"),
		CompletedPercentage:     component.GetRecordIntInfo(0, "number"),
		RejectList:              rejectList,
		OrderStatus:             component.GetRecordInfo("", "text"),
		WarningMessage:          component.GetRecordInfo("", "text"),
		ColorCode:               component.GetRecordInfo("", "text"),
		StoppedBy:               component.GetRecordInfo("-", "text"),
		StartedBy:               component.GetRecordInfo("-", "text"),
		MachineConnectColorCode: component.GetRecordInfo(connectStatusColorCode, "text"),
		IsBatchManaged:          component.GetRecordObjectInfo(false, "bool"),
		Location:                locationList,
	}
	if partInfo != nil {
		hmiInfoResponse.PartId = component.GetRecordInfo(partInfo.PartNumber, "text")
		hmiInfoResponse.PartDescription = component.GetRecordInfo(partInfo.Description, "text")
	}
	if scheduledEventInfo != nil {
		hmiInfoResponse.ScheduledQty = component.GetRecordIntInfo(scheduledEventInfo.ScheduledQty, "number")
		hmiInfoResponse.CompletedQty = component.GetRecordIntInfo(scheduledEventInfo.CompletedQty, "number")
		hmiInfoResponse.CompletedPercentage = component.GetRecordIntInfo(int(scheduledEventInfo.PercentDone), "number")
		hmiInfoResponse.StartTime = component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", scheduledEventInfo.StartDate), "text")
		hmiInfoResponse.EndTime = component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", scheduledEventInfo.EndDate), "text")
		hmiInfoResponse.ProductionOrderId = component.GetRecordInfo(scheduledEventInfo.ProductionOrder, "text")
		hmiInfoResponse.RejectQuantity = component.GetRecordIntInfo(scheduledEventInfo.RejectedQty, "number")
	}
	if eventStatusInfo != nil {
		hmiInfoResponse.OrderStatus = component.GetRecordInfo(eventStatusInfo.Status, "text")
		hmiInfoResponse.ColorCode = component.GetRecordInfo(eventStatusInfo.ColorCode, "text")
	}
	serializedHMIResponse, _ := json.Marshal(hmiInfoResponse)
	return serializedHMIResponse

}

func getUnassignedToolingHMIResponse(machineId int, machineMasterInfo map[string]interface{},
	operatorInfo common.UserBasicInfo,
	hmiStopReasons component.RecordInfo,
	scheduledEventInfo map[string]interface{}, partInfo *PartInfo, eventStatusInfo *ProductionOrderStatusInfo, actualStartTime string, actualEndTime string, connectStatusColor string, statusName string) datatypes.JSON {

	var stopList = make([]map[string]interface{}, 0)
	var rejectList = make([]map[string]interface{}, 0)
	var setupTimeList = make([]map[string]interface{}, 0)

	hmiInfoResponse := HMIInfoResponse{
		MachineImage:            component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineImage"]), "text"),
		MachineId:               component.GetRecordIntInfo(machineId, "number"),
		MachineName:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["newMachineId"]), "text"),
		Status:                  component.GetRecordInfo(statusName, "text"),
		HMIStatus:               component.GetRecordInfo("stopped", "text"),
		OperatorName:            component.GetRecordInfo(operatorInfo.FullName, "text"),
		OperatorAvatarUrl:       component.GetRecordInfo(operatorInfo.AvatarUrl, "text"),
		ProductionOrderId:       component.GetRecordInfo("", "text"),
		PartId:                  component.GetRecordInfo("", "text"),
		Description:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineDescription"]), "text"),
		StartTime:               component.GetRecordInfo("", "text"),
		EndTime:                 component.GetRecordInfo("", "text"),
		EventId:                 component.GetRecordIntInfo(0, "number"),
		StopReason:              hmiStopReasons,
		RejectQuantity:          component.GetRecordIntInfo(0, "number"),
		ActualStartTime:         component.GetRecordInfo(actualStartTime, "date"),
		ActualEndTime:           component.GetRecordInfo(actualEndTime, "date"),
		EventName:               component.GetRecordInfo("", "text"),
		PartDescription:         component.GetRecordInfo("", "text"),
		StopList:                stopList,
		StartUserUrl:            component.GetRecordInfo("", "text"),
		EndUserUrl:              component.GetRecordInfo("", "text"),
		ScheduledQty:            component.GetRecordIntInfo(0, "number"),
		CompletedQty:            component.GetRecordInfo("0 Hours", "text"),
		CompletedPercentage:     component.GetRecordIntInfo(0, "number"),
		RejectList:              rejectList,
		OrderStatus:             component.GetRecordInfo("", "text"),
		WarningMessage:          component.GetRecordInfo("", "text"),
		ColorCode:               component.GetRecordInfo("", "text"),
		CanSetupTime:            component.GetDefaultBoolRecordInfo(),
		SetupTime:               component.GetRecordInfo(util.InterfaceToString(scheduledEventInfo["setupTime"]), "text"),
		BomId:                   component.GetRecordInfo("", "text"),
		TotalDuration:           component.GetRecordInfo("0 Hours", "text"),
		StartedBy:               component.GetRecordInfo("-", "text"),
		StoppedBy:               component.GetRecordInfo("-", "text"),
		MachineConnectColorCode: component.GetRecordInfo(connectStatusColor, "text"),
		IsBatchManaged:          component.GetRecordObjectInfo(false, "bool"),
		SetupTimeList:           setupTimeList,
	}
	if partInfo != nil {
		hmiInfoResponse.PartId = component.GetRecordInfo(partInfo.PartNumber, "text")
		hmiInfoResponse.PartDescription = component.GetRecordInfo(partInfo.Description, "text")
	}
	if scheduledEventInfo != nil {
		hmiInfoResponse.ScheduledQty = component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["scheduledQty"]), "number")
		hmiInfoResponse.CompletedQty = component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["completedQty"]), "number")
		hmiInfoResponse.CompletedPercentage = component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["percentDone"]), "number")
		hmiInfoResponse.StartTime = component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", util.InterfaceToString(scheduledEventInfo["startDate"])), "text")
		hmiInfoResponse.EndTime = component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", util.InterfaceToString(scheduledEventInfo["endDate"])), "text")
		hmiInfoResponse.ProductionOrderId = component.GetRecordInfo(util.InterfaceToString(scheduledEventInfo["productionOrder"]), "text")
		hmiInfoResponse.RejectQuantity = component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["rejectedQty"]), "number")
	}
	if eventStatusInfo != nil {
		hmiInfoResponse.OrderStatus = component.GetRecordInfo(eventStatusInfo.Status, "text")
		hmiInfoResponse.ColorCode = component.GetRecordInfo(eventStatusInfo.ColorCode, "text")
	}
	serializedHMIResponse, _ := json.Marshal(hmiInfoResponse)
	return serializedHMIResponse

}

func buildHMIResponse(eventId int, machineMasterInfo map[string]interface{}, scheduledEventInfo *ScheduledOrderEventInfo, partInfo *PartInfo, eventStatusInfo *ProductionOrderStatusInfo,
	hmiStatus string,
	operatorInfo common.UserBasicInfo,
	actualStartTime string,
	actualEndTime string,
	hmiStopReasons component.RecordInfo,
	stopList []map[string]interface{},
	rejectList []map[string]interface{},
	warningMessage string,
	mouldList component.RecordInfo,
	mouldDescription string,
	mouldId string,
	hmiId int,
	startedBy string,
	stoppedBy string,
	connectStatusColor string,
	machineConnectStatus string,
	startUrl string,
	endUrl string,
	operatorId int,
	locationList component.RecordInfo,
	isItAlreadyStarted bool,
	machineParameterId int,
) datatypes.JSON {

	hmiInfoResponse := HMIInfoResponse{
		MachineImage:            component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineImage"]), "text"),
		MachineId:               component.GetRecordIntInfo(scheduledEventInfo.MachineId, "number"),
		MachineName:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["newMachineId"]), "text"),
		Status:                  component.GetRecordInfo(machineConnectStatus, "text"),
		HMIStatus:               component.GetRecordInfo(hmiStatus, "text"),
		OperatorName:            component.GetRecordInfo(operatorInfo.FullName, "text"),
		OperatorAvatarUrl:       component.GetRecordInfo(operatorInfo.AvatarUrl, "text"),
		ProductionOrderId:       component.GetRecordInfo(scheduledEventInfo.ProductionOrder, "text"),
		PartId:                  component.GetRecordInfo(partInfo.PartNumber, "text"),
		Description:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineDescription"]), "text"),
		StartTime:               component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", scheduledEventInfo.StartDate), "text"),
		EndTime:                 component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", scheduledEventInfo.EndDate), "text"),
		EventId:                 component.GetRecordIntInfo(eventId, "number"),
		StopReason:              hmiStopReasons,
		RejectQuantity:          component.GetRecordIntInfo(scheduledEventInfo.RejectedQty, "number"),
		ActualStartTime:         component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", actualStartTime), "date"),
		ActualEndTime:           component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", actualEndTime), "date"),
		EventName:               component.GetRecordInfo(scheduledEventInfo.Name, "text"),
		PartDescription:         component.GetRecordInfo(partInfo.Description, "text"),
		StopList:                stopList,
		StartUserUrl:            component.GetRecordInfo(startUrl, "text"),
		EndUserUrl:              component.GetRecordInfo(endUrl, "text"),
		ScheduledQty:            component.GetRecordIntInfo(scheduledEventInfo.ScheduledQty, "number"),
		CompletedQty:            component.GetRecordIntInfo(scheduledEventInfo.CompletedQty, "number"),
		CompletedPercentage:     component.GetRecordIntInfo(int(scheduledEventInfo.PercentDone), "number"),
		RejectList:              rejectList,
		OrderStatus:             component.GetRecordInfo(eventStatusInfo.Status, "text"),
		WarningMessage:          component.GetRecordInfo(warningMessage, "text"),
		ColorCode:               component.GetRecordInfo(eventStatusInfo.ColorCode, "text"),
		MouldList:               mouldList,
		Cavity:                  component.GetRecordIntInfo(scheduledEventInfo.CustomCavity, "number"),
		EnableCustomCavity:      component.GetRecordObjectInfo(scheduledEventInfo.EnableCustomCavity, "bool"),
		MouldUp:                 component.GetRecordObjectInfo(scheduledEventInfo.MouldUp, "text"),
		MouldDown:               component.GetRecordObjectInfo(scheduledEventInfo.MouldDown, "text"),
		MouldDescription:        component.GetRecordObjectInfo(mouldDescription, "text"),
		MouldId:                 component.GetRecordObjectInfo(mouldId, "text"),
		HmiId:                   component.GetRecordIntInfo(hmiId, "number"),
		StartedBy:               component.GetRecordInfo(startedBy, "text"),
		StoppedBy:               component.GetRecordInfo(stoppedBy, "text"),
		MachineConnectColorCode: component.GetRecordInfo(connectStatusColor, "text"),
		IsBatchManaged:          component.GetRecordObjectInfo(partInfo.IsBatchManaged, "bool"),
		OperatorId:              component.GetRecordIntInfo(operatorId, "number"),
		MouldBatchResourceId:    component.GetRecordIntInfo(scheduledEventInfo.MouldBatchResourceId, "number"),
		CanComplete:             component.GetBoolRecordInfo(scheduledEventInfo.CanComplete, "bool"),
		IsShowBatchManaged:      component.GetBoolRecordInfo(!isItAlreadyStarted, "bool"),
		MachineParameterId:      component.GetRecordIntInfo(machineParameterId, "number"),
		Location:                locationList,
	}
	serializedObject, _ := json.Marshal(hmiInfoResponse)
	return serializedObject
}

func buildToolingHMIResponse(eventId int, machineMasterInfo map[string]interface{}, scheduledEventInfo map[string]interface{}, partInfo map[string]interface{}, eventStatusInfo *ProductionOrderStatusInfo,
	hmiStatus string,
	operatorInfo common.UserBasicInfo,
	actualStartTime string,
	actualEndTime string,
	hmiStopReasons component.RecordInfo,
	stopList []map[string]interface{},
	rejectList []map[string]interface{},
	warningMessage string,
	totalDuration string,
	bomId string,
	hmiId int,
	startedBy string,
	stoppedBy string,
	setupTimeList []map[string]interface{},
	connectStatusColor string,
	statusName string,
	startUrl string,
	stopUrl string,
	operatorId int,
	locationList component.RecordInfo,
) datatypes.JSON {

	completedQty := util.InterfaceToString(scheduledEventInfo["completedQty"])
	if completedQty == "" {
		completedQty = "0 Hours"
	}
	hmiInfoResponse := HMIInfoResponse{
		MachineImage:            component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineImage"]), "text"),
		MachineId:               component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["machineId"]), "number"),
		MachineName:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["newMachineId"]), "text"),
		Status:                  component.GetRecordInfo(statusName, "text"),
		HMIStatus:               component.GetRecordInfo(hmiStatus, "text"),
		OperatorName:            component.GetRecordInfo(operatorInfo.FullName, "text"),
		OperatorAvatarUrl:       component.GetRecordInfo(operatorInfo.AvatarUrl, "text"),
		ProductionOrderId:       component.GetRecordInfo(util.InterfaceToString(scheduledEventInfo["productionOrder"]), "text"),
		PartId:                  component.GetRecordInfo(util.InterfaceToString(partInfo["name"]), "text"),
		Description:             component.GetRecordInfo(util.InterfaceToString(machineMasterInfo["machineDescription"]), "text"),
		StartTime:               component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", util.InterfaceToString(scheduledEventInfo["startDate"])), "text"),
		EndTime:                 component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", util.InterfaceToString(scheduledEventInfo["endDate"])), "text"),
		EventId:                 component.GetRecordIntInfo(eventId, "number"),
		StopReason:              hmiStopReasons,
		RejectQuantity:          component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["rejectedQty"]), "number"),
		ActualStartTime:         component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", actualStartTime), "date"),
		ActualEndTime:           component.GetRecordInfo(util.ConvertTimeToZeroZone("Asia/Singapore", actualEndTime), "date"),
		EventName:               component.GetRecordInfo(util.InterfaceToString(scheduledEventInfo["name"]), "text"),
		PartDescription:         component.GetRecordInfo(util.InterfaceToString(partInfo["description"]), "text"),
		StopList:                stopList,
		StartUserUrl:            component.GetRecordInfo(startUrl, "text"),
		EndUserUrl:              component.GetRecordInfo(stopUrl, "text"),
		ScheduledQty:            component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["scheduledQty"]), "number"),
		CompletedQty:            component.GetRecordInfo(completedQty, "text"),
		CompletedPercentage:     component.GetRecordIntInfo(util.InterfaceToInt(scheduledEventInfo["percentDone"]), "number"),
		RejectList:              rejectList,
		OrderStatus:             component.GetRecordInfo(eventStatusInfo.Status, "text"),
		WarningMessage:          component.GetRecordInfo(warningMessage, "text"),
		ColorCode:               component.GetRecordInfo(eventStatusInfo.ColorCode, "text"),
		CanSetupTime:            component.GetBoolRecordInfo(util.InterfaceToBool(scheduledEventInfo["canSetupTime"]), "bool"),
		SetupTime:               component.GetRecordInfo(util.InterfaceToString(scheduledEventInfo["setupTime"]), "text"),
		TotalDuration:           component.GetRecordInfo(totalDuration, "text"),
		BomId:                   component.GetRecordInfo(bomId, "text"),
		HmiId:                   component.GetRecordIntInfo(hmiId, "number"),
		StartedBy:               component.GetRecordInfo(startedBy, "text"),
		StoppedBy:               component.GetRecordInfo(stoppedBy, "text"),
		MachineConnectColorCode: component.GetRecordInfo(connectStatusColor, "text"),
		SetupTimeList:           setupTimeList,
		IsBatchManaged:          component.GetRecordObjectInfo(util.InterfaceToBool(partInfo["isBatchManaged"]), "bool"),
		OperatorId:              component.GetRecordIntInfo(operatorId, "number"),
		Location:                locationList,
	}
	serializedObject, _ := json.Marshal(hmiInfoResponse)
	return serializedObject
}

func setOperatorInfo(operatorId int) OperatorResponse {
	operatorInfo := getOperatorAvatarUrl(operatorId)

	if cmp.Equal(common.UserBasicInfo{}, operatorInfo) {
		return OperatorResponse{}
	}

	// if (common.UserBasicInfo{}) == operatorInfo {
	// 	return OperatorResponse{}
	// }

	return OperatorResponse{
		Id:           operatorId,
		OperatorName: operatorInfo.Username,
		AvatarUrl:    operatorInfo.AvatarUrl,
	}

}

func getOperatorAvatarUrl(operatorId int) common.UserBasicInfo {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	operatorInfo := authService.GetUserInfoById(operatorId)

	return operatorInfo
}

type OperatorResponse struct {
	Id           int    `json:"id"`
	OperatorName string `json:"operatorName"`
	AvatarUrl    string `json:"avatarUrl"`
}

type HmiStopReasonInfo struct {
	Data   []HMIStopReasonResponse `json:"data"`
	IsEdit bool                    `json:"isEdit"`
	Type   string                  `json:"type"`
	Value  string                  `json:"value"`
}

type HMIStopReasonResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IconCls string `json:"iconCls"`
}

type HMIInfoResponse struct {
	MachineImage            component.RecordInfo     `json:"machineImage"`
	MachineId               component.RecordInfo     `json:"machineId"`
	MachineName             component.RecordInfo     `json:"machineName"`
	Status                  component.RecordInfo     `json:"status"`
	HMIStatus               component.RecordInfo     `json:"hmiStatus"`
	OperatorId              component.RecordInfo     `json:"operatorId"`
	OperatorName            component.RecordInfo     `json:"operatorName"`
	OperatorAvatarUrl       component.RecordInfo     `json:"operatorAvatarUrl"`
	ProductionOrderId       component.RecordInfo     `json:"productionOrderId"`
	PartId                  component.RecordInfo     `json:"partId"`
	Description             component.RecordInfo     `json:"description"`
	StartTime               component.RecordInfo     `json:"startTime"`
	EndTime                 component.RecordInfo     `json:"endTime"`
	EventId                 component.RecordInfo     `json:"eventId"`
	ActualStartTime         component.RecordInfo     `json:"actualStartTime"`
	ActualEndTime           component.RecordInfo     `json:"actualEndTime"`
	StopReason              component.RecordInfo     `json:"stopReason"`
	RejectQuantity          component.RecordInfo     `json:"rejectQuantity"`
	EventName               component.RecordInfo     `json:"eventName"`
	PartDescription         component.RecordInfo     `json:"partDescription"`
	ScheduledQty            component.RecordInfo     `json:"scheduledQty"`
	CompletedQty            component.RecordInfo     `json:"completedQty"`
	CompletedPercentage     component.RecordInfo     `json:"completedPercentage"`
	OrderStatus             component.RecordInfo     `json:"orderStatus"`
	StopList                []map[string]interface{} `json:"stopList"`
	StartUserUrl            component.RecordInfo     `json:"startUserUrl"`
	EndUserUrl              component.RecordInfo     `json:"endUserUrl"`
	SetupTimeList           []map[string]interface{} `json:"setupTimeList"`
	RejectList              []map[string]interface{} `json:"rejectList"`
	WarningMessage          component.RecordInfo     `json:"warningMessage"`
	ColorCode               component.RecordInfo     `json:"colorCode"`
	MouldList               component.RecordInfo     `json:"mouldList"`
	Cavity                  component.RecordInfo     `json:"cavity"`
	EnableCustomCavity      component.RecordInfo     `json:"enableCustomCavity"`
	MouldDescription        component.RecordInfo     `json:"mouldDescription"`
	MouldUp                 component.RecordInfo     `json:"mouldup"`
	MouldDown               component.RecordInfo     `json:"mouldDown"`
	MouldId                 component.RecordInfo     `json:"mouldId"`
	CanSetupTime            component.RecordInfo     `json:"canSetupTime"`
	SetupTime               component.RecordInfo     `json:"setupTime"`
	TotalDuration           component.RecordInfo     `json:"totalDuration"`
	BomId                   component.RecordInfo     `json:"bomId"`
	HmiId                   component.RecordInfo     `json:"hmiId"`
	StartedBy               component.RecordInfo     `json:"startedBy"`
	StoppedBy               component.RecordInfo     `json:"stoppedBy"`
	MachineConnectColorCode component.RecordInfo     `json:"machineConnectColorCode"`
	IsBatchManaged          component.RecordInfo     `json:"isBatchManaged"`
	Location                component.RecordInfo     `json:"location"`
	MouldBatchResourceId    component.RecordInfo     `json:"mouldBatchResourceId"`
	IsShowBatchManaged      component.RecordInfo     `json:"isShowBatchManaged"`
	CanComplete             component.RecordInfo     `json:"canComplete"`
	MachineParameterId      component.RecordInfo     `json:"machineParameterId"`
}

type StopList struct {
	StopTime  string `json:"stopTime"`
	StartTime string `json:"startTime"`
	Reason    string `json:"reason"`
}

type MachineStatResponse struct {
	Daily  int   `json:"Daily"`
	Actual int   `json:"Actual"`
	Oee    int   `json:"Oee"`
	TS     int64 `json:"Ts"`
}

type MachineDashboardResponse struct {
	Actual                     int     `json:"actual"`
	Brand                      string  `json:"brand"`
	MachineImage               string  `json:"machineImage"`
	Model                      string  `json:"model"`
	NewMachineId               string  `json:"newMachineId"`
	CurrentStatus              string  `json:"currentStatus"`
	ProductionOrder            string  `json:"productionOrder"`
	PartName                   string  `json:"partName"`
	PartDescription            string  `json:"partDescription"`
	ScheduleStartTime          string  `json:"scheduleStartTime"`
	ScheduleEndTime            string  `json:"scheduleEndTime"`
	ActualStartTime            string  `json:"actualStartTime"`
	ActualEndTime              string  `json:"actualEndTime"`
	EstimatedEndTime           string  `json:"estimatedEndTime"`
	Oee                        float64 `json:"oee"`
	Availability               int     `json:"availability"`
	Performance                int     `json:"performance"`
	Quality                    int     `json:"quality"`
	PlannedQuality             int     `json:"plannedQuality"`
	DailyPlannedQuality        int     `json:"dailyPlannedQuality"`
	Completed                  int     `json:"completed"`
	Rejects                    int     `json:"rejects"`
	CompletedPercentage        float64 `json:"completedPercentage"`
	OeeDiff                    float64 `json:"oeeDiff"`
	AvailabilityDiff           float64 `json:"availabilityDiff"`
	PerformanceDiff            float64 `json:"performanceDiff"`
	QualityDiff                float64 `json:"qualityDiff"`
	CompletedPercentageDiff    float64 `json:"completedPercentageDiff"`
	OverallRejectedQty         int     `json:"overallRejectedQty"`
	OverallCompletedPercentage float64 `json:"overallCompletedPercentage"`
	PlannedQualityHrs          string  `json:"plannedQualityHrs"`
	DailyPlannedQualityHrs     string  `json:"dailyPlannedQualityHrs"`
	ActualHrs                  string  `json:"actualHrs"`
	CompletedHrs               string  `json:"completedHrs"`
	RejectsHrs                 string  `json:"rejectsHrs"`
	OverallRejectedHrs         int     `json:"overallRejectedHrs"`
	BomId                      string  `json:"bomId"`
	Remark                     string  `json:"remark"`
	IsUnderMaintenance         bool    `json:"isUnderMaintenance"`
	ColorCode                  string  `json:"colorCode"`
	HmiStartUserName           string  `json:"hmiStartUserName"`
	HmiStartUserAvatarUrl      string  `json:"hmiStartUserAvatarUrl"`
	HmiStopUserName            string  `json:"hmiStopUserName"`
	HmiStopUserAvatarUrl       string  `json:"hmiStopUserAvatarUrl"`
}

type MachineParamResults struct {
	MouldId    string
	OldToolId  string
	ProgramNo  string
	ApprovedBy string
}

type MouldingMachineOverview struct {
	WorkOrderType struct {
		Series []struct {
			Data [][]interface{} `json:"data"`
			Name string          `json:"name"`
			Type string          `json:"type"`
		} `json:"series"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Subtitle struct {
			Text string `json:"text"`
		} `json:"subtitle"`
		Chart struct {
			Width int `json:"width"`
		} `json:"chart"`
	} `json:"workOrderType"`
	WorkOrders struct {
		Series []struct {
			Data  [][]interface{} `json:"data"`
			Name  string          `json:"name"`
			Type  string          `json:"type"`
			YAxis int             `json:"yAxis"`
		} `json:"series"`
		YAxis struct {
			Min   int `json:"min"`
			Title struct {
				Text string `json:"text"`
			} `json:"title"`
		} `json:"yAxis"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Subtitle struct {
			Text string `json:"text"`
		} `json:"subtitle"`
		XAxis struct {
			Categories []string `json:"categories"`
		} `json:"xAxis"`
		Chart struct {
			Width int `json:"width"`
		} `json:"chart"`
	} `json:"workOrders"`
	UnscheduledDowntime struct {
		Series []struct {
			Data  [][]interface{} `json:"data"`
			Name  string          `json:"name"`
			Type  string          `json:"type"`
			YAxis int             `json:"yAxis"`
		} `json:"series"`
		YAxis struct {
			Min   int `json:"min"`
			Title struct {
				Text string `json:"text"`
			} `json:"title"`
		} `json:"yAxis"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Subtitle struct {
			Text string `json:"text"`
		} `json:"subtitle"`
		XAxis struct {
			Categories []string `json:"categories"`
		} `json:"xAxis"`
		Chart struct {
			Width int `json:"width"`
		} `json:"chart"`
	} `json:"unscheduledDowntime"`
	OeePercentage struct {
		Series []struct {
			Data  [][]interface{} `json:"data"`
			Name  string          `json:"name"`
			Type  string          `json:"type"`
			YAxis int             `json:"yAxis"`
		} `json:"series"`
		YAxis struct {
			Min   int `json:"min"`
			Title struct {
				Text string `json:"text"`
			} `json:"title"`
		} `json:"yAxis"`
		Title struct {
			Text string `json:"text"`
		} `json:"title"`
		Subtitle struct {
			Text string `json:"text"`
		} `json:"subtitle"`
		XAxis struct {
			Categories []string `json:"categories"`
		} `json:"xAxis"`
		Chart struct {
			Width int `json:"width"`
		} `json:"chart"`
	} `json:"oeePercentage"`
	OpenWorkOrder struct {
		Data int `json:"data"`
	} `json:"openWorkOrder"`
	ClosedWorkOrder struct {
		Data int `json:"data"`
	} `json:"closedWorkOrder"`
	TotalDownTime struct {
		Data string `json:"data"`
	} `json:"totalDownTime"`
}

func (v *MachineService) getMouldingMachineOverview(ctx *gin.Context) {
	// projectId := ctx.Param("projectId")
	// recordId := ctx.Param("recordId")

	// dbConnection := ms.BaseService.ServiceDatabases[projectId]

	// openWorkOrderQuery := "select count(*) from maintenance_work_order where object_info->>'$.assetId'=" + recordId + " AND object_info->>'$.assetId' != " + strconv.Itoa(WorkOrderDoneStatus)
	// var countOpenWorkOrders int
	// dbConnection.Raw(openWorkOrderQuery).Scan(&countOpenWorkOrders)

	// cloasedWorkOrderQuery := "select count(*) from maintenance_work_order where object_info->>'$.assetId'=" + recordId + " AND object_info->>'$.assetId' = " + strconv.Itoa(WorkOrderDoneStatus)
	// var countClosedWorkOrders int
	// dbConnection.Raw(cloasedWorkOrderQuery).Scan(&countClosedWorkOrders)

	// downtimeHoursQuery := "select sum(timestampdiff(HOUR, t1.object_info ->> '$.workOrderScheduledEndDate', t2.object_info ->> '$.workOrderScheduledStartDate')) as di from fuyu_mes.maintenance_work_order t1 left join fuyu_mes.maintenance_work_order t2 on t1.id + 1 = t2.id where t1.object_info -> '$.assetId' =" + recordId + " AND t2.object_info -> '$.assetId' = " + recordId + " order by t1.id"
	// var downtimeWorkOrders int
	// dbConnection.Raw(downtimeHoursQuery).Scan(&downtimeWorkOrders)

	// oeePercentageQuery := "select round(avg(stats_info ->> '$.oee'), 2) / 100 as oeePercentage, MONTH(from_unixtime(ts)) as month from fuyu_mes.machine_statistics where machine_id=" + recordId +  " group by MONTH(from_unixtime(ts))"
	// var oeePercentage int
	// dbConnection.Raw(downtimeHoursQuery).Scan(&downtimeWorkOrders)

	response := `
{
    "workOrderType":
    {
        "series":
        [
            {
                "data":
                [
                    [
                        "Corrective Work Order",
                        1672
                    ],
                    [
                        "Preventive Work Order",
                        211
                    ],
                    [
                        "Testing Work Order",
                        1590
                    ]
                ],
                "name": "Work Order Type",
                "type": "pie"
            }
        ],
        "title":
        {
            "text": "Work Order Type Percentage"
        },
        "subtitle":
        {
            "text": "Average Work Order Type/Machine"
        },
        "chart":
        {
            "width": 800
        }
    },
    "workOrders":
    {
        "series":
        [
            {
                "data":
                [
                    [
                        "Jan",
                        45
                    ],
                    [
                        "Feb",
                        23
                    ],
                    [
                        "Mar",
                        12
                    ],
                    [
                        "April",
                        87
                    ],
                    [
                        "May",
                        34
                    ]
                ],
                "name": "New/Closed Work Order",
                "type": "column",
                "yAxis": 0
            }
        ],
        "yAxis":
        {
            "min": 0,
            "title":
            {
                "text": "Work Orders"
            }
        },
        "title":
        {
            "text": "Total New/Closed Work Orders"
        },
        "subtitle":
        {
            "text": "Average work orders"
        },
        "xAxis":
        {
            "categories":
            [
                "Jan",
                "Feb",
                "Mar",
                "Apr",
                "May"
            ]
        },
        "chart":
        {
            "width": 800
        }
    },
    "unscheduledDowntime":
    {
        "series":
        [
            {
                "data":
                [
                    [
                        "Jan",
                        5
                    ],
                    [
                        "Feb",
                        8
                    ],
                    [
                        "Mar",
                        12
                    ],
                    [
                        "April",
                        16
                    ],
                    [
                        "May",
                        1
                    ]
                ],
                "name": "Unscheduled Downtime",
                "type": "column",
                "yAxis": 0
            }
        ],
        "yAxis":
        {
            "min": 0,
            "title":
            {
                "text": "Downtime (Hours)"
            }
        },
        "title":
        {
            "text": "Machines Total Unscheduled Downtime"
        },
        "subtitle":
        {
            "text": "Average Downtime"
        },
        "xAxis":
        {
            "categories":
            [
                "Jan",
                "Feb",
                "Mar",
                "Apr",
                "May"
            ]
        },
        "chart":
        {
            "width": 800
        }
    },
    "oeePercentage":
    {
        "series":
        [
            {
                "data":
                [
                    [
                        "Jan",
                        97
                    ],
                    [
                        "Feb",
                        92
                    ],
                    [
                        "Mar",
                        65
                    ],
                    [
                        "April",
                        87
                    ],
                    [
                        "May",
                        98
                    ]
                ],
                "name": "OEE Percentage",
                "type": "column",
                "yAxis": 0
            }
        ],
        "yAxis":
        {
            "min": 0,
            "title":
            {
                "text": "OEE"
            }
        },
        "title":
        {
            "text": "Machines OEE Performance"
        },
        "subtitle":
        {
            "text": "Average OEE Performance"
        },
        "xAxis":
        {
            "categories":
            [
                "Jan",
                "Feb",
                "Mar",
                "Apr",
                "May"
            ]
        },
        "chart":
        {
            "width": 800
        }
    },
    "openWorkOrder":
    {
        "data": 12
    },
    "closedWorkOrder":
    {
        "data": 20
    },
    "totalDownTime":
    {
        "data": "20.12 Hrs"
    }
}
`
	var rMessage MouldingMachineOverview
	json.Unmarshal([]byte(response), &rMessage)

	ctx.JSON(http.StatusOK, rMessage)
}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
