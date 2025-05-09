package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/database"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/machines/handler/const_util"
	"cx-micro-flake/services/machines/handler/model"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func (v *MachineService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, MachineComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *MachineService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, MachineComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, MachineComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *MachineService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, MachineComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *MachineService) CreateMachineParam(projectId string, testId int, userId int) (error, int) {
	machineParamInfo := MachineParamInfo{}
	machineParamInfo.TestId = testId
	rawParamInfo, _ := json.Marshal(machineParamInfo)
	generalObject := component.GeneralObject{ObjectInfo: rawParamInfo}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, recordId := Create(dbConnection, MachineParameterTable, generalObject)
	if err != nil {
		v.BaseService.Logger.Error("creating machine param had failed", zap.String("error", err.Error()))
	}
	var machineParamReferenceId string
	if recordId < 10 {
		machineParamReferenceId = "MP0000" + strconv.Itoa(recordId)
	} else if recordId < 100 {
		machineParamReferenceId = "MP000" + strconv.Itoa(recordId)
	} else if recordId < 1000 {
		machineParamReferenceId = "MP00" + strconv.Itoa(recordId)
	} else if recordId < 10000 {
		machineParamReferenceId = "MP0" + strconv.Itoa(recordId)
	} else if recordId < 100000 {
		machineParamReferenceId = "MP" + strconv.Itoa(recordId)
	}
	updatingData := make(map[string]interface{})
	_, machineParamGeneralObject := Get(dbConnection, MachineParameterTable, recordId)
	machinePram := MachineParameter{ObjectInfo: machineParamGeneralObject.ObjectInfo}
	machinePramObject := machinePram.getMachineParamInfo()
	machinePramObject.MachineParamReferenceId = machineParamReferenceId
	machinePramObject.MoldTemperatureCoolingMedium = 1
	machinePramObject.ObjectStatus = common.ObjectStatusActive
	machinePramObject.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machinePramObject.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	machinePramObject.CreatedBy = userId
	machinePramObject.LastUpdatedBy = userId
	machinePramObject.CanEdit = true

	serializedParamInfo, _ := json.Marshal(machinePramObject)
	updatingData["object_info"] = serializedParamInfo
	err = Update(dbConnection, MachineParameterTable, recordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("updating machine param had failed", zap.String("error", err.Error()))
	}
	v.BaseService.Logger.Info("updating machine param", zap.Any("machine_param_id", recordId))
	return err, recordId
}

func (v *MachineService) IsMachineAlreadyUnderMaintenance(projectId string, machineId int) (error, bool) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, MachineMasterTable, machineId)
	if err == nil {
		machineMaster := MachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getMachineMasterInfo()
		if machineMasterInfo.MachineStatus == MachineStatusMaintenance {
			return nil, true
		} else {
			return nil, false
		}
	}
	return err, false
}

func (v *MachineService) MoveMachineToMaintenance(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, MachineMasterTable, machineId)
	if err == nil {
		machineMaster := MachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusMaintenance
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, MachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveMachineLiveStatusToMaintenance(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, MachineMasterTable, machineId)
	if err == nil {
		machineMaster := MachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getMachineMasterInfo()
		machineMasterInfo.MachineConnectStatus = machineConnectStatusMaintenance
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, MachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveMachineLiveStatusToActive(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, MachineMasterTable, machineId)
	if err == nil {
		machineMaster := MachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getMachineMasterInfo()
		machineMasterInfo.MachineConnectStatus = machineConnectStatusWaitingForFeed
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, MachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveAssemblyMachineToMaintenance(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	if err == nil {
		machineMaster := AssemblyMachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getAssemblyMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusMaintenance
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, AssemblyMachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveToolingMachineToMaintenance(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, ToolingMachineMasterTable, machineId)
	if err == nil {
		machineMaster := ToolingMachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getToolingMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusMaintenance
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, ToolingMachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveMachineToActive(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, MachineMasterTable, machineId)
	if err == nil {
		machineMaster := MachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusActive
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, MachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveAssemblyMachineToActive(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, AssemblyMachineMasterTable, machineId)
	if err == nil {
		machineMaster := AssemblyMachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getAssemblyMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusActive
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, AssemblyMachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) MoveToolingMachineToActive(projectId string, machineId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, ToolingMachineMasterTable, machineId)
	if err == nil {
		machineMaster := ToolingMachineMaster{ObjectInfo: generalObject.ObjectInfo}
		machineMasterInfo := machineMaster.getToolingMachineMasterInfo()
		machineMasterInfo.MachineStatus = MachineStatusActive
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(machineMasterInfo)
		updatingData["object_info"] = serializedObject
		err = Update(v.BaseService.ReferenceDatabase, ToolingMachineMasterTable, machineId, updatingData)
	}
	return err
}

func (v *MachineService) GetListOfMachines(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := GetObjects(dbConnection, MachineMasterTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}

func (v *MachineService) GetListOfToolingMachines(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := GetObjects(dbConnection, ToolingMachineMasterTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}

func (v *MachineService) GetListOfAssemblyMachines(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := GetObjects(dbConnection, AssemblyMachineMasterTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}

func (v *MachineService) GetListOfAssemblyLines(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := GetObjects(dbConnection, AssemblyMachineLineTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}
func (v *MachineService) GetAssemblyLineFromId(projectId string, resourceId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, assemblyLineObject := Get(dbConnection, AssemblyMachineLineTable, resourceId)

	return err, assemblyLineObject
}

func (v *MachineService) GetListOfMachineSubCategory(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := GetObjects(dbConnection, MachineSubCategoryTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}

func (v *MachineService) GetDepartmentDisplayEnabledMachines(projectId string, listOfDepartments []int, allowedDepartment []component.OrderedData) datatypes.JSON {
	// select the department machines
	var departmentList string
	departmentList = " IN ("
	for index, departmentId := range listOfDepartments {
		departmentList += strconv.Itoa(departmentId) + ","
		if index == len(listOfDepartments)-1 {
			departmentList = strings.TrimSuffix(departmentList, ",")
			departmentList += ")"
		}
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	condition := " object_info ->>'$.department'  " + departmentList

	selectedDepartmentMachines, _ := GetConditionalObjects(dbConnection, MachineMasterTable, condition)

	condition = " object_info ->>'$.displayEnabled'= 'true'"
	allDisplayEnabledMachines, _ := GetConditionalObjects(dbConnection, MachineDisplaySettingTable, condition)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var machineIds = make(map[int][]int, 0)
	for _, machineInterface := range *selectedDepartmentMachines {
		for _, displayEnabledMachines := range *allDisplayEnabledMachines {
			if machineInterface.Id == displayEnabledMachines.Id {
				machineMaster := MachineMaster{ObjectInfo: machineInterface.ObjectInfo}
				err, _ := productionOrderInterface.GetCurrentScheduledEvent(projectId, machineInterface.Id)
				if err != nil {
					if machineMaster.getMachineMasterInfo().MachineConnectStatus == machineConnectStatusWaitingForFeed {
						continue
					}
				}

				if _, ok := machineIds[machineMaster.getMachineMasterInfo().Department]; ok {
					machineIds[machineMaster.getMachineMasterInfo().Department] = append(machineIds[machineMaster.getMachineMasterInfo().Department], machineInterface.Id)
				} else {
					var machineList []int
					machineList = append(machineList, machineInterface.Id)
					machineIds[machineMaster.getMachineMasterInfo().Department] = machineList
				}
			}
		}
	}
	var arrayOfDepartmentDisplayMachines = make([]DepartmentDisplayMachines, 0)
	for key, value := range machineIds {
		departmentDisplayMachines := DepartmentDisplayMachines{}
		departmentDisplayMachines.DepartmentId = key
		departmentDisplayMachines.Name = getDepartmentName(allowedDepartment, key)
		departmentDisplayMachines.ListOfMachines = value
		arrayOfDepartmentDisplayMachines = append(arrayOfDepartmentDisplayMachines, departmentDisplayMachines)
	}

	machineModuleSetting, err := GetObjects(dbConnection, MachineModuleSettingTable)
	var displayInterval = 25
	if err == nil {
		machineModuleInterface := (*machineModuleSetting)[0]
		machineModule := MachineModuleSetting{ObjectInfo: machineModuleInterface.ObjectInfo}
		displayInterval = machineModule.getMachineModuleSettingInfo().DisplayRotateInterval
	}

	var displaySettingResponse = make(map[string]interface{})
	displaySettingResponse["displayEnabledMachines"] = arrayOfDepartmentDisplayMachines
	displaySettingResponse["displayInterval"] = displayInterval
	displaySettingResponse["displayEnabledAssemblyMachines"] = v.GetDepartmentDisplayEnabledAssemblyMachines(projectId, listOfDepartments, allowedDepartment)
	displaySettingResponse["displayEnabledToolingMachines"] = v.GetDepartmentDisplayEnabledToolingMachines(projectId, listOfDepartments, allowedDepartment)

	rawData, _ := json.Marshal(displaySettingResponse)

	return rawData

}

func getDepartmentName(allowedDepartments []component.OrderedData, departmentId int) string {
	var departmentName string

	for _, department := range allowedDepartments {
		if department.Id == departmentId {
			return department.Value
		}
	}

	return departmentName
}

func (v *MachineService) GetDepartmentDisplayEnabledAssemblyMachines(projectId string, listOfDepartments []int, allowedDepartment []component.OrderedData) []DepartmentDisplayMachines {
	// select the department machines
	var departmentList string
	departmentList = " IN ("
	for index, departmentId := range listOfDepartments {
		departmentList += strconv.Itoa(departmentId) + ","
		if index == len(listOfDepartments)-1 {
			departmentList = strings.TrimSuffix(departmentList, ",")
			departmentList += ")"
		}
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	condition := " object_info ->>'$.department'  " + departmentList

	selectedDepartmentMachines, _ := GetConditionalObjects(dbConnection, AssemblyMachineMasterTable, condition)

	condition = " object_info ->>'$.displayEnabled'= 'true'"
	allDisplayEnabledMachines, _ := GetConditionalObjects(dbConnection, AssemblyMachineDisplaySettingTable, condition)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var machineIds = make(map[int][]int, 0)
	for _, machineInterface := range *selectedDepartmentMachines {
		for _, displayEnabledMachines := range *allDisplayEnabledMachines {
			if machineInterface.Id == displayEnabledMachines.Id {
				err, _ := productionOrderInterface.GetCurrentAssemblyScheduledEvent(projectId, machineInterface.Id)
				machineMaster := AssemblyMachineMaster{ObjectInfo: machineInterface.ObjectInfo}
				if err != nil {
					if machineMaster.getAssemblyMachineMasterInfo().MachineConnectStatus == machineConnectStatusWaitingForFeed {
						continue
					}

				}

				if _, ok := machineIds[machineMaster.getAssemblyMachineMasterInfo().Department]; ok {
					machineIds[machineMaster.getAssemblyMachineMasterInfo().Department] = append(machineIds[machineMaster.getAssemblyMachineMasterInfo().Department], machineInterface.Id)
				} else {
					var machineList []int
					machineList = append(machineList, machineInterface.Id)
					machineIds[machineMaster.getAssemblyMachineMasterInfo().Department] = machineList
				}
			}
		}
	}
	var arrayOfDepartmentDisplayMachines = make([]DepartmentDisplayMachines, 0)
	for key, value := range machineIds {
		departmentDisplayMachines := DepartmentDisplayMachines{}
		departmentDisplayMachines.DepartmentId = key
		departmentDisplayMachines.Name = getDepartmentName(allowedDepartment, key)
		departmentDisplayMachines.ListOfMachines = value
		arrayOfDepartmentDisplayMachines = append(arrayOfDepartmentDisplayMachines, departmentDisplayMachines)
	}

	return arrayOfDepartmentDisplayMachines

}

func (v *MachineService) GetDepartmentDisplayEnabledToolingMachines(projectId string, listOfDepartments []int, allowedDepartment []component.OrderedData) []DepartmentDisplayMachines {
	// select the department machines
	var departmentList string
	departmentList = " IN ("
	for index, departmentId := range listOfDepartments {
		departmentList += strconv.Itoa(departmentId) + ","
		if index == len(listOfDepartments)-1 {
			departmentList = strings.TrimSuffix(departmentList, ",")
			departmentList += ")"
		}
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	condition := " object_info ->>'$.department'  " + departmentList

	selectedDepartmentMachines, _ := GetConditionalObjects(dbConnection, ToolingMachineMasterTable, condition)

	condition = " object_info ->>'$.displayEnabled'= 'true'"
	allDisplayEnabledMachines, _ := GetConditionalObjects(dbConnection, ToolingMachineDisplaySettingTable, condition)

	productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)

	var machineIds = make(map[int][]int, 0)
	for _, machineInterface := range *selectedDepartmentMachines {
		for _, displayEnabledMachines := range *allDisplayEnabledMachines {
			if machineInterface.Id == displayEnabledMachines.Id {
				err, _ := productionOrderInterface.GetCurrentToolingScheduledEvent(projectId, machineInterface.Id)
				machineMaster := ToolingMachineMaster{ObjectInfo: machineInterface.ObjectInfo}
				if err != nil {
					if machineMaster.getToolingMachineMasterInfo().MachineConnectStatus == machineConnectStatusWaitingForFeed {
						continue
					}

				}

				if _, ok := machineIds[machineMaster.getToolingMachineMasterInfo().Department]; ok {
					machineIds[machineMaster.getToolingMachineMasterInfo().Department] = append(machineIds[machineMaster.getToolingMachineMasterInfo().Department], machineInterface.Id)
				} else {
					var machineList []int
					machineList = append(machineList, machineInterface.Id)
					machineIds[machineMaster.getToolingMachineMasterInfo().Department] = machineList
				}
			}
		}
	}
	var arrayOfDepartmentDisplayMachines = make([]DepartmentDisplayMachines, 0)
	for key, value := range machineIds {
		departmentDisplayMachines := DepartmentDisplayMachines{}
		departmentDisplayMachines.DepartmentId = key
		departmentDisplayMachines.Name = getDepartmentName(allowedDepartment, key)
		departmentDisplayMachines.ListOfMachines = value
		arrayOfDepartmentDisplayMachines = append(arrayOfDepartmentDisplayMachines, departmentDisplayMachines)
	}

	return arrayOfDepartmentDisplayMachines

}

func (v *MachineService) GetMachineInfoById(projectId string, machineId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, machineData := Get(dbConnection, MachineMasterTable, machineId)

	if err != nil {
		return getError("Machine record not available"), machineData
	}

	return err, machineData
}

func (v *MachineService) GetAssemblyMachineInfoById(machineId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, machineData := Get(dbConnection, AssemblyMachineMasterTable, machineId)

	if err != nil {
		return getError("Machine record not available"), machineData
	}

	return err, machineData
}
func (v *MachineService) GetAssemblyMachineDefaultManpower(machineId int) int {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	assemblymanMasterInfo := GetAssemblyMachineFromId(dbConnection, machineId)

	if assemblymanMasterInfo == nil {
		return 0
	}
	err, assemblyLineObject := v.GetAssemblyLineFromId(ProjectID, assemblymanMasterInfo.AssemblyLineOption)
	if err != nil {
		return 0
	}
	return getAssemblyLinesInfo(assemblyLineObject.ObjectInfo).DefaultPlannedManpower
}
func (v *MachineService) AddMouldingForceStop(projectId string, machineId int, eventId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	objectInfo := make(map[string]interface{})
	objectInfo["eventId"] = eventId
	objectInfo["machineId"] = machineId
	objectInfo["remark"] = "Operator force to stop it"
	objectInfo["hmiStatus"] = "stop"

	serialisedObject, _ := json.Marshal(objectInfo)
	commonObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, _ := Create(dbConnection, MachineHMITable, commonObject)

	if err != nil {
		return getError("Can't insert HMI record")
	}

	return err
}

func (v *MachineService) AddToolingForceStop(projectId string, machineId int, eventId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	objectInfo := make(map[string]interface{})
	objectInfo["eventId"] = eventId
	objectInfo["machineId"] = machineId
	objectInfo["remark"] = "Operator force to stop it"
	objectInfo["hmiStatus"] = "stop"

	serialisedObject, _ := json.Marshal(objectInfo)
	commonObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, _ := Create(dbConnection, ToolingMachineHmiTable, commonObject)

	if err != nil {
		return getError("Can't insert HMI record")
	}

	return err
}

func (v *MachineService) AddAssemblyForceStop(projectId string, machineId int, eventId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	objectInfo := make(map[string]interface{})
	objectInfo["eventId"] = eventId
	objectInfo["machineId"] = machineId
	objectInfo["remark"] = "Operator force to stop it"
	objectInfo["hmiStatus"] = "stop"

	serialisedObject, _ := json.Marshal(objectInfo)
	commonObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, _ := Create(dbConnection, AssemblyMachineHmiTable, commonObject)

	if err != nil {
		return getError("Can't insert HMI record")
	}

	return err
}

func (v *MachineService) GetListOfAllowedAssemblyLines(projectId string, userId int) []int {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfAssemblyMachines, err := GetObjects(dbConnection, AssemblyMachineLineTable)
	if err != nil {
		return make([]int, 0)
	}
	var lineUserCache = make([]int, 0)
	for _, assemblyMachineInterface := range *listOfAssemblyMachines {
		assemblyInfo := getAssemblyLinesInfo(assemblyMachineInterface.ObjectInfo)
		var arrayOfTVUsers = assemblyInfo.TvUsers
		if util.HasInt(userId, arrayOfTVUsers) {
			lineUserCache = append(lineUserCache, assemblyMachineInterface.Id)
		}
	}
	return lineUserCache
}

func (v *MachineService) CreateAssemblyHmiEntry(projectId string, machineHmiInfo map[string]interface{}) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	serializedMachineInfo, _ := json.Marshal(machineHmiInfo)
	object := component.GeneralObject{
		ObjectInfo: serializedMachineInfo,
	}
	err, _ := Create(dbConnection, AssemblyMachineHmiTable, object)
	return err
}

func (v *MachineService) ArchivedAssemblyHmiEntry(projectId string, eventId int) error {
	var conditionString = "object_info ->> '$.eventId' = " + strconv.Itoa(eventId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfObjects, err := GetConditionalObjects(dbConnection, AssemblyMachineHmiTable, conditionString)

	if err != nil {
		return err
	}

	for _, assemblyObject := range *listOfObjects {
		var objectInfo = make(map[string]interface{})
		err = json.Unmarshal(assemblyObject.ObjectInfo, &objectInfo)

		if err != nil {
			continue
		}

		objectInfo["objectStatus"] = "Archived"
		var updatingData = make(map[string]interface{})
		_, serializedParamInfo := json.Marshal(objectInfo)
		updatingData["object_info"] = serializedParamInfo
		err = Update(dbConnection, AssemblyMachineHmiTable, assemblyObject.Id, updatingData)
		if err != nil {
			v.BaseService.Logger.Error("updating assembly machine hmi had failed", zap.String("error", err.Error()))
		}

	}
	return err
}

func (v *MachineService) GetActualValueFromAssemblyStatistic(projectId string, eventId int, timeEpoch int64) int {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var queryString = "SELECt * FROM assembly_machine_statistics where stats_info->>'$.eventId' = " + strconv.Itoa(eventId) + " and ts  > " + strconv.FormatInt(timeEpoch, 10) + " order by ts asc limit 1"
	v.BaseService.Logger.Info("select query for shift actual value", zap.Any("query", queryString))
	var statsData = make([]MachineStatistics, 0)
	dbConnection.Raw(queryString).Scan(&statsData)

	if len(statsData) > 0 {
		var lastStatData = MachineStatisticsInfo{}
		var err = json.Unmarshal(statsData[0].StatsInfo, &lastStatData)

		if err != nil {
			return 0
		}
		return lastStatData.Actual
	}
	return 0
}

func (v *MachineService) GetLastNHoursAssemblyStats(hours int) datatypes.JSON {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	if hours > 0 {
		var queryString = "SELECT ts,stats_info  FROM assembly_machine_statistics WHERE ts >= UNIX_TIMESTAMP(NOW() - INTERVAL " + strconv.Itoa(hours) + " HOUR) ORDER BY ts desc"
		type statsDate struct {
			Ts        int64          `json:"ts"`
			StatsInfo datatypes.JSON `json:"statsInfo"`
		}
		var statsData = make([]statsDate, 0)
		dbConnection.Raw(queryString).Scan(&statsData)

		serialisedData, _ := json.Marshal(statsData)
		return serialisedData
	} else {
		v.BaseService.Logger.Error("invalid hours passed in the getting last N hours stats")
		return []byte{}
	}

}

func (v *MachineService) GetListOfMachinesNeededHelp(timestamp int64) (error, []datatypes.JSON) {
	var fullMachineList []datatypes.JSON
	var dbConnection = v.BaseService.ServiceDatabases[ProjectID]
	if timestamp == 0 {
		// this is the first time happening, so return the last one
		err, assemblyMachineHelpSignalViewObjects := database.GetObjects(dbConnection, const_util.AssemblyMachineHelpSignalViewTable)
		if err != nil {
			return err, fullMachineList
		}
		for _, object := range assemblyMachineHelpSignalViewObjects {
			var assemblyMachineHelpSignalViewInfo = model.GetAssemblyMachineHelpSignalViewInfo(object.ObjectInfo)
			signalledMachine := map[string]interface{}{
				"id":                  object.Id,
				"signalGeneratedTime": assemblyMachineHelpSignalViewInfo.HelpButtonPressedTime,
			}
			serialisedSignalledMachine, err := json.Marshal(signalledMachine)
			if err != nil {
				v.BaseService.Logger.Error("Failed to marshal final object_info", zap.Int("machine_id", object.Id), zap.Error(err))
				continue
			}
			fullMachineList = append(fullMachineList, serialisedSignalledMachine)
		}
	} else {
		// now check all the machine updated greater than given time
		var condition = " object_info->>'$.helpButtonPressedTime' >  " + strconv.Itoa(int(timestamp))
		v.BaseService.Logger.Info("condition string to get the machine id ", zap.String("conditionString", condition))
		err, assemblyMachineHelpSignalViewObjects := database.GetConditionalObjects(dbConnection, const_util.AssemblyMachineHelpSignalViewTable, condition)
		if err != nil {
			return err, fullMachineList
		}
		v.BaseService.Logger.Info("list of machines which are signalled", zap.Int("length", len(assemblyMachineHelpSignalViewObjects)))
		for _, object := range assemblyMachineHelpSignalViewObjects {
			var assemblyMachineHelpSignalViewInfo = model.GetAssemblyMachineHelpSignalViewInfo(object.ObjectInfo)
			signalledMachine := map[string]interface{}{
				"id":                  object.Id,
				"signalGeneratedTime": assemblyMachineHelpSignalViewInfo.HelpButtonPressedTime,
			}
			serialisedSignalledMachine, err := json.Marshal(signalledMachine)
			if err != nil {
				v.BaseService.Logger.Error("Failed to marshal final object_info", zap.Int("machine_id", object.Id), zap.Error(err))
				continue
			}
			fullMachineList = append(fullMachineList, serialisedSignalledMachine)
		}
	}

	v.BaseService.Logger.Info("Final Machine List", zap.Any("Machines", fullMachineList))
	return nil, fullMachineList
}
func (v *MachineService) GetAssemblyMachineInfoFromEquipmentId(projectId string, equipmentId string) (error, int) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	var conditionString = " object_info->>'$.equipmentId'='" + equipmentId + "'"
	v.BaseService.Logger.Info("condition string to get the machine id from equipmentId", zap.String("conditionString", conditionString))
	err, listOfMachines := database.GetConditionalObjects(dbConnection, AssemblyMachineMasterTable, conditionString)
	if err != nil || len(listOfMachines) == 0 {
		return err, -1
	}
	return nil, listOfMachines[0].Id
}
func (v *MachineService) UpdateMachineParamEditStatus(projectId string, canEdit bool, paramId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	updatingData := make(map[string]interface{})
	_, machineParamGeneralObject := Get(dbConnection, MachineParameterTable, paramId)
	machinePram := MachineParameter{ObjectInfo: machineParamGeneralObject.ObjectInfo}
	machinePramObject := machinePram.getMachineParamInfo()
	machinePramObject.CanEdit = canEdit
	machinePramObject.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")

	serializedParamInfo, _ := json.Marshal(machinePramObject)
	updatingData["object_info"] = serializedParamInfo
	err := Update(dbConnection, MachineParameterTable, paramId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("updating machine param had failed", zap.String("error", err.Error()))
	}
	v.BaseService.Logger.Info("updating machine param", zap.Any("machine_param_id", paramId))

	return err
}

func (v *MachineService) GetListOfEnergyManagementMachines(projectId string) []int {

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	condition := " object_info ->>'$.enableCurrentSensorIntegration' = 'true'"

	mouldMachines, _ := GetConditionalObjects(dbConnection, MouldMachineSettingTable, condition)
	var machineList []int

	for _, mouldMachine := range *mouldMachines {
		machineList = append(machineList, mouldMachine.Id)
	}

	return machineList

}
