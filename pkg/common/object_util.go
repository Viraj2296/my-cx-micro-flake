package common

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func AddFieldJSONObject(jsonObject datatypes.JSON, fieldName string, value interface{}) datatypes.JSON {
	var data = make(map[string]interface{})
	json.Unmarshal(jsonObject, &data)
	data[fieldName] = value
	rawModifiedJSON, _ := json.Marshal(data)
	return rawModifiedJSON
}

func InitMetaInfoFromSerializedObject(jsonObject datatypes.JSON, ctx *gin.Context) datatypes.JSON {
	var data = make(map[string]interface{})
	userId := GetUserId(ctx)
	json.Unmarshal(jsonObject, &data)
	data["objectStatus"] = ObjectStatusActive
	data["createdAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	data["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	data["createdBy"] = userId
	data["lastUpdatedBy"] = userId
	rawModifiedJSON, _ := json.Marshal(data)
	return rawModifiedJSON
}

func InitMetaInfoFromInterface(object map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	userId := GetUserId(ctx)
	object["objectStatus"] = ObjectStatusActive
	object["createdAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	object["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	object["createdBy"] = userId
	object["lastUpdatedBy"] = userId
	return object
}

func ValidateObjectStatus(jsonObject datatypes.JSON) bool {
	var data = make(map[string]interface{})
	json.Unmarshal(jsonObject, &data)
	if objectStatus, ok := data["objectStatus"]; ok {
		if objectStatus == ObjectStatusActive {
			return true
		} else {
			return false
		}
	}
	// no object status field, so let the modification happen, but this is not good, change it later
	return true
}

func UpdateMetaInfoFromSerializedObject(jsonObject datatypes.JSON, ctx *gin.Context) datatypes.JSON {
	var data = make(map[string]interface{})
	userId := GetUserId(ctx)
	json.Unmarshal(jsonObject, &data)
	data["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	data["lastUpdatedBy"] = userId
	rawModifiedJSON, _ := json.Marshal(data)
	return rawModifiedJSON
}

func UpdateMetaInfoFromInterface(object map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	userId := GetUserId(ctx)
	object["lastUpdatedAt"] = util.GetCurrentTime("2006-01-02T15:04:05.000Z")
	object["lastUpdatedBy"] = userId
	return object
}
