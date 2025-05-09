package handler

import (
	"cx-micro-flake/pkg/common/component"
	"errors"
	"gorm.io/datatypes"
)

func (v *ManufacturingService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, ManufacturingComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *ManufacturingService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, ManufacturingComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, ManufacturingComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *ManufacturingService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, ManufacturingComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *ManufacturingService) GetAssemblyPartInfo(projectId string, partId int) (error, component.GeneralObject) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, partGeneralObject := Get(dbConnection, ManufacturingAssemblyPartTable, partId)
	return err, partGeneralObject
}
