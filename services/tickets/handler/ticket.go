package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type TicketsServiceRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type TicketsServiceComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type TicketsServiceRequest struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type TicketsServiceRequestCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type TicketsServiceRequestSubCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type TicketsServiceEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type TicketsServiceEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}
type TicketsServiceRequestInfo struct {
	Name                 string          `json:"name"`
	Version              string          `json:"version"`
	ResolveTime          string          `json:"resolveTime"`
	Labels               []string        `json:"labels"`
	Remarks              string          `json:"remarks"`
	Category             int             `json:"category"`
	CreatedAt            string          `json:"createdAt"`
	CreatedBy            int             `json:"createdBy"`
	Attachment           string          `json:"attachment"`
	ActionStatus         string          `json:"actionStatus"`
	Description          string          `json:"description"`
	RequestType          int             `json:"requestType"`
	SubCategory          int             `json:"subCategory"`
	HasAttachment        bool            `json:"hasAttachment"`
	LastUpdatedAt        string          `json:"lastUpdatedAt"`
	LastUpdatedBy        int             `json:"lastUpdatedBy"`
	ServiceStatus        int             `json:"serviceStatus"`
	ServiceRequestId     string          `json:"serviceRequestId"`
	DetailedDescription  string          `json:"detailedDescription"`
	ActionRemarks        []ActionRemarks `json:"actionRemarks"`
	CanUserSubmit        bool            `json:"canUserSubmit"`
	CanEdit              bool            `json:"canEdit"`
	CanUserAcknowledge   bool            `json:"canUserAcknowledge"`
	CanHODApprove        bool            `json:"canHODApprove"`
	CanHODReturn         bool            `json:"canHODReturn"`
	CanHODReject         bool            `json:"canHODReject"`
	CanITApprove         bool            `json:"canITApprove"`
	CanITReject          bool            `json:"canITReject"`
	CanExecPartyComplete bool            `json:"canExecPartyComplete"`
	ObjectStatus         string          `json:"objectStatus"`
}

func (v *TicketsServiceRequestInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *TicketsServiceRequestInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (mwo *TicketsServiceRequest) getServiceRequestInfo() *TicketsServiceRequestInfo {
	ticketsServiceRequestInfo := TicketsServiceRequestInfo{}
	json.Unmarshal(mwo.ObjectInfo, &ticketsServiceRequestInfo)
	return &ticketsServiceRequestInfo
}

type TicketsServiceRequestStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *TicketsServiceRequestStatus) getRequestStatusInfo() *TicketsServiceRequestStatusInfo {
	ticketsServiceRequestStatusInfo := TicketsServiceRequestStatusInfo{}
	json.Unmarshal(mwo.ObjectInfo, &ticketsServiceRequestStatusInfo)
	return &ticketsServiceRequestStatusInfo
}

type TicketsServiceRequestStatusInfo struct {
	Status           string      `json:"status"`
	ColorCode        string      `json:"colorCode"`
	CreatedAt        time.Time   `json:"createdAt"`
	CreatedBy        int         `json:"createdBy"`
	Description      string      `json:"description"`
	ObjectStatus     string      `json:"objectStatus"`
	LastUpdatedAt    time.Time   `json:"lastUpdatedAt"`
	LastUpdatedBy    int         `json:"lastUpdatedBy"`
	EmailTemplateId  interface{} `json:"emailTemplateId"`
	RoutingCondition []struct {
		Id    string `json:"id"`
		Query struct {
			Rules []struct {
				Field    string `json:"field"`
				Value    int    `json:"value"`
				Operator string `json:"operator"`
			} `json:"rules"`
			Condition string `json:"condition"`
		} `json:"query"`
		NotificationUserList []int `json:"notificationUserList"`
	} `json:"routingCondition"`
	AuthorisedUserList         []int `json:"authorisedUserList"`
	UnmatchedRoutingUsers      []int `json:"unmatchedRoutingUsers"`
	IsEmailNotificationEnabled bool  `json:"isEmailNotificationEnabled"`
}

func contain(arr []int, searchedValue int) bool {
	for _, intObj := range arr {
		if intObj == searchedValue {
			return true
		}
	}
	return false
}
