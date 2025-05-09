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

	ManufacturingRecordTrailTable  = "manufacturing_record_trail"
	ManufacturingComponentTable    = "manufacturing_component"
	ManufacturingMouldingPartTable = "manufacturing_moulding_part"
	ManufacturingAssemblyPartTable = "manufacturing_assembly_part"
	ManufacturingVendorMasterTable = "manufacturing_vendor_master"

	InvalidSourceError = "Invalid Source"

	InvalidSchedulePosition = "Invalid Schedule Position"
	InvalidComponent        = 6010

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	DecodingFailed       = 6070
	InvalidMachineStatus = 6071

	ProjectID  = "906d0fd569404c59956503985b330132"
	TimeLayout = "2006-01-02T15:04:05.000Z"

	ModuleName = "manufacturing"
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
