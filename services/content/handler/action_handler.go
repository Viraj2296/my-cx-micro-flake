package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ContentPropertiesResponse struct {
	Name         string        `json:"name"`
	FileTypeIcon string        `json:"fileTypeIcon"`
	Description  string        `json:"description"`
	Type         string        `json:"type"`
	Size         string        `json:"size"`
	Location     string        `json:"location"`
	LocationId   int           `json:"locationId"`
	Owner        string        `json:"owner"`
	Modified     string        `json:"modified"`
	Opened       string        `json:"opened"` // this will be useful for shared files
	Created      string        `json:"created"`
	Tags         []string      `json:"tags"`
	PreviewImage string        `json:"previewImage"`
	AccessDetail []interface{} `json:"accessDetail"`
}

func (v *ContentService) handleGetIndividualRecordAction(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)

	if actionName == GetMetaInfoAction {
		recordId := util.GetRecordId(ctx)
		projectId := util.GetProjectId(ctx)
		v.BaseService.Logger.Info("downloading requested file name", zap.Any("record id ", recordId))

		dbConnection := v.BaseService.GetDatabase(projectId)
		if dbConnection == nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("internal server error establishing database connection"), InternalServerErrorLocatingDatabaseConnection)
			return

		}
		err, generalContentObject := Get(dbConnection, ContentMasterTable, recordId)

		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("requested file is not found in the server"), RequestFileIdNotFoundInDatabase)
			return
		}
		content := ContentMaster{ObjectInfo: generalContentObject.ObjectInfo}
		ctx.JSON(http.StatusOK, content.GetContentInfo())
	} else if actionName == ContentProperties {
		recordId := util.GetRecordId(ctx)
		projectId := util.GetProjectId(ctx)
		v.BaseService.Logger.Info("downloading requested file name", zap.Any("record id ", recordId))

		dbConnection := v.BaseService.GetDatabase(projectId)
		err, generalContentObject := Get(dbConnection, ContentMasterTable, recordId)

		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Invalid Resource Id",
					Description: "Failed to find your file, should have been removed or doesn't exist",
				})
			return
		}
		content := ContentMaster{ObjectInfo: generalContentObject.ObjectInfo}
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		contentPropertiesResponse := ContentPropertiesResponse{}
		contentInfo := content.GetContentInfo()
		userId := common.GetUserId(ctx)
		var ownerName = ""
		if userId == contentInfo.CreatedBy {
			ownerName = "me"
		} else {
			userInfoCreated := authService.GetUserInfoById(contentInfo.CreatedBy)
			ownerName = userInfoCreated.FullName
		}

		userInfoLastUpdated := authService.GetUserInfoById(contentInfo.LastUpdatedBy)
		contentPropertiesResponse.Type = contentInfo.MIMEType
		contentPropertiesResponse.Created = util.ConvertTimeToTimeZonLongCorrected("Asia/Singapore", contentInfo.CreatedAt)
		contentPropertiesResponse.Modified = util.ConvertTimeToTimeZonLongCorrected("Asia/Singapore", contentInfo.LastUpdatedAt) + " by " + userInfoLastUpdated.FullName
		contentPropertiesResponse.Owner = ownerName
		contentPropertiesResponse.Size = contentInfo.Size
		contentPropertiesResponse.Tags = contentInfo.Tags
		contentPropertiesResponse.Opened = "-"
		contentPropertiesResponse.PreviewImage = contentInfo.FilePreviewImage
		// we don't need to send the ../content/storage/
		actualLocation := contentInfo.Path
		replacedLocation := strings.Replace(actualLocation, "../content/storage/", "", 1)

		contentPropertiesResponse.Location = replacedLocation
		contentPropertiesResponse.Description = contentInfo.Description
		contentPropertiesResponse.FileTypeIcon = contentInfo.FileTypeIcon
		if contentInfo.ChainReference != "" {
			// get the last element to indicate where it is located
			arrayOfReferencePath := strings.Split(contentInfo.ChainReference, ":")
			referenceFolderId, _ := strconv.Atoi(arrayOfReferencePath[0])
			contentPropertiesResponse.LocationId = referenceFolderId
		}
		accessDetails := make([]interface{}, 0)
		if len(contentInfo.Share) > 0 {
			// this has been shared with others
			for _, sharedUserId := range contentInfo.Share {
				userInfo := authService.GetUserInfoById(sharedUserId.SharedWith)
				accessDetails = append(accessDetails, userInfo)
			}
		}
		contentPropertiesResponse.AccessDetail = accessDetails
		ctx.JSON(http.StatusOK, contentPropertiesResponse)

	} else if actionName == SharingInfo {
		// we need to provide the sharing info
		recordId := util.GetRecordId(ctx)
		projectId := util.GetProjectId(ctx)
		v.BaseService.Logger.Info("downloading requested file name", zap.Any("record id ", recordId))

		dbConnection := v.BaseService.GetDatabase(projectId)
		err, generalContentObject := Get(dbConnection, ContentMasterTable, recordId)

		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Invalid Resource Id",
					Description: "Failed to find your file, should have been removed or doesn't exist",
				})
			return
		}
		content := ContentMaster{ObjectInfo: generalContentObject.ObjectInfo}
		createdBy := content.GetContentInfo().CreatedBy
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		ownerUserInfo := authService.GetUserInfoById(createdBy)

		var accessDetails []interface{}
		if len(content.GetContentInfo().Share) > 0 {
			// this has been shared with others
			for _, sharedUserId := range content.GetContentInfo().Share {
				userInfo := authService.GetUserInfoById(sharedUserId.SharedWith)
				accessDetails = append(accessDetails, userInfo)
			}
		}
		accessDetails = append(accessDetails, ownerUserInfo)
		listOfUsers := authService.GetUserList()
		var sharingInfoResponse = make(map[string]interface{})
		sharingInfoResponse["userListMaster"] = component.RecordInfo{
			Data:   listOfUsers,
			IsEdit: false,
		}
		sharingInfoResponse["existingSharingInfo"] = component.RecordInfo{
			Data:   accessDetails,
			IsEdit: false,
		}

		ctx.JSON(http.StatusOK, sharingInfoResponse)

	}

}

func (v *ContentService) handleGetAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	if actionName == GetFavorites {
		// TODO implement , send the files which contains the user id, get the user id from contxt ,and search what are the files and then send list getRecords format.(you can reused the get records)
		userId := common.GetUserId(ctx)
		projectId := util.GetProjectId(ctx)
		componentName := util.GetComponentName(ctx)
		dbConnection := v.BaseService.GetDatabase(projectId)

		conditionString := " JSON_CONTAINS(object_info," + "'" + strconv.Itoa(userId) + "', '$.favoriteList') "
		listOfContentObjects, err := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)

		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("unknown user id"), RequestFileIdNotFoundInDatabase)
			return
		}
		var tableRecordsResponse datatypes.JSON
		if listOfContentObjects != nil {
			totalRecords := int64(len(*listOfContentObjects))
			_, tableRecordsResponse = v.ComponentManager.GetTableRecords(dbConnection, listOfContentObjects, totalRecords, componentName, "", zone)

		}

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	} else if actionName == RecentFiles {
		projectId := util.GetProjectId(ctx)
		componentName := util.GetComponentName(ctx)
		dbConnection := v.BaseService.GetDatabase(projectId)

		conditionString := "(CAST((JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.lastUpdatedAt\"))) AS DATETIME)  <  NOW() - INTERVAL 1 DAY)  and JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isFile\")) = 'true'"
		listOfContentObjects, err := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)

		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("unknown user id"), RequestFileIdNotFoundInDatabase)
			return
		}
		var tableRecordsResponse datatypes.JSON
		if listOfContentObjects != nil {
			totalRecords := int64(len(*listOfContentObjects))
			_, tableRecordsResponse = v.ComponentManager.GetTableRecords(dbConnection, listOfContentObjects, totalRecords, componentName, "", zone)

		}

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	} else if actionName == SharedWithMe {
		projectId := util.GetProjectId(ctx)
		componentName := util.GetComponentName(ctx)
		dbConnection := v.BaseService.GetDatabase(projectId)

		conditionString := "(CAST((JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.lastUpdatedAt\"))) AS DATETIME)  <  NOW() - INTERVAL 1 DAY)  and JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isFile\")) = 'true'"
		listOfContentObjects, err := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)

		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("unknown user id"), RequestFileIdNotFoundInDatabase)
			return
		}
		var tableRecordsResponse datatypes.JSON
		if listOfContentObjects != nil {
			totalRecords := int64(len(*listOfContentObjects))
			_, tableRecordsResponse = v.ComponentManager.GetTableRecords(dbConnection, listOfContentObjects, totalRecords, componentName, "", zone)

		}

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *ContentService) handleIndividualRecordPOSTAction(ctx *gin.Context) {
	actionName := util.GetActionName(ctx)
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.GetDatabase(projectId)
	if actionName == AddToFavorite {
		recordId := util.GetRecordId(ctx)
		userId := common.GetUserId(ctx)
		err, contentMasterGeneralObject := Get(dbConnection, ContentMasterTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Invalid Operation",
					Description: "Failed to find your file, should have been removed or doesn't exist",
				})
			return
		}
		contentMaster := ContentMaster{ObjectInfo: contentMasterGeneralObject.ObjectInfo}
		contentInfo := contentMaster.GetContentInfo()
		contentInfo.FavoriteList = append(contentInfo.FavoriteList, userId)
		var updatingData = make(map[string]interface{})

		updatingData["object_info"] = contentInfo.Serialize()
		Update(dbConnection, ContentMasterTable, recordId, updatingData)
		ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 200, Message: "Great, Now the content  [" + contentInfo.Name + "] is your favorite"})
	} else if actionName == RemoveFromFavorite {
		recordId := util.GetRecordId(ctx)
		userId := common.GetUserId(ctx)
		err, contentMasterGeneralObject := Get(dbConnection, ContentMasterTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Invalid Operation",
					Description: "Failed to find your file, should have been removed or doesn't exist",
				})
			return
		}
		contentMaster := ContentMaster{ObjectInfo: contentMasterGeneralObject.ObjectInfo}
		contentInfo := contentMaster.GetContentInfo()
		contentInfo.FavoriteList = append(contentInfo.FavoriteList, userId)
		removedList := util.RemoveFromIntArray(contentInfo.FavoriteList, userId)
		contentInfo.FavoriteList = removedList
		var updatingData = make(map[string]interface{})

		updatingData["object_info"] = contentInfo.Serialize()
		Update(dbConnection, ContentMasterTable, recordId, updatingData)
		ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 200, Message: "[ " + contentInfo.Name + "] is now removed from favorite"})
	} else if actionName == SendShareNotification {
		recordId := util.GetRecordId(ctx)
		var updateRequest = make(map[string]interface{})
		if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err, contentMasterGeneralObject := Get(dbConnection, ContentMasterTable, recordId)
		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Invalid Operation",
					Description: "Failed to find your file, should have been removed or doesn't exist",
				})
			return
		}
		contentMaster := ContentMaster{ObjectInfo: contentMasterGeneralObject.ObjectInfo}
		contentInfo := contentMaster.GetContentInfo()

		s := reflect.ValueOf(updateRequest["sharedUserList"])
		userIdList := make([]int, 0)
		for i := 0; i < s.Len(); i++ {
			userIdList = append(userIdList, util.InterfaceToInt(s.Index(i).Interface()))
		}

		userEmailList := make([]string, 0)

		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)

		for _, userId := range userIdList {
			userBasicInfo := authService.GetUserInfoById(int(userId))

			if userBasicInfo.Email != "" {
				userEmailList = append(userEmailList, userBasicInfo.Email)
			}
		}

		var emailMessages []common.Message

		// cs.BaseService.Logger.Infow("email template:", "emailTemplate: ", emailTemplateInfo.Template)
		emailMessage := common.Message{
			To:          userEmailList,
			SingleEmail: false,
			Subject:     "Share File",
			Body: map[string]string{
				"text/html": "Share this " + contentInfo.Url,
			},
			Info:          "",
			ReplyTo:       make([]string, 0),
			EmbeddedFiles: nil,
			AttachedFiles: nil,
		}

		emailMessages = append(emailMessages, emailMessage)
		notificationService := common.GetService("notification_module").ServiceInterface.(common.NotificationInterface)
		err = notificationService.CreateMessages("906d0fd569404c59956503985b330132", emailMessages)

		fmt.Println("Error in sending email ", err)

		if err != nil {
			response.DispatchDetailedError(ctx, FolderCreationFailed,
				&response.DetailedError{
					Header:      "Failed to send email",
					Description: "Can't send notification to given user",
				})
			return
		}

		roleId := util.InterfaceToInt(updateRequest["role"])
		for _, userId := range userIdList {
			sharePermission := SharePermission{SharedWith: int(userId), RoleId: roleId}
			contentInfo.Share = append(contentInfo.Share, sharePermission)
		}

		fmt.Println("Content info ", contentInfo)

		var updatingData = make(map[string]interface{})

		updatingData["object_info"] = contentInfo.Serialize()
		Update(dbConnection, ContentMasterTable, recordId, updatingData)

	}
}
func (v *ContentService) handlePOSTAction(ctx *gin.Context) {

	actionName := util.GetActionName(ctx)
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.GetDatabase(projectId)
	if actionName == CreateFolderAction {
		createFolderRequest := CreateFolderRequest{}
		if err := ctx.ShouldBindBodyWith(&createFolderRequest, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		storageDirectory := v.ContentConfig.StorageDirectory
		referencePathId := util.InterfaceToInt(createFolderRequest.ReferencePathId)
		if referencePathId == 0 {
			//we are creating the root one
			folderInfo, err := os.Stat(storageDirectory + "/" + createFolderRequest.FolderName)
			if os.IsNotExist(err) {
				err = os.Mkdir(storageDirectory+"/"+createFolderRequest.FolderName, 0755)
				if err != nil {
					response.DispatchDetailedError(ctx, FolderCreationFailed,
						&response.DetailedError{
							Header:      "Folder creation failed",
							Description: "Invalid folder structure, please check the reference folder",
						})
					return
				}

				contentInfo := ContenetMasterInfo{
					Name:          createFolderRequest.FolderName,
					Size:          "",
					Path:          storageDirectory + "/" + createFolderRequest.FolderName,
					IsFile:        false,
					MIMEType:      "folder",
					CreatedAt:     util.GetCurrentTime(TimeLayout),
					CreatedBy:     0,
					LastUpdatedAt: util.GetCurrentTime(TimeLayout),
					LastUpdatedBy: 0,
					ObjectStatus:  common.ObjectStatusActive,
				}

				serializedContentMasterInfo, _ := json.Marshal(contentInfo)
				contentMaster := component.GeneralObject{
					ObjectInfo: serializedContentMasterInfo,
				}

				err, recordId := Create(dbConnection, ContentMasterTable, contentMaster)
				if err != nil {
					//TODO delete the folder created
					response.DispatchDetailedError(ctx, FolderCreationFailed,
						&response.DetailedError{
							Header:      "Internal System Error",
							Description: "System error occured during folder creation",
						})
					return
				}
				var updatingData = make(map[string]interface{})

				contentInfo.ChainReference = "0:" + strconv.Itoa(recordId)
				updatingData["object_info"] = contentInfo.Serialize()
				Update(dbConnection, ContentMasterTable, recordId, updatingData)
				ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 200, Message: "Your folder [" + createFolderRequest.FolderName + "] is successfully created"})

			} else {
				v.BaseService.Logger.Info("folder info", zap.Any("info", folderInfo.Name()), zap.Any("size", folderInfo.Size()), zap.Any("mode", folderInfo.Mode()), zap.Any("time", folderInfo.ModTime()))
				response.DispatchDetailedError(ctx, FolderExistError,
					&response.DetailedError{
						Header:      "Folder already exist",
						Description: "The folder that you are trying to create is already exist, please choose different name or use the existing one",
					})
				return
			}

		} else {
			//Adding parent path#################################################################
			conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) like '%" + strconv.Itoa(referencePathId) + "' and JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.isFile\")) = 'false' "
			listOfObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)

			var chainReference string
			if listOfObjects == nil {
				v.BaseService.Logger.Info("Parent folder doesn't exist")
				response.DispatchDetailedError(ctx, FolderExistError,
					&response.DetailedError{
						Header:      "Incorrect reference Id",
						Description: "Can't create folder under given parent directory",
					})
				return
			} else {
				parentContenetMasterInfo := ContenetMasterInfo{}
				json.Unmarshal((*listOfObjects)[0].ObjectInfo, &parentContenetMasterInfo)
				chainReference = parentContenetMasterInfo.ChainReference
			}

			referencePathIdList := strings.Split(chainReference, ":")

			recordIdList := "("
			for _, recordId := range referencePathIdList {
				recordIdList += recordId + ","
			}
			recordIdList = util.TrimSuffix(recordIdList, ",")
			recordIdList += ")"

			contentConditionString := " id in " + recordIdList
			listOfContentObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, contentConditionString)

			for _, folderObject := range *listOfContentObjects {
				contenetMasterInfo := ContenetMasterInfo{}
				json.Unmarshal([]byte(folderObject.ObjectInfo), &contenetMasterInfo)
				storageDirectory += "/" + contenetMasterInfo.Name
			}

			//End#################################################################

			folderInfo, err := os.Stat(storageDirectory + "/" + createFolderRequest.FolderName)
			if os.IsNotExist(err) {
				err = os.Mkdir(storageDirectory+"/"+createFolderRequest.FolderName, 0755)
				if err != nil {
					response.DispatchDetailedError(ctx, FolderCreationFailed,
						&response.DetailedError{
							Header:      "Folder creation failed",
							Description: "Invalid folder structure, please check the reference folder",
						})
					return
				}

				contentInfo := ContenetMasterInfo{
					Name:          createFolderRequest.FolderName,
					Size:          "",
					Path:          storageDirectory + "/" + createFolderRequest.FolderName,
					IsFile:        false,
					MIMEType:      "folder",
					CreatedAt:     util.GetCurrentTime(TimeLayout),
					CreatedBy:     0,
					LastUpdatedAt: util.GetCurrentTime(TimeLayout),
					LastUpdatedBy: 0,
					ObjectStatus:  common.ObjectStatusActive,
				}

				serializedContentMasterInfo, _ := json.Marshal(contentInfo)
				contentMaster := component.GeneralObject{
					ObjectInfo: serializedContentMasterInfo,
				}

				err, recordId := Create(dbConnection, ContentMasterTable, contentMaster)
				if err != nil {
					//TODO delete the folder created
					response.DispatchDetailedError(ctx, FolderCreationFailed,
						&response.DetailedError{
							Header:      "Internal System Error",
							Description: "System error occured during folder creation",
						})
					return
				}
				var updatingData = make(map[string]interface{})

				contentInfo.ChainReference = chainReference + ":" + strconv.Itoa(recordId)
				updatingData["object_info"] = contentInfo.Serialize()
				Update(dbConnection, ContentMasterTable, recordId, updatingData)
				ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 200, Message: "Your folder [" + createFolderRequest.FolderName + "] is successfully created"})

			} else {
				v.BaseService.Logger.Info("folder info", zap.Any("info", folderInfo.Name()), zap.Any("size", folderInfo.Size()), zap.Any("mode", folderInfo.Mode()), zap.Any("time", folderInfo.ModTime()))
				response.DispatchDetailedError(ctx, FolderExistError,
					&response.DetailedError{
						Header:      "Folder already exist",
						Description: "The folder that you are trying to create is already exist, please choose different name or use the existing one",
					})
				return
			}

		}

	} else if actionName == UpdateDirectoryName {
		requestParam := make(map[string]interface{})
		if err := ctx.ShouldBindBodyWith(&requestParam, binding.JSON); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		err := updateFileFolderName(dbConnection, requestParam)

		if err != nil {
			response.DispatchDetailedError(ctx, FolderExistError,
				&response.DetailedError{
					Header:      "Requested file or folder doesn't exist",
					Description: "The folder or file that you are trying to rename isn't already exist",
				})
			return
		}

		ctx.JSON(http.StatusOK, response.GeneralResponse{Code: 200, Message: "Your folder or file [" + util.InterfaceToString(requestParam["newName"]) + "] is successfully renamed"})
	}

}

func updateFileFolderName(dbConnection *gorm.DB, renameRequest map[string]interface{}) error {
	referencePathId := util.InterfaceToInt(renameRequest["referencePathId"])
	err, existingData := Get(dbConnection, ContentMasterTable, referencePathId)

	if err != nil {
		return nil
	}

	contenetMasterInfo := ContenetMasterInfo{}
	json.Unmarshal(existingData.ObjectInfo, &contenetMasterInfo)

	oldName := contenetMasterInfo.Name

	newName := util.InterfaceToString(renameRequest["newName"])
	contenetMasterInfo.Name = newName

	if !contenetMasterInfo.IsFile {
		//Updating folder path
		pathList := strings.Split(contenetMasterInfo.Path, "/")
		if len(pathList) > 0 {
			pathList = pathList[:len(pathList)-1]
		}
		newPath := ""
		for index, folderName := range pathList {
			if index == 0 {
				continue
			} else {
				newPath += "/" + folderName
			}

		}
		newPath += "/" + newName

		contenetMasterInfo.Path = newPath

		//Updating child paths
		contentConditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) like '%" + strconv.Itoa(existingData.Id) + ":%'"
		listOfContentObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, contentConditionString)

		if listOfContentObjects != nil {
			for _, contents := range *listOfContentObjects {
				childContenetMasterInfo := ContenetMasterInfo{}
				json.Unmarshal(contents.ObjectInfo, &childContenetMasterInfo)
				childContenetMasterInfo.Path = strings.Replace(childContenetMasterInfo.Path, "/"+oldName+"/", "/"+newName+"/", -1)

				var updatingChildData = make(map[string]interface{})
				updatingChildData["object_info"] = childContenetMasterInfo.Serialize()
				Update(dbConnection, ContentMasterTable, contents.Id, updatingChildData)

			}
		}

	}

	var updatingData = make(map[string]interface{})
	updatingData["object_info"] = contenetMasterInfo.Serialize()
	Update(dbConnection, ContentMasterTable, referencePathId, updatingData)

	return nil
}
