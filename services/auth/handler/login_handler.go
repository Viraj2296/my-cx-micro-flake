package handler

import (
	"cx-micro-flake/pkg/auth"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"

	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (as *AuthService) createNewPassword(ctx *gin.Context) {

	newPasswordRequest := NewPassword{}

	if err := ctx.ShouldBindBodyWith(&newPasswordRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	isValid, err := auth.IsTokenStringValid(newPasswordRequest.Token)
	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Invalid Token",
				Description: "Token is not valid, please make sure you are accessing the link sent from system",
			})
	}
	if !isValid {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Token Expired",
				Description: "Token is not valid, please make sure you are accessing the link sent from system",
			})
	}
	userId, err := auth.ExtractResourceId(newPasswordRequest.Token)
	dbConnection := as.BaseService.ReferenceDatabase

	err, userObject := Get(dbConnection, UserTable, userId)

	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Invalid Id",
				Description: "Token is invalid for further processing, try again later",
			})
	}
	userInfo := GetUserInfo(userObject.ObjectInfo)
	generatedHashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPasswordRequest.Password), bcrypt.DefaultCost)
	err = VerifyPassword(userInfo.Password, newPasswordRequest.Password)
	if err == nil {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Invalid Password",
				Description: "Use a new password not used before to increase your security",
			})
	}
	if userInfo.InvitationStatus == "Invited" {
		// this means, new user is creating the password, so he accepted our invitation
		userInfo.InvitationStatus = "Accepted"
		userInfo.InvitationToken = "" // invalidate the token we sent, so that others won't misuse the token
	}
	userInfo.Password = string(generatedHashedPassword)
	userInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	userInfo.ResetPasswordKey = ""
	err = Update(as.BaseService.ReferenceDatabase, UserTable, userId, userInfo.DatabaseSerialize())
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Your password has been successfully updated",
	})
}
func (as *AuthService) validateResetToken(ctx *gin.Context) {
	validateResetToken := ValidateResetToken{}

	if err := ctx.ShouldBindBodyWith(&validateResetToken, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	isValid, err := auth.IsTokenStringValid(validateResetToken.Token)
	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Invalid Token",
				Description: "Token is not valid, please make sure you are accessing the link sent from system",
			})
	}
	if !isValid {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Token Expired",
				Description: "Token is not valid, please make sure you are accessing the link sent from system",
			})
	}
	userId, err := auth.ExtractResourceId(validateResetToken.Token)
	dbConnection := as.BaseService.ReferenceDatabase

	err, generalUserObject := Get(dbConnection, UserTable, userId)

	if err != nil {
		response.DispatchDetailedError(ctx, InvalidEmailSetting,
			&response.DetailedError{
				Header:      "Invalid User",
				Description: "The link you are using is not valid as your credentials are removed from system, Please contact admin",
			})
	}
	userInfo := GetUserInfo(generalUserObject.ObjectInfo)
	var basicUserInfo = make(map[string]interface{})
	basicUserInfo["email"] = userInfo.Email
	basicUserInfo["fullName"] = userInfo.FullName
	ctx.JSON(http.StatusOK, basicUserInfo)
	return
}

func (as *AuthService) forgetPassword(ctx *gin.Context) {

	forgetPassword := ForgetPassword{}

	if err := ctx.ShouldBindBodyWith(&forgetPassword, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := as.BaseService.ReferenceDatabase

	err, listOfUsers := GetObjects(dbConnection, UserTable)

	if err != nil {
		as.BaseService.Logger.Error("Loading users had failed due to :", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("loading users had failed"), LoadingUsersFailed)
		return
	}

	as.BaseService.Logger.Info("loaded users", zap.Int("count", len(listOfUsers)))
	for _, userInterface := range listOfUsers {
		userInfo := GetUserInfo(userInterface.ObjectInfo)
		if userInfo.Email == forgetPassword.Email {
			as.BaseService.Logger.Info("checking user email ", zap.String("user", userInfo.Username))
			//found the email
			// we need to send an email to reset the password
			// update the password token , and generate the reset url

			// Specify the path to your HTML file
			resetPasswordEmailTemplate := "../services/auth/resource/reset_password.html"

			// Read the file content
			content, err := ioutil.ReadFile(resetPasswordEmailTemplate)
			if err != nil {
				as.BaseService.Logger.Error("error reading email template", zap.String("error", err.Error()))
				response.DispatchDetailedError(ctx, InvalidEmailSetting,
					&response.DetailedError{
						Header:      "Invalid Email Setting",
						Description: "Invalid email setting configured, please check the admin user setting to configure the correct email template setting",
					})
				return
			}

			// Convert the content to a string
			htmlString := string(content)

			zAuthToken, _ := auth.CreateToken(userInterface.Id, as.DefaultTokenExpiry)

			resetPasswordLink := as.EmailNotificationDomain + "/reset_password?token=" + zAuthToken
			htmlString = strings.Replace(htmlString, "[RESET_PASSWORD_LINK]", resetPasswordLink, -1)
			htmlString = strings.Replace(htmlString, "[USER]", userInfo.FullName, -1)
			var emailMessages []common.Message

			emailMessage := common.Message{
				To:          []string{userInfo.Email},
				SingleEmail: false,
				Subject:     "MES Reset Password",
				Body: map[string]string{
					"text/html": htmlString,
				},
				Info:          "",
				EmbeddedFiles: nil,
				AttachedFiles: nil,
			}

			emailMessages = append(emailMessages, emailMessage)
			notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
			err = notificationService.CreateMessages("906d0fd569404c59956503985b330132", emailMessages)
			if err != nil {
				response.DispatchDetailedError(ctx, InvalidEmailSetting,
					&response.DetailedError{
						Header:      "Internal Server Error",
						Description: "System couldn't able to send an email at this moment, this is due to necessary configurations are not in place",
					})
				return
			}

			// user has successfully logged in, so update the last login and last active

			userInfo.LastPasswordResetDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			userInfo.LastUpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			userInfo.ResetPasswordKey = resetPasswordLink
			Update(as.BaseService.ReferenceDatabase, UserTable, userInterface.Id, userInfo.DatabaseSerialize())

			ctx.JSON(http.StatusOK, response.GeneralResponse{
				Code:    0,
				Message: "Reset password link is sent to your registered email address",
			})
			return
		}
	}

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("invalid email, check your email is valid"), AuthenticationFailed, loginError)
}

//TODO we need to implement the force logout when the same user trying login but previous session is active
// this will apply only when user is allowed not to have multiple session, otherwise, we shouldn't user to login

func (as *AuthService) login(ctx *gin.Context) {

	loginCommand := LoginCommand{}

	if err := ctx.ShouldBindBodyWith(&loginCommand, binding.JSON); err != nil {
		as.BaseService.Logger.Error("Failed to bind data", zap.Error(err))
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	as.BaseService.Logger.Info("login command", zap.String("user", loginCommand.Username))

	systemReleaseVersion, err := as.getLatestVersion()
	if err != nil {
		as.BaseService.Logger.Error("Failed to fetch latest version", zap.Error(err))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("internal server error"), LoadingUsersFailed)
		return
	}
	as.BaseService.Logger.Info("system release version", zap.Any("systemReleaseVersion", systemReleaseVersion))
	err, listOfUsers := GetObjects(as.BaseService.ReferenceDatabase, UserTable)

	if err != nil {
		as.BaseService.Logger.Error("Loading users had failed due to :", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("loading users had failed"), LoadingUsersFailed)
		return
	}
	as.BaseService.Logger.Info("loaded users ", zap.Int("no_of_users", len(listOfUsers)))
	//hashPassword, _ := util.Hash("123456")
	for _, userInterface := range listOfUsers {
		userInfo := GetUserInfo(userInterface.ObjectInfo)
		if userInfo.Username == loginCommand.Username {
			as.BaseService.Logger.Info("checking user password ", zap.String("user", userInfo.Username))
			// check the objectStatus == Archived, we need to send error
			if userInfo.ObjectStatus == "Archived" {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Account is archived"), PasswordVerificationFailed, "Please note that your account has been already archived. Kindly contact your system administrator for assistance in this matter")
				return
			}
			if userInfo.Status == "disabled" {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Account Disabled"), PasswordVerificationFailed, "Please note that your account has not been enabled yet. Kindly contact your system administrator for assistance in this matter")
				return
			}
			if userInfo.InvitationStatus == "Pending" {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Activation Pending"), PasswordVerificationFailed, "Please note that your account is not activated yet, wait until system administrator send you the invitation link to access the system.")
				return
			}

			err = VerifyPassword(userInfo.Password, loginCommand.Password)
			if err != nil {
				as.BaseService.Logger.Info("verify password had failed ", zap.String("user", err.Error()))
				if userInfo.InvitationStatus == "Invited" {
					response.SendDetailedError(ctx, http.StatusBadRequest, getError("Password Reset"), PasswordVerificationFailed, "Please note that your account password is reset now, System administrator has already sent you the temporary password to login")
				} else {
					response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Username or Password"), PasswordVerificationFailed, "Please check your username or password again")
				}
				return

			}

			if len(userInfo.Department) == 0 && len(userInfo.SecondaryDepartmentList) == 0 {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Department"), PasswordVerificationFailed, "Please note that your account is not linked with any department, please ask system administrator to assign corresponding departments")
				return
			}

			// create a token based on record id, not based on user id (as we have record id as a primary key, not by userId)
			var refreshToken string
			var authToken string
			if userInfo.IsSessionTimeOut {
				authToken, _ = auth.CreateToken(userInterface.Id, as.DefaultTokenExpiry)
				refreshToken, _ = auth.CreateRefreshToken(userInterface.Id, loginCommand.Platform, as.DefaultRefreshTokenExpiry)
			} else {
				authToken, _ = auth.CreateInfToken(userInterface.Id)
				refreshToken, _ = auth.CreateRefreshInfToken(userInterface.Id, loginCommand.Platform)
				as.BaseService.Logger.Info("user session is not timeoutauth token", zap.String("user", userInfo.Username), zap.String("token", authToken), zap.String("refresh token", refreshToken))
			}

			// user has successfully logged in, so update the last login and last active

			if userInfo.InvitationStatus == "Invited" {
				userInfo.InvitationStatus = "Accepted"
				userInfo.PlainPassword = ""
			}
			userInfo.LastLogin = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			userInfo.LastActive = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			userInfo.UiApplicationVersion = loginCommand.UiApplicationVersion
			Update(as.BaseService.ReferenceDatabase, UserTable, userInterface.Id, userInfo.DatabaseSerialize())

			// save the user session
			userSession := UserSession{UserId: userInterface.Id, SessionId: authToken}
			sessionInfo := SessionInfo{SessionStatus: "logged-in", RequestHeader: ctx.Request.Header}
			var clientIpArray []string
			clientIp, _ := util.GetClientIPHelper(ctx.Request)
			clientIpArray = append(clientIpArray, clientIp)
			sessionInfo.RequestHeader["ClientIp"] = clientIpArray
			//TODO currently add the device finger print in the session, later move the database table
			sessionInfo.Fingerprint = loginCommand.Fingerprint
			sessionInfo.Platform = loginCommand.Platform
			userSession.SessionTime = time.Now().Unix()
			userSession.SessionInfo = sessionInfo.Serialize()
			as.BaseService.ReferenceDatabase.Model(&UserSession{}).Create(userSession)
			as.BaseService.Logger.Info("login user details  ", zap.Any("user_id", userInterface.Id), zap.Any("token", authToken), zap.Any("refresh_token", refreshToken))

			// this point, check the platform, and route there, the below logic is for web authentication
			if loginCommand.Platform == "mobile" && loginCommand.AppName == "" {
				// currently only allow supervisor can access the mobile
				labourManagementInterface := common.GetService("labour_management_module").ServiceInterface.(common.LabourManagementInterface)
				if labourManagementInterface != nil {
					var allowedMobileJobRoles = labourManagementInterface.GetMobileAllowedJobRoles()
					if len(allowedMobileJobRoles) > 0 {
						if util.HasInt(userInfo.JobRoleId, allowedMobileJobRoles) {
							loginResponse := LoginResponse{
								Token:        authToken,
								RefreshToken: refreshToken,
								SystemConfig: SystemConfig{},
							}

							as.BaseService.Logger.Info("mobile authentication successfully, roles has been found from the labour management configured list ", zap.String("user", loginCommand.Username))
							ctx.JSON(http.StatusOK, loginResponse)
							return
						}
					}
				}

				as.BaseService.Logger.Error("user role is not assigned or created for authentication")
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Job Role"), PasswordVerificationFailed, "Please note that your account is not linked with supervisor role, please ask system administrator to assign corresponding job role")
				return

			} else if loginCommand.Platform == "mobile" && loginCommand.AppName == "machine_downtime" {
				// insert the token
				// Create a new DeviceToken instance
				notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
				var listOfSubscriptions = notificationService.GetOneSignalSubscriptionId(loginCommand.DeviceToken)
				serialisedSubscriptionIds, _ := json.Marshal(listOfSubscriptions)
				deviceToken := UserDeviceToken{
					UserId:         userInterface.Id,
					DeviceToken:    loginCommand.DeviceToken,
					CreatedAt:      time.Now(),
					LastUsedAt:     time.Now(),
					CreatedBy:      userInterface.Id,
					LastUsedBy:     userInterface.Id,
					SubscriptionId: serialisedSubscriptionIds,
				}

				// Check if a record exists with the given UserId
				var existingToken UserDeviceToken
				err := as.BaseService.ReferenceDatabase.Where("user_id = ?", userInterface.Id).First(&existingToken).Error

				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						// Record does not exist, so insert it
						result := as.BaseService.ReferenceDatabase.Create(&deviceToken)
						if result.Error != nil {
							as.BaseService.Logger.Error("Failed to insert device token", zap.Error(result.Error))
						} else {
							as.BaseService.Logger.Info("user device token is successfully created", zap.Any("device_token", deviceToken))
						}
					} else {
						// Handle other errors
						as.BaseService.Logger.Error("Error querying for existing device token", zap.Error(err))
					}
				} else {
					// Record exists, so update it
					existingToken.DeviceToken = loginCommand.DeviceToken
					existingToken.LastUsedAt = time.Now()
					existingToken.LastUsedBy = userInterface.Id
					existingToken.SubscriptionId = serialisedSubscriptionIds

					result := as.BaseService.ReferenceDatabase.Save(&existingToken)
					if result.Error != nil {
						as.BaseService.Logger.Error("Failed to update device token", zap.Error(result.Error))
					} else {
						as.BaseService.Logger.Info("user device token is successfully updated", zap.Any("token", existingToken))
					}
				}
				loginResponse := LoginResponse{
					Token:        authToken,
					RefreshToken: refreshToken,
					SystemConfig: SystemConfig{},
				}

				as.BaseService.Logger.Info("machine downtime mobile authentication successfully", zap.String("user", loginCommand.Username))
				ctx.JSON(http.StatusOK, loginResponse)
				return
			} else {
				if userInfo.EnforceVersionUpdate {
					if loginCommand.UiApplicationVersion != systemReleaseVersion {
						response.SendDetailedError(ctx, http.StatusBadRequest, getError("Version Mismatch"), AuthenticationFailed, "A new version is available. Please refresh your browser to update the application.")
						as.BaseService.Logger.Warn("version mismatch detected", zap.String("expected", systemReleaseVersion), zap.String("provided", loginCommand.UiApplicationVersion))
						return
					}
				} else {
					as.BaseService.Logger.Warn("force version update is set false, so ignoring version update to user ", zap.Int("user_id", userInterface.Id))
				}
			}

			factoryService := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
			recordInfo := component.RecordInfo{}
			var dropDownArray []component.OrderedData
			index := 0
			var listOfDepartments []int
			var dropdownValue string

			for _, departmentId := range userInfo.Department {
				id := departmentId
				dropdownValue := factoryService.GetDepartmentName(departmentId)
				dropDownArray = append(dropDownArray, component.OrderedData{
					Id:    id,
					Value: dropdownValue,
				})
				listOfDepartments = append(listOfDepartments, id)

				if index == 0 {
					recordInfo.Index = id
					recordInfo.Value = dropdownValue
				}
				index = index + 1
			}

			for _, secondaryDepartmentId := range userInfo.SecondaryDepartmentList {
				var isExist bool
				isExist = false
				listOfDepartments = append(listOfDepartments, secondaryDepartmentId)
				for _, dropDownId := range dropDownArray {
					if dropDownId.Id == secondaryDepartmentId {
						isExist = true
					}
				}
				if !isExist {
					dropdownValue = factoryService.GetDepartmentName(secondaryDepartmentId)

					dropDownArray = append(dropDownArray, component.OrderedData{
						Id:    secondaryDepartmentId,
						Value: dropdownValue,
					})
				}

			}
			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)

			tvDisplayResponse := machineService.GetDepartmentDisplayEnabledMachines(ProjectID, listOfDepartments, dropDownArray)
			machineList := machineService.GetListOfEnergyManagementMachines(ProjectID)

			recordInfo.Data = dropDownArray
			systemConfig := SystemConfig{}
			systemConfig.TVDisplay = tvDisplayResponse
			systemConfig.EnergyManagementMachines = machineList
			systemConfig.SchedulerViewStartDate = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
			systemConfig.SchedulerViewEndDate = util.GetCurrentTimeWithOffSet(time.Hour*24*7, "2006-01-02T15:04:05.000Z") // always send 1 week schedule
			listOfAllowedAssemblyLines := machineService.GetListOfAllowedAssemblyLines(ProjectID, userInterface.Id)
			systemConfig.ListOfTVLabourManagementAssemblyLines = listOfAllowedAssemblyLines
			if userInfo.Type == "system_admin" {
				systemConfig.AllowedMenusIds = []string{"*"}
			} else {
				systemConfig.AllowedMenusIds = as.getListOfAllowedMenus(userInterface.Id)
			}
			var listOfMachineReports []ModuleReports
			listOfMachineReports = append(listOfMachineReports, ModuleReports{
				MenuId:      as.Report[0].MenuId,
				DashboardId: as.Report[0].DashboardId,
			})
			listOfMachineReports = append(listOfMachineReports, ModuleReports{
				MenuId:      as.Report[1].MenuId,
				DashboardId: as.Report[1].DashboardId,
			})
			var listOfReports []Reports
			reports := Reports{ModuleId: 1, ModuleReports: listOfMachineReports}
			listOfReports = append(listOfReports, reports)
			systemConfig.AllowedDepartments = recordInfo
			systemConfig.Reports = listOfReports

			var listOfModules []common.SystemModuleInfo
			err, listOfModulesObjects := GetObjects(as.BaseService.ServiceDatabases[ProjectID], "system_module")
			if err == nil {
				for _, moduleObject := range listOfModulesObjects {
					systemModuleInfo := GetSystemModuleInfo(moduleObject.ObjectInfo)
					listOfModules = append(listOfModules, common.SystemModuleInfo{Id: moduleObject.Id, DisplayName: systemModuleInfo.DisplayName, Description: systemModuleInfo.Description, ModuleName: systemModuleInfo.Name})
				}
			}

			if len(listOfModules) == 0 {
				systemConfig.AllowedModules = make([]int, 0)
			} else {
				if userInfo.Type == "system_admin" {
					for _, moduleInfo := range listOfModules {
						systemConfig.AllowedModules = append(systemConfig.AllowedModules, moduleInfo.Id)
					}
				} else {

					allowedModules := as.getUserModuleAccess(userInterface.Id)
					systemConfig.AllowedModules = allowedModules
				}

			}
			as.BaseService.Logger.Info("generated system config to user", zap.Any("system_config", systemConfig), zap.Any("user_id", userInterface.Id))
			systemConfig.MaintenanceMessage = ""
			loginResponse := LoginResponse{
				Token:        authToken,
				RefreshToken: refreshToken,
				SystemConfig: systemConfig,
			}

			as.BaseService.Logger.Info("authentication successfully ", zap.String("user", loginCommand.Username))
			ctx.JSON(http.StatusOK, loginResponse)
			return
		}
	}

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("authentication failed"), AuthenticationFailed, loginError)
	as.BaseService.Logger.Info("authentication has failed ", zap.String("user", loginCommand.Username))
}

func (as *AuthService) userLogout(ctx *gin.Context) {
	userId := common.GetUserId(ctx)
	err, userObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)

	if err != nil {
		response.DispatchDetailedError(ctx, ErrorGettingIndividualObjectInformation,
			&response.DetailedError{
				Header:      getError("Invalid User Id").Error(),
				Description: "Requested user id is not found in the system, check the user id, looks like your session is invalid or altered",
			})

		return

	}

	userInfo := GetUserInfo(userObject.ObjectInfo)
	userInfo.LastActive = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	Update(as.BaseService.ReferenceDatabase, UserTable, userId, userInfo.DatabaseSerialize())

	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:    0,
		Message: "Thank you for using the system. See you again",
	})

}

func (as *AuthService) getUserProfile(ctx *gin.Context) {
	userRecordId := common.GetUserId(ctx)

	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userRecordId)

	if err != nil {
		response.DispatchDetailedError(ctx, UserNotFoundInDatabase,
			&response.DetailedError{
				Header:      getError("Invalid User Id").Error(),
				Description: "Requested user id is not found in the system, check the user id, looks like your session is invalid or altered",
			})

		return
	}
	userInfo := GetUserInfo(generalObject.ObjectInfo)
	userInfo.Password = "********"
	ctx.JSON(http.StatusOK, userInfo)

}
func (as *AuthService) updateProfile(ctx *gin.Context) {
	userRecordId := common.GetUserId(ctx)
	var updateRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err, userObject := Get(as.BaseService.ReferenceDatabase, UserTable, userRecordId)
	updatingData := make(map[string]interface{})

	if err != nil {
		response.DispatchDetailedError(ctx, ErrorGettingIndividualObjectInformation,
			&response.DetailedError{
				Header:      getError("Invalid User Id").Error(),
				Description: "Requested user id is not found in the system, check the user id, looks like your session is invalid or altered",
			})

		return

	}

	updatingData["object_info"] = as.ComponentManager.GetUpdateRequest(updateRequest, userObject.ObjectInfo, UserComponent)
	err = Update(as.BaseService.ReferenceDatabase, UserTable, userRecordId, updatingData)
	if err != nil {
		response.DispatchDetailedError(ctx, ErrorUpdatingObjectInformation,
			&response.DetailedError{
				Header:      getError("Error Updating User Profile").Error(),
				Description: "Internal system error updating user profile information, please try again later.",
			})

		return
	}

}

func (as *AuthService) renewToken(ctx *gin.Context) {
	userId, err := auth.ExtractRefreshTokenID(ctx.Request)

	if err != nil {
		as.BaseService.Logger.Error("invalid refresh token, check the token format again", zap.String("error", err.Error()), zap.Any("user_id", userId))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid refresh token, check the token format again"), InvalidToken)
		return
	}
	platform, err := auth.ExtractRefreshTokenPlatform(ctx.Request)
	if err != nil {
		as.BaseService.Logger.Error("invalid refresh token, check the token format again", zap.String("error", err.Error()), zap.Any("user_id", userId))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid refresh token, check the token format again"), InvalidToken)
		return
	}
	err, generalObject := Get(as.BaseService.ReferenceDatabase, UserTable, userId)

	if err != nil {
		as.BaseService.Logger.Error("invalid user id, user not found", zap.String("error", err.Error()), zap.Any("user_id", userId))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid user id, user not found"), UserNotFoundInDatabase)
		return
	}

	userInfo := GetUserInfo(generalObject.ObjectInfo)
	if platform == "web" {
		systemReleaseVersion, err := as.getLatestVersion()
		as.BaseService.Logger.Info("system release version", zap.Any("system release version", systemReleaseVersion))
		if err != nil {
			as.BaseService.Logger.Error("Failed to fetch latest version", zap.Error(err))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("internal server error"), LoadingUsersFailed)
			return
		}
		if userInfo.EnforceVersionUpdate {
			if userInfo.UiApplicationVersion != systemReleaseVersion {
				response.SendDetailedError(ctx, http.StatusBadRequest, getError("Version Mismatch"), VersionMismatch, "A new version is available. Please refresh your browser to update the application.")
				as.BaseService.Logger.Warn("version mismatch detected", zap.String("expected", systemReleaseVersion), zap.String("provided", userInfo.UiApplicationVersion))
				return
			}
		} else {
			as.BaseService.Logger.Warn("force version update is set false, so ignoring version update to user ", zap.Int("user_id", userId))
		}

	}

	if userInfo.Status == UserStatusEnabled {

		var refreshToken string
		var token string
		if userInfo.IsSessionTimeOut {
			token, _ = auth.CreateToken(userId, as.DefaultTokenExpiry)
			refreshToken, _ = auth.CreateRefreshToken(userId, platform, as.DefaultRefreshTokenExpiry)
		} else {
			token, _ = auth.CreateInfToken(userId)
			refreshToken, _ = auth.CreateRefreshInfToken(userId, platform)
		}

		loginResponse := LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
		}
		as.BaseService.Logger.Info("user token renew successfully ", zap.Any("token", token), zap.String("refresh_token", refreshToken), zap.Any("user_id", userId))
		ctx.JSON(http.StatusOK, loginResponse)
	} else {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid user status"), InvalidUserStatus)
		return
	}

}

func (as *AuthService) healthHandler(ctx *gin.Context) {
	as.BaseService.Logger.Info("Sends an empty response with HTTP 200 status")
	ctx.Status(200)
}
