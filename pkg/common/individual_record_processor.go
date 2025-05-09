package common

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (cm *ComponentManager) GenerateAdditionalRecordResponse(zone string, dbConnection *gorm.DB, objectInfo datatypes.JSON, componentName string, generalResponse map[string]interface{}) {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	additionalRecords := cm.ComponentSchema[int64ComponentId].AdditionalRecords

	if len(additionalRecords) > 0 {
		// we have configured additional record schema
		for _, individualRecordSchema := range additionalRecords {
			recordInfo := component.RecordInfo{}
			if individualRecordSchema.ObjectMapping.Predefined != nil {
				var queryResults []map[string]interface{}
				queryResults = individualRecordSchema.ObjectMapping.Predefined.Data
				var objectFields = make(map[string]interface{})
				json.Unmarshal(objectInfo, &objectFields)
				if individualRecordSchema.ObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0
					var referenceIndex int
					if individualRecordSchema.ObjectMapping.Builder.SingleDropdown.ReferenceField != "" {
						// reference field configured, get the data from there to set it
						referenceIndex = util.InterfaceToInt(objectFields[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.ReferenceField])
					}
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.Index])
						dropdownValue := queryResult[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if referenceIndex != 0 {
							if id == referenceIndex {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}
						} else {
							if index == 0 {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}
							index = index + 1
						}

					}
					recordInfo.Data = dropDownArray
				}
			} else if individualRecordSchema.ObjectMapping.Query != nil {
				var queryResults []map[string]interface{}
				var objectFields = make(map[string]interface{})
				json.Unmarshal(objectInfo, &objectFields)
				executeQuery := individualRecordSchema.ObjectMapping.Query.Query
				// before run the query, do the replacements
				for _, replacementField := range individualRecordSchema.ObjectMapping.Query.ReplacementFields {
					if replacementField.Format == component.JsonToStringArray {
						value := objectFields[replacementField.Property]
						if value != nil {
							// then we have empty array
							replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
							executeQuery = strings.Replace(executeQuery, "["+replacementField.Field+"]", replacementValue, -1)
						}

					} else {
						value := objectFields[replacementField.Property]
						replacementValue := util.InterfaceToString(value)
						executeQuery = strings.Replace(executeQuery, "["+replacementField.Field+"]", replacementValue, -1)
					}

				}
				fmt.Println("")
				fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Executing Query ", executeQuery)
				fmt.Println("")
				dbConnection.Raw(executeQuery).Scan(&queryResults)
				fmt.Println("")
				fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Result size :", len(queryResults))
				fmt.Println("")

				// before goes to builder , look for catalyst to boost the response, because it may contain the service call
				if len(individualRecordSchema.ObjectMapping.Catalyst) > 0 {
					for _, catalyst := range individualRecordSchema.ObjectMapping.Catalyst {
						authService := GetService(catalyst.Name).ServiceInterface.(AuthInterface)
						method := reflect.ValueOf(authService).MethodByName(catalyst.Call)

						for index, queryResult := range queryResults {
							inputs := make([]reflect.Value, len(catalyst.ServiceParam))
							for key, value := range queryResult {
								for serviceParamIndex, serviceParam := range catalyst.ServiceParam {
									if key == serviceParam.Field {
										if serviceParam.Type == "int" {
											inputs[serviceParamIndex] = reflect.ValueOf(util.InterfaceToInt(value))
										} else if serviceParam.Type == "text" {
											inputs[serviceParamIndex] = reflect.ValueOf(value)
										} else if serviceParam.Type == "int_array" {
											inputs[serviceParamIndex] = reflect.ValueOf(value)
										}

									}
								}

							}
							if method.IsValid() {
								returnValues := method.Call(inputs)
								fin := returnValues[0].Interface().(component.TableDataResponse)
								var responseField = make(map[string]interface{})
								json.Unmarshal(fin.Data[0], &responseField)
								for key, value := range responseField {
									queryResults[index][key] = value
								}

							}

						}

					}
				}

				if individualRecordSchema.ObjectMapping.Builder.SingleDropdown != nil {
					var dropDownArray []component.OrderedData
					index := 0
					var referenceIndex int
					if individualRecordSchema.ObjectMapping.Builder.SingleDropdown.ReferenceField != "" {
						// reference field configured, get the data from there to set it
						referenceIndex = util.InterfaceToInt(objectFields[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.ReferenceField])
					}
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.Index])
						dropdownValue := queryResult[individualRecordSchema.ObjectMapping.Builder.SingleDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})
						if referenceIndex != 0 {
							if id == referenceIndex {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}
						} else {
							if index == 0 {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}
							index = index + 1
						}

					}
					recordInfo.Data = dropDownArray
				} else if individualRecordSchema.ObjectMapping.Builder.SingleDropdownObject != nil {
					// we need to order based on id starting from low.

					existingId := util.InterfaceToInt(objectFields[individualRecordSchema.Property])
					var existingValue interface{}
					for _, queryResult := range queryResults {
						existingId = util.InterfaceToInt(queryResult["id"])

						existingValue = queryResult

					}
					recordInfo.Index = existingId
					recordInfo.Data = queryResults
					recordInfo.Value = existingValue
					recordInfo.InterfaceType = "singleDropdownObject"

				} else if individualRecordSchema.ObjectMapping.Builder.TableFieldsToObject != nil {

					// this will convert table fields into object
					var objectResponse = make(map[string]interface{}, len(individualRecordSchema.ObjectMapping.Builder.TableFieldsToObject.Fields))
					if len(queryResults) > 0 {
						for _, individualField := range individualRecordSchema.ObjectMapping.Builder.TableFieldsToObject.Fields {
							dataType := individualField.Type
							fieldName := individualField.Name
							individualFieldValue := queryResults[0][fieldName]
							if dataType == "bool" {
								objectResponse[fieldName] = util.InterfaceToBool(individualFieldValue)
							} else if dataType == "int" {
								objectResponse[fieldName] = util.InterfaceToInt(individualFieldValue)
							} else if dataType == "double" {
								objectResponse[fieldName] = util.InterfaceToFloat(individualField)
							} else {
								objectResponse[fieldName] = individualFieldValue
							}

						}
					}

					recordInfo.Data = objectResponse

				} else if individualRecordSchema.ObjectMapping.Builder.Table != nil {
					header := individualRecordSchema.ObjectMapping.Builder.Table.Schema
					for index, tableFieldSchema := range header {
						if tableFieldSchema.HeaderObjectMapping != nil {
							if tableFieldSchema.HeaderObjectMapping.Query != nil {
								var queryResults []map[string]interface{}
								dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)
								var objectKeyValues = make(map[string]string, len(queryResults))
								for _, result := range queryResults {

									switch v := result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(type) {
									case int:
										objectKeyValues[strconv.Itoa(v)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
									case int32:
										objectKeyValues[strconv.Itoa(int(v))] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
									case float64:
										objectKeyValues[strconv.Itoa(int(v))] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
									default:
										objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
									}

								}
								raw, _ := json.Marshal(objectKeyValues)
								header[index].ObjectList = raw
								// set the object to null, so that response won't send that field
								header[index].HeaderObjectMapping = nil
							} else if tableFieldSchema.HeaderObjectMapping.Predefined != nil {
								predefinedResults := tableFieldSchema.HeaderObjectMapping.Predefined.Data
								var objectKeyValues = make(map[string]string, len(predefinedResults))
								if tableFieldSchema.HeaderObjectMapping.Builder.KeyValue != nil {
									for _, result := range predefinedResults {
										switch v := result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(type) {
										case int:
											objectKeyValues[strconv.Itoa(v)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
										case int32:
											objectKeyValues[strconv.Itoa(int(v))] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
										case float64:
											objectKeyValues[strconv.Itoa(int(v))] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
										default:
											objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
										}

									}
								}

								raw, _ := json.Marshal(objectKeyValues)
								header[index].ObjectList = raw
								// set the object to null, so that response won't send that field
								//header[index].HeaderObjectMapping = nil
								//header[index].ReferenceObjectMapping = nil
							}

						}
					}

					recordInfo.Header = header
					if individualRecordSchema.ObjectMapping.Builder.Table.CommonRouteLink != "" {
						// we have configured the common route link, so send that
						recordInfo.CommonRouteLink = individualRecordSchema.ObjectMapping.Builder.Table.CommonRouteLink
					}
					// iterate through each results
					// if not results obtained , then send the empty array rather not sending data object
					var canExport = false
					if len(queryResults) == 0 {
						recordInfo.Data = make([]interface{}, 0)
					} else {
						var dataObject []interface{}
						for _, results := range queryResults {
							// we need to check the schema for any route link
							var internalRecords = make(map[string]interface{}, 0)
							for _, individualHeader := range header {
								if individualHeader.ReferenceObjectMapping != nil {
									// we have specified the mapping
									if individualHeader.ReferenceObjectMapping.Query != nil {
										query := individualHeader.ReferenceObjectMapping.Query.Query
										var queryResults []map[string]interface{}
										for _, replacementField := range individualHeader.ReferenceObjectMapping.Query.ReplacementFields {
											if replacementField.Format == component.JsonToStringArray {
												value := results[replacementField.Property]
												replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
												query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
											} else {
												value := results[replacementField.Property]
												replacementValue := util.InterfaceToString(value)
												query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
											}

										}

										dbConnection.Raw(query).Scan(&queryResults)
										if individualHeader.ReferenceObjectMapping.Builder.SingleValue != nil {
											if len(queryResults) > 0 {
												internalRecords[individualHeader.Property] = queryResults[0][individualHeader.ReferenceObjectMapping.Builder.SingleValue.Field]
											}
										}
									} else if individualHeader.ReferenceObjectMapping.Predefined != nil {
										predefinedResults := individualHeader.ReferenceObjectMapping.Predefined.Data
										singleValueFromIdBuilder := individualHeader.ReferenceObjectMapping.Builder
										if singleValueFromIdBuilder.SingleValueFromId != nil {
											if len(predefinedResults) > 0 {
												idField := util.InterfaceToInt(results[singleValueFromIdBuilder.SingleValueFromId.IdField])
												var isFieldFound bool
												for index, internalPredefinedResult := range predefinedResults {
													id := util.InterfaceToInt(internalPredefinedResult["id"])
													if id == idField {
														internalRecords[individualHeader.Property] = predefinedResults[index][singleValueFromIdBuilder.SingleValueFromId.Field]
														isFieldFound = true
													}
												}
												if !isFieldFound {
													internalRecords[individualHeader.Property] = "-"
												}
											}
										}
									}

								} else {
									// once the field is updated, lets check the type

									// once the field is updated, lets check the type
									if individualHeader.Type == component.TableDataTypeDateTime {
										// should be in the user timezone format
										if value, ok := results[individualHeader.Property]; ok {
											if value != nil {
												currentValue := value.(string)
												var correctedValue string
												if currentValue != "" {
													correctedValue = util.ISO2TableDateTimeFormat(zone, currentValue)
												}
												internalRecords[individualHeader.Property] = correctedValue
											}

										}

									} else if individualHeader.Type == component.TableDataTypeDate {
										// should be in the user timezone format
										if value, ok := results[individualHeader.Property]; ok {
											if value != nil {
												currentValue := value.(string)
												var correctedValue string
												if currentValue != "" {
													correctedValue = util.ISO2TableDateTimeFormat(zone, currentValue)
												}
												internalRecords[individualHeader.Property] = correctedValue
											}

										}

									} else if individualHeader.Type == "bool" {
										internalRecords[individualHeader.Property] = util.InterfaceToBool(results[individualHeader.Property])
									} else if individualHeader.Type == "object_array" {
										var jsonObjectFields = make([]map[string]interface{}, 0)
										var stringObject = util.InterfaceToString(results[individualHeader.Property])
										json.Unmarshal([]byte(stringObject), &jsonObjectFields)
										internalRecords[individualHeader.Property] = jsonObjectFields
									} else {
										internalRecords[individualHeader.Property] = results[individualHeader.Property]
									}

								}
							}

							dataObject = append(dataObject, internalRecords)
							if len(dataObject) > 0 {
								canExport = true
							}
						}

						recordInfo.Data = dataObject

					}

					if individualRecordSchema.InternalExport != nil {
						recordInfo.InternalExport = individualRecordSchema.InternalExport
						recordInfo.InternalExport.CanExport = canExport
					}

				} else if individualRecordSchema.ObjectMapping.Builder.SingleValue != nil {
					if len(queryResults) > 0 {
						if len(individualRecordSchema.ObjectMapping.Query.OutFieldsMapping) > 0 {
							// first apply the out field mapping logics using formattters
							for _, outField := range individualRecordSchema.ObjectMapping.Query.OutFieldsMapping {
								for index, result := range queryResults {
									if queryResultValue, ok := result[outField.Field]; ok {
										queryResults[index][outField.Field] = applyFormatters(outField.Formatter, queryResultValue)
									}
								}
							}
						}
						if individualRecordSchema.ObjectMapping.Builder.SingleValue.Type == "datetime" {
							assignedDateTime := util.InterfaceToString(queryResults[0][individualRecordSchema.ObjectMapping.Builder.SingleValue.Field])
							recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", assignedDateTime)

						} else if individualRecordSchema.ObjectMapping.Builder.SingleValue.Type == "date" {
							assignedDate := util.InterfaceToString(queryResults[0][individualRecordSchema.ObjectMapping.Builder.SingleValue.Field])
							recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", assignedDate)
						} else if individualRecordSchema.ObjectMapping.Builder.SingleValue.Type == "image" {
							imageUrl := util.InterfaceToString(queryResults[0][individualRecordSchema.ObjectMapping.Builder.SingleValue.Field])
							upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
							fmt.Println("upstreamResponse:", string(upstreamResponse))
							if err == nil {
								var upstreamResponseFields = make(map[string]interface{})
								json.Unmarshal(upstreamResponse, &upstreamResponseFields)
								var fileMetaInfo = make(map[string]interface{})
								fileMetaInfo["name"] = upstreamResponseFields["name"]
								fileMetaInfo["size"] = upstreamResponseFields["size"]
								fileMetaInfo["url"] = upstreamResponseFields["url"]
								recordInfo.Data = fileMetaInfo
								recordInfo.Value = queryResults[0][individualRecordSchema.ObjectMapping.Builder.SingleValue.Field]
							}
						} else {
							recordInfo.Value = queryResults[0][individualRecordSchema.ObjectMapping.Builder.SingleValue.Field]
						}

					} else {
						if individualRecordSchema.ObjectMapping.Builder.SingleValue.Type == "int" {
							recordInfo.Value = 0
						} else if individualRecordSchema.ObjectMapping.Builder.SingleValue.Type == "text" {
							recordInfo.Value = ""
						}

						fmt.Println(" error in query resutls :", queryResults, " field mapping :", individualRecordSchema.ObjectMapping.Builder.SingleValue.Field)
					}

				} else if individualRecordSchema.ObjectMapping.Builder.GroupBy != nil {
					executeQuery := individualRecordSchema.ObjectMapping.Query.Query
					// before run the query, do the replacements
					for _, replacementField := range individualRecordSchema.ObjectMapping.Query.ReplacementFields {
						fmt.Println("replacementField:", replacementField)
						if replacementField.Format == component.JsonToStringArray {
							value := objectFields[replacementField.Property]
							replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
							executeQuery = strings.Replace(executeQuery, "["+replacementField.Field+"]", replacementValue, 1)
						} else {
							value := objectFields[replacementField.Property]
							fmt.Println("replacement value : ", value)
							replacementValue := util.InterfaceToString(value)
							executeQuery = strings.Replace(executeQuery, "["+replacementField.Field+"]", replacementValue, 1)
						}

					}
					groupQueryResults := make([]map[string]interface{}, 0)
					fmt.Println("executing query:", executeQuery)
					dbConnection.Raw(executeQuery).Scan(&groupQueryResults)
					groupByColumn := individualRecordSchema.ObjectMapping.Builder.GroupBy.GroupByColumnName
					idColumn := individualRecordSchema.ObjectMapping.Builder.GroupBy.Id
					valueColumn := individualRecordSchema.ObjectMapping.Builder.GroupBy.Value
					var groupByResults = make(map[string][]component.RecordInfo, 20)
					for _, results := range groupQueryResults {
						groupByValue := util.InterfaceToString(results[groupByColumn])
						idValue := util.InterfaceToInt(results[idColumn])
						valueValue := util.InterfaceToString(results[valueColumn])

						groupByResults[groupByValue] = append(groupByResults[groupByValue], component.RecordInfo{
							Value: valueValue,
							Index: idValue,
						})

					}
					recordInfo.Data = groupByResults
					recordInfo.IsExternal = true

				}
			} else if individualRecordSchema.ObjectMapping.Service != nil {
				serviceName := individualRecordSchema.ObjectMapping.Service.Name
				serviceInterface := GetService(serviceName).ServiceInterface.(AuthInterface)

				method := reflect.ValueOf(serviceInterface).MethodByName(individualRecordSchema.ObjectMapping.Service.Call)
				fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Calling service ", individualRecordSchema.ObjectMapping.Service.Call)
				var objectFields = make(map[string]interface{})
				json.Unmarshal(objectInfo, &objectFields)
				if method.IsValid() {
					fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Service Param", individualRecordSchema.ObjectMapping.Service.ServiceParam)

					inputs := make([]reflect.Value, len(individualRecordSchema.ObjectMapping.Service.ServiceParam))

					for serviceParamIndex, serviceParam := range individualRecordSchema.ObjectMapping.Service.ServiceParam {
						if serviceParam.Type == "int_array" {
							intSlice := make([]int, 0)
							sliceType := reflect.TypeOf(intSlice)
							intSliceReflect := reflect.MakeSlice(sliceType, 0, 0)
							if objectValue, ok := objectFields[serviceParam.Field]; ok {
								if objectValue == nil {
									inputs[serviceParamIndex] = reflect.ValueOf(intSliceReflect)
								} else {
									intArray := util.InterfaceToIntArray(objectValue)
									for _, intValue := range intArray {
										intSliceReflect = reflect.Append(intSliceReflect, reflect.ValueOf(intValue))
									}

									inputs[serviceParamIndex] = reflect.ValueOf(intSliceReflect)
									fmt.Println("inputs[serviceParamIndex]:", inputs[serviceParamIndex])
									fmt.Println("inputs[inputs]:", inputs)
								}

							} else {
								inputs[serviceParamIndex] = reflect.ValueOf(intSliceReflect)
							}
						} else if serviceParam.Type == "text" {
							var emptyString = ""
							if objectValue, ok := objectFields[serviceParam.Field]; ok {

								if objectValue == nil {
									inputs[serviceParamIndex] = reflect.ValueOf(emptyString)
								} else {
									inputs[serviceParamIndex] = reflect.ValueOf(objectValue)
								}

								inputs[serviceParamIndex] = reflect.ValueOf(objectValue)
							} else {
								inputs[serviceParamIndex] = reflect.ValueOf(emptyString)
							}
						} else if serviceParam.Type == "int" {
							var emptyInt = 0
							if objectValue, ok := objectFields[serviceParam.Field]; ok {
								if objectValue == nil {
									inputs[serviceParamIndex] = reflect.ValueOf(emptyInt)
								} else {
									inputs[serviceParamIndex] = reflect.ValueOf(objectValue)
								}
								inputs[serviceParamIndex] = reflect.ValueOf(util.InterfaceToInt(objectValue))
							} else {

								inputs[serviceParamIndex] = reflect.ValueOf(emptyInt)
							}
						}
					}

					fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Calling service with inputs ", len(inputs))
					fmt.Println("[ADDITIONAL_RECORD_RESPONSE] Calling service with inputs ", inputs)
					returnValues := method.Call(inputs)
					fmt.Println("returnValues: ", returnValues)
					fin := returnValues[0].Interface().(component.TableDataResponse)

					recordInfo.Data = fin.Data
					recordInfo.Header = fin.Header
					recordInfo.IsExternal = true
					recordInfo.IsEdit = false

				} else {
					recordInfo.Data = make([]string, 0)
					recordInfo.Header = make([]string, 0)
					recordInfo.IsExternal = true
					recordInfo.IsEdit = false
				}

			} else if len(individualRecordSchema.RecordSchema) > 0 {
				// we have record schema defined as an object
				// we should have record schema configured here.
				var objectFields = make(map[string]interface{})
				json.Unmarshal(objectInfo, &objectFields)
				internalObjectField := objectFields[individualRecordSchema.Property]
				serialisedInternalObjectField, _ := json.Marshal(internalObjectField)
				responseInternal := make(map[string]interface{}, 0)
				json.Unmarshal(serialisedInternalObjectField, &responseInternal)
				//responseInternal := make(map[string]interface{}, 0)
				//var insideObject = make(map[string]interface{}, 0)
				//serializedObject, _ := json.Marshal(recordMap[individualRecordSchema.Property])
				//json.Unmarshal(serializedObject, &insideObject)
				//property := individualRecordSchema.Property
				//newResponse := cm.BuildGeneralRecordResponse(dbConnection, insideObject, responseInternal, individualRecordSchema.RecordSchema)
				//response[property] = newResponse
				for _, internalRecordSchema := range individualRecordSchema.RecordSchema {
					internalRecordInfo := component.RecordInfo{}
					internalRecordInfo.Type = internalRecordSchema.Type
					internalRecordInfo.IsEdit = internalRecordSchema.IsEdit
					if internalRecordSchema.Default != nil {
						internalRecordInfo.Value = internalRecordSchema.Default
					} else {
						internalRecordInfo.Data = responseInternal[internalRecordSchema.Property]
					}
				}
				recordInfo.Data = responseInternal
				recordInfo.Value = nil
				recordInfo.IsExternal = true
			}

			generalResponse[individualRecordSchema.Property] = recordInfo

		}
	}

}

func (cm *ComponentManager) GenerateTemplateResponse(dbConnection *gorm.DB, objectInfo datatypes.JSON, componentName string, generalResponse map[string]interface{}) {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	referenceTemplateComponent := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent
	var objectFields = make(map[string]interface{})
	json.Unmarshal(objectInfo, &objectFields)
	var templateId = -1
	fmt.Println("generate template response")
	// we are going to add any template fields configured
	if referenceTemplateComponent != nil {
		templateFieldName := referenceTemplateComponent.TemplateFieldName
		if value, ok := objectFields[templateFieldName]; ok {
			templateId = util.InterfaceToInt(value)
			referenceTable := cm.GetTargetTable(referenceTemplateComponent.ReferenceComponent) //it_service_category_template
			err, templateTableFieldsObject := Get(dbConnection, referenceTable, templateId)
			if err == nil {

				// got all the template fields
				templateFields := GetObjectFields(templateTableFieldsObject.ObjectInfo)
				templateFieldRecords := templateFields[templateFieldName]
				serialisedRecords := GetInterfaceToSerialisation(templateFieldRecords)
				var listOfTableRecords []TemplateRecords
				json.Unmarshal(serialisedRecords, &listOfTableRecords)

				//We have to sort the field then send new record object
				//sort.Slice(listOfTableRecords, func(i, j int) bool {
				//	return listOfTableRecords[i].Id < listOfTableRecords[j].Id
				//})

				var listOfDynamicRecords []component.RecordInfo
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

						if fieldValue, valueOk := objectFields[templateRecord.Property]; valueOk {
							recordInfo.Value = fieldValue
						}

						// if it index array. Checkbox values are index array
						if templateRecord.InterfaceFieldList == CheckBox {
							if fieldValue, valueOk := objectFields[templateRecord.Property]; valueOk {
								recordInfo.IndexArray = util.InterfaceToIntArray(fieldValue)
							}
						} else if templateRecord.InterfaceFieldList == FieldTypeFile {
							if templateRecord.InterfaceType == SingleDropdown {
								imageUrl := util.InterfaceToString(objectFields[templateRecord.Property])
								upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
								fmt.Println("upstreamResponse:", string(upstreamResponse))
								if err == nil {
									var upstreamResponseFields = make(map[string]interface{})
									json.Unmarshal(upstreamResponse, &upstreamResponseFields)
									var fileMetaInfo = make(map[string]interface{})
									fileMetaInfo["name"] = upstreamResponseFields["name"]
									fileMetaInfo["size"] = upstreamResponseFields["size"]
									fileMetaInfo["url"] = upstreamResponseFields["url"]
									recordInfo.Data = fileMetaInfo
								}
							} else {
								listOfImageUrl := util.InterfaceToStringArray(objectFields[templateRecord.Property])
								listOfDataUrl := make([]map[string]interface{}, 0)
								firstIndex := 0
								for index, imageUrl := range listOfImageUrl {
									upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
									fmt.Println("upstreamResponse:", string(upstreamResponse))
									if err == nil {
										if index == 0 {
											firstIndex = 1
										}
										var upstreamResponseFields = make(map[string]interface{})
										json.Unmarshal(upstreamResponse, &upstreamResponseFields)
										var fileMetaInfo = make(map[string]interface{})
										fileMetaInfo["name"] = upstreamResponseFields["name"]
										fileMetaInfo["size"] = upstreamResponseFields["size"]
										fileMetaInfo["url"] = upstreamResponseFields["url"]
										listOfDataUrl = append(listOfDataUrl, fileMetaInfo)

									}
								}
								recordInfo.Data = listOfDataUrl
								recordInfo.Index = firstIndex
							}

						} else if templateRecord.InterfaceFieldList == FieldTypeDate {
							if fieldValue, valueOk := objectFields[templateRecord.Property]; valueOk {
								recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", util.InterfaceToString(fieldValue))
							}
						}

						listOfTemplateFields[templateRecord.Property] = recordInfo
					}
				}
				for _, templateRecord := range listOfTableRecords {
					var isRenderEnabled bool
					isRenderEnabled = false
					// If the templateRecord.EnabledActionFieldNames  is exist in the object, and the value is true, then render it otherwise
					// don't render that.
					if templateRecord.EnabledAfterWorkflowStatusLevel == nil {
						isRenderEnabled = true
					} else {
						// any of the configured field exist, then render true

						if actionFieldValue, ok := objectFields["levelCounter"]; ok {
							entityLevel := util.InterfaceToInt(actionFieldValue) // assume object has the level 3
							if templateRecord.EnabledAfterWorkflowStatusLevel != nil {
								if entityLevel >= *templateRecord.EnabledAfterWorkflowStatusLevel {
									isRenderEnabled = true
								} else {
									isRenderEnabled = false
								}
							}
						} else {
							// if the fields are not configured with correct status field
							isRenderEnabled = true
						}

					}

					// got each of the template records
					recordInfo := component.RecordInfo{}
					recordInfo.Property = templateRecord.Property
					recordInfo.Type = getInt2DataType(templateRecord.DataType)
					recordInfo.Label = templateRecord.Label
					recordInfo.GridSystem = getGridSystem2Str(templateRecord.GridSystem)
					recordInfo.Render = isRenderEnabled
					recordInfo.IsDynamic = true
					recordInfo.IsMandatoryField = templateRecord.IsMandatoryField
					recordInfo.Description = templateRecord.Description
					recordInfo.InterfaceType = getInterfaceType2Str(templateRecord.InterfaceTypeList)
					recordInfo.InterfaceField = getInt2FieldType(templateRecord.InterfaceFieldList)

					if objectValue, ok := objectFields[templateRecord.Property]; ok {
						recordInfo.Value = objectValue
					}

					// if it index array. Checkbox values are index array
					if templateRecord.InterfaceFieldList == CheckBox {
						if fieldValue, valueOk := objectFields[templateRecord.Property]; valueOk {
							recordInfo.IndexArray = util.InterfaceToIntArray(fieldValue)
						}
					} else if templateRecord.InterfaceFieldList == FieldTypeFile {
						if templateRecord.InterfaceTypeList == SingleDropdown {
							imageUrl := util.InterfaceToString(objectFields[templateRecord.Property])
							upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
							if err == nil {
								var upstreamResponseFields = make(map[string]interface{})
								json.Unmarshal(upstreamResponse, &upstreamResponseFields)
								var fileMetaInfo = make(map[string]interface{})
								fileMetaInfo["name"] = upstreamResponseFields["name"]
								fileMetaInfo["size"] = upstreamResponseFields["size"]
								fileMetaInfo["url"] = upstreamResponseFields["url"]
								recordInfo.Data = fileMetaInfo
							}
						} else {
							listOfImageUrl := util.InterfaceToStringArray(objectFields[templateRecord.Property])
							listOfDataUrl := make([]map[string]interface{}, 0)
							firstIndex := 0
							for index, imageUrl := range listOfImageUrl {
								upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
								fmt.Println("upstreamResponse:", string(upstreamResponse))
								if err == nil {
									if index == 0 {
										firstIndex = 1
									}
									var upstreamResponseFields = make(map[string]interface{})
									json.Unmarshal(upstreamResponse, &upstreamResponseFields)
									var fileMetaInfo = make(map[string]interface{})
									fileMetaInfo["name"] = upstreamResponseFields["name"]
									fileMetaInfo["size"] = upstreamResponseFields["size"]
									fileMetaInfo["url"] = upstreamResponseFields["url"]
									listOfDataUrl = append(listOfDataUrl, fileMetaInfo)

								}
							}
							recordInfo.Data = listOfDataUrl
							recordInfo.Index = firstIndex
						}

					} else if templateRecord.InterfaceFieldList == FieldTypeDate {
						if fieldValue, valueOk := objectFields[templateRecord.Property]; valueOk {
							recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", util.InterfaceToString(fieldValue))

						}
					}

					var condFieldMapProperty = make(map[string][]component.RecordInfo)
					for _, templateRecordOfExisting := range listOfTableRecords {
						if templateRecordOfExisting.InterfaceTypeList == 5 || templateRecord.InterfaceTypeList == 6 {
							if len(templateRecordOfExisting.DynamicDroppedDownAttributes.ConditionalFields) > 0 {
								conditionalFields := templateRecordOfExisting.DynamicDroppedDownAttributes.ConditionalFields
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

					if templateRecord.InterfaceTypeList == 5 || templateRecord.InterfaceTypeList == 6 {
						//recordInfo.DynamicConditionalFields = templateRecord.DynamicDroppedDownAttributes.ConditionalFields // later we need to remove sending some critical information
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
									index := 1
									existingId := util.InterfaceToInt(objectFields[templateRecord.Property])
									for _, sourceObject := range listOfSourceObjects {
										id := sourceObject.Id
										var objectDropFields = make(map[string]interface{})
										json.Unmarshal(sourceObject.ObjectInfo, &objectDropFields)
										if valueDrop, okDrop := objectDropFields[sourceField]; okDrop {
											// now check the value is
											if conValue, ok := condFieldMapProperty[util.InterfaceToString(valueDrop)]; ok {
												dropDownArray = append(dropDownArray, component.OrderedData{
													Id:                       index,
													Value:                    util.InterfaceToString(valueDrop),
													OnValueConditionalFields: conValue,
												})
											} else {
												dropDownArray = append(dropDownArray, component.OrderedData{
													Id:    index,
													Value: util.InterfaceToString(valueDrop),
												})
											}
											if index == existingId {
												recordInfo.Index = id
												recordInfo.Value = util.InterfaceToString(valueDrop)
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
							existingId := util.InterfaceToInt(objectFields[templateRecord.Property])

							for index, valueDrop := range templateRecord.DynamicDroppedDownAttributes.ManualFieldsSource {
								index = index + 1

								if conValue, ok := condFieldMapProperty[util.InterfaceToString(valueDrop)]; ok {
									dropDownArray = append(dropDownArray, component.OrderedData{
										Id:                       index,
										Value:                    util.InterfaceToString(valueDrop),
										OnValueConditionalFields: conValue,
									})
								} else {
									dropDownArray = append(dropDownArray, component.OrderedData{
										Id:    index,
										Value: util.InterfaceToString(valueDrop),
									})
								}
								if index == existingId {
									recordInfo.Index = existingId
									recordInfo.Value = util.InterfaceToString(valueDrop)
								}
								//index = index + 1
							}
							recordInfo.Data = dropDownArray
						}
					}
					recordInfo.Display = true
					fmt.Println("Final", recordInfo.Value)

					if !templateRecord.IsDroppedDownConditionalField {
						listOfDynamicRecords = append(listOfDynamicRecords, recordInfo)
					}
				}

				generalResponse["dynamicFields"] = listOfDynamicRecords
			}
		}

	}

}

func (cm *ComponentManager) BuildGeneralRecordResponse(dbConnection *gorm.DB, recordMap map[string]interface{}, response map[string]interface{}, recordSchema []component.RecordSchema) map[string]interface{} {
	// normally if the field is linked field, in the record, we will see the id
	for _, individualRecordSchema := range recordSchema {
		recordInfo := component.RecordInfo{}
		// this is configured as not send to front-end during , if so, just ignore
		if individualRecordSchema.IgnoreEmptyInSending {
			continue
		}
		if individualRecordSchema.ResponseObjectMapping != nil {

			if individualRecordSchema.ResponseObjectMapping.Query != nil {

				query := individualRecordSchema.ResponseObjectMapping.Query.Query
				for _, replacementField := range individualRecordSchema.ResponseObjectMapping.Query.ReplacementFields {

					if replacementField.Format == component.JsonToStringArray {
						value := recordMap[replacementField.Property]
						replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
						query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
					} else {
						value := recordMap[replacementField.Property]
						replacementValue := util.InterfaceToString(value)
						query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
					}

				}
				fmt.Println("")
				fmt.Println("[BUILD_GENERAL_RECORD] Composed Query :", query)
				fmt.Println("")
				var queryResults []map[string]interface{}
				dbConnection.Raw(query).Scan(&queryResults)

				fmt.Println("")
				fmt.Println("[BUILD_GENERAL_RECORD] Query results size :", len(queryResults))
				fmt.Println("")

				if individualRecordSchema.ResponseObjectMapping.Builder.SingleValue != nil {
					if len(queryResults) > 0 {
						fieldValue := queryResults[0][individualRecordSchema.ResponseObjectMapping.Builder.SingleValue.Field]
						recordInfo.Value = fieldValue
					} else {
						recordInfo.Value = "-"
					}

				} else if individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown != nil {
					// we need to order based on id starting from low.
					dropDownArray := make([]component.OrderedData, 0)
					existingId := util.InterfaceToInt(recordMap[individualRecordSchema.Property])
					for _, queryResult := range queryResults {
						id := int(queryResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Index].(int32))
						if dropdownValue, ok := queryResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Value]; ok {
							if dropdownValue != nil {
								dropDownArray = append(dropDownArray, component.OrderedData{
									Id:    id,
									Value: util.InterfaceToString(dropdownValue),
								})
								if id == existingId {
									recordInfo.Index = id
									recordInfo.Value = dropdownValue
								}
							}

						}

					}
					recordInfo.IsDynamic = individualRecordSchema.IsDynamic
					recordInfo.DynamicMappingField = individualRecordSchema.DynamicMappingField
					recordInfo.DynamicComponent = individualRecordSchema.DynamicComponent
					recordInfo.Data = dropDownArray
				} else if individualRecordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.
					existingIds := util.InterfaceToIntArray(recordMap[individualRecordSchema.Property])

					if len(existingIds) == 0 {
						recordInfo.IndexArray = make([]int, 0)
					} else {
						recordInfo.IndexArray = existingIds
					}
					recordInfo.Data = queryResults
					recordInfo.InterfaceType = "multiDropdownObject"

				} else if individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdownObject != nil {
					// we need to order based on id starting from low.
					existingId := util.InterfaceToInt(recordMap[individualRecordSchema.Property])
					var existingValue interface{}
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult["id"])
						if id == existingId {
							existingValue = queryResult
						}
					}
					recordInfo.Index = existingId
					recordInfo.Data = queryResults
					recordInfo.Value = existingValue
					recordInfo.InterfaceType = "singleDropdownObject"

				} else if individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					existingIds := util.InterfaceToIntArray(recordMap[individualRecordSchema.Property])
					recordInfo.IndexArray = existingIds
					var dropDownArray []component.OrderedData
					for _, queryResult := range queryResults {
						id := util.InterfaceToInt(queryResult[individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := queryResult[individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.InterfaceType = individualRecordSchema.InterfaceType
					recordInfo.Data = dropDownArray
				} else if individualRecordSchema.ResponseObjectMapping.Builder.Table != nil {
					header := individualRecordSchema.ResponseObjectMapping.Builder.Table.Schema
					for index, tableFieldSchema := range header {
						if tableFieldSchema.HeaderObjectMapping != nil {
							if tableFieldSchema.HeaderObjectMapping.Query != nil {
								var queryResults []map[string]interface{}
								dbConnection.Raw(tableFieldSchema.HeaderObjectMapping.Query.Query).Scan(&queryResults)
								var objectKeyValues = make(map[string]string, len(queryResults))
								for _, result := range queryResults {
									objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)

								}
								raw, _ := json.Marshal(objectKeyValues)
								header[index].ObjectList = raw
								// set the object to null, so that response won't send that field
								header[index].HeaderObjectMapping = nil
							} else if tableFieldSchema.HeaderObjectMapping.Predefined != nil {
								predefinedResults := tableFieldSchema.HeaderObjectMapping.Predefined.Data
								var objectKeyValues = make(map[string]string, len(predefinedResults))
								if tableFieldSchema.HeaderObjectMapping.Builder.KeyValue != nil {
									for _, result := range predefinedResults {
										objectKeyValues[result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Key].(string)] = result[tableFieldSchema.HeaderObjectMapping.Builder.KeyValue.Value].(string)
									}
								}

								raw, _ := json.Marshal(objectKeyValues)
								header[index].ObjectList = raw
								// set the object to null, so that response won't send that field
								header[index].HeaderObjectMapping = nil
								header[index].ReferenceObjectMapping = nil
							}

						}
					}

					recordInfo.Header = header
					if individualRecordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink != "" {
						// we have configured the common route link, so send that
						recordInfo.CommonRouteLink = individualRecordSchema.ResponseObjectMapping.Builder.Table.CommonRouteLink
					}
					// iterate through each results
					// if not results obtained , then send the empty array rather not sending data object
					if len(queryResults) == 0 {
						recordInfo.Data = make([]interface{}, 0)
					} else {
						var dataObject []interface{}
						for _, results := range queryResults {
							// we need to check the schema for any route link
							var internalRecords = make(map[string]interface{}, 0)
							for _, individualHeader := range header {

								if individualHeader.ReferenceObjectMapping != nil {
									// we have specified the mapping
									query := individualHeader.ReferenceObjectMapping.Query.Query
									var queryResults []map[string]interface{}
									for _, replacementField := range individualHeader.ReferenceObjectMapping.Query.ReplacementFields {
										if replacementField.Format == component.JsonToStringArray {
											value := results[replacementField.Property]
											replacementValue := util.InterfaceArrayToCommaSeperatedString(value)
											query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
										} else {
											value := results[replacementField.Property]
											replacementValue := util.InterfaceToString(value)
											query = strings.Replace(query, "["+replacementField.Field+"]", replacementValue, 1)
										}

									}

									dbConnection.Raw(query).Scan(&queryResults)
									fmt.Println("query :", query, " queryResults : ", queryResults)
									if individualHeader.ReferenceObjectMapping.Builder.SingleValue != nil {
										if len(queryResults) > 0 {
											internalRecords[individualHeader.Property] = queryResults[0][individualHeader.ReferenceObjectMapping.Builder.SingleValue.Field]
										}
									}

								} else {
									if individualHeader.Type == "datetime" {
										assignedDateTime := util.InterfaceToString(results[individualHeader.Property])
										internalRecords[individualHeader.Property] = util.ConvertTimeToTimeZonCorrectedPrimeNgTable("Asia/Singapore", assignedDateTime)
									} else if individualHeader.Type == "date" {
										assignedDate := util.InterfaceToString(results[individualHeader.Property])
										internalRecords[individualHeader.Property] = util.ConvertTimeToTimeZonCorrectedPrimeNgTable("Asia/Singapore", assignedDate)
									} else {
										internalRecords[individualHeader.Property] = results[individualHeader.Property]
									}

								}
							}

							dataObject = append(dataObject, internalRecords)
						}
						recordInfo.Data = dataObject

					}

				}
			} else if individualRecordSchema.ResponseObjectMapping.Predefined != nil {
				if individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown != nil {
					// we need to order based on id starting from low.
					if individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.IsReverseMapping {
						// since it is reverse mapped enabled, we should look for value and get the id
						var dropDownArray []component.OrderedData
						existingValue := util.InterfaceToString(recordMap[individualRecordSchema.Property])
						for _, predefinedResult := range individualRecordSchema.ResponseObjectMapping.Predefined.Data {
							id := util.InterfaceToInt(predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Index])
							dropdownValue := predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Value]
							dropDownArray = append(dropDownArray, component.OrderedData{
								Id:    id,
								Value: dropdownValue.(string),
							})
							if dropdownValue.(string) == existingValue {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}

						}
						recordInfo.Data = dropDownArray
					} else {
						var dropDownArray []component.OrderedData
						existingId := util.InterfaceToInt(recordMap[individualRecordSchema.Property])
						for _, predefinedResult := range individualRecordSchema.ResponseObjectMapping.Predefined.Data {
							id := util.InterfaceToInt(predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Index])
							dropdownValue := predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.SingleDropdown.Value]
							dropDownArray = append(dropDownArray, component.OrderedData{
								Id:    id,
								Value: dropdownValue.(string),
							})
							if id == existingId {
								recordInfo.Index = id
								recordInfo.Value = dropdownValue
							}

						}
						recordInfo.Data = dropDownArray
					}

				} else if individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown != nil {
					// we need to order based on id starting from low.
					existingIds := util.InterfaceToIntArray(recordMap[individualRecordSchema.Property])
					recordInfo.IndexArray = existingIds
					var dropDownArray []component.OrderedData
					for _, predefinedResult := range individualRecordSchema.ResponseObjectMapping.Predefined.Data {
						id := util.InterfaceToInt(predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Index])
						dropdownValue := predefinedResult[individualRecordSchema.ResponseObjectMapping.Builder.MultiValueDropdown.Value]
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    id,
							Value: dropdownValue.(string),
						})

					}
					recordInfo.Data = dropDownArray
				}
			} else if individualRecordSchema.ResponseObjectMapping.Service != nil {
				authService := GetService(individualRecordSchema.ResponseObjectMapping.Service.Name).ServiceInterface.(AuthInterface)
				method := reflect.ValueOf(authService).MethodByName(individualRecordSchema.ResponseObjectMapping.Service.Call)

				inputs := make([]reflect.Value, len(individualRecordSchema.ResponseObjectMapping.Service.ServiceParam))
				var queryResults []datatypes.JSON
				if method.IsValid() {
					returnValues := method.Call(inputs)
					queryResults = returnValues[0].Interface().([]datatypes.JSON)

				}
				fmt.Println("queryResults: queryResults: queryResults: queryResults: ", queryResults)
				if individualRecordSchema.ResponseObjectMapping.Builder.MultiDropdownObject != nil {
					// we need to order based on id starting from low.

					existingIds := util.InterfaceToIntArray(recordMap[individualRecordSchema.Property])
					//valueArray := individualRecordSchema.ResponseObjectMapping.Builder.MultiDropdownObject.Values

					recordInfo.IndexArray = existingIds
					recordInfo.Data = queryResults
					recordInfo.InterfaceType = "multiDropdownObject"
					fmt.Println("recordInfo: recordInfo: recordInfo: recordInfo: ", recordInfo)
				}
			}

		} else {
			if individualRecordSchema.Type == "datetime" {
				assignedDateTime := util.InterfaceToString(recordMap[individualRecordSchema.Property])
				recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", assignedDateTime)

			} else if individualRecordSchema.Type == "date" {
				assignedDate := util.InterfaceToString(recordMap[individualRecordSchema.Property])
				if assignedDate != "" {
					recordInfo.Value = util.ConvertTimeToTimeZonCorrectedPrimeNg("Asia/Singapore", assignedDate)
				} else {
					recordInfo.Value = ""
				}

			} else if individualRecordSchema.Type == "image" {

				imageUrl := util.InterfaceToString(recordMap[individualRecordSchema.Property])
				upstreamResponse, err := util.UpstreamGet(imageUrl + "/action/meta_info")
				fmt.Println("upstreamResponse:", string(upstreamResponse))
				if err == nil {
					var upstreamResponseFields = make(map[string]interface{})
					json.Unmarshal(upstreamResponse, &upstreamResponseFields)
					var fileMetaInfo = make(map[string]interface{})
					fileMetaInfo["name"] = upstreamResponseFields["name"]
					fileMetaInfo["size"] = upstreamResponseFields["size"]
					fileMetaInfo["url"] = upstreamResponseFields["url"]
					recordInfo.Data = fileMetaInfo
				}
				recordInfo.Value = imageUrl
			} else if individualRecordSchema.Type == "object_array" {
				// we should have record schema configured here.
				responseInternal := make(map[string]interface{}, 0)
				var insideObject = make([]map[string]interface{}, 0)
				serializedObject, _ := json.Marshal(recordMap[individualRecordSchema.Property])
				json.Unmarshal(serializedObject, &insideObject)
				var listOfResponse = make([]interface{}, 0)

				for _, individualObject := range insideObject {
					individualResponse := cm.BuildGeneralRecordResponse(dbConnection, individualObject, responseInternal, individualRecordSchema.RecordSchema)
					var finalResponse = make(map[string]interface{})
					for key, value := range individualResponse {
						recordInfoValue := value.(component.RecordInfo)
						if recordInfoValue.Data != nil {
							finalResponse[key] = recordInfoValue.Data
						} else {
							finalResponse[key] = recordInfoValue.Value
						}
					}
					listOfResponse = append(listOfResponse, finalResponse)
				}

				recordInfo.Data = listOfResponse
				var headerObject = make([]interface{}, 0)
				for _, internalRecordSchema := range individualRecordSchema.RecordSchema {
					tableSchema := component.TableSchema{}
					tableSchema.Property = internalRecordSchema.Property
					tableSchema.Name = internalRecordSchema.Name
					tableSchema.Type = internalRecordSchema.Type
					tableSchema.InterfaceType = internalRecordSchema.InterfaceType
					tableSchema.Render = internalRecordSchema.Render
					tableSchema.Display = internalRecordSchema.Display
					tableSchema.LinkedProperty = internalRecordSchema.LinkedProperty
					tableSchema.LinkedDataType = internalRecordSchema.LinkedDataType
					tableSchema.GridSystem = internalRecordSchema.GridSystem
					tableSchema.Label = internalRecordSchema.Label

					if internalRecordSchema.HeaderObjectMapping != nil {
						if internalRecordSchema.HeaderObjectMapping.Predefined != nil {
							headerObjectMapping := internalRecordSchema.HeaderObjectMapping
							predefinedResults := headerObjectMapping.Predefined.Data
							var objectKeyValues = make(map[string]string, len(predefinedResults))
							if headerObjectMapping.Builder.KeyValue != nil {
								for _, result := range predefinedResults {
									objectKeyValues[result[headerObjectMapping.Builder.KeyValue.Key].(string)] = result[headerObjectMapping.Builder.KeyValue.Value].(string)
								}
							}

							raw, _ := json.Marshal(objectKeyValues)
							tableSchema.ObjectList = raw
							// set the object to null, so that response won't send that field
							tableSchema.HeaderObjectMapping = nil
						}
					}
					if internalRecordSchema.HeaderFontColorObjectMapping != nil {
						if internalRecordSchema.HeaderFontColorObjectMapping != nil {
							headerFontColorObjectMapping := internalRecordSchema.HeaderFontColorObjectMapping
							if headerFontColorObjectMapping.Predefined != nil {
								predefinedResults := headerFontColorObjectMapping.Predefined.Data
								var objectKeyValues = make(map[string]string, len(predefinedResults))
								if headerFontColorObjectMapping.Builder.KeyValue != nil {
									for _, result := range predefinedResults {
										objectKeyValues[result[headerFontColorObjectMapping.Builder.KeyValue.Key].(string)] = result[headerFontColorObjectMapping.Builder.KeyValue.Value].(string)
									}
								}

								raw, _ := json.Marshal(objectKeyValues)
								tableSchema.FontColorList = raw
								// set the object to null, so that response won't send that field
								tableSchema.HeaderFontColorObjectMapping = nil
							}

						}
					}

					headerObject = append(headerObject, tableSchema)
				}
				recordInfo.Header = headerObject

			} else if individualRecordSchema.Type == "object" {
				// we should have record schema configured here.
				responseInternal := make(map[string]interface{}, 0)
				var insideObject = make(map[string]interface{}, 0)
				serializedObject, _ := json.Marshal(recordMap[individualRecordSchema.Property])
				json.Unmarshal(serializedObject, &insideObject)
				newResponse := cm.BuildGeneralRecordResponse(dbConnection, insideObject, responseInternal, individualRecordSchema.RecordSchema)
				if _, ok := recordMap[individualRecordSchema.Property]; ok {
					recordInfo.Data = recordMap[individualRecordSchema.Property]
				} else {
					// we don't have field, check any default values configured.
					if individualRecordSchema.Default != nil {
						recordInfo.Data = *individualRecordSchema.Default
					} else {
						recordInfo.Data = recordMap[individualRecordSchema.Property]
					}
				}

				recordInfo.Header = newResponse

			} else {
				recordInfo.Value = recordMap[individualRecordSchema.Property]
			}

		}

		// indicating what is the id it is referring to
		recordInfo.IsEdit = individualRecordSchema.IsEdit
		recordInfo.Type = individualRecordSchema.Type
		recordInfo.Label = individualRecordSchema.Label
		recordInfo.InterfaceType = individualRecordSchema.InterfaceType
		recordInfo.Icon = individualRecordSchema.Icon
		recordInfo.IconType = individualRecordSchema.IconType
		response[individualRecordSchema.Property] = recordInfo

	}
	return response
}

func (cm *ComponentManager) GetIndividualRecordResponse(zone string, dbConnection *gorm.DB, id int, componentName string, rawObjectInfo datatypes.JSON) map[string]interface{} {
	var recordMap = make(map[string]interface{})
	json.Unmarshal(rawObjectInfo, &recordMap)

	// if the size of the map is 0, then no response is decoded, so we can send the empty respose,
	// otherwise send the general record response
	if len(recordMap) == 0 {
		var response = make(map[string]interface{}, 0)
		return response
	} else {
		response := make(map[string]interface{}, 0)
		int64ComponentId := cm.ComponentNameIdMapping[componentName]
		recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
		generalResponse := cm.BuildGeneralRecordResponse(dbConnection, recordMap, response, recordSchema)
		if cm.ComponentSchema[int64ComponentId].EnabledAfterEntityLevel != nil {
			var enableAfterEntityLevel = *cm.ComponentSchema[int64ComponentId].EnabledAfterEntityLevel
			var objectFields = make(map[string]interface{})

			if cm.ComponentSchema[int64ComponentId].EntityField != nil {
				var entityField = *cm.ComponentSchema[int64ComponentId].EntityField
				json.Unmarshal(rawObjectInfo, &objectFields)
				if value, ok := objectFields[entityField]; ok {
					var existingValue = util.InterfaceToInt(value)
					if existingValue > enableAfterEntityLevel {
						if canEditValue, ok := generalResponse["canEdit"]; ok {
							canEditRecordInfo := canEditValue.(component.RecordInfo)
							canEditRecordInfo.Value = true
							generalResponse["canEdit"] = canEditRecordInfo
						}
					}
				}
			}

		}

		cm.GenerateAdditionalRecordResponse(zone, dbConnection, rawObjectInfo, componentName, generalResponse)

		cm.GenerateTemplateResponse(dbConnection, rawObjectInfo, componentName, generalResponse)

		// get the next and previous records
		targetTable := cm.GetTargetTable(componentName)
		recordId := strconv.Itoa(id)
		var nextQueryResults []map[string]interface{}
		nextIdQuery := "SELECT id FROM " + targetTable + " WHERE id > " + recordId + " ORDER BY id LIMIT 1"

		dbConnection.Raw(nextIdQuery).Scan(&nextQueryResults)
		if len(nextQueryResults) == 0 {
			intRecordId, _ := strconv.Atoi(recordId)
			generalResponse["nextId"] = intRecordId
		} else {
			nextId := util.InterfaceToInt(nextQueryResults[0]["id"])
			generalResponse["nextId"] = nextId
		}
		var previousQueryResults []map[string]interface{}
		previousQuery := "SELECT id FROM " + targetTable + "  WHERE id < " + recordId + " ORDER BY id DESC LIMIT 1"
		dbConnection.Raw(previousQuery).Scan(&previousQueryResults)
		if len(previousQueryResults) == 0 {
			intRecordId, _ := strconv.Atoi(recordId)
			generalResponse["previousId"] = intRecordId
		} else {
			previousId := util.InterfaceToInt(previousQueryResults[0]["id"])
			generalResponse["previousId"] = previousId
		}
		var totalRecordsQueryResults []map[string]interface{}
		totalRecordsQuery := "SELECT COUNT(*) AS total_records FROM  " + targetTable

		dbConnection.Raw(totalRecordsQuery).Scan(&totalRecordsQueryResults)
		if len(totalRecordsQueryResults) == 0 {
			generalResponse["totalRecords"] = 0
		} else {
			totalRecords := util.InterfaceToInt(totalRecordsQueryResults[0]["total_records"])
			generalResponse["totalRecords"] = totalRecords
		}
		return generalResponse
	}
}
