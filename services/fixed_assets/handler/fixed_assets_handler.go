package handler

import (
	"cx-micro-flake/pkg/auth"
	"cx-micro-flake/pkg/common"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// loadFile ShowAccount godoc
// @Summary load the file and get the schema information with data(currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.LoadDataFileCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/loadFile [post]
func (v *FixedAssetService) loadFile(ctx *gin.Context) {
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

// importObjects ShowAccount godoc
// @Summary import machine register information (currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.ImportDataCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/import [get]
func (v *FixedAssetService) importObjects(ctx *gin.Context) {
	// we will get the uploaded url
	//projectId := ctx.Param("projectId")

	//componentName := ctx.Param("componentName")
	//targetTable := ms.ComponentManager.GetTargetTable(componentName)
	importDataCommand := component.ImportDataCommand{}
	if err := ctx.ShouldBindBodyWith(&importDataCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

}

// exportObjects ShowAccount godoc
// @Summary export machine related information (currently only csv format)
// @Description based on user permission, user will allow importing csv file url to populate machine register
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param User body  component.ExportDataCommand true "Send the following fields"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/import [get]
func (v *FixedAssetService) exportObjects(ctx *gin.Context) {
	// we will get the uploaded url
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	exportCommand := component.ExportDataCommand{}

	if err := ctx.ShouldBindBodyWith(&exportCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var condition string
	err, errorCode, exportDataResponse := v.ComponentManager.ExportData(dbConnection, componentName, exportCommand, condition)
	if err != nil {
		v.BaseService.Logger.Error("unable to handle export", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, err, errorCode)
		return
	}
	ctx.JSON(http.StatusOK, exportDataResponse)
}

// getTableSchema ShowAccount godoc
// @Summary Get the table schema
// @Description based on user permission, user will get the table related fields
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/table_import_schema [get]
func (v *FixedAssetService) getTableImportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	tableImportSchema := v.ComponentManager.GetTableImportSchema(componentName)
	ctx.JSON(http.StatusOK, tableImportSchema)
}

// getExportSchema ShowAccount godoc
// @Summary Get the table schema
// @Description based on user permission, user will get the table related fields
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/getExportSchema [get]
func (v *FixedAssetService) getExportSchema(ctx *gin.Context) {
	componentName := ctx.Param("componentName")
	exportSchema := v.ComponentManager.GetTableExportSchema(componentName)
	ctx.JSON(http.StatusOK, exportSchema)
}

// getMachineRegister ShowAccount godoc
// @Summary Get all the machine related information
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/records [get]
func (v *FixedAssetService) getObjects(ctx *gin.Context) {

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
	//Have to next flag
	isNext := true
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var totalRecords int64
	var err error
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
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
		fmt.Println("target table :", targetTable)
		listOfObjects, err = GetObjects(dbConnection, targetTable)
		fmt.Println("target table :", listOfObjects)
		if listOfObjects != nil {
			totalRecords = int64(len(*listOfObjects))
		}

	} else {
		totalRecords = Count(dbConnection, targetTable)
		if limitValue == "" {
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition))

		} else {
			limitVal, _ := strconv.Atoi(limitValue)
			listOfObjects, err = GetConditionalObjects(dbConnection, targetTable, component.TableCondition(offsetValue, fields, values, condition), limitVal)

			currentRecordCount := len(*listOfObjects)
			conditionString := component.TableCondition(offsetValue, fields, values, condition)
			if currentRecordCount < limitVal {
				isNext = false
			} else if currentRecordCount == limitVal {
				andClauses := strings.Split(conditionString, "AND")
				var totalRecordObjects *[]component.GeneralObject
				if len(andClauses) > 1 {
					totalRecordObjects, _ = GetConditionalObjects(dbConnection, targetTable, conditionString)

				} else {
					totalRecordObjects, _ = GetObjects(dbConnection, targetTable)
				}
				lenTotalRecord := len(*totalRecordObjects)
				if (*listOfObjects)[currentRecordCount-1].Id == (*totalRecordObjects)[lenTotalRecord-1].Id {
					isNext = false
				}
			}
		}

	}
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if format == "array" {
		arrayResponseError, arrayResponse := v.ComponentManager.TableRecordsToArray(dbConnection, listOfObjects, componentName, outFields)
		if arrayResponseError != nil {
			v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, arrayResponse)

	} else {
		_, tableRecordsResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, outFields, zone)
		tableObjectResponse := component.TableObjectResponse{}
		json.Unmarshal(tableRecordsResponse, &tableObjectResponse)

		tableObjectResponse.IsNext = isNext
		tableRecordsResponse, _ = json.Marshal(tableObjectResponse)

		ctx.JSON(http.StatusOK, tableRecordsResponse)
	}
}

// getCardView ShowAccount godoc
// @Summary Get all the machine information in a card view
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/card_view [get]
func (v *FixedAssetService) getCardView(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	projectId := ctx.Param("projectId")
	offsetValue := ctx.Query("offset")
	limitValue := ctx.Query("limit")

	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var listOfObjects *[]component.GeneralObject
	var err error
	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))

	if offsetValue != "" && limitValue != "" {
		limitVal, _ := strconv.Atoi(limitValue)

		// requesting to search fields for table
		listOfObjects, err = GetObjects(dbConnection, targetTable, limitVal)
	} else {
		listOfObjects, err = GetObjects(dbConnection, targetTable)
	}

	v.BaseService.Logger.Info("parameter info", zap.Any("project_id", projectId), zap.Any("component_id", componentName), zap.Any("target_table", targetTable), zap.Any("offset_table", offsetValue), zap.Any("limit_value", limitValue))
	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingObjectsInformation)
		return
	}
	cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)

	ctx.JSON(http.StatusOK, cardViewResponse)

}

// deleteResource ShowAccount godoc
// @Summary Delete the any given resource using resource id
// @Description based on user permission, user can perform delete operations
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   offset     query    int     false        "offset"
// @Param   limit     query    int     false        "limit"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [delete]
func (v *FixedAssetService) deleteResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		v.BaseService.Logger.Error("error getting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting records information"), ErrorGettingIndividualObjectInformation)
		return
	}

	err = ArchiveObject(dbConnection, targetTable, generalObject)

	if err != nil {
		v.BaseService.Logger.Error("error deleting records", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error removing records information"), ErrorRemovingObjectInformation)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// updateResource ShowAccount godoc
// @Summary update given resource based on resource id
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   resourceId     path    string     true        "Resource Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [put]
func (v *FixedAssetService) updateResource(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	var updateRequest = make(map[string]interface{})

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}

	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//Adding update preprocess request
	serializedObject := v.ComponentManager.GetUpdateRequest(updateRequest, objectInterface.ObjectInfo, componentName)

	err = v.ComponentManager.DoFieldValidationOnSerializedObject(componentName, "update", dbConnection, serializedObject)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	initializedObject := common.UpdateMetaInfoFromSerializedObject(serializedObject, ctx)
	updatingData["object_info"] = initializedObject

	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		fmt.Println(err)
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating resource information"), ErrorUpdatingObjectInformation)
		return
	}

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully updated",
		Error:   0,
	})

}

// createNewResource ShowAccount godoc
// @Summary create new resource
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   componentId     path    string     true        "Component Id"
// @Param   recordId     path    string     true        "Record Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/records [post]
func (v *FixedAssetService) createNewResource(ctx *gin.Context) {

	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	userId := common.GetUserId(ctx)

	targetTable := v.ComponentManager.GetTargetTable(componentName)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	var createRequest = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&createRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createRequest["objectStatus"] = common.ObjectStatusActive

	var createdRecordId int
	var err error

	updatedRequest := v.ComponentManager.PreprocessCreateRequestFields(createRequest, componentName)
	// here we should do the validation
	err = v.ComponentManager.DoFieldValidation(componentName, "create", dbConnection, createRequest)
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Validation Failed"), ErrorCreatingObjectInformation, err.Error())
		return
	}

	rawCreateRequest, _ := json.Marshal(updatedRequest)
	preprocessedRequest := common.InitMetaInfoFromSerializedObject(rawCreateRequest, ctx)
	object := component.GeneralObject{
		ObjectInfo: preprocessedRequest,
	}
	err, createdRecordId = Create(dbConnection, targetTable, object)
	if err != nil {
		fmt.Println("err : ", err.Error())
		response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
			&response.DetailedError{
				Header:      "Internal Server Error",
				Description: "Error creating new resource in the system due to internal server error",
			})
		return
	}

	switch componentName {
	case FixedAssetContractComponent:

		_, generalContractOrder := Get(dbConnection, FixedAssetContractTable, createdRecordId)
		fixedAssetContractInfo := make(map[string]interface{})
		json.Unmarshal(generalContractOrder.ObjectInfo, &fixedAssetContractInfo)
		AssetNumberPrefix := "CN"
		if createdRecordId < 10 {
			fixedAssetContractInfo["contractNumber"] = AssetNumberPrefix + "0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			fixedAssetContractInfo["contractNumber"] = AssetNumberPrefix + "000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			fixedAssetContractInfo["contractNumber"] = AssetNumberPrefix + "00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			fixedAssetContractInfo["contractNumber"] = AssetNumberPrefix + "0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			fixedAssetContractInfo["contractNumber"] = AssetNumberPrefix + "SR" + strconv.Itoa(createdRecordId)
		}

		// Create the barcode
		resourceToken, _ := auth.CreateResourceToken(projectId, FixedAssetContractComponent, createdRecordId)
		url := v.QRCodeDomainUrl + "?lms_token=" + resourceToken
		qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
		v.BaseService.Logger.Info("created QR token", zap.String("token", resourceToken))

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 200, 200)

		// create the output file
		fileSuffixName := uuid.New().String()
		fileName := v.ComponentContentConfig.DownloadDirectory + "/" + fileSuffixName + ".png"
		file, err := os.Create(fileName)
		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "QR Code Save Failed",
					Description: "Internal system error during QR code storing, this is normally happen when the backend is not properly configured the path",
				})
			return
		}
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)
		var params map[string]string
		err, httpResponse, rawResponse := util.FileUploadFromDisk(v.ComponentContentConfig.UpStream, params, "file", fileName)

		if err == nil && httpResponse.StatusCode == 200 {
			// no error
			fmt.Println("raw response :", string(rawResponse))
			fileUploadResponse := FileUploadResponse{}
			json.Unmarshal(rawResponse, &fileUploadResponse)
			fmt.Println("contentResponse: ", fileUploadResponse)
			fixedAssetContractInfo["qrCodeUrl"] = fileUploadResponse.Url
		} else {
			fmt.Println("raw response :", string(rawResponse))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "QR Code Generation Failed",
					Description: "Internal system error during QR code generation, this is normally happen when the backend is not properly configured the path",
				})
			return
		}

		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(fixedAssetContractInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		Update(dbConnection, FixedAssetContractTable, createdRecordId, updatingData)

		// Email sending part
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		listOfDepartmentHod := authService.GetHeadOfDepartments(userId)

		// Sending email to HODs
		for _, hodAccountId := range listOfDepartmentHod {
			v.emailGenerator(dbConnection, FixedAssetContractCreationEmailTemplate, hodAccountId.UserId, FixedAssetContractComponent, createdRecordId)
		}

		conditionString := " object_info ->> '$.status' = 'CONTRACT REGISTERED' "
		inUseObject, _ := GetConditionalObjects(dbConnection, FixedAssetContractStatusTable, conditionString)

		fixedAssetContractStatusInfo := make(map[string]interface{})
		json.Unmarshal((*inUseObject)[0].ObjectInfo, &fixedAssetContractStatusInfo)

		userList := util.InterfaceToIntArray(fixedAssetContractStatusInfo["userList"])
		// Sending email to configured user
		for _, internalAccountId := range userList {
			v.emailGenerator(dbConnection, FixedAssetContractCreationEmailTemplate, internalAccountId, FixedAssetContractComponent, createdRecordId)
		}

	case FixedAssetMasterComponent:
		_, generalWorkOrder := Get(dbConnection, FixedAssetMasterTable, createdRecordId)

		fixedAssetMaster := FixedAssetMaster{ObjectInfo: generalWorkOrder.ObjectInfo}
		fixedAssetInfo := fixedAssetMaster.getFixedAssetMasterInfo()

		if createdRecordId < 10 {
			fixedAssetInfo.AssetNumber = AssetNumberPrefix + "0000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100 {
			fixedAssetInfo.AssetNumber = AssetNumberPrefix + "000" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 1000 {
			fixedAssetInfo.AssetNumber = AssetNumberPrefix + "00" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 10000 {
			fixedAssetInfo.AssetNumber = AssetNumberPrefix + "0" + strconv.Itoa(createdRecordId)
		} else if createdRecordId < 100000 {
			fixedAssetInfo.AssetNumber = AssetNumberPrefix + "SR" + strconv.Itoa(createdRecordId)
		}

		userId := common.GetUserId(ctx)
		v.CreateUserRecordMessage(ProjectID, componentName, "New resource is created", createdRecordId, userId, nil, nil)

		// Create the barcode
		resourceToken, _ := auth.CreateResourceToken(projectId, FixedAssetMasterComponent, createdRecordId)
		url := v.QRCodeDomainUrl + "?lms_token=" + resourceToken
		qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
		v.BaseService.Logger.Info("created QR token", zap.String("token", resourceToken))

		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 200, 200)

		// create the output file
		fileSuffixName := uuid.New().String()
		fileName := v.ComponentContentConfig.DownloadDirectory + "/" + fileSuffixName + ".png"

		file, err := os.Create(fileName)

		if err != nil {
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "QR Code Save Failed",
					Description: "Internal system error during QR code storing, this is normally happen when the backend is not properly configured the path",
				})
			return
		}
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)
		var params map[string]string
		err, httpResponse, rawResponse := util.FileUploadFromDisk(v.ComponentContentConfig.UpStream, params, "file", fileName)

		if err == nil && httpResponse.StatusCode == 200 {
			// no error
			fmt.Println("raw response :", string(rawResponse))
			fileUploadResponse := FileUploadResponse{}
			json.Unmarshal(rawResponse, &fileUploadResponse)
			fmt.Println("contentResponse: ", fileUploadResponse)
			fixedAssetInfo.QRCodeUrl = fileUploadResponse.Url
		} else {
			fmt.Println("raw response :", string(rawResponse))
			response.DispatchDetailedError(ctx, common.InvalidObjectStatus,
				&response.DetailedError{
					Header:      "QR Code Generation Failed",
					Description: "Internal system error during QR code generation, this is normally happen when the backend is not properly configured the path",
				})
			return
		}
		updatingData := make(map[string]interface{})
		rawWorkOrderInfo, _ := json.Marshal(fixedAssetInfo)
		updatingData["object_info"] = rawWorkOrderInfo

		Update(dbConnection, FixedAssetMasterTable, createdRecordId, updatingData)

		conditionString := " object_info ->> '$.status' = 'CONTRACT REGISTERED' "
		inUseObject, _ := GetConditionalObjects(dbConnection, FixedAssetStatusTable, conditionString)

		fixedAssetStatusInfo := make(map[string]interface{})
		json.Unmarshal((*inUseObject)[0].ObjectInfo, &fixedAssetStatusInfo)

		userList := util.InterfaceToIntArray(fixedAssetStatusInfo["userList"])
		for _, internalAccountId := range userList {
			v.emailGenerator(dbConnection, FixedAssetCreationEmailTemplate, internalAccountId, FixedAssetMasterComponent, createdRecordId)
		}

		// Email sending part
		authService := common.GetService("general_auth").ServiceInterface.(common.AuthInterface)
		listOfDepartmentHod := authService.GetHeadOfDepartments(userId)

		// Sending email to HODs
		for _, hodAccountId := range listOfDepartmentHod {
			v.emailGenerator(dbConnection, FixedAssetCreationEmailTemplate, hodAccountId.UserId, FixedAssetMasterComponent, createdRecordId)
		}
	case FixedAssetMyTransferComponent:
		assignee := util.InterfaceToInt(createRequest["assignee"])
		fmt.Println("========================================================================")
		fmt.Println("assignee", assignee)
		// Sending email to HODs

		v.emailGenerator(dbConnection, FixedAssetTransferEmailTemplate, assignee, FixedAssetMyTransferComponent, createdRecordId)

	}

	ctx.JSON(http.StatusOK, component.GeneralResponse{
		Message: "Successfully created",
		Error:   0,
	})
}

// getNewRecord ShowAccount godoc
// @Summary Get the new record based on record schema
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/new_record [get]
func (v *FixedAssetService) getNewRecord(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	projectId := ctx.Param("projectId")
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	newRecordResponse := v.ComponentManager.GetNewRecordResponse(zone, dbConnection, componentName)
	ctx.JSON(http.StatusOK, newRecordResponse)

}

// getRecordFormData ShowAccount godoc
// @Summary Get the record form data to facilitate the update
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentName}/record/{recordId} [get]
func (v *FixedAssetService) getRecordFormData(ctx *gin.Context) {

	// first get the record
	componentName := ctx.Param("componentName")

	projectId := ctx.Param("projectId")
	recordId := ctx.Param("recordId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	intRecordId, _ := strconv.Atoi(recordId)
	dbConnection := v.BaseService.ServiceDatabases[projectId]

	fmt.Println("target table : ", targetTable, "record : ", intRecordId)
	err, generalObject := Get(dbConnection, targetTable, intRecordId)

	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("getting record information has failed"), ErrorGettingIndividualObjectInformation)
		return
	}
	rawObjectInfo := generalObject.ObjectInfo
	rawJSONObject := common.AddFieldJSONObject(rawObjectInfo, "id", recordId)
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	response := v.ComponentManager.GetIndividualRecordResponse(zone, dbConnection, intRecordId, componentName, rawJSONObject)

	ctx.JSON(http.StatusOK, response)

}

type Re struct {
	Data []struct {
		Id    int    `json:"id"`
		Value string `json:"value"`
	} `json:"data"`
	IsEdit        bool   `json:"isEdit"`
	Value         string `json:"value"`
	IsExternal    bool   `json:"isExternal"`
	Property      string `json:"property"`
	Index         int    `json:"index"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	InterfaceType string `json:"interfaceType"`
}

var some = `[
    {
        "data":
        [
            {
                "id": 1,
                "value": "Windows"
            },
            {
                "id": 2,
                "value": "Linux"
            },
            {
                "id": 3,
                "value": "MS Server"
            }
        ],
        "isEdit": true,
        "value": "Windows",
        "isExternal": false,
		"property":"osTypes",
        "index": 1,
        "type": "int",
        "label": "OS Types",
        "interfaceType": "singleDropdown"
    },
    {
        "data":
        [
            {
                "id": 1,
                "value": "2019"
            },
            {
                "id": 2,
                "value": "2018"
            },
            {
                "id": 3,
                "value": "2020"
            }
        ],
        "isEdit": true,
        "value": "2019",
        "isExternal": false,
		"property":"officeVersions",
        "index": 1,
        "type": "int",
        "label": "Office Version",
        "interfaceType": "singleDropdown"
    }
]`

func (v *FixedAssetService) getDynamicFields(ctx *gin.Context) {

	//componentName := util.GetComponentName(ctx)
	projectId := util.GetProjectId(ctx)
	recordId := util.GetRecordId(ctx)
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	fmt.Println("recordId : ", recordId)
	_, fixedAssetObject := Get(dbConnection, FixedAssetMasterTable, recordId)
	var fixedAssetFields = make(map[string]interface{})
	json.Unmarshal(fixedAssetObject.ObjectInfo, &fixedAssetFields)

	assetConfigurationInterface := fixedAssetFields["assetConfiguration"].(map[string]interface{})
	_, generalObject := Get(dbConnection, FixedAssetDynamicFieldConfigurationTable, 1)

	_, generalDynamicFieldTable := Get(dbConnection, FixedAssetDynamicFieldTable, 1)

	var dynamicFieldDataFields = make(map[string]interface{})
	json.Unmarshal(generalDynamicFieldTable.ObjectInfo, &dynamicFieldDataFields)

	dynamicFieldConfiguration := FixedAssetDynamicFieldConfiguration{ObjectInfo: generalObject.ObjectInfo}
	dynamicFields := dynamicFieldConfiguration.getFixedAssetDynamicFieldInfo().DynamicFields
	var arrayOfDynamicFieldResponse []component.RecordInfo
	fmt.Println("dynamicFields : ", assetConfigurationInterface)
	for _, dynamicFieldObject := range dynamicFields {
		recordInfo := component.RecordInfo{}
		fmt.Println("dynamicFieldObject.Property :l ", dynamicFieldObject.Property)
		if value, ok := assetConfigurationInterface[dynamicFieldObject.Property]; ok {
			if dynamicFieldObject.InterfaceType == "singleDropdown" {
				recordInfo.Index = util.InterfaceToInt(value)
				recordInfo.InterfaceType = dynamicFieldObject.InterfaceType
				recordInfo.IsExternal = false
				recordInfo.Type = "int"
				recordInfo.IsEdit = true
				recordInfo.Property = dynamicFieldObject.Property
				recordInfo.Label = dynamicFieldObject.Label
				recordInfo.GridSystem = dynamicFieldObject.GridSystem
				arrayOfData := dynamicFieldDataFields[dynamicFieldObject.Property].([]interface{})
				var dropDownArray []component.OrderedData
				for _, dataElements := range arrayOfData {
					dataElementsMap := dataElements.(map[string]interface{})

					if value, ok := dataElementsMap["id"]; ok {
						if recordInfo.Index == util.InterfaceToInt(value) {
							recordInfo.Value = util.InterfaceToString(dataElementsMap["name"])
						}
						dropDownArray = append(dropDownArray, component.OrderedData{
							Id:    util.InterfaceToInt(value),
							Value: util.InterfaceToString(dataElementsMap["name"]),
						})
					}

				}
				recordInfo.Data = dropDownArray
			} else {
				if dynamicFieldObject.InterfaceType == "textInput" {

					recordInfo.InterfaceType = dynamicFieldObject.InterfaceType
					recordInfo.IsExternal = false
					recordInfo.Type = "text"
					recordInfo.IsEdit = true
					recordInfo.Property = dynamicFieldObject.Property
					recordInfo.Label = dynamicFieldObject.Label
					recordInfo.Value = value
					recordInfo.GridSystem = dynamicFieldObject.GridSystem

				}
			}
			arrayOfDynamicFieldResponse = append(arrayOfDynamicFieldResponse, recordInfo)
		}

	}
	//	dynamicFieldObject := FixedAssetDynamicField{ObjectInfo: (*listOfObjects)[0].ObjectInfo}
	//var re []Re
	//json.Unmarshal([]byte(some), &re)
	ctx.JSON(http.StatusOK, arrayOfDynamicFieldResponse)

}

// getSearchResults ShowAccount godoc
// @Summary Get the search results based on given input
// @Description based on user permission, user will get the list of machine assigned and all the details
// @Tags MachineManagement
// @Accept  json
// @Param SearchField body SearchKeys true "Pass the array of key and values"
// @Param   projectId     path    string     true        "Project Id"
// @Param   ComponentId     path    string     true        "Component Id"
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /project/{projectId}/machines/component/{componentId}/search [post]
func (v *FixedAssetService) getSearchResults(ctx *gin.Context) {

	var searchFieldCommand []component.SearchKeys
	projectId := ctx.Param("projectId")
	componentName := ctx.Param("componentName")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	var totalRecords int64
	if err := ctx.ShouldBindBodyWith(&searchFieldCommand, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dbConnection := v.BaseService.ServiceDatabases[projectId]
	userId := common.GetUserId(ctx)
	zone := getUserTimezone(userId)
	if len(searchFieldCommand) == 0 {
		// reset the search
		listOfObjects, err := GetObjects(dbConnection, targetTable)
		totalRecords = int64(len(*listOfObjects))
		err, tableObjectResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, totalRecords, componentName, "", zone)
		if err != nil {
			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
			return
		}
		ctx.JSON(http.StatusOK, tableObjectResponse)
		return
	}

	format := ctx.Query("format")
	searchQuery := v.ComponentManager.GetSearchQuery(componentName, searchFieldCommand)
	listOfObjects, err := GetConditionalObjects(dbConnection, targetTable, searchQuery)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting machines register information"), ErrorGettingObjectsInformation)
		return
	}
	if format != "" {
		if format == "card_view" {
			cardViewResponse := v.ComponentManager.GetCardViewResponse(listOfObjects, componentName)
			ctx.JSON(http.StatusOK, cardViewResponse)
			return
		} else {

			response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("invalid format, only card_view format is available"), ErrorGettingObjectsInformation)
			return

		}
	}

	_, searchResponse := v.ComponentManager.GetTableRecords(dbConnection, listOfObjects, int64(len(*listOfObjects)), componentName, "", zone)
	ctx.JSON(http.StatusOK, searchResponse)

}

func (v *FixedAssetService) removeInternalArrayReference(ctx *gin.Context) {

	componentName := ctx.Param("componentName")
	//projectId := ctx.Param("projectId")
	targetTable := v.ComponentManager.GetTargetTable(componentName)
	recordIdString := ctx.Param("recordId")
	intRecordId, _ := strconv.Atoi(recordIdString)
	dbConnection := v.BaseService.ReferenceDatabase

	var removeInternalReferenceRequest = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&removeInternalReferenceRequest, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatingData := make(map[string]interface{})
	err, objectInterface := Get(dbConnection, targetTable, intRecordId)
	if err != nil {
		v.BaseService.Logger.Error("error getting given record", zap.String("error", err.Error()))
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error getting requested record"), ErrorGettingIndividualObjectInformation)
		return
	}
	fmt.Println("objectInterface:", objectInterface)
	serializedObject := v.ComponentManager.ProcessInternalArrayReferenceRequest(removeInternalReferenceRequest, objectInterface.ObjectInfo, componentName)
	updatingData["object_info"] = serializedObject
	err = Update(v.BaseService.ReferenceDatabase, targetTable, intRecordId, updatingData)
	if err != nil {
		response.SendSimpleError(ctx, http.StatusBadRequest, errors.New("error updating machines register information"), ErrorUpdatingObjectInformation)
		return
	}
	var updatingObjectFields map[string]interface{}
	json.Unmarshal(serializedObject, &updatingObjectFields)
	ctx.JSON(http.StatusOK, updatingObjectFields)

}
