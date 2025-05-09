package handler

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/services/it/handler/const_util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func SendInvalidComponentSchema(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, const_util.InvalidSchema,
		&response.DetailedError{
			Header:      "Invalid Component Schema",
			Description: "The requested component is missing elements which is required to render. Please report to developers",
		})
	return
}

func SendResourceNotFound(ctx *gin.Context) {
	response.DispatchDetailedError(ctx, const_util.ErrorGettingObjectsInformation,
		&response.DetailedError{
			Header:      "Invalid Resource",
			Description: "The resource that system is trying process not found, it should be due to either other process deleted it before it access or not created yet",
		})
	return
}

func GetObjectFields(serialised datatypes.JSON) map[string]interface{} {
	var objectFields = make(map[string]interface{})
	json.Unmarshal(serialised, &objectFields)
	return objectFields
}

func GetInterfaceToSerialisation(objectInterface interface{}) datatypes.JSON {
	serialisedObject, _ := json.Marshal(objectInterface)
	return serialisedObject
}

func MergeObjects(src map[string]interface{}, dst map[string]interface{}) map[string]interface{} {
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func AddObjectField(serialisedData datatypes.JSON, key string, value interface{}) map[string]interface{} {
	var objectFields = make(map[string]interface{})
	json.Unmarshal(serialisedData, &objectFields)
	objectFields[key] = value
	return objectFields
}
func IsEmpty(src string) bool {
	if src == "" {
		return true
	}
	return false
}

func reverseSlice(generalObjects *[]component.GeneralObject) *[]component.GeneralObject {
	if generalObjects != nil {
		for i, j := 0, len(*generalObjects)-1; i < j; i, j = i+1, j-1 {
			(*generalObjects)[i], (*generalObjects)[j] = (*generalObjects)[j], (*generalObjects)[i]
		}
	}

	return generalObjects
}
