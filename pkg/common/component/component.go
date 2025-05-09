package component

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	"gorm.io/datatypes"
)

const (
	JsonToStringArray    = "json_array_to_string_array"
	ObjectStatusActive   = "Active"
	ObjectStatusArchived = "Archived"
	ISOTimeLayout        = "2006-01-02T15:04:05.000Z"
	// table data types

	TableDataTypeDate     = "date"
	TableDataTypeDateTime = "datetime"
)

type GroupByCardView struct {
	GroupByField string                   `json:"groupByField"`
	Cards        []map[string]interface{} `json:"cards"`
}

type GroupByView struct {
	GroupByField string      `json:"groupByField"`
	DisplayField string      `json:"displayField"`
	Cards        interface{} `json:"cards"`
}

type OverviewData struct {
	Value           []map[string]interface{} `json:"value"`
	Type            string                   `json:"type"`
	IsVisible       bool                     `json:"isVisible"`
	Label           string                   `json:"label"`
	Icon            string                   `json:"icon"`
	BackgroundColor string                   `json:"backgroundColor"`
	IsClickable     bool                     `json:"isClickable"`
	Module          string                   `json:"module"`
	Component       string                   `json:"component"`
	ApiQuery        string                   `json:"apiQuery"`
	MenuId          string                   `json:"menuId"`
}
type OverviewResponse struct {
	Data  []OverviewData `json:"data"`
	Label string         `json:"label"`
}
type ServiceConfig struct {
	ServiceType        string      `json:"serviceType"`
	ServiceName        string      `json:"serviceName"`
	ServiceDescription string      `json:"serviceDescription"`
	ServiceInterface   interface{} `json:"serviceInterface"`
}

type FieldValidation struct {
	Create []*ObjectMapping `json:"create"`
	Update []*ObjectMapping `json:"update"`
}
type CreateFieldValidator struct {
	Validator string `json:"validator"`
}
type FieldValidator struct {
	Create []CreateFieldValidator `json:"create"`
}
type ValidationResult struct {
	ResultStatus bool
	Error        error
}

func InitGeneralObject(objectInfo datatypes.JSON) GeneralObject {
	generalObject := GeneralObject{
		ObjectInfo: objectInfo,
	}
	return generalObject
}

type RecordSchema struct {
	IsEdit                bool             `json:"isEdit,omitempty"`
	Name                  string           `json:"name,omitempty"`
	Type                  string           `json:"type,omitempty"`
	InterfaceType         string           `json:"interfaceType,omitempty"`
	Property              string           `json:"property,omitempty"`
	Default               *interface{}     `json:"default,omitempty"`
	IsMandatory           bool             `json:"isMandatory"`
	LinkedObjectMapping   *ObjectMapping   `json:"linkedObjectMapping"`
	ResponseObjectMapping *ObjectMapping   `json:"responseObjectMapping"`
	CreateValidation      *ObjectMapping   `json:"createValidation,omitempty"`
	DefaultObjectMapping  *ObjectMapping   `json:"defaultObjectMapping,omitempty"`
	RecordSchema          []RecordSchema   `json:"recordSchema"`
	FieldValidation       *FieldValidation `json:"fieldValidation"`
	Formatter             string           `json:"formatter"`
	Overwrite             *bool            `json:"overwrite"`
	ReferenceComponent    string           `json:"referenceComponent"`
	FieldValidator        *FieldValidator  `json:"fieldValidator"`
	IgnoreEmptyInCreate   bool             `json:"ignoreEmptyInCreate"`
	IgnoreEmptyInSending  bool             `json:"ignoreEmptyInSending"`
	UniqueIndex           []string         `json:"uniqueIndex"`

	LinkedDataType *string `json:"linkedDataType"`
	LinkedProperty *string `json:"linkedProperty"`
	Display        bool    `json:"display"`
	Render         bool    `json:"render"`
	GridSystem     string  `json:"gridSystem"`
	Label          string  `json:"label"`

	Icon     string `json:"icon"`
	IconType string `json:"iconType"`

	CanDisplayTable bool `json:"canDisplayTable"`

	IsDynamic                    bool           `json:"isDynamic"`
	DynamicMappingField          string         `json:"dynamicMappingField"`
	DynamicComponent             string         `json:"dynamicComponent"`
	DefaultType                  *string        `json:"defaultType"`
	HeaderObjectMapping          *ObjectMapping `json:"headerObjectMapping,omitempty"`
	HeaderFontColorObjectMapping *ObjectMapping `json:"headerFontColorObjectMapping,omitempty"`
}
type CardViewSchema struct {
	Template string `json:"template"`
	Field    []struct {
		Property     string `json:"property"`
		MappingField string `json:"mappingField"`
	} `json:"field"`
}

type CreateDependencyResourceInjection struct {
	ResourceInjection []ObjectMapping `json:"resourceInjection"`
}
type DeleteDependencyResourceInjection struct {
	ResourceInjection []ObjectMapping `json:"resourceInjection"`
}
type TableSchema struct {
	Name                         string         `json:"name"`
	Type                         string         `json:"type"`
	Display                      bool           `json:"display"`
	Property                     string         `json:"property"`
	RouteEnabled                 bool           `json:"routeEnabled"`
	RouteLink                    string         `json:"routeLink,omitempty"`
	RouteRecordIdProperty        string         `json:"routeRecordIdProperty,omitempty"`
	ObjectList                   datatypes.JSON `json:"objectList,omitempty"`
	FontColorList                datatypes.JSON `json:"fontColorList,omitempty"`
	ReferenceObjectMapping       *ObjectMapping `json:"referenceObjectMapping,omitempty"`
	HeaderObjectMapping          *ObjectMapping `json:"headerObjectMapping,omitempty"`
	HeaderFontColorObjectMapping *ObjectMapping `json:"headerFontColorObjectMapping,omitempty"`
	ColorCode                    string         `json:"colorCode,omitempty"`
	ColumSize                    int            `json:"columSize"`
	LinkedDataType               *string        `json:"linkedDataType"`
	LinkedProperty               *string        `json:"linkedProperty"`
	InterfaceType                string         `json:"interfaceType"`
	Render                       bool           `json:"render"`
	GridSystem                   string         `json:"gridSystem"`
	Label                        string         `json:"label"`
	IsGroupByField               bool           `json:"isGroupByField"`
	Unit                         string         `json:"unit"`
}

func GetTableHeaderSchema(name, property string) TableSchema {
	tableSchema := TableSchema{}
	tableSchema.Name = name
	tableSchema.Property = property
	tableSchema.Display = true

	return tableSchema
}

type LinkedField struct {
	Field        string `json:"field"`
	Query        string `json:"query"`
	LinkedColumn string `json:"linkedColumn"`
}

type OrderedData struct {
	Id                       int          `json:"id"`
	Value                    string       `json:"value"`
	OnValueConditionalFields []RecordInfo `json:"onValueConditionalFields"`
}
type CRUDTableSchema struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Display  bool   `json:"display"`
	Property string `json:"property"`
}
type DynamicConditionalFields struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}
type DynamicDroppedDownAttributes struct {
	Type               int                        `json:"type"`
	AutoFieldsSource   string                     `json:"autoFieldsSource"`
	ManualFieldsSource []string                   `json:"manualFieldsSource"`
	ConditionalFields  []DynamicConditionalFields `json:"conditionalFields"`
}
type RecordInfo struct {
	Data                interface{}     `json:"data,omitempty"`
	IsEdit              bool            `json:"isEdit"`
	Value               interface{}     `json:"value,omitempty"`
	IsExternal          bool            `json:"isExternal"`
	ValueArray          []interface{}   `json:"valueArray,omitempty"`
	Index               int             `json:"index,omitempty"`
	IndexArray          []int           `json:"indexArray"`
	Type                string          `json:"type,omitempty"`
	Header              interface{}     `json:"header,omitempty"`
	CommonRouteLink     string          `json:"commonRouteLink,omitempty"`
	InterfaceType       string          `json:"interfaceType,omitempty"`
	DynamicMappingField string          `json:"dynamicMappingField"`
	DynamicComponent    string          `json:"dynamicComponent"`
	IsDynamic           bool            `json:"isDynamic"`
	GridSystem          string          `json:"gridSystem"`
	Label               string          `json:"label"`
	Icon                string          `json:"icon"`
	IconType            string          `json:"iconType"`
	Property            string          `json:"property"`
	Render              bool            `json:"render"`
	Display             bool            `json:"display"`
	Description         string          `json:"description"`
	InterfaceField      string          `json:"interfaceField"`
	IsMandatoryField    bool            `json:"isMandatoryField"`
	InternalExport      *InternalExport `json:"internalExport,omitempty"`
}

type InternalExport struct {
	IsExportEnabled     bool   `json:"isExportEnabled"`
	CanExport           bool   `json:"canExport"`
	ExportSchema        string `json:"exportSchema"`
	ExportEndpoint      string `json:"exportEndpoint"`
	RefreshExportSchema string `json:"refreshExportSchema"`
}

type CreateAction struct {
	QueryActions []ObjectMapping `json:"queryActions"`
	EmailAction  *ObjectMapping  `json:"emailAction"`
}

type ActionFields struct {
	Type          string         `json:"type"`
	Value         interface{}    `json:"value"`
	Default       *interface{}   `json:"default"`
	ArrayDataType string         `json:"arrayDataType"`
	Property      string         `json:"property"`
	IsEdit        bool           `json:"isEdit"`
	ObjectMapping *ObjectMapping `json:"objectMapping"`
}
type DependencyComponent struct {
	Delete struct {
		LinkedComponent string `json:"linkedComponent"`
		Property        string `json:"property"`
		Query           string `json:"query"`
	} `json:"delete"`
	Add struct {
		LinkedComponent string `json:"linkedComponent"`
		Property        string `json:"property"`
		Query           string `json:"query"`
	} `json:"add"`
}

type TableImportFields struct {
	DataType         string         `json:"dataType"`
	Validation       *ObjectMapping `json:"validation,omitempty"`
	Replacement      *ObjectMapping `json:"replacement,omitempty"`
	DisplayName      string         `json:"displayName"`
	DestinationField string         `json:"destinationField"`
}
type TableImportExtraField struct {
	Field         string         `json:"field"`
	ObjectMapping *ObjectMapping `json:"objectMapping"`
	Default       int            `json:"default,omitempty"`
	DataType      string         `json:"dataType,omitempty"`
}
type TableImportSchema struct {
	ExtraField []TableImportExtraField `json:"extraField"`
	Fields     []TableImportFields     `json:"fields"`
	Update     *ObjectMapping          `json:"update"`
}

type Constraints struct {
	Reference                     string `json:"reference"`
	ReferenceProperty             string `json:"referenceProperty"`
	IsObjectField                 bool   `json:"isObjectField"`
	ReferenceComponentDisplayName string `json:"referenceComponentDisplayName"`
}

type ComponentSchema struct {
	ModuleId                          int                                `json:"moduleId"`
	CreatedAt                         string                             `json:"createdAt"`
	LastUpdatedAt                     string                             `json:"lastUpdatedAt"`
	CreatedBy                         int                                `json:"createdBy"`
	Description                       string                             `json:"description"`
	TableImportSchema                 TableImportSchema                  `json:"tableImportSchema"` // this schema indicates, how to render the front-end crud table
	TableSchema                       []TableSchema                      `json:"tableSchema"`       // how to import the table
	ComponentName                     string                             `json:"componentName"`
	TargetTable                       string                             `json:"targetTable"`
	RecordSchema                      []RecordSchema                     `json:"recordSchema"`
	CardViewSchema                    CardViewSchema                     `json:"cardViewSchema"`
	AdditionalRecords                 []AdditionalRecordSchema           `json:"additionalRecords"` // this will be used to compose additional fields need extra configuration
	TableAdditionalRecords            []AdditionalRecordSchema           `json:"tableAdditionalRecords"`
	ExportSchema                      []ExportSchema                     `json:"exportSchema"`
	CreateDependencyResourceInjection *CreateDependencyResourceInjection `json:"createDependencyResourceInjection"`
	DeleteDependencyResourceInjection *DeleteDependencyResourceInjection `json:"deleteDependencyResourceInjection"`
	Constraints                       []Constraints                      `json:"constraints"`
	LinkedObjectMap                   []LinkedObjectMap                  `json:"linkedObjectMap"`
	ParentFields                      *ParentFields                      `json:"parentFields"`
	ReferenceTemplateComponent        *ReferenceTemplateSchemaMapping    `json:"referenceTemplateComponent"`
	EnabledAfterEntityLevel           *int                               `json:"enabledAfterWorkflowStatusLevel"`
	EntityField                       *string                            `json:"entityField"`
}

type ParentFields struct {
	Level1 string `json:"level1"`
	Level2 string `json:"level2"`
	Level3 string `json:"level3"`
}

type ReferenceTemplateSchemaMapping struct {
	ReferenceComponent string `json:"referenceComponent"`
	TemplateFieldName  string `json:"templateFieldName"`
}

type LinkedObjectMap struct {
	Property       string `json:"property"`
	LinkedTable    string `json:"linkedTable"`
	TargetProperty string `json:"targetProperty"`
}

func GetComponentSchema(config []byte) ComponentSchema {
	componentConfig := ComponentSchema{}
	json.Unmarshal(config, &componentConfig)
	return componentConfig
}

func (v *GeneralObject) Serialised() datatypes.JSON {
	var objectFields = make(map[string]interface{})
	json.Unmarshal(v.ObjectInfo, &objectFields)
	objectFields["id"] = v.Id
	serialisedObject, _ := json.Marshal(objectFields)
	return serialisedObject
}

type AdditionalRecordSchema struct {
	ObjectMapping  ObjectMapping   `json:"objectMapping"`
	Property       string          `json:"property"`
	RecordSchema   []RecordSchema  `json:"recordSchema"`
	InternalExport *InternalExport `json:"internalExport,omitempty"`
}
type GeneralObject struct {
	Id         int
	ObjectInfo datatypes.JSON
}

func GetDefaultValueArrayRecordInfo() RecordInfo {
	var valueArray []interface{}
	valueArray = append(valueArray, "new")
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		ValueArray:      valueArray,
		Index:           0,
		Type:            "bool",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}
func GetDefaultBoolRecordInfo() RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           true,
		Index:           0,
		Type:            "bool",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}
func GetRecordInfo(value, dataType string) RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           value,
		Index:           0,
		Type:            dataType,
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetBoolRecordInfo(value bool, dataType string) RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           value,
		Index:           0,
		Type:            dataType,
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetRecordIntInfo(value int, dataType string) RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           value,
		Index:           0,
		Type:            dataType,
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetRecordObjectInfo(value interface{}, dataType string) RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           value,
		Index:           0,
		Type:            dataType,
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetEmptyDateRecordInfo() RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           "",
		Index:           0,
		Type:            "date",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetDefaultDateRecordInfo() RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           "06/19/2022 10:00:00",
		Index:           0,
		Type:            "date",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetCurrentDateRecordInfo() RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           time.Now().Format("2006-01-02T15:04:05.999Z"),
		Index:           0,
		Type:            "date",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

func GetDefaultRecordInfo() RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           "",
		Index:           0,
		Type:            "",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}
func GetDefaultRecordInfoWithValue(value interface{}) RecordInfo {
	recordInfo := RecordInfo{
		Data:            nil,
		IsEdit:          true,
		Value:           value,
		Index:           0,
		Type:            "",
		Header:          nil,
		CommonRouteLink: "",
	}
	return recordInfo
}

type UpstreamContentConfig struct {
	DownloadDirectory string `json:"downloadDirectory"`
	Insecure          bool   `json:"insecure"`
	Getter            string `json:"getter"`
	UpStream          string `json:"upStream"`
}

type FileData struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Url  string `json:"url"`
}

func (cc *UpstreamContentConfig) GetGetter() map[string]getter.Getter {
	////provide the getter needed to download the files
	if cc.Getter == "https" {
		getter := map[string]getter.Getter{
			"https": &getter.HttpGetter{},
		}
		return getter
	}

	return map[string]getter.Getter{
		"http": &getter.HttpGetter{},
	}
}

func TableCondition(offset, fields, values, condition string) string {
	basicCondition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$."
	var searchQuery string
	if fields != "" && values != "" {
		arrayOfFields := strings.Split(fields, ",")
		arrayOfValues := strings.Split(values, ",")
		if len(arrayOfFields) == len(arrayOfValues) {
			for index, field := range arrayOfFields {
				if len(arrayOfFields)-1 == index {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\""

				} else {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\" " + condition

				}
			}
		}
	}
	if offset == "" {
		offset = "0"
	}
	if searchQuery != "" {
		return "id >= " + offset + " " + condition + " " + searchQuery
	} else {
		return "id >= " + offset
	}

}

func TableConditionV1(offset, fields, values, condition string) string {
	basicCondition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$."
	var searchQuery string
	if fields != "" && values != "" {
		arrayOfFields := strings.Split(fields, ",")
		arrayOfValues := strings.Split(values, ",")
		if len(arrayOfFields) == len(arrayOfValues) {
			for index, field := range arrayOfFields {
				if len(arrayOfFields)-1 == index {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\""

				} else {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\" " + condition

				}
			}
		}
	}
	if offset == "" {
		offset = "0"
	}
	if searchQuery != "" {
		if offset == "-1" {
			return condition + " " + searchQuery
		} else {
			return "id > " + offset + " " + condition + " " + searchQuery
		}

	} else {
		if offset == "-1" {
			return ""
		} else {
			return "id > " + offset
		}

	}

}

func TableDecendingOrderCondition(offset, fields, values, condition string) string {
	basicCondition := "JSON_UNQUOTE(JSON_EXTRACT(object_info, \"$."
	var searchQuery string
	if fields != "" && values != "" {
		arrayOfFields := strings.Split(fields, ",")
		arrayOfValues := strings.Split(values, ",")
		if len(arrayOfFields) == len(arrayOfValues) {
			for index, field := range arrayOfFields {
				if len(arrayOfFields)-1 == index {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\""

				} else {
					searchQuery = searchQuery + basicCondition + field + "\")) = \"" + arrayOfValues[index] + "\" " + condition

				}
			}
		}
	}
	if offset == "" {
		offset = "0"
	}

	if searchQuery != "" {
		if offset == "-1" {
			return condition + " " + searchQuery
		} else {
			return "id < " + offset + " " + condition + " " + searchQuery
		}

	} else {
		if offset == "-1" {
			return ""
		} else {
			return "id < " + offset
		}

	}

}

func IsArchived(rawObjectData datatypes.JSON) bool {
	var objectFields map[string]interface{}
	json.Unmarshal(rawObjectData, &objectFields)
	if value, ok := objectFields["objectStatus"]; ok {
		var objectStatus = util.InterfaceToString(value)
		if objectStatus == "Archived" {
			return true
		}
	}
	return false
}

type ExportSchema struct {
	Label               string         `json:"label"`
	ExpandedIcon        string         `json:"expandedIcon"`
	CollapsedIcon       string         `json:"collapsedIcon"`
	Droppable           bool           `json:"droppable"`
	Data                string         `json:"data"`
	Id                  float32        `json:"id"`
	TargetTable         string         `json:"targetTable"`
	Children            []Childern     `json:"children,omitempty"`
	DataType            string         `json:"dataType"`
	LinkedMapFlag       bool           `json:"linkedMapFlag"`
	LinkedObjectMapping *ObjectMapping `json:"linkedObjectMapping,omitempty"`
	Property            string         `json:"property"`
	Type                string         `json:"type"`
}

type Childern struct {
	Label       string     `json:"label"`
	Icon        string     `json:"icon"`
	Data        string     `json:"data"`
	Id          float32    `json:"id"`
	Children    []Childern `json:"children,omitempty"`
	TargetTable string     `json:"targetTable"`
}

type ContentResponse struct {
	ContentId   string `json:"contentId"`
	ContentInfo struct {
		Url      string `json:"url"`
		Name     string `json:"name"`
		Size     string `json:"size"`
		Path     string `json:"path"`
		MetaInfo struct {
			Created   string `json:"created"`
			CreatedBy string `json:"createdBy"`
			Updated   string `json:"updated"`
			UpdatedBy string `json:"updatedBy"`
		} `json:"metaInfo"`
	} `json:"contentInfo"`
}

type ArrayOfRows struct {
	Rows []interface{} `json:"rows"`
}

type SchedulerResponse struct {
	Resources   ArrayOfRows `json:"resources"`
	Events      ArrayOfRows `json:"events"`
	Assignments ArrayOfRows `json:"assignments"`
}

type GeneralResponse struct {
	Error    uint32 `json:"error,omitempty"`
	Message  string `json:"message"`
	RecordId int    `json:"recordId"`
}

type GeneralResourceCreateResponse struct {
	Error   uint32 `json:"error,omitempty"`
	Message string `json:"message"`
	Id      int    `json:"id"`
}

func GetError(errorString string) error {
	return errors.New(errorString)
}

func TableConditionForFilter(filter string) string {

	pattern := `^(\w+)\s*(=|!=)\s*([\w\d_-]+)$`
	re := regexp.MustCompile(pattern)
	var conditions []string
	filters := strings.Split(filter, ",")

	for _, f := range filters {
		matches := re.FindStringSubmatch(strings.TrimSpace(f))
		if len(matches) != 4 {
			return ""
		}

		key, operator, value := matches[1], matches[2], matches[3]

		condition := fmt.Sprintf("object_info ->> '$.%s' %s '%s'", key, operator, value)
		conditions = append(conditions, condition)
	}

	return strings.Join(conditions, " AND ")
}

type SchedulerEventResponse struct {
	MachineId string        `json:"machineId"`
	Events    []interface{} `json:"events"`
}
