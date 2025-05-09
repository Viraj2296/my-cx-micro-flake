package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/labour_management/handler/const_util"
	"cx-micro-flake/services/labour_management/handler/database"
	"errors"
	"fmt"
	"strings"

	"gorm.io/datatypes"
)

func (v *LabourManagementService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, listOfObjects := database.GetObjects(dbConnection, const_util.LabourManagementComponentTable)
	if err == nil {
		for _, objectInterface := range listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *LabourManagementService) UpdateComponent(componentName, targetTable string, serialisedObject datatypes.JSON) error {
	var conditionQuery = " object_info->>'$.componentName' = '" + componentName + "' AND object_info->'$.targetTable' = '" + targetTable + "'"
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	err, listOfObjects := database.GetConditionalObjects(dbConnection, const_util.LabourManagementComponentTable, conditionQuery)
	if err == nil {
		if len(listOfObjects) == 0 {
			return errors.New("system couldn't able to find any components")
		}
		if len(listOfObjects) > 1 {
			return errors.New("system found more than required resources with the same component name")
		}

		updatingData := make(map[string]interface{})
		updatingData["object_info"] = serialisedObject
		err := database.Update(dbConnection, const_util.LabourManagementComponentTable, (listOfObjects)[0].Id, updatingData)
		v.LoadInitComponents()
		return err
	}
	return err

}

func (v *LabourManagementService) CreateComponent(serialisedObject datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	generalObject := component.GeneralObject{ObjectInfo: serialisedObject}
	err, recordId := database.CreateFromGeneralObject(dbConnection, const_util.LabourManagementComponentTable, generalObject)
	if err == nil {
		v.LoadInitComponents()
	}
	return recordId, err
}

func (v *LabourManagementService) GetMobileAllowedJobRoles() []int {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	var listOfMobileAllowedRoles = make([]int, 0)
	err, c := database.Get(dbConnection, const_util.LabourManagementSettingTable, 1)
	if err != nil {
		return listOfMobileAllowedRoles
	}
	shiftSettingInfo := database.GetLabourManagementSettingInfo(c.ObjectInfo)
	for _, allowedRoles := range shiftSettingInfo.ShiftCreationRoles {
		listOfMobileAllowedRoles = append(listOfMobileAllowedRoles, allowedRoles)
	}

	return listOfMobileAllowedRoles
}

func (v *LabourManagementService) GetLabourManagementShiftTemplate(templateId []int) []component.GeneralObject {
	dbConnection := v.BaseService.ServiceDatabases[const_util.ProjectID]
	var idStrings []string
	for _, id := range templateId {
		idStrings = append(idStrings, fmt.Sprintf("%d", id))
	}
	conditionString := fmt.Sprintf("id IN (%s)", strings.Join(idStrings, ", "))
	err, listOfMouldBatch := database.GetConditionalObjects(dbConnection, const_util.LabourManagementShiftTemplateTable, conditionString)
	if err != nil {
		return make([]component.GeneralObject, 0)
	}
	return listOfMouldBatch
}
