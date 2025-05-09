package database

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/it/handler/const_util"
	"encoding/json"
	"gorm.io/gorm"
	"time"

	"gorm.io/datatypes"
)

type ITServiceRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type ITServiceComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceRequest struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceRequestCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}
type ITServiceRequestSubCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceSAPChangeReason struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceSAPAuthorizationFunction struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceCategoryTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type ITServiceWorkflowEngine struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *ITServiceWorkflowEngine) GetWorkFlowEngineInfo() *ITServiceWorkflowEngineInfo {
	itServiceRequestStatusInfo := ITServiceWorkflowEngineInfo{}
	json.Unmarshal(mwo.ObjectInfo, &itServiceRequestStatusInfo)
	return &itServiceRequestStatusInfo
}

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
}

//type ITServiceRequestInfo struct {
//	Name                 string          `json:"name"`
//	Labels               []string        `json:"labels"`
//	Remarks              string          `json:"remarks"`
//	CategoryId           int             `json:"categoryId"`
//	CreatedAt            string          `json:"createdAt"`
//	CreatedBy            int             `json:"createdBy"`
//	Attachment           string          `json:"attachment"`
//	ActionStatus         string          `json:"actionStatus"`
//	Description          string          `json:"description"`
//	RequestType          int             `json:"requestType"`
//	SubCategory          []int           `json:"subCategory"`
//	HasAttachment        bool            `json:"hasAttachment"`
//	LastUpdatedAt        string          `json:"lastUpdatedAt"`
//	LastUpdatedBy        int             `json:"lastUpdatedBy"`
//	ServiceStatus        int             `json:"serviceStatus"`
//	ServiceRequestId     string          `json:"serviceRequestId"`
//	DetailedDescription  string          `json:"detailedDescription"`
//	ActionRemarks        []ActionRemarks `json:"actionRemarks"`
//	CanUserSubmit        bool            `json:"canUserSubmit"`
//	CanUserAcknowledge   bool            `json:"canUserAcknowledge"`
//	CanHODApprove        bool            `json:"canHODApprove"`
//	CanHODReturn         bool            `json:"canHODReturn"`
//	CanHODReject         bool            `json:"canHODReject"`
//	CanEdit              bool            `json:"canEdit"`
//	CanExecPartyComplete bool            `json:"canExecPartyComplete"`
//	IsAssignedToEnabled  bool            `json:"isAssignedToEnabled"`
//	AssignedUser         int             `json:"assignedUser"`
//	ObjectStatus         string          `json:"objectStatus"`
//	isOnbehalfEnabled    bool            `json:"isOnbehalfEnabled"`
//	OnbehalfUserId       int             `json:"onbehalfUserId"`
//	TemplateFields       int             `json:"templateFields"`
//	HodEmail             string          `json:"hodEmail"`
//	IsAssignable         bool            `json:"isAssignable"`
//}

//func (v *ITServiceRequestInfo) Serialize() []byte {
//	rawData, _ := json.Marshal(v)
//	return rawData
//}
//
//func (v *ITServiceRequestInfo) DatabaseSerialize(userId int) map[string]interface{} {
//	updatingData := make(map[string]interface{})
//	v.LastUpdatedBy = userId
//	v.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
//	updatingData["object_info"] = v.Serialize()
//	return updatingData
//}
//
//func (mwo *ITServiceRequest) getServiceRequestInfo() *ITServiceRequestInfo {
//	itServiceRequestInfo := ITServiceRequestInfo{}
//	json.Unmarshal(mwo.ObjectInfo, &itServiceRequestInfo)
//	return &itServiceRequestInfo
//}

type ITServiceRequestStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func (mwo *ITServiceRequestStatus) GetRequestStatusInfo() *ITServiceRequestStatusInfo {
	itServiceRequestStatusInfo := ITServiceRequestStatusInfo{}
	json.Unmarshal(mwo.ObjectInfo, &itServiceRequestStatusInfo)
	return &itServiceRequestStatusInfo
}

type ITServiceRequestStatusInfo struct {
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

type ITServiceWorkflowEngineInfo struct {
	Name     string `json:"name"`
	Entities []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Parent           int    `json:"parent"`
		RoutingCondition []struct {
			ID    int64 `json:"id"`
			Query struct {
				Rules     []Rules `json:"rules"`
				Condition string  `json:"condition"`
			} `json:"query"`
			Config struct {
				Fields struct {
					Category struct {
						Name    string `json:"name"`
						Type    string `json:"type"`
						Value   string `json:"value"`
						Options []struct {
							Name  string `json:"name"`
							Value int    `json:"value"`
						} `json:"options"`
					} `json:"category"`
					SubCategory struct {
						Name    string `json:"name"`
						Type    string `json:"type"`
						Value   string `json:"value"`
						Options []struct {
							Name  string `json:"name"`
							Value int    `json:"value"`
						} `json:"options"`
					} `json:"sub_category"`
				} `json:"fields"`
			} `json:"config"`
			NotificationUserList     []int         `json:"notificationUserList"`
			NotificationUserListInfo []interface{} `json:"notificationUserListInfo"`
		} `json:"routingCondition"`
		IsEmailNotificationEnabled bool `json:"isEmailNotificationEnabled"`
	} `json:"entities"`
	CreatedAt       time.Time `json:"createdAt"`
	CreatedBy       int       `json:"createdBy"`
	Description     string    `json:"description"`
	ObjectStatus    string    `json:"objectStatus"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	LastUpdatedBy   int       `json:"lastUpdatedBy"`
	ListOfTemplates []int     `json:"listOfTemplates"`
}

type Rules struct {
	Field    string `json:"field"`
	Value    int    `json:"value"`
	Operator string `json:"operator"`
}

func GetWorkFlowUsers(dbConnection *gorm.DB, workflowEngineId int, userId int, categoryId int, entityIndex int) []int {
	err, serviceStatusObject := Get(dbConnection, const_util.ITServiceWorkflowEngineTable, workflowEngineId)
	authInterface := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userInfo := authInterface.GetUserInfoById(userId)
	siteId := userInfo.Site
	listOfRequestedUsers := make([]int, 0)

	if err == nil {
		serviceRequestStatus := ITServiceWorkflowEngine{ObjectInfo: serviceStatusObject.ObjectInfo}
		entityList := serviceRequestStatus.GetWorkFlowEngineInfo().Entities
		if len(entityList) >= entityIndex {

			routingConditionList := entityList[entityIndex-1].RoutingCondition
			//unmatchedRoutingUsers := serviceRequestStatus.getRequestStatusInfo().UnmatchedRoutingUsers
			if len(routingConditionList) > 0 {

				for _, routingCondition := range routingConditionList {
					categoryCheck := false
					siteCheck := false
					for _, rule := range routingCondition.Query.Rules {
						if rule.Field == "category" {
							if rule.Value == categoryId {
								categoryCheck = true
							}
						} else if rule.Field == "site" {
							if rule.Value == siteId {
								siteCheck = true
							}
						}

					}

					if categoryCheck && siteCheck {
						listOfRequestedUsers = append(listOfRequestedUsers, routingCondition.NotificationUserList...)
					}

				}
			}

		}

	}

	return listOfRequestedUsers
}
