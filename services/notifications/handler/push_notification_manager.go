package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// NotificationPayload is the structure of the notification payload for OneSignal
type NotificationPayload struct {
	AppID            string                 `json:"app_id"`
	IncludePlayerIDs []string               `json:"IncludePlayerIDs"`
	Headings         map[string]string      `json:"headings"`
	Contents         map[string]string      `json:"contents"`
	Data             map[string]interface{} `json:"data,omitempty"`
}

type NotificationMessage struct {
	IncludePlayerIDs  []string               `json:"IncludePlayerIDs"`
	Headings          map[string]string      `json:"headings"`
	Contents          map[string]string      `json:"contents"`
	Data              map[string]interface{} `json:"data,omitempty"`
	CreatedAt         string                 `json:"created_at,omitempty"`
	CreatedBy         string                 `json:"created_by,omitempty"`
	LastUpdatedAt     string                 `json:"LastUpdatedAt,omitempty"`
	LastUpdatedBy     string                 `json:"lastUpdatedBy,omitempty"`
	RetryCount        int                    `json:"retryCount"`
	DeliveryStatus    string                 `json:"deliveryStatus"`
	Diagnostics       map[string]interface{} `json:"diagnostics,omitempty"`
	DiagnosticMessage string                 `json:"diagnosticMessage,omitempty"`
	ReferenceData     map[string]string      `json:"referenceData"` // this is reference object ID we can pass, so that later we can use this
}

func (v *NotificationMessage) getSerialised() datatypes.JSON {
	serialisedData, _ := json.Marshal(v)
	return serialisedData
}
func getNotificationMessage(serialisedData datatypes.JSON) *NotificationMessage {
	notificationMessage := &NotificationMessage{}
	json.Unmarshal(serialisedData, notificationMessage)
	return notificationMessage
}

func (v *NotificationManager) getNotificationPayload(notificationMessage *NotificationMessage) NotificationPayload {
	payload := NotificationPayload{
		AppID:            v.PushNotificationConfig.AppID,
		IncludePlayerIDs: notificationMessage.IncludePlayerIDs,
		Headings:         notificationMessage.Headings,
		Contents:         notificationMessage.Contents,
	}
	return payload
}

type SendingPayload struct {
	AppId                  string            `json:"app_id"`
	Contents               map[string]string `json:"contents"`
	Headings               map[string]string `json:"headings"`
	TargetChannel          string            `json:"target_channel"`
	IncludeSubscriptionIds []string          `json:"include_subscription_ids"`
}

func (v *NotificationManager) sendPushNotification(payload NotificationPayload) error {
	var sendingPayload SendingPayload
	if len(payload.IncludePlayerIDs) == 0 {
		v.Logger.Warn("player ids are empty, can not dispatch push notification ...")
		return nil
	}
	sendingPayload.AppId = v.PushNotificationConfig.AppID
	sendingPayload.Contents = payload.Contents
	sendingPayload.Headings = payload.Headings
	sendingPayload.TargetChannel = "push"
	sendingPayload.IncludeSubscriptionIds = payload.IncludePlayerIDs

	payloadBytes, err := json.Marshal(sendingPayload)
	if err != nil {
		v.Logger.Error("error unmarshalling payload ", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", v.PushNotificationConfig.OneSignalAPI, bytes.NewBuffer(payloadBytes))
	if err != nil {
		v.Logger.Error("error creating request ", zap.Error(err))
		return err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+v.PushNotificationConfig.APIKey)

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
		v.Logger.Error("error sending request ", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		v.Logger.Error("notification send failed with status", zap.String("status", resp.Status))
		return errors.New("Notification send failed with status: " + resp.Status)
	}

	v.Logger.Info("Notification sent successfully!")
	return nil
}
func (v *NotificationManager) PoolPushNotification() {
	v.Logger.Info("starting push notification pooling")
	tDuration, err := time.ParseDuration(v.PushNotificationConfig.PoolingInterval)
	if err != nil {
		tDuration, err = time.ParseDuration("10s")
		if err != nil {
			v.Logger.Error("error parsing pooling interval, using default value", zap.String("error", err.Error()))
			tDuration = 10 * time.Second
		}
	}
	for {
		v.Logger.Info("getting push notification objects")
		listOfObjects, err := GetObjects(v.ServiceDatabase, PushNotificationTable)
		if err == nil {
			for _, object := range *listOfObjects {
				notificationMessage := getNotificationMessage(object.ObjectInfo)
				if notificationMessage.DeliveryStatus == "pending" && notificationMessage.RetryCount < 2 {
					payload := v.getNotificationPayload(notificationMessage)
					sendingError := v.sendPushNotification(payload)
					updatingData := make(map[string]interface{})
					if sendingError == nil {
						notificationMessage.DeliveryStatus = "sent"
						notificationMessage.DiagnosticMessage = "successfully sent"

						updatingData["object_info"] = notificationMessage.getSerialised()
						err = Update(v.ServiceDatabase, PushNotificationTable, object.Id, updatingData)
						if err != nil {
							v.Logger.Error("error updating push notification sent status notification", zap.String("error", err.Error()))
						} else {
							v.Logger.Info("Push notification has been successfully sent", zap.Any("push_notification_id", object.Id))
						}

					} else {
						v.Logger.Error("error sending push notification", zap.String("error", sendingError.Error()))
						notificationMessage.RetryCount = notificationMessage.RetryCount + 1
						notificationMessage.DiagnosticMessage = sendingError.Error()
						updatingData["object_info"] = notificationMessage.getSerialised()
						err = Update(v.ServiceDatabase, PushNotificationTable, object.Id, updatingData)
						if err != nil {
							v.Logger.Error("error updating push notification sent status notification", zap.String("error", err.Error()))
						} else {
							v.Logger.Info("Push notification has been successfully updated", zap.Any("push_notification_id", object.Id))
						}
					}
				}
			}
		} else {
			v.Logger.Error("error getting push notification objects", zap.String("error", err.Error()))
		}
		time.Sleep(tDuration)
	}
}
