package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/moulds/const_util"
	"cx-micro-flake/services/moulds/handler/database"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func (v *MouldService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	listOfObjects, err := database.GetObjects(dbConnection, const_util.MouldComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *MouldService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.MouldComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := database.Update(dbConnection, const_util.MouldComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *MouldService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := database.Create(dbConnection, const_util.MouldComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *MouldService) GetNoOfCavity(projectId string, mouldId int) (error, int) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

	if err != nil {
		return err, -1
	}

	var recordMap = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &recordMap)

	numberOfCavity := util.InterfaceToInt(recordMap["noOfCav"])
	return err, numberOfCavity
}

func (v *MouldService) PutToMaintenanceMode(projectId string, mouldId int, userId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

	if err != nil {
		return err
	}
	updatingData := make(map[string]interface{})

	var objectFields = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &objectFields)
	objectFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusMaintenance

	objectFields[const_util.CommonFieldLastUpdatedBy] = userId
	objectFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	serialisedData, _ := json.Marshal(objectFields)
	updatingData["object_info"] = serialisedData

	err = database.Update(dbConnection, const_util.MouldMasterTable, mouldId, updatingData)
	return err
}
func (v *MouldService) PutToActiveMode(projectId string, mouldId int, userId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

	if err != nil {
		return err
	}
	updatingData := make(map[string]interface{})

	var objectFields = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &objectFields)
	objectFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusActive

	objectFields[const_util.CommonFieldLastUpdatedBy] = userId
	objectFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	serialisedData, _ := json.Marshal(objectFields)
	updatingData["object_info"] = serialisedData

	err = database.Update(dbConnection, const_util.MouldMasterTable, mouldId, updatingData)
	return err
}
func (v *MouldService) PutToRepairMode(projectId string, mouldId int, userId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)

	if err != nil {
		return err
	}

	updatingData := make(map[string]interface{})

	var objectFields = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &objectFields)
	objectFields[const_util.MouldMasterFieldMouldStatus] = const_util.MouldStatusRepair

	objectFields[const_util.CommonFieldLastUpdatedBy] = userId
	objectFields[const_util.CommonFieldLastUpdatedAt] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	serialisedData, _ := json.Marshal(objectFields)
	updatingData["object_info"] = serialisedData

	err = database.Update(dbConnection, const_util.MouldMasterTable, mouldId, updatingData)

	return err
}

func (v *MouldService) UpdateShotCount(projectId string, mouldId int, shotCount int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	listOfObjects, err := database.GetMouldShoutCountByID(dbConnection, mouldId)

	if listOfObjects == nil || err != nil {
		var currentTime = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		newRecord := database.MouldShotCountViewInfo{
			CreatedAt:          currentTime,
			UpdatedAt:          currentTime,
			CreatedBy:          1,
			UpdatedBy:          1,
			CurrentShotCount:   shotCount,
			IsNotificationSent: false,
		}
		shotCountObject := component.GeneralObject{Id: mouldId, ObjectInfo: newRecord.Serialised()}
		err, _ = database.Create(dbConnection, const_util.MouldShoutCountViewTable, shotCountObject)
		if err != nil {
			v.BaseService.Logger.Error("error creating new mould shot count record", zap.Error(err))
			return fmt.Errorf("error creating new mould shout count record: %v", err)
		} else {
			v.BaseService.Logger.Info("creating the mould shot count", zap.Any("mould_id", mouldId), zap.Any("shotCount", shotCount))
		}

		return nil
	}
	mouldShoutCountInfo := database.GetMouldShoutCountViewInfo(listOfObjects.ObjectInfo)
	v.BaseService.Logger.Info("updating the mould shot count", zap.Any("mould_id", mouldId), zap.Any("time_since", mouldShoutCountInfo.UpdatedAt), zap.Any("shotCount", shotCount))
	mouldShoutCountInfo.CurrentShotCount = mouldShoutCountInfo.CurrentShotCount + shotCount
	mouldShoutCountInfo.UpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	mouldShoutCountInfo.UpdatedBy = 1

	var updateObject = make(map[string]interface{})
	updateObject["object_info"] = mouldShoutCountInfo.Serialised()

	err = database.Update(dbConnection, const_util.MouldShoutCountViewTable, mouldId, updateObject)
	if err != nil {
		v.BaseService.Logger.Error("error updating mould shot count record", zap.Error(err))
		return fmt.Errorf("error updating mould shout count: %v", err)
	}

	return nil
}

func (v *MouldService) GetMouldTestEventsForScheduler(projectId string) (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var conditionalStr = "object_info->>'$.objectStatus' = 'Active'"
	listOfEvents, err := database.GetConditionalObjects(dbConnection, const_util.MouldTestRequestTable, conditionalStr)
	listOfStatusCode, err := database.GetObjects(dbConnection, const_util.MouldTestStatusTable)

	var mouldTestRequestStatusCache = make(map[int]*database.MouldTestStatusInfo)
	for _, statusGeneralObject := range *listOfStatusCode {
		mouldStatus := database.MouldTestStatus{ObjectInfo: statusGeneralObject.ObjectInfo}
		mouldTestRequestStatusCache[statusGeneralObject.Id] = mouldStatus.GetMouldTestStatusInfo()
	}
	mouldIdTuple := "("
	for _, scheduledEventObject := range *listOfEvents {
		tempScheduledEvent := make(map[string]interface{})
		json.Unmarshal(scheduledEventObject.ObjectInfo, &tempScheduledEvent)
		mouldIdTuple += strconv.Itoa(util.InterfaceToInt(tempScheduledEvent["mouldId"])) + ","
	}
	mouldIdTuple = util.TrimSuffix(mouldIdTuple, ",")
	mouldIdTuple += ")"
	allChildOrderCondition := " id in " + mouldIdTuple
	listOfParentObjects, err := database.GetConditionalObjects(dbConnection, const_util.MouldMasterTable, allChildOrderCondition)

	var scheduledEvent = make(map[string]interface{})
	var arrayOfStatusObject []component.GeneralObject
	// here we need to send the action current status as front-end doesn't know the status id, and no need to know about the status id ,
	// so based on current order status id, we need to generate two flag, whether front-end action can be performed or not
	for _, scheduledEventObject := range *listOfEvents {
		json.Unmarshal(scheduledEventObject.ObjectInfo, &scheduledEvent)
		composedScheduleEvent := composeEventByMould(scheduledEvent)

		if value, ok := scheduledEvent["actionStatus"]; ok {
			actionStatus := util.InterfaceToString(value)
			mouldName := getMouldNameById(util.InterfaceToInt(scheduledEvent["mouldId"]), listOfParentObjects)
			composedScheduleEvent["canUpdate"] = true
			if actionStatus == const_util.MouldTestRequestActionCreated {
				composedScheduleEvent["canRelease"] = true
			} else {
				composedScheduleEvent["canRelease"] = false
			}
			if actionStatus == const_util.MouldTestRequestActionScheduled {
				composedScheduleEvent["canHold"] = true
			} else {
				composedScheduleEvent["canHold"] = false
			}
			if actionStatus == const_util.MouldTestRequestActionTestInProgress {
				composedScheduleEvent["canUpdate"] = false
			}
			if actionStatus == const_util.MouldTestRequestActionFailed {
				composedScheduleEvent["canUpdate"] = false
			}
			if actionStatus == const_util.MouldTestRequestActionPassed {
				composedScheduleEvent["canUpdate"] = false
			}
			mouldTestStatus := util.InterfaceToInt(scheduledEvent["mouldTestStatus"])

			composedScheduleEvent["canRelease"] = scheduledEvent["canRelease"]
			composedScheduleEvent["canUnRelease"] = scheduledEvent["canUnRelease"]
			composedScheduleEvent["canForceStop"] = scheduledEvent["canForceStop"]
			composedScheduleEvent["canComplete"] = false
			composedScheduleEvent["isAbortEnabled"] = scheduledEvent["isAbortEnabled"]
			testStatusInfo := mouldTestRequestStatusCache[mouldTestStatus]
			if testStatusInfo != nil {
				composedScheduleEvent["eventStatusName"] = testStatusInfo.Status
				composedScheduleEvent["eventColor"] = testStatusInfo.ColorCode
				composedScheduleEvent["eventSourceId"] = scheduledEventObject.Id
				composedScheduleEvent["machineId"] = util.InterfaceToInt(scheduledEvent["machineId"])
				composedScheduleEvent["name"] = mouldName

				serializedEventObject, _ := json.Marshal(composedScheduleEvent)
				arrayOfStatusObject = append(arrayOfStatusObject, component.GeneralObject{Id: scheduledEventObject.Id, ObjectInfo: serializedEventObject})
			}

		}

	}
	return err, &arrayOfStatusObject

}

func (v *MouldService) GetPartInfo(projectId string, partId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, partGeneralObject := database.Get(dbConnection, const_util.PartMasterTable, partId)
	return err, partGeneralObject
}

func (v *MouldService) GetMouldInfoById(projectId string, mouldId int) (error, component.GeneralObject) {
	if mouldId == 0 {
		return nil, component.GeneralObject{}
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, mouldGeneralObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)
	return err, mouldGeneralObject
}

func (v *MouldService) GetListOfMoulds(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := database.GetObjects(dbConnection, const_util.MouldMasterTable)
	var listOfObjects []component.GeneralObject
	if err == nil {
		if listOfMachines != nil {
			for _, object := range *listOfMachines {
				listOfObjects = append(listOfObjects, object)
			}
		}
	} else {
		v.BaseService.Logger.Error("error getting list of moulds", zap.Error(err))
	}

	return err, listOfObjects
}

func (v *MouldService) GetMouldsByPartNo(projectId string, partNo string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	condition := " JSON_CONTAINS(object_info, '{\"partNumber\":" + partNo + "}' , '$.partListArray') AND object_info ->> '$.objectStatus' = 'Active'"
	listOfMachines, err := database.GetConditionalObjects(dbConnection, const_util.MouldMasterTable, condition)
	var listOfObjects []component.GeneralObject
	if err == nil {
		if listOfMachines != nil {
			for _, object := range *listOfMachines {
				listOfObjects = append(listOfObjects, object)
			}
		}
	} else {
		v.BaseService.Logger.Error("error getting moulds by part number", zap.Error(err))
	}

	return err, listOfObjects
}

func (v *MouldService) GetMouldMachineTestParam(machineId int, mouldId int) int {
	// ask from proudction order service to get the scheduler order event
	var conditionString = " object_info->>'$.mouldId' = " + strconv.Itoa(mouldId) + " AND object_info->>'$.machineId' = " + strconv.Itoa(machineId) + " ORDER BY object_info->>'$.createdAt' DESC LIMIT 1"
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	objects, err := database.GetConditionalObjects(dbConnection, const_util.MouldTestRequestTable, conditionString)
	if err != nil {
		return -1
	}
	if objects != nil {
		if len(*objects) == 1 {
			var machineParamId = database.GetMouldTestRequestInfo((*objects)[0].ObjectInfo).MachineParamId
			v.BaseService.Logger.Info("got the machine param ID ", zap.Any("machine_param_id", machineParamId))
			return machineParamId
		}
	}

	return -1
}
func (v *MouldService) GetListOfMouldCategory(projectId string) (error, []component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	listOfMachines, err := database.GetObjects(dbConnection, const_util.MouldCategoryTable)
	var listOfObjects []component.GeneralObject
	for _, object := range *listOfMachines {
		listOfObjects = append(listOfObjects, object)
	}

	return err, listOfObjects
}
func (v *MouldService) GetBrandName(projectId string, brandId int) (error, string) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMachineMasterTable, brandId)

	if err != nil {
		return err, ""
	}

	var recordMap = make(map[string]interface{})
	json.Unmarshal(generalObject.ObjectInfo, &recordMap)

	machineBrand := util.InterfaceToString(recordMap["machineBrand"])
	return err, machineBrand
}
func (v *MouldService) MoveMouldToRepair(projectId string, mouldId int) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, const_util.MouldMasterTable, mouldId)
	if err == nil {
		mouldMaster := database.MouldMaster{ObjectInfo: generalObject.ObjectInfo}
		mouldMasterInfo := mouldMaster.GetMouldMasterInfo()
		mouldMasterInfo.MouldStatus = const_util.MouldStatusRepair
		updatingData := make(map[string]interface{})
		serializedObject, _ := json.Marshal(mouldMasterInfo)
		updatingData["object_info"] = serializedObject
		err = database.Update(v.BaseService.ReferenceDatabase, const_util.MouldMasterTable, mouldId, updatingData)
	}
	return err
}
func (v *MouldService) GetMouldShotCountViewForNotification() (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	condition := "is_notification_sent = 0"
	generalObject, err := database.GetConditionalObjects(dbConnection, const_util.MouldShoutCountViewTable, condition, 1)
	if err != nil {
		return err, nil
	}
	return err, generalObject

}
func (v *MouldService) GetMouldMaster() (error, *[]component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	condition := "is_notification_sent = 0"
	generalObject, err := database.GetConditionalObjects(dbConnection, const_util.MouldShoutCountViewTable, condition)
	if err != nil {
		return err, nil
	}
	return err, generalObject

}

func (v *MouldService) GetMouldSettingById(projectId string) (error, component.GeneralObject) {

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, mouldGeneralObject := database.Get(dbConnection, const_util.MouldSettingTable, 1)
	return err, mouldGeneralObject
}
func (v *MouldService) GetMouldNameByMachineId(machineId int) (error, string) { //err, mould name, start date, end date
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	condition := "object_info ->> '$.machineId' = " + strconv.Itoa(machineId) + " AND CURDATE() BETWEEN DATE(JSON_UNQUOTE(object_info->> '$.requestTestStartDate')) AND DATE(JSON_UNQUOTE(object_info->> '$.requestTestEndDate')) ORDER BY id DESC"
	generalObject, err := database.GetConditionalObjects(dbConnection, const_util.MouldTestRequestTable, condition, 1)

	if err != nil {
		return err, ""
	}
	if generalObject == nil || len(*generalObject) == 0 {
		return fmt.Errorf("no data found"), ""
	}
	mouldTestRequest := database.GetMouldTestRequestInfo((*generalObject)[0].ObjectInfo)

	startDate := util.ConvertStringToDateTime(mouldTestRequest.RequestTestStartDate)
	endDate := util.ConvertStringToDateTime(mouldTestRequest.RequestTestEndDate)

	timeNow := time.Now().Unix()

	if (startDate.DateTimeEpoch <= timeNow) && (endDate.DateTimeEpoch > timeNow) {

		err, mouldMasterInfo := database.Get(dbConnection, const_util.MouldMasterTable, mouldTestRequest.MouldId)

		if err != nil {

			return err, ""
		}

		mouldMaster := database.MouldMaster{ObjectInfo: mouldMasterInfo.ObjectInfo}
		mouldMasterInfoData := mouldMaster.GetMouldMasterInfo().ToolNo

		if util.InterfaceToString(mouldMasterInfoData) == "" {
			return fmt.Errorf("no data found"), ""
		}
		return nil, util.InterfaceToString(mouldMasterInfoData)
	}
	return err, ""

}
