package handler

import (
	"crypto/tls"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/services/labour_management/handler/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func (v *NotificationService) GetComponents() []datatypes.JSON {
	var arrayOfObject []datatypes.JSON
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	listOfObjects, err := GetObjects(dbConnection, NotificationComponentTable)
	if err == nil {
		for _, objectInterface := range *listOfObjects {
			arrayOfObject = append(arrayOfObject, objectInterface.ObjectInfo)
		}
	}

	return arrayOfObject
}

func (v *NotificationService) CreateMessages(projectId string, messages []common.Message) error {
	// write into database
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	for _, message := range messages {
		var objectInfo = make(map[string]interface{})
		message.From = v.ServiceConfig.Email.FromEmail
		objectInfo["message"] = message
		objectInfo["retryCount"] = 0
		objectInfo["createdTs"] = time.Now().Unix()
		objectInfo["status"] = "pending"
		rawObject, _ := json.Marshal(objectInfo)
		generalObject := component.GeneralObject{ObjectInfo: rawObject}
		Create(dbConnection, NotificationTable, generalObject)
	}

	return nil
}

func (v *NotificationService) CreateSystemNotification(projectId string, notificationMessage datatypes.JSON) error {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	generalObject := component.GeneralObject{ObjectInfo: notificationMessage}
	err, notificationId := Create(dbConnection, SystemNotificationTable, generalObject)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	systemNotification := common.SystemNotification{}
	json.Unmarshal(notificationMessage, &systemNotification)
	for _, userId := range systemNotification.TargetUsers {
		v.BaseService.Logger.Info("notification is added by system to user", zap.Any("user_id", userId), zap.Any("notification_id", notificationId))
		authService.AddNotificationIds(userId, notificationId)
	}

	return err
}

func (v *NotificationService) CreatePushNotification(projectId string, notificationMessage datatypes.JSON) (int, error) {
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	generalObject := component.GeneralObject{ObjectInfo: notificationMessage}
	fmt.Println("notificationMessage", notificationMessage)
	err, notificationId := Create(dbConnection, PushNotificationTable, generalObject)
	if err != nil {
		v.BaseService.Logger.Error("error creating push notification", zap.Error(err))
		return 0, err
	}
	v.BaseService.Logger.Info("Push notification has been successfully created", zap.Any("push_notification_id", notificationId))
	return notificationId, nil
}

func (v *NotificationService) GetNotificationHistory(listOfId []int) (error, []component.GeneralObject) {
	stringIDs := make([]string, len(listOfId))
	for i, id := range listOfId {
		stringIDs[i] = fmt.Sprintf("%d", id)
	}
	condition := fmt.Sprintf(" id IN (%s)", strings.Join(stringIDs, ", "))
	dbConnection := v.BaseService.ServiceDatabases[ProjectID]
	err, i := database.GetConditionalObjects(dbConnection, PushNotificationTable, condition)
	return err, i
}

type OneSignalSubscriptionResponse struct {
	Properties struct {
		Language    string `json:"language"`
		TimezoneId  string `json:"timezone_id"`
		Country     string `json:"country"`
		FirstActive int    `json:"first_active"`
		LastActive  int    `json:"last_active"`
		Ip          string `json:"ip"`
	} `json:"properties"`
	Identity struct {
		OnesignalId string `json:"onesignal_id"`
	} `json:"identity"`
	Subscriptions []struct {
		Id                string `json:"id"`
		AppId             string `json:"app_id"`
		Type              string `json:"type"`
		Token             string `json:"token"`
		Enabled           bool   `json:"enabled"`
		NotificationTypes int    `json:"notification_types"`
		SessionTime       int    `json:"session_time"`
		SessionCount      int    `json:"session_count"`
		Sdk               string `json:"sdk"`
		DeviceModel       string `json:"device_model"`
		DeviceOs          string `json:"device_os"`
		Rooted            bool   `json:"rooted"`
		TestType          int    `json:"test_type"`
		AppVersion        string `json:"app_version"`
		NetType           int    `json:"net_type"`
		Carrier           string `json:"carrier"`
		WebAuth           string `json:"web_auth"`
		WebP256           string `json:"web_p256"`
	} `json:"subscriptions"`
}

func (v *NotificationService) GetOneSignalSubscriptionId(deviceToken string) []string {
	var listOfSubscriptions []string
	//https://api.onesignal.com/apps/c2be4cf0-9c68-4e9f-8446-25fd206f7261/users/by/onesignal_id/8070a246-a5e5-4a60-80dc-0fe24bd35a19
	var Url = v.ServiceConfig.PushNotificationConfig.OneSignalAppsAPI + "/" + v.ServiceConfig.PushNotificationConfig.AppID + "/users/by/onesignal_id/" + deviceToken
	v.BaseService.Logger.Info("sending request URL", zap.String("url", Url))
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		v.BaseService.Logger.Error("error creating request ", zap.Error(err))
		return listOfSubscriptions
	}

	req.Header.Set("Authorization", "Basic "+v.ServiceConfig.PushNotificationConfig.APIKey)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		v.BaseService.Logger.Error("error sending request", zap.Error(err))
		return listOfSubscriptions
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		v.BaseService.Logger.Error("getting subscription failed with status", zap.String("status", resp.Status))
		return listOfSubscriptions
	}

	var oneSignalSubscriptionResponse OneSignalSubscriptionResponse
	if body, err := io.ReadAll(resp.Body); err != nil {
		v.BaseService.Logger.Error("failed to decode response body", zap.Error(err))
		return listOfSubscriptions
	} else {
		json.Unmarshal(body, &oneSignalSubscriptionResponse)
		for _, subscription := range oneSignalSubscriptionResponse.Subscriptions {
			listOfSubscriptions = append(listOfSubscriptions, subscription.Id)
		}
	}
	v.BaseService.Logger.Info("retried list of subscriptions", zap.Any("subscriptions", listOfSubscriptions))
	return listOfSubscriptions

}
