package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

type NewPassword struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}
type ValidateResetToken struct {
	Token string `json:"token"`
}

type PasswordUpdateRequest struct {
	NewPassword string `json:"newPassword"`
	Password    string `json:"password"`
}

type IndividualUserPasswordUpdateRequest struct {
	NewPassword string `json:"newPassword"`
}

func (as *AuthService) resetIndividualUserPassword(ctx *gin.Context) {

	userId := common.GetUserId(ctx)
	dbConnection := as.BaseService.ReferenceDatabase
	userRecordId := util.GetRecordId(ctx)
	individualUserPasswordUpdateRequest := IndividualUserPasswordUpdateRequest{}

	if err := ctx.ShouldBindBodyWith(&individualUserPasswordUpdateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, UserTable, userRecordId)
	if err != nil {
		as.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, ErrorGettingIndividualObjectInformation,
			&response.DetailedError{
				Header:      "Invalid Resource ID",
				Description: "Error getting requested resource, please check your resource id again. This might be due to required resource id already archived or not existed",
			})

		return
	}
	userInfo := GetUserInfo(objectInterface.ObjectInfo)
	// check this update request contains newPassword field filled, if so do the password validations.
	currentPassword := userInfo.Password
	err = VerifyPassword(util.InterfaceToString(currentPassword), individualUserPasswordUpdateRequest.NewPassword)

	if err == nil {
		response.DispatchDetailedError(ctx, ErrorGettingIndividualObjectInformation,
			&response.DetailedError{
				Header:      "Invalid Password",
				Description: "Please use the different password, you are trying to update the same password as previous",
			})

		return
	}
	generatedPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(individualUserPasswordUpdateRequest.NewPassword), bcrypt.DefaultCost)
	userInfo.Password = string(generatedPasswordHash)
	// before update anything, we need to check necessary settings are there to send email
	userInfo.LastUpdatedBy = userId
	userInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	userInfo.PlainPassword = individualUserPasswordUpdateRequest.NewPassword
	userInfo.InvitationStatus = "Pending"
	serialisedData, _ := json.Marshal(userInfo)
	updatingData["object_info"] = serialisedData
	err = Update(as.BaseService.ReferenceDatabase, UserTable, userRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Your password has been successfully updated",
	})

}
func (as *AuthService) updatePassword(ctx *gin.Context) {
	dbConnection := as.BaseService.ReferenceDatabase
	userRecordId := common.GetUserId(ctx)

	passwordUpdateRequest := PasswordUpdateRequest{}

	if err := ctx.ShouldBindBodyWith(&passwordUpdateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var updateRequestFields = make(map[string]interface{})
	err, objectInterface := Get(dbConnection, UserTable, userRecordId)
	if err != nil {
		as.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.DispatchDetailedError(ctx, ErrorGettingIndividualObjectInformation,
			&response.DetailedError{
				Header:      "Invalid Resource ID",
				Description: "Error getting requested resource, please check your resource id again. This might be due to required resource id already archived or not existed",
			})

		return
	}

	json.Unmarshal(objectInterface.ObjectInfo, &updateRequestFields)
	// check this update request contains newPassword field filled, if so do the password validations.

	userInfo := GetUserInfo(objectInterface.ObjectInfo)
	err = VerifyPassword(userInfo.Password, passwordUpdateRequest.Password)
	if err == nil {
		// password is matched
		generatedHashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordUpdateRequest.NewPassword), bcrypt.DefaultCost)
		userInfo.Password = string(generatedHashedPassword)
	} else {
		as.BaseService.Logger.Error("invalid password, check your password again", zap.String("error", err.Error()))
		//send error saying, it is already scheduled
		response.DispatchDetailedError(ctx, InvalidCredentials,
			&response.DetailedError{
				Header:      "Invalid Credentials",
				Description: "Invalid password, check your password again. If you couldn't remember old password, please do proceed reset password",
			})
		return

	}
	// before update anything, we need to check necessary settings are there to send email

	err = Update(as.BaseService.ReferenceDatabase, UserTable, userRecordId, userInfo.DatabaseSerialize())
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Your password has been successfully updated",
	})
}

func (as *AuthService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	recordId := util.GetRecordId(ctx)
	targetTable := as.ComponentManager.GetTargetTable(componentName)
	err, generalObject := Get(as.BaseService.ReferenceDatabase, targetTable, recordId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Resource",
				Description: "The resource that you are trying to delete doesn't exist, Please check refresh page and try again",
			})
		return
	}
	if component.IsArchived(generalObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Resource Archived",
				Description: "The resource that you are trying to delete is already archived. This operation is not allowed",
			})
		return
	}
	var dependencyComponents []string
	var dependencyRecords int
	as.checkReference(as.BaseService.ReferenceDatabase, componentName, recordId, &dependencyComponents, &dependencyRecords)
	if dependencyRecords > 0 {
		var dependencyString string
		dependencyComponents = util.RemoveDuplicateString(dependencyComponents)
		dependencyString = " ["
		for index, dependencyComponent := range dependencyComponents {
			if index == len(dependencyComponents)-1 {
				dependencyString += dependencyComponent
			} else {
				dependencyString += dependencyComponent + " ->"
			}
		}
		dependencyString += " ]"
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			CanDelete: false,
			Code:      100,
			Message:   "There are dependencies bound to the resource that you are trying to remove. Removing this resource would create the chain removal on following resources " + dependencyString + " in " + strconv.Itoa(dependencyRecords) + " places, Please understand the risk of deleting this resource as all the dependencies would be achieved immediately, and this process is not reversible",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		CanDelete: true,
		Code:      100,
		Message:   "There are no dependencies bound to the resource that you are trying to remove. So, removing this resource won't affect others resource now, you can proceed !!",
	})

}
