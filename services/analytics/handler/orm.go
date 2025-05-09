package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

var emptyObject interface{}

func GetConditionalObjectsV2(database *gorm.DB, table string, condition string, objectCount ...int) (error, []component.GeneralObject) {

	var err error

	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Debug().Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Debug().Where(condition).Find(&dbObjects).Error
	}

	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}

}

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == AnalyticsRecordTrailTable:
		var dbObjects []AnalyticsRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&AnalyticsRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err

	default:
		return nil, errors.New(GetUnknownObjectType)
	}
}
func CreateRecordTrail(database *gorm.DB, objectInterface AnalyticsRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == AnalyticsDashboardTable:
		object := AnalyticsDashboard{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AnalyticsWidgetTable:
		object := AnalyticsWidget{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AnalyticsDataSourceTable:
		object := AnalyticsDatasource{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AnalyticsReportTable:
		object := AnalyticsReport{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AnalyticsComponentTable:
		object := AnalyticsComponent{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == SPCStatDataSourceTable:
		object := SpcStatDatasource{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == SPCResourceDataSourceTable:
		object := SpcResourceDatasource{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == AnalyticsDashboardPermissionTable:
		object := AnalyticsDashboardPermission{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == AnalyticsDashboardTable:
		dbObject := AnalyticsDashboard{Id: recordId}
		err = database.Debug().Model(&AnalyticsDashboard{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AnalyticsDatasourcesMasterTable:
		dbObject := AnalyticsDatasourcesMaster{Id: recordId}
		err = database.Debug().Model(&AnalyticsDatasourcesMaster{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AnalyticsWidgetTable:
		dbObject := AnalyticsWidget{Id: recordId}
		err = database.Debug().Model(&AnalyticsWidget{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AnalyticsDataSourceTable:
		dbObject := AnalyticsDatasource{Id: recordId}
		err = database.Debug().Model(&AnalyticsDatasource{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AnalyticsReportTable:
		dbObject := AnalyticsReport{Id: recordId}
		err = database.Debug().Model(&AnalyticsReport{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == AnalyticsDashboardPermissionTable:
		dbObject := AnalyticsDashboardPermission{Id: recordId}
		err = database.Debug().Model(&AnalyticsDashboardPermission{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == AnalyticsComponentTable:

		var dbObjects []AnalyticsComponent
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsDashboardTable:

		var dbObjects []AnalyticsDashboard
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsDashboard{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDashboard{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsDatasourcesMasterTable:

		var dbObjects []AnalyticsDatasourcesMaster
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsDatasourcesMaster{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDatasourcesMaster{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsWidgetTable:
		var dbObjects []AnalyticsWidget
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsWidget{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsWidget{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsDataSourceTable:
		var dbObjects []AnalyticsDatasource
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsDatasource{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDatasource{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsReportTable:
		var dbObjects []AnalyticsDatasource
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsReport{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsReport{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsDashboardPermissionTable:
		var dbObjects []AnalyticsDashboardPermission
		if len(objectCount) > 0 {
			err = database.Model(&AnalyticsDashboardPermission{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDashboardPermission{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	default:
		return nil, errors.New(GetUnknownObjectType)
	}

}

func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == AnalyticsDashboardTable:
		var dbObjects []AnalyticsDashboard
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsDashboard{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDashboard{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AnalyticsDatasourcesMasterTable:

		var dbObjects []AnalyticsDatasourcesMaster
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsDatasourcesMaster{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDatasourcesMaster{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsWidgetTable:
		var dbObjects []AnalyticsWidget
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsWidget{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsWidget{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AnalyticsDataSourceTable:
		var dbObjects []AnalyticsDatasource
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsDatasource{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDatasource{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AnalyticsReportTable:
		var dbObjects []AnalyticsReport
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsReport{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsReport{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == AnalyticsComponentTable:
		var dbObjects []AnalyticsComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsComponent{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsComponent{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == AnalyticsDashboardPermissionTable:
		var dbObjects []AnalyticsDashboardPermission
		if len(objectCount) > 0 {
			err = database.Debug().Model(&AnalyticsDashboardPermission{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&AnalyticsDashboardPermission{}).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	default:
		return nil, errors.New(GetUnknownObjectType)
	}
}
func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == AnalyticsDashboardTable:
		err = database.Debug().Model(&AnalyticsDashboard{}).Count(&numberOfRecords).Error
	case table == AnalyticsWidgetTable:
		err = database.Debug().Model(&AnalyticsWidget{}).Count(&numberOfRecords).Error
	case table == AnalyticsDataSourceTable:
		err = database.Debug().Model(&AnalyticsDatasource{}).Count(&numberOfRecords).Error
	case table == AnalyticsReportTable:
		err = database.Debug().Model(&AnalyticsReport{}).Count(&numberOfRecords).Error
	case table == AnalyticsDashboardPermissionTable:
		err = database.Debug().Model(&AnalyticsDashboardPermission{}).Count(&numberOfRecords).Error
	default:
		return -1
	}
	if err != nil {
		return -1
	}
	return numberOfRecords

}
func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == AnalyticsDashboardTable:
		err = database.Debug().Model(&AnalyticsDashboard{}).Delete(&AnalyticsDashboard{Id: objectInterface.Id}).Error
	case table == AnalyticsWidgetTable:
		err = database.Debug().Model(&AnalyticsWidget{}).Delete(&AnalyticsWidget{Id: objectInterface.Id}).Error
	case table == AnalyticsDataSourceTable:
		err = database.Debug().Model(&AnalyticsDatasource{}).Delete(&AnalyticsDatasource{Id: objectInterface.Id}).Error
	case table == AnalyticsReportTable:
		err = database.Debug().Model(&AnalyticsReport{}).Delete(&AnalyticsReport{Id: objectInterface.Id}).Error
	case table == SPCStatDataSourceTable:
		err = database.Debug().Model(&SpcStatDatasource{}).Delete(&SpcStatDatasource{Id: objectInterface.Id}).Error
	case table == SPCResourceDataSourceTable:
		err = database.Debug().Model(&SpcResourceDatasource{}).Delete(&SpcResourceDatasource{Id: objectInterface.Id}).Error
	case table == AnalyticsDashboardPermissionTable:
		err = database.Debug().Model(&AnalyticsDashboardPermission{}).Delete(&AnalyticsDashboardPermission{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == AnalyticsDashboardTable:
		err = database.Debug().Model(&AnalyticsDashboard{}).Take(&AnalyticsDashboard{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AnalyticsWidgetTable:
		err = database.Debug().Model(&AnalyticsWidget{}).Take(&AnalyticsWidget{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AnalyticsDataSourceTable:
		err = database.Debug().Model(&AnalyticsDatasource{}).Take(&AnalyticsDatasource{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AnalyticsReportTable:
		err = database.Debug().Model(&AnalyticsReport{}).Take(&AnalyticsReport{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AnalyticsComponentTable:
		err = database.Debug().Model(&AnalyticsComponent{}).Take(&AnalyticsComponent{Id: recordId}).UpdateColumns(updateObject).Error
	case table == AnalyticsDashboardPermissionTable:
		err = database.Debug().Model(&AnalyticsDashboardPermission{}).Take(&AnalyticsDashboardPermission{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func NoOfReferenceObjects(database *gorm.DB, table string, referenceField string, id int) int {
	conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$." + referenceField + "\")) =" + strconv.Itoa(id)
	listOfObjects, err := GetConditionalObjects(database, table, conditionString)
	if err != nil {
		return -1
	} else {
		return len(*listOfObjects)
	}
}

func ArchiveReferenceObjects(database *gorm.DB, table string, referenceField string, id int) error {
	var err error

	conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$." + referenceField + "\")) =" + strconv.Itoa(id)
	listOfObjects, err := GetConditionalObjects(database, table, conditionString)
	for _, object := range *listOfObjects {
		updateObject := make(map[string]interface{})
		var objectFields map[string]interface{}
		json.Unmarshal(object.ObjectInfo, &objectFields)
		objectFields["objectStatus"] = "Archived"
		objectFields["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
		serializedObject, _ := json.Marshal(objectFields)
		updateObject["object_info"] = serializedObject
		switch {
		case table == AnalyticsDashboardTable:
			err = database.Debug().Model(&AnalyticsDashboard{}).Take(&AnalyticsDashboard{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == AnalyticsWidgetTable:
			err = database.Debug().Model(&AnalyticsWidget{}).Take(&AnalyticsWidget{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == AnalyticsDataSourceTable:
			err = database.Debug().Model(&AnalyticsDatasource{}).Take(&AnalyticsDatasource{Id: object.Id}).UpdateColumns(updateObject).Error
		case table == AnalyticsDashboardPermissionTable:
			err = database.Debug().Model(&AnalyticsDashboardPermission{}).Take(&AnalyticsDashboardPermission{Id: object.Id}).UpdateColumns(updateObject).Error
		default:
			err = errors.New(UpdateUnknownObjectType)
		}
	}

	return err

}

func ArchiveObject(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	updateObject := make(map[string]interface{})
	var objectFields map[string]interface{}
	json.Unmarshal(objectInterface.ObjectInfo, &objectFields)
	objectFields["objectStatus"] = "Archived"
	serializedObject, _ := json.Marshal(objectFields)
	updateObject["object_info"] = serializedObject

	switch {
	case table == AnalyticsDashboardTable:
		err = database.Debug().Model(&AnalyticsDashboard{}).Take(&AnalyticsDashboard{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AnalyticsWidgetTable:
		err = database.Debug().Model(&AnalyticsWidget{}).Take(&AnalyticsWidget{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AnalyticsDataSourceTable:
		err = database.Debug().Model(&AnalyticsDatasource{}).Take(&AnalyticsDatasource{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == AnalyticsDashboardPermissionTable:
		err = database.Debug().Model(&AnalyticsDashboardPermission{}).Take(&AnalyticsDashboardPermission{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err

}

func (as *AnalyticsService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(*listOfObjects)
					for _, referenceObject := range *listOfObjects {
						as.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (as *AnalyticsService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := as.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := as.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					for _, referenceObject := range *listOfObjects {
						fmt.Println("referenceTable : ", referenceTable, " id :", referenceObject)
						ArchiveObject(dbConnection, referenceTable, referenceObject)
						as.CreateUserRecordMessage(ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						as.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
