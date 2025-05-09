package notification

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

func (v *EmailHandler) CreateMyDepartmentSystemNotification(projectId string, hodEmail string, name string, requestId string, resourceId int) {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	hodUserId := authService.EmailToUserId(hodEmail)

	systemNotification := common.SystemNotification{}
	systemNotification.Name = name
	systemNotification.ColorCode = "#14F44E"
	systemNotification.IconCls = "icon-park-outline:transaction-order"
	systemNotification.RecordId = resourceId
	systemNotification.RouteLinkComponent = "it_service_my_department_request"
	systemNotification.Component = "IT Service"
	systemNotification.Description = "New request [" + requestId + "] has been submitted for your approval."
	systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
	systemNotification.TargetUsers = append(systemNotification.TargetUsers, hodUserId)
	rawSystemNotification, _ := json.Marshal(systemNotification)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	notificationService.CreateSystemNotification(projectId, rawSystemNotification)
	v.Logger.Info("system notification is generated for HOD Users", zap.Any("hod_email", hodEmail), zap.Any("resource_id", resourceId))

}

func (v *EmailHandler) CreateHodApproveNotification(projectId string, hodEmail string, name string, requestId string, resourceId int) {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	hodUserId := authService.EmailToUserId(hodEmail)

	systemNotification := common.SystemNotification{}
	systemNotification.Name = name
	systemNotification.ColorCode = "#14F44E"
	systemNotification.IconCls = "icon-park-outline:transaction-order"
	systemNotification.RecordId = resourceId
	systemNotification.RouteLinkComponent = "it_service_my_department_request"
	systemNotification.Component = "IT Service"
	systemNotification.Description = "New request [" + requestId + "] has been submitted for your approval."
	systemNotification.GeneratedTime = util.GetCurrentTime(time.RFC822)
	systemNotification.TargetUsers = append(systemNotification.TargetUsers, hodUserId)
	rawSystemNotification, _ := json.Marshal(systemNotification)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	notificationService.CreateSystemNotification(projectId, rawSystemNotification)
	v.Logger.Info("system notification is generated for HOD Users", zap.Any("hod_email", hodEmail), zap.Any("resource_id", resourceId))

}
