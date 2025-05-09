package handler

import "gorm.io/datatypes"

type Notification struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type NotificationComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type SystemNotification struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type PushNotification struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type NotificationRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type NotificationListResponse struct {
	ViewCount     int         `json:"viewCount"`
	Notifications interface{} `json:"notifications"`
}
