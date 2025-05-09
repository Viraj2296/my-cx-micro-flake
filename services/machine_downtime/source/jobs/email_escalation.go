package jobs

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/machine_downtime/source/consts"
	"cx-micro-flake/services/machine_downtime/source/models"
	"fmt"
	"go.cerex.io/transcendflow/component"
	"go.cerex.io/transcendflow/orm"
	"go.cerex.io/transcendflow/util"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CheckIntervalPassed checks if the interval in minutes has passed since the fault creation time.
func CheckIntervalPassed(faultCreatedTime string, thresholdMinutes int) (bool, float64, error) {
	parsedTime, err := time.Parse(time.RFC3339, faultCreatedTime)
	if err != nil {
		return false, 0, fmt.Errorf("invalid time format: %w", err)
	}
	currentTime := time.Now()
	timeDiff := currentTime.Sub(parsedTime).Minutes()
	return timeDiff > float64(thresholdMinutes), timeDiff, nil
}
func (v *Jobs) SendEmailEscalationJob() {
	var condition = " object_info->>'$.canCheckIn' = 'true' and object_info ->> '$.objectStatus' = 'Active'"
	err, listOfFaultsInterface := orm.GetConditionalObjects(v.Database, consts.MachineDownTimeMasterTable, condition)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	var duration = time.Duration(v.DowntimeConfig.JobServiceConfig.EscalationPollingInterval) * time.Second
	if err == nil {
		for {
			for _, downtimeObject := range listOfFaultsInterface {
				downtimeJobInfo := models.GetMachineDowntimeInfo(downtimeObject.ObjectInfo)
				v.Logger.Info("processing downtime job", zap.Any("job", downtimeJobInfo))
				var downTimeDataCount = v.Repository.GetEmailEscalationCount(downtimeObject.Id)

				var isPassed bool
				var timeDifference float64
				if downTimeDataCount > 0 {
					escalationInfo := v.Repository.GetEmailEscalationInfo(downtimeObject.Id)
					isPassed, timeDifference, err = CheckIntervalPassed(escalationInfo.CreatedAt, v.machineDowntimeSettingInfo.EscalationWaitingPeriod)
				} else {
					// If there is no data for downtime then we need to create new data
					// nothing created, lets insert it
					isPassed, timeDifference, err = CheckIntervalPassed(downtimeJobInfo.CreatedAt, v.machineDowntimeSettingInfo.InitialWaitingPeriod)

				}
				v.Logger.Debug("time difference", zap.Any("diff", timeDifference), zap.Any("is_passed", isPassed))
				if err == nil {
					if isPassed {
						// now send the email to configured users
						var targetUsers = v.machineDowntimeSettingInfo.PrimaryEmailRecipients
						for _, userId := range targetUsers {
							userInfo := authService.GetUserInfoById(userId)
							var emailTemplate = v.EscalationEmailTemplate
							emailTemplate = strings.Replace(emailTemplate, "[USER]", userInfo.Username, 1)
							difference := fmt.Sprintf("%f", timeDifference)
							emailTemplate = strings.Replace(emailTemplate, "[MIN]", difference, 1)
							v.sendEmail(userInfo.Email, emailTemplate)
						}
						emailEscalationInfo := models.MachineDowntimeEmailEscalationInfo{
							DowntimeId:      downtimeObject.Id,
							CreatedAt:       util.GetCurrentTime("2006-01-02T15:04:05.000Z"),
							CreatedBy:       1,
							EmailRecipients: targetUsers,
							ObjectStatus:    common.ObjectStatusActive,
						}
						generalObject := component.GeneralObject{
							ObjectInfo: emailEscalationInfo.Serialised(),
						}
						err, i := orm.CreateFromGeneralObject(v.Database, consts.MachineDownTimeEmailEscalationTable, generalObject)
						if err == nil {
							fmt.Println("Action service is ok", err)
							v.Logger.Info("downtime job is successfully created", zap.Any("job", i))
						} else {
							v.Logger.Error("downtime job creation failed", zap.Error(err))
						}
					} else {
						v.Logger.Info("email triggering interval is not yet passed ...")
					}
				}

			}
			time.Sleep(duration)
		}

	} else {
		v.Logger.Error("error getting jobs")
	}
}

func (v *Jobs) sendEmail(targetEmail string, emailContent string) {
	var emailMessages []common.Message

	emailMessage := common.Message{
		To:          []string{targetEmail},
		SingleEmail: false,
		Subject:     "Assembly Fault Repair Escalation",
		Body: map[string]string{
			"text/html": emailContent,
		},
		Info:          "",
		EmbeddedFiles: nil,
		AttachedFiles: nil,
	}

	emailMessages = append(emailMessages, emailMessage)
	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	err := notificationService.CreateMessages("906d0fd569404c59956503985b330132", emailMessages)
	if err != nil {
		v.Logger.Error("error creating notification messages", zap.Error(err))
	} else {
		v.Logger.Info("notification messages successfully created")
	}

}
