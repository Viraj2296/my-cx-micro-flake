package common

import (
	"cx-micro-flake/pkg/common/component"
	"encoding/json"
	"gorm.io/datatypes"
	"time"
)

type TemplateRecords struct {
	Id                              int                                    `json:"id"`
	Label                           string                                 `json:"label"`
	Description                     string                                 `json:"description"`
	Unit                            string                                 `json:"unit"`
	Property                        string                                 `json:"property"`
	CreatedAt                       time.Time                              `json:"createdAt"`
	GridSystem                      int                                    `json:"gridSystem"`
	InterfaceType                   int                                    `json:"interfaceType"`
	InterfaceTypeList               int                                    `json:"interfaceTypeList"`
	InterfaceFieldList              int                                    `json:"interfaceFieldList"`
	LastUpdatedAt                   time.Time                              `json:"lastUpdatedAt"`
	DataType                        int                                    `json:"dataType"`
	DynamicDroppedDownAttributes    component.DynamicDroppedDownAttributes `json:"dynamicDroppedDownAttributes"`
	IsDroppedDownConditionalField   bool                                   `json:"isDroppedDownConditionalField"`
	IsMandatoryField                bool                                   `json:"isMandatoryField"`
	EnabledAfterWorkflowStatusLevel *int                                   `json:"enabledAfterWorkflowStatusLevel"`
}

const (
	SingleDropdown   = 1
	NumberInputField = 2
	TextInputField   = 3
	DateInputField   = 4
	DroppedDown      = 5
	CheckBox         = 6
	ArrayInputField  = 7

	DataTypeString  = 1
	DataTypeInteger = 2
	DataTypeDouble  = 3
	DataTypeDate    = 4
	DataTypeSingle  = 5

	FieldTypeNumber   = 2
	FieldTypeText     = 3
	FieldTypeDate     = 4
	FieldTypeDropDown = 5
	FieldTypeCheckBox = 6
	FieldTypeFile     = 7

	Col12MdCol4LgCol4SmCol4   = 1
	Col12MdCol4LgCol3SmCol4   = 2
	Col12MdCol12LgCol12SmCol4 = 3
	Col12MdCol6LgCol6mCol4    = 4

	TemplateIdQuery = "templateId"
)

func getInt2DataType(dataType int) string {
	if dataType == DataTypeString {
		return "text"
	}
	if dataType == DataTypeInteger {
		return "int"
	}
	if dataType == DataTypeDouble {
		return "double"
	}
	if dataType == DataTypeDate {
		return "date"
	}

	return "text"

}

func getInt2FieldType(dataType int) string {
	if dataType == FieldTypeDropDown {
		return "DropDown"
	}
	if dataType == FieldTypeNumber {
		return "Number"
	}
	if dataType == FieldTypeText {
		return "Text"
	}
	if dataType == FieldTypeDate {
		return "Date"
	}
	if dataType == FieldTypeCheckBox {
		return "CheckBox"
	}
	if dataType == FieldTypeFile {
		return "File"
	}

	return "Text"

}

func getGridSystem2Str(interfaceType int) string {
	if interfaceType == Col12MdCol4LgCol4SmCol4 {
		return "col-12 md:col-4 lg:col-4 sm:col-12"
	}
	if interfaceType == Col12MdCol4LgCol3SmCol4 {
		return "col-12 md:col-4 lg:col-3 sm:col-12"
	}
	if interfaceType == Col12MdCol12LgCol12SmCol4 {
		return "col-12 md:col-12 lg:col-12 sm:col-12"
	}
	if interfaceType == Col12MdCol6LgCol6mCol4 {
		return "col-12 md:col-6 lg:col-6 sm:col-12"
	}

	return "col-12 md:col-4 lg:col-4 sm:col-12"
}

func getInterfaceType2Str(interfaceType int) string {
	if interfaceType == SingleDropdown {
		return "singleDropdown"
	}
	if interfaceType == NumberInputField {
		return "numberInputField"
	}
	if interfaceType == TextInputField {
		return "textInputField"
	}
	if interfaceType == DateInputField {
		return "dateInputField"
	}
	if interfaceType == DroppedDown {
		return "droppedDownField"
	}
	if interfaceType == CheckBox {
		return "checkBoxField"
	}
	if interfaceType == ArrayInputField {
		return "arrayField"
	}

	return "textInputField"
}

func (cm *ComponentManager) HandleNewDynamicRecords(componentName string, level1Index string, responseMappingField string, generalResponse []component.GeneralObject, templateTableFieldsObject component.GeneralObject) map[string]interface{} {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	componentSchema := cm.ComponentSchema[int64ComponentId]
	var responseMapping = make(map[string]interface{})

	//level1Index := header_parser.GetQueryField(ctx, "level1") // level 1 field indicates, we needed to get the records with the condition
	//templateId := header_parser.GetQueryField(ctx, "templateId")
	//responseMappingField := header_parser.GetQueryField(ctx, "responseMappingField")
	tableObjectResponse := component.TableObjectResponse{}
	if level1Index != "" {

		templateFieldName := componentSchema.ReferenceTemplateComponent.TemplateFieldName
		var dynamicTableHeader []component.TableSchema
		if templateFieldName != "" {
			// use that to extract the header
			var templateFields = make(map[string]interface{})
			json.Unmarshal(templateTableFieldsObject.ObjectInfo, &templateFields)

			templateFieldRecords := templateFields[templateFieldName]
			serialisedRecords, _ := json.Marshal(templateFieldRecords)
			var listOfTableRecords []TemplateRecords
			json.Unmarshal(serialisedRecords, &listOfTableRecords)
			// we got the fields

			for _, templateRecords := range listOfTableRecords {
				dynamicTableHeader = append(dynamicTableHeader, component.TableSchema{
					Name:          templateRecords.Label,
					Type:          getInt2DataType(templateRecords.DataType),
					Display:       true,
					Property:      templateRecords.Property,
					Unit:          templateRecords.Unit,
					RouteEnabled:  false,
					InterfaceType: getInterfaceType2Str(templateRecords.InterfaceTypeList),
					Render:        true,
					GridSystem:    getGridSystem2Str(templateRecords.GridSystem),
					Label:         templateRecords.Label,
				})
			}

		}
		tableObjectResponse.Header = dynamicTableHeader
		var arrayObjectData []datatypes.JSON
		for _, generalObject := range generalResponse {
			arrayObjectData = append(arrayObjectData, generalObject.Serialised())
		}

	} else {
		// level is not given, we needed to looking

		templateFieldName := componentSchema.ReferenceTemplateComponent.TemplateFieldName
		var dynamicTableHeader []component.TableSchema
		if templateFieldName != "" {
			// use that to extract the header

			var templateFields = make(map[string]interface{})
			json.Unmarshal(templateTableFieldsObject.ObjectInfo, &templateFields)

			templateFieldRecords := templateFields[templateFieldName]
			serialisedRecords, _ := json.Marshal(templateFieldRecords)
			var listOfTableRecords []TemplateRecords
			json.Unmarshal(serialisedRecords, &listOfTableRecords)
			// we got the fields

			for _, templateRecords := range listOfTableRecords {
				dynamicTableHeader = append(dynamicTableHeader, component.TableSchema{
					Name:          templateRecords.Label,
					Type:          getInt2DataType(templateRecords.DataType),
					Display:       true,
					Property:      templateRecords.Property,
					RouteEnabled:  false,
					InterfaceType: getInterfaceType2Str(templateRecords.InterfaceTypeList),
					Render:        true,
					GridSystem:    getGridSystem2Str(templateRecords.GridSystem),
					Label:         templateRecords.Label,
				})
			}

		}
		tableObjectResponse.Header = dynamicTableHeader
		var arrayObjectData []datatypes.JSON
		for _, generalObject := range generalResponse {
			arrayObjectData = append(arrayObjectData, generalObject.Serialised())
		}
		tableObjectResponse.Data = arrayObjectData
	}

	if responseMappingField != "" {
		// responseMapping Field is parsed, we should send the response under this field

		responseMapping[responseMappingField] = tableObjectResponse
		return responseMapping
	}

	return nil

}
