package common

import (
	"context"
	"cx-micro-flake/pkg/common/analytics"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dimchansky/utfbom"
	"github.com/google/uuid"
	"github.com/hashicorp/go-getter"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ComponentManager struct {
	ComponentLinkedFields  map[int]LinkedValues              // load the linked fields during loading
	ComponentSchema        map[int]component.ComponentSchema // keep the component schema based on id for fast lookup, we don't need to keep loading component schema again and again
	ComponentNameIdMapping map[string]int
	ComponentContentConfig component.UpstreamContentConfig
	ComponentTables        []string
}

func (cm *ComponentManager) GetTargetTable(componentName string) string {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	return cm.ComponentSchema[int64ComponentId].TargetTable
}

func (cm *ComponentManager) GetConstraints(componentName string) []component.Constraints {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	return cm.ComponentSchema[int64ComponentId].Constraints
}

func (cm *ComponentManager) GetTableImportSchema(componentName string) component.TableImportSchema {
	intComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[intComponentId]
	return componentSchema.TableImportSchema
}
func (cm *ComponentManager) GetTableSchema(componentName string) []component.TableSchema {
	intComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[intComponentId]
	return componentSchema.TableSchema
}

func (cm *ComponentManager) GetRecordSchema(componentName string) []component.RecordSchema {
	intComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[intComponentId]
	return componentSchema.RecordSchema
}
func (cm *ComponentManager) GetLinkedObjectMapSchema(componentName string) []component.LinkedObjectMap {
	intComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[intComponentId]
	return componentSchema.LinkedObjectMap
}

func (cm *ComponentManager) GetTableExportSchema(componentName string) []component.ExportSchema {
	intComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[intComponentId]
	exportSchema := addIndexExportSchema(componentSchema.ExportSchema)
	return exportSchema
}

func (cm *ComponentManager) GetTableExportSchemaV2(componentName string, objects *[]component.GeneralObject) []component.ExportSchema {
	intComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[intComponentId]
	exportSchema := addIndexExportSchema(componentSchema.ExportSchema)
	return exportSchema
}

func addIndexExportSchema(exportSchema []component.ExportSchema) []component.ExportSchema {
	mainId := float32(1.0)
	for index, schema := range exportSchema {
		schema.Id = mainId
		if len(schema.Children) > 0 {
			subId := 1
			for childIndex, schemaChildren := range schema.Children {
				//Index look like 7.12
				stringIndex := strconv.Itoa(int(mainId)) + "." + strconv.Itoa(subId)
				childrenIndex, _ := strconv.ParseFloat(stringIndex, 2)
				schemaChildren.Id = float32(childrenIndex)
				schema.Children[childIndex] = schemaChildren

				subId += 1
			}

		}
		exportSchema[index] = schema
		mainId += 1

	}
	return exportSchema
}

func (cm *ComponentManager) ProcessInternalArrayReferenceRequest(request map[string]interface{}, recordInfo datatypes.JSON, componentName string) []byte {
	var objectFields map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[int64ComponentId]
	baseRecordSchema := componentSchema.RecordSchema
	json.Unmarshal(recordInfo, &objectFields)
	for _, fields := range baseRecordSchema {
		if val, ok := request[fields.Property]; ok {
			//do something here
			if fields.Type == "int_array" {
				removingData := request[fields.Property]
				updatedValue := util.RemoveFromIntArray(objectFields[fields.Property], removingData)
				request[fields.Property] = updatedValue
			} else if fields.Type == "array" {
				removingData := request[fields.Property]
				updatedValue := util.RemoveFromIntArray(objectFields[fields.Property], removingData)
				request[fields.Property] = updatedValue
			} else if fields.Type == "object_array" {
				removingDataObjects := request[fields.Property]
				var arrayOfRemovingData = removingDataObjects.([]interface{})
				var modifiedObjects []interface{}
				// we already have an array of objects to update
				var isRemove bool
				isRemove = false
				var arrayExistingObjects = objectFields[fields.Property].([]interface{})
				for _, existingObject := range arrayExistingObjects {

					var existingObjectFields = existingObject.(map[string]interface{})
					existingId := existingObjectFields["id"]
					for _, removingData := range arrayOfRemovingData {
						var removingObjectFields = removingData.(map[string]interface{})
						removingId := removingObjectFields["id"]
						if existingId == removingId {
							// we are removing this
							isRemove = true
						}
					}
					if !isRemove {
						modifiedObjects = append(modifiedObjects, existingObject)
					}
					isRemove = false
				}
				request[fields.Property] = modifiedObjects
			} else {
				request[fields.Property] = val
			}
		} else {
			// field is not there or not send from front-end , so update the data
			request[fields.Property] = objectFields[fields.Property]
		}

	}
	rawUpdateRequest, _ := json.Marshal(request)
	return rawUpdateRequest
}

func (cm *ComponentManager) GetUpdateRequest(updateRequest map[string]interface{}, recordInfo datatypes.JSON, componentName string) []byte {
	var recordFields map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]

	componentSchema := cm.ComponentSchema[int64ComponentId]
	baseRecordSchema := componentSchema.RecordSchema
	json.Unmarshal(recordInfo, &recordFields)

	for _, recordSchema := range baseRecordSchema {
		// skip fields if overwrite is false, don't overwrite those fields.
		if recordSchema.Overwrite != nil {
			if !*recordSchema.Overwrite {
				updateRequest[recordSchema.Property] = recordFields[recordSchema.Property]
				continue
			}
		}
		if val, ok := updateRequest[recordSchema.Property]; ok {

			//do something here
			if recordSchema.Type == "int_array" {
				requestedIntArray := updateRequest[recordSchema.Property]
				updateRequest[recordSchema.Property] = requestedIntArray
			} else if recordSchema.Type == "object_array" {
				updatingData := updateRequest[recordSchema.Property]
				switch x := updatingData.(type) {

				case []interface{}:
					fmt.Println("[GetUpdateRequest] existing data:", recordFields[recordSchema.Property], " [Property] :", recordSchema.Property, " type :", x)
					if recordFields[recordSchema.Property] == nil {

						// which mean, existing data array is null, it shouldn't happen, but if that happen, this how we should deal with
						if len(updatingData.([]interface{})) > 0 {
							// front-end sending empty data
							var arrayOfUpdatingData = updatingData.([]interface{})
							var recordId int
							recordId = 1
							var modifiedObjectArray []interface{}
							fmt.Println("[GetUpdateRequest] recordFields[recordSchema.Property] is Null :", arrayOfUpdatingData)
							for _, internalObject := range arrayOfUpdatingData {
								var internalObjectFields = internalObject.(map[string]interface{})
								internalObjectFields["id"] = recordId
								modifiedObjectArray = append(modifiedObjectArray, internalObjectFields)
								recordId = recordId + 1
							}
							updateRequest[recordSchema.Property] = modifiedObjectArray
						} else {
							updateRequest[recordSchema.Property] = make([]interface{}, 0)
						}

					} else {

						var modifiedObjects []interface{}
						// we receive an array of objects to update
						var arrayOfUpdatingData = updatingData.([]interface{})

						// we already have an array of objects to update
						var arrayExistingObjects = recordFields[recordSchema.Property].([]interface{})
						if len(arrayExistingObjects) == 0 {
							fmt.Println("[GetUpdateRequest] Existing array is empty")
							// existing one is already empty, so we can just add all the objects
							// if we have the updating data is empty
							if len(arrayOfUpdatingData) == 0 {
								modifiedObjects = make([]interface{}, 0)
							} else {
								var recordId int
								recordId = 1
								for _, updatingObject := range arrayOfUpdatingData {
									var updatingObjectFields = updatingObject.(map[string]interface{})
									updatingObjectFields["id"] = recordId
									modifiedObjects = append(modifiedObjects, updatingObjectFields)
									recordId = recordId + 1
								}
							}
							fmt.Println("[GetUpdateRequest] Modified Objects:", modifiedObjects)

						} else {
							// if we have existing one, we should be carefully on update if any matched, otherwise add
							var updatingObjectFromRequest []interface{}
							// Get maximum id number the add new id
							//var recordId = len(arrayExistingObjects)
							var recordId int
							for _, existingObject := range arrayExistingObjects {
								var existingObjectField = existingObject.(map[string]interface{})
								existingObjectId := util.InterfaceToInt(existingObjectField["id"])
								if existingObjectId > recordId {
									recordId = existingObjectId
								}
							}
							for _, objects := range arrayOfUpdatingData { // 4, and 5
								var updateObjectFields = objects.(map[string]interface{})

								if _, ok := updateObjectFields["id"]; !ok {
									recordId = recordId + 1
									updateObjectFields["id"] = recordId
									modifiedObjects = append(modifiedObjects, updateObjectFields) // it should be new ones

								} else {
									updatingObjectFromRequest = append(updatingObjectFromRequest, objects) //4,5
								}

							}

							for _, existingObject := range arrayExistingObjects {
								var existingObjectFields = existingObject.(map[string]interface{})
								existingId := existingObjectFields["id"]
								for _, objects := range updatingObjectFromRequest {
									var updateObjectFields = objects.(map[string]interface{})
									updatingId := updateObjectFields["id"]
									if existingId == updatingId {
										for key, value := range updateObjectFields {
											if key != CreatedAt && key != LastUpdatedAt {
												existingObjectFields[key] = value
											}

										}
									}
								}

								modifiedObjects = append(modifiedObjects, existingObject) // this one now have the new one and update object fields
							}

							for _, updateObject := range updatingObjectFromRequest {
								var updateObjectFields = updateObject.(map[string]interface{})
								updatingId := updateObjectFields["id"]
								var isExist bool
								isExist = false
								for _, existingObject := range arrayExistingObjects {
									var existingObjectFields = existingObject.(map[string]interface{})
									existingId := existingObjectFields["id"]
									if updatingId == existingId {
										isExist = true
									}
								}
								if !isExist {
									modifiedObjects = append(modifiedObjects, updateObject)
								}
							}

						}

						// array is defined,
						var defaultProcessedObjects []interface{}
						for _, modifiedObject := range modifiedObjects {
							var internalFields = modifiedObject.(map[string]interface{})
							for _, insideField := range recordSchema.RecordSchema {
								if insideField.DefaultType != nil {
									if _, ok := internalFields[insideField.Property]; !ok {
										if *insideField.DefaultType == DefaultDataTypeCurrentTimestamp {
											defaultValue := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
											internalFields[insideField.Property] = defaultValue
										}

									}

								}
							}
							defaultProcessedObjects = append(defaultProcessedObjects, internalFields)
						}

						updateRequest[recordSchema.Property] = defaultProcessedObjects

					}
					// we shouldn't handle this case this field is array
				case interface{}:
					if recordFields[recordSchema.Property] == nil {
						var objectArray []interface{}
						objectArray = append(objectArray, recordFields[recordSchema.Property])
						updateRequest[recordSchema.Property] = objectArray
					} else {
						var objectFields = updatingData.(map[string]interface{})
						objectFields["id"] = util.GetMD5OfUUID()
						updatedValue := util.AppendToObjectArray(recordFields[recordSchema.Property], objectFields)
						updateRequest[recordSchema.Property] = updatedValue
					}
				default:
					fmt.Println("error in object array")
				}

			} else {
				updateRequest[recordSchema.Property] = val
			}
		} else {
			// field is not there or not send from front-end , so update the data
			updateRequest[recordSchema.Property] = recordFields[recordSchema.Property]
		}

	}

	// createdAt and createdBy fields are default fields. If that is available in the object, include them, otherwise, it will be gone.
	if createAt, ok := recordFields[CreatedAt]; ok {
		updateRequest["createdAt"] = createAt
	}
	if objectStatus, ok := recordFields[ObjectStatus]; ok {
		updateRequest["objectStatus"] = objectStatus
	}

	if createdBy, ok := recordFields[CreatedBy]; ok {
		updateRequest["createdBy"] = createdBy
	}
	rawUpdateRequest, _ := json.Marshal(updateRequest)
	return rawUpdateRequest
}
func (cm *ComponentManager) InitComponentManager(totalComponents int) {

	cm.ComponentLinkedFields = make(map[int]LinkedValues, totalComponents)
	cm.ComponentSchema = make(map[int]component.ComponentSchema, totalComponents)
	cm.ComponentNameIdMapping = make(map[string]int, totalComponents)
}
func (cm *ComponentManager) LoadTableSchema(listOfComponents *[]component.GeneralObject) {

	for _, componentInterface := range *listOfComponents {
		componentSchema := component.GetComponentSchema(componentInterface.ObjectInfo)
		cm.ComponentSchema[componentInterface.Id] = componentSchema
		cm.ComponentNameIdMapping[componentSchema.ComponentName] = componentInterface.Id
		cm.ComponentTables = append(cm.ComponentTables, componentSchema.TargetTable)
	}

}

func (cm *ComponentManager) LoadTableSchemaV1(listOfComponents []component.GeneralObject) {

	for _, componentInterface := range listOfComponents {
		componentSchema := component.GetComponentSchema(componentInterface.ObjectInfo)
		cm.ComponentSchema[componentInterface.Id] = componentSchema
		cm.ComponentNameIdMapping[componentSchema.ComponentName] = componentInterface.Id
		cm.ComponentTables = append(cm.ComponentTables, componentSchema.TargetTable)
	}

}
func (cm *ComponentManager) DownloadFile(contentUrl string) (error, string) {
	savedFileName := cm.ComponentContentConfig.DownloadDirectory + "/" + uuid.New().String()
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: savedFileName,
		Dir: false,
		//the repository with a subdirectory I would like to clone only
		Src:  contentUrl,
		Mode: getter.ClientModeFile,
		////define the type of detectors go getter should use, in this case only github is needed

		////provide the getter needed to download the files
		Getters:  cm.ComponentContentConfig.GetGetter(),
		Insecure: cm.ComponentContentConfig.Insecure,
	}
	//download the files

	if err := client.Get(); err != nil {
		return err, ""
	}

	return nil, savedFileName
}

func (cm *ComponentManager) ProcessDefaultObjectMapping(dbConnection *gorm.DB, requestFields map[string]interface{}, componentName string) map[string]interface{} {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	for _, recordSchemaField := range recordSchema {
		if recordSchemaField.DefaultObjectMapping != nil {
			var queryResults []map[string]interface{}
			if recordSchemaField.DefaultObjectMapping.Query != nil {
				// configured the query
				query := recordSchemaField.DefaultObjectMapping.Query.Query

				for _, replacementField := range recordSchemaField.DefaultObjectMapping.Query.ReplacementFields {
					if replacementField.Format == component.JsonToStringArray {
						value := requestFields[replacementField.Property]
						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						value := requestFields[replacementField.Property]
						replacementValue := util.InterfaceToString(value)
						query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}

				dbConnection.Raw(query).Scan(&queryResults)
			}

			if recordSchemaField.DefaultObjectMapping.Query == nil {
				if recordSchemaField.DefaultObjectMapping.Builder.SingleValue != nil {
					requestFields[recordSchemaField.Property] = requestFields[recordSchemaField.DefaultObjectMapping.Builder.SingleValue.Field]
				}
			} else {
				if recordSchemaField.DefaultObjectMapping.Builder.SingleValue != nil {
					requestFields[recordSchemaField.Property] = queryResults[0][recordSchemaField.DefaultObjectMapping.Builder.SingleValue.Field]
				}
			}

		}
	}
	return requestFields
}

func (cm *ComponentManager) ApplyFormatting(requestFields map[string]interface{}, componentName string) map[string]interface{} {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	for _, recordSchemaField := range recordSchema {
		if recordSchemaField.Formatter == "bcrypt" {
			fieldValue := requestFields[recordSchemaField.Property]
			generatedHashedPassword, _ := bcrypt.GenerateFromPassword([]byte(util.InterfaceToString(fieldValue)), bcrypt.DefaultCost)
			requestFields[recordSchemaField.Property] = string(generatedHashedPassword)
		}
	}
	return requestFields
}

func (cm *ComponentManager) PreprocessCreateRequestFields(requestFields map[string]interface{}, componentName string) map[string]interface{} {
	// first check any ignoreEmpty true, then filter them out
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	var preProcessedData = make(map[string]interface{})
	var propertyMap = make(map[string]int)

	for _, recordSchemaField := range recordSchema {
		// here we have default value is configured in the record schema
		if recordSchemaField.Default != nil {
			// default values are specified, and front-end is not sending this, so we will set this.
			_, ok := requestFields[recordSchemaField.Property]
			if !ok {
				preProcessedData[recordSchemaField.Property] = *recordSchemaField.Default
				propertyMap[recordSchemaField.Property] = 0
			}
		}
		if recordSchemaField.Type == "object_array" {
			if len(recordSchemaField.RecordSchema) > 0 {
				// array is defined,
				if requestFields[recordSchemaField.Property] != nil {
					var arrayOfUpdatingData = requestFields[recordSchemaField.Property].([]interface{})
					var modifiedInternalFields []map[string]interface{}
					var recordId int
					recordId = 1
					for _, updatingData := range arrayOfUpdatingData {
						var internalFields = updatingData.(map[string]interface{})
						for _, insideField := range recordSchemaField.RecordSchema {
							if insideField.DefaultType != nil {
								if *insideField.DefaultType == "current_timestamp" {
									defaultValue := util.GetCurrentTime("2006-01-02T15:04:05.000Z")
									internalFields[insideField.Property] = defaultValue
								}
							}
						}
						internalFields["id"] = recordId
						recordId = recordId + 1
						modifiedInternalFields = append(modifiedInternalFields, internalFields)
					}
					preProcessedData[recordSchemaField.Property] = modifiedInternalFields
					propertyMap[recordSchemaField.Property] = 0
				}

			}
		}

	}

	for _, recordSchemaField := range recordSchema {
		fieldValue, ok := requestFields[recordSchemaField.Property]
		// we already filled with default values
		if _, isAlreadyProcessed := preProcessedData[recordSchemaField.Property]; !isAlreadyProcessed {
			if recordSchemaField.IgnoreEmptyInCreate {
				// ignore empty is configured
				if ok {
					if fieldValue != nil {
						preProcessedData[recordSchemaField.Property] = fieldValue
						propertyMap[recordSchemaField.Property] = 0
					}
				}
			} else {
				preProcessedData[recordSchemaField.Property] = fieldValue
				propertyMap[recordSchemaField.Property] = 0
			}
		}

	}

	for requestField, requestValue := range requestFields {
		if _, ok := propertyMap[requestField]; !ok {
			preProcessedData[requestField] = requestValue
		}
	}

	return preProcessedData
}

//	func (cm *ComponentManager) PreprocessDynamicRecord (int64ComponentId int, templateId int, dbConnection *gorm.DB, generalResponse map[string]interface{}) map[string]interface{}{
//		if cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent != nil {
//			referenceComponentName := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent.ReferenceComponent
//			templateFieldName := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent.TemplateFieldName
//			referenceTable := cm.GetTargetTable(referenceComponentName) //it_service_category_template
//			err, templateTableFieldsObject := Get(dbConnection, referenceTable, templateId)
//
//			if err == nil {
//
//				// got all the template fields
//				templateFields := GetObjectFields(templateTableFieldsObject.ObjectInfo)
//				templateFieldRecords := templateFields[templateFieldName]
//				serialisedRecords := GetInterfaceToSerialisation(templateFieldRecords)
//				var listOfTableRecords []TemplateRecords
//				json.Unmarshal(serialisedRecords, &listOfTableRecords)
//				// we got the fields
//
//				for _, templateRecord := range listOfTableRecords {
//					// got each of the template records
//					recordInfo := component.RecordInfo{}
//					recordInfo.Property = templateRecord.Property
//					_, ok := generalResponse[templateRecord.Property]
//					if ok {
//						preProcessedData[recordSchemaField.Property] = *recordSchemaField.Default
//					}
//				}
//
//			}
//
//		}
//
//		return generalResponse
//	}
func (cm *ComponentManager) ProcessDeleteDependencyInjection(dbConnection *gorm.DB, recordId int, componentName string) error {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	dependencyResourceInjection := cm.ComponentSchema[int64ComponentId].DeleteDependencyResourceInjection
	if dependencyResourceInjection != nil {
		for _, resourceInjection := range dependencyResourceInjection.ResourceInjection {
			executionQuery := resourceInjection.Query.Query
			for _, replacementField := range resourceInjection.Query.ReplacementFields {
				executionQuery = strings.Replace(executionQuery, "["+replacementField.Field+"]", strconv.Itoa(recordId), 1)
			}

			if resourceInjection.Builder.DBExecution != nil {
				// lets do the database insert
				err := dbConnection.Exec(executionQuery).Error
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (cm *ComponentManager) ProcessCreateDependencyInjection(dbConnection *gorm.DB, objectFields map[string]interface{}, componentName string) error {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	dependencyResourceInjection := cm.ComponentSchema[int64ComponentId].CreateDependencyResourceInjection

	if dependencyResourceInjection != nil {
		for _, resourceInjection := range dependencyResourceInjection.ResourceInjection {
			executionQuery := resourceInjection.Query.Query

			for _, replacementField := range resourceInjection.Query.ReplacementFields {
				if replacementField.Format == component.JsonToStringArray {
					value := objectFields[replacementField.Property]
					replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
					executionQuery = strings.Replace(executionQuery, "["+replacementField.Field+"]", replacementValue, 1)
				} else {
					value := objectFields[replacementField.Property]
					replacementValue := util.InterfaceToString(value)
					executionQuery = strings.Replace(executionQuery, "["+replacementField.Field+"]", replacementValue, 1)
				}

			}

			if resourceInjection.Builder.DBExecution != nil {
				// lets do the database insert
				err := dbConnection.Exec(executionQuery).Error
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (cm *ComponentManager) GetQueryResponse(dbConnection *gorm.DB, queryResponseFields map[string]interface{}) (error, interface{}) {
	var queryResults []map[string]interface{}
	tableObjectResponse := component.TableObjectResponse{}

	serialisedData, err := json.Marshal(queryResponseFields)

	// lets check the basic structure
	queryResponseRequest := analytics.QueryResponseRequest{}
	json.Unmarshal(serialisedData, &queryResponseRequest)
	query := queryResponseRequest.Query
	//It will replace all occurrences of oldSubstring with newSubstring in the originalString.
	if len(queryResponseRequest.Param) > 0 {
		// no param defined
		for key, value := range queryResponseRequest.Param {
			query = strings.Replace(query, "{{"+key+"}}", util.InterfaceToString(value), -1)
		}
	}
	err = dbConnection.Raw(query).Scan(&queryResults).Error
	if err != nil {
		return err, tableObjectResponse
	}

	baseBuilder := analytics.BaseBuilder{}
	baseBuilder.QueryResults = queryResults
	baseBuilder.QueryResponseRequest = serialisedData
	baseBuilder.Init()
	if queryResponseRequest.Format == "table" {
		responseBuilder := analytics.TableResponseBuilder{BaseBuilder: &baseBuilder}
		return responseBuilder.BuildResponse()
	} else if queryResponseRequest.Format == "chart" {
		responseBuilder := analytics.ChartResponseBuilder{BaseBuilder: &baseBuilder}
		return responseBuilder.BuildResponse()
	} else if queryResponseRequest.Format == "timeline" {
		responseBuilder := analytics.TimelineResponseBuilder{BaseBuilder: &baseBuilder}
		return responseBuilder.BuildResponse()
	}

	return nil, tableObjectResponse
}

type Result struct {
	ObjectInfo datatypes.JSON
}
type DatabaseSchemaNodeChildren struct {
	Data  string `json:"data"`
	Label string `json:"label"`
	Icon  string `json:"icon"`
}
type DatabaseSchemaNode struct {
	Data         string                       `json:"data"`
	Label        string                       `json:"label"`
	ExpandedIcon string                       `json:"expandedIcon"`
	CollapseIcon string                       `json:"collapseIcon"`
	Children     []DatabaseSchemaNodeChildren `json:"children"`
}

func (cm *ComponentManager) GetDatabaseTableSchema(dbConnection *gorm.DB, componentName string, connectedTables []string) interface{} {
	tableInformationSchema := "SELECT * FROM information_schema.tables WHERE table_schema =\"" + componentName + "\""
	var queryResults []map[string]interface{}
	dbConnection.Raw(tableInformationSchema).Scan(&queryResults)
	var listOfTables []string
	for _, result := range queryResults {
		tableName := result["TABLE_NAME"]
		for _, connectedTable := range connectedTables {
			if tableName == connectedTable {
				listOfTables = append(listOfTables, util.InterfaceToString(tableName))
			}
		}

	}
	var arrayOfSchemaNodes []DatabaseSchemaNode
	for _, tableName := range listOfTables {
		if tableName == "machine_statistics" {
			internalDataQuery := "SELECT stats_info FROM " + tableName + " LIMIT 1"
			type StatsResult struct {
				StatsInfo datatypes.JSON
			}
			re := StatsResult{}
			err := dbConnection.Raw(internalDataQuery).Scan(&re).Error

			if err == nil {

				var objectFields map[string]interface{}
				json.Unmarshal(re.StatsInfo, &objectFields)
				var arrayOfChildren []DatabaseSchemaNodeChildren
				databaseSchemaNode := DatabaseSchemaNode{}
				databaseSchemaNode.Label = tableName
				databaseSchemaNode.CollapseIcon = "pi pi-table"
				databaseSchemaNode.ExpandedIcon = "pi pi-folder-open"
				databaseSchemaNode.Data = tableName
				for key, _ := range objectFields {
					databaseSchemaNodeChildren := DatabaseSchemaNodeChildren{}
					databaseSchemaNodeChildren.Icon = "pi pi-minus-circle"
					databaseSchemaNodeChildren.Label = key
					databaseSchemaNodeChildren.Data = key
					arrayOfChildren = append(arrayOfChildren, databaseSchemaNodeChildren)
				}
				databaseSchemaNode.Children = arrayOfChildren
				arrayOfSchemaNodes = append(arrayOfSchemaNodes, databaseSchemaNode)

			}
		} else {
			internalDataQuery := "SELECT object_info FROM " + tableName + " LIMIT 1"
			re := Result{}
			err := dbConnection.Raw(internalDataQuery).Scan(&re).Error
			if err == nil {

				var objectFields map[string]interface{}
				json.Unmarshal(re.ObjectInfo, &objectFields)
				var arrayOfChildren []DatabaseSchemaNodeChildren
				databaseSchemaNode := DatabaseSchemaNode{}
				databaseSchemaNode.Label = tableName
				databaseSchemaNode.CollapseIcon = "pi pi-table"
				databaseSchemaNode.ExpandedIcon = "pi pi-folder-open"
				databaseSchemaNode.Data = tableName
				for key, _ := range objectFields {
					databaseSchemaNodeChildren := DatabaseSchemaNodeChildren{}
					databaseSchemaNodeChildren.Icon = "pi pi-minus-circle"
					databaseSchemaNodeChildren.Label = key
					databaseSchemaNodeChildren.Data = key
					arrayOfChildren = append(arrayOfChildren, databaseSchemaNodeChildren)
				}
				databaseSchemaNode.Children = arrayOfChildren
				arrayOfSchemaNodes = append(arrayOfSchemaNodes, databaseSchemaNode)

			}
		}

	}
	return arrayOfSchemaNodes

}
func (cm *ComponentManager) GetTableRecordsV1(dbConnection *gorm.DB, listOfObjects []component.GeneralObject, totalRecords int64, componentName string, outFields string, zone string) (error, datatypes.JSON) {

	var err error
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	tableObjectResponse := component.TableObjectResponse{}
	rawJSONResponse, err := json.Marshal(tableObjectResponse)

	if err != nil {
		return err, rawJSONResponse
	}

	if listOfObjects == nil {
		return err, rawJSONResponse
	}
	// based on permission, we should send the project id, contact permission service to get the list of projects
	var modifiedObjectList []datatypes.JSON
	tableObjectResponse.TotalRowCount = totalRecords
	tableObjectResponse.CurrentRowCount = int64(len(listOfObjects))

	// now prepare the table header object mapping data if any
	arrayOfTableSchema := cm.ComponentSchema[int64ComponentId].TableSchema
	for index, tableFieldSchema := range arrayOfTableSchema {
		if tableFieldSchema.HeaderObjectMapping != nil {
			if tableFieldSchema.HeaderObjectMapping.Query != nil {
				var queryResults []map[string]interface{}
				dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)

				var objectKeyValues = make(map[string]string, len(queryResults))
				for _, result := range queryResults {
					objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
				}
				raw, _ := json.Marshal(objectKeyValues)
				arrayOfTableSchema[index].ObjectList = raw
				// set the object to null, so that response won't send that field
				//arrayOfTableSchema[index].HeaderObjectMapping = nil
			} else if tableFieldSchema.HeaderObjectMapping.Predefined != nil {
				predefinedResults := tableFieldSchema.HeaderObjectMapping.Predefined.Data
				var objectKeyValues = make(map[string]string, len(predefinedResults))
				if tableFieldSchema.HeaderObjectMapping.Builder.KeyValue != nil {
					for _, result := range predefinedResults {
						objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
					}
				}

				raw, _ := json.Marshal(objectKeyValues)
				arrayOfTableSchema[index].ObjectList = raw
				// set the object to null, so that response won't send that field
				arrayOfTableSchema[index].HeaderObjectMapping = nil
			}
			if tableFieldSchema.HeaderFontColorObjectMapping != nil {
				if tableFieldSchema.HeaderFontColorObjectMapping.Predefined != nil {
					predefinedResults := tableFieldSchema.HeaderFontColorObjectMapping.Predefined.Data
					var objectKeyValues = make(map[string]string, len(predefinedResults))
					if tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue != nil {
						for _, result := range predefinedResults {
							objectKeyValues[result[tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Value].(string)
						}
					}

					raw, _ := json.Marshal(objectKeyValues)
					arrayOfTableSchema[index].FontColorList = raw
					// set the object to null, so that response won't send that field
					//arrayOfTableSchema[index].HeaderFontColorObjectMapping = nil
				}

			}

		}
	}

	arrayOfRecordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	var recordSchemaMap = make(map[string]component.RecordSchema)
	for _, recordSchema := range arrayOfRecordSchema {
		recordSchemaMap[recordSchema.Property] = recordSchema
	}

	var outFieldsTableMap = make(map[string]interface{})
	var outFieldsArray []string
	if outFields != "" {
		outFieldsArray = strings.Split(outFields, ",")
	}
	for _, objectInterface := range listOfObjects {
		var tableMap = make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &tableMap)

		if len(outFieldsArray) > 0 {
			for _, outField := range outFieldsArray {
				recordSchema := recordSchemaMap[outField]
				if recordSchema.LinkedObjectMapping != nil {

					if recordSchema.LinkedObjectMapping.Query != nil {
						// if the query object is configured
						queryTemplate := recordSchema.LinkedObjectMapping.Query.Query
						var queryResults []map[string]interface{}

						for _, replacementField := range recordSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							} else {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							}

						}

						dbConnection.Raw(queryTemplate).Scan(&queryResults)
						if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							outFieldsTableMap[outField] = queryResults[0][recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
						} else {
							outFieldsTableMap[outField] = tableMap[outField]
						}

					}

				}
			}
			rawData, _ := json.Marshal(outFieldsTableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)

		} else {
			tableMap["id"] = objectInterface.Id

			for field, _ := range tableMap {

				recordSchema := recordSchemaMap[field]

				// This block stop to fetch unnessary data
				if recordSchema.CanDisplayTable {
					delete(tableMap, field)
					continue
				}

				if recordSchema.LinkedObjectMapping != nil {

					if recordSchema.LinkedObjectMapping.Query != nil {
						queryTemplate := recordSchema.LinkedObjectMapping.Query.Query
						var queryResults []map[string]interface{}
						for _, replacementField := range recordSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							} else {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							}

						}
						dbConnection.Raw(queryTemplate).Scan(&queryResults)
						fmt.Println("queryTemplate:", queryTemplate)
						if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(queryResults) > 0 {
								tableMap[field] = queryResults[0][recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
							} else {
								tableMap[field] = "-"
							}

						}
					} else {
						if recordSchema.LinkedObjectMapping.Predefined != nil {
							if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
								existingValue := tableMap[field]
								for _, internalPredefinedData := range recordSchema.LinkedObjectMapping.Predefined.Data {
									if value, ok := internalPredefinedData["id"]; ok {
										if existingValue == value {
											tableMap[field] = internalPredefinedData[recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
										}
									}
								}
							}
						}
					}

				}

				// once the field is updated, lets check the type
				if recordSchema.Type == component.TableDataTypeDateTime {
					// should be in the user timezone format
					if value, ok := tableMap[field]; ok {
						if value != nil {
							currentValue := value.(string)
							var correctedValue string
							if currentValue != "" {
								correctedValue = util.ISO2TableDateTimeFormat(zone, currentValue)
							}

							tableMap[field] = correctedValue
						}

					}

				} else if recordSchema.Type == component.TableDataTypeDate {
					// should be in the user timezone format
					if value, ok := tableMap[field]; ok {
						if value != nil {
							currentValue := util.InterfaceToString(value)
							var correctedValue string
							if currentValue != "" {
								correctedValue = util.ISO2TableDateFormat(zone, currentValue)
							}
							tableMap[field] = correctedValue
						}

					}

				}
			}

			// check we can build the additional table records, this is normally happen we want to show the data
			// but not in the table records

			tableAdditionalRecords := cm.ComponentSchema[int64ComponentId].TableAdditionalRecords
			for _, tableAdditionalRecord := range tableAdditionalRecords {
				queryTemplate := tableAdditionalRecord.ObjectMapping.Query.Query
				var queryResults []map[string]interface{}
				for _, replacementField := range tableAdditionalRecord.ObjectMapping.Query.ReplacementFields {
					if replacementField.Format == component.JsonToStringArray {
						value := tableMap[replacementField.Property]
						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						value := tableMap[replacementField.Property]
						replacementValue := util.InterfaceToString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}
				dbConnection.Raw(queryTemplate).Scan(&queryResults)
				if len(tableAdditionalRecord.ObjectMapping.Query.OutFieldsMapping) > 0 {
					// first apply the out field mapping logics using formattters
					for _, outField := range tableAdditionalRecord.ObjectMapping.Query.OutFieldsMapping {
						for index, result := range queryResults {
							if queryResultValue, ok := result[outField.Field]; ok {
								queryResults[index][outField.Field] = applyFormatters(outField.Formatter, queryResultValue)
							}
						}
					}
				}

				if tableAdditionalRecord.ObjectMapping.Builder.SingleValue != nil {
					field := tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Field
					if len(queryResults) > 0 {

						if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == component.TableDataTypeDateTime {
							// should be in the user timezone format
							currentValue := queryResults[0][field].(string)
							correctedValue := util.ISO2TableDateTimeFormat(zone, currentValue)
							tableMap[field] = correctedValue

						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == component.TableDataTypeDate {
							// should be in the user timezone format
							currentValue := queryResults[0][field].(string)
							correctedValue := util.ISO2TableDateTimeFormat(zone, currentValue)
							tableMap[field] = correctedValue

						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "transform" {
							var currentValue string
							for _, rowResult := range queryResults {
								if val, ok := rowResult["id"]; ok {
									if util.InterfaceToInt(val) == objectInterface.Id {
										currentValue = rowResult[field].(string)
										break
									}
								}
							}
							tableMap[field] = currentValue

						} else {
							tableMap[field] = queryResults[0][field]
						}

					} else {
						//something wrong in the query , or actually results are empty
						if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "int" {
							tableMap[field] = 0
						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "text" {
							tableMap[field] = ""
						}
					}

				} else if tableAdditionalRecord.ObjectMapping.Builder.ObjectArray != nil {
					valuesArray := tableAdditionalRecord.ObjectMapping.Builder.ObjectArray.Values
					var listOfArrayObjects []map[string]interface{}
					for _, result := range queryResults {
						responseArray := make(map[string]interface{}, len(valuesArray))
						for _, objectArrayValue := range valuesArray {
							objectResponse := util.InterfaceToString(result[objectArrayValue])
							responseArray[objectArrayValue] = objectResponse
						}
						listOfArrayObjects = append(listOfArrayObjects, responseArray)
					}

					tableMap[tableAdditionalRecord.Property] = listOfArrayObjects
				} else if tableAdditionalRecord.ObjectMapping.Builder.SingleObjectValueArray != nil {
					field := tableAdditionalRecord.ObjectMapping.Builder.SingleObjectValueArray.Field
					var arrayOfValues []string
					for _, result := range queryResults {
						value := util.InterfaceToString(result[field])
						arrayOfValues = append(arrayOfValues, value)
					}
					tableMap[tableAdditionalRecord.Property] = arrayOfValues
				}
			}

			rawData, _ := json.Marshal(tableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)
		}

	}

	if len(listOfObjects) == 0 {
		tableObjectResponse.Data = make([]datatypes.JSON, 0)
	} else {
		tableObjectResponse.Data = modifiedObjectList
	}
	tableObjectResponse.Header = arrayOfTableSchema
	rawJSONResponse, _ = json.Marshal(tableObjectResponse)
	return nil, rawJSONResponse
}

func (cm *ComponentManager) GetTableRecords(dbConnection *gorm.DB, listOfObjects *[]component.GeneralObject, totalRecords int64, componentName string, outFields string, zone string) (error, datatypes.JSON) {

	var err error
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	tableObjectResponse := component.TableObjectResponse{}
	rawJSONResponse, err := json.Marshal(tableObjectResponse)

	var listCount = 0
	if listOfObjects != nil {
		listCount = len(*listOfObjects)
	}

	if err != nil {
		return err, rawJSONResponse
	}

	if listOfObjects == nil {
		return err, rawJSONResponse
	}
	// based on permission, we should send the project id, contact permission service to get the list of projects
	var modifiedObjectList []datatypes.JSON
	tableObjectResponse.TotalRowCount = totalRecords
	tableObjectResponse.CurrentRowCount = int64(listCount)

	// now prepare the table header object mapping data if any
	arrayOfTableSchema := cm.ComponentSchema[int64ComponentId].TableSchema
	for index, tableFieldSchema := range arrayOfTableSchema {
		if tableFieldSchema.HeaderObjectMapping != nil {
			if tableFieldSchema.HeaderObjectMapping.Query != nil {
				var queryResults []map[string]interface{}
				dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)

				var objectKeyValues = make(map[string]string, len(queryResults))
				for _, result := range queryResults {
					objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
				}
				raw, _ := json.Marshal(objectKeyValues)
				arrayOfTableSchema[index].ObjectList = raw
				// set the object to null, so that response won't send that field
				//arrayOfTableSchema[index].HeaderObjectMapping = nil
			} else if tableFieldSchema.HeaderObjectMapping.Predefined != nil {
				predefinedResults := tableFieldSchema.HeaderObjectMapping.Predefined.Data
				var objectKeyValues = make(map[string]string, len(predefinedResults))
				if tableFieldSchema.HeaderObjectMapping.Builder.KeyValue != nil {
					for _, result := range predefinedResults {
						objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
					}
				}

				raw, _ := json.Marshal(objectKeyValues)
				arrayOfTableSchema[index].ObjectList = raw
				// set the object to null, so that response won't send that field
				arrayOfTableSchema[index].HeaderObjectMapping = nil
			}
			if tableFieldSchema.HeaderFontColorObjectMapping != nil {
				if tableFieldSchema.HeaderFontColorObjectMapping.Predefined != nil {
					predefinedResults := tableFieldSchema.HeaderFontColorObjectMapping.Predefined.Data
					var objectKeyValues = make(map[string]string, len(predefinedResults))
					if tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue != nil {
						for _, result := range predefinedResults {
							objectKeyValues[result[tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Value].(string)
						}
					}

					raw, _ := json.Marshal(objectKeyValues)
					arrayOfTableSchema[index].FontColorList = raw
					// set the object to null, so that response won't send that field
					//arrayOfTableSchema[index].HeaderFontColorObjectMapping = nil
				}

			}

		}
	}

	arrayOfRecordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	var recordSchemaMap = make(map[string]component.RecordSchema)
	for _, recordSchema := range arrayOfRecordSchema {
		recordSchemaMap[recordSchema.Property] = recordSchema
	}

	var outFieldsTableMap = make(map[string]interface{})
	var outFieldsArray []string
	if outFields != "" {
		outFieldsArray = strings.Split(outFields, ",")
	}
	for _, objectInterface := range *listOfObjects {
		var tableMap = make(map[string]interface{})
		json.Unmarshal(objectInterface.ObjectInfo, &tableMap)

		if len(outFieldsArray) > 0 {
			for _, outField := range outFieldsArray {
				recordSchema := recordSchemaMap[outField]
				if recordSchema.LinkedObjectMapping != nil {

					if recordSchema.LinkedObjectMapping.Query != nil {
						// if the query object is configured
						queryTemplate := recordSchema.LinkedObjectMapping.Query.Query
						var queryResults []map[string]interface{}

						for _, replacementField := range recordSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							} else {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							}

						}

						dbConnection.Raw(queryTemplate).Scan(&queryResults)
						if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							outFieldsTableMap[outField] = queryResults[0][recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
						} else {
							outFieldsTableMap[outField] = tableMap[outField]
						}

					}

				}
			}
			rawData, _ := json.Marshal(outFieldsTableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)

		} else {
			tableMap["id"] = objectInterface.Id

			for field, _ := range tableMap {

				recordSchema := recordSchemaMap[field]

				// This block stop to fetch unnessary data
				if recordSchema.CanDisplayTable {
					delete(tableMap, field)
					continue
				}

				if recordSchema.LinkedObjectMapping != nil {

					if recordSchema.LinkedObjectMapping.Query != nil {
						queryTemplate := recordSchema.LinkedObjectMapping.Query.Query
						var queryResults []map[string]interface{}
						for _, replacementField := range recordSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							} else {
								value := tableMap[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
							}

						}
						dbConnection.Raw(queryTemplate).Scan(&queryResults)
						fmt.Println("queryTemplate:", queryTemplate)
						if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(queryResults) > 0 {
								tableMap[field] = queryResults[0][recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
							} else {
								tableMap[field] = "-"
							}

						}
					} else {
						if recordSchema.LinkedObjectMapping.Predefined != nil {
							if recordSchema.LinkedObjectMapping.Builder.SingleValue != nil {
								existingValue := tableMap[field]
								for _, internalPredefinedData := range recordSchema.LinkedObjectMapping.Predefined.Data {
									if value, ok := internalPredefinedData["id"]; ok {
										if existingValue == value {
											tableMap[field] = internalPredefinedData[recordSchema.LinkedObjectMapping.Builder.SingleValue.Field]
										}
									}
								}
							}
						}
					}

				}

				// once the field is updated, lets check the type
				if recordSchema.Type == component.TableDataTypeDateTime {
					// should be in the user timezone format
					if value, ok := tableMap[field]; ok {
						if value != nil {
							currentValue := value.(string)
							var correctedValue string
							if currentValue != "" {
								correctedValue = util.ISO2TableDateTimeFormat(zone, currentValue)
							}

							tableMap[field] = correctedValue
						}

					}

				} else if recordSchema.Type == component.TableDataTypeDate {
					// should be in the user timezone format

					if value, ok := tableMap[field]; ok {

						if value != nil {
							currentValue := value.(string)

							var correctedValue string
							if currentValue != "" {
								correctedValue = util.ISO2TableDateFormat(zone, currentValue)

							}

							tableMap[field] = correctedValue
						}

					}

				}
			}

			// check we can build the additional table records, this is normally happen we want to show the data
			// but not in the table records

			tableAdditionalRecords := cm.ComponentSchema[int64ComponentId].TableAdditionalRecords
			for _, tableAdditionalRecord := range tableAdditionalRecords {
				queryTemplate := tableAdditionalRecord.ObjectMapping.Query.Query
				var queryResults []map[string]interface{}
				for _, replacementField := range tableAdditionalRecord.ObjectMapping.Query.ReplacementFields {
					if replacementField.Format == component.JsonToStringArray {
						value := tableMap[replacementField.Property]
						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						value := tableMap[replacementField.Property]
						replacementValue := util.InterfaceToString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}
				dbConnection.Raw(queryTemplate).Scan(&queryResults)
				if len(tableAdditionalRecord.ObjectMapping.Query.OutFieldsMapping) > 0 {
					// first apply the out field mapping logics using formattters
					for _, outField := range tableAdditionalRecord.ObjectMapping.Query.OutFieldsMapping {
						for index, result := range queryResults {
							if queryResultValue, ok := result[outField.Field]; ok {
								queryResults[index][outField.Field] = applyFormatters(outField.Formatter, queryResultValue)
							}
						}
					}
				}

				if tableAdditionalRecord.ObjectMapping.Builder.SingleValue != nil {
					field := tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Field
					if len(queryResults) > 0 {

						if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == component.TableDataTypeDateTime {
							// should be in the user timezone format
							currentValue := queryResults[0][field].(string)
							correctedValue := util.ISO2TableDateTimeFormat(zone, currentValue)
							tableMap[field] = correctedValue

						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == component.TableDataTypeDate {
							// should be in the user timezone format
							currentValue := queryResults[0][field].(string)
							correctedValue := util.ISO2TableDateTimeFormat(zone, currentValue)
							tableMap[field] = correctedValue

						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "transform" {
							var currentValue string
							for _, rowResult := range queryResults {
								if val, ok := rowResult["id"]; ok {
									if util.InterfaceToInt(val) == objectInterface.Id {
										currentValue = rowResult[field].(string)
										break
									}
								}
							}
							tableMap[field] = currentValue

						} else {
							tableMap[field] = queryResults[0][field]
						}

					} else {
						//something wrong in the query , or actually results are empty
						if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "int" {
							tableMap[field] = 0
						} else if tableAdditionalRecord.ObjectMapping.Builder.SingleValue.Type == "text" {
							tableMap[field] = ""
						}
					}

				} else if tableAdditionalRecord.ObjectMapping.Builder.ObjectArray != nil {
					valuesArray := tableAdditionalRecord.ObjectMapping.Builder.ObjectArray.Values
					var listOfArrayObjects []map[string]interface{}
					for _, result := range queryResults {
						responseArray := make(map[string]interface{}, len(valuesArray))
						for _, objectArrayValue := range valuesArray {
							objectResponse := util.InterfaceToString(result[objectArrayValue])
							responseArray[objectArrayValue] = objectResponse
						}
						listOfArrayObjects = append(listOfArrayObjects, responseArray)
					}

					tableMap[tableAdditionalRecord.Property] = listOfArrayObjects
				} else if tableAdditionalRecord.ObjectMapping.Builder.SingleObjectValueArray != nil {
					field := tableAdditionalRecord.ObjectMapping.Builder.SingleObjectValueArray.Field
					var arrayOfValues []string
					for _, result := range queryResults {
						value := util.InterfaceToString(result[field])
						arrayOfValues = append(arrayOfValues, value)
					}
					tableMap[tableAdditionalRecord.Property] = arrayOfValues
				}
			}

			rawData, _ := json.Marshal(tableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)
		}

	}

	if listCount == 0 {
		tableObjectResponse.Data = make([]datatypes.JSON, 0)
	} else {
		tableObjectResponse.Data = modifiedObjectList
	}
	tableObjectResponse.Header = arrayOfTableSchema
	rawJSONResponse, _ = json.Marshal(tableObjectResponse)
	return nil, rawJSONResponse
}

func (cm *ComponentManager) TableRecordsToArray(dbConnection *gorm.DB, listOfObjects *[]component.GeneralObject, componentName string, outFields string) (error, datatypes.JSON) {

	var err error
	var totalRecords int64
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	totalRecords = int64(len(*listOfObjects))
	tableObjectResponse := component.TableObjectResponse{}
	rawJSONResponse, err := json.Marshal(tableObjectResponse)
	if err != nil {
		return err, rawJSONResponse
	}
	// based on permission, we should send the project id, contact permission service to get the list of projects
	var modifiedObjectList []datatypes.JSON
	tableObjectResponse.TotalRowCount = totalRecords
	tableObjectResponse.CurrentRowCount = int64(len(*listOfObjects))

	// we need to prepare the linked data if any
	linkedData := cm.ComponentLinkedFields[int64ComponentId]

	// now prepare the table header object mapping data if any
	arrayOfTableSchema := cm.ComponentSchema[int64ComponentId].TableSchema
	for index, tableFieldSchema := range arrayOfTableSchema {
		if tableFieldSchema.HeaderObjectMapping != nil {

			var queryResults []map[string]interface{}
			dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)

			var objectKeyValues = make(map[string]string, len(queryResults))
			for _, result := range queryResults {
				objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
			}
			raw, _ := json.Marshal(objectKeyValues)
			arrayOfTableSchema[index].ObjectList = raw
			// set the object to null, so that response won't send that field
			arrayOfTableSchema[index].HeaderObjectMapping = nil
		}
	}

	var tableMap = make(map[string]interface{})
	var outFieldsTableMap = make(map[string]interface{})
	var arrayResponse = make(map[string][]interface{}, 20)
	outFieldsArray := strings.Split(outFields, ",")
	for _, objectInterface := range *listOfObjects {
		json.Unmarshal(objectInterface.ObjectInfo, &tableMap)

		if len(outFieldsArray) > 0 {
			for _, outField := range outFieldsArray {
				outFieldsTableMap[outField] = tableMap[outField]
			}
			rawData, _ := json.Marshal(outFieldsTableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)

		} else {
			tableMap["recordId"] = objectInterface.Id
			for key, value := range linkedData.FieldMapping {
				currentKey := util.InterfaceToInt(tableMap[key])
				replacementData := value[currentKey]
				tableMap[key] = replacementData

			}
			rawData, _ := json.Marshal(tableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)
		}
	}

	for _, modifiedObject := range modifiedObjectList {
		var ddd = make(map[string]interface{})
		rawData, _ := modifiedObject.MarshalJSON()
		json.Unmarshal(rawData, &ddd)
		for key, value := range ddd {
			if val, ok := arrayResponse[key]; ok {
				//do something here
				val = append(val, value)
				arrayResponse[key] = val

			} else {
				var arrayOfData []interface{}
				arrayOfData = append(arrayOfData, value)
				arrayResponse[key] = arrayOfData
			}
		}
	}
	rawJSONResponse, err = json.Marshal(arrayResponse)

	return err, rawJSONResponse

}

func (cm *ComponentManager) TableRecordsToArrayV1(dbConnection *gorm.DB, listOfObjects []component.GeneralObject, componentName string, outFields string) (error, datatypes.JSON) {

	var err error
	var totalRecords int64
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	totalRecords = int64(len(listOfObjects))
	tableObjectResponse := component.TableObjectResponse{}
	rawJSONResponse, err := json.Marshal(tableObjectResponse)
	if err != nil {
		return err, rawJSONResponse
	}
	// based on permission, we should send the project id, contact permission service to get the list of projects
	var modifiedObjectList []datatypes.JSON
	tableObjectResponse.TotalRowCount = totalRecords
	tableObjectResponse.CurrentRowCount = int64(len(listOfObjects))

	// we need to prepare the linked data if any
	linkedData := cm.ComponentLinkedFields[int64ComponentId]

	// now prepare the table header object mapping data if any
	arrayOfTableSchema := cm.ComponentSchema[int64ComponentId].TableSchema
	for index, tableFieldSchema := range arrayOfTableSchema {
		if tableFieldSchema.HeaderObjectMapping != nil {

			var queryResults []map[string]interface{}
			dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)

			var objectKeyValues = make(map[string]string, len(queryResults))
			for _, result := range queryResults {
				objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
			}
			raw, _ := json.Marshal(objectKeyValues)
			arrayOfTableSchema[index].ObjectList = raw
			// set the object to null, so that response won't send that field
			arrayOfTableSchema[index].HeaderObjectMapping = nil
		}
	}

	var tableMap = make(map[string]interface{})
	var outFieldsTableMap = make(map[string]interface{})
	var arrayResponse = make(map[string][]interface{}, 20)
	outFieldsArray := strings.Split(outFields, ",")
	for _, objectInterface := range listOfObjects {
		json.Unmarshal(objectInterface.ObjectInfo, &tableMap)

		if len(outFieldsArray) > 0 {
			for _, outField := range outFieldsArray {
				outFieldsTableMap[outField] = tableMap[outField]
			}
			rawData, _ := json.Marshal(outFieldsTableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)

		} else {
			tableMap["recordId"] = objectInterface.Id
			for key, value := range linkedData.FieldMapping {
				currentKey := util.InterfaceToInt(tableMap[key])
				replacementData := value[currentKey]
				tableMap[key] = replacementData

			}
			rawData, _ := json.Marshal(tableMap)
			modifiedObjectList = append(modifiedObjectList, rawData)
		}
	}

	for _, modifiedObject := range modifiedObjectList {
		var ddd = make(map[string]interface{})
		rawData, _ := modifiedObject.MarshalJSON()
		json.Unmarshal(rawData, &ddd)
		for key, value := range ddd {
			if val, ok := arrayResponse[key]; ok {
				//do something here
				val = append(val, value)
				arrayResponse[key] = val

			} else {
				var arrayOfData []interface{}
				arrayOfData = append(arrayOfData, value)
				arrayResponse[key] = arrayOfData
			}
		}
	}
	rawJSONResponse, err = json.Marshal(arrayResponse)

	return err, rawJSONResponse

}

func GetObjects(database *gorm.DB, table string, objectCount ...int) (error, []component.GeneralObject) {
	var err error

	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Find(&dbObjects).Error
	}

	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}

}
func GetConditionalObjects(database *gorm.DB, table string, condition string, objectCount ...int) (error, []component.GeneralObject) {
	var err error
	var dbObjects []component.GeneralObject
	if len(objectCount) > 0 {
		err = database.Table(table).Where(condition).Limit(objectCount[0]).Find(&dbObjects).Error
	} else {
		err = database.Table(table).Where(condition).Find(&dbObjects).Error
	}
	if err != nil {
		return err, dbObjects
	} else {
		return nil, dbObjects
	}
}
func Get(database *gorm.DB, table string, recordId int) (error, component.GeneralObject) {
	generalObject := component.GeneralObject{Id: recordId}
	err := database.Table(table).Find(&generalObject).Error
	if err != nil {
		return err, generalObject
	} else {
		return nil, generalObject
	}
}

func (cm *ComponentManager) GetNewRecordResponse(zone string, dbConnection *gorm.DB, componentName string) map[string]interface{} {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	newRecordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	response := make(map[string]interface{}, len(newRecordSchema))
	// normally if the field is linked field, in the record, we will see the id
	for _, recordSchema := range newRecordSchema {
		recordInfo := component.RecordInfo{}
		if recordSchema.ResponseObjectMapping != nil {

			if recordSchema.ResponseObjectMapping.Query != nil {
				var queryResults []map[string]interface{}
				dbConnection.Raw(recordSchema.ResponseObjectMapping.Query.Query).Scan(&queryResults)
				if recordSchema.ResponseObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0

					for _, queryResult := range queryResults {
						id := int(queryResult[recordSchema.ResponseObjectMapping.Builder.SingleDropdown.Index].(int32))
						dropdownValue := queryResult[recordSchema.ResponseObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if index == 0 {
							recordInfo.Index = id
							recordInfo.Value = dropdownValue
						}
						index = index + 1
					}
					recordInfo.Data = dropDownArray
					recordInfo.IsDynamic = recordSchema.IsDynamic
					recordInfo.DynamicMappingField = recordSchema.DynamicMappingField
					recordInfo.DynamicComponent = recordSchema.DynamicComponent
				} else if recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.

					valueArray := recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Values
					for _, queryResult := range queryResults {
						//id := int(queryResult[recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Index].(int32))
						for _, dropDownValueField := range valueArray {
							fmt.Println("queryResult[dropDownValueField]:", queryResult[dropDownValueField])
							//dropdownValue := queryResult[dropDownValueField]
							recordInfo.ValueArray = append(recordInfo.ValueArray, queryResult)
						}
						recordInfo.IndexArray = make([]int, 0)
					}
					recordInfo.Data = queryResults
				} else if recordSchema.ResponseObjectMapping.Builder.SingleDropdownObject != nil {
					// we need to order based on id starting from low.

					valueArray := recordSchema.ResponseObjectMapping.Builder.SingleDropdownObject.Values
					var indexId int
					indexId = -1
					for index, queryResult := range queryResults {
						if index == 0 {
							indexId = 1
							var valurMap = make(map[string]interface{})
							for _, dropDownValueField := range valueArray {
								fmt.Println("queryResult[dropDownValueField]:", queryResult[dropDownValueField])
								valurMap[dropDownValueField] = queryResult[dropDownValueField]
							}
							recordInfo.Value = valurMap
						}

					}

					recordInfo.Index = indexId
					recordInfo.Data = queryResults
				} else if recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					recordInfo.IndexArray = make([]int, 0)
					var dropDownArray []component.OrderedData
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := queryResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.ResponseObjectMapping.Builder.Table != nil {
					header := recordSchema.ResponseObjectMapping.Builder.Table.Schema

					recordInfo.Header = header
					if recordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink != "" {
						// we have configured the common route link, so send that
						recordInfo.CommonRouteLink = recordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink
					}
					// iterate through each results
					//for _, results := range queryResults {
					//	// we need to check the schema for any route link
					//	for index, individualHeader := range header {
					//		if individualHeader.RouteEnabled {
					//			recordIdQuery := individualHeader.RecordIdQuery
					//			recordIdQuery = strings.Replace(recordIdQuery, "["+individualHeader.Property+"]", results[individualHeader.Property].(string), 1)
					//			var linkedRecordId int
					//
					//			dbConnection.Raw(recordIdQuery).Scan(&linkedRecordId)
					//			if linkedRecordId == 0 {
					//				// no record found, so we should make the header route link false, and the id -
					//				results[individualHeader.RouteRecordIdProperty] = -1
					//				header[index].RouteEnabled = false
					//			} else {
					//				results[individualHeader.RouteRecordIdProperty] = -1
					//			}
					//			header[index].RecordIdQuery = ""
					//		}
					//	}
					//	recordInfo.Data = queryResults
					//}
				}
			} else if recordSchema.ResponseObjectMapping.Service != nil {
				authService := GetService(recordSchema.ResponseObjectMapping.Service.Name).ServiceInterface.(AuthInterface)
				method := reflect.ValueOf(authService).MethodByName(recordSchema.ResponseObjectMapping.Service.Call)

				inputs := make([]reflect.Value, len(recordSchema.ResponseObjectMapping.Service.ServiceParam))
				var queryResults []datatypes.JSON
				if method.IsValid() {
					returnValues := method.Call(inputs)
					queryResults = returnValues[0].Interface().([]datatypes.JSON)

				}
				if recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.

					existingIds := make([]int, 0)
					//valueArray := individualRecordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Values

					recordInfo.IndexArray = existingIds
					recordInfo.Data = queryResults
					recordInfo.InterfaceType = "multiDropdownObject"
				}
			}

		} else if recordSchema.Type == "object_array" {
			var header = make([]component.TableSchema, 0)
			for _, internalRecordSchema := range recordSchema.RecordSchema {
				headerObject := component.TableSchema{}
				headerObject.Type = internalRecordSchema.Type
				headerObject.Display = internalRecordSchema.Display
				headerObject.GridSystem = internalRecordSchema.GridSystem
				headerObject.Label = internalRecordSchema.Label
				headerObject.Property = internalRecordSchema.Property
				headerObject.InterfaceType = internalRecordSchema.InterfaceType
				headerObject.Render = internalRecordSchema.Render
				headerObject.Name = internalRecordSchema.Name
				headerObject.LinkedProperty = internalRecordSchema.LinkedProperty
				headerObject.LinkedDataType = internalRecordSchema.LinkedDataType
				if internalRecordSchema.HeaderFontColorObjectMapping != nil {
					if internalRecordSchema.HeaderFontColorObjectMapping.Predefined != nil {
						predefinedResults := internalRecordSchema.HeaderFontColorObjectMapping.Predefined.Data
						var objectKeyValues = make(map[string]string, len(predefinedResults))
						if internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue != nil {
							for _, result := range predefinedResults {
								objectKeyValues[result[internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Key].(string)] = result[internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Value].(string)
							}
						}

						raw, _ := json.Marshal(objectKeyValues)
						headerObject.FontColorList = raw
					}
				}
				if internalRecordSchema.HeaderObjectMapping != nil {
					if internalRecordSchema.HeaderObjectMapping.Predefined != nil {
						predefinedResults := internalRecordSchema.HeaderObjectMapping.Predefined.Data
						var objectKeyValues = make(map[string]string, len(predefinedResults))
						if internalRecordSchema.HeaderObjectMapping.Builder.KeyValue != nil {
							for _, result := range predefinedResults {
								objectKeyValues[result[internalRecordSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[internalRecordSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
							}
						}

						raw, _ := json.Marshal(objectKeyValues)
						headerObject.ObjectList = raw
					}
				}

				header = append(header, headerObject)

			}
			recordInfo.Header = header
			recordInfo.Data = make([]interface{}, 0)
		} else if recordSchema.Type == "object" {
			type EmptyObject struct {
			}
			recordInfo.Data = EmptyObject{}
		}
		// indicating what is the id it is referring to
		recordInfo.IsEdit = recordSchema.IsEdit
		recordInfo.Type = recordSchema.Type
		if recordSchema.Default != nil {
			recordInfo.Value = recordSchema.Default
		}

		// we should deal with defaultObjectMapping
		if recordSchema.DefaultObjectMapping != nil {
			if recordSchema.DefaultObjectMapping.Query != nil {
				query := recordSchema.DefaultObjectMapping.Query.Query
				var queryResults []map[string]interface{}
				dbConnection.Raw(query).Scan(&queryResults)
				if recordSchema.DefaultObjectMapping.Builder.SingleValue != nil {
					recordInfo.Value = recordSchema.DefaultObjectMapping.Builder.SingleValue.Field
				}
			} else if recordSchema.DefaultObjectMapping.Predefined != nil {
				if recordSchema.DefaultObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0
					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						id := util.InterfaceToInt(predefinedResult[recordSchema.DefaultObjectMapping.Builder.SingleDropdown.Index])
						dropdownValue := predefinedResult[recordSchema.DefaultObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if index == 0 {
							recordInfo.Index = id
							recordInfo.Value = dropdownValue
						}
						index = index + 1
					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.DefaultObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					recordInfo.IndexArray = make([]int, 0)
					var dropDownArray []component.OrderedData
					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						id := util.InterfaceToInt(predefinedResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := predefinedResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.DefaultObjectMapping.Builder.TableFieldsToObjectArray != nil {

					// this will convert table fields into object
					var objectResponse []map[string]interface{}

					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						var individualResponse = make(map[string]interface{})
						for _, individualField := range recordSchema.DefaultObjectMapping.Builder.TableFieldsToObjectArray.Fields {
							dataType := individualField.Type
							fieldName := individualField.Name
							if dataType == "bool" {
								individualResponse[fieldName] = util.InterfaceToBool(predefinedResult[fieldName])
							} else if dataType == "int" {
								individualResponse[fieldName] = util.InterfaceToInt(predefinedResult[fieldName])
							} else if dataType == "double" {
								individualResponse[fieldName] = util.InterfaceToFloat(predefinedResult[fieldName])
							} else {
								individualResponse[fieldName] = predefinedResult[fieldName]
							}
						}
						objectResponse = append(objectResponse, individualResponse)
					}
					recordInfo.Data = objectResponse

				}
			} else {
				// in that case, check if any builders configured
				if recordSchema.DefaultObjectMapping.Builder.SingleValue != nil {
					recordInfo.Value = recordSchema.DefaultObjectMapping.Builder.SingleValue.Field
				}
			}

		}

		response[recordSchema.Property] = recordInfo

	}

	var emptyObject = make(map[string]interface{})
	serializedEmptyObject, _ := json.Marshal(emptyObject)
	cm.GenerateAdditionalRecordResponse(zone, dbConnection, serializedEmptyObject, componentName, response)

	return response

}

func (cm *ComponentManager) GetNewRecordResponse_v1(zone string, dbConnection *gorm.DB, componentName string, templateId int) map[string]interface{} {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	newRecordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	response := make(map[string]interface{}, len(newRecordSchema))
	// normally if the field is linked field, in the record, we will see the id
	for _, recordSchema := range newRecordSchema {
		recordInfo := component.RecordInfo{}
		if recordSchema.ResponseObjectMapping != nil {

			if recordSchema.ResponseObjectMapping.Query != nil {
				var queryResults []map[string]interface{}
				dbConnection.Raw(recordSchema.ResponseObjectMapping.Query.Query).Scan(&queryResults)
				if recordSchema.ResponseObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0

					for _, queryResult := range queryResults {
						id := int(queryResult[recordSchema.ResponseObjectMapping.Builder.SingleDropdown.Index].(int32))
						dropdownValue := queryResult[recordSchema.ResponseObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if index == 0 {
							recordInfo.Index = id
							recordInfo.Value = dropdownValue
						}
						index = index + 1
					}
					recordInfo.Data = dropDownArray
					recordInfo.IsDynamic = recordSchema.IsDynamic
					recordInfo.DynamicMappingField = recordSchema.DynamicMappingField
					recordInfo.DynamicComponent = recordSchema.DynamicComponent
				} else if recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.

					valueArray := recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Values
					for _, queryResult := range queryResults {
						//id := int(queryResult[recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Index].(int32))
						for _, dropDownValueField := range valueArray {
							fmt.Println("queryResult[dropDownValueField]:", queryResult[dropDownValueField])

							//dropdownValue := queryResult[dropDownValueField]
							recordInfo.ValueArray = append(recordInfo.ValueArray, queryResult)
						}
						recordInfo.IndexArray = make([]int, 0)
					}
					recordInfo.Data = queryResults
				} else if recordSchema.ResponseObjectMapping.Builder.SingleDropdownObject != nil {
					// we need to order based on id starting from low.

					valueArray := recordSchema.ResponseObjectMapping.Builder.SingleDropdownObject.Values
					var indexId int
					indexId = -1
					for index, queryResult := range queryResults {
						if index == 0 {
							indexId = 1
							var valurMap = make(map[string]interface{})
							for _, dropDownValueField := range valueArray {
								fmt.Println("queryResult[dropDownValueField]:", queryResult[dropDownValueField])
								valurMap[dropDownValueField] = queryResult[dropDownValueField]
							}
							recordInfo.Value = valurMap
						}

					}

					recordInfo.Index = indexId
					recordInfo.Data = queryResults
				} else if recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					recordInfo.IndexArray = make([]int, 0)
					var dropDownArray []component.OrderedData
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := queryResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.ResponseObjectMapping.Builder.Table != nil {
					header := recordSchema.ResponseObjectMapping.Builder.Table.Schema

					recordInfo.Header = header
					if recordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink != "" {
						// we have configured the common route link, so send that
						recordInfo.CommonRouteLink = recordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink
					}
				}
			} else if recordSchema.ResponseObjectMapping.Service != nil {
				authService := GetService(recordSchema.ResponseObjectMapping.Service.Name).ServiceInterface.(AuthInterface)
				method := reflect.ValueOf(authService).MethodByName(recordSchema.ResponseObjectMapping.Service.Call)

				inputs := make([]reflect.Value, len(recordSchema.ResponseObjectMapping.Service.ServiceParam))
				var queryResults []datatypes.JSON
				if method.IsValid() {
					returnValues := method.Call(inputs)
					queryResults = returnValues[0].Interface().([]datatypes.JSON)

				}
				if recordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.

					existingIds := make([]int, 0)
					//valueArray := individualRecordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Values

					recordInfo.IndexArray = existingIds
					recordInfo.Data = queryResults
					recordInfo.InterfaceType = "multiDropdownObject"
				}
			}

		} else if recordSchema.Type == "object_array" {
			var header = make([]component.TableSchema, 0)
			for _, internalRecordSchema := range recordSchema.RecordSchema {
				headerObject := component.TableSchema{}
				headerObject.Type = internalRecordSchema.Type
				headerObject.Display = internalRecordSchema.Display
				headerObject.GridSystem = internalRecordSchema.GridSystem
				headerObject.Label = internalRecordSchema.Label
				headerObject.Property = internalRecordSchema.Property
				headerObject.InterfaceType = internalRecordSchema.InterfaceType
				headerObject.Render = internalRecordSchema.Render
				headerObject.Name = internalRecordSchema.Name
				headerObject.LinkedProperty = internalRecordSchema.LinkedProperty
				headerObject.LinkedDataType = internalRecordSchema.LinkedDataType
				if internalRecordSchema.HeaderFontColorObjectMapping != nil {
					if internalRecordSchema.HeaderFontColorObjectMapping.Predefined != nil {
						predefinedResults := internalRecordSchema.HeaderFontColorObjectMapping.Predefined.Data
						var objectKeyValues = make(map[string]string, len(predefinedResults))
						if internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue != nil {
							for _, result := range predefinedResults {
								objectKeyValues[result[internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Key].(string)] = result[internalRecordSchema.HeaderFontColorObjectMapping.Builder.KeyValue.Value].(string)
							}
						}

						raw, _ := json.Marshal(objectKeyValues)
						headerObject.FontColorList = raw
					}
				}
				if internalRecordSchema.HeaderObjectMapping != nil {
					if internalRecordSchema.HeaderObjectMapping.Predefined != nil {
						predefinedResults := internalRecordSchema.HeaderObjectMapping.Predefined.Data
						var objectKeyValues = make(map[string]string, len(predefinedResults))
						if internalRecordSchema.HeaderObjectMapping.Builder.KeyValue != nil {
							for _, result := range predefinedResults {
								objectKeyValues[result[internalRecordSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[internalRecordSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
							}
						}

						raw, _ := json.Marshal(objectKeyValues)
						headerObject.ObjectList = raw
					}
				}

				header = append(header, headerObject)

			}
			recordInfo.Header = header
			recordInfo.Data = make([]interface{}, 0)
		} else if recordSchema.Type == "object" {
			type EmptyObject struct {
			}
			recordInfo.Data = EmptyObject{}
		}
		// indicating what is the id it is referring to
		recordInfo.IsEdit = recordSchema.IsEdit
		recordInfo.Type = recordSchema.Type
		if recordSchema.Default != nil {
			recordInfo.Value = recordSchema.Default
		}

		// we should deal with defaultObjectMapping
		if recordSchema.DefaultObjectMapping != nil {
			if recordSchema.DefaultObjectMapping.Query != nil {
				query := recordSchema.DefaultObjectMapping.Query.Query
				var queryResults []map[string]interface{}
				dbConnection.Raw(query).Scan(&queryResults)
				if recordSchema.DefaultObjectMapping.Builder.SingleValue != nil {
					recordInfo.Value = recordSchema.DefaultObjectMapping.Builder.SingleValue.Field
				}
			} else if recordSchema.DefaultObjectMapping.Predefined != nil {
				if recordSchema.DefaultObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0
					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						id := util.InterfaceToInt(predefinedResult[recordSchema.DefaultObjectMapping.Builder.SingleDropdown.Index])
						dropdownValue := predefinedResult[recordSchema.DefaultObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if index == 0 {
							recordInfo.Index = id
							recordInfo.Value = dropdownValue
						}
						index = index + 1
					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.DefaultObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					recordInfo.IndexArray = make([]int, 0)
					var dropDownArray []component.OrderedData
					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						id := util.InterfaceToInt(predefinedResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := predefinedResult[recordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.Data = dropDownArray
				} else if recordSchema.DefaultObjectMapping.Builder.TableFieldsToObjectArray != nil {

					// this will convert table fields into object
					var objectResponse []map[string]interface{}

					for _, predefinedResult := range recordSchema.DefaultObjectMapping.Predefined.Data {
						var individualResponse = make(map[string]interface{})
						for _, individualField := range recordSchema.DefaultObjectMapping.Builder.TableFieldsToObjectArray.Fields {
							dataType := individualField.Type
							fieldName := individualField.Name
							if dataType == "bool" {
								individualResponse[fieldName] = util.InterfaceToBool(predefinedResult[fieldName])
							} else if dataType == "int" {
								individualResponse[fieldName] = util.InterfaceToInt(predefinedResult[fieldName])
							} else if dataType == "double" {
								individualResponse[fieldName] = util.InterfaceToFloat(predefinedResult[fieldName])
							} else {
								individualResponse[fieldName] = predefinedResult[fieldName]
							}
						}
						objectResponse = append(objectResponse, individualResponse)
					}
					recordInfo.Data = objectResponse

				}
			} else {
				// in that case, check if any builders configured
				if recordSchema.DefaultObjectMapping.Builder.SingleValue != nil {
					recordInfo.Value = recordSchema.DefaultObjectMapping.Builder.SingleValue.Field
				}
			}

		}

		response[recordSchema.Property] = recordInfo

	}

	// we are going to add any template fields configured
	if cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent != nil {
		referenceComponentName := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent.ReferenceComponent
		templateFieldName := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent.TemplateFieldName
		referenceTable := cm.GetTargetTable(referenceComponentName) //it_service_category_template
		err, templateTableFieldsObject := Get(dbConnection, referenceTable, templateId)
		if err == nil {

			// got all the template fields
			templateFields := GetObjectFields(templateTableFieldsObject.ObjectInfo)
			templateFieldRecords := templateFields[templateFieldName]
			serialisedRecords := GetInterfaceToSerialisation(templateFieldRecords)
			var listOfTableRecords []TemplateRecords
			err = json.Unmarshal(serialisedRecords, &listOfTableRecords)

			//We have to sort the field then send new record object
			//sort.Slice(listOfTableRecords, func(i, j int) bool {
			//	return listOfTableRecords[i].Id < listOfTableRecords[j].Id
			//})

			// we got the fields
			var listOfTemplateFields = make(map[string]component.RecordInfo)
			for _, templateRecord := range listOfTableRecords {
				// only add this fields enabled as is not dropped down condition field
				if templateRecord.IsDroppedDownConditionalField {
					recordInfo := component.RecordInfo{}
					recordInfo.Property = templateRecord.Property
					recordInfo.Type = getInt2DataType(templateRecord.DataType)
					recordInfo.Label = templateRecord.Label
					recordInfo.GridSystem = getGridSystem2Str(templateRecord.GridSystem)
					recordInfo.Render = true
					recordInfo.IsDynamic = true
					recordInfo.IsMandatoryField = templateRecord.IsMandatoryField
					recordInfo.Description = templateRecord.Description
					recordInfo.InterfaceType = getInterfaceType2Str(templateRecord.InterfaceTypeList)
					recordInfo.InterfaceField = getInt2FieldType(templateRecord.InterfaceFieldList)
					listOfTemplateFields[templateRecord.Property] = recordInfo
				}
			}
			var condFieldMapProperty = make(map[string][]component.RecordInfo)
			for _, templateRecord := range listOfTableRecords {
				if templateRecord.InterfaceTypeList == 5 || templateRecord.InterfaceTypeList == 6 {
					if len(templateRecord.DynamicDroppedDownAttributes.ConditionalFields) > 0 {
						conditionalFields := templateRecord.DynamicDroppedDownAttributes.ConditionalFields
						for _, condField := range conditionalFields {
							existingRecordInfo := listOfTemplateFields[condField.Property]
							if condValue, ok := condFieldMapProperty[condField.Value]; ok {
								// already here
								condValue = append(condValue, existingRecordInfo)
								condFieldMapProperty[condField.Value] = condValue
							} else {
								// already not here
								var listOfRecordInfo []component.RecordInfo
								listOfRecordInfo = append(listOfRecordInfo, existingRecordInfo)
								condFieldMapProperty[condField.Value] = listOfRecordInfo
							}

						}
					}

				}
			}

			var listOfDynamicRecords []component.RecordInfo
			for _, templateRecord := range listOfTableRecords {
				// only add this fields enabled as is not dropped down condition field
				if !templateRecord.IsDroppedDownConditionalField {

					recordInfo := component.RecordInfo{}
					recordInfo.Property = templateRecord.Property
					recordInfo.Type = getInt2DataType(templateRecord.DataType)
					recordInfo.Label = templateRecord.Label
					recordInfo.GridSystem = getGridSystem2Str(templateRecord.GridSystem)
					// initially don't render field after entity level configured
					if templateRecord.EnabledAfterWorkflowStatusLevel == nil {
						recordInfo.Render = true
					} else {
						if templateRecord.EnabledAfterWorkflowStatusLevel != nil {
							if *templateRecord.EnabledAfterWorkflowStatusLevel > 1 {
								recordInfo.Render = false
							} else {
								recordInfo.Render = true
							}
						}

					}
					recordInfo.IsDynamic = true
					recordInfo.IsMandatoryField = templateRecord.IsMandatoryField
					recordInfo.Description = templateRecord.Description
					recordInfo.InterfaceType = getInterfaceType2Str(templateRecord.InterfaceTypeList)
					recordInfo.InterfaceField = getInt2FieldType(templateRecord.InterfaceFieldList)
					if templateRecord.InterfaceTypeList == 5 || templateRecord.InterfaceTypeList == 6 {

						// if it is drooped down, use the dropped attributes to load the data
						if templateRecord.DynamicDroppedDownAttributes.Type == 2 {
							// auto fields are needed to generate the output

							autoFieldSource := templateRecord.DynamicDroppedDownAttributes.AutoFieldsSource

							fieldSourceMap := strings.Split(autoFieldSource, ".")

							if len(fieldSourceMap) == 3 {
								sourceComponent := fieldSourceMap[1]
								sourceField := fieldSourceMap[2]
								sourceTargetTable := cm.GetTargetTable(sourceComponent)
								err, listOfSourceObjects := GetObjects(dbConnection, sourceTargetTable)

								if err == nil {
									var dropDownArray []component.OrderedData
									index := 0

									for _, sourceObject := range listOfSourceObjects {
										id := sourceObject.Id
										var objectFields = make(map[string]interface{})
										json.Unmarshal(sourceObject.ObjectInfo, &objectFields)
										if value, ok := objectFields[sourceField]; ok {

											// now check the value is
											if conValue, ok := condFieldMapProperty[util.InterfaceToString(value)]; ok {
												dropDownArray = append(dropDownArray, component.OrderedData{
													Id:                       index + 1,
													Value:                    util.InterfaceToString(value),
													OnValueConditionalFields: conValue,
												})
											} else {
												dropDownArray = append(dropDownArray, component.OrderedData{
													Id:    index + 1,
													Value: util.InterfaceToString(value),
												})
											}

											if index == 0 {
												recordInfo.Index = id
												recordInfo.Value = util.InterfaceToString(value)
											}
											index = index + 1
										}

									}
									recordInfo.Data = dropDownArray
								}
							}
						}
						if templateRecord.DynamicDroppedDownAttributes.Type == 1 {
							var dropDownArray []component.OrderedData
							for index, valueDrop := range templateRecord.DynamicDroppedDownAttributes.ManualFieldsSource {
								if conValue, ok := condFieldMapProperty[util.InterfaceToString(valueDrop)]; ok {
									dropDownArray = append(dropDownArray, component.OrderedData{
										Id:                       index + 1,
										Value:                    util.InterfaceToString(valueDrop),
										OnValueConditionalFields: conValue,
									})
								} else {
									dropDownArray = append(dropDownArray, component.OrderedData{
										Id:    index + 1,
										Value: util.InterfaceToString(valueDrop),
									})
								}

								if index == 0 {
									recordInfo.Index = index + 1
									recordInfo.Value = util.InterfaceToString(valueDrop)
								}
								index = index + 1
							}
							recordInfo.Data = dropDownArray
						}
					}
					recordInfo.Display = true
					listOfDynamicRecords = append(listOfDynamicRecords, recordInfo)
				}
			}
			response["dynamicFields"] = listOfDynamicRecords
		}

	}

	var emptyObject = make(map[string]interface{})
	serializedEmptyObject, _ := json.Marshal(emptyObject)
	cm.GenerateAdditionalRecordResponse(zone, dbConnection, serializedEmptyObject, componentName, response)

	return response

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

func (cm *ComponentManager) GetCardViewArrayOfMapInterface(listOfObjects *[]component.GeneralObject, componentName string) []map[string]interface{} {
	var dataRecords map[string]interface{}
	var cardViewResponseMap []map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var mappingFields = make(map[string]string, len(componentSchema.CardViewSchema.Field))

	var isAnyArray bool
	for _, mappingField := range componentSchema.CardViewSchema.Field {
		mappingFields[mappingField.Property] = mappingField.MappingField
		splitFields := strings.Split(mappingField.Property, ".")
		if len(splitFields) > 1 {
			isAnyArray = true
		}
	}
	var internalRecords []map[string]interface{}
	if isAnyArray {
		// first generate all the array fields
		for _, records := range *listOfObjects {
			json.Unmarshal(records.ObjectInfo, &dataRecords)
			arrayOfObjects := dataRecords["taskInfo"]
			serializedData, _ := json.Marshal(arrayOfObjects)
			json.Unmarshal(serializedData, &internalRecords)
			for _, internalElement := range internalRecords {
				var individualCardViewResponse = make(map[string]interface{})
				for property, mapField := range mappingFields {
					splitFields := strings.Split(property, ".")
					if len(splitFields) == 1 {
						individualCardViewResponse[mapField] = dataRecords[property]
					} else {
						var internalElementRecord map[string]interface{}
						sri, _ := json.Marshal(internalElement)
						json.Unmarshal(sri, &internalElementRecord)
						secondElement := splitFields[1]
						individualCardViewResponse[mapField] = internalElementRecord[secondElement]
					}
				}
				individualCardViewResponse["id"] = records.Id
				cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

			}
		}
	} else {
		cardViewResponse := component.CardViewResponse{}
		cardViewResponse.Template = componentSchema.CardViewSchema.Template
		for _, records := range *listOfObjects {
			json.Unmarshal(records.ObjectInfo, &dataRecords)
			var individualCardViewResponse = make(map[string]interface{}, len(mappingFields)+1)
			for property, mapField := range mappingFields {
				individualCardViewResponse[mapField] = dataRecords[property]

			}
			individualCardViewResponse["id"] = records.Id
			cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

		}

	}

	return cardViewResponseMap
}

func (cm *ComponentManager) GetCardViewArrayOfMapInterfaceV1(listOfObjects []component.GeneralObject, componentName string) []map[string]interface{} {
	var dataRecords map[string]interface{}
	var cardViewResponseMap []map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var mappingFields = make(map[string]string, len(componentSchema.CardViewSchema.Field))

	var isAnyArray bool
	for _, mappingField := range componentSchema.CardViewSchema.Field {
		mappingFields[mappingField.Property] = mappingField.MappingField
		splitFields := strings.Split(mappingField.Property, ".")
		if len(splitFields) > 1 {
			isAnyArray = true
		}
	}
	var internalRecords []map[string]interface{}
	if isAnyArray {
		// first generate all the array fields
		for _, records := range listOfObjects {
			json.Unmarshal(records.ObjectInfo, &dataRecords)
			arrayOfObjects := dataRecords["taskInfo"]
			serializedData, _ := json.Marshal(arrayOfObjects)
			json.Unmarshal(serializedData, &internalRecords)
			for _, internalElement := range internalRecords {
				var individualCardViewResponse = make(map[string]interface{})
				for property, mapField := range mappingFields {
					splitFields := strings.Split(property, ".")
					if len(splitFields) == 1 {
						individualCardViewResponse[mapField] = dataRecords[property]
					} else {
						var internalElementRecord map[string]interface{}
						sri, _ := json.Marshal(internalElement)
						json.Unmarshal(sri, &internalElementRecord)
						secondElement := splitFields[1]
						individualCardViewResponse[mapField] = internalElementRecord[secondElement]
					}
				}
				individualCardViewResponse["id"] = records.Id
				cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

			}
		}
	} else {
		cardViewResponse := component.CardViewResponse{}
		cardViewResponse.Template = componentSchema.CardViewSchema.Template
		for _, records := range listOfObjects {
			json.Unmarshal(records.ObjectInfo, &dataRecords)
			var individualCardViewResponse = make(map[string]interface{}, len(mappingFields)+1)
			for property, mapField := range mappingFields {
				individualCardViewResponse[mapField] = dataRecords[property]

			}
			individualCardViewResponse["id"] = records.Id
			cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

		}

	}

	return cardViewResponseMap
}

func (cm *ComponentManager) GetCardViewFromListOfInterface(listOfComponents []interface{}, componentName string) []interface{} {
	var dataRecords map[string]interface{}
	var cardViewResponseArray []interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var mappingFields = make(map[string]string, len(componentSchema.CardViewSchema.Field))

	for _, mappingField := range componentSchema.CardViewSchema.Field {
		mappingFields[mappingField.Property] = mappingField.MappingField
	}
	cardViewResponse := component.CardViewResponse{}
	cardViewResponse.Template = componentSchema.CardViewSchema.Template
	for _, recordInterface := range listOfComponents {
		componentObject := component.GeneralObject{ObjectInfo: recordInterface.(datatypes.JSON)}
		json.Unmarshal(componentObject.ObjectInfo, &dataRecords)
		var individualCardViewResponse = make(map[string]interface{}, len(mappingFields)+1)
		for property, mapField := range mappingFields {
			individualCardViewResponse[mapField] = dataRecords[property]
		}
		individualCardViewResponse["id"] = componentObject.Id
		cardViewResponseArray = append(cardViewResponseArray, individualCardViewResponse)

	}

	return cardViewResponseArray
}
func (cm *ComponentManager) GetCardViewResponse(listOfComponents *[]component.GeneralObject, componentName string) component.CardViewResponse {
	var dataRecords map[string]interface{}
	var cardViewResponseMap []map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var mappingFields = make(map[string]string, len(componentSchema.CardViewSchema.Field))

	for _, mappingField := range componentSchema.CardViewSchema.Field {
		mappingFields[mappingField.Property] = mappingField.MappingField
	}
	cardViewResponse := component.CardViewResponse{}
	cardViewResponse.Template = componentSchema.CardViewSchema.Template
	for _, records := range *listOfComponents {
		json.Unmarshal(records.ObjectInfo, &dataRecords)
		var individualCardViewResponse = make(map[string]interface{}, len(mappingFields)+1)
		for property, mapField := range mappingFields {
			individualCardViewResponse[mapField] = dataRecords[property]
		}

		// Added default image
		if val, ok := individualCardViewResponse["field2"]; ok {
			if val == nil {
				individualCardViewResponse["field2"] = "-"
			}
		}

		individualCardViewResponse["id"] = records.Id
		cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

	}
	cardsResponse, _ := json.Marshal(cardViewResponseMap)
	cardViewResponse.Cards = cardsResponse
	return cardViewResponse
}

func (cm *ComponentManager) GetCardViewResponseV1(listOfComponents []component.GeneralObject, componentName string) component.CardViewResponse {
	var dataRecords map[string]interface{}
	var cardViewResponseMap []map[string]interface{}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var mappingFields = make(map[string]string, len(componentSchema.CardViewSchema.Field))

	for _, mappingField := range componentSchema.CardViewSchema.Field {
		mappingFields[mappingField.Property] = mappingField.MappingField
	}
	cardViewResponse := component.CardViewResponse{}
	cardViewResponse.Template = componentSchema.CardViewSchema.Template
	for _, records := range listOfComponents {
		json.Unmarshal(records.ObjectInfo, &dataRecords)
		var individualCardViewResponse = make(map[string]interface{}, len(mappingFields)+1)
		for property, mapField := range mappingFields {
			individualCardViewResponse[mapField] = dataRecords[property]
		}

		// Added default image
		if val, ok := individualCardViewResponse["field2"]; ok {
			if val == nil {
				individualCardViewResponse["field2"] = "-"
			}
		}

		individualCardViewResponse["id"] = records.Id
		cardViewResponseMap = append(cardViewResponseMap, individualCardViewResponse)

	}
	cardsResponse, _ := json.Marshal(cardViewResponseMap)
	cardViewResponse.Cards = cardsResponse
	return cardViewResponse
}

func (cm *ComponentManager) ExportDataFromQueryResults(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, queryResults []component.GeneralObject) (error, int, component.ExportDataResponse) {
	// first select all the fields based on given fields
	// first lets fix the flat query, and then we will focus on children query
	int64ComponentId := cm.ComponentNameIdMapping[componentName]

	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)
	var savedFileName string
	var err error

	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsv(dbConnection, queryResults, childMasterData, parentChildMap, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// no error
		contentResponse := make(map[string]interface{})
		//contentResponse := component.ContentResponse{}
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("export file has been successfully created, but system error during uploading to content service"), UnableToCreateExportFile, exportDataResponse
	}
	return nil, 0, exportDataResponse
}

func (cm *ComponentManager) ExportDataProductionOrder(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, condition string) (error, int, component.ExportDataResponse) {
	// first select all the fields based on given fields
	// first lets fix the flat query, and then we will focus on children query
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	targetTable := componentSchema.TargetTable
	var queryResults []component.GeneralObject
	var selectQuery string
	if condition == "" {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " ORDER BY id DESC"
	} else {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " where " + condition + " ORDER BY id DESC"
	}

	dbConnection.Raw(selectQuery).Scan(&queryResults)

	// fmt.Println(childMasterData)

	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)
	var savedFileName string
	var err error

	// write the custom extraction
	var newQueryObject = make([]component.GeneralObject, 0)

	for _, records := range queryResults {
		var modifiedObject datatypes.JSON
		var isFieldFound bool
		objectFields := make(map[string]interface{})

		// Parse the object once
		json.Unmarshal(records.ObjectInfo, &objectFields)

		// Create a temporary map to hold any new fields
		newFields := make(map[string]interface{})

		for _, fieldName := range exportFields.Data {
			fmt.Println("checking field name: ", fieldName)
			if fieldName.Data == "partNumber" {
				eventSourceId := util.InterfaceToInt(objectFields["eventSourceId"])
				query := "select object_info->>'$.partNumber' as partNumber from part_master where id = (select object_info->>'$.partNumber' from production_order_master where id = " + strconv.Itoa(eventSourceId) + ")"
				var partNumber string
				dbConnection.Raw(query).Scan(&partNumber)

				// Add partNumber to newFields map
				newFields["partNumber"] = partNumber
				isFieldFound = true
				fmt.Println("found the part number field name: ", partNumber)
			}
			if fieldName.Data == "cycleTime" {
				eventSourceId := util.InterfaceToInt(objectFields["eventSourceId"])
				query := "select object_info->>'$.cycleTime' as cycleTime from production_order_master where id = " + strconv.Itoa(eventSourceId)
				var cycleTime string
				dbConnection.Raw(query).Scan(&cycleTime)

				// Add cycleTime to newFields map
				newFields["cycleTime"] = cycleTime
				isFieldFound = true
				fmt.Println("found the cycleTime field name: ", cycleTime)
			}
		}

		// If we found fields to update, apply them to the original object
		if isFieldFound {
			for key, value := range newFields {
				objectFields[key] = value
			}
			// Marshal updated fields into modifiedObject
			updatedJSON, _ := json.Marshal(objectFields)
			modifiedObject = datatypes.JSON(updatedJSON)
		} else {
			modifiedObject = records.ObjectInfo
		}

		newQueryObject = append(newQueryObject, component.GeneralObject{
			Id:         records.Id,
			ObjectInfo: modifiedObject,
		})
	}

	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsv(dbConnection, newQueryObject, childMasterData, parentChildMap, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, newQueryObject, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// no error
		contentResponse := make(map[string]interface{})
		//contentResponse := component.ContentResponse{}
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("export file has been successfully created, but system error during uploading to content service"), UnableToCreateExportFile, exportDataResponse
	}
	return nil, 0, exportDataResponse
}

func (cm *ComponentManager) ExportData(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, condition string) (error, int, component.ExportDataResponse) {
	// first select all the fields based on given fields
	// first lets fix the flat query, and then we will focus on children query
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	targetTable := componentSchema.TargetTable
	var queryResults []component.GeneralObject
	var selectQuery string
	if condition == "" {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " ORDER BY id DESC"
	} else {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " where " + condition + " ORDER BY id DESC"
	}

	dbConnection.Raw(selectQuery).Scan(&queryResults)

	childMasterData, parentChildMap := getChildrenDataV2(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)
	var savedFileName string
	var err error

	if exportFields.Format == "csv" {
		_, savedFileName, err = writeInCsvV3(cm, dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// no error
		contentResponse := make(map[string]interface{})
		//contentResponse := component.ContentResponse{}
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("export file has been successfully created, but system error during uploading to content service"), UnableToCreateExportFile, exportDataResponse
	}
	return nil, 0, exportDataResponse
}

func (cm *ComponentManager) ITServiceExportData(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, condition string) (error, int, component.ExportDataResponse) {
	// first select all the fields based on given fields
	// first lets fix the flat query, and then we will focus on children query
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	targetTable := componentSchema.TargetTable
	var queryResults []component.GeneralObject
	var selectQuery string
	if condition == "" {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " ORDER BY id DESC"
	} else {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " where " + condition + " ORDER BY id DESC"
	}

	dbConnection.Raw(selectQuery).Scan(&queryResults)

	fmt.Println(exportFields.Data)

	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)
	var savedFileName string
	var err error

	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsv(dbConnection, queryResults, childMasterData, parentChildMap, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// no error
		contentResponse := make(map[string]interface{})
		//contentResponse := component.ContentResponse{}
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("export file has been successfully created, but system error during uploading to content service"), UnableToCreateExportFile, exportDataResponse
	}
	return nil, 0, exportDataResponse
}
func (cm *ComponentManager) GeneralExportData(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, condition string) (error, int, component.ExportDataResponse) {
	// first select all the fields based on given fields
	// first lets fix the flat query, and then we will focus on children query
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	targetTable := componentSchema.TargetTable
	var queryResults []component.GeneralObject
	var selectQuery string
	if condition == "" {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " ORDER BY id DESC"
	} else {
		selectQuery = "SELECT id, object_info FROM " + targetTable + " where " + condition + " ORDER BY id DESC"
	}

	dbConnection.Raw(selectQuery).Scan(&queryResults)

	fmt.Println(exportFields.Data)

	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)
	var savedFileName string
	var err error

	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsv(dbConnection, queryResults, childMasterData, parentChildMap, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// no error
		contentResponse := make(map[string]interface{})
		//contentResponse := component.ContentResponse{}
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("export file has been successfully created, but system error during uploading to content service"), UnableToCreateExportFile, exportDataResponse
	}
	return nil, 0, exportDataResponse
}
func (cm *ComponentManager) ExportDataV3(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, queryResults []component.GeneralObject) (error, int, component.ExportDataResponse) {
	// Retrieve child master data and parent-child mapping
	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)

	var savedFileName string
	var err error
	// Determine the format and create the export file accordingly
	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsv(dbConnection, queryResults, childMasterData, parentChildMap, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		//_, savedFileName, err = writeInExcel(cm, dbConnection, queryResults, childMasterData, parentChildMap, exportFields)
	} else {
		UnsupportedExportFormat := 1003
		return errors.New("unsupported export format"), UnsupportedExportFormat, component.ExportDataResponse{}
	}

	// Handle any errors during file creation
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, component.ExportDataResponse{}
	}

	// Upload the file
	exportDataResponse := component.ExportDataResponse{}

	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		// Handle successful file upload
		contentResponse := make(map[string]interface{})
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("unable to upload the export file"), UnableToCreateExportFile, exportDataResponse
	}

	return nil, 0, exportDataResponse
}

func (cm *ComponentManager) ExportDataV4(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, condition string) (error, int, component.ExportDataResponse) {
	// First, select all the fields based on given fields
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	targetTable := componentSchema.TargetTable

	var queryResults []component.GeneralObject
	selectQuery := "SELECT id, object_info FROM " + targetTable

	// Add condition if provided
	if condition != "" {
		selectQuery += " WHERE " + condition
	}

	selectQuery += " ORDER BY id DESC"
	dbConnection.Raw(selectQuery).Scan(&queryResults)

	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)

	var savedFileName string
	var err error
	if exportFields.Format == "csv" {
		_, savedFileName, err = writeInCsvV3(cm, dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields, componentName)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = writeInExcelV2(cm, dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	}

	exportDataResponse := component.ExportDataResponse{}
	fmt.Println("export errot", err)
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, exportDataResponse
	}

	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)

	if err == nil && httpResponse.StatusCode == 200 {
		fmt.Println("raw response :", string(rawResponse))
		contentResponse := make(map[string]interface{})
		json.Unmarshal(rawResponse, &contentResponse)

		exportDataResponse.Url = util.InterfaceToString(contentResponse["url"])
		exportDataResponse.Size = util.InterfaceToString(contentResponse["size"])
		exportDataResponse.Name = util.InterfaceToString(contentResponse["name"])
	} else {
		return errors.New("unable to create an export file"), UnableToCreateExportFile, exportDataResponse
	}

	return nil, 0, exportDataResponse
}

// ExportDataV2 function integrated for pass filtered value from it_service_request table export option
func (cm *ComponentManager) ExportDataV2(dbConnection *gorm.DB, componentName string, exportFields component.ExportDataCommand, queryResults []component.GeneralObject) (error, int, component.ExportDataResponse) {
	// Retrieve component information
	int64ComponentId := cm.ComponentNameIdMapping[componentName]

	// Extract serviceStatus IDs
	serviceStatusIDs := make([]int, 0)
	for _, result := range queryResults {
		var objectInfo map[string]interface{}
		if err := json.Unmarshal([]byte(result.ObjectInfo), &objectInfo); err == nil {
			if serviceStatusID, ok := objectInfo["serviceStatus"].(float64); ok {
				serviceStatusIDs = append(serviceStatusIDs, int(serviceStatusID))
			}
		}
	}

	// Remove duplicates
	uniqueServiceStatusIDs := unique(serviceStatusIDs)

	// Fetch status values
	var statusResults []struct {
		ID     int    `gorm:"column:id"`
		Status string `gorm:"column:status"`
	}
	statusQuery := "SELECT id, object_info->>'$.status' as status FROM it_service_request_status WHERE id IN (?)"
	if err := dbConnection.Raw(statusQuery, uniqueServiceStatusIDs).Scan(&statusResults).Error; err != nil {
		QueryExecutionFailed := 1003
		return errors.New("failed to fetch service statuses"), QueryExecutionFailed, component.ExportDataResponse{}
	}

	// Map statuses to IDs
	statusMap := make(map[int]string)
	for _, status := range statusResults {
		statusMap[status.ID] = status.Status
	}

	// Enrich main query results with statuses
	for i, result := range queryResults {
		var objectInfo map[string]interface{}
		if err := json.Unmarshal([]byte(result.ObjectInfo), &objectInfo); err == nil {
			if serviceStatusID, ok := objectInfo["serviceStatus"].(float64); ok {
				if status, found := statusMap[int(serviceStatusID)]; found {
					objectInfo["serviceStatus"] = status
				}
				newObjectInfo, _ := json.Marshal(objectInfo)
				queryResults[i].ObjectInfo = datatypes.JSON(newObjectInfo)
			}
		}

	}

	// Retrieve child data
	childMasterData, parentChildMap := getChildrenData(exportFields.Data, componentName, cm, dbConnection)

	// Prepare export schema
	arrayOfExportSchema := cm.ComponentSchema[int64ComponentId].ExportSchema
	arrayOfExportSchema = addIndexExportSchema(arrayOfExportSchema)

	// Prepare file name and write data based on the export format
	var savedFileName string
	var err error
	if exportFields.Format == "csv" {
		_, savedFileName, err = cm.writeInCsvV2(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)
	} else if exportFields.Format == "excel" {
		_, savedFileName, err = cm.writeInExcel(dbConnection, queryResults, childMasterData, parentChildMap, arrayOfExportSchema, exportFields)

	}

	// Handle file creation errors
	if err != nil {
		return errors.New("unable to create file"), UnableToReadCSVFile, component.ExportDataResponse{}
	}

	// Upload the file and handle the response
	var params map[string]string
	err, httpResponse, rawResponse := util.FileUploadFromDisk(cm.ComponentContentConfig.UpStream, params, "file", savedFileName)
	if err == nil && httpResponse.StatusCode == 200 {
		// Parse content response
		contentResponse := make(map[string]interface{})
		json.Unmarshal(rawResponse, &contentResponse)

		// Prepare export data response
		exportDataResponse := component.ExportDataResponse{
			Url:  util.InterfaceToString(contentResponse["url"]),
			Size: util.InterfaceToString(contentResponse["size"]),
			Name: util.InterfaceToString(contentResponse["name"]),
		}
		return nil, 0, exportDataResponse
	} else {
		return errors.New("unable to create export file"), UnableToCreateExportFile, component.ExportDataResponse{}
	}
}

// Helper function to remove duplicates from a slice of integers
func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// findExportDataInArray searches for an export schema with a given ID in an array of export schemas.
// If the schema is found, it is returned; otherwise, an empty export schema is returned.
func findExportDataInArray(arrayOfExportSchema []component.ExportSchema, serachId float32) component.ExportSchema {
	for _, exportField := range arrayOfExportSchema {
		if exportField.Id == serachId {
			return exportField
		}

		if len(exportField.Children) > 0 {
			for _, childrenField := range exportField.Children {
				if childrenField.Id == serachId {
					return component.ExportSchema{
						Data:  childrenField.Data,
						Label: childrenField.Label,
						Id:    childrenField.Id,
					}
				}
			}
		}
	}

	return component.ExportSchema{}
}

// writeInCsvV2 creates a CSV file from the query results and child master data according to the specified export fields.
func (cm *ComponentManager) writeInCsvV2(dbConnection *gorm.DB, queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, arrayOfExportSchema []component.ExportSchema, exportFields component.ExportDataCommand) (*os.File, string, error) {
	var header string
	for index, exportField := range exportFields.Data {
		if index == len(exportFields.Data)-1 {
			header = header + exportField.Label
		} else {
			header = header + exportField.Label + ","
		}
	}
	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)
	if err != nil {
		return nil, "cannot create file", err
	}

	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, uuid.New().String()+".csv")

	file, err := os.Create(savedFileName)
	fmt.Println("savedFileName", savedFileName)
	if err != nil {

		fmt.Println("error", err)
		return file, savedFileName, err
		// return err, 0, component.ExportDataResponse{}
	}
	defer file.Close()
	file.Write([]byte(header + "\n"))
	for _, generalObject := range queryResults {
		var keyValueData map[string]interface{}

		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		var composedString string
		for _, exportField := range exportFields.Data {
			originalExportSchema := findExportDataInArray(arrayOfExportSchema, exportField.Id)
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")

			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				//childMaster [Id:1, value:childInfo]
				childMaster := childMasterData[float32(numericId)]

				//cateory
				parentLabel := parentChildMap[fmt.Sprintf("%v", exportField.Id)]

				//Get foreign key from master data 1
				foreignKey := keyValueData[parentLabel]

				//Pick related row from parent master
				rowData := childMaster[util.InterfaceToInt(foreignKey)]
				if parentLabel != "eventId" {
				}

				if rowData == nil {
					composedString = composedString + "" + ","
				} else {
					composedString = composedString + util.InterfaceToString(rowData[exportField.Data]) + ","
				}

			} else {
				if exportField.LinkedMapFlag {
					if originalExportSchema.LinkedObjectMapping.Query != nil {
						query := originalExportSchema.LinkedObjectMapping.Query.Query
						for _, replacementField := range originalExportSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							} else {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							}
						}
						var linkedQueryResults []map[string]interface{}
						dbConnection.Raw(query).Scan(&linkedQueryResults)
						if originalExportSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(linkedQueryResults) > 0 {
								fieldValue := linkedQueryResults[0][originalExportSchema.LinkedObjectMapping.Builder.SingleValue.Field]
								composedString = composedString + util.InterfaceToString(fieldValue) + ","
							} else {
								composedString = composedString + "" + ","
							}
						}
					}
				} else {
					data := keyValueData[exportField.Data]
					if data == nil {
						composedString = composedString + "" + ","
					} else {
						// Check the data type
						_, timeErr := time.Parse(util.ISOTimeLayout, util.InterfaceToString(data))
						if timeErr == nil {
							composedString = composedString + util.ConvertTimeToTimeZonCorrectedFormat("Asia/Singapore", util.InterfaceToString(data)) + ","
						} else {
							result := strings.ReplaceAll(util.InterfaceToString(data), ",", "")
							composedString = composedString + result + ","
						}

					}

				}

			}
		}

		file.Write([]byte(composedString + "\n"))
	}
	return file, savedFileName, nil
}

func (cm *ComponentManager) writeInCsvV4(queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, exportFields component.ExportDataCommand) (*os.File, string, error) {
	var header string
	for index, exportField := range exportFields.Data {
		if index == len(exportFields.Data)-1 {
			header = header + exportField.Label
		} else {
			header = header + exportField.Label + ","
		}
	}
	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)
	if err != nil {
		return nil, "cannot create file", err
	}

	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, uuid.New().String()+".csv")

	file, err := os.Create(savedFileName)
	if err != nil {
		return file, savedFileName, err
	}
	defer file.Close()
	file.Write([]byte(header + "\n"))

	for _, generalObject := range queryResults {
		var keyValueData map[string]interface{}
		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)

		var composedString string
		for _, exportField := range exportFields.Data {
			data := keyValueData[exportField.Data]
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")
			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				//childMaster [Id:1, value:childInfo]
				childMaster := childMasterData[float32(numericId)]

				//cateory
				parentLabel := parentChildMap[exportField.Label]

				//Get foreign key from master data 1
				foreignKey := keyValueData[parentLabel]

				//Pick related row from parent master
				rowData := childMaster[util.InterfaceToInt(foreignKey)]
				if rowData == nil {
					composedString = composedString + "" + ","

				} else {
					rowDataString := util.InterfaceToString(rowData[exportField.Data])

					var filteredString = strings.Replace(rowDataString, ",", " ", -1)
					composedString = composedString + filteredString + ","

				}

			} else {
				// Handle 'list_of_attendance' specially by joining values into one column
				if exportField.Data == "list_of_attendance" {
					if data != nil {
						// Assuming data is a slice of strings, join it into one string
						attendanceList := data.([]string)                          // Convert the data into a string slice
						composedString += strings.Join(attendanceList, ", ") + "," // Join with commas
					} else {
						composedString += ","
					}
				} else {
					// Handle other fields as usual
					if data == nil {
						composedString += "" + ","
					} else {
						composedString += util.InterfaceToString(data) + ","
					}
				}
			}

		}
		file.Write([]byte(composedString + "\n"))
	}
	return file, savedFileName, nil
}

func (cm *ComponentManager) writeInCsv(dbConnection *gorm.DB, queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, exportFields component.ExportDataCommand, componentName string) (*os.File, string, error) {
	var header string
	for index, exportField := range exportFields.Data {
		if index == len(exportFields.Data)-1 {
			header = header + exportField.Label
		} else {
			header = header + exportField.Label + ","
		}
	}
	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)
	if err != nil {
		return nil, "cannot create file", err
	}

	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, RenameExportedFile(componentName)+".csv")

	file, err := os.Create(savedFileName)

	if err != nil {
		return file, savedFileName, err
	}
	defer file.Close()
	file.Write([]byte(header + "\n"))
	for _, generalObject := range queryResults {
		var keyValueData map[string]interface{}

		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		var composedString string
		for _, exportField := range exportFields.Data {
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")

			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				//childMaster [Id:1, value:childInfo]
				childMaster := childMasterData[float32(numericId)]

				//cateory
				parentLabel := parentChildMap[fmt.Sprintf("%v", exportField.Id)]
				//Get foreign key from master data 1
				foreignKey := keyValueData[parentLabel]

				//Pick related row from parent master
				rowData := childMaster[util.InterfaceToInt(foreignKey)]
				if parentLabel != "eventId" {
				}

				if exportField.Data == "assemblyLine" && componentName == "labour_management_shift_production" {
					var machineLineResults []component.GeneralObject
					machineMasterInfo := rowData["assemblyLineOption"]
					childQuery := "select * from fuyu_mes.assembly_machine_lines where id = " + util.InterfaceToString(machineMasterInfo)
					dbConnection.Raw(childQuery).Scan(&machineLineResults)
					var machineLineData map[string]interface{}

					json.Unmarshal(machineLineResults[0].ObjectInfo, &machineLineData)
					rowData[exportField.Data] = machineLineData["name"]

				}
				if rowData == nil {
					composedString = composedString + "" + ","
				} else {
					composedString = composedString + util.InterfaceToString(rowData[exportField.Data]) + ","
				}
				fmt.Println("parent label :", parentLabel, " foreignKey :", foreignKey, " row data : ", rowData, " exportField.Data : ", exportField.Data, " value : ", rowData[exportField.Data])
			} else {
				data := keyValueData[exportField.Data]

				if data == nil {
					composedString = composedString + "" + ","
				} else {
					data := keyValueData[exportField.Data]
					if data == nil {
						composedString = composedString + "" + ","
					} else {
						// Check the data type
						_, timeErr := time.Parse(util.ISOTimeLayout, util.InterfaceToString(data))
						if timeErr == nil {
							composedString = composedString + util.ConvertTimeToTimeZonCorrectedFormat("Asia/Singapore", util.InterfaceToString(data)) + ","
						} else {
							composedString = composedString + util.InterfaceToString(data) + ","
						}

					}

				}

			}
		}

		file.Write([]byte(composedString + "\n"))
	}
	return file, savedFileName, nil
}

func (cm *ComponentManager) writeInExcel(dbConnection *gorm.DB, queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, arrayOfExportSchema []component.ExportSchema, exportFields component.ExportDataCommand) (*excelize.File, string, error) {

	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)

	if err != nil {
		return nil, "cannot create file", err
	}

	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, uuid.New().String()+".xlsx")
	file := excelize.NewFile()
	index := file.NewSheet("Sheet1")

	//Adding headers
	for index, exportField := range exportFields.Data {
		headerCoordinateCell, _ := excelize.CoordinatesToCellName(1, index+1)
		file.SetCellValue("Sheet1", headerCoordinateCell, exportField.Label)
	}

	//Adding data
	for rowIndex, generalObject := range queryResults {
		var keyValueData map[string]interface{}

		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		var composedString string
		for colIndex, exportField := range exportFields.Data {
			//Row 1 is already occupied for headers
			coordinateCell, _ := excelize.CoordinatesToCellName(rowIndex+2, colIndex+1)
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")
			originalExportSchema := findExportDataInArray(arrayOfExportSchema, exportField.Id)
			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				//childMaster [Id:1, value:childInfo]
				childMaster := childMasterData[float32(numericId)]

				//cateory
				parentLabel := parentChildMap[exportField.Data]

				//Get foreign key from master data 1
				foreignKey := keyValueData[parentLabel]

				//Pick related row from parent master
				rowData := childMaster[util.InterfaceToInt(foreignKey)]

				if rowData == nil {
					composedString = ""
				} else {
					composedString = util.InterfaceToString(rowData[exportField.Data])
				}
				file.SetCellValue("Sheet1", coordinateCell, composedString)
			} else {
				if originalExportSchema.LinkedMapFlag {
					if originalExportSchema.LinkedObjectMapping.Query != nil {
						query := originalExportSchema.LinkedObjectMapping.Query.Query
						for _, replacementField := range originalExportSchema.LinkedObjectMapping.Query.ReplacementFields {

							if replacementField.Format == component.JsonToStringArray {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							} else {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							}

						}

						var linkedQueryResults []map[string]interface{}
						dbConnection.Raw(query).Scan(&linkedQueryResults)
						if originalExportSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(linkedQueryResults) > 0 {
								fieldValue := linkedQueryResults[0][originalExportSchema.LinkedObjectMapping.Builder.SingleValue.Field]
								composedString = composedString + util.InterfaceToString(fieldValue) + ","
							} else {
								composedString = composedString + "" + ","
							}

						}
					}
				} else {
					data := keyValueData[exportField.Data]
					if data == nil {
						composedString = ""
					} else {
						if exportField.DataType == "datetime" {
							dateTimeData := util.ConvertTimeToTimeZonCorrectedPrimeNgTable("Asia/Singapore", util.InterfaceToString(data))
							composedString = composedString + dateTimeData + ","
						} else {
							composedString = util.InterfaceToString(data)
						}

					}
				}

				file.SetCellValue("Sheet1", coordinateCell, composedString)
			}
		}

	}
	file.SetActiveSheet(index)
	if err := file.SaveAs(savedFileName); err != nil {

		return file, savedFileName, err
	}
	return file, savedFileName, nil
}

// formatCSVValue ensures proper CSV formatting
func formatCSVValue(value string) string {
	if strings.Contains(value, ",") || strings.Contains(value, "\"") {
		// Escape double quotes by replacing " with ""
		escaped := strings.ReplaceAll(value, "\"", "\"\"")
		return fmt.Sprintf("\"%s\"", escaped) // Wrap in double quotes
	}
	return value
}

func writeInCsvV3(cm *ComponentManager, dbConnection *gorm.DB, queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, arrayOfExportSchema []component.ExportSchema, exportFields component.ExportDataCommand, componentName string) (*os.File, string, error) {
	var header string

	for index, exportField := range exportFields.Data {
		if index == len(exportFields.Data)-1 {
			header = header + exportField.Label
		} else {
			header = header + exportField.Label + ","
		}
	}
	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)
	if err != nil {
		return nil, "cannot create file", err
	}

	// savedFileName := cm.ComponentContentConfig.DownloadDirectory + "/" + uuid.New().String() + ".csv"
	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, RenameExportedFile(componentName)+".csv")

	file, err := os.Create(savedFileName)

	if err != nil {
		fmt.Println("error", err)
		return file, savedFileName, err
	}
	defer file.Close()

	file.Write([]byte(header + "\n"))

	for _, generalObject := range queryResults {
		var keyValueData map[string]interface{}
		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		var composedString string

		for _, exportField := range exportFields.Data {
			originalExportSchema := findExportDataInArray(arrayOfExportSchema, exportField.Id)
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")

			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				childMaster := childMasterData[float32(numericId)]
				parentLabel := parentChildMap[exportField.Label]
				foreignKey := keyValueData[parentLabel]

				rowData := childMaster[util.InterfaceToInt(foreignKey)]
				if rowData == nil {
					composedString += "," // If rowData is nil, add an empty value
				} else {
					var formatedValue = formatValue(rowData[exportField.Data])
					formatedValue = formatCSVValue(formatCSVValue(formatedValue))
					composedString += formatedValue + ","
				}

			} else {
				if exportField.LinkedMapFlag {
					if originalExportSchema.LinkedObjectMapping.Query != nil {
						query := originalExportSchema.LinkedObjectMapping.Query.Query
						for _, replacementField := range originalExportSchema.LinkedObjectMapping.Query.ReplacementFields {
							if replacementField.Format == component.JsonToStringArray {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							} else {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							}
						}
						var linkedQueryResults []map[string]interface{}
						dbConnection.Raw(query).Scan(&linkedQueryResults)
						if originalExportSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(linkedQueryResults) > 0 {
								fieldValue := linkedQueryResults[0][originalExportSchema.LinkedObjectMapping.Builder.SingleValue.Field]
								var formatedValue = formatValue(fieldValue)
								formatedValue = formatCSVValue(formatCSVValue(formatedValue))
								composedString += formatedValue + ","
							} else {
								composedString += ","
							}
						} else if originalExportSchema.LinkedObjectMapping.Builder.SingleObjectValueArray != nil {
							if len(linkedQueryResults) > 0 {
								for _, linkedResult := range linkedQueryResults {
									fieldValue := linkedResult[originalExportSchema.LinkedObjectMapping.Builder.SingleObjectValueArray.Field]
									composedString += formatValue(fieldValue) + "; "
								}
								composedString += ","
							} else {
								composedString += ","
							}
						}
					}
				} else {
					data := keyValueData[exportField.Data]
					if data == nil {
						composedString += "," // If data is nil, add an empty value
					} else {
						_, timeErr := time.Parse(util.ISOTimeLayout, util.InterfaceToString(data))
						if timeErr == nil {
							composedString += util.ConvertTimeToTimeZonCorrectedFormat("Asia/Singapore", util.InterfaceToString(data)) + ","
						} else {
							var formatedValue = formatValue(data)
							formatedValue = formatCSVValue(formatCSVValue(formatedValue))
							composedString += formatedValue + ","
						}
					}
				}
			}
		}

		// Write the composed string for this row to the file
		file.Write([]byte(strings.TrimSuffix(composedString, ",") + "\n"))
	}

	return file, savedFileName, nil
}

// Helper function to format values
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%.6f", v) // Customize the precision as needed
	case float32:
		return fmt.Sprintf("%.6f", v) // Customize the precision as needed
	default:
		return fmt.Sprintf("%v", v)
	}
}

func writeInExcelV2(cm *ComponentManager, dbConnection *gorm.DB, queryResults []component.GeneralObject, childMasterData map[float32]map[int]map[string]interface{}, parentChildMap map[string]string, arrayOfExportSchema []component.ExportSchema, exportFields component.ExportDataCommand) (*excelize.File, string, error) {

	err := os.MkdirAll(cm.ComponentContentConfig.DownloadDirectory, os.ModePerm)
	if err != nil {
		return nil, "cannot create file", err
	}

	savedFileName := filepath.Join(cm.ComponentContentConfig.DownloadDirectory, uuid.New().String()+".xlsx")
	file := excelize.NewFile()
	index := file.NewSheet("Sheet1")

	//Adding headers
	for index, exportField := range exportFields.Data {
		headerCoordinateCell, _ := excelize.CoordinatesToCellName(1, index+1)
		file.SetCellValue("Sheet1", headerCoordinateCell, exportField.Label)
	}

	//Adding data
	for rowIndex, generalObject := range queryResults {
		var keyValueData map[string]interface{}

		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		var composedString string
		for colIndex, exportField := range exportFields.Data {
			//Row 1 is already occupied for headers
			coordinateCell, _ := excelize.CoordinatesToCellName(rowIndex+2, colIndex+1)
			findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")
			originalExportSchema := findExportDataInArray(arrayOfExportSchema, exportField.Id)
			if len(findDataType) > 1 {
				numericId, err := strconv.Atoi(findDataType[0])

				if err != nil {
					continue
				}
				//childMaster [Id:1, value:childInfo]
				childMaster := childMasterData[float32(numericId)]

				//cateory
				parentLabel := parentChildMap[exportField.Data]

				//Get foreign key from master data 1
				foreignKey := keyValueData[parentLabel]

				//Pick related row from parent master
				rowData := childMaster[util.InterfaceToInt(foreignKey)]

				if rowData == nil {
					composedString = ""
				} else {
					composedString = util.InterfaceToString(rowData[exportField.Data])
				}
				file.SetCellValue("Sheet1", coordinateCell, composedString)
			} else {
				if originalExportSchema.LinkedMapFlag {
					if originalExportSchema.LinkedObjectMapping.Query != nil {
						query := originalExportSchema.LinkedObjectMapping.Query.Query
						for _, replacementField := range originalExportSchema.LinkedObjectMapping.Query.ReplacementFields {

							if replacementField.Format == component.JsonToStringArray {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							} else {
								value := keyValueData[replacementField.Property]
								replacementValue := util.InterfaceToString(value)
								query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, -1)
							}
						}
						var linkedQueryResults []map[string]interface{}
						dbConnection.Raw(query).Scan(&linkedQueryResults)
						if originalExportSchema.LinkedObjectMapping.Builder.SingleValue != nil {
							if len(linkedQueryResults) > 0 {
								fieldValue := linkedQueryResults[0][originalExportSchema.LinkedObjectMapping.Builder.SingleValue.Field]
								composedString = composedString + util.InterfaceToString(fieldValue) + ","
							} else {
								composedString = composedString + "" + ","
							}
						}
					}
				} else {
					data := keyValueData[exportField.Data]
					if data == nil {
						composedString = ""
					} else {
						if exportField.DataType == "datetime" {
							dateTimeData := util.ConvertTimeToTimeZonCorrectedPrimeNgTable("Asia/Singapore", util.InterfaceToString(data))
							composedString = composedString + dateTimeData + ","
						} else {
							composedString = util.InterfaceToString(data)
						}

					}
				}

				file.SetCellValue("Sheet1", coordinateCell, composedString)
			}
		}

	}
	file.SetActiveSheet(index)
	if err := file.SaveAs(savedFileName); err != nil {

		return file, savedFileName, err
	}
	return file, savedFileName, nil
}

func getChildrenData(exportFieldData []component.ExportSchema, componentName string, cm *ComponentManager, dbConnection *gorm.DB) (map[float32]map[int]map[string]interface{}, map[string]string) {
	childData := make(map[float32]map[int]map[string]interface{})
	parentChildMap := make(map[string]string)
	//var fetchedChild []string
	targetTable := cm.GetTableExportSchema(componentName)
	var queryResults []component.GeneralObject
	for _, exportField := range exportFieldData {
		findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")
		if len(findDataType) > 0 {
			numericId, err := strconv.Atoi(findDataType[0])
			if err != nil {
				continue
			}
			schema := findBySchemaId(float32(numericId), targetTable)

			if schema.TargetTable == "" {
				continue
			}
			// 4.1 ->
			parentChildMap[fmt.Sprintf("%v", exportField.Id)] = schema.Data
			fmt.Println("parent child map, key [", fmt.Sprintf("%v", exportField.Id), " schema.data [", schema.Data, "]")
			childQuery := "select * from " + schema.TargetTable
			dbConnection.Raw(childQuery).Scan(&queryResults)
			childData[schema.Id] = formateChildMasterData(queryResults)
		}
	}

	return childData, parentChildMap
}

func getChildrenDataV2(exportFieldData []component.ExportSchema, componentName string, cm *ComponentManager, dbConnection *gorm.DB) (map[float32]map[int]map[string]interface{}, map[string]string) {
	childData := make(map[float32]map[int]map[string]interface{})
	parentChildMap := make(map[string]string)
	//var fetchedChild []string
	targetTable := cm.GetTableExportSchema(componentName)
	var queryResults []component.GeneralObject
	for _, exportField := range exportFieldData {
		findDataType := strings.Split(fmt.Sprintf("%v", exportField.Id), ".")

		if len(findDataType) > 0 {
			numericId, err := strconv.Atoi(findDataType[0])

			if err != nil {
				continue
			}
			schema := findBySchemaId(float32(numericId), targetTable)
			if schema.TargetTable == "" {
				continue
			}
			parentChildMap[exportField.Label] = schema.Data

			childQuery := "select * from " + schema.TargetTable
			fmt.Println("[export] childQuery under getChildrenDataV2 :", childQuery)
			dbConnection.Raw(childQuery).Scan(&queryResults)
			childData[schema.Id] = formateChildMasterData(queryResults)
		}
	}

	return childData, parentChildMap
}

func formateChildMasterData(queryResult []component.GeneralObject) map[int]map[string]interface{} {
	childData := make(map[int]map[string]interface{})
	//var keyValueData map[string]interface{}
	for _, generalObject := range queryResult {
		keyValueData := make(map[string]interface{})
		json.Unmarshal(generalObject.ObjectInfo, &keyValueData)
		childData[generalObject.Id] = keyValueData
	}
	return childData
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func findBySchemaId(id float32, exportSchema []component.ExportSchema) component.ExportSchema {
	for _, schema := range exportSchema {
		if schema.Id == id {
			return schema
		}
	}
	return component.ExportSchema{}
}

func (cm *ComponentManager) ProcessLoadFile(loadFileCommand component.LoadDataFileCommand) (error, int, component.LoadFileResponse) {
	savedFileName := cm.ComponentContentConfig.DownloadDirectory + "/" + uuid.New().String()
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: savedFileName,
		Dir: false,
		//the repository with a subdirectory I would like to clone only
		Src:  loadFileCommand.ContentUrl,
		Mode: getter.ClientModeFile,
		////define the type of detectors go getter should use, in this case only github is needed

		////provide the getter needed to download the files
		Getters:  cm.ComponentContentConfig.GetGetter(),
		Insecure: cm.ComponentContentConfig.Insecure,
	}
	//download the files
	if err := client.Get(); err != nil {
		return errors.New("failed to process given file, check the file again"), FailedToDownloadTheImportFileUrl, component.LoadFileResponse{}
	}

	file, err := os.Open(savedFileName)
	if err != nil {
		return errors.New("unable to read csv file"), UnableToReadCSVFile, component.LoadFileResponse{}
	}
	defer file.Close()

	utfData, _ := utfbom.Skip(file)
	csvReader := csv.NewReader(utfData)
	records, err := csvReader.ReadAll()
	if err != nil {
		return errors.New("unable to parse csv file"), ParsingCSVFileFailed, component.LoadFileResponse{}
	}
	// validate all the columns are available
	var schemaMapping = make(map[int]component.CRUDTableSchema, len(records[0]))

	for i, record := range records {
		fmt.Println(record)
		if i == 0 {
			// column definitions
			for position, individualColumn := range record {

				crudTableSchema := component.CRUDTableSchema{
					Property: util.CamelCase(individualColumn),
					Name:     individualColumn,
					Display:  true,
					Type:     "text",
				}
				schemaMapping[position] = crudTableSchema
			}
		} else {
			break
		}
	}
	// now compose the json record and insert the record to database
	var totalRecords int64
	var modifiedObjectList []datatypes.JSON
	for i, record := range records {
		if i > 0 {
			// now focus on the records
			var insertRecord = make(map[string]interface{})

			for j, value := range record {
				if crudTableSchema, ok := schemaMapping[j]; ok {
					insertRecord[crudTableSchema.Property] = value
				}
			}

			rawRecordInfo, _ := json.Marshal(insertRecord)
			modifiedObjectList = append(modifiedObjectList, rawRecordInfo)
			totalRecords = totalRecords + 1
		}
	}
	loadFileResponse := component.LoadFileResponse{}
	var arrayOfTableSchema []component.CRUDTableSchema
	for _, tableSchema := range schemaMapping {
		arrayOfTableSchema = append(arrayOfTableSchema, tableSchema)
	}
	loadFileResponse.TotalRowCount = totalRecords
	loadFileResponse.Data = modifiedObjectList
	loadFileResponse.TableSchema = arrayOfTableSchema
	return nil, 0, loadFileResponse
}

func getUpdateActionId(dbConnection *gorm.DB, updateMapping *component.ObjectMapping, recordData map[string]interface{}) int {
	queryTemplate := updateMapping.Builder.SingleValueCondition.Action.Query
	// if it is update, we need to check the replacement fields if any
	for _, replacementField := range updateMapping.Builder.SingleValueCondition.Action.ReplacementFields {
		value := recordData[replacementField.Property]
		if replacementField.Format == component.JsonToStringArray {

			replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
			queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
		} else {
			replacementValue := util.InterfaceToString(value)
			queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
		}

	}
	var innerQueryResults []map[string]interface{}
	dbConnection.Raw(queryTemplate).Scan(&innerQueryResults)
	idColumn := updateMapping.Builder.SingleValueCondition.Action.IdColumn
	objectIdValue := util.InterfaceToInt(innerQueryResults[0][idColumn])
	return objectIdValue
}
func (cm *ComponentManager) ImportData(dbConnection *gorm.DB, componentName string, importDataCommand component.ImportDataCommand) (error, int, component.ImportDataObjects) {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	var insertObject []component.GeneralObject
	var updateObject []component.GeneralObject
	skippedRecordInfo := component.RecordInfo{}
	var totalRecords int
	var totalSkippedRecords int
	importDataObjects := component.ImportDataObjects{
		TotalRecords:        0,
		TotalSkippedRecords: 0,
		SkippedData:         nil,
		InsertObjects:       insertObject,
		UpdateObjects:       updateObject,
	}
	err, savedFileName := cm.DownloadFile(importDataCommand.ContentUrl)
	if err != nil {
		return errors.New("failed to process given file, check the file again"), FailedToDownloadTheImportFileUrl, importDataObjects
	}
	componentSchema := cm.ComponentSchema[int64ComponentId]
	tableImportSchema := componentSchema.TableImportSchema
	extraFields := componentSchema.TableImportSchema.ExtraField
	recordSchema := componentSchema.RecordSchema
	tableImportFields := tableImportSchema.Fields
	updateMapping := tableImportSchema.Update
	var destinationFieldMap = make(map[string]component.TableImportFields, len(tableImportSchema.Fields))
	for _, tc := range tableImportFields {
		destinationFieldMap[tc.DestinationField] = tc
	}

	file, err := os.Open(savedFileName)
	if err != nil {
		return errors.New("unable to read csv file"), UnableToReadCSVFile, importDataObjects
	}

	//csvReader := csv.NewReader(f)
	utfData, _ := utfbom.Skip(file)
	csvReader := csv.NewReader(utfData)
	records, err := csvReader.ReadAll()

	if err != nil {
		return errors.New("unable to parse csv file"), ParsingCSVFileFailed, importDataObjects
	}
	//defer f.Close()
	// validate all the columns are available
	var schemaMapping = make(map[int]component.ImportSchemaRequest, len(records[0]))

	for i, record := range records {
		fmt.Println(record)
		if i == 0 {
			// column definitions
			for _, schemaElement := range importDataCommand.Schema {
				isExist, position := util.StringContainsWithPos(record, strings.TrimSpace(schemaElement.SourceField))
				if isExist {
					schemaMapping[position] = schemaElement
				} else {
					return errors.New("schema is not matched with uploaded CSV"), SchemaIsNotMatchedWithUploadedCSV, importDataObjects
				}
			}
		} else {
			break
		}
	}
	var defaultValueFields = make(map[string]string, 10)
	for _, recordSchemaField := range recordSchema {
		if recordSchemaField.Default != nil {
			defaultValueFields[recordSchemaField.Property] = util.InterfaceToString(*recordSchemaField.Default)
		}
	}

	// now compose the json record and insert the record to database

	//lets create the array of map interface
	var extractedRecords []map[string]interface{}
	for i, record := range records {
		if i > 0 {
			var recordData = make(map[string]interface{})
			for j, value := range record {
				if val, ok := schemaMapping[j]; ok {
					recordData[val.DestinationField] = value

				}
			}
			extractedRecords = append(extractedRecords, recordData)
		}
	}
	var skippedData []map[string]interface{}
	var headerData []component.TableSchema
	var insertRecords []map[string]interface{}
	var updateRecords []map[string]interface{}
	var updateRecordIds []int
	var skipRecordsNames string
	skipRecordsNames = " [ "
	for index, recordData := range extractedRecords {
		var isSkipEntireRecord bool
		isSkipEntireRecord = false
		if index == 0 {
			// build the header
			for key, _ := range recordData {
				fieldSchema := destinationFieldMap[key]
				headerData = append(headerData, component.TableSchema{
					Name:     fieldSchema.DisplayName,
					Type:     fieldSchema.DataType,
					Display:  true,
					Property: fieldSchema.DestinationField,
				})
			}
		}

		for key, value := range recordData {
			var isValidationPass bool
			isValidationPass = true
			fieldSchema := destinationFieldMap[key]
			if fieldSchema.Validation != nil {
				// we have validation defined.
				var queryResults []map[string]interface{}
				queryObject := fieldSchema.Validation.Query
				queryTemplate := queryObject.Query
				for _, replacementField := range queryObject.ReplacementFields {
					if replacementField.Format == component.JsonToStringArray {
						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						replacementValue := util.InterfaceToString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}
				dbConnection.Raw(queryTemplate).Scan(&queryResults)
				if len(queryResults) == 0 {
					isValidationPass = false
				}
				if fieldSchema.Validation.Builder.SingleValueCondition != nil {
					checkField := fieldSchema.Validation.Builder.SingleValueCondition.Field
					checkFieldValue := util.InterfaceToInt(queryResults[0][checkField])
					condition := fieldSchema.Validation.Builder.SingleValueCondition.Condition
					if condition == "=" {
						if checkFieldValue == fieldSchema.Validation.Builder.SingleValueCondition.Value {
							if fieldSchema.Validation.Builder.SingleValueCondition.Action.Type == "continue" {
								isValidationPass = true
							}
						}
					} else if condition == ">" {
						if checkFieldValue > fieldSchema.Validation.Builder.SingleValueCondition.Value {
							if fieldSchema.Validation.Builder.SingleValueCondition.Action.Type == "continue" {
								isValidationPass = true
							}
						}
					} else if condition == "<" {
						if checkFieldValue < fieldSchema.Validation.Builder.SingleValueCondition.Value {
							if fieldSchema.Validation.Builder.SingleValueCondition.Action.Type == "continue" {
								isValidationPass = true
							}
						}
					}
				}
			}
			// validation pass only , we will proceed to next step
			if isValidationPass {
				if fieldSchema.Replacement != nil {
					var queryResults []map[string]interface{}
					queryObject := fieldSchema.Replacement.Query
					queryTemplate := queryObject.Query
					for _, replacementField := range queryObject.ReplacementFields {
						if replacementField.Format == component.JsonToStringArray {
							replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
							queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
						} else {
							replacementValue := util.InterfaceToString(value)
							queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
						}

					}
					dbConnection.Raw(queryTemplate).Scan(&queryResults)
					if fieldSchema.Replacement.Builder.FieldPropertyAssignment != nil {

						if len(queryResults) > 0 {
							field := fieldSchema.Replacement.Builder.FieldPropertyAssignment.Field
							propertyField := fieldSchema.Replacement.Builder.FieldPropertyAssignment.Property
							fieldValue := util.InterfaceToString(queryResults[0][field])
							recordData[propertyField] = fieldValue
						} else {
							isSkipEntireRecord = true
							skipRecordsNames += fieldSchema.DisplayName + ","

						}

					}
				}
			}
			if isSkipEntireRecord {
				// we found that the entire record should be skipped. so now create the record info and put into array
				skippedData = append(skippedData, recordData)
				break
			}
			recordDataValue := recordData[key]
			if fieldSchema.DataType == "int" {
				recordData[key] = util.InterfaceToInt(recordDataValue)
			} else if fieldSchema.DataType == "double" {
				recordData[key] = util.InterfaceToFloat(recordDataValue)
			} else {
				recordData[key] = recordDataValue
			}

		}
		if !isSkipEntireRecord {
			// this is not skipped, so we need to check whether update condition is met, otherwise, it should goes
			// insert , else update
			if updateMapping != nil {
				queryTemplate := updateMapping.Query.Query
				for _, replacementField := range updateMapping.Query.ReplacementFields {
					value := recordData[replacementField.Property]
					if replacementField.Format == component.JsonToStringArray {

						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						replacementValue := util.InterfaceToString(value)
						queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}
				var queryResults []map[string]interface{}
				dbConnection.Raw(queryTemplate).Scan(&queryResults)

				checkField := updateMapping.Builder.SingleValueCondition.Field
				checkFieldValue := util.InterfaceToInt(queryResults[0][checkField])
				condition := updateMapping.Builder.SingleValueCondition.Condition

				if condition == "=" {
					if checkFieldValue == updateMapping.Builder.SingleValueCondition.Value {
						if updateMapping.Builder.SingleValueCondition.Action.Type == "update" {
							queryTemplate = updateMapping.Builder.SingleValueCondition.Action.Query
							// if it is update, we need to check the replacement fields if any
							updateId := getUpdateActionId(dbConnection, updateMapping, recordData)
							updateRecords = append(updateRecords, recordData)
							updateRecordIds = append(updateRecordIds, updateId)
						}
					} else {
						insertRecords = append(insertRecords, recordData)
					}
				} else if condition == ">" {
					if checkFieldValue > updateMapping.Builder.SingleValueCondition.Value {
						if updateMapping.Builder.SingleValueCondition.Action.Type == "update" {
							updateId := getUpdateActionId(dbConnection, updateMapping, recordData)
							updateRecords = append(updateRecords, recordData)
							updateRecordIds = append(updateRecordIds, updateId)
							updateRecords = append(updateRecords, recordData)
							updateRecordIds = append(updateRecordIds, updateId)
						}
					} else {
						insertRecords = append(insertRecords, recordData)
					}
				} else if condition == "<" {
					if checkFieldValue < updateMapping.Builder.SingleValueCondition.Value {
						if updateMapping.Builder.SingleValueCondition.Action.Type == "update" {
							updateId := getUpdateActionId(dbConnection, updateMapping, recordData)
							updateRecords = append(updateRecords, recordData)
							updateRecordIds = append(updateRecordIds, updateId)
							updateRecords = append(updateRecords, recordData)
							updateRecordIds = append(updateRecordIds, updateId)
						}
					} else {
						insertRecords = append(insertRecords, recordData)
					}
				} else {
					insertRecords = append(insertRecords, recordData)
				}

			} else {
				insertRecords = append(insertRecords, recordData)
			}
		}

		// lets work on extra fields
		for _, extraField := range extraFields {
			if extraField.ObjectMapping != nil {
				recordData[extraField.Field] = recordData[extraField.ObjectMapping.Builder.ObjectFieldAssignment.Property]
			} else {
				recordData[extraField.Field] = extraField.Default

			}
		}

	}

	skippedRecordInfo.Data = skippedData
	skippedRecordInfo.Header = headerData

	totalRecords = len(extractedRecords)
	totalSkippedRecords = len(skippedData)
	for _, insertRecordData := range insertRecords {
		rawObject, _ := json.Marshal(insertRecordData)
		generalObject := component.GeneralObject{ObjectInfo: rawObject}
		insertObject = append(insertObject, generalObject)
	}
	for index, updateRecordData := range updateRecords {
		rawObject, _ := json.Marshal(updateRecordData)
		generalObject := component.GeneralObject{Id: updateRecordIds[index], ObjectInfo: rawObject}
		updateObject = append(updateObject, generalObject)
	}

	skipRecordsNames += "]"
	importDataObjects.TotalSkippedRecords = totalSkippedRecords
	importDataObjects.TotalRecords = totalRecords
	importDataObjects.SkippedData = skippedRecordInfo
	importDataObjects.InsertObjects = insertObject
	importDataObjects.UpdateObjects = updateObject
	importDataObjects.SkippedRecordNames = skipRecordsNames
	return nil, 0, importDataObjects
}

func (cm *ComponentManager) GetAbsoluteSearchQuery(componentName string, searchFields []component.SearchKeys) string {

	//tableSchema := cm.GetTableSchema(componentName)
	basicCondition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$."
	var searchQuery string

	for index, searchField := range searchFields {
		if len(searchFields)-1 == index {
			searchQuery = searchQuery + basicCondition + searchField.Field + "\")) =" + searchField.Value + ""

		} else {
			searchQuery = searchQuery + basicCondition + searchField.Field + "\")) = " + searchField.Value + "' or "

		}
	}
	return searchQuery
}

func (cm *ComponentManager) GetSearchQuery(componentName string, searchFields []component.SearchKeys) string {

	basicCondition := " object_info ->>'$."
	var searchQuery string

	for index, searchField := range searchFields {
		if len(searchFields)-1 == index {
			searchQuery = searchQuery + basicCondition + searchField.Field + "' like '%" + searchField.Value + "%' "

		} else {
			searchQuery = searchQuery + basicCondition + searchField.Field + "' like '%" + searchField.Value + "%' and "

		}
	}
	return searchQuery
}

func (cm *ComponentManager) GetWhereLikeQuery(componentName string, searchFields []component.SearchKeys) string {

	basicCondition := " " + componentName + ".object_info ->>'$."
	var searchQuery string

	for index, searchField := range searchFields {
		if len(searchFields)-1 == index {
			searchQuery = searchQuery + basicCondition + searchField.Field + "' like '%" + searchField.Value + "%' "

		} else {
			searchQuery = searchQuery + basicCondition + searchField.Field + "' like '%" + searchField.Value + "%' and "

		}
	}
	return searchQuery
}

func (cm *ComponentManager) GetSearchQueryV2(dbConnection *gorm.DB, componentName string, searchFields []component.SearchKeys) string {
	selectQuery := componentName + ".id "
	var joinQuery string

	componentWhereClause := " where (" + cm.GetWhereLikeQuery(componentName, searchFields) + ") or ("

	arrayOfLinkedObjectMapSchema := cm.GetLinkedObjectMapSchema(componentName)
	fmt.Println(componentName)
	fmt.Println(arrayOfLinkedObjectMapSchema)

	for _, linkedObjectMap := range arrayOfLinkedObjectMapSchema {
		joinQuery = joinQuery + " join " + linkedObjectMap.LinkedTable + " on " + componentName + ".object_info ->> '$." + linkedObjectMap.Property + "' = " + linkedObjectMap.LinkedTable + ".id "
		selectQuery = selectQuery + ", " + linkedObjectMap.LinkedTable + ".object_info ->> '$." + linkedObjectMap.TargetProperty + "' "
		for index, searchField := range searchFields {
			if len(searchFields)-1 == index {
				componentWhereClause = " " + componentWhereClause + linkedObjectMap.LinkedTable + ".object_info ->> '$." + linkedObjectMap.TargetProperty + "' like '%" + searchField.Value + "%'"

			} else {
				componentWhereClause = " " + componentWhereClause + linkedObjectMap.LinkedTable + ".object_info ->> '$." + linkedObjectMap.TargetProperty + "' like '%" + searchField.Value + "%' and "

			}

		}
		componentWhereClause = componentWhereClause + ") or ("
	}
	searchQuery := "select " + selectQuery + " from " + componentName + joinQuery + util.TrimSuffix(componentWhereClause, "or (")

	var selectedIds []map[string]interface{}
	dbConnection.Raw(searchQuery).Scan(&selectedIds)

	idList := "("
	for _, idObjects := range selectedIds {
		id := util.InterfaceToString(idObjects["id"])
		idList += id + ","
	}

	idList = util.TrimSuffix(idList, ",") + ")"
	fmt.Println("Length:", idList)
	return idList
}

func (cm *ComponentManager) GetSearchResponse(componentName string, listOfObjects *[]component.GeneralObject) component.TableObjectResponse {

	arrayOfTableSchema := cm.GetTableSchema(componentName)
	var modifiedObjectList []datatypes.JSON
	var decodedBody = make(map[string]interface{})
	for _, generalObjects := range *listOfObjects {
		json.Unmarshal(generalObjects.ObjectInfo, &decodedBody)
		decodedBody["id"] = generalObjects.Id
		rawData, _ := json.Marshal(decodedBody)
		modifiedObjectList = append(modifiedObjectList, rawData)
	}

	tableObjectResponse := component.TableObjectResponse{}
	tableObjectResponse.TotalRowCount = int64(len(*listOfObjects))
	tableObjectResponse.CurrentRowCount = int64(len(*listOfObjects))
	if len(*listOfObjects) == 0 {
		modifiedObjectList = make([]datatypes.JSON, 0)
		tableObjectResponse.Data = modifiedObjectList
	} else {
		tableObjectResponse.Data = modifiedObjectList
	}
	tableObjectResponse.Header = arrayOfTableSchema

	return tableObjectResponse
}

func GetRecordTrailResponse(zone string, listOfRecords *[]component.GeneralObject) interface{} {
	var recordMessageMap = make(map[string][]RecordMessageData)
	for _, recordMessageInterface := range *listOfRecords {
		resourceInfo := ResourceInfo{}
		json.Unmarshal(recordMessageInterface.ObjectInfo, &resourceInfo)
		timeObject := util.ConvertStringToDateTime(resourceInfo.ResourceMeta.CreatedAt)

		dateString := timeObject.DateTime.Month().String() + " " + strconv.Itoa(timeObject.DateTime.Day()) + " " + strconv.Itoa(timeObject.DateTime.Year())
		recordMessageData := RecordMessageData{}
		authService := GetService("general_auth").ServiceInterface.(AuthInterface)
		basicUserInfo := authService.GetUserInfoById(resourceInfo.ResourceMeta.UserId)
		// first split the time
		referencedTime := time.Now().Unix() - timeObject.DateTimeEpoch
		recordMessageData.AvatarUrl = basicUserInfo.AvatarUrl
		recordMessageData.UserId = basicUserInfo.UserId
		recordMessageData.Username = basicUserInfo.FullName
		recordMessageData.UserProfileRouterLink = "user_reg_edit"
		recordMessageData.ReferenceTime = util.ConvertReferenceTimeToString(referencedTime)
		recordMessageData.Message = resourceInfo.Message
		recordMessageData.CreatedAt = util.ISO2TableDateTimeFormat(zone, resourceInfo.ResourceMeta.CreatedAt)
		if len(resourceInfo.TrackingFields) == 0 {
			recordMessageData.TrackingFields = []TrackingFieldsResponse{}
		} else {
			var trackingFieldsResponse []TrackingFieldsResponse
			for _, trackingField := range resourceInfo.TrackingFields {
				if trackingField.ChangedField == "lastUpdatedBy" {
					basicUserInfoOld := authService.GetUserInfoById(util.InterfaceToInt(trackingField.OldValue))
					basicUserInfoNew := authService.GetUserInfoById(util.InterfaceToInt(trackingField.NewValue))
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     basicUserInfoOld.FullName,
						NewValue:     basicUserInfoNew.FullName,
					})
				} else if trackingField.ChangedField == "lastUpdatedAt" {
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
						NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
					})
				} else {
					switch trackingField.NewValue.(type) {
					case string:
						if util.IsDateString(util.InterfaceToString(trackingField.OldValue)) {
							trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
								ChangedField: trackingField.ChangedField,
								OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
								NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
							})
						} else {
							trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
								ChangedField: trackingField.ChangedField,
								OldValue:     util.InterfaceToString(trackingField.OldValue),
								NewValue:     util.InterfaceToString(trackingField.NewValue),
							})
						}

					case map[string]interface{}:
						var objectFields = trackingField.NewValue.(map[string]interface{})
						var modifiedValue string
						modifiedValue = "["
						var index = 0
						for key, value := range objectFields {
							index = index + 1
							if len(objectFields) == index {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "]"
							} else {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "] ["
							}

						}
						trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
							ChangedField: trackingField.ChangedField,
							OldValue:     util.InterfaceToString(trackingField.OldValue),
							NewValue:     modifiedValue,
						})
					}
				}
			}

			recordMessageData.TrackingFields = trackingFieldsResponse
		}

		if arrayOfTimelines, ok := recordMessageMap[dateString]; ok {
			arrayOfTimelines = append(arrayOfTimelines, recordMessageData)
			recordMessageMap[dateString] = arrayOfTimelines
		} else {
			var timelineArray []RecordMessageData
			timelineArray = append(timelineArray, recordMessageData)
			recordMessageMap[dateString] = timelineArray
		}
		//April 26 2022

	}

	if len(recordMessageMap) == 0 {
		// we need to return the empty response
		var emptyArray = make([]string, 0)
		return emptyArray
	} else {
		var response []RecordMessageResponse
		for val, r := range recordMessageMap {
			recordMessageResponse := RecordMessageResponse{
				Date: val,
				Data: r,
			}
			response = append(response, recordMessageResponse)
		}
		return response
	}
}

func GetRecordTrailResponseV1(zone string, listOfRecords []component.GeneralObject) interface{} {
	var recordMessageMap = make(map[string][]RecordMessageData)
	for _, recordMessageInterface := range listOfRecords {
		resourceInfo := ResourceInfo{}
		json.Unmarshal(recordMessageInterface.ObjectInfo, &resourceInfo)
		timeObject := util.ConvertStringToDateTime(resourceInfo.ResourceMeta.CreatedAt)

		dateString := timeObject.DateTime.Month().String() + " " + strconv.Itoa(timeObject.DateTime.Day()) + " " + strconv.Itoa(timeObject.DateTime.Year())
		recordMessageData := RecordMessageData{}
		authService := GetService("general_auth").ServiceInterface.(AuthInterface)
		basicUserInfo := authService.GetUserInfoById(resourceInfo.ResourceMeta.UserId)
		// first split the time
		referencedTime := time.Now().Unix() - timeObject.DateTimeEpoch
		recordMessageData.AvatarUrl = basicUserInfo.AvatarUrl
		recordMessageData.UserId = basicUserInfo.UserId
		recordMessageData.Username = basicUserInfo.FullName
		recordMessageData.UserProfileRouterLink = "user_reg_edit"
		recordMessageData.ReferenceTime = util.ConvertReferenceTimeToString(referencedTime)
		recordMessageData.Message = resourceInfo.Message
		recordMessageData.CreatedAt = util.ISO2TableDateTimeFormat(zone, resourceInfo.ResourceMeta.CreatedAt)
		if len(resourceInfo.TrackingFields) == 0 {
			recordMessageData.TrackingFields = []TrackingFieldsResponse{}
		} else {
			var trackingFieldsResponse []TrackingFieldsResponse
			for _, trackingField := range resourceInfo.TrackingFields {
				if trackingField.ChangedField == "lastUpdatedBy" {
					basicUserInfoOld := authService.GetUserInfoById(util.InterfaceToInt(trackingField.OldValue))
					basicUserInfoNew := authService.GetUserInfoById(util.InterfaceToInt(trackingField.NewValue))
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     basicUserInfoOld.FullName,
						NewValue:     basicUserInfoNew.FullName,
					})
				} else if trackingField.ChangedField == "lastUpdatedAt" {
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
						NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
					})
				} else {
					switch trackingField.NewValue.(type) {
					case string:
						if util.IsDateString(util.InterfaceToString(trackingField.OldValue)) {
							trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
								ChangedField: trackingField.ChangedField,
								OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
								NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
							})
						} else {
							trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
								ChangedField: trackingField.ChangedField,
								OldValue:     util.InterfaceToString(trackingField.OldValue),
								NewValue:     util.InterfaceToString(trackingField.NewValue),
							})
						}

					case map[string]interface{}:
						var objectFields = trackingField.NewValue.(map[string]interface{})
						var modifiedValue string
						modifiedValue = "["
						var index = 0
						for key, value := range objectFields {
							index = index + 1
							if len(objectFields) == index {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "]"
							} else {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "] ["
							}

						}
						trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
							ChangedField: trackingField.ChangedField,
							OldValue:     util.InterfaceToString(trackingField.OldValue),
							NewValue:     modifiedValue,
						})
					}
				}
			}

			recordMessageData.TrackingFields = trackingFieldsResponse
		}

		if arrayOfTimelines, ok := recordMessageMap[dateString]; ok {
			arrayOfTimelines = append(arrayOfTimelines, recordMessageData)
			recordMessageMap[dateString] = arrayOfTimelines
		} else {
			var timelineArray []RecordMessageData
			timelineArray = append(timelineArray, recordMessageData)
			recordMessageMap[dateString] = timelineArray
		}
		//April 26 2022

	}

	if len(recordMessageMap) == 0 {
		// we need to return the empty response
		var emptyArray = make([]string, 0)
		return emptyArray
	} else {
		var response []RecordMessageResponse
		for val, r := range recordMessageMap {
			recordMessageResponse := RecordMessageResponse{
				Date: val,
				Data: r,
			}
			response = append(response, recordMessageResponse)
		}
		return response
	}
}

func GetRecordTrailResponse_V1(zone string, dbConnection *gorm.DB, listOfRecords *[]component.GeneralObject, recordSchema []component.RecordSchema) interface{} {
	var recordMessageMap = make(map[string][]RecordMessageData)
	for _, recordMessageInterface := range *listOfRecords {
		resourceInfo := ResourceInfo{}
		json.Unmarshal(recordMessageInterface.ObjectInfo, &resourceInfo)
		timeObject := util.ConvertStringToDateTime(resourceInfo.ResourceMeta.CreatedAt)

		dateString := timeObject.DateTime.Month().String() + " " + strconv.Itoa(timeObject.DateTime.Day()) + " " + strconv.Itoa(timeObject.DateTime.Year())
		recordMessageData := RecordMessageData{}
		authService := GetService("general_auth").ServiceInterface.(AuthInterface)
		basicUserInfo := authService.GetUserInfoById(resourceInfo.ResourceMeta.UserId)
		// first split the time
		referencedTime := time.Now().Unix() - timeObject.DateTimeEpoch
		recordMessageData.AvatarUrl = basicUserInfo.AvatarUrl
		recordMessageData.UserId = basicUserInfo.UserId
		recordMessageData.Username = basicUserInfo.FullName
		recordMessageData.UserProfileRouterLink = "user_reg_edit"
		recordMessageData.ReferenceTime = util.ConvertReferenceTimeToString(referencedTime)
		recordMessageData.Message = resourceInfo.Message
		recordMessageData.CreatedAt = util.ISO2TableDateTimeFormat(zone, resourceInfo.ResourceMeta.CreatedAt)
		if len(resourceInfo.TrackingFields) == 0 {
			recordMessageData.TrackingFields = []TrackingFieldsResponse{}
		} else {
			var trackingFieldsResponse []TrackingFieldsResponse
			for _, trackingField := range resourceInfo.TrackingFields {
				if trackingField.ChangedField == "lastUpdatedBy" {
					basicUserInfoOld := authService.GetUserInfoById(util.InterfaceToInt(trackingField.OldValue))
					basicUserInfoNew := authService.GetUserInfoById(util.InterfaceToInt(trackingField.NewValue))
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     basicUserInfoOld.FullName,
						NewValue:     basicUserInfoNew.FullName,
					})
				} else if trackingField.ChangedField == "lastUpdatedAt" {
					trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
						ChangedField: trackingField.ChangedField,
						OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
						NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
					})
				} else {
					switch trackingField.NewValue.(type) {
					case string:
						if util.IsDateString(util.InterfaceToString(trackingField.OldValue)) {
							trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
								ChangedField: trackingField.ChangedField,
								OldValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.OldValue)),
								NewValue:     util.ISO2TableDateTimeFormat(zone, util.InterfaceToString(trackingField.NewValue)),
							})
						} else {

							// check the field has anything linkedObjectMapping, so that we can get the linked data
							isLinkedField := false
							var linkedObjectMapping *component.ObjectMapping
							for _, schema := range recordSchema {
								if schema.Property == trackingField.ChangedField {
									if schema.LinkedObjectMapping != nil {
										isLinkedField = true
										linkedObjectMapping = schema.LinkedObjectMapping
									}
								}
							}
							if !isLinkedField {
								trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
									ChangedField: trackingField.ChangedField,
									OldValue:     util.InterfaceToString(trackingField.OldValue),
									NewValue:     util.InterfaceToString(trackingField.NewValue),
								})
							} else {
								var oldValue string = "-"
								var newValue string = "-"
								if linkedObjectMapping.Query != nil {
									// if the query object is configured
									queryTemplate := linkedObjectMapping.Query.Query
									var queryResults []map[string]interface{}

									for _, replacementField := range linkedObjectMapping.Query.ReplacementFields {
										if replacementField.Format == component.JsonToStringArray {
											value := util.InterfaceToString(trackingField.OldValue)
											replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
											queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
										} else {
											value := util.InterfaceToString(trackingField.OldValue)
											replacementValue := util.InterfaceToString(value)
											queryTemplate = strings.Replace(queryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
										}

									}
									dbConnection.Raw(queryTemplate).Scan(&queryResults)
									if linkedObjectMapping.Builder.SingleValue != nil {
										if len(queryResults) > 0 {
											oldValue = util.InterfaceToString(queryResults[0][linkedObjectMapping.Builder.SingleValue.Field])
										}

									}
									newQueryTemplate := linkedObjectMapping.Query.Query
									for _, replacementField := range linkedObjectMapping.Query.ReplacementFields {
										if replacementField.Format == component.JsonToStringArray {
											value := util.InterfaceToString(trackingField.NewValue)
											replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
											newQueryTemplate = strings.Replace(newQueryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
										} else {
											value := util.InterfaceToString(trackingField.NewValue)
											replacementValue := util.InterfaceToString(value)
											newQueryTemplate = strings.Replace(newQueryTemplate, "["+replacementField.Field+"]", replacementValue, 1)
										}

									}
									var newQueryResults []map[string]interface{}
									dbConnection.Raw(newQueryTemplate).Scan(&newQueryResults)
									if linkedObjectMapping.Builder.SingleValue != nil {
										newValue = util.InterfaceToString(newQueryResults[0][linkedObjectMapping.Builder.SingleValue.Field])
									}

								}
								trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
									ChangedField: trackingField.ChangedField,
									OldValue:     oldValue,
									NewValue:     newValue,
								})
							}

						}

					case map[string]interface{}:
						var objectFields = trackingField.NewValue.(map[string]interface{})
						var modifiedValue string
						modifiedValue = "["
						var index = 0
						for key, value := range objectFields {
							index = index + 1
							if len(objectFields) == index {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "]"
							} else {
								modifiedValue += key + " : " + util.InterfaceToString(value) + "] ["
							}

						}
						trackingFieldsResponse = append(trackingFieldsResponse, TrackingFieldsResponse{
							ChangedField: trackingField.ChangedField,
							OldValue:     util.InterfaceToString(trackingField.OldValue),
							NewValue:     modifiedValue,
						})
					}
				}
			}

			recordMessageData.TrackingFields = trackingFieldsResponse
		}

		if arrayOfTimelines, ok := recordMessageMap[dateString]; ok {
			arrayOfTimelines = append(arrayOfTimelines, recordMessageData)
			recordMessageMap[dateString] = arrayOfTimelines
		} else {
			var timelineArray []RecordMessageData
			timelineArray = append(timelineArray, recordMessageData)
			recordMessageMap[dateString] = timelineArray
		}
		//April 26 2022

	}

	if len(recordMessageMap) == 0 {
		// we need to return the empty response
		var emptyArray = make([]string, 0)
		return emptyArray
	} else {
		var response []RecordMessageResponse
		for val, r := range recordMessageMap {
			recordMessageResponse := RecordMessageResponse{
				Date: val,
				Data: r,
			}
			response = append(response, recordMessageResponse)
		}
		return response
	}
}
func GetRecordVersionResponse(listOfRecords *[]component.GeneralObject, zone string) interface{} {
	recordInfo := make(map[string]interface{})
	var dropDownArray []map[string]interface{}
	for index, recordMessageInterface := range *listOfRecords {
		// Response Single drop down
		resourceInfo := ResourceInfo{}
		json.Unmarshal(recordMessageInterface.ObjectInfo, &resourceInfo)

		versionStr := resourceInfo.SourceObject.Version
		versionSplit := strings.Split(versionStr, "_")
		timeConvertedVersion := "version_" + util.ISO2TableDateTimeFormat(zone, versionSplit[1])

		dropDownArray = append(dropDownArray, map[string]interface{}{
			"id":    versionStr,
			"value": timeConvertedVersion,
		})
		if index == 0 {
			recordInfo["index"] = versionStr
			recordInfo["value"] = timeConvertedVersion
		}

	}
	recordInfo["data"] = dropDownArray

	return recordInfo
}
func RenameExportedFile(componentName string) string {

	sgtTime := util.GetCurrentTimeInSingapore("2006-01-02T15:04:05")

	if componentName == "" {
		return sgtTime
	}
	uniqueFileName := componentName + "_" + sgtTime

	return uniqueFileName
}
