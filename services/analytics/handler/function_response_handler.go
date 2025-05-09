package handler

import (
	"cx-micro-flake/pkg/common/component"
)

func (as *AnalyticsService) getAdditionalRecordFunctionResponse(componentName string, objectFields map[string]interface{}, generalResponse map[string]interface{}) {
	int64ComponentId := as.ComponentManager.ComponentNameIdMapping[componentName]
	additionalRecords := as.ComponentManager.ComponentSchema[int64ComponentId].AdditionalRecords

	if len(additionalRecords) > 0 {
		// we have configured additional record schema
		for _, individualRecordSchema := range additionalRecords {
			recordInfo := component.RecordInfo{}
			if individualRecordSchema.ObjectMapping.Function != nil {

				if individualRecordSchema.ObjectMapping.Function.Name == FunctionGetExistingDatabaseTables {
					functionResponse := getExistingDatabaseTables(as, objectFields, individualRecordSchema.ObjectMapping.Function.Arguments)
					recordInfo.Data = functionResponse
					recordInfo.IsExternal = true
				}

				generalResponse[individualRecordSchema.Property] = recordInfo

			}
		}

	}
}
