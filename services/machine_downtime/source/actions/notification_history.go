package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/dto"
	"cx-micro-flake/services/machine_downtime/source/models"
	"encoding/json"
	"fmt"
	"go.cerex.io/transcendflow/orm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type NotificationMessage struct {
	IncludePlayerIDs  []string               `json:"IncludePlayerIDs"`
	Headings          map[string]string      `json:"headings"`
	Contents          map[string]string      `json:"contents"`
	Data              map[string]interface{} `json:"data,omitempty"`
	CreatedAt         string                 `json:"created_at,omitempty"`
	CreatedBy         string                 `json:"created_by,omitempty"`
	LastUpdatedAt     string                 `json:"last_updated_at,omitempty"`
	LastUpdatedBy     string                 `json:"last_updated_by,omitempty"`
	RetryCount        int                    `json:"retryCount"`
	DeliveryStatus    string                 `json:"deliveryStatus"`
	Diagnostics       map[string]interface{} `json:"diagnostics,omitempty"`
	DiagnosticMessage string                 `json:"diagnosticMessage,omitempty"`
	ReferenceData     map[string]string      `json:"referenceData"` // this is reference object ID we can pass, so that later we can use this
}

func getNotificationMessage(serialisedData datatypes.JSON) *NotificationMessage {
	notificationMessage := &NotificationMessage{}
	json.Unmarshal(serialisedData, notificationMessage)
	return notificationMessage
}

type Notification struct {
	MachineImage   string `json:"machineImage"`
	Header         string `json:"header"`
	Content        string `json:"content"`
	SequenceId     int    `json:"sequenceId"`
	DowntimeJobId  int    `json:"downtimeJobId"`
	IsView         bool   `json:"isView"`
	Id             int    `json:"id"`
	CreatedAt      string `json:"createdAt"`
	ReferenceTime  string `json:"referenceTime"`
	IsJobCompleted bool   `json:"isJobCompleted"`
}

type NotificationHistory struct {
	TotalNotification int            `json:"totalNotification"`
	ViewCount         int            `json:"viewCount"`
	Notification      []Notification `json:"notification"`
}

func (v *Actions) HandleGetNotificationHistory(ctx *gin.Context) {

	userId := common.GetUserId(ctx)
	v.Logger.Info("HandleGetNotificationHistory", zap.Any("userId", userId))
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	// first get all the assigned notification ID, and view notification iD
	var listOfNotifications = authService.GetPushNotificationList(userId)
	var listOfViewedNotifications = authService.GetViewNotificationList(userId)

	var notificationHistory NotificationHistory
	notificationHistory.TotalNotification = len(listOfNotifications)
	notificationHistory.ViewCount = len(listOfViewedNotifications)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	var arrayOfNotificationIds = common.ExtractNotificationIds(listOfNotifications)
	err, listOfNotificationObjects := notificationService.GetNotificationHistory(arrayOfNotificationIds)
	if err != nil {
		v.Logger.Error("error getting notification history", zap.Any("err", err))
		emptyNotificationHistory := NotificationHistory{
			TotalNotification: 0,
			ViewCount:         0,
			Notification:      make([]Notification, 0),
		}
		ctx.JSON(http.StatusOK, emptyNotificationHistory)
		return
	}
	if len(listOfNotificationObjects) == 0 {
		v.Logger.Warn("HandleGetNotificationHistory: no notification history")
		emptyNotificationHistory := NotificationHistory{
			TotalNotification: 0,
			ViewCount:         0,
			Notification:      make([]Notification, 0),
		}
		ctx.JSON(http.StatusOK, emptyNotificationHistory)
		return
	}
	machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

	var notificationList = make([]Notification, 0)
	for index, notification := range listOfNotificationObjects {
		notificationInfo := getNotificationMessage(notification.ObjectInfo)
		// get the job, and get the machine information
		if value, ok := notificationInfo.ReferenceData["downtimeId"]; ok {
			var downtimeId = util.InterfaceToInt(value)
			err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, downtimeId)
			if err == nil {
				downtimeInfo := models.GetMachineDowntimeInfo(c.ObjectInfo)
				err, c2 := machineService.GetAssemblyMachineInfoById(downtimeInfo.MachineId)
				if err == nil {
					machineMasterInfo, _ := dto.GetAssemblyMachineMasterInfo(c2.ObjectInfo)
					newNotification := Notification{}
					newNotification.MachineImage = machineMasterInfo.MachineImage
					newNotification.Header = notificationInfo.Headings["en"]
					newNotification.Content = notificationInfo.Contents["en"]
					newNotification.DowntimeJobId = downtimeId
					// newNotification.IsView = true
					//check viewed notificatins and accordin to that set the values
					if ContainsNotificationId(listOfViewedNotifications, notification.Id) {
						newNotification.IsView = true
					} else {
						newNotification.IsView = false
					}
					newNotification.SequenceId = index + 1
					newNotification.CreatedAt = notificationInfo.CreatedAt
					notificationList = append(notificationList, newNotification)
					newNotification.Id = notification.Id
					parsedTime, err := time.Parse(time.RFC3339, notificationInfo.CreatedAt)
					if err == nil {
						newNotification.ReferenceTime = GetNotificationTime(parsedTime)
					}
					newNotification.IsJobCompleted = v.hasJobCompleted(downtimeId)

				}
			}
		}
	}
	notificationHistory.Notification = notificationList
	v.Logger.Info("sending notification history", zap.Any("history", notificationHistory))
	ctx.JSON(http.StatusOK, notificationHistory)

}

func (v *Actions) hasJobCompleted(jobId int) bool {
	err, c := orm.Get(v.Database, consts.MachineDownTimeMasterTable, jobId)
	if err != nil {
		return false
	} else {
		return !models.GetMachineDowntimeInfo(c.ObjectInfo).CanCheckOut
	}
}
func GetNotificationTime(notificationTime time.Time) string {
	now := time.Now()
	diff := now.Sub(notificationTime)

	switch {
	case diff.Seconds() < 60:
		return fmt.Sprintf("%d sec ago", int(diff.Seconds()))
	case diff.Minutes() < 60:
		return fmt.Sprintf("%d min ago", int(diff.Minutes()))
	case diff.Hours() < 24:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff.Hours() < 24*30:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	case diff.Hours() < 24*365:
		months := int(diff.Hours() / (24 * 30))
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / (24 * 365))
		return fmt.Sprintf("%d years ago", years)
	}
}

func ContainsNotificationId(listOfViewedNotifications []common.NotificationMetaInfo, notificationId int) bool {
	for _, notification := range listOfViewedNotifications {
		if notification.Id == notificationId {
			return true
		}
	}
	return false

}
