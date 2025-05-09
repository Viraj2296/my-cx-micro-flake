package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func (v *TraceabilityService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, TraceabilityComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *TraceabilityService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, TraceabilityComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, TraceabilityComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *TraceabilityService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, TraceabilityComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *TraceabilityService) CreateTraceabilityResource(mouldBatchId string, qaStartTime, qaEndTime string, qaPerson int) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	batchManagementInterface := common.GetService("batch_management_module").ServiceInterface.(common.BatchManagementInterface)
	var listOfMouldBatch = batchManagementInterface.GetListOfMouldBatch(mouldBatchId)
	v.BaseService.Logger.Info("loading all the mould batch records", zap.Any("mould_batch_records", listOfMouldBatch))
	for _, mouldBatchInterface := range listOfMouldBatch {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(mouldBatchInterface.ObjectInfo, &objectFields)
		objectFields["qaStartTime"] = qaStartTime
		objectFields["qaCompleteTime"] = qaEndTime
		objectFields["createdBy"] = qaPerson
		serialisedObject, _ := json.Marshal(objectFields)
		generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
		v.BaseService.Logger.Info("creating the traceability record", zap.Any("object", serialisedObject))
		err, recordId := Create(dbConnection, TraceabilityOrdersTable, generalObject)
		if err != nil {
			v.BaseService.Logger.Error("error creating traceability record", zap.String("error", err.Error()))
		} else {
			v.BaseService.Logger.Info("created traceability record", zap.Int("record", recordId))
		}

	}
}
