package handler

import (
	"context"
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/hashicorp/go-getter"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TestDatasourceResponse struct {
	Code     int                  `json:"code"`
	Message  string               `json:"message"`
	MetaInfo component.RecordInfo `json:"metaInfo"`
}

type ConnectionStatus interface {
	TestConnection(connectionParameter datatypes.JSON) (bool, error)
	GetRecordInfo() component.RecordInfo
}

type MySQLConnection struct {
	RecordInfo component.RecordInfo
}

type CSVFileConnection struct {
	RecordInfo component.RecordInfo
}

type DataSchema struct {
	Column    string      `json:"column"`
	Droppable bool        `json:"droppable"`
	Data      string      `json:"data"`
	Id        int         `json:"id"`
	DataType  string      `json:"dataType"`
	IsDefault bool        `json:"isDefault"`
	Default   interface{} `json:"default"`
}

func GetGetter(transport string) map[string]getter.Getter {
	////provide the getter needed to download the files
	if transport == "https" {
		getter := map[string]getter.Getter{
			"https": &getter.HttpGetter{},
		}
		return getter
	}

	return map[string]getter.Getter{
		"http": &getter.HttpGetter{},
	}
}

func (v *CSVFileConnection) TestConnection(connectionParameter datatypes.JSON) (bool, error) {
	csvFileConnectionParameter := CSVFileConnectionParameters{}
	json.Unmarshal(connectionParameter, &csvFileConnectionParameter)
	// now download the file in /tmp folder
	savedFileName := "/tmp/" + uuid.New().String()

	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: savedFileName,
		Dir: false,
		//the repository with a subdirectory I would like to clone only
		Src:  csvFileConnectionParameter.URL,
		Mode: getter.ClientModeFile,
		////define the type of detectors go getter should use, in this case only github is needed

		////provide the getter needed to download the files
		Getters:  GetGetter("https"),
		Insecure: false,
	}
	//download the files
	if err := client.Get(); err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}

	file, err := os.Open(savedFileName)
	if err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}
	defer file.Close()
	// now we saved the file

	// read the first line to extract the colum
	savedCSVFilePointer, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}
	defer savedCSVFilePointer.Close()

	reader := csv.NewReader(savedCSVFilePointer)
	record, err := reader.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}
	fmt.Println("record: ", record)
	var csvSchema []DataSchema
	for index, individualRecord := range record {
		csvSchema = append(csvSchema, DataSchema{
			Column:    individualRecord,
			Droppable: false,
			Data:      individualRecord,
			Id:        index,
			DataType:  "string",
		})
	}

	var recordInfo component.RecordInfo
	type MetaInfo struct {
		Schema  []DataSchema `json:"schema"`
		Records int          `json:"records"`
	}
	var metaInfo = MetaInfo{
		Schema:  csvSchema,
		Records: 100,
	}
	recordInfo.Data = metaInfo
	v.RecordInfo = recordInfo
	return true, nil
}

func (v *CSVFileConnection) ImportRefreshData(url string, csvDatasource CSVDataSource, dbConnection *gorm.DB) (bool, error) {

	dstTable := csvDatasource.ConnectedDatabaseTables[0]
	// now download the file in /tmp folder
	savedFileName := "/tmp/" + uuid.New().String()

	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: savedFileName,
		Dir: false,
		//the repository with a subdirectory I would like to clone only
		Src:  url,
		Mode: getter.ClientModeFile,
		////define the type of detectors go getter should use, in this case only github is needed

		////provide the getter needed to download the files
		Getters:  GetGetter("https"),
		Insecure: false,
	}
	//download the files
	if err := client.Get(); err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}

	file, err := os.Open(savedFileName)
	if err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}
	defer file.Close()
	// now we saved the file

	// read the first line to extract the colum
	savedCSVFilePointer, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}
	defer savedCSVFilePointer.Close()

	reader := csv.NewReader(savedCSVFilePointer)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}

	var insertStatement string

	// get the colum insert
	var insertColumns string
	insertColumns = "("
	for _, schemaElement := range csvDatasource.DatasourceAttribute.DestinationSchema {
		insertColumns += "`" + schemaElement.Data + "`" + ","
	}
	insertColumns = insertColumns[:len(insertColumns)-1]
	insertColumns += ") "
	for index, individualRecord := range records {
		if index == 0 {
			// ignore the first colum, to avoid the colum names to repeat
			continue
		}
		// now we needed to check the unique keys, and if not , we will insert, otherwise we will update
		var uniqueKey string
		// we need to create record, and insert to database
		insertStatement = "INSERT INTO " + dstTable + insertColumns + " VALUES("
		// we should only process the data schema loop, and get the record accordingly
		for _, dataSchemaElement := range csvDatasource.DatasourceAttribute.DestinationSchema {
			// check if it is unique
			individualCellValue := individualRecord[dataSchemaElement.Id]
			if util.ArrayContains(csvDatasource.DatasourceAttribute.UniqueIndexColumns, dataSchemaElement.Column) {
				uniqueKey += dataSchemaElement.Data + "='" + util.InterfaceToString(individualCellValue) + "' AND "
			}

			if dataSchemaElement.DataType == "string" {
				if individualCellValue == "" {
					if dataSchemaElement.Default == nil {
						insertStatement += "'-',"
					} else {
						insertStatement += " '" + dataSchemaElement.Default.(string) + "',"
					}

				} else {
					insertStatement += " '" + individualCellValue + "',"
				}

			}
			if dataSchemaElement.DataType == "double" {
				if individualCellValue == "" {
					doubleValue := util.InterfaceToFloat(dataSchemaElement.Default)
					doubleStringValue := strconv.FormatFloat(doubleValue, 'f', -1, 64)
					insertStatement += doubleStringValue + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
			if dataSchemaElement.DataType == "boolean" {
				if individualCellValue == "" {
					insertStatement += dataSchemaElement.Default.(string) + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
			if dataSchemaElement.DataType == "integer" {
				if individualCellValue == "" {
					insertStatement += strconv.Itoa(dataSchemaElement.Default.(int)) + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
		}
		lastIndex := strings.LastIndex(uniqueKey, "AND ")
		if lastIndex != -1 {
			// Remove the last occurrence of "AND" from the string
			uniqueKey = uniqueKey[:lastIndex]
		}
		err, listOfObjects := GetConditionalObjectsV2(dbConnection, dstTable, uniqueKey)
		if err == nil {
			if len(listOfObjects) == 0 {
				// insert it
				insertStatement = insertStatement[:len(insertStatement)-1]
				insertStatement += ")"
				fmt.Println("insertStatement : ", insertStatement)
				err = dbConnection.Exec(insertStatement).Error
				if err != nil {
					fmt.Println("error inserting into database", "query", insertStatement)

				}
			} else if len(listOfObjects) == 1 {
				//update it
			}
		} else {
			fmt.Println("error getting conditional data", err.Error())
		}

	}

	return true, err
}

func (v *CSVFileConnection) ImportData(url string, targetTable string, dataSchema []DataSchema, dbConnection *gorm.DB) (bool, error) {
	// now download the file in /tmp folder
	savedFileName := "/tmp/" + uuid.New().String()

	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: savedFileName,
		Dir: false,
		//the repository with a subdirectory I would like to clone only
		Src:  url,
		Mode: getter.ClientModeFile,
		////define the type of detectors go getter should use, in this case only github is needed

		////provide the getter needed to download the files
		Getters:  GetGetter("https"),
		Insecure: false,
	}
	//download the files
	if err := client.Get(); err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}

	file, err := os.Open(savedFileName)
	if err != nil {
		return false, errors.New("failed to process given file, check the file again")
	}
	defer file.Close()
	// now we saved the file

	// read the first line to extract the colum
	savedCSVFilePointer, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}
	defer savedCSVFilePointer.Close()

	reader := csv.NewReader(savedCSVFilePointer)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return false, errors.New("failed to process given file, check the file again")
	}
	var insertStatement string

	// get the colum insert
	var insertColumns string
	insertColumns = "("
	for _, schemaElement := range dataSchema {
		insertColumns += "`" + schemaElement.Data + "`" + ","
	}
	insertColumns = insertColumns[:len(insertColumns)-1]
	insertColumns += ") "
	for index, individualRecord := range records {
		if index == 0 {
			// ignore the first colum, to avoid the colum names to repeat
			continue
		}
		// we need to create record, and insert to database
		insertStatement = "INSERT INTO " + targetTable + insertColumns + " VALUES("
		// we should only process the data schema loop, and get the record accordingly
		for index, dataSchemaElement := range dataSchema {
			individualCellValue := individualRecord[index]
			if dataSchemaElement.DataType == "string" {
				if individualCellValue == "" {
					if dataSchemaElement.Default == nil {
						insertStatement += "'-',"
					} else {
						insertStatement += " '" + dataSchemaElement.Default.(string) + "',"
					}

				} else {
					insertStatement += " '" + individualCellValue + "',"
				}

			}
			if dataSchemaElement.DataType == "double" {
				if individualCellValue == "" {
					doubleValue := util.InterfaceToFloat(dataSchemaElement.Default)
					doubleStringValue := strconv.FormatFloat(doubleValue, 'f', -1, 64)
					insertStatement += doubleStringValue + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
			if dataSchemaElement.DataType == "boolean" {
				if individualCellValue == "" {
					insertStatement += dataSchemaElement.Default.(string) + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
			if dataSchemaElement.DataType == "integer" {
				if individualCellValue == "" {
					insertStatement += strconv.Itoa(dataSchemaElement.Default.(int)) + ","
				} else {
					insertStatement += individualCellValue + ","
				}
			}
		}
		insertStatement = insertStatement[:len(insertStatement)-1]
		insertStatement += ")"
		fmt.Println("insertStatement : ", insertStatement)
		err = dbConnection.Exec(insertStatement).Error
		if err != nil {
			os.Exit(0)
		}

	}

	return true, err
}
func (v *CSVFileConnection) GetRecordInfo() component.RecordInfo {
	return v.RecordInfo
}
func (v *MySQLConnection) GetRecordInfo() component.RecordInfo {
	return v.RecordInfo
}
func (v *MySQLConnection) TestConnection(connectionParameter datatypes.JSON) (bool, error) {
	mysqlConnectionParameter := MysqlConnectionParameters{}

	json.Unmarshal(connectionParameter, &mysqlConnectionParameter)

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConnectionParameter.Username, mysqlConnectionParameter.Password, mysqlConnectionParameter.HostName, mysqlConnectionParameter.Port, mysqlConnectionParameter.Schema)
	dbConnection, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return false, err
	} else {
		tableInformationSchema := "SELECT * FROM information_schema.tables WHERE table_schema =\"" + mysqlConnectionParameter.Schema + "\""
		var queryResults []map[string]interface{}
		dbConnection.Raw(tableInformationSchema).Scan(&queryResults)
		var recordInfo component.RecordInfo
		var dropDownArray []component.OrderedData
		var index = 0
		for _, result := range queryResults {
			index += 1
			tableName := result["TABLE_NAME"]
			dropDownArray = append(dropDownArray, component.OrderedData{
				Id:    index,
				Value: tableName.(string),
			})
			if index == 1 {
				recordInfo.Index = index
				recordInfo.Value = tableName
			}
		}
		recordInfo.Data = dropDownArray
		v.RecordInfo = recordInfo
	}
	return true, nil
}
func (as *AnalyticsService) testDatasourceConnection(ctx *gin.Context) {
	datasourceName := ctx.Param("datasourceName")
	var connectionFields = make(map[string]interface{})

	if err := ctx.ShouldBindBodyWith(&connectionFields, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if datasourceName == "mysql" {
		var connectionStatus ConnectionStatus
		connectionStatus = &MySQLConnection{}
		serialisedData, _ := json.Marshal(connectionFields)
		status, err := connectionStatus.TestConnection(serialisedData)
		if !status {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Connection Failed"), ConnectingDatasourceFailed, err.Error())
			return
		} else {
			var datasourceDBResponse = TestDatasourceResponse{}
			datasourceDBResponse.Code = 0
			datasourceDBResponse.Message = "Your parameters are successfully tested, now you can proceed to save the datasource"
			datasourceDBResponse.MetaInfo = connectionStatus.GetRecordInfo()

			ctx.JSON(http.StatusOK, datasourceDBResponse)
		}

	}
	if datasourceName == "csv" {
		var connectionStatus ConnectionStatus
		connectionStatus = &CSVFileConnection{}
		serialisedData, _ := json.Marshal(connectionFields)
		status, err := connectionStatus.TestConnection(serialisedData)
		if !status {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Connection Failed"), ConnectingDatasourceFailed, err.Error())
			return
		} else {
			var datasourceDBResponse = TestDatasourceResponse{}
			datasourceDBResponse.Code = 0
			datasourceDBResponse.Message = "Your parameters are successfully tested, now you can proceed to save the datasource"
			datasourceDBResponse.MetaInfo = connectionStatus.GetRecordInfo()

			ctx.JSON(http.StatusOK, datasourceDBResponse)
		}

	}
	if datasourceName == "csv_link" {
		var connectionStatus ConnectionStatus
		connectionStatus = &CSVFileConnection{}
		serialisedData, _ := json.Marshal(connectionFields)
		status, err := connectionStatus.TestConnection(serialisedData)
		if !status {
			response.SendDetailedError(ctx, http.StatusBadRequest, getError("Connection Failed"), ConnectingDatasourceFailed, err.Error())
			return
		} else {
			var datasourceDBResponse = TestDatasourceResponse{}
			datasourceDBResponse.Code = 0
			datasourceDBResponse.Message = "Your parameters are successfully tested, now you can proceed to save the datasource"
			datasourceDBResponse.MetaInfo = connectionStatus.GetRecordInfo()

			ctx.JSON(http.StatusOK, datasourceDBResponse)
		}

	}

}
