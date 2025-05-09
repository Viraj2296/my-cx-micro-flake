package notification

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/facility/handler/const_util"
	"cx-micro-flake/services/facility/handler/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

type EmailHandler struct {
	Logger                  *zap.Logger
	EmailNotificationDomain string
	ComponentManager        *common.ComponentManager
}

func (v *EmailHandler) IsEmailTemplateExist(dbConnection *gorm.DB, templateType int) bool {
	// this object will contains all the details
	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateType) + " AND object_info ->> '$.isTemplateEnabled' = 'true'"
	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceEmailTemplateTable, condition)
	if err == nil {
		if len(*listOfObjects) == 0 {
			return false
		} else {
			return true
		}
	}
	return false
}

func (v *EmailHandler) EmailGeneratorByEmail(dbConnection *gorm.DB, templateId int, targetEmailId string, primaryComponent string, primaryObjectId int) error {

	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateId) + " AND object_info ->>'$.isTemplateEnabled' = 'true'"
	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceEmailTemplateTable, condition)

	if err == nil {
		if listOfObjects == nil {
			return errors.New("no email template defined in the system to proceed this action")
		}
		if len(*listOfObjects) == 0 {
			return errors.New("no email template defined in the system to proceed this action")
		}
		templateFieldId := (*listOfObjects)[0].Id

		emailTemplateInfo := common.GetEmailTemplateInfo((*listOfObjects)[0].ObjectInfo)
		emailTemplate := emailTemplateInfo.Template
		err, commonObject := database.Get(dbConnection, const_util.FacilityServiceEmailTemplateFieldTable, templateFieldId)

		targetTable := v.ComponentManager.GetTargetTable(primaryComponent)
		_, primaryObject := database.Get(dbConnection, targetTable, primaryObjectId)
		var primaryObjectFields = make(map[string]interface{})
		json.Unmarshal(primaryObject.ObjectInfo, &primaryObjectFields)

		if err == nil {
			serviceFields := common.GetEmailTemplateFields(commonObject.ObjectInfo).ServiceFields
			if len(serviceFields) > 0 {
				for _, serviceField := range serviceFields {
					if serviceField.Name == "general_auth" {
						authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
						userId := authService.EmailToUserId(targetEmailId)
						userObject := authService.GetUserInfoById(userId)
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
						fmt.Println("moduleField.ComponentName: ", moduleField.ComponentName, " primaryComponent :", primaryComponent)
						index := strings.Index(emailTemplateInfo.Template, moduleField.Display)
						fmt.Println("index", index)
						fmt.Println("moduleField.Display: ", moduleField.Display)

						if index != -1 {
							// yes we have that field

							if moduleField.ComponentName == primaryComponent {
								// this is the primary component
								if moduleField.Type == "href" {
									var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa(primaryObjectId)
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
								} else {
									replacementValue := util.InterfaceToString(primaryObjectFields[moduleField.Property])
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
								}
							} else if moduleField.Type == "href" { // this is not checking the component, but if href is configured , then add the field

								if moduleField.TokenAttributes != nil {
									// we have token attributes defined, if the generate the token
									claims := jwt.MapClaims{}
									claims["fields"] = moduleField.TokenAttributes.Fields
									claims["menuId"] = moduleField.TokenAttributes.MenuId
									claims["inAppRouting"] = moduleField.TokenAttributes.InAppRouting
									claims["recordId"] = primaryObjectId
									token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
									generatedToken, err := token.SignedString([]byte(os.Getenv("api_secret")))
									if err != nil {
										// This is for safe side
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, v.EmailNotificationDomain, -1)
									} else {
										var objectLink = v.EmailNotificationDomain + "?routeLink=email_link_validation&token=" + generatedToken
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
									}
								}
								var objectLink = v.EmailNotificationDomain
								emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
							} else {
								// some other component, using the primary id in the object, it is possible objects are using primary key or internal objects
								if moduleField.IsObjectField {
									// this is the object field
									targetTable = v.ComponentManager.GetTargetTable(moduleField.ComponentName)
									condition = " object_info ->>'$." + moduleField.LinkedField + "'= " + strconv.Itoa(primaryObjectId)
									listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, condition)
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
									err, commonObject = database.Get(dbConnection, targetTable, primaryObjectId)
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
			v.Logger.Error("error getting template field", zap.String("error", err.Error()))
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
		notificationService.CreateMessages(const_util.ProjectID, emailMessages)
	} else {
		fmt.Println("error getting email template :", err.Error())
		v.Logger.Error("error getting email template", zap.String("error", err.Error()))
	}
	return err

}

func (v *EmailHandler) EmailGenerator(dbConnection *gorm.DB, templateId int, user int, primaryComponent string, primaryObjectId int) error {

	condition := " object_info ->>'$.templateType'= " + strconv.Itoa(templateId) + " AND object_info ->>'$.isTemplateEnabled' = 'true'"
	listOfObjects, err := database.GetConditionalObjects(dbConnection, const_util.FacilityServiceEmailTemplateTable, condition)
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
		fmt.Println("emailTemplate", emailTemplate)
		err, commonObject := database.Get(dbConnection, const_util.FacilityServiceEmailTemplateFieldTable, templateFieldId)
		fmt.Println("err", err)
		targetTable := v.ComponentManager.GetTargetTable(primaryComponent)
		_, primaryObject := database.Get(dbConnection, targetTable, primaryObjectId)
		var primaryObjectFields = make(map[string]interface{})
		json.Unmarshal(primaryObject.ObjectInfo, &primaryObjectFields)
		fmt.Println("err2", err)
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
						fmt.Println("moduleField.ComponentName: ", moduleField.ComponentName, " primaryComponent :", primaryComponent)
						index := strings.Index(emailTemplateInfo.Template, moduleField.Display)
						fmt.Println("index", index)
						fmt.Println("moduleField.Display: ", moduleField.Display)

						if index != -1 {
							// yes we have that field

							if moduleField.ComponentName == primaryComponent {
								// this is the primary component
								if moduleField.Type == "href" {
									var objectLink = v.EmailNotificationDomain + "?routeLink=" + moduleField.MenuId + "&recordId=" + strconv.Itoa(primaryObjectId)
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
								} else {
									replacementValue := util.InterfaceToString(primaryObjectFields[moduleField.Property])
									emailTemplate = strings.Replace(emailTemplate, moduleField.Display, replacementValue, -1)
								}
							} else if moduleField.Type == "href" { // this is not checking the component, but if href is configured , then add the field

								if moduleField.TokenAttributes != nil {
									// we have token attributes defined, if the generate the token
									claims := jwt.MapClaims{}
									claims["fields"] = moduleField.TokenAttributes.Fields
									claims["menuId"] = moduleField.TokenAttributes.MenuId
									claims["inAppRouting"] = moduleField.TokenAttributes.InAppRouting
									claims["recordId"] = primaryObjectId
									token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
									generatedToken, err := token.SignedString([]byte(os.Getenv("api_secret")))
									if err != nil {
										// This is for safe side
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, v.EmailNotificationDomain, -1)
									} else {
										var objectLink = v.EmailNotificationDomain + "?routeLink=email_link_validation&token=" + generatedToken
										emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
									}
								}
								var objectLink = v.EmailNotificationDomain
								emailTemplate = strings.Replace(emailTemplate, moduleField.Display, objectLink, -1)
							} else {
								// some other component, using the primary id in the object, it is possible objects are using primary key or internal objects
								if moduleField.IsObjectField {
									// this is the object field
									targetTable = v.ComponentManager.GetTargetTable(moduleField.ComponentName)
									condition = " object_info ->>'$." + moduleField.LinkedField + "'= " + strconv.Itoa(primaryObjectId)
									listOfObjects, err = database.GetConditionalObjects(dbConnection, targetTable, condition)
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
									err, commonObject = database.Get(dbConnection, targetTable, primaryObjectId)
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
			v.Logger.Error("error getting template field", zap.String("error", err.Error()))
			return err
		}
		emailTemplate = util.FormatStringHTML(emailTemplate)
		v.Logger.Info("generated email template", zap.Any("template", emailTemplate))
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
		notificationService.CreateMessages(const_util.ProjectID, emailMessages)
	} else {
		fmt.Println("error getting email template :", err.Error())
		v.Logger.Error("error getting email template", zap.String("error", err.Error()))
	}
	return err

}

func (v *EmailHandler) isWorkflowEmailExist(dbConnection *gorm.DB, tableName string, budgetWorkflowId int) []int {
	err, budgetWorkflowStatus := database.Get(dbConnection, tableName, budgetWorkflowId)
	if err == nil {
		var workflowFields = make(map[string]interface{})
		json.Unmarshal(budgetWorkflowStatus.ObjectInfo, &workflowFields)
		if value, ok := workflowFields["userList"]; ok {
			return util.InterfaceToIntArray(value)
		}
	}
	return []int{}
}

func (v *EmailHandler) notifyWorkflowEmail(dbConnection *gorm.DB, userList []int, componentName string, templateId int, id int) {
	for _, user := range userList {
		v.EmailGenerator(dbConnection, templateId, user, componentName, id)
	}
}
