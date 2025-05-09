package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type WorkHourAssignmentRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentRequest struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentTask struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type WorkHourAssignmentRequestSubCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentMasterMR struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentMasterTL struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type WorkHourAssignmentMasterJR struct {
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
type WorkHourAssignmentRequestInfo struct {
	Labels                []string        `json:"labels"`
	Remarks               string          `json:"remarks"`
	CreatedAt             string          `json:"createdAt"`
	CreatedBy             int             `json:"createdBy"`
	Attachment            string          `json:"attachment"`
	ActionStatus          string          `json:"actionStatus"`
	Description           string          `json:"description"`
	HasAttachment         bool            `json:"hasAttachment"`
	LastUpdatedAt         string          `json:"lastUpdatedAt"`
	LastUpdatedBy         int             `json:"lastUpdatedBy"`
	AssignmentStatus      int             `json:"assignmentStatus"`
	AssignmentReferenceId string          `json:"assignmentReferenceId"`
	WorkHour              int             `json:"workHour"`
	WorkMinute            int             `json:"workMinute"`
	AssignmentId          int             `json:"assignmentId"`
	AssignmentBase        int             `json:"assignmentBase"`
	ReportTime            string          `json:"reportTime"`
	TaskId                int             `json:"taskId"`
	DetailedDescription   string          `json:"detailedDescription"`
	ActionRemarks         []ActionRemarks `json:"actionRemarks"`
	CanUserSubmit         bool            `json:"canUserSubmit"`
	CanHODApprove         bool            `json:"canHODApprove"`
	CanHODReturn          bool            `json:"canHODReturn"`
	CanHODReject          bool            `json:"canHODReject"`
	CanEdit               bool            `json:"canEdit"`
	ObjectStatus          string          `json:"objectStatus"`
}

func (v *WorkHourAssignmentRequestInfo) Serialize() []byte {
	rawData, _ := json.Marshal(v)
	return rawData
}

func (v *WorkHourAssignmentRequestInfo) DatabaseSerialize(userId int) map[string]interface{} {
	updatingData := make(map[string]interface{})
	v.LastUpdatedBy = userId
	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	updatingData["object_info"] = v.Serialize()
	return updatingData
}

func (mwo *WorkHourAssignmentRequest) getServiceRequestInfo() *WorkHourAssignmentRequestInfo {
	workHourAssignmentRequestInfo := WorkHourAssignmentRequestInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workHourAssignmentRequestInfo)
	return &workHourAssignmentRequestInfo
}

type WorkHourAssignmentRequestStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *WorkHourAssignmentRequestStatus) getRequestStatusInfo() *WorkHourAssignmentRequestStatusInfo {
	workHourAssignmentRequestInfo := WorkHourAssignmentRequestStatusInfo{}
	json.Unmarshal(mwo.ObjectInfo, &workHourAssignmentRequestInfo)
	return &workHourAssignmentRequestInfo
}

type WorkHourAssignmentRequestStatusInfo struct {
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
