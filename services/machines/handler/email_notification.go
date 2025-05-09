package handler

import (
	"bufio"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"fmt"
	"os"
	"strings"
)

func (v *MachineService) emailGenerator(productionOrderName string, machineName string, targetEmailId string, ccEmailIds []string) error {

	file, err := os.Open(v.EmailConfig.FilePath)

	if err != nil {
		fmt.Println("unable to open the template file..")
		v.BaseService.Logger.Error("unable to open the template file..should be resources folder under notification service..")
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the file contents into a string
	var templateString string
	for scanner.Scan() {
		templateString += scanner.Text() + "\n"
	}

	emailTemplate := strings.Replace(templateString, "{{PRODUCTION_ORDER}}", productionOrderName, -1)
	emailTemplate = strings.Replace(emailTemplate, "{{MACHINE_NAME}}", machineName, -1)

	emailTemplate = util.FormatStringHTML(emailTemplate)
	var emailList []string
	emailList = append(emailList, targetEmailId)
	emailMessages := make([]common.Message, 0)
	emailMessage := common.Message{
		To:          emailList,
		SingleEmail: false,
		Subject:     "HMI Warning Message",
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

	return err

}
