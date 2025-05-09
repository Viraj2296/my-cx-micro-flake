package handler

import "errors"

const (
	CreateUnknownObjectType = "unknown object type in creating object in database"
	GetUnknownObjectType    = "unknown object type in getting object in database"
	DeleteUnknownObjectType = "unknown object type in deleting object in database"
	UpdateUnknownObjectType = "unknown object type in updating object in database"

	ProjectRecordTrailTable = "project_record_trail"
	ProjectTable            = "project"
	ProjectComponentTable   = "project_component"

	ProjectID = "906d0fd569404c59956503985b330132"

	ErrorGettingObjectsInformation          = 5000
	ErrorUpdatingObjectInformation          = 5005
	ErrorRemovingObjectInformation          = 5006
	ErrorGettingIndividualObjectInformation = 5007
	ErrorCreatingObjectInformation          = 5008

	ModuleName = "project"
)

func getError(errorString string) error {
	return errors.New(errorString)
}
