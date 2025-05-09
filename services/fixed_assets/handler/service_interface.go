package handler

import (
	"cx-micro-flake/pkg/common/component"
	"errors"
	"gorm.io/datatypes"
)

func (v *FixedAssetService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, FixedAssetsComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *FixedAssetService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, FixedAssetsComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, FixedAssetsComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *FixedAssetService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, FixedAssetsComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}
