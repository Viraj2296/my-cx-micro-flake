package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"github.com/gin-gonic/gin"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	NotificationTable            = "notification"
	SystemNotificationTable      = "system_notification"
	NotificationComponentTable   = "notification_component"
	NotificationRecordTrailTable = "notification_record_trail"

	SystemNotificationComponent = "system_notification"

	ErrorGettingObjectsInformation = 10000
	ErrorGettingObjectInformation  = 10001
	ErrorUpdatingObjectInformation = 10002
	PushNotificationTable          = "push_notification"
	ProjectID                      = "906d0fd569404c59956503985b330132"

	ActionMarkAllRead = "mark_all_read"
)

func sendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func sendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}
