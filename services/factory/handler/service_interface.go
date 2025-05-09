package handler

import (
	"cx-micro-flake/pkg/common/component"
	"errors"
	"strconv"

	"gorm.io/datatypes"
)

func (v *FactoryService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, FactoryComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *FactoryService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetConditionalObjects(dbConnection, FactoryComponentTable, conditionQuery)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(*listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := Update(dbConnection, FactoryComponentTable, (*listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *FactoryService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := Create(dbConnection, FactoryComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *FactoryService) GetSections(departmentId int) ([]int, error) {
	listofSections := make([]int, 0)
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	var conditionQuery = " object_info->>'$.department' = '" + strconv.Itoa(departmentId) + "'"
	listOfObjects, err := GetConditionalObjects(dbConnection, FactoryDepartmentSectionTable, conditionQuery)

	if err != nil {
		return listofSections, err
	}
	for _, factoryDepartmentSectionInterfaceObject := range *listOfObjects {
		listofSections = append(listofSections, factoryDepartmentSectionInterfaceObject.Id)
	}

	return listofSections, err
}
