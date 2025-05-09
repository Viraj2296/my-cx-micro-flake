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

	FixedAssetRecordTrailTable = "fixed_asset_record_trail"

	FixedAssetsComponentTable  = "fixed_assets_component"
	FixedAssetMasterTable      = "fixed_asset_master"
	FixedAssetContractTable    = "fixed_asset_contract"
	FixedAssetSetupMasterTable = "fixed_asset_setup_master"

	FixedAssetDynamicFieldTable              = "fixed_asset_dynamic_field"
	FixedAssetDynamicFieldConfigurationTable = "fixed_asset_dynamic_field_configuration"

	FixedAssetDisposalComponent = "fixed_asset_disposal"

	FixedAssetMyDisposalComponent = "fixed_asset_my_disposal"
	FixedAssetMyTransferComponent = "fixed_asset_my_transfer"

	FixedAssetSetupMasterComponent = "fixed_asset_setup_master"

	FixedAssetTransferComponent = "fixed_asset_transfer"

	FixedAssetDisposalTable = "fixed_asset_disposal"
	FixedAssetTransferTable = "fixed_asset_transfer"
	FixedAssetLocationTable = "fixed_asset_location"

	FixedAssetStatusTable         = "fixed_asset_status"
	FixedAssetContractStatusTable = "fixed_asset_contract_status"
	FixedAssetTransferStatusTable = "fixed_asset_transfer_status"
	FixedAssetDisposalStatusTable = "fixed_asset_disposal_status"
	FixedAssetCategoryTable       = "fixed_asset_category"
	FixedAssetSubCategoryTable    = "fixed_asset_sub_category"
	FixedAssetClassTable          = "fixed_asset_class"

	FixedAssetEmailTemplateTable      = "fixed_asset_email_template"
	FixedAssetEmailTemplateFieldTable = "fixed_asset_email_template_field"

	ITAssetCategoryTable = "it_asset_category"
	AssetClassTable      = "fixed_asset_class"

	FixedAssetMasterComponent   = "fixed_asset_master"
	FixedAssetContractComponent = "fixed_asset_contract"
	// AssetMasterComponent      = "fixed_asset_master"

	FixedAssetCreationEmailTemplate         = 1
	FixedAssetContractCreationEmailTemplate = 2
	FixedAssetDisposalEmailTemplate         = 3
	FixedAssetTransferEmailTemplate         = 4

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	DisposalStatusCreated     = 1
	DisposalStatusSubmitted   = 2
	DisposalStatusHODApproved = 3
	DisposalStatusHODRejected = 4
	DisposalStatusCEOApproved = 5
	DisposalStatusCEORejected = 6
	DisposalStatusCEONotified = 7

	DisposalUserWorkflowLevel = 1
	DisposalHODWorkflowLevel  = 2
	DisposalCEOWorkflowLevel  = 3

	TransferStatusCreated     = 1
	TransferStatusSubmitted   = 2
	TransferStatusHODApproved = 3
	TransferStatusHODRejected = 4
	TransferStatusCEOApproved = 5
	TransferStatusCEORejected = 6

	TransferUserWorkflowLevel = 1
	TransferHODWorkflowLevel  = 2
	TransferCEOWorkflowLevel  = 3

	AssetNumberPrefix = "FUYU_"
	ISOTimeLayout     = "2006-01-02T15:04:05.000Z"

	ProjectID = "906d0fd569404c59956503985b330132"

	ScheduleStatusPreferenceThree = 3
	ScheduleStatusPreferenceFour  = 4

	ModuleName = "fixed_assets"
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
