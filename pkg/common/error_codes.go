package common

const (
	InvalidTokenOrTokenMissing                = 900
	ExtractingUserIdHadFailedFromToken        = 901
	UnableToCreateExportFile                  = 905
	UnableToReadCSVFile                       = 906
	FailedToDownloadTheImportFileUrl          = 907
	ParsingCSVFileFailed                      = 908
	SchemaIsNotMatchedWithUploadedCSV         = 909
	AccessDenied                              = 910
	InterModuleCommunicationProblem           = 911
	AccessDeniedDescriptionForGetResources    = "You are not authorised to access the requested resources, permission is not given to explore the resources, please check the permission level, and access principal assigned to you by admin"
	AccessDeniedDescriptionForUpdateResources = "You are not authorised to update the requested resources, permission is not given to update the resources, please check the permission level, and access principal assigned to you by admin"
	AccessDeniedDescriptionForCreateResources = "You are not authorised to create new resources, permission is not granted to create the resources, please check the permission level, and access principal assigned to you by admin"
	AccessDeniedDescriptionForActionResources = "You are not authorised to do this action on this resource, permission is not granted, please check the permission level, and access principal assigned to you by admin"

	AccessDeniedDescriptionForDeleteResources = "You are not authorised to delete resources, permission is not granted to delete this resources, please check the permission level, and access principal assigned to you by admin"

	InvalidObjectStatus       = 911
	ResourceCreationFailed    = 913
	UpdateResourceFailed      = 914
	OperationNotPermitted     = 915
	UnmarshalingError         = 916
	ErrorCreatingResource     = "Error Creating New Resource"
	UpdateResourceFailedError = "Error Updating Resource"

	ObjectNotFound                       = 912
	InvalidObjectStatusError             = "Invalid Resource Status"
	OperationNotPermittedError           = "Invalid Operation"
	InterModuleCommunicationProblemError = "Internal Service Communication Failed"
	ObjectNotFoundError                  = "Requested Resource Not Found"
)
