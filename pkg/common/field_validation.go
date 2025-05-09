package common

import (
	"cx-micro-flake/pkg/common/component"
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// here we define the validatitors
type Validator interface {
	Validate() (bool, error)
}

type BaseValidator struct {
	ObjectFields map[string]interface{}
	DbConnection *gorm.DB
	TargetTable  string
}

type EmptyFieldValidator struct {
	RecordSchema  *component.RecordSchema
	BaseValidator *BaseValidator
}

type DuplicateFieldValidator struct {
	RecordSchema  *component.RecordSchema
	BaseValidator *BaseValidator
}

type MandatoryFieldValidator struct {
	RecordSchema  *component.RecordSchema
	BaseValidator *BaseValidator
}
type NotificationThresholdFieldValidator struct {
	RecordSchema  *component.RecordSchema
	BaseValidator *BaseValidator
}

func (v *MandatoryFieldValidator) Validate() (bool, error) {
	if _, ok := v.BaseValidator.ObjectFields[v.RecordSchema.Property]; !ok {
		return false, errors.New("Field " + v.RecordSchema.Name + " is mandatory to create the resource. Please enter information!!")
	}
	return true, nil
}
func (v *DuplicateFieldValidator) Validate() (bool, error) {
	// get the field and check for already availability
	if val, ok := v.BaseValidator.ObjectFields[v.RecordSchema.Property]; ok {
		fieldName := v.RecordSchema.Property
		// count is it already there
		var queryResults = make(map[string]interface{})
		var query string
		value := util.InterfaceToString(val)
		var conditionIndex = 0
		var combinedCondition = ""
		fmt.Println("v.RecordSchema: ", v.RecordSchema.Property, " v.RecordSchema.UniqueIndex : ", v.RecordSchema.UniqueIndex)
		if len(v.RecordSchema.UniqueIndex) > 0 {

			for _, indexField := range v.RecordSchema.UniqueIndex {
				fmt.Println("v.indexField: ", indexField)
				if conditionIndex == 0 {
					combinedCondition = "object_info ->>'$." + indexField + "' = '"
				} else {
					combinedCondition = "object_info ->>'$." + indexField + "' = ' AND "
				}
				combinedCondition += util.InterfaceToString(v.BaseValidator.ObjectFields[indexField]) + "'"
			}

		}

		fmt.Println("v.BaseValidator.ConditionString: ", combinedCondition)
		if combinedCondition != "" {
			query = "SELECT COUNT(*) as numberOfRecords FROM " + v.BaseValidator.TargetTable + " WHERE  object_info ->>'$." + fieldName + "' =" + "\"" + value + "\" and  object_info ->>'$.objectStatus' != 'Archived' AND " + combinedCondition
		} else {
			query = "SELECT COUNT(*) as numberOfRecords FROM " + v.BaseValidator.TargetTable + " WHERE  object_info ->>'$." + fieldName + "' =" + "\"" + value + "\" and  object_info ->>'$.objectStatus' != 'Archived'"
		}
		v.BaseValidator.DbConnection.Raw(query).Scan(&queryResults)

		if queryResults["numberOfRecords"] != nil {

			numberOfRecords := util.InterfaceToInt(queryResults["numberOfRecords"])
			if numberOfRecords > 0 {
				return false, errors.New("Field " + v.RecordSchema.Name + " has already contain the same information in the system, please choose new one !!")
			}
		}

	}
	return true, nil
}
func (v *EmptyFieldValidator) Validate() (bool, error) {
	if val, ok := v.BaseValidator.ObjectFields[v.RecordSchema.Property]; ok {
		if v.RecordSchema.Type == "int" {
			if val == 0 {
				return false, errors.New("Field " + v.RecordSchema.Name + " can not be zero !!")
			}
		} else if v.RecordSchema.Type == "text" {
			if val == "" {
				return false, errors.New("Field " + v.RecordSchema.Name + " can not be empty !!")
			}
		} else if v.RecordSchema.Type == "double" {
			if val == 0 {
				return false, errors.New("Field " + v.RecordSchema.Name + " can not be zero !!")
			}
		}
	}
	return true, nil

}
func (cm *ComponentManager) DoFieldValidationOnSerializedObject(componentName string, action string, dbConnection *gorm.DB, serializedObject datatypes.JSON) error {
	// before go into any validation, check the objectStatus validation, if the object is archived
	var payload = make(map[string]interface{})
	err := json.Unmarshal(serializedObject, &payload)
	if err != nil {
		return errors.New("Invalid serialization, internal system error")
	}
	if objectStatusValue, ok := payload["objectStatus"]; ok {
		if objectStatusValue == ObjectStatusArchived {
			return errors.New("This resource is already archived, no modifications are allowed further. Archived resources only allowed to view only")
		}
	}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	targetTable := cm.ComponentSchema[int64ComponentId].TargetTable
	baseFieldValidator := &BaseValidator{}
	baseFieldValidator.ObjectFields = payload
	baseFieldValidator.DbConnection = dbConnection
	baseFieldValidator.TargetTable = targetTable

	for _, recordSchemaField := range recordSchema {
		if recordSchemaField.FieldValidator != nil {
			if action == "create" {

				if len(recordSchemaField.FieldValidator.Create) > 0 {
					var validationResult []component.ValidationResult
					for _, individualValidator := range recordSchemaField.FieldValidator.Create {
						if individualValidator.Validator == "emptyField" {
							validator := EmptyFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}
							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "duplicateField" {
							validator := DuplicateFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}

							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "mandatory" {
							validator := MandatoryFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}

							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "notificationThresholdField" {
							validator := NotificationThresholdFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}
							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						}

					}
					// now check
					for _, result := range validationResult {
						if result.ResultStatus == false {
							return result.Error
						}
					}
				}
			}
		}
	}
	return nil
}

func (cm *ComponentManager) DoFieldValidation(componentName string, action string, dbConnection *gorm.DB, payload map[string]interface{}) error {

	// before go into any validation, check the objectStatus validation, if the object is archived
	if objectStatusValue, ok := payload["objectStatus"]; ok {
		if objectStatusValue == ObjectStatusArchived {
			return errors.New("The resource is already archived, no modifications are allowed further")
		}
	}
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	recordSchema := cm.ComponentSchema[int64ComponentId].RecordSchema
	targetTable := cm.ComponentSchema[int64ComponentId].TargetTable
	baseFieldValidator := &BaseValidator{}
	baseFieldValidator.ObjectFields = payload
	baseFieldValidator.DbConnection = dbConnection
	baseFieldValidator.TargetTable = targetTable

	for _, recordSchemaField := range recordSchema {
		if recordSchemaField.FieldValidator != nil {
			if action == "create" {
				if len(recordSchemaField.FieldValidator.Create) > 0 {
					var validationResult []component.ValidationResult
					for _, individualValidator := range recordSchemaField.FieldValidator.Create {

						if individualValidator.Validator == "emptyField" {
							validator := EmptyFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}
							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "duplicateField" {
							validator := DuplicateFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}

							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "mandatory" {
							validator := MandatoryFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}

							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						} else if individualValidator.Validator == "notificationThresholdField" {
							validator := NotificationThresholdFieldValidator{
								RecordSchema:  &recordSchemaField,
								BaseValidator: baseFieldValidator,
							}
							resultStatus, err := validator.Validate()
							validationResult = append(validationResult, component.ValidationResult{
								ResultStatus: resultStatus,
								Error:        err,
							})
						}
					}
					// now check
					for _, result := range validationResult {
						if result.ResultStatus == false {
							return result.Error
						}
					}
				}
			}
		}
	}
	return nil
}

func applyFormatters(formatter string, srcValue interface{}) interface{} {
	if formatter == "string_array_to_json_array" {
		// we have string array, now convert that to json array
		var interfaceArray []string
		rawResultValue := []byte(srcValue.(string))
		json.Unmarshal(rawResultValue, &interfaceArray)
		return interfaceArray
	} else if formatter == "string_to_bool" {
		// we have string array, now convert that to json array
		stringValue := srcValue.(string)
		if stringValue == "false" {
			return false
		} else {
			return true
		}
	}
	return nil
}

func (cm *ComponentManager) TemplateMandatoryFieldValidation(dbConnection *gorm.DB, objectInfo datatypes.JSON, componentName string) error {
	int64ComponentId := cm.ComponentNameIdMapping[componentName]
	referenceTemplateComponent := cm.ComponentSchema[int64ComponentId].ReferenceTemplateComponent
	var objectFields = make(map[string]interface{})
	json.Unmarshal(objectInfo, &objectFields)
	var templateId = -1
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
				for _, templateRecord := range listOfTableRecords {
					// EnabledAfterWorkflowStatusLevel will be any set to default value 1, if the value is set more than, we don't need
					// to process
					if templateRecord.EnabledAfterWorkflowStatusLevel != nil {
						if *templateRecord.EnabledAfterWorkflowStatusLevel <= 1 {
							if templateRecord.IsMandatoryField {
								datatype := getInt2DataType(templateRecord.DataType)
								if datatype == "text" {
									if propertyValue, isProperty := objectFields[templateRecord.Property]; isProperty {
										if propertyValue == "" {
											return errors.New("system validation in place, please fill the mandatory field (marked in red * mark)")
										}
									} else {
										return errors.New("Please fill the mandatory field")
									}
								}
							}
						}
					}

				}

			}
		}
	}

	return nil
}
func (v *NotificationThresholdFieldValidator) Validate() (bool, error) {
	if val, ok := v.BaseValidator.ObjectFields[v.RecordSchema.Property]; ok {
		if v.RecordSchema.Type == "integer" {
			intVal, ok := val.(float64)
			if !ok {
				return false, errors.New("Field " + v.RecordSchema.Property + " is not a valid number")
			}

			if intVal <= 0 || intVal > 100 {
				return false, errors.New("Field " + v.RecordSchema.Property + " should be greater than 0 and less than or equal to 100!!")
			}
		}
	}

	return true, nil
}
