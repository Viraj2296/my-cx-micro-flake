package handler

import (
	"cx-micro-flake/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	EmailSendingFailed = 9000
)

func (v *NotificationService) sendTestEmail(ctx *gin.Context) {
	//emailSettingRequest := common.EmailSettingRequest{}
	//
	//if err := ctx.ShouldBindBodyWith(&emailSettingRequest, binding.JSON); err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	//ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 0, Message: "Email has been successfully sent, please check your inbox to verify before save the setting"})
	//return
	////userId := common.GetUserId(ctx)
	////authService := services.GetService("general_auth").ServiceInterface.(services.AuthInterface)
	//////userInfo := authService.GetUserInfoById(userId)
	//emailSettingRequest.SkipVerify = false
	//
	//emailSettingRequest.StartTLSPolicy = "NoStartTLS"
	//emailSettingRequest.ContentTypes = append(emailSettingRequest.ContentTypes, "text/html")
	//emailSettingRequest.ContentTypes = append(emailSettingRequest.ContentTypes, "text/plain")
	//var emailMessage = common.Message{
	//	To:          []string{emailSettingRequest.ToAddress},
	//	SingleEmail: false,
	//	Subject:     "Email Setting Verification",
	//	Body: map[string]string{
	//		"text/html": emailSettingRequest.SampleContent,
	//	},
	//	Info:          "",
	//	ReplyTo:       nil,
	//	EmbeddedFiles: nil,
	//	AttachedFiles: nil,
	//}
	//emailMessage.From = emailSettingRequest.FromAddress
	//err := ns.SendTestMailFromDynamicConfig(emailMessage, emailSettingRequest)
	//if err != nil {
	//	response.DispatchDetailedError(ctx, EmailSendingFailed,
	//		&response.DetailedError{
	//			Header:      "Error Sending Email",
	//			Description: err.Error(),
	//		})
	//	return
	//}
	ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 0, Message: "Email has been successfully sent, please check your inbox to verify before save the setting"})
}
