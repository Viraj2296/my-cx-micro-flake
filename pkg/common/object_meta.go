package common

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type MetaInfo struct {
	Created   string `json:"created"`
	CreatedBy int    `json:"createdBy"`
	Updated   string `json:"updated"`
	UpdatedBy int    `json:"updatedBy"`
}

func (metaInfo *MetaInfo) Serialize() []byte {
	rawData, _ := json.Marshal(metaInfo)
	return rawData
}

func GetUserId(ctx *gin.Context) int {
	value, isExist := ctx.Get("id")
	var userRecordId int
	if isExist {
		userRecordId = value.(int)
	} else {
		userRecordId = 0
	}
	return userRecordId
}

func CreateMetaInfo(ctx *gin.Context) MetaInfo {
	metaInfo := MetaInfo{}

	metaInfo.CreatedBy = GetUserId(ctx)
	metaInfo.Created = time.Now().String()
	metaInfo.Updated = time.Now().String()
	metaInfo.UpdatedBy = GetUserId(ctx)
	return metaInfo
}

func (metaInfo *MetaInfo) UpdateMetaInfo(ctx *gin.Context) {
	userId := GetUserId(ctx)
	metaInfo.Updated = time.Now().String()
	metaInfo.UpdatedBy = userId
}
