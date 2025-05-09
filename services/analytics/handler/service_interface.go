package handler

import (
	"cx-micro-flake/pkg/common/component"
	"errors"
	"gorm.io/datatypes"
)

func (as *AnalyticsService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := as.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, AnalyticsComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (as *AnalyticsService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := as.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, AnalyticsComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, AnalyticsComponentTable, (*listOfObjects)[0].Id, updatingData)
		as.LoadInitComponents()
		return err
	}
	return err

}

func (as *AnalyticsService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := as.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, AnalyticsComponentTable, generalObject)
	if err == nil {
		as.LoadInitComponents()
	}
	return recordId, err
}
