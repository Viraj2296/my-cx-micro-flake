package handler

import (
	"bufio"
	"context"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	notificationTable = "notification"
)

type EmailConfig struct {
	ConnectionProfile string `json:"connectionProfile"`
	SendingInterval   string `json:"sendingInterval"`
	FromEmail         string `json:"fromEmail"`
	PoolingEnabled    bool   `json:"poolingEnabled"`
	TestEmail         struct {
		IsEnabled         bool   `json:"isEnabled"`
		To                string `json:"to"`
		TestEmailTemplate string `json:"testEmailTemplate"`
	} `json:"testEmail"`
}

type PushNotificationConfig struct {
	OneSignalAPI     string `json:"oneSignalAPI"` //https://onesignal.com/api/v1/notifications
	OneSignalAppsAPI string `json:"oneSignalAppsAPI"`
	AppID            string `json:"appID"`
	APIKey           string `json:"apiKey"`
	PoolingEnabled   bool   `json:"poolingEnabled"`
	PoolingInterval  string `json:"poolingInterval"`
}
type NotificationManager struct {
	ServiceDatabase        *gorm.DB
	Client                 *ses.Client
	EmailConfig            EmailConfig
	Logger                 *zap.Logger
	PushNotificationConfig PushNotificationConfig
}

func (v *NotificationManager) Init() error {

	if v.EmailConfig.PoolingEnabled {
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithSharedConfigProfile(v.EmailConfig.ConnectionProfile))
		if err != nil {
			v.Logger.Error("Failed to load configuration", zap.String("error", err.Error()))
			v.Client = nil
			return err
		}
		client := ses.NewFromConfig(cfg)
		v.Client = client
	} else {
		v.Logger.Warn("email pooling is disabled, so setting null")
		v.Client = nil
	}

	return nil
}
func (v *NotificationManager) PoolNotification() {
	tDuration, err := time.ParseDuration(v.EmailConfig.SendingInterval)
	if err != nil {
		tDuration, err = time.ParseDuration("10s")
	}
	for {

		listOfObjects, err := GetObjects(v.ServiceDatabase, notificationTable)
		if err == nil {
			for _, object := range *listOfObjects {
				var emailMessage = make(map[string]interface{})
				json.Unmarshal(object.ObjectInfo, &emailMessage)

				retryCount := util.InterfaceToInt(emailMessage["retryCount"])
				if emailMessage["status"] == "pending" && retryCount < 2 {
					emailObject := emailMessage["message"]
					newEmailMessage := common.Message{}
					rawEmailObject, _ := json.Marshal(emailObject)
					json.Unmarshal(rawEmailObject, &newEmailMessage)
					newEmailMessage.From = v.EmailConfig.FromEmail
					sendingError := v.sendEmail(newEmailMessage)

					updatingData := make(map[string]interface{})
					if sendingError == nil {
						emailMessage["status"] = "sent"
						emailMessage["response"] = "successfully sent"
						rawObject, _ := json.Marshal(emailMessage)
						updatingData["object_info"] = rawObject
						err = Update(v.ServiceDatabase, NotificationTable, object.Id, updatingData)
						if err != nil {
							v.Logger.Error("error updating email sent status notification", zap.String("error", err.Error()))
						}

					} else {
						v.Logger.Error("error sending email", zap.String("error", sendingError.Error()))
						retryCount = retryCount + 1
						emailMessage["retryCount"] = retryCount
						emailMessage["response"] = sendingError.Error()
						rawObject, _ := json.Marshal(emailMessage)
						updatingData["object_info"] = rawObject
						err = Update(v.ServiceDatabase, NotificationTable, object.Id, updatingData)
						if err != nil {
							v.Logger.Error("error updating email sent status notification", zap.String("error", err.Error()))
						}
					}
				}

			}
		} else {
			v.Logger.Error("error getting emails from tables", zap.String("error", err.Error()))
		}

		time.Sleep(tDuration)
	}
}
func (v *NotificationManager) sendEmail(emailMessage common.Message) error {
	// Create the email message
	message := &types.Message{
		Body: &types.Body{
			Html: &types.Content{
				Data: aws.String(emailMessage.Body["text/html"]),
			},
		},
		Subject: &types.Content{
			Data: aws.String(emailMessage.Subject),
		},
	}

	// Create the email input object
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: emailMessage.To,
		},
		Message: message,
		Source:  aws.String(emailMessage.From),
	}
	// Send the email
	result, err := v.Client.SendEmail(context.TODO(), input)
	if err != nil {
		v.Logger.Error("email sending failed to ..", zap.Any("email", emailMessage.To), zap.String("error", err.Error()))
		return err
	}
	v.Logger.Info("email is successfully dispatched..", zap.Any("email", emailMessage.To), zap.Any("result", result), zap.String("message_id", *result.MessageId))
	return nil

}

func (v *NotificationManager) SendTestEmail() {
	// read the email template and send
	// Open the template file
	if v.EmailConfig.TestEmail.IsEnabled {
		file, err := os.Open(v.EmailConfig.TestEmail.TestEmailTemplate)
		if err != nil {
			v.Logger.Error("unable to open the template file..should be resources folder under notification service..")
			return
		}
		defer file.Close()

		// Create a scanner to read the file line-by-line
		scanner := bufio.NewScanner(file)

		// Read the file contents into a string
		var templateString string
		for scanner.Scan() {
			templateString += scanner.Text() + "\n"
		}
		emailMessage := common.Message{}
		emailMessage.To = append(emailMessage.To, v.EmailConfig.TestEmail.To)
		emailMessage.Body = make(map[string]string)
		emailMessage.Body["text/html"] = templateString
		emailMessage.Subject = "Test email"
		emailMessage.From = "notifications@cerex.io"
		// Create the email message
		message := &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(emailMessage.Body["text/html"]),
				},
			},
			Subject: &types.Content{
				Data: aws.String(emailMessage.Subject),
			},
		}
		// Create the email input object
		input := &ses.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: emailMessage.To,
			},
			Message: message,
			Source:  aws.String(emailMessage.From),
		}

		// Send the email
		result, err := v.Client.SendEmail(context.TODO(), input)
		if err != nil {
			v.Logger.Debug("email sending failed to ..", zap.Any("email", emailMessage.To), zap.String("error", err.Error()))
			return
		}
		v.Logger.Debug("email is successfully dispatched..", zap.Any("email", emailMessage.To), zap.Any("result", result), zap.String("message_id", *result.MessageId))
		return
	}

}
