package jobs

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/services/spare_part_inventory/source/consts"
	"cx-micro-flake/services/spare_part_inventory/source/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (v *Jobs) GenerateNotificationMessage() {
	var duration = time.Duration(v.JobConfig.InventoryLimitPollingInterval) * time.Second
	var sparePartInventoryMasterTable = v.ComponentManager.GetTargetTable(consts.SparePartInventoryMasterComponent)
	// var sparePartInventoryEmailEscalationTable = v.ComponentManager.GetTargetTable(consts.SparePartInventoryEmailEscalationComponent)
	var _ = v.ComponentManager.GetTargetTable(consts.SparePartInventoryEmailEscalationComponent)
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

	for {
		time.Sleep(duration)
		v.Logger.Info("Checking inventory master table to generate the onHand quantity alert")
		var condition = " object_info->>'$.onHandQty' < " + strconv.Itoa(v.settingInfo.InventoryThresholdQuantity)

		err, listOfThreshold := v.Repository.GetResourceWithCondition(sparePartInventoryMasterTable, condition)
		if err != nil {
			v.Logger.Error("error getting records", zap.String("error", err.Error()))
			continue
		}

		if len(listOfThreshold) > 0 {
			for _, thresholdDetailInterface := range listOfThreshold {
				err, inventoryMasterInfo := models.GetSparePartInventoryMasterInfo(thresholdDetailInterface.ObjectInfo)
				if err != nil {
					v.Logger.Error("error getting spare part inventory master info", zap.String("error", err.Error()))
					continue
				}
				condition := "object_info->>'$.sparePartId' = " + strconv.Itoa(thresholdDetailInterface.Id)

				var targetUsers = v.settingInfo.InventoryStockLimitAlertUsers
				count := v.Repository.GetCountByCondition("spare_part_inventory_email_escalation", condition)
				if count == 0 {
					for _, userId := range targetUsers {
						userInfo := authService.GetUserInfoById(userId)
						var emailTemplate = v.EscalationEmailTemplate
						emailTemplate = strings.Replace(emailTemplate, "[USER]", userInfo.Username, 1)
						emailTemplate = strings.Replace(emailTemplate, "[DESCRIPTION]", inventoryMasterInfo.SparePartNumber, 1)
						v.sendEmail(userInfo.Email, emailTemplate)
					}
					emailEscalationInfo := models.SparePartInventoryEmailEscalationInfo{
						SparePartId:     thresholdDetailInterface.Id,
						EmailRecipients: targetUsers,
						ObjectStatus:    common.ObjectStatusActive,
					}
					var serialisedData = emailEscalationInfo.Serialised()
					err, _ = v.Repository.CreateResource("spare_part_inventory_email_escalation", serialisedData, 1)
					if err == nil {
						fmt.Println("Action service is ok", err)
						v.Logger.Info("spare part job is successfully created")
					} else {
						v.Logger.Error("spare part job creation failed", zap.Error(err))
					}
				} else {
					continue
				}
			}
		}

	}
}

func (v *Jobs) sendEmail(targetEmail string, emailContent string) {
	var emailMessages []common.Message

	emailMessage := common.Message{
		To:          []string{targetEmail},
		SingleEmail: false,
		Subject:     "Spare Part Inventory Escalation",
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
