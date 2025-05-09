package handler

import (
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/error_util"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"cx-micro-flake/services/batch_management/handler/const_util"
	"cx-micro-flake/services/batch_management/handler/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func (v *BatchManagementService) loadFile(ctx *gin.Context) {
	loadDataFileCommand := component.LoadDataFileCommand{}
	if err := ctx.ShouldBindBodyWith(&loadDataFileCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err, errorCode, loadFileResponse := v.ComponentManager.ProcessLoadFile(loadDataFileCommand)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	ctx.JSON(http.StatusOK, loadFileResponse)
	return
}

func (v *BatchManagementService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")
	//
	//componentName := ctx.Param("componentName")

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")

}

func (v *BatchManagementService) exportObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")
	//componentName := ctx.Param("componentName")

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")

}

func (v *BatchManagementService) getTableImportSchema(ctx *gin.Context) {
	//componentName := ctx.Param("componentName")

	response.SendDetailedError(ctx, http.StatusBadRequest, errors.New("Invalid Component"), const_util.InvalidMouldComponent, "Requested component name doesn't exist or function is not supported yet")

}

func (v *BatchManagementService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)

}

func (v *BatchManagementService) getObjects(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")
	fields := ctx.Query("fields")
	values := ctx.Query("values")
	condition := ctx.Query("condition")
	outFields := ctx.Query("out_fields")
	format := ctx.Query("format")
	searchFields := ctx.Query("search")
	orderValue := ctx.Query("order")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	//Have to next flag
	isNext := true
	var listOfObjects []component.GeneralObject
	var totalRecords int64
	var err error
	userId := common.GetUserId(ctx)
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
		err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, searchWithBaseQuery, limitVal)
	} else if offsetValue == "" && limitValue == "" && fields == "" && values == "" && condition == "" && outFields == "" {
		err, listOfObjects = database.GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(listOfObjects))
	} else {

		totalRecords = database.Count(dbConnection, targetTable)
		if limitValue == "" {
			err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			var conditionString string
			limitVal, _ := strconv.Atoi(limitValue)

			//currentRecordCount := len(*listOfObjects)

			if orderValue == "desc" {
				offsetVal, _ := strconv.Atoi(offsetValue)
				var tableCondition string
				if conditionString != "" {
					if offsetVal == -1 {
						tableCondition = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						tableCondition = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}
					if tableCondition != "" {
						conditionString = tableCondition + " AND " + conditionString
					}
				} else {
					if offsetVal == -1 {
						conditionString = component.TableConditionV1(offsetValue, fields, values, condition)
					} else {
						conditionString = component.TableDecendingOrderCondition(offsetValue, fields, values, condition)
					}

				}

				orderBy := "object_info ->> '$.createdAt' desc"

				err, listOfObjects = database.GetConditionalObjectsOrderBy(dbConnection, targetTable, conditionString, orderBy, limitVal)

				currentRecordCount := len(listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(conditionString, "AND")
					var totalRecordObjects []component.GeneralObject
					if len(andClauses) > 1 {
						_, totalRecordObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

					} else {
						_, totalRecordObjects = database.GetObjects(dbConnection, targetTable)
					}

					if (listOfObjects)[currentRecordCount-1].Id == (totalRecordObjects)[0].Id {
						isNext = false
					}
				}

			} else {
				err, listOfObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString, limitVal)
				currentRecordCount := len(listOfObjects)

				if currentRecordCount < limitVal {
					isNext = false
				} else if currentRecordCount == limitVal {
					andClauses := strings.Split(conditionString, "AND")
					var totalRecordObjects []component.GeneralObject
					if len(andClauses) > 1 {
						_, totalRecordObjects = database.GetConditionalObjects(dbConnection, targetTable, conditionString)

					} else {
						err, totalRecordObjects = database.GetObjects(dbConnection, targetTable)
					}
					lenTotalRecord := len(totalRecordObjects)
					if (listOfObjects)[currentRecordCount-1].Id == (totalRecordObjects)[lenTotalRecord-1].Id {
						isNext = false
					}
				}
			}
		}

	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArrayV1(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		userId = common.GetUserId(ctx)
		zone := getUserTimezone(userId)
		_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

func (v *BatchManagementService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects []component.GeneralObject
	var err error
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if offsetValue != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		err, listOfObjects = database.GetObjects(dbConnection, targetTable, limitVal)
	} else {
		err, listOfObjects = database.GetObjects(dbConnection, targetTable)
	}

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)
}

func (v *BatchManagementService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, recordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		error_util.SendResourceNotFound(ctx)
		return
	}
	// first archive all the dependency
	userId := common.GetUserId(ctx)
	v.archiveReferences(userId, dbConnection, componentName, recordId)

	err = database.ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		error_util.SendArchiveFailed(ctx)
		return
	}

	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource is archived, no further modification allowed", recordId, userId, nil, nil)
	ctx.Status(http.StatusNoContent)

}

func (v *BatchManagementService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var userId = common.GetUserId(ctx)
	var removeInternalReferenceRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&removeInternalReferenceRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}
	initializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = initializedObject
	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), const_util.ErrorUpdatingObjectInformation)
		return
	}
	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "Resource get updated", intRecordId, userId, &objectInterface, &component.GeneralObject{ObjectInfo: initializedObject})

	ctx.JSON(http.StatusOK, updatingData)

}

func (v *BatchManagementService) updateResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	err, objectInterface := database.Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	if !common.ValidateObjectStatus(objectInterface.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      const_util.GetError(common.InvalidObjectStatusError).Error(),
				Description: "This resource is already archived, no further modifications are allowed.",
			})
		return
	}

	var updateRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})

	//Adding update process request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)
	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		v.BaseService.Logger.Error("error updating resource", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating resource information"), const_util.ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		Code:    0,
		Message: "Successfully updated the resource",
	})

}

func (v *BatchManagementService) createNewResource(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, createRequest := component.GetRequestFields(ctx)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Validation Failed"), const_util.FieldValidationFailed, err.Error())
		return
	}
	processedDefaultValues := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	processedDefaultObjectMappingFields := v.ComponentManager.ProcessDefaultObjectMapping(dbConnection, processedDefaultValues, componentName)

	processedDefaultObjectMappingFields["objectStatus"] = common.ObjectStatusActive
	rawCreateRequest, _ := json.Marshal(processedDefaultObjectMappingFields)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)

	var createdRecordId int
	var userId = common.GetUserId(ctx)

	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}

	err, createdRecordId = database.CreateFromGeneralObject(dbConnection, targetTable, object)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Error in Resource Creation "), const_util.ErrorCreatingObjectInformation, err.Error())
		return
	}

	if targetTable == const_util.BatchManagementRawMaterialTable {
		// check the printer is defined, if not throw an error saying printer is not configured for that location

		sidUpdateRequest := make(map[string]interface{})
		json.Unmarshal(preprocessedRequest, &sidUpdateRequest)
		var locationId int
		if locationInterface, ok := sidUpdateRequest["location"]; ok {
			factoryInterface := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
			isBuilding := factoryInterface.IsFactoryBuildingExist(util.InterfaceToInt(locationInterface))
			if !isBuilding {
				response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Invalid Location ID"), const_util.ErrorCreatingObjectInformation, "Batch management need location defined, please make sure correction location is defined")
				database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
				return
			}
			locationId = util.InterfaceToInt(locationInterface)
		} else {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Invalid Location"), const_util.ErrorCreatingObjectInformation, "Batch management need location defined, please make sure correction location is defined")
			database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
			return
		}

		if createdRecordId < 10 {
			sidUpdateRequest["materialBatchNo"] = "MAT0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			sidUpdateRequest["materialBatchNo"] = "MAT000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			sidUpdateRequest["materialBatchNo"] = "MAT00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			sidUpdateRequest["materialBatchNo"] = "MAT0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			sidUpdateRequest["materialBatchNo"] = "MAT" + strconv.Itoa(createdRecordId)
		}
		var batchId = ""
		if util.InterfaceToInt(sidUpdateRequest["materialType"]) == const_util.RasinType {

			if util.InterfaceToString(sidUpdateRequest["vendorLotNo"]) == "" {
				v.BaseService.Logger.Error("error creating label", zap.String("error", err.Error()))
				response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal Error"), const_util.ErrorCreatingObjectInformation, "Please enter Vendor Lot Number")
				database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
				return
			}
			batchId = util.InterfaceToString(sidUpdateRequest["vendorLotNo"])
		} else {
			_, generalObject := database.Get(dbConnection, const_util.ManufacturingVendorMasterTable, util.InterfaceToInt(sidUpdateRequest["vendorId"]))
			vendorMasterInfo := make(map[string]interface{})
			json.Unmarshal(generalObject.ObjectInfo, &vendorMasterInfo)
			if util.InterfaceToString(sidUpdateRequest["vendorLotNo"]) != "" {
				batchId = util.InterfaceToString(vendorMasterInfo["name"]) + "_" + util.InterfaceToString(sidUpdateRequest["vendorLotNo"])
			} else {
				zone := getUserTimezone(userId)
				batchId = util.InterfaceToString(vendorMasterInfo["name"]) + "_" + util.GetZoneCurrentTimeInPMFormat(zone)
			}
		}
		sidUpdateRequest["batchId"] = batchId
		// now generate the label
		// get the printer setting
		conditionalString := "object_info ->> '$.location' = " + strconv.Itoa(locationId)
		_, locationInterface := database.GetConditionalObjects(dbConnection, const_util.BatchManagementPrinterTable, conditionalString)

		if len(locationInterface) > 0 {
			printerSetting := database.GetPrinterInfo((locationInterface)[0].ObjectInfo)
			err, generatedQRLabel := generateQRCodeLabel(batchId, printerSetting)
			if err != nil {
				v.BaseService.Logger.Error("error creating label", zap.String("error", err.Error()))
				response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal Error"), const_util.ErrorCreatingObjectInformation, "Label creation failed due to ["+err.Error()+"]")
				database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
				return
			}
			sidUpdateRequest["label"] = generatedQRLabel
			err, generatedQRImage := generateQRCode(batchId)
			if err != nil {
				v.BaseService.Logger.Error("error creating barcode", zap.String("error", err.Error()))
				response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal Error"), const_util.ErrorCreatingObjectInformation, "QR code creation failed due to ["+err.Error()+"]")
				database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
				return
			}
			sidUpdateRequest["labelImage"] = generatedQRImage
		}

		v.BaseService.Logger.Info("update request for raw material", zap.Any("update_request", sidUpdateRequest))

		preprocessedRequest, _ = json.Marshal(sidUpdateRequest)
		//err = orm.UpdateSerialisedResourceFromId(dbConnection, InventorySupplierTable, createdRecordId, userId, rawWorkOrderInfo)
	} else if targetTable == const_util.BatchManagementMouldTable {

		sidUpdateRequest := make(map[string]interface{})
		json.Unmarshal(preprocessedRequest, &sidUpdateRequest)

		if locationInterface, ok := sidUpdateRequest["location"]; ok {
			factoryInterface := common.GetService("factory_module").ServiceInterface.(common.FactoryServiceInterface)
			isBuilding := factoryInterface.IsFactoryBuildingExist(util.InterfaceToInt(locationInterface))
			if !isBuilding {
				fmt.Println("error", isBuilding)
				response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Invalid Location ID"), const_util.ErrorCreatingObjectInformation, "Batch management need location defined, please make sure correction location is defined")
				database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
				return
			}

		} else {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Invalid Location"), const_util.ErrorCreatingObjectInformation, "Batch management need location defined, please make sure correction location is defined")
			database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
			return
		}
		var mouldBatchNumber = ""
		if createdRecordId < 10 {
			mouldBatchNumber = "MOL0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			mouldBatchNumber = "MOL000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			mouldBatchNumber = "MOL00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			mouldBatchNumber = "MOL0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			mouldBatchNumber = "MOL" + strconv.Itoa(createdRecordId)
		}
		sidUpdateRequest["mouldBatchNumber"] = mouldBatchNumber

		var rawMaterialId = util.InterfaceToInt(sidUpdateRequest["rawMaterialId"])
		// take the mould batch reference ID
		var mouldBatchId = ""

		// take it at least one
		err, listOfBatches := database.GetConditionalObjects(dbConnection, const_util.BatchManagementMouldTable, " object_info->>'$.rawMaterialId' ="+strconv.Itoa(rawMaterialId))
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal System Error"), const_util.ErrorCreatingObjectInformation, "System couldn't able to determine the raw material batch, please check the raw material")
			database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
			return
		}
		if len(listOfBatches) == 1 {
			var newMachineId string
			var scheduleOrderName string

			machineId := util.InterfaceToInt(sidUpdateRequest["machineId"])
			scheduleEventId := util.InterfaceToInt(sidUpdateRequest["scheduleEventId"])

			machineService := common.GetService("machines_module").ServiceInterface.(common.MachineInterface)
			machineError, machineGeneralComponent := machineService.GetMachineInfoById(projectId, machineId)

			if machineError == nil {
				machineInfo := make(map[string]interface{})
				json.Unmarshal(machineGeneralComponent.ObjectInfo, &machineInfo)

				newMachineId = util.InterfaceToString(machineInfo["newMachineId"])
			}

			productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
			scheduleError, scheduleOrderInfo := productionOrderInterface.GetScheduledOrderInfo(projectId, scheduleEventId)

			if scheduleError == nil {
				machineInfo := make(map[string]interface{})
				json.Unmarshal(scheduleOrderInfo.ObjectInfo, &machineInfo)

				scheduleOrderName = util.InterfaceToString(machineInfo["name"])
			}

			if createdRecordId < 10 {
				mouldBatchId = newMachineId + "_" + scheduleOrderName + "_" + "0000" + strconv.Itoa(createdRecordId)
			} else if createdRecordId < 100 {
				mouldBatchId = newMachineId + "_" + scheduleOrderName + "_" + "000" + strconv.Itoa(createdRecordId)
			} else if createdRecordId < 1000 {
				mouldBatchId = newMachineId + "_" + scheduleOrderName + "_" + "00" + strconv.Itoa(createdRecordId)
			} else if createdRecordId < 10000 {
				mouldBatchId = newMachineId + "_" + scheduleOrderName + "_" + "0" + strconv.Itoa(createdRecordId)
			} else if createdRecordId < 100000 {
				mouldBatchId = newMachineId + "_" + scheduleOrderName + "_" + strconv.Itoa(createdRecordId)
			}
			//now := time.Now()
			//// Format the time to a string in the desired format
			//mouldBatchId = now.Format("20060102_150405")
			sidUpdateRequest["mouldBatchId"] = mouldBatchId
		} else if len(listOfBatches) > 1 {
			var mouldBatchFields = make(map[string]interface{})
			json.Unmarshal((listOfBatches)[0].ObjectInfo, &mouldBatchFields)
			mouldBatchId = util.InterfaceToString(mouldBatchFields["mouldBatchId"])
			sidUpdateRequest["mouldBatchId"] = mouldBatchId

			locationId := util.InterfaceToInt(sidUpdateRequest["location"])
			conditionalString := "object_info ->> '$.location' = " + strconv.Itoa(locationId)
			_, locationInterface := database.GetConditionalObjects(dbConnection, const_util.BatchManagementPrinterTable, conditionalString)

			if len(locationInterface) > 0 {
				printerSetting := database.GetPrinterInfo((locationInterface)[0].ObjectInfo)
				err, generatedQRLabel := generateQRCodeLabel(mouldBatchId, printerSetting)
				if err != nil {
					v.BaseService.Logger.Error("error creating label", zap.String("error", err.Error()))
					response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal Error"), const_util.ErrorCreatingObjectInformation, "Label creation failed due to ["+err.Error()+"]")
					database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
					return
				}
				sidUpdateRequest["label"] = generatedQRLabel
				err, generatedQRImage := generateQRCode(mouldBatchId)
				if err != nil {
					v.BaseService.Logger.Error("error creating barcode", zap.String("error", err.Error()))
					response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Internal Error"), const_util.ErrorCreatingObjectInformation, "QR code creation failed due to ["+err.Error()+"]")
					database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
					return
				}
				sidUpdateRequest["labelImage"] = generatedQRImage
			}
		}
		// we don't generate the label if the production order continue, it is only created the label if it is completed or stopped only
		preprocessedRequest, _ = json.Marshal(sidUpdateRequest)

		// now update the mould scheduled order event table with mould batch ID
		productionOrderInterface := common.GetService("production_order_module").ServiceInterface.(common.ProductionOrderInterface)
		scheduleEventId := sidUpdateRequest["scheduleEventId"]
		var scheduleEventResourceId = util.InterfaceToInt(scheduleEventId)
		err = productionOrderInterface.UpdateMouldScheduleOrderEventMouldBatchId(projectId, scheduleEventResourceId, createdRecordId)
		if err != nil {
			response.SendDetailedError(ctx, http.StatusBadRequest, const_util.GetError("Invalid Schedule Event"), const_util.ErrorCreatingObjectInformation, "Batch management need location defined, please make sure correction scheduler event is defined")
			database.Delete(dbConnection, const_util.BatchManagementRawMaterialTable, createdRecordId)
			return
		}

	}
	updatingData := make(map[string]interface{})
	updatingData["object_info"] = preprocessedRequest

	err = database.Update(v.BaseService.ReferenceDatabase, targetTable, createdRecordId, updatingData)
	v.BaseService.Logger.Info("object is created successfully", zap.Any("record_id", createdRecordId), zap.Any("component_name", componentName))

	v.CreateUserRecordMessage(const_util.ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:     0,
		RecordId: createdRecordId,
		Message:  "New resource is successfully created",
	})
}

func (v *BatchManagementService) getNewRecord(ctx *gin.Context) {
	componentName := util.GetComponentName(ctx)
	projectId := util.GetProjectId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

func (v *BatchManagementService) getRecordFormData(ctx *gin.Context) {
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting table information", zap.String("error", err.Error()), zap.String("table", targetTable))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), const_util.ErrorGettingIndividualObjectInformation)
		return
	}

	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	ctx.JSON(http.StatusOK, response)

}

func (v *BatchManagementService) getSearchResults(ctx *gin.Context) {
	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)

	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)

	if len(searchFieldCommand) == 0 {
		// reset the search
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}

	format := ctx.Query("format")
	searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
	err, listOfObjects := database.GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), const_util.ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponseV1(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), const_util.ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, int64(len(listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)

}

func (v *BatchManagementService) deleteValidation(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := util.GetRecordId(ctx)
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := database.Get(dbConnection, targetTable, recordId)

	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Invalid Resource",
				Description: "The resource that you are trying to delete doesn't exist, Please check refresh page and try again",
			})
		return
	}
	if component.IsArchived(generalObject.ObjectInfo) {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Resource Archived",
				Description: "The resource that you are trying to delete is already archived. This operation is not allowed",
			})
		return
	}
	var dependencyComponents []string
	var dependencyRecords int
	v.checkReference(dbConnection, componentName, recordId, &dependencyComponents, &dependencyRecords)
	if dependencyRecords > 0 {
		var dependencyString string
		dependencyComponents = util.RemoveDuplicateString(dependencyComponents)
		dependencyString = " ["
		for index, dependencyComponent := range dependencyComponents {
			if index == len(dependencyComponents)-1 {
				dependencyString += dependencyComponent
			} else {
				dependencyString += dependencyComponent + " ->"
			}
		}
		dependencyString += " ]"
		ctx.JSON(http.StatusOK, response.GeneralResponse{
			CanDelete: false,
			Code:      100,
			Message:   "There are dependencies bound to the resource that you are trying to remove. Removing this resource would create the chain removal on following resources " + dependencyString + " in " + strconv.Itoa(dependencyRecords) + " places, Please understand the risk of deleting this resource as all the dependencies would be achieved immediately, and this process is not reversible",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GeneralResponse{
		CanDelete: true,
		Code:      100,
		Message:   "There are no dependencies bound to the resource that you are trying to remove. So, removing this resource won't affect others resource now, you can proceed !!",
	})

}

func (v *BatchManagementService) getGroupBy(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")

	groupByAction := component.GroupByAction{}
	if err := ctx.ShouldBindBodyWith(&groupByAction, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// if the group by is empty, then return the normal default 30 records per page results.
	if len(groupByAction.GroupBy) == 0 {
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getObjects(ctx)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, listOfObjects := database.GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	fmt.Println("groupByAction : ", groupByAction)
	var totalRecords = database.Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
	tableResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableResponse)
	groupByColumns := groupByAction.GroupBy
	finalResponse := component.TableObjectResponse{}
	if len(groupByColumns) == 1 {
		results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range results {
			groupByChildren := component.GroupByChildren{}

			groupByChildren.Data = level1Value
			groupByChildren.Type = "json"

			tableGroupResponse := component.TableGroupByResponse{}
			tableGroupResponse.Label = level1Key
			tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Data = append(finalResponse.Data, rawData)
		}
	} else if len(groupByColumns) == 2 {
		level1Results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range level1Results {
			tableGroupResponse := component.TableGroupByResponse{}
			if len(level1Value) > 1 {
				var internalGroupResponse []interface{}
				// here we need to group again

				level2Results := component.GetGroupByResultsFromInterface(groupByColumns[1], level1Value)
				for level2Key, level2Value := range level2Results {
					level2Children := component.GroupByChildren{}
					level2Children.Data = level2Value
					level2Children.Type = "json"

					internalTableGroupResponse := component.TableGroupByResponse{}
					internalTableGroupResponse.Label = level2Key
					internalTableGroupResponse.Children = append(internalTableGroupResponse.Children, level2Children)
					internalGroupResponse = append(internalGroupResponse, internalTableGroupResponse)
				}
				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = internalGroupResponse
			} else {
				groupByChildren := component.GroupByChildren{}
				groupByChildren.Data = level1Value
				groupByChildren.Type = "json"

				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			}

			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Data = append(finalResponse.Data, rawData)
		}
	}

	finalResponse.Header = tableResponse.Header
	finalResponse.TotalRowCount = tableResponse.TotalRowCount
	finalResponse.CurrentRowCount = tableResponse.CurrentRowCount

	ctx.JSON(http.StatusOK, finalResponse)
}
func (v *BatchManagementService) getCardViewGroupBy(ctx *gin.Context) {

	//userId := common.GetUserId(ctx)
	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")

	groupByAction := component.GroupByAction{}
	if err := ctx.ShouldBindBodyWith(&groupByAction, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// if the group by is empty, then return the normal default 30 records per page results.
	if len(groupByAction.GroupBy) == 0 {
		ctx.Set("offset", 1)
		ctx.Set("limit", 30)
		v.getCardView(ctx)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, listOfObjects := database.GetObjects(dbConnection, targetTable)
	if err != nil {
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "System could not able to get all requested resources due to internal system exception. Please report this error to system admin",
			})
		return
	}
	fmt.Println("groupByAction : ", groupByAction)
	var totalRecords = database.Count(dbConnection, targetTable)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	_, tableRecordsResponse := v.ComponentManager.GetTableRecordsV1(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
	tableResponse := component.TableObjectResponse{}
	json.Unmarshal(tableRecordsResponse, &tableResponse)
	groupByColumns := groupByAction.GroupBy
	finalResponse := component.CardViewGroupResponse{}
	if len(groupByColumns) == 1 {
		results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range results {
			groupByChildren := component.GroupByChildren{}

			groupByChildren.Data = v.ComponentManager.GetCardViewFromListOfInterface(level1Value, componentName)
			groupByChildren.Type = "json"

			tableGroupResponse := component.TableGroupByResponse{}
			tableGroupResponse.Label = level1Key
			tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Cards = append(finalResponse.Cards, rawData)
		}
	} else if len(groupByColumns) == 2 {
		level1Results := component.GetGroupByResults(groupByColumns[0], tableResponse.Data)
		for level1Key, level1Value := range level1Results {
			tableGroupResponse := component.TableGroupByResponse{}
			if len(level1Value) > 1 {
				var internalGroupResponse []interface{}
				// here we need to group again

				level2Results := component.GetGroupByResultsFromInterface(groupByColumns[1], level1Value)
				for level2Key, level2Value := range level2Results {
					level2Children := component.GroupByChildren{}
					level2Children.Data = v.ComponentManager.GetCardViewFromListOfInterface(level2Value, componentName)
					level2Children.Type = "json"

					internalTableGroupResponse := component.TableGroupByResponse{}
					internalTableGroupResponse.Label = level2Key
					internalTableGroupResponse.Children = append(internalTableGroupResponse.Children, level2Children)
					internalGroupResponse = append(internalGroupResponse, internalTableGroupResponse)
				}
				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = internalGroupResponse
			} else {
				groupByChildren := component.GroupByChildren{}
				groupByChildren.Data = level1Value
				groupByChildren.Type = "json"

				tableGroupResponse.Label = level1Key
				tableGroupResponse.Children = append(tableGroupResponse.Children, groupByChildren)
			}

			rawData, _ := json.Marshal(tableGroupResponse)
			finalResponse.Cards = append(finalResponse.Cards, rawData)
		}
	}

	ctx.JSON(http.StatusOK, finalResponse)
}

func (v *BatchManagementService) checkReference(dbConnection *gorm.DB, componentName string, recordId int, dependencyComponents *[]string, dependencyRecords *int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					*dependencyComponents = append(*dependencyComponents, constraint.ReferenceComponentDisplayName)
					*dependencyRecords = *dependencyRecords + len(listOfObjects)
					for _, referenceObject := range listOfObjects {
						v.checkReference(dbConnection, referenceComponent, referenceObject.Id, dependencyComponents, dependencyRecords)
					}
				}
			}
		}
	}
}

func (v *BatchManagementService) archiveReferences(userId int, dbConnection *gorm.DB, componentName string, recordId int) {
	listOfConstraints := v.ComponentManager.GetConstraints(componentName)
	if len(listOfConstraints) > 0 {
		for _, constraint := range listOfConstraints {
			referenceComponent := constraint.Reference
			referenceField := constraint.ReferenceProperty
			referenceTable := v.ComponentManager.GetTargetTable(referenceComponent)
			conditionString := " object_info ->>'$." + referenceField + "'=" + strconv.Itoa(recordId)
			err, listOfObjects := database.GetConditionalObjects(dbConnection, referenceTable, conditionString)
			if err == nil {
				if len(listOfObjects) > 0 {
					for _, referenceObject := range listOfObjects {
						err := database.ArchiveObject(dbConnection, referenceTable, referenceObject)
						if err == nil {
							err := v.CreateUserRecordMessage(const_util.ProjectID, referenceComponent, "Resource is deleted", referenceObject.Id, userId, nil, nil)
							if err == nil {
								v.archiveReferences(userId, dbConnection, referenceComponent, referenceObject.Id)
							}
						}
					}
				}
			}
		}
	}
}
