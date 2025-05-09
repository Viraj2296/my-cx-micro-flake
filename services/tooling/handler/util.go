package handler

import (
	"cx-micro-flake/pkg/util"
	"encoding/json"
	"gorm.io/datatypes"
)

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

func getGroupByResults(groupByColumn string, dataArray []datatypes.JSON) map[string][]interface{} {
	var groupByResults = make(map[string][]interface{})
	for _, objectInterface := range dataArray {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface, &objectFields)
		var groupByField = util.InterfaceToString(objectFields[groupByColumn])
		groupByResults[groupByField] = append(groupByResults[groupByField], objectInterface)
	}
	return groupByResults
}

func getGroupByResultsFromInterface(groupByColumn string, dataArray []interface{}) map[string][]interface{} {
	var groupByResults = make(map[string][]interface{})
	for _, objectInterface := range dataArray {
		var objectFields = make(map[string]interface{})
		json.Unmarshal(objectInterface.(datatypes.JSON), &objectFields)
		var groupByField = util.InterfaceToString(objectFields[groupByColumn])
		groupByResults[groupByField] = append(groupByResults[groupByField], objectInterface)
	}
	return groupByResults
}
