package component

import (
	"cx-micro-flake/pkg/response"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/datatypes"
	"net/http"
	"strconv"
)

func GetRequestFields(ctx *gin.Context) (error, map[string]interface{}) {
	var requestFields = make(map[string]interface{})
	if err := ctx.ShouldBindBodyWith(&requestFields, binding.JSON); err != nil {
		return err, requestFields
	}
	return nil, requestFields

}
func GetGeneralObject(fields map[string]interface{}) GeneralObject {
	fields["objectStatus"] = ObjectStatusActive
	rawCreateRequest, _ := json.Marshal(fields)
	object := GeneralObject{
		ObjectInfo: rawCreateRequest,
	}
	return object
}

func GetGeneralObjectFromInterface(object interface{}) GeneralObject {
	serializedObject, _ := json.Marshal(object)
	generalObject := GeneralObject{ObjectInfo: serializedObject}
	return generalObject
}
func SendObjectCreationResponse(ctx *gin.Context, projectId, componentName string, recordId int) {
	url := "/project/" + projectId + "/moulds/component/" + componentName + "/record/" + strconv.Itoa(recordId)
	ctx.Writer.Header().Set("Location", url)
	ctx.JSON(http.StatusCreated, response.GeneralResponse{
		Code:     0,
		RecordId: recordId,
		Message:  "New resource is successfully created",
	})
}

type SearchKeys struct {
	Field string `json:"key"`
	Value string `json:"value"`
}

type TableObjectResponse struct {
	TotalRowCount   int64            `json:"totalRowCount,omitempty"`
	CurrentRowCount int64            `json:"currentRowCount,omitempty"`
	Header          []TableSchema    `json:"header,omitempty"`
	Data            []datatypes.JSON `json:"data"`
	IsNext          bool             `json:"isNext"`
}

type WidgetTableObjectResponse struct {
	TotalRowCount   int64            `json:"totalRowCount,omitempty"`
	CurrentRowCount int64            `json:"currentRowCount,omitempty"`
	Header          []TableSchema    `json:"header,omitempty"`
	Data            []datatypes.JSON `json:"data"`
	XColumnList     []string         `json:"xColumnList"`
	YColumnList     []string         `json:"yColumnList"`
}

type LoadFileResponse struct {
	TotalRowCount int64             `json:"totalRowCount"`
	Data          []datatypes.JSON  `json:"data"`
	TableSchema   []CRUDTableSchema `json:"header"`
}

type ImportDataObjects struct {
	TotalRecords        int             `json:"totalRecords"`
	TotalSkippedRecords int             `json:"totalSkippedRecords"`
	SkippedData         interface{}     `json:"skippedData"`
	InsertObjects       []GeneralObject `json:"insertObjects"`
	SkippedRecordNames  string          `json:"skippedRecordNames"`
	UpdateObjects       []GeneralObject `json:"updateObjects"`
}

type ImportDataResponse struct {
	TotalRecords        int         `json:"totalRecords"`
	FailedRecords       int         `json:"failedRecords"`
	TotalSkippedRecords int         `json:"totalSkippedRecords"`
	SkippedData         interface{} `json:"skippedData"`
	Message             string      `json:"message"`
}

type ImportSchemaRequest struct {
	SourceField      string `json:"sourceField"`
	DestinationField string `json:"destinationField"`
}

type ImportDataCommand struct {
	ContentUrl string                `json:"contentUrl"`
	Schema     []ImportSchemaRequest `json:"schema"`
}
type ExportFieldAttributes struct {
	Operator string      `json:"operator"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
}
type ExportDataCommand struct {
	Format     string                           `json:"format"`
	Data       []ExportSchema                   `json:"data"`
	Attributes map[string]ExportFieldAttributes `json:"attributes"` // { "machineId":4}
}

type LoadDataFileCommand struct {
	ContentUrl string `json:"contentUrl"`
}

type CardViewResponse struct {
	Template string         `json:"template"`
	Cards    datatypes.JSON `json:"cards"`
}

type CardViewGroupResponse struct {
	Template string           `json:"template"`
	Cards    []datatypes.JSON `json:"cards"`
}

type TableDataResponse struct {
	Header        []TableSchema    `json:"header"`
	Data          []datatypes.JSON `json:"data"`
	TotalRowCount int              `json:"totalRowCount"`
}

type ExportDataResponse struct {
	Url     string `json:"url"`
	Name    string `json:"name"`
	Size    string `json:"size"`
	SiteMap map[string]string
}

type GroupByAction struct {
	GroupBy []string `json:"groupBy"`
}
type GroupByChildren struct {
	Data []interface{} `json:"data"`
	Type string        `json:"type"`
}

type TableGroupByResponse struct {
	Label    string        `json:"label"`
	Children []interface{} `json:"children"`
}

func GetGroupByResults(groupByColumn string, dataArray []datatypes.JSON) map[string][]interface{} {
	var groupByResults = make(map[string][]interface{})
	for _, objectInterface := range dataArray {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface, &objectFields)
		var groupByField = util.InterfaceToString(objectFields[groupByColumn])
		groupByResults[groupByField] = append(groupByResults[groupByField], objectInterface)
	}
	return groupByResults
}

func GetGroupByResultsFromInterface(groupByColumn string, dataArray []interface{}) map[string][]interface{} {
	var groupByResults = make(map[string][]interface{})
	for _, objectInterface := range dataArray {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface.(datatypes.JSON), &objectFields)
		var groupByField = util.InterfaceToString(objectFields[groupByColumn])
		groupByResults[groupByField] = append(groupByResults[groupByField], objectInterface)
	}
	return groupByResults
}
