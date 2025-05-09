package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (as *AuthService) isEmailTemplateExist(dbConnection *gorm.DB, templateType int) bool {
	// this object will contains all the details
	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateType) + " AND object_info ->> '$.isTemplateEnabled' = 'true'"
	err, listOfObjects := GetConditionalObjects(dbConnection, UserEmailTemplateTable, condition)
	if err == nil {
		if len(listOfObjects) == 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

func (as *AuthService) emailGenerator(dbConnection *gorm.DB, templateId int, user int, primaryComponent string, primaryObjectId int) error {

	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateId) + " AND object_info ->>'$.isTemplateEnabled' = 'true'"
	err, listOfObjects := GetConditionalObjects(dbConnection, UserEmailTemplateTable, condition)

	if err == nil {
		var targetEmailId string
		emailTemplateInfo := common.GetEmailTemplateInfo((listOfObjects)[0].ObjectInfo)
		emailTemplate := emailTemplateInfo.Template
		err, commonObject := Get(dbConnection, UserEmailTemplateFieldsTable, emailTemplateInfo.TemplateType)

		targetTable := as.ComponentManager.GetTargetTable(primaryComponent)
		_, primaryObject := Get(dbConnection, targetTable, primaryObjectId)
		var primaryObjectFields = make(map[string]interface{})
		json.Unmarshal(primaryObject.ObjectInfo, &primaryObjectFields)

		if err == nil {
			serviceFields := common.GetEmailTemplateFields(commonObject.ObjectInfo).ServiceFields
			if len(serviceFields) > 0 {
				for _, serviceField := range serviceFields {
					if serviceField.Name == "general_auth" {
						userObject := as.GetUserInfoById(user)
						targetEmailId = userObject.Email
						for _, internalServiceField := range serviceField.Fields {
							index := strings.Index(emailTemplateInfo.Template, internalServiceField.Display)
							if index != -1 {
								serialisedObject, _ := json.Marshal(userObject)
								var internalServiceObjectFields = make(map[string]interface{}, 0)
								json.Unmarshal(serialisedObject, &internalServiceObjectFields)
								replacementValue := util.InterfaceToString(internalServiceObjectFields[internalServiceField.Property])
								emailTemplate = strings.Replace(emailTemplate, internalServiceField.Display, replacementValue, -1)
							}
						}

					}
				}
			}
			moduleFields := common.GetEmailTemplateFields(commonObject.ObjectInfo).ModuleFields
			if len(moduleFields) > 0 {

				if len(moduleFields) > 0 {
					for _, moduleField := range moduleFields {

						index := strings.Index(emailTemplateInfo.Template, moduleField.Display)
						if index != -1 {
							// yes we have that field
							if moduleField.ComponentName == primaryComponent {
								// this is the primary component
								if moduleField.Type == "href" {
									var objectLink string
									if moduleField.MenuId != "" {
										objectLink = as.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa(primaryObjectId)
									} else {
										objectLink = as.EmailNotificationDomain
									}
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
								} else {

									if value, ok := primaryObjectFields[moduleField.Property]; ok {
										replacementValue := util.InterfaceToString(value)
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
										if moduleField.IsUsedForTargetEmail {
											targetEmailId = replacementValue
										}
									} else {
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, "[Invalid Input]", -1)
										if moduleField.IsUsedForTargetEmail {
											targetEmailId = ""
										}
									}

								}
							} else {
								// some other component, using the primary id in the object, it is possible objects are using primary key or internal objects
								if moduleField.IsObjectField {
									// this is the object field
									targetTable = as.ComponentManager.GetTargetTable(moduleField.ComponentName)
									condition = " object_info ->>'$." + moduleField.LinkedField + "'= " + strconv.Itoa(primaryObjectId)
									err, listOfObjects = GetConditionalObjects(dbConnection, targetTable, condition)
									if err == nil {
										var objectFields = make(map[string]interface{})
										json.Unmarshal((listOfObjects)[0].ObjectInfo, &objectFields)
										if moduleField.Type == "href" {
											var objectLink = as.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa((listOfObjects)[0].Id)
											emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
										} else {
											replacementValue := util.InterfaceToString(objectFields[moduleField.Property])
											emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
										}

										if moduleField.IsUsedForTargetEmail {
											targetEmailId = util.InterfaceToString(objectFields[moduleField.Property])
										}
									}
								} else {
									// we can get the fields by accessing the primary key
									targetTable = as.ComponentManager.GetTargetTable(moduleField.ComponentName)
									err, commonObject = Get(dbConnection, targetTable, primaryObjectId)
									if err == nil {
										var objectFields = make(map[string]interface{})
										json.Unmarshal(commonObject.ObjectInfo, &objectFields)
										if moduleField.Type == "href" {
											var objectLink = as.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa((listOfObjects)[0].Id)
											emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
										} else {
											replacementValue := util.InterfaceToString(objectFields[moduleField.Property])
											emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
										}
									}
								}
							}

						}

					}
				}

			}
		} else {
			fmt.Println("error getting :", err.Error())
			as.BaseService.Logger.Error("error getting template field", zap.String("error", err.Error()))
			return err
		}

		if targetEmailId == "" {
			return errors.New("no target email found")
		}
		emailTemplate = util.FormatStringHTML(emailTemplate)

		as.BaseService.Logger.Info("generated email template", zap.String("template", emailTemplate))

		var emailList []string
		emailList = append(emailList, targetEmailId)
		emailMessages := make([]common.Message, 0)
		emailMessage := common.Message{
			To:          emailList,
			SingleEmail: false,
			Subject:     emailTemplateInfo.Subject,
			Body: map[string]string{
				"text/html": emailTemplate,
			},
			Info:          "",
			ReplyTo:       make([]string, 0),
			EmbeddedFiles: nil,
			AttachedFiles: nil,
		}

		emailMessages = append(emailMessages, emailMessage)
		notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
		notificationService.CreateMessages(ProjectID, emailMessages)
	} else {
		fmt.Println("error getting email template :", err.Error())
		as.BaseService.Logger.Error("error getting email template", zap.String("error", err.Error()))
	}
	return err

}

func (as *AuthService) isWorkflowEmailExist(dbConnection *gorm.DB, tableName string, budgetWorkflowId int) []int {
	err, budgetWorkflowStatus := Get(dbConnection, tableName, budgetWorkflowId)
	if err == nil {
		var workflowFields = make(map[string]interface{})
		json.Unmarshal(budgetWorkflowStatus.ObjectInfo, &workflowFields)
		if value, ok := workflowFields["userList"]; ok {
			return util.InterfaceToIntArray(value)
		}
	}
	return []int{}
}

func (as *AuthService) notifyWorkflowEmail(dbConnection *gorm.DB, userList []int, componentName string, templateId int, id int) {

	for _, user := range userList {
		as.emailGenerator(dbConnection, templateId, user, componentName, id)
	}

}
