package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/it/handler/const_util"
	"cx-micro-flake/services/it/handler/database"
	"errors"
	"fmt"
	"gorm.io/datatypes"
)

func (v *ITService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	listOfObjects, err := database.GetObjects(dbConnection, const_util.ITServiceComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *ITService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.ITServiceComponentTable, conditionQuery)
	fmt.Println("condition query :", conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := database.Update(dbConnection, const_util.ITServiceComponentTable, (*listOfObjects)[0].Id, updatingData)

		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *ITService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := database.Create(dbConnection, const_util.ITServiceComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}
