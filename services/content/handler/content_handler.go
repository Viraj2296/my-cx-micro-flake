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
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getMIMEType(mimeType string) string {
	if mimeType == "application/octet-stream" {
		return "file"
	} else if mimeType == "image/x-xpixmap" {
		return "xpm"
	} else if mimeType == "application/x-7z-compressed" {
		return "7z"
	} else if mimeType == "application/zip" {
		return "zip"
	} else if mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return "xlsx"
	} else if mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		return "docx"
	} else if mimeType == "application/vnd.openxmlformats-officedocument.presentationml.presentation" {
		return "pptx"
	} else if mimeType == "application/epub+zip" {
		return "epub"
	} else if mimeType == "application/jar" {
		return "jar"
	} else if mimeType == "application/vnd.oasis.opendocument.text" {
		return "odt"
	} else if mimeType == "application/vnd.oasis.opendocument.text-template" {
		return "ott"
	} else if mimeType == "application/vnd.oasis.opendocument.spreadsheet" {
		return "ods"
	} else if mimeType == "application/vnd.oasis.opendocument.spreadsheet-template" {
		return "ots"
	} else if mimeType == "application/vnd.oasis.opendocument.presentation" {
		return "odp"
	} else if mimeType == "application/vnd.oasis.opendocument.presentation-template" {
		return "otp"
	} else if mimeType == "application/vnd.oasis.opendocument.graphics" {
		return "odg"
	} else if mimeType == "application/vnd.oasis.opendocument.graphics-template" {
		return "otg"
	} else if mimeType == "application/vnd.oasis.opendocument.formula" {
		return "odf"
	} else if mimeType == "application/vnd.oasis.opendocument.chart" {
		return "odc"
	} else if mimeType == "application/vnd.sun.xml.calc" {
		return "sxc"
	} else if mimeType == "application/pdf" {
		return "pdf"
	} else if mimeType == "application/vnd.fdf" {
		return "fdf"
	} else if mimeType == "application/x-ole-storage" {
		return "file"
	} else if mimeType == "application/x-ms-installer" {
		return "msi"
	} else if mimeType == "application/octet-stream" {
		return "aaf"
	} else if mimeType == "application/vnd.ms-outlook" {
		return "msg"
	} else if mimeType == "application/vnd.ms-excel" {
		return "xls"
	} else if mimeType == "application/vnd.ms-publisher" {
		return "pub"
	} else if mimeType == "application/vnd.ms-powerpoint" {
		return "ppt"
	} else if mimeType == "application/msword" {
		return "doc"
	} else if mimeType == "application/postscript" {
		return "ps"
	} else if mimeType == "image/vnd.adobe.photoshop" {
		return "psd"
	} else if mimeType == "application/pkcs7-signature" {
		return "p7s"
	} else if mimeType == "application/ogg" {
		return "ogg"
	} else if mimeType == "audio/ogg" {
		return "oga"
	} else if mimeType == "video/ogg" {
		return "ogv"
	} else if mimeType == "image/png" {
		return "png"
	} else if mimeType == "image/vnd.mozilla.apng" {
		return "png"
	} else if mimeType == "image/jpeg" {
		return "jpg"
	} else if mimeType == "image/jxl" {
		return "jxl"
	} else if mimeType == "image/jp2" {
		return "jp2"
	} else if mimeType == "image/jpx" {
		return "jpf"
	} else if mimeType == "image/jpm" {
		return "jpm"
	} else if mimeType == "image/jxs" {
		return "jxs"
	} else if mimeType == "image/gif" {
		return "gif"
	} else if mimeType == "image/webp" {
		return "webp"
	} else if mimeType == "application/vnd.microsoft.portable-executable" {
		return "exe"
	} else if mimeType == "application/x-elf" {
		return "file"
	} else if mimeType == "application/x-object" {
		return "file"
	} else if mimeType == "application/x-executable" {
		return "file"
	} else if mimeType == "application/x-sharedlib" {
		return "so"
	} else if mimeType == "application/x-coredump" {
		return "file"
	} else if mimeType == "application/x-archive" {
		return "a"
	} else if mimeType == "application/vnd.debian.binary-package" {
		return "deb"
	} else if mimeType == "application/x-tar" {
		return "tar"
	} else if mimeType == "application/x-xar" {
		return "xar"
	} else if mimeType == "application/x-bzip2" {
		return "bz2"
	} else if mimeType == "application/fits" {
		return "fits"
	} else if mimeType == "image/tiff" {
		return "tiff"
	} else if mimeType == "image/bmp" {
		return "bmp"
	} else if mimeType == "image/x-icon" {
		return "ico"
	} else if mimeType == "audio/mpeg" {
		return "mp3"
	} else if mimeType == "audio/flac" {
		return "flac"
	} else if mimeType == "audio/midi" {
		return "midi"
	} else if mimeType == "audio/ape" {
		return "ape"
	} else if mimeType == "audio/musepack" {
		return "mpc"
	} else if mimeType == "audio/amr" {
		return "amr"
	} else if mimeType == "audio/wav" {
		return "wav"
	} else if mimeType == "audio/aiff" {
		return "aiff"
	} else if mimeType == "audio/basic" {
		return "au"
	} else if mimeType == "video/mpeg" {
		return "mpeg"
	} else if mimeType == "video/quicktime" {
		return "mov"
	} else if mimeType == "video/quicktime" {
		return "mqv"
	} else if mimeType == "video/mp4" {
		return "mp4"
	} else if mimeType == "video/webm" {
		return "webm"
	} else if mimeType == "video/3gpp" {
		return "3gp"
	} else if mimeType == "video/3gpp2" {
		return "3g2"
	} else if mimeType == "video/x-msvideo" {
		return "avi"
	} else if mimeType == "video/x-flv" {
		return "flv"
	} else if mimeType == "video/x-matroska" {
		return "mkv"
	} else if mimeType == "video/x-ms-asf" {
		return "asf"
	} else if mimeType == "audio/aac" {
		return "aac"
	} else if mimeType == "audio/x-unknown" {
		return "voc"
	} else if mimeType == "audio/mp4" {
		return "mp4"
	} else if mimeType == "audio/x-m4a" {
		return "m4a"
	} else if mimeType == "application/vnd.apple.mpegurl" {
		return "m3u"
	} else if mimeType == "video/x-m4v" {
		return "m4v"
	} else if mimeType == "application/vnd.rn-realmedia-vbr" {
		return "rmvb"
	} else if mimeType == "application/gzip" {
		return "gz"
	} else if mimeType == "application/x-java-applet" {
		return "class"
	} else if mimeType == "application/x-shockwave-flash" {
		return "swf"
	} else if mimeType == "application/x-chrome-extension" {
		return "crx"
	} else if mimeType == "font/ttf" {
		return "ttf"
	} else if mimeType == "font/woff" {
		return "woff"
	} else if mimeType == "font/woff2" {
		return "woff2"
	} else if mimeType == "font/otf" {
		return "otf"
	} else if mimeType == "font/collection" {
		return "ttc"
	} else if mimeType == "application/vnd.ms-fontobject" {
		return "eot"
	} else if mimeType == "application/wasm" {
		return "wasm"
	} else if mimeType == "application/vnd.shx" {
		return "shx"
	} else if mimeType == "application/vnd.shp" {
		return "shp"
	} else if mimeType == "application/x-dbf" {
		return "dbf"
	} else if mimeType == "application/dicom" {
		return "dcm"
	} else if mimeType == "application/x-rar-compressed" {
		return "rar"
	} else if mimeType == "image/vnd.djvu" {
		return "djvu"
	} else if mimeType == "application/x-mobipocket-ebook" {
		return "mobi"
	} else if mimeType == "application/x-ms-reader" {
		return "lit"
	} else if mimeType == "image/bpg" {
		return "bpg"
	} else if mimeType == "application/vnd.sqlite3" {
		return "sqlite"
	} else if mimeType == "image/vnd.dwg" {
		return "dwg"
	} else if mimeType == "application/vnd.nintendo.snes.rom" {
		return "nes"
	} else if mimeType == "application/x-ms-shortcut" {
		return "lnk"
	} else if mimeType == "application/x-mach-binary" {
		return "macho"
	} else if mimeType == "audio/qcelp" {
		return "qcp"
	} else if mimeType == "image/x-icns" {
		return "icns"
	} else if mimeType == "image/heic" {
		return "heic"
	} else if mimeType == "image/heic-sequence" {
		return "heic"
	} else if mimeType == "image/heif" {
		return "heif"
	} else if mimeType == "image/heif-sequence" {
		return "heif"
	} else if mimeType == "image/vnd.radiance" {
		return "hdr"
	} else if mimeType == "application/marc" {
		return "mrc"
	} else if mimeType == "application/x-msaccess" {
		return "mdb"
	} else if mimeType == "application/x-msaccess" {
		return "accdb"
	} else if mimeType == "application/zstd" {
		return "zst"
	} else if mimeType == "application/vnd.ms-cab-compressed" {
		return "cab"
	} else if mimeType == "application/x-rpm" {
		return "rpm"
	} else if mimeType == "application/x-xz" {
		return "xz"
	} else if mimeType == "application/lzip" {
		return "lz"
	} else if mimeType == "application/x-bittorrent" {
		return "torrent"
	} else if mimeType == "application/x-cpio" {
		return "cpio"
	} else if mimeType == "application/tzif" {
		return "file"
	} else if mimeType == "image/x-xcf" {
		return "xcf"
	} else if mimeType == "image/x-gimp-pat" {
		return "pat"
	} else if mimeType == "image/x-gimp-gbr" {
		return "gbr"
	} else if mimeType == "model/gltf-binary" {
		return "glb"
	} else if mimeType == "image/avif" {
		return "avif"
	} else if mimeType == "application/x-installshield" {
		return "cab"
	} else if mimeType == "image/jxr" {
		return "jxr"
	} else if mimeType == "text/plain" {
		return "txt"
	} else if mimeType == "text/html" {
		return "html"
	} else if mimeType == "image/svg+xml" {
		return "svg"
	} else if mimeType == "text/xml" {
		return "xml"
	} else if mimeType == "application/rss+xml" {
		return "rss"
	} else if mimeType == "application/atom+xml" {
		return "atom"
	} else if mimeType == "model/x3d+xml" {
		return "x3d"
	} else if mimeType == "application/vnd.google-earth.kml+xml" {
		return "kml"
	} else if mimeType == "application/x-xliff+xml" {
		return "xlf"
	} else if mimeType == "model/vnd.collada+xml" {
		return "dae"
	} else if mimeType == "application/gml+xml" {
		return "gml"
	} else if mimeType == "application/gpx+xml" {
		return "gpx"
	} else if mimeType == "application/vnd.garmin.tcx+xml" {
		return "tcx"
	} else if mimeType == "application/x-amf" {
		return "amf"
	} else if mimeType == "application/vnd.ms-package.3dmanufacturing-3dmodel+xml" {
		return "3mf"
	} else if mimeType == "application/vnd.adobe.xfdf" {
		return "xfdf"
	} else if mimeType == "application/owl+xml" {
		return "owl"
	} else if mimeType == "text/x-php" {
		return "php"
	} else if mimeType == "application/javascript" {
		return "js"
	} else if mimeType == "text/x-lua" {
		return "lua"
	} else if mimeType == "text/x-perl" {
		return "pl"
	} else if mimeType == "text/x-python" {
		return "py"
	} else if mimeType == "application/json" {
		return "json"
	} else if mimeType == "application/geo+json" {
		return "geojson"
	} else if mimeType == "application/json" {
		return "har"
	} else if mimeType == "application/x-ndjson" {
		return "ndjson"
	} else if mimeType == "text/rtf" {
		return "rtf"
	} else if mimeType == "application/x-subrip" {
		return "srt"
	} else if mimeType == "text/x-tcl" {
		return "tcl"
	} else if mimeType == "text/csv" {
		return "csv"
	} else if mimeType == "text/tab-separated-values" {
		return "tsv"
	} else if mimeType == "text/vcard" {
		return "vcf"
	} else if mimeType == "text/calendar" {
		return "ics"
	} else if mimeType == "application/warc" {
		return "warc"
	} else if mimeType == "text/vtt" {
		return "vtt"
	}
	return "file"
}

func getFileTypeIcon(fileType string) string {
	if fileType == "pdf" {
		return "ci:file-pdf"
	} else if fileType == "doc" {
		return "bxs:file-doc"
	} else if fileType == "docx" {
		return "bi:filetype-docx"
	} else if fileType == "xls" {
		return "bi:filetype-xls"
	} else if fileType == "xlsx" {
		return "bi:filetype-xlsx"
	} else if fileType == "png" {
		return "ci:file-png"
	} else if fileType == "jpg" {
		return "bxs:file-jpg"
	}
	return "akar-icons:file"
}

type CreateFolderRequest struct {
	ReferencePathId int    `json:"referencePathId"`
	FolderName      string `json:"folderName"`
}

// getRecordFormData ShowAccount godoc
// @Summary Get the record form data to facilitate the update
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/content/component/{componentName}/record/{recordId} [get]
func (v *ContentService) getRecordFormData(ctx *gin.Context) {

	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	componentName := util.GetComponentName(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.GetDatabase(projectId)

	err, generalObject := Get(dbConnection, targetTable, recordId)
	fmt.Println("record id : ", recordId, "target table ;", targetTable, "error", err)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, getError("Requested file is not found in the server"), RequestFileIdNotFoundInDatabase)
		return
	}
	if componentName == ContentMasterComponent {
		contentMaster := ContentMaster{ObjectInfo: generalObject.ObjectInfo}
		v.BaseService.Logger.Info("serving content info", zap.Any("content_info", contentMaster.GetContentInfo()))
		contentDisposition := "filename=\"" + contentMaster.GetContentInfo().Name + "\""
		ctx.Writer.Header().Set("Content-Disposition", contentDisposition)
		mimeType, _ := mimetype.DetectFile(contentMaster.GetContentInfo().Path)
		ctx.Writer.Header().Set("content-type", mimeType.String())
		ctx.File(contentMaster.GetContentInfo().Path)
	}

}

// createNewResource ShowAccount godoc
// @Summary create new resource
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   recordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/content/component/{componentName}/records  [post]
func (v *ContentService) createNewResource(ctx *gin.Context) {
	var err error
	projectId := util.GetProjectId(ctx)
	uploadFile, err := ctx.FormFile("file")
	if err != nil {
		v.BaseService.Logger.Error("file parameter is not found in the request", zap.Any("ctx", ctx))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("extracting uploaded file, required parameter (file) not available, check the API specifications"), RequiredFileParamNotAvailable)
		return
	}
	dbConnection := v.BaseService.GetDatabase(projectId)
	referencePathId := ctx.Query("referencePathId")
	var dstDirectory string
	contentId := uuid.New().String()
	var chainReference string
	chainReference = "0"

	if referencePathId != "" {
		// get the folder path
		recordId := util.InterfaceToInt(referencePathId)
		if recordId == 0 {
			// upload to main directory
			//err, generalContentObject := Get(dbConnection, ContentMasterTable, recordId)
			//if err != nil {
			//	response.DispatchDetailedError(ctx, FolderExistError,
			//		&response.DetailedError{
			//			Header:      "Folder doesn't exist",
			//			Description: "Invalid folder structure, please check the reference folder",
			//		})
			//	return
			//}
			//contentMaster := ContentMaster{ObjectInfo: generalContentObject.ObjectInfo}
			//dstDirectory = contentMaster.GetContentInfo().Path
			//chainReference = contentMaster.GetContentInfo().ChainReference
			dstDirectory = v.ContentConfig.ApplicationStorageDirectory
			chainReference = "0"
		} else {
			//Adding parent path#################################################################
			conditionString := " object_info->>'$.chainReference' like '%" + strconv.Itoa(recordId) + "' and object_info->>'$.isFile' = 'false' "
			listOfObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)
			parentContenetMasterInfo := ContenetMasterInfo{}
			if listOfObjects == nil {
				v.BaseService.Logger.Info("Parent folder doesn't exist")
				response.DispatchDetailedError(ctx, FolderExistError,
					&response.DetailedError{
						Header:      "Incorrect reference Id",
						Description: "Can't create folder under given parent directory",
					})
				return
			} else {

				json.Unmarshal((*listOfObjects)[0].ObjectInfo, &parentContenetMasterInfo)
				chainReference = parentContenetMasterInfo.ChainReference
			}

			dstDirectory = parentContenetMasterInfo.Path
		}

	} else {
		// this is not sent in the URL
		dstDirectory = v.ContentConfig.ApplicationStorageDirectory
		chainReference = "0:1"
	}

	dstPath := path.Join(dstDirectory, contentId)
	v.BaseService.Logger.Info("destination path", zap.String("path", dstPath))
	userId := common.GetUserId(ctx)

	err = ctx.SaveUploadedFile(uploadFile, dstPath)
	mimetype.SetLimit(0)
	mimeType, err := mimetype.DetectFile(dstPath)

	var convertedFileType = "file"
	if err == nil {
		convertedFileType = getMIMEType(mimeType.String())
	}

	v.BaseService.Logger.Info("extracted MIME Type", zap.String("type", mimeType.String()))
	fileTypeIcon := getFileTypeIcon(convertedFileType)

	//now send the response
	contentInfo := ContenetMasterInfo{
		Name:             uploadFile.Filename,
		Size:             util.ByteCountSI(uploadFile.Size),
		Path:             dstPath,
		Share:            make([]SharePermission, 0),
		IsFile:           true,
		ChainReference:   chainReference,
		MIMEType:         convertedFileType,
		CreatedAt:        util.GetCurrentTime(TimeLayout),
		CreatedBy:        userId,
		LastUpdatedAt:    util.GetCurrentTime(TimeLayout),
		LastUpdatedBy:    userId,
		ObjectStatus:     common.ObjectStatusActive,
		FilePreviewImage: v.ContentConfig.DefaultPreviewUrl,
		FileTypeIcon:     fileTypeIcon,
	}

	serializedContentMasterInfo, _ := json.Marshal(contentInfo)
	contentMaster := component.GeneralObject{
		ObjectInfo: serializedContentMasterInfo,
	}

	err, recordId := Create(dbConnection, ContentMasterTable, contentMaster)
	var updatingData = make(map[string]interface{})

	url := v.ContentConfig.DomainUrl + "/project/" + projectId + "/content/component/content_master/record/" + strconv.Itoa(recordId)
	contentInfo.Url = url
	contentInfo.ChainReference += ":" + strconv.Itoa(recordId)
	updatingData["object_info"] = contentInfo.Serialize()
	Update(dbConnection, ContentMasterTable, recordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("Error creating content object in database", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, getError("Error creating file in the system"), ErrorInsertingNewContentRecord)
		return
	}

	//userInfo := authService.GetUserInfoById(userId)

	//resourceInfo := common.ResourceInfo{}
	//recordResource := common.RecordResource{}
	//
	//recordResource.UserId = userId
	//recordResource.Username = userInfo.Username
	//recordResource.CreatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	//recordResource.UpdatedAt = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	//
	//resourceInfo.ResourceMeta = recordResource
	//resourceInfo.Message = "Content " + uploadFile.Filename + " is created"
	//
	//recordMessageService := services.GetService("record_message").ServiceInterface.(services.RecordMessageInterface)
	//err = recordMessageService.CreateUserRecordMessage(projectId, strconv.Itoa(recordId), componentName, common.MessageTypeNotification, resourceInfo)
	//
	//if err != nil {
	//	cs.BaseService.Logger.Error("Error creating record message ", err.Error())
	//}
	//
	//rawInfo, _ := json.Marshal(contentInfo)
	//recordMessage := componentName + " is created"
	//err = CreateUserRecordTrail(cs, projectId, strconv.Itoa(recordId), componentName, recordMessage, &component.GeneralObject{Id: recordId, ObjectInfo: rawInfo}, &component.GeneralObject{Id: recordId, ObjectInfo: rawInfo})
	//if err != nil {
	//	cs.BaseService.Logger.Error("error in create record trail", zap.String("error", err.Error()))
	//}

	ctx.JSON(http.StatusOK, contentInfo)

}

// createNewResource ShowAccount godoc
// @Summary create new resource
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   recordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/content/component/{componentName}/records  [post]
func (v *ContentService) createMultipleNewResource(ctx *gin.Context) {
	var err error
	projectId := util.GetProjectId(ctx)
	form, _ := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Set the maximum memory allocated to parse multipart form data (5MB)
	// and the maximum file size allowed (5MB)
	if err = ctx.Request.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("uploaded file size is more than 5MB"), ErrorInsertingNewContentRecord)
		return
	}

	uploadFiles := form.File["file"]
	//uploadFile, err := ctx.FormFile("file")
	if err != nil {
		v.BaseService.Logger.Error("file parameter is not found in the request", zap.Any("ctx", ctx))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("extracting uploaded file, required parameter (file) not available, check the API specifications"), RequiredFileParamNotAvailable)
		return
	}
	dbConnection := v.BaseService.GetDatabase(projectId)
	referencePathId := ctx.Query("referencePathId")
	var dstDirectory string
	contentId := uuid.New().String()
	var chainReference string
	chainReference = "0"
	listOfContentInfo := make([]ContenetMasterInfo, 0)
	for _, uploadFile := range uploadFiles {
		if referencePathId != "" {
			// get the folder path
			recordId := util.InterfaceToInt(referencePathId)
			if recordId == 0 {
				// upload to main directory
				//err, generalContentObject := Get(dbConnection, ContentMasterTable, recordId)
				//if err != nil {
				//	response.DispatchDetailedError(ctx, FolderExistError,
				//		&response.DetailedError{
				//			Header:      "Folder doesn't exist",
				//			Description: "Invalid folder structure, please check the reference folder",
				//		})
				//	return
				//}
				//contentMaster := ContentMaster{ObjectInfo: generalContentObject.ObjectInfo}
				//dstDirectory = contentMaster.GetContentInfo().Path
				//chainReference = contentMaster.GetContentInfo().ChainReference
				dstDirectory = v.ContentConfig.ApplicationStorageDirectory
				chainReference = "0"
			} else {
				//Adding parent path#################################################################
				conditionString := " object_info->>'$.chainReference' like '%" + strconv.Itoa(recordId) + "' and object_info->>'$.isFile' = 'false' "
				listOfObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)
				parentContenetMasterInfo := ContenetMasterInfo{}
				if listOfObjects == nil {
					v.BaseService.Logger.Info("Parent folder doesn't exist")
					response.DispatchDetailedError(ctx, FolderExistError,
						&response.DetailedError{
							Header:      "Incorrect reference Id",
							Description: "Can't create folder under given parent directory",
						})
					return
				} else {

					json.Unmarshal((*listOfObjects)[0].ObjectInfo, &parentContenetMasterInfo)
					chainReference = parentContenetMasterInfo.ChainReference
				}

				dstDirectory = parentContenetMasterInfo.Path
			}

		} else {
			// this is not sent in the URL
			dstDirectory = v.ContentConfig.ApplicationStorageDirectory
			chainReference = "0:1"
		}

		dstPath := path.Join(dstDirectory, contentId)
		v.BaseService.Logger.Info("destination path", zap.String("path", dstPath))
		userId := common.GetUserId(ctx)

		err = ctx.SaveUploadedFile(uploadFile, dstPath)
		mimetype.SetLimit(0)
		mimeType, err := mimetype.DetectFile(dstPath)

		var convertedFileType = "file"
		if err == nil {
			convertedFileType = getMIMEType(mimeType.String())
		}

		v.BaseService.Logger.Info("extracted MIME Type", zap.String("type", mimeType.String()))
		fileTypeIcon := getFileTypeIcon(convertedFileType)

		//now send the response
		contentInfo := ContenetMasterInfo{
			Name:             uploadFile.Filename,
			Size:             util.ByteCountSI(uploadFile.Size),
			Path:             dstPath,
			Share:            make([]SharePermission, 0),
			IsFile:           true,
			ChainReference:   chainReference,
			MIMEType:         convertedFileType,
			CreatedAt:        util.GetCurrentTime(TimeLayout),
			CreatedBy:        userId,
			LastUpdatedAt:    util.GetCurrentTime(TimeLayout),
			LastUpdatedBy:    userId,
			ObjectStatus:     common.ObjectStatusActive,
			FilePreviewImage: v.ContentConfig.DefaultPreviewUrl,
			FileTypeIcon:     fileTypeIcon,
		}

		serializedContentMasterInfo, _ := json.Marshal(contentInfo)
		contentMaster := component.GeneralObject{
			ObjectInfo: serializedContentMasterInfo,
		}

		err, recordId := Create(dbConnection, ContentMasterTable, contentMaster)
		var updatingData = make(map[string]interface{})

		url := v.ContentConfig.DomainUrl + "/project/" + projectId + "/content/component/content_master/record/" + strconv.Itoa(recordId)
		contentInfo.Url = url
		contentInfo.ChainReference += ":" + strconv.Itoa(recordId)
		updatingData["object_info"] = contentInfo.Serialize()
		Update(dbConnection, ContentMasterTable, recordId, updatingData)
		if err != nil {
			v.BaseService.Logger.Error("Error creating content object in database", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, getError("Error creating file in the system"), ErrorInsertingNewContentRecord)
			return
		}

		listOfContentInfo = append(listOfContentInfo, contentInfo)
	}

	ctx.JSON(http.StatusOK, listOfContentInfo)

}

func (v *ContentService) getChildObjects(ctx *gin.Context) {
	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	componentName := ctx.Param("componentName")
	outFields := ctx.Query("out_fields")
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//'92:[0-9]+$'
	regularExpression := recordId + ":[0-9]+$"

	conditionString := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) REGEXP '" + regularExpression + "' "
	listOfObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)
	totalRecords := int64(len(*listOfObjects))

	_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
	ctx.JSON(http.StatusOK, tableRecordsResponse)

}
func (v *ContentService) getObjects(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	format := ctx.Query("format")
	searchFields := ctx.Query("search")
	objectStatus := ctx.Query("objectStatus")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	regularExpression := "0:[0-9]+$"

	rootDirectoryCondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) REGEXP '" + regularExpression + "' "
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	statusQuery := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.objectStatus\")) = '" + objectStatus + "'"
	if searchFields != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)
		baseCondition := component.TableCondition(offsetValue, fields, values, condition)
		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		searchWithBaseQuery := searchQuery + " AND " + baseCondition + " AND " + rootDirectoryCondition
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		if objectStatus == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, rootDirectoryCondition)
			totalRecords = int64(len(*listOfObjects))
		} else {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, statusQuery+" AND "+rootDirectoryCondition)
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			if objectStatus == "" {
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))
			} else {
				baseCondition := component.TableCondition(offsetValue, fields, values, condition)
				listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, baseCondition+" AND "+statusQuery+" AND "+rootDirectoryCondition)
			}

		} else {
			baseCondition := component.TableCondition(offsetValue, fields, values, condition) + " AND " + rootDirectoryCondition
			if objectStatus != "" {
				baseCondition += " AND " + statusQuery
			}
			limitVal, _ := strconv.Atoi(limitValue)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, baseCondition, limitVal)
		}

	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *ContentService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordIdString(ctx)
	recordIdInt, _ := strconv.Atoi(recordId)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, recordIdInt)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	contentMasterInfo := ContenetMasterInfo{}
	json.Unmarshal(generalObject.ObjectInfo, &contentMasterInfo)

	if contentMasterInfo.IsFile {
		// check constraints and proceed
		listOfConstraints := v.ComponentManager.GetConstraints(componentName)
		for _, constraint := range listOfConstraints {

			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)

			ArchiveReferenceObjects(dbConnection, referenceTable, referenceField, recordIdInt)
		}
		v.ComponentManager.ProcessDeleteDependencyInjection(dbConnection, recordIdInt, componentName)
		updatedObjectInfo := common.UpdateMetaInfoFromSerializedObject(generalObject.ObjectInfo, ctx)
		err = ArchiveObject(dbConnection, targetTable, component.GeneralObject{Id: generalObject.Id, ObjectInfo: updatedObjectInfo})

		if err != nil {
			v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
			return
		}

		// recordMessage := componentName + " is deleted"
		// err = CreateUserRecordTrail(cs, projectId, recordId, componentName, recordMessage, nil, nil)
		if err != nil {
			v.BaseService.Logger.Error("error in create record trail", zap.String("error", err.Error()))
		}
	} else {
		chainReference := contentMasterInfo.ChainReference
		referenceList := strings.Split(chainReference, ":")
		conditionString := " object_info ->>'$.chainReference' like '%" + referenceList[len(referenceList)-1] + "%' "
		listOfObjects, _ := GetConditionalObjects(dbConnection, ContentMasterTable, conditionString)

		for _, contentObject := range *listOfObjects {
			recordIdInt = contentObject.Id
			// check constraints and proceed
			listOfConstraints := v.ComponentManager.GetConstraints(componentName)
			for _, constraint := range listOfConstraints {

				referenceComponent := constraint.Reference
				referenceField := constraint.ReferenceProperty
				referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)

				ArchiveReferenceObjects(dbConnection, referenceTable, referenceField, recordIdInt)
			}
			v.ComponentManager.ProcessDeleteDependencyInjection(dbConnection, recordIdInt, componentName)
			updatedObjectInfo := common.UpdateMetaInfoFromSerializedObject(contentObject.ObjectInfo, ctx)
			err = ArchiveObject(dbConnection, targetTable, component.GeneralObject{Id: contentObject.Id, ObjectInfo: updatedObjectInfo})

			if err != nil {
				v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
				response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
				return
			}

			if err != nil {
				v.BaseService.Logger.Error("error in create record trail", zap.String("error", err.Error()))
			}
		}

		ctx.Status(http.StatusNoContent)
	}

}

// getSearchResults ShowAccount godoc
// @Summary Get the search results based on given input
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param SearchField body SearchKeys true "Pass the array of key and values"
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/search [post]
func (v *ContentService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	var totalRecords int64
	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	if len(searchFieldCommand) == 0 {
		// reset the search
		listOfObjects, err := GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
		err, tableObjectResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, tableObjectResponse)
		return
	}

	format := ctx.Query("format")
	searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
	listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)
}

// getCardView ShowAccount godoc
// @Summary Get all the machine information in a card view
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/card_view [get]
func (v *ContentService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	regularExpression := "0:[0-9]+$"
	rootDirectorycondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) REGEXP '" + regularExpression + "' "

	if searchFields != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchQuery+" AND "+rootDirectorycondition, limitVal)
	} else {
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, rootDirectorycondition)
	}

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	// cardViewResponse := cs.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
	var dd []GroupByCardView

	for _, object := range *listOfObjects {
		var contentMasterInfo map[string]interface{}
		json.Unmarshal(object.ObjectInfo, &contentMasterInfo)
		if util.InterfaceToString(contentMasterInfo["objectStatus"]) != "Archived" {
			if reflect.ValueOf(contentMasterInfo["isFile"]).Bool() {
				idx := searchGroupField(dd, "File")

				if idx != -1 {

					contentMasterInfo["id"] = object.Id

					dd[idx].Cards = append(dd[idx].Cards, contentMasterInfo)
				} else {
					newCard := make([]map[string]interface{}, 0)

					contentMasterInfo["id"] = object.Id

					newCard = append(newCard, contentMasterInfo)
					fileGroupByCardView := GroupByCardView{GroupByField: "File", Cards: newCard}

					dd = append(dd, fileGroupByCardView)
				}
			} else {
				idx := searchGroupField(dd, "Folder")
				if idx != -1 {

					contentMasterInfo["id"] = object.Id

					dd[idx].Cards = append(dd[idx].Cards, contentMasterInfo)
				} else {
					newCard := make([]map[string]interface{}, 0)

					contentMasterInfo["id"] = object.Id

					newCard = append(newCard, contentMasterInfo)
					fileGroupByCardView := GroupByCardView{GroupByField: "Folder", Cards: newCard}

					dd = append(dd, fileGroupByCardView)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, dd)

}

// getCardView ShowAccount godoc
// @Summary Get all the machine information in a card view
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/card_view [get]
func (v *ContentService) getChildCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := util.GetProjectId(ctx)
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	searchFields := ctx.Query("search")
	recordId := ctx.Param("recordId")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error

	regularExpression := recordId + ":[0-9]+$"
	rootDirectorycondition := " JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$.chainReference\")) REGEXP '" + regularExpression + "' "

	if searchFields != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		listOfSearchFields := strings.Split(searchFields, ",")
		var searchFieldCommand []component.SearchKeys
		for _, searchFieldObject := range listOfSearchFields {
			keyValueObject := strings.Split(searchFieldObject, ":")
			searchFieldCommand = append(searchFieldCommand, component.SearchKeys{Field: keyValueObject[0], Value: keyValueObject[1]})
		}
		searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchQuery+" AND "+rootDirectorycondition, limitVal)
	} else {
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, rootDirectorycondition)
	}

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	// cardViewResponse := cs.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
	var dd []GroupByCardView

	for _, object := range *listOfObjects {
		var contentMasterInfo map[string]interface{}
		json.Unmarshal(object.ObjectInfo, &contentMasterInfo)
		if util.InterfaceToString(contentMasterInfo["objectStatus"]) != "Archived" {
			if reflect.ValueOf(contentMasterInfo["isFile"]).Bool() {
				idx := searchGroupField(dd, "File")

				if idx != -1 {

					contentMasterInfo["id"] = object.Id

					dd[idx].Cards = append(dd[idx].Cards, contentMasterInfo)
				} else {
					newCard := make([]map[string]interface{}, 0)

					contentMasterInfo["id"] = object.Id

					newCard = append(newCard, contentMasterInfo)
					fileGroupByCardView := GroupByCardView{GroupByField: "File", Cards: newCard}

					dd = append(dd, fileGroupByCardView)
				}
			} else {
				idx := searchGroupField(dd, "Folder")
				if idx != -1 {

					contentMasterInfo["id"] = object.Id

					dd[idx].Cards = append(dd[idx].Cards, contentMasterInfo)
				} else {
					newCard := make([]map[string]interface{}, 0)

					contentMasterInfo["id"] = object.Id

					newCard = append(newCard, contentMasterInfo)
					fileGroupByCardView := GroupByCardView{GroupByField: "Folder", Cards: newCard}

					dd = append(dd, fileGroupByCardView)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, dd)

}

func (v *ContentService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	err, objectInterface := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      getError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})

	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	err = Update(dbConnection, targetTable, intRecordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("error updating content master information", zap.String("error", err.Error()))
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Updating Resource Failed"), ErrorUpdatingObjectInformation, "Failed to update resource due to internal system error")
		return
	}

}

func searchGroupField(dd []GroupByCardView, serachField string) int {
	for idx, v := range dd {
		if v.GroupByField == serachField {
			return idx
		}
	}
	return -1
}
