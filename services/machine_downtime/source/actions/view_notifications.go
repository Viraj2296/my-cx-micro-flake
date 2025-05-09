package actions

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type ViewNotifications struct {
	NotificationsId []int `json:"notificationsId"`
}

func (v *Actions) HandleViewNotifications(ctx *gin.Context) {
	v.Logger.Info("handleViewNotifications")

	var viewNotifications ViewNotifications
	if err := ctx.ShouldBindJSON(&viewNotifications); err != nil {
		ctx.JSON(http.StatusBadRequest, response.GeneralResponse{
			Code:    404,
			Message: "Invalid request body",
		})
		return
	}

	// Get auth service and user ID
	authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
	userId := common.GetUserId(ctx)

	// Update the view notification IDs
	if len(viewNotifications.NotificationsId) > 0 {
		for _, notificationId := range viewNotifications.NotificationsId {
			err := authService.AddViewPushNotificationIds(userId, notificationId)
			if err != nil {
				v.Logger.Error("error adding view push notification id",
					zap.Error(err),
					zap.Int("notificationId", notificationId))
				ctx.JSON(http.StatusOK, response.GeneralResponse{
					Code:    400,
					Message: "Failed to update",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{
			Code:    200,
			Message: "Successfully updated",
		})
		return
	}
}
