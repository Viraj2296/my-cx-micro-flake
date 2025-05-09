package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"errors"
	"gorm.io/datatypes"
)

func (v *BatchManagementService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, listOfObjects := database.GetObjects(dbConnection, const_util.BatchManagementComponentTable)
	if err == nil {
		for _, objectInterface := range listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *BatchManagementService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, listOfObjects := database.GetConditionalObjects(dbConnection, const_util.BatchManagementComponentTable, conditionQuery)
	if err == nil {
		if len(listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := database.Update(dbConnection, const_util.BatchManagementComponentTable, (listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *BatchManagementService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := database.CreateFromGeneralObject(dbConnection, const_util.BatchManagementComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *BatchManagementService) GetListOfMouldBatch(mouldBatchId string) []component.GeneralObject {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	var conditionString = " object_info->>'$.mouldBatchId' = '" + mouldBatchId + "'"
	err, listOfMouldBatch := database.GetConditionalObjects(dbConnection, const_util.BatchManagementMouldTable, conditionString)
	if err != nil {
		return make([]component.GeneralObject, 0)
	}
	return listOfMouldBatch
}

func (v *BatchManagementService) GetBatchRawMaterial(ramMaterialId int) component.GeneralObject {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, generalObject := database.Get(dbConnection, const_util.BatchManagementRawMaterialTable, ramMaterialId)
	if err != nil {
		return component.GeneralObject{}
	}
	return generalObject
}

func (v *BatchManagementService) GetMouldBatch(recordId int) component.GeneralObject {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, generalObject := database.Get(dbConnection, const_util.BatchManagementMouldTable, recordId)
	if err != nil {
		return component.GeneralObject{}
	}
	return generalObject
}
