package common

import (
	"cx-micro-flake/pkg/common/component"
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	MessageTypeComment      = "comment"
	MessageTypeNotification = "notification"
)

type TrackingFields struct {
	ChangedField string      `json:"changedField"`
	OldValue     interface{} `json:"oldValue"`
	NewValue     interface{} `json:"newValue"`
}

type TrackingFieldsResponse struct {
	ChangedField string `json:"changedField"`
	OldValue     string `json:"oldValue"`
	NewValue     string `json:"newValue"`
}

type ResourceMeta struct {
	UserId    int    `json:"userId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
type ResourceInfo struct {
	TrackingFields []TrackingFields `json:"trackingFields"`
	HasAttachment  bool             `json:"hasAttachment"`
	AttachmentList []string         `json:"attachmentList"`
	Message        string           `json:"message"`
	ResourceMeta   ResourceMeta     `json:"resourceMeta"`
	MessageType    string           `json:"messageType"`
	SourceObject   SourceObject     `json:"sourceObject"`
}

type SourceObject struct {
	Version    string                 `json:"version"`
	ObjectInfo map[string]interface{} `json:"objectInfo"`
}

type BotResource struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

type RecordMessageResponse struct {
	Date string              `json:"date"`
	Data []RecordMessageData `json:"data"`
}

type RecordMessageData struct {
	AvatarUrl             string                   `json:"avatarUrl"`
	UserId                int                      `json:"userId"`
	UserProfileRouterLink string                   `json:"userProfileRouterLink"`
	Username              string                   `json:"username"`
	Message               string                   `json:"message"`
	TrackingFields        []TrackingFieldsResponse `json:"trackingFields"`
	ReferenceTime         string                   `json:"referenceTime"`
	CreatedAt             string                   `json:"createdAt"`
}

func GetActionTrackingFields(actionFields map[string]interface{}) []TrackingFields {
	var trackFieldList = make([]TrackingFields, 0)

	for key, value := range actionFields {
		trackField := TrackingFields{
			ChangedField: key,
			OldValue:     "-",
			NewValue:     value,
		}
		trackFieldList = append(trackFieldList, trackField)
	}

	return trackFieldList
}
func GetTrackingFields(existingData *component.GeneralObject, updatedData *component.GeneralObject) []TrackingFields {
	existingKeyValue := make(map[string]interface{})
	var trackFieldList = make([]TrackingFields, 0)
	if existingData != nil {
		err := json.Unmarshal(existingData.ObjectInfo, &existingKeyValue)

		if err != nil {
			return trackFieldList
		}

	}

	updatingKeyValue := make(map[string]interface{})
	if updatedData != nil {
		err := json.Unmarshal(updatedData.ObjectInfo, &updatingKeyValue)

		if err != nil {
			return trackFieldList
		}
	}

	for key, existElement := range existingKeyValue {
		if updateElement, ok := updatingKeyValue[key]; ok {
			if !reflect.DeepEqual(existElement, updateElement) {
				trackField := TrackingFields{
					ChangedField: key,
					OldValue:     fmt.Sprintf("%v", existElement),
					NewValue:     fmt.Sprintf("%v", updateElement),
				}
				trackFieldList = append(trackFieldList, trackField)
			}
		}
	}

	return trackFieldList
}
