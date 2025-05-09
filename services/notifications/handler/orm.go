package handler

import (
	"cx-micro-flake/pkg/common/component"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func GetConditionalObjectsOrderBy(database *gorm.DB, table string, condition string, orderBy string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == NotificationRecordTrailTable:
		var dbObjects []NotificationRecordTrail
		if len(objectCount) > 0 {
			err = database.Debug().Model(&NotificationRecordTrail{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&NotificationRecordTrail{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
		}
		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}

		return nil, err
	case table == SystemNotificationTable:
		var dbObjects []SystemNotification
		if len(objectCount) > 0 {
			err = database.Debug().Model(&SystemNotification{}).Order(orderBy).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&SystemNotification{}).Order(orderBy).Where(condition).Find(&dbObjects).Error
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
func CreateRecordTrail(database *gorm.DB, objectInterface NotificationRecordTrail) (error, int) {
	err := database.Create(&objectInterface).Error
	return err, objectInterface.Id
}

func Create(database *gorm.DB, table string, objectInterface component.GeneralObject) (error, int) {
	var err error
	switch {
	case table == NotificationTable:
		object := Notification{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == PushNotificationTable:
		object := PushNotification{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	case table == SystemNotificationTable:
		object := SystemNotification{ObjectInfo: objectInterface.ObjectInfo}
		err = database.Create(&object).Error
		return err, object.Id
	default:
		return errors.New(CreateUnknownObjectType), -1
	}
}

func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	var err error

	switch {
	case table == NotificationTable:
		dbObject := Notification{Id: recordId}
		err = database.Debug().Model(&Notification{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	case table == SystemNotificationTable:
		dbObject := SystemNotification{Id: recordId}
		err = database.Debug().Model(&SystemNotification{}).Limit(100).Find(&dbObject).Error
		generalObject := component.GeneralObject{Id: dbObject.Id, ObjectInfo: dbObject.ObjectInfo}
		return err, generalObject
	default:
		return errors.New(GetUnknownObjectType), component.GeneralObject{}
	}

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (*[]component.GeneralObject, error) {
	var err error

	switch {
	case table == NotificationTable:

		var dbObjects []Notification
		if len(objectCount) > 0 {
			err = database.Model(&Notification{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&Notification{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == SystemNotificationTable:

		var dbObjects []SystemNotification
		if len(objectCount) > 0 {
			err = database.Debug().Model(&SystemNotification{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&SystemNotification{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == NotificationComponentTable:

		var dbObjects []NotificationComponent
		if len(objectCount) > 0 {
			err = database.Debug().Model(&NotificationComponent{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&NotificationComponent{}).Find(&dbObjects).Error
		}

		if err == nil {
			generalObjects := make([]component.GeneralObject, len(dbObjects))
			for i, v := range dbObjects {
				generalObjects[i] = component.GeneralObject{Id: v.Id, ObjectInfo: v.ObjectInfo}
			}
			return &generalObjects, err
		}
		return nil, err
	case table == PushNotificationTable:

		var dbObjects []PushNotification
		if len(objectCount) > 0 {
			err = database.Debug().Model(&PushNotification{}).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Debug().Model(&PushNotification{}).Find(&dbObjects).Error
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
	case table == NotificationTable:
		var dbObjects []Notification
		if len(objectCount) > 0 {
			err = database.Debug().Model(&Notification{}).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
		} else {
			err = database.Model(&Notification{}).Where(condition).Find(&dbObjects).Error
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

func Delete(database *gorm.DB, table string, objectInterface component.GeneralObject) error {
	var err error
	switch {
	case table == NotificationTable:
		err = database.Debug().Model(&Notification{}).Delete(&Notification{Id: objectInterface.Id}).Error
	default:
		return errors.New(GetUnknownObjectType)
	}

	return err

}

func Count(database *gorm.DB, table string) int64 {
	var err error
	var numberOfRecords int64
	switch {
	case table == NotificationTable:
		err = database.Debug().Model(&Notification{}).Count(&numberOfRecords).Error
	case table == SystemNotificationTable:
		err = database.Debug().Model(&SystemNotification{}).Count(&numberOfRecords).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	if err != nil {
		return -1
	} else {
		return numberOfRecords
	}
}

func Update(database *gorm.DB, table string, recordId int, updateObject map[string]interface{}) error {
	var err error
	switch {
	case table == NotificationTable:
		err = database.Debug().Model(&Notification{}).Take(&Notification{Id: recordId}).UpdateColumns(updateObject).Error
	case table == SystemNotificationTable:
		err = database.Debug().Model(&SystemNotification{}).Take(&SystemNotification{Id: recordId}).UpdateColumns(updateObject).Error
	case table == PushNotificationTable:
		err = database.Debug().Model(&PushNotification{}).Take(&PushNotification{Id: recordId}).UpdateColumns(updateObject).Error
	default:
		err = errors.New(UpdateUnknownObjectType)
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
	case table == NotificationTable:
		err = database.Debug().Model(&Notification{}).Take(&Notification{Id: objectInterface.Id}).UpdateColumns(updateObject).Error
	case table == SystemNotificationTable:
		err = database.Debug().Model(&SystemNotification{}).Take(&SystemNotification{Id: objectInterface.Id}).UpdateColumns(updateObject).Error

	default:
		err = errors.New(UpdateUnknownObjectType)
	}

	return err
}

func (v *NotificationService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(*listOfObjects)
					for _, referenceObject := range *listOfObjects {
						v.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}

		}
	}
}

func (v *NotificationService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			listOfObjects, err := GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(*listOfObjects) > 0 {
					for _, referenceObject := range *listOfObjects {
						fmt.Println("referenceTable : ", referenceTable, " id :", referenceObject)
						ArchiveObject(dbConnection, referenceTable, referenceObject)
						v.CreateUserRecordMessage(ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
						v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
					}
				}
			}

		}
	}
}
