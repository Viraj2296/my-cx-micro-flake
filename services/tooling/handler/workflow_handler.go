package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
)

func getTimeDifference(dst string) string {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch
	if difference < 60 {
		// this is seconds
		return strconv.Itoa(int(difference)) + "  seconds"
	} else if difference < 3600 {
		minutes := difference / 60
		return strconv.Itoa(int(minutes)) + "  minutes"
	} else {
		minutes := difference / 3600
		return strconv.Itoa(int(minutes)) + "  hour"
	}
}

func isLessThan24Hours(dst string) bool {
	currentTime := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	var difference = util.ConvertStringToDateTime(currentTime).DateTimeEpoch - util.ConvertStringToDateTime(dst).DateTimeEpoch

	hours := difference / 3600
	if hours < 24 {
		return true
	}
	return false
}

func isTaskAssignNotificationEmailConfigured(dbConnection *gorm.DB) bool {
	condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.templateType\")) = " + strconv.Itoa(ToolingProjectTaskAssignmentEmailTemplateType) + " AND JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isTemplateEnabled\")) = 'true'"
	listOfObjects, _ := GetConditionalObjects(dbConnection, ToolingEmailTemplateTable, condition)
	if len(*listOfObjects) == 1 {
		return true
	}
	return false
}

// this should be generated based on component name, now, we will only consider for main component
func (v *ToolingService) generateEmailObject(dbConnection *gorm.DB, componentName string, userId int, id int) map[string]interface{} {
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userObject := authService.GetUserInfoById(userId)
	var emailObjectFields = make(map[string]interface{})
	if componentName == ToolingProjectTaskComponent {
		err, toolingProjectTaskInterface := Get(dbConnection, ToolingProjectTaskTable, id)
		if err == nil {
			json.Unmarshal(toolingProjectTaskInterface.ObjectInfo, &emailObjectFields)
			routingLink := v.EmailNotificationDomain + "?routeLink=tooling_task&recordId=" + strconv.Itoa(id)
			emailObjectFields["taskLink"] = routingLink
		}
	}
	serialisedUserFields, _ := json.Marshal(userObject)
	json.Unmarshal(serialisedUserFields, &emailObjectFields)

	return emailObjectFields
}

func (v *ToolingService) notifyProjectTaskAssignment(dbConnection *gorm.DB, userId, id int) {
	// this object will contains all the details
	//emailObjectFields := ts.generateEmailObject(dbConnection, ToolingProjectTaskComponent, userId, id)
	//condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.templateType\")) = " + strconv.Itoa(ToolingProjectTaskAssignmentEmailTemplateType) + " AND JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isTemplateEnabled\")) = 'true'"
	//listOfObjects, _ := GetConditionalObjects(dbConnection, ToolingEmailTemplateTable, condition)
	//
	//if len(*listOfObjects) == 1 {
	//	emailTemplateInfo := common.GetEmailTemplateInfo((*listOfObjects)[0].ObjectInfo)
	//	emailTemplate := emailTemplateInfo.Template
	//	listOfEmailTemplateFields, _ := GetObjects(dbConnection, ToolingEmailTemplateFieldTable)
	//	templateFields := common.GetEmailTemplateFields((*listOfEmailTemplateFields)[0].ObjectInfo)
	//	fmt.Println("templateFields: templateFields:", templateFields)
	//	for _, templateField := range templateFields {
	//		if value, ok := emailObjectFields[templateField.Property]; ok {
	//			fmt.Println("templateField.Property:", templateField.Property)
	//			emailTemplate = strings.Replace(emailTemplate, templateField.Display, value.(string), -1)
	//		}
	//	}
	//
	//	emailTemplate = util.FormatStringHTML(emailTemplate)
	//	ts.BaseService.Logger.Infow("generated email template", "template", emailTemplate)
	//	var emailList []string
	//	emailList = append(emailList, emailObjectFields["email"].(string))
	//	emailMessages := make([]common.Message, 0)
	//	emailMessage := common.Message{
	//		To:          emailList,
	//		SingleEmail: false,
	//		Subject:     emailTemplateInfo.Subject,
	//		Body: map[string]string{
	//			"text/html": emailTemplate,
	//		},
	//		Info:          "",
	//		ReplyTo:       make([]string, 0),
	//		EmbeddedFiles: nil,
	//		AttachedFiles: nil,
	//	}
	//
	//	emailMessages = append(emailMessages, emailMessage)
	//	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	//	notificationService.CreateMessages(ProjectID, emailMessages)
	//}
}

func (v *ToolingService) notifyTaskCompletion(dbConnection *gorm.DB, userId, id int) {
	// this object will contain all the details
	//emailObjectFields := ts.generateEmailObject(dbConnection, ToolingProjectTaskComponent, userId, id)
	//condition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.templateType\")) = " + strconv.Itoa(ToolingProjectTaskStatusEmailTemplateType) + " AND JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isTemplateEnabled\")) = 'true'"
	//listOfObjects, _ := GetConditionalObjects(dbConnection, ToolingEmailTemplateTable, condition)
	//
	//if len(*listOfObjects) == 1 {
	//	emailTemplateInfo := common.GetEmailTemplateInfo((*listOfObjects)[0].ObjectInfo)
	//	emailTemplate := emailTemplateInfo.Template
	//	listOfEmailTemplateFields, _ := GetObjects(dbConnection, ToolingEmailTemplateFieldTable)
	//	templateFields := common.GetEmailTemplateFields((*listOfEmailTemplateFields)[0].ObjectInfo)
	//	fmt.Println("templateFields: templateFields:", templateFields)
	//	for _, templateField := range templateFields {
	//		if value, ok := emailObjectFields[templateField.Property]; ok {
	//			fmt.Println("templateField.Property:", templateField.Property)
	//			emailTemplate = strings.Replace(emailTemplate, templateField.Display, value.(string), -1)
	//		}
	//	}
	//
	//	emailTemplate = util.FormatStringHTML(emailTemplate)
	//	ts.BaseService.Logger.Infow("generated email template", "template", emailTemplate)
	//	var emailList []string
	//	emailList = append(emailList, emailObjectFields["email"].(string))
	//	emailMessages := make([]common.Message, 0)
	//	emailMessage := common.Message{
	//		To:          emailList,
	//		SingleEmail: false,
	//		Subject:     emailTemplateInfo.Subject,
	//		Body: map[string]string{
	//			"text/html": emailTemplate,
	//		},
	//		Info:          "",
	//		ReplyTo:       make([]string, 0),
	//		EmbeddedFiles: nil,
	//		AttachedFiles: nil,
	//	}
	//
	//	emailMessages = append(emailMessages, emailMessage)
	//	notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
	//	notificationService.CreateMessages(ProjectID, emailMessages)
	//}
}
