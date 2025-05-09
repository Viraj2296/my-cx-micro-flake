package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	ContentRecordTrailTable = "content_record_trail"

	ContentComponentTable = "content_component"
	ContentMasterTable    = "content_master"

	TimeLayout = "2006-01-02T15:04:05.000Z"
	DateLayout = "2006-01-02"

	GetMetaInfoAction     = "meta_info"
	ContentProperties     = "content_properties"
	SharingInfo           = "sharing_info"
	GetFavorites          = "get_favorites"
	RecentFiles           = "recent_files"
	SharedWithMe          = "shared_with_me"
	CreateFolderAction    = "create_folder"
	RemoveFromFavorite    = "remove_from_favorite"
	AddToFavorite         = "add_to_favorite"
	UpdateDirectoryName   = "rename_directory"
	SendShareNotification = "send_share_notification"

	ContentMasterComponent = "content_master"

	FolderExistError     = 1000
	FolderCreationFailed = 1001

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ProjectID = "906d0fd569404c59956503985b330132"

	ContentShareRoleIdOwner  = 1
	ContentShareRoleIdViewer = 2
	ContentShareRoleIdEditor = 3

	SharePermissionViewerRoleId    = 1 // Viewer // default one
	SharePermissionCommenterRoleId = 2 // Commenter
	SharePermissionEditorRoleId    = 3 // Editor

)

type InvalidRequest struct {
	Message string `json:"message"`
}

func getError(errorString string) error {
	return errors.New(errorString)
}

func sendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func sendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, common.ObjectNotFound,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}
