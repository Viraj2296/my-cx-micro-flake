package header_parser

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetMiddlewareUserId(ctx *gin.Context) int {
	value, isExist := ctx.Get("id")
	var userRecordId int
	if isExist {
		userRecordId = value.(int)
	} else {
		userRecordId = 0
	}
	return userRecordId
}

func GetMiddlewareProjectId(ctx *gin.Context) string {
	value, isExist := ctx.Get("projectId")
	var projectId string
	if isExist {
		projectId = value.(string)
	} else {
		projectId = ""
	}
	return projectId
}

func GetMiddlewareComponentName(ctx *gin.Context) string {
	value, isExist := ctx.Get("componentName")
	var componentName string
	if isExist {
		componentName = value.(string)
	} else {
		componentName = ""
	}
	return componentName
}

func GetMiddlewareRecordId(ctx *gin.Context) int {
	value, isExist := ctx.Get("id")
	var recordId int
	if isExist {
		recordId = value.(int)
	} else {
		recordId = 0
	}
	return recordId
}
func GetQueryField(ctx *gin.Context, key string) string {
	return ctx.Query(key)
}

func GetComponentName(ctx *gin.Context) string {
	return ctx.Param("componentName")
}

func GetActionName(ctx *gin.Context) string {
	return ctx.Param("actionName")
}

func GetModuleName(ctx *gin.Context) string {
	return ctx.Param("moduleName")
}

func GetToken(ctx *gin.Context) string {
	return ctx.Param("token")
}

func GetRecordId(ctx *gin.Context) int {
	recordId := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordId)
	return intRecordId
}

func GetLevel1RecordId(ctx *gin.Context) int {
	recordId := ctx.Param("level1RecordId")
	intRecordId, _ := strconv.Atoi(recordId)
	return intRecordId
}

func GetOriginalRecordId(ctx *gin.Context) string {
	return ctx.Param("recordId")
}

func GetMachineId(ctx *gin.Context) int {
	id := ctx.Param("machineId")
	intRecordId, _ := strconv.Atoi(id)
	return intRecordId
}

func GetEventId(ctx *gin.Context) int {
	id := ctx.Param("eventId")
	intRecordId, _ := strconv.Atoi(id)
	return intRecordId
}

func GetRecordIdString(ctx *gin.Context) string {
	recordId := ctx.Param("recordId")
	return recordId
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
