package database

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/facility/handler/const_util"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"

	"gorm.io/datatypes"
)

type ActionRemarks struct {
	ExecutedTime  string `json:"executedTime"`
	Status        string `json:"status"`
	UserId        int    `json:"userId"`
	Remarks       string `json:"remarks"`
	ProcessedTime string `json:"processedTime"`
	RemarkType    string `json:"remarkType"`
}
type FacilityServiceRecordTrail struct {
	Id            int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	RecordId      int            `json:"recordId" gorm:"primary_key;not_null"`
	ComponentName string         `json:"componentName" gorm:"primary_key;not_null"`
	ObjectInfo    datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceComponent struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceRequest struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceRequestCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceWorkflowEngine struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceRequestSubCategory struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceSAPChangeReason struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceSAPAuthorizationFunction struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceEmailTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceCategoryTemplate struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceEmailTemplateField struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceAdminSetting struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

type FacilityServiceAdminSettingInfo struct {
	SecurityTeams   []int `json:"securityTeams"`
	EmailTemplateId int   `json:"emailTemplateId"`
}

func GetFacilityServiceAdminSettingInfo(objectInfo datatypes.JSON) *FacilityServiceAdminSettingInfo {
	facilityServiceAdminSettingInfo := FacilityServiceAdminSettingInfo{}
	json.Unmarshal(objectInfo, &facilityServiceAdminSettingInfo)
	return &facilityServiceAdminSettingInfo
}

func GetWorkFlowEngineInfo(objectInfo datatypes.JSON) *FacilityServiceWorkflowEngineInfo {
	facilityServiceRequestStatusInfo := FacilityServiceWorkflowEngineInfo{}
	json.Unmarshal(objectInfo, &facilityServiceRequestStatusInfo)
	return &facilityServiceRequestStatusInfo
}

type FacilityServiceRequestStatus struct {
	Id         int            `json:"id" gorm:"primary_key;auto_increment;not_null"`
	ObjectInfo datatypes.JSON `json:"objectInfo"`
}

func GetRequestStatusInfo(objectInfo datatypes.JSON) *FacilityServiceRequestStatusInfo {
	facilityServiceRequestStatusInfo := FacilityServiceRequestStatusInfo{}
	json.Unmarshal(objectInfo, &facilityServiceRequestStatusInfo)
	return &facilityServiceRequestStatusInfo
}

type FacilityServiceRequestStatusInfo struct {
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

type FacilityServiceRequestInfo struct {
	Name                      string          `json:"name"`
	Labels                    []string        `json:"labels"`
	Remarks                   string          `json:"remarks"`
	Category                  int             `json:"category"`
	SubCategory               []int           `json:"subCategory"`
	CreatedAt                 string          `json:"createdAt"`
	DeliveryDate              string          `json:"deliveryDate"`
	CreatedBy                 int             `json:"createdBy"`
	Attachment                string          `json:"attachment"`
	ActionStatus              string          `json:"actionStatus"`
	Description               string          `json:"description"`
	Building                  string          `json:"building"`
	RequestType               int             `json:"requestType"`
	HasAttachment             bool            `json:"hasAttachment"`
	LastUpdatedAt             string          `json:"lastUpdatedAt"`
	LastUpdatedBy             int             `json:"lastUpdatedBy"`
	ServiceStatus             int             `json:"serviceStatus"`
	ServiceRequestId          string          `json:"serviceRequestId"`
	DetailedDescription       string          `json:"detailedDescription"`
	ActionRemarks             []ActionRemarks `json:"actionRemarks"`
	CanUserSubmit             bool            `json:"canUserSubmit"`
	CanEdit                   bool            `json:"canEdit"`
	CanUserCancel             bool            `json:"canUserCancel"`
	CanUserAcknowledge        bool            `json:"canUserAcknowledge"`
	CanExecPartyComplete      bool            `json:"canExecPartyComplete"`
	CanEHSManagerApprove      bool            `json:"canEHSManagerApprove"`
	CanEHSManagerReject       bool            `json:"canEHSManagerReject"`
	CanFacilityManagerApprove bool            `json:"canFacilityManagerApprove"`
	CanFacilityManagerReject  bool            `json:"canFacilityManagerReject"`
	CanHODApprove             bool            `json:"canHODApprove"`
	CanHODReject              bool            `json:"canHODReject"`
	Comment                   string          `json:"comment"`
	ObjectStatus              string          `json:"objectStatus"`
	HoursTaken                string          `json:"hoursTaken"`
	AssignmentId              int             `json:"assignmentId"`
	Priority                  string          `json:"priority"`
	OnbehalfUserId            int             `json:"onbehalfUserId"`
	IsOnbehalfEnabled         bool            `json:"isOnbehalfEnabled"`
}
type FacilityServiceWorkflowEngineInfo struct {
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
	err, serviceStatusObject := Get(dbConnection, const_util.FacilityServiceWorkflowEngineTable, workflowEngineId)
	authInterface := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userInfo := authInterface.GetUserInfoById(userId)
	siteId := userInfo.Site
	listOfRequestedUsers := make([]int, 0)

	if err == nil {
		workFlowEngineInfo := GetWorkFlowEngineInfo(serviceStatusObject.ObjectInfo)
		entityList := workFlowEngineInfo.Entities
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
						fmt.Println("listOfRequestedUsers", listOfRequestedUsers)

					}

				}
			}

		}

	}

	return listOfRequestedUsers
}
