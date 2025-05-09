package handler

import (
	"cx-micro-flake/pkg/response"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSVConnectionRequest struct {
	ConnectionParam struct {
		Url string `json:"url"`
	} `json:"connectionParam"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	DatasourceAttribute struct {
		TargetTable        string       `json:"targetTable"`
		SourceSchema       []DataSchema `json:"sourceSchema"`
		DestinationSchema  []DataSchema `json:"destinationSchema"`
		UniqueIndexColumns []string     `json:"uniqueIndexColumns"`
	} `json:"datasourceAttribute"`
}

type CSVDataSource struct {
	Name            string        `json:"name"`
	Status          string        `json:"status"`
	CreatedAt       time.Time     `json:"createdAt"`
	CreatedBy       int           `json:"createdBy"`
	Description     string        `json:"description"`
	Permissions     []interface{} `json:"permissions"`
	ObjectStatus    string        `json:"objectStatus"`
	LastUpdatedAt   time.Time     `json:"lastUpdatedAt"`
	LastUpdatedBy   int           `json:"lastUpdatedBy"`
	ConnectionParam struct {
		Url string `json:"url"`
	} `json:"connectionParam"`
	DatasourceMaster    int `json:"datasourceMaster"`
	DatasourceAttribute struct {
		TargetTable        string        `json:"targetTable"`
		SourceSchema       []interface{} `json:"sourceSchema"`
		DestinationSchema  []DataSchema  `json:"destinationSchema"`
		UniqueIndexColumns []string      `json:"uniqueIndexColumns"`
	} `json:"datasourceAttribute"`
	ConnectedDatabaseTables []string `json:"connectedDatabaseTables"`
}

func (as *AnalyticsService) handleCSV2DatabaseRecordCreation(ctx *gin.Context, dbConnection *gorm.DB, request map[string]interface{}) error {

	serialisedRequest, _ := json.Marshal(request)
	fmt.Println("serialised request : ", string(serialisedRequest))
	csvFileConnectionRequest := CSVConnectionRequest{}
	json.Unmarshal(serialisedRequest, &csvFileConnectionRequest)
	// create the table
	fmt.Println("csvFileConnectionRequest request : ", csvFileConnectionRequest)
	type CSVDatasourceAttributes struct {
		TargetTable       string       `json:"targetTable"`
		SourceSchema      []DataSchema `json:"sourceSchema"`
		DestinationSchema []DataSchema `json:"destinationSchema"`
	}

	// check table is there created already
	//CREATE TABLE `fuyu_mes`.`new_table` (`id` INT NOT NULL AUTO_INCREMENT,`student_marks` VARCHAR(45) NULL,PRIMARY KEY (`id`));
	var tableExecutionString = "CREATE TABLE `" + csvFileConnectionRequest.DatasourceAttribute.TargetTable + "` (`record_id` INT NOT NULL AUTO_INCREMENT,"
	for _, dstSchema := range csvFileConnectionRequest.DatasourceAttribute.DestinationSchema {
		if dstSchema.DataType == "integer" {
			tableExecutionString += "`" + dstSchema.Data + "` INT,"
		}
		if dstSchema.DataType == "string" {
			tableExecutionString += "`" + dstSchema.Data + "` VARCHAR(512) NULL,"
		}
		if dstSchema.DataType == "double" {
			tableExecutionString += "`" + dstSchema.Data + "` DOUBLE NULL,"
		}

		if dstSchema.DataType == "boolean" {
			tableExecutionString += "`" + dstSchema.Data + "` INT(1) NULL,"
		}
	}
	tableExecutionString = tableExecutionString[:len(tableExecutionString)-1]
	tableExecutionString += ",PRIMARY KEY (`record_id`));"
	err := dbConnection.Exec(tableExecutionString).Error

	as.BaseService.Logger.Info("created table creation string", zap.String("execution_string", tableExecutionString))
	if err != nil {
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Invalid Schema"), ErrorCreatingObjectInformation, err.Error())
		return err
	}
	var connectedTables []string
	connectedTables = append(connectedTables, csvFileConnectionRequest.DatasourceAttribute.TargetTable)
	request["connectedDatabaseTables"] = connectedTables

	// import the data
	fileUrl := csvFileConnectionRequest.ConnectionParam.Url
	csvConnection := CSVFileConnection{}
	_, err = csvConnection.ImportData(fileUrl, csvFileConnectionRequest.DatasourceAttribute.TargetTable, csvFileConnectionRequest.DatasourceAttribute.DestinationSchema, dbConnection)
	if err != nil {
		dbConnection.Exec("DROP TABLE " + csvFileConnectionRequest.DatasourceAttribute.TargetTable)
		response.SendDetailedError(ctx, http.StatusBadRequest, getError("Failed to process data"), ErrorCreatingObjectInformation, err.Error())
		return err
	}
	return nil
}
