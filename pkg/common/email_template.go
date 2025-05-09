package common

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type EmailTemplateInfo struct {
	ReplyTo           string `json:"replyTo"`
	Subject           string `json:"subject"`
	Template          string `json:"template"`
	Description       string `json:"description"`
	ObjectStatus      string `json:"objectStatus"`
	TemplateName      string `json:"templateName"`
	TemplateType      int    `json:"templateType"`
	LastUpdatedAt     string `json:"lastUpdatedAt"`
	LastUpdatedBy     int    `json:"lastUpdatedBy"`
	IsTemplateEnabled bool   `json:"isTemplateEnabled"`
}

func GetEmailTemplateInfo(objectInfo datatypes.JSON) *EmailTemplateInfo {
	emailTemplateInfo := EmailTemplateInfo{}
	json.Unmarshal(objectInfo, &emailTemplateInfo)
	return &emailTemplateInfo
}

type ServiceFields struct {
	Name   string `json:"name"`
	Fields []struct {
		Type          string `json:"type"`
		Display       string `json:"display"`
		Property      string `json:"property"`
		ComponentName string `json:"componentName"`
	} `json:"fields"`
}

type TokenAttributes struct {
	MenuId       string   `json:"menuId"`
	InAppRouting bool     `json:"inAppRouting"`
	Fields       []string `json:"fields"`
}

type ModuleFields struct {
	Type                 string           `json:"type"`
	MenuId               string           `json:"menuId"`
	TokenAttributes      *TokenAttributes `json:"tokenAttributes"`
	Display              string           `json:"display"`
	IsUsedForTargetEmail bool             `json:"isUsedForTargetEmail"`
	Property             string           `json:"property"`
	LinkedField          string           `json:"linkedField"`
	ComponentName        string           `json:"componentName"`
	IsObjectField        bool             `json:"isObjectField"`
	IsForeignObjectMap   bool             `json:"isForeignObjectMap"`
	ForeignLinkedField   string           `json:"foreignLinkedField"`
}

type EmailTemplateField struct {
	ServiceFields []ServiceFields `json:"serviceFields"`
	ModuleFields  []ModuleFields  `json:"moduleFields"`
}

func GetEmailTemplateFields(objectInfo datatypes.JSON) EmailTemplateField {
	var emailTemplateFields EmailTemplateField
	json.Unmarshal(objectInfo, &emailTemplateFields)
	return emailTemplateFields
}
