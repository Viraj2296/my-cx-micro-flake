package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func (v *MachineService) getFilterObjects(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")
	var payloadObject common.FilterCriteria
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	orderValue := ctx.Query("order")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	format := ctx.Query("format")
	searchFields := ctx.Query("search")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//Have to next flag
	isNext := true
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error

	if err := ctx.ShouldBindBodyWith(&payloadObject, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

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
		searchWithBaseQuery := searchQuery + " AND " + baseCondition
		listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
	} else {
		if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
			statusCondition := targetTable + ".object_info ->> '$.hmiStatus' = 'started' OR " + targetTable + ".object_info ->> '$.hmiStatus' = 'stopped'"
			totalRecords = v.ComponentManager.GetFilterCount(dbConnection, targetTable, payloadObject, statusCondition)
		} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
			statusCondition := targetTable + ".object_info ->> '$.rejectedQuantity' > 0 "
			totalRecords = v.ComponentManager.GetFilterCount(dbConnection, targetTable, payloadObject, statusCondition)
		} else {
			totalRecords = v.ComponentManager.GetFilterCount(dbConnection, targetTable, payloadObject, "")
		}

		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			queryCondition := OffsetCondition(offsetValue, targetTable)

			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				offsetValue = strconv.Itoa(int(totalRecords) - limitVal + 1)

				limitVal = limitVal - offsetVal
				queryCondition = OffsetCondition(offsetValue, targetTable)
				if componentName == MachineHMITable || componentName == AssemblyMachineHmiTable || componentName == ToolingMachineHmiTable {
					queryCondition = queryCondition + " AND (object_info ->> '$.hmiStatus' = 'started' OR object_info ->> '$.hmiStatus' = 'stopped')"

				} else if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
					queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

				}
				orderBy := componentName + ".id desc"
				var idListCondition = v.getFilterQueryCondition(dbConnection, payloadObject, queryCondition, limitVal, targetTable)
				if idListCondition != "()" {
					listOfObjects, err = GetConditionalObjectsOrderBy(dbConnection, targetTable, idListCondition, orderBy, limitVal)
				} else {
					listOfObject := make([]component.GeneralObject, 0)
					listOfObjects = &listOfObject
				}

				var currentRecordCount = 0
				if listOfObjects != nil {
					currentRecordCount = len(*listOfObjects)
				}

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(queryCondition, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

					} else {
						totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
					}

					if listOfObjects != nil {
						if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[0].Id {
							isNext = false
						}
					} else {
						isNext = false
					}

				}

			} else {
				if componentName == MachineHMIRejectedComponent || componentName == AssemblyMachineHMIRejectedComponent || componentName == ToolingMachineHMIRejectedComponent {
					queryCondition = queryCondition + " AND (object_info ->> '$.rejectedQuantity' > 0)"

				}

				var idListCondition = v.getFilterQueryCondition(dbConnection, payloadObject, queryCondition, limitVal, targetTable)
				if idListCondition != "()" {
					listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, idListCondition, limitVal)
				} else {
					listOfObject := make([]component.GeneralObject, 0)
					listOfObjects = &listOfObject
				}

				var currentRecordCount = 0
				if listOfObjects != nil {
					currentRecordCount = len(*listOfObjects)
				}

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(queryCondition, "AND")
					var totalRecordObjects *[]component.GeneralObject
					if len(andClauses) > 1 {
						totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, queryCondition)

					} else {
						totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
					}
					lenTotalRecord := len(*totalRecordObjects)

					if listOfObjects != nil {
						if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
							isNext = false
						}
					} else {
						isNext = false
					}

				}
			}
			//limitVal, _ := strconv.Atoi(limitValue)
			//listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)
		}

	}
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		userId := common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)

		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		if listOfObjects == nil {
			tableObjectResponse.Data = make([]datatypes.JSON, 0)
		}

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)
		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func OffsetCondition(offset, componentName string) string {

	if offset == "" {
		offset = "0"
	}

	if offset == "-1" {
		return ""
	} else {
		return componentName + ".id > " + offset
	}

}

func (v *MachineService) getFilterQueryCondition(dbConnection *gorm.DB, searchFieldCommand common.FilterCriteria, offsetCondition string, limitValue int, componentName string) string {
	limitVal := strconv.Itoa(limitValue)
	searchList := v.ComponentManager.ExecuteFilter(dbConnection, componentName, searchFieldCommand, offsetCondition, limitVal)
	fmt.Println("searchList: ", searchList)
	if searchList != "()" {
		searchQuery := " id in " + searchList
		fmt.Println("searchQuery: ", searchQuery)
		return searchQuery

	}

	return searchList
}
