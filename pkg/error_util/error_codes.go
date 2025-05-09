package error_util

import (
	"cx-micro-flake/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	InvalidTokenOrTokenMissing         = 900
	ExtractingUserIdHadFailedFromToken = 901
	UnableToCreateExportFile           = 905
	UnableToReadCSVFile                = 906
	FailedToDownloadTheImportFileUrl   = 907
	ParsingCSVFileFailed               = 908
	SchemaIsNotMatchedWithUploadedCSV  = 909
	AccessDenied                       = 910
	InterModuleCommunicationProblem    = 911
	UnmarshlingError                   = 911

	InvalidObjectStatus    = 911
	ResourceCreationFailed = 913
	UpdateResourceFailed   = 914

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008
	InvalidSchema                           = 5009
	InvalidQueryField                       = 5010
	QueryExecutionFailed                    = 5011
	ErrorInExportingObject                  = 5012
	ErrorInLoadingFile                      = 5013
	InternalSystemErrorCode                 = 5012

	GettingServiceFailed = 5013

	ConnectingToDatasourceFailed = 5014
	InvalidTargetTable           = 5015
	InvalidRequestParameter      = 5015
	ParentFolderDoesntExist      = 5016
	FolderCreationFailed         = 5017
	FolderExistError             = 5018
	ErrorResetTable              = 5019
	InvalidWorkflow              = 5020
	WorkflowActionNotDefined     = 5021
	WorkflowValidationFailed     = 5022
)

func GetError(errorString string) error {
	return errors.New(errorString)
}

func SendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorGettingObjectsInformation,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}
func SendInvalidWorkflow(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, InvalidWorkflow,
		&response.DetailedError{
			Header:      "Invalid Workflow",
			Description: "Invalid workflow execution. it should be due to either other workflow execution not defined, or execution in a different way that defined",
		})
	return
}
func SendWorkflowActionNotDefined(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, WorkflowActionNotDefined,
		&response.DetailedError{
			Header:      "Invalid Workflow Action",
			Description: "Invalid workflow execution. it should be due to either other workflow execution not defined, or execution in a different way that defined",
		})
	return
}

func SendResourceCreationFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorCreatingObjectInformation,
		&response.DetailedError{
			Header:      "Resource Creation Failed",
			Description: "The resource that system is trying create had failed due to internal system error",
		})
	return
}
func SendResourceUpdateFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorUpdatingObjectInformation,
		&response.DetailedError{
			Header:      "Resource Update Failed",
			Description: "The resource that system is trying update had failed due to internal system error",
		})
	return
}

func SendUnmarshlingFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorUpdatingObjectInformation,
		&response.DetailedError{
			Header:      "Invalid JSON",
			Description: "The system is unable to process the request. Please review your request body",
		})
	return
}
func SendArchiveFailed(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorRemovingObjectInformation,
		&response.DetailedError{
			Header:      "Archived Failed",
			Description: "Internal system error during archive process. This is normally happen when the system is not configured properly. Please report to system administrator",
		})
	return
}

func SendInvalidRequestParameter(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, InvalidSchema,
		&response.DetailedError{
			Header:      "Invalid Request Parameter",
			Description: "Request contains invalid request parameter, check the request API specifications for further information",
		})
	return

}
func SendInvalidComponentSchema(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, InvalidSchema,
		&response.DetailedError{
			Header:      "Invalid Component Schema",
			Description: "The requested component is missing elements which is required to render. Please report to developers",
		})
	return
}

func SendInternalSystemError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, InternalSystemErrorCode,
		&response.DetailedError{
			Header:      "Internal System Error",
			Description: "A malfunction has occurred within the API system. Please try again later or contact support if the issue persists",
		})
	return
}

func SendExportObjectError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorInExportingObject,
		&response.DetailedError{
			Header:      "Error Exporting Objects",
			Description: "A malfunction has occurred within the API system. Please try again later or contact support if the issue persists",
		})
	return
}

func SendLoadFileError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorInLoadingFile,
		&response.DetailedError{
			Header:      "Error Loading File",
			Description: "A malfunction has occurred within the API system. Please try again later or contact support if the issue persists",
		})
	return
}

func SendInvalidTargetTableError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, ErrorInLoadingFile,
		&response.DetailedError{
			Header:      "Internal Server Error",
			Description: "Invalid target table, it can happen when the component is not loaded properly or wrongly configured by the user",
		})
	return
}

func SendInvalidTokenError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, InvalidTokenOrTokenMissing,
		&response.DetailedError{
			Header:      "Invalid Token",
			Description: "Token is not valid, please make sure you are accessing the link sent from system",
		})
}

func SendDeserialiseError(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, UnmarshlingError,
		&response.DetailedError{
			Header:      "Invalid Object",
			Description: "System couldn't able deserialise objectm this will happen when the input is not correctly formated as expected",
		})
}
