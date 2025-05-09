package handler

import (
	"cx-micro-flake/pkg/common/component"
	"errors"
	"gorm.io/datatypes"
)

func (ts *TicketsService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := ts.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, TicketsServiceComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (ts *TicketsService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := ts.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, TicketsServiceComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, TicketsServiceComponentTable, (*listOfObjects)[0].Id, updatingData)
		ts.LoadInitComponents()
		return err
	}
	return err

}

func (ts *TicketsService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := ts.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, TicketsServiceComponentTable, generalObject)
	if err == nil {
		ts.LoadInitComponents()
	}
	return recordId, err
}
