package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"gorm.io/gorm"
)

func (v *MaintenanceService) isEmailTemplateExist(dbConnection *gorm.DB, templateType int) bool {
	// this object will contains all the details
	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateType) + " AND object_info ->> '$.isTemplateEnabled' = 'true'"
	listOfObjects, err := GetConditionalObjects(dbConnection, MaintenanceEmailTemplateTable, condition)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

func (v *MaintenanceService) emailGenerator(dbConnection *gorm.DB, templateId int, user int, primaryComponent string, primaryObjectId int) error {

	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateId) + " AND object_info ->>'$.isTemplateEnabled' = 'true'"
	listOfObjects, err := GetConditionalObjects(dbConnection, MaintenanceEmailTemplateTable, condition)
	fmt.Println("MErr", err)
	if err == nil {
		if listOfObjects == nil {
			return errors.New("no email template defined in the system to proceed this action")
		}
		if len(*listOfObjects) == 0 {
			return errors.New("no email template defined in the system to proceed this action")
		}
		templateFieldId := (*listOfObjects)[0].Id
		var targetEmailId string
		emailTemplateInfo := common.GetEmailTemplateInfo((*listOfObjects)[0].ObjectInfo)
		emailTemplate := emailTemplateInfo.Template
		err, commonObject := Get(dbConnection, MaintenanceEmailTemplateFieldTable, templateFieldId)

		targetTable := v.ComponentManager.GetTargetTable(primaryComponent)
		_, primaryObject := Get(dbConnection, targetTable, primaryObjectId)
		var primaryObjectFields = make(map[string]interface{})
		json.Unmarshal(primaryObject.ObjectInfo, &primaryObjectFields)

		if err == nil {
			serviceFields := common.GetEmailTemplateFields(commonObject.ObjectInfo).ServiceFields
			if len(serviceFields) > 0 {
				for _, serviceField := range serviceFields {
					if serviceField.Name == "general_auth" {
						authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
						userObject := authService.GetUserInfoById(user)
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
				moduleFields := common.GetEmailTemplateFields(commonObject.ObjectInfo).ModuleFields

				if len(moduleFields) > 0 {
					for _, moduleField := range moduleFields {
						fmt.Println("===========================================================================")
						fmt.Println("moduleField.ComponentName: ", moduleField.ComponentName, " primaryComponent :", primaryComponent)
						index := strings.Index(emailTemplateInfo.Template, moduleField.Display)
						fmt.Println("index", index)
						fmt.Println("moduleField.Display: ", moduleField.Display)

						if index != -1 {
							// yes we have that field

							if moduleField.ComponentName == primaryComponent {
								// this is the primary component
								if moduleField.Type == "href" {
									fmt.Println("href")
									var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa(primaryObjectId)
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
								} else if moduleField.Type == "time" { // this is not checking the component, but if href is configured , then add the field
									var objectLink = util.GetZoneCurrentTimeInPMFormat("Asia/Singapore")

									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
								} else {

									replacementValue := util.InterfaceToString(primaryObjectFields[moduleField.Property])
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
								}
							} else if moduleField.Type == "href" { // this is not checking the component, but if href is configured , then add the field
								var objectLink = v.EmailNotificationDomain

								emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
							} else if moduleField.IsForeignObjectMap {
								targetTable = moduleField.ComponentName
								recordId := util.InterfaceToInt(primaryObjectFields[moduleField.ForeignLinkedField])
								var object component.GeneralObject
								var err error

								if primaryComponent == MouldMaintenanceCorrectiveWorkOrderTable {
									machineInterface := common.GetService("moulds_module").ServiceInterface.(common.MouldInterface)
									err, object = machineInterface.GetMouldInfoById(ProjectID, recordId)

								} else {
									machineInterface := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
									err, object = machineInterface.GetMachineInfoById(ProjectID, recordId)

								}
								if err == nil {
									var objectFields = make(map[string]interface{})
									json.Unmarshal(object.ObjectInfo, &objectFields)
									if moduleField.Type == "href" {
										var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa(object.Id)
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
									} else {
										replacementValue := util.InterfaceToString(objectFields[moduleField.Property])

										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
									}
								}

							} else {

								// some other component, using the primary id in the object, it is possible objects are using primary key or internal objects
								if moduleField.IsObjectField {
									// this is the object field
									targetTable = v.ComponentManager.GetTargetTable(moduleField.ComponentName)
									condition = " object_info ->>'$." + moduleField.LinkedField + "'= " + strconv.Itoa(primaryObjectId)
									listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, condition)
									if err == nil {
										// needed to check whether we have listOfObject is null or not
										if listOfObjects != nil {
											if len(*listOfObjects) > 0 {
												var objectFields = make(map[string]interface{})
												json.Unmarshal((*listOfObjects)[0].ObjectInfo, &objectFields)
												if moduleField.Type == "href" {
													var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa((*listOfObjects)[0].Id)
													emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
												} else {
													replacementValue := util.InterfaceToString(objectFields[moduleField.Property])
													emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
												}
											} else {
												fmt.Println("no objects found, query condition", condition, " target Table :", targetTable)
											}
										} else {
											fmt.Println("list of object is null query condition", condition, " target Table :", targetTable)
										}

									}
								} else {
									// we can get the fields by accessing the primary key
									targetTable = v.ComponentManager.GetTargetTable(moduleField.ComponentName)
									err, commonObject = Get(dbConnection, targetTable, primaryObjectId)
									if err == nil {
										var objectFields = make(map[string]interface{})
										json.Unmarshal(commonObject.ObjectInfo, &objectFields)
										if moduleField.Type == "href" {
											var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa((*listOfObjects)[0].Id)
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
			v.BaseService.Logger.Error("error getting template field", zap.String("error", err.Error()))
			return err
		}
		emailTemplate = util.FormatStringHTML(emailTemplate)
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
		fmt.Println("sending email teamplate to :", emailList)
		emailMessages = append(emailMessages, emailMessage)
		//fmt.Println("emailMessage: ", emailMessage)
		notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
		notificationService.CreateMessages(ProjectID, emailMessages)
	} else {
		fmt.Println("error getting email template :", err.Error())
		v.BaseService.Logger.Error("error getting email template", zap.String("error", err.Error()))
	}
	return err

}

func (v *MaintenanceService) isWorkflowEmailExist(dbConnection *gorm.DB, tableName string, budgetWorkflowId int) []int {
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

func (v *MaintenanceService) notifyWorkflowEmail(dbConnection *gorm.DB, userList []int, componentName string, templateId int, id int) {

	for _, user := range userList {
		v.emailGenerator(dbConnection, templateId, user, componentName, id)
	}

}
