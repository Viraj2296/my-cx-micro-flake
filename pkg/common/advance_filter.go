package common

import (
	"cx-micro-flake/pkg/util"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type FilterValue struct {
	ID    int         `json:"id"`
	Value interface{} `json:"value"` // Interface{} for flexibility
}

type Filter struct {
	Property       string        `json:"property"`
	EvaluationType string        `json:"evaluationType"`
	InterfaceType  string        `json:"interfaceType"`
	Id             int           `json:"id"`
	Values         []FilterValue `json:"values"`
}

type FilterInfo struct {
	ComponentName     string   `json:"componentName"`
	CreatedBy         int      `json:"createdBy"`
	FilterName        string   `json:"filterName"`
	FilterDescription string   `json:"filterDescription"`
	Type              string   `json:"type"`
	Condition         string   `json:"condition"`
	Filters           []Filter `json:"filters"`
}

type FilterCriteria struct {
	ComponentName     string         `json:"componentName"`
	CreatedBy         int            `json:"createdBy"`
	FilterName        string         `json:"filterName"`
	FilterDescription string         `json:"filterDescription"`
	Type              string         `json:"type"`
	Condition         string         `json:"condition"`
	Filters           []FilterOption `json:"filters"`
}

type FilterOption struct {
	Property       string      `json:"property"`
	EvaluationType string      `json:"evaluationType"`
	InterfaceType  string      `json:"interfaceType"`
	Value          interface{} `json:"value"`
}

func (cm *ComponentManager) ExecuteFilter(dbConnection *gorm.DB, componentName string, searchFields FilterCriteria, additionalClause string, limitValue string) string {
	var filterCriteria = searchFields.Filters

	selectQuery := componentName + ".id "
	var componentWhereClause = " where "

	for index, filterOption := range filterCriteria {
		if len(filterCriteria)-1 == index {
			switch filterOption.InterfaceType {
			case "input_text":
				if filterOption.EvaluationType == "equal" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' "
				} else if filterOption.EvaluationType == "contains" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' like '%" + searchValue + "%' "
				}
			case "input_number":
				var searchValue = util.InterfaceToInt(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "contains" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' like '%" + strconv.Itoa(searchValue) + "%' "
				}
			case "date_input":
				var searchValue = util.InterfaceToString(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > '" + searchValue + "' "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < '" + searchValue + "' "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != '" + searchValue + "' "
				}
			case "datetime_input":
				var searchValueStr = util.InterfaceToString(filterOption.Value)
				var searchValue = util.ConvertSingaporeTimeToUTC(searchValueStr)
				if searchValue != "" {
					if filterOption.EvaluationType == "equal" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' "
					} else if filterOption.EvaluationType == "greater_than" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > '" + searchValue + "' "
					} else if filterOption.EvaluationType == "lesser_than" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < '" + searchValue + "' "
					} else if filterOption.EvaluationType == "not_equal" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != '" + searchValue + "' "
					}
				}

			}
		} else {
			switch filterOption.InterfaceType {
			case "input_text":
				if filterOption.EvaluationType == "equal" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "contains" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' like '%" + searchValue + "%' " + searchFields.Condition + " "
				}
			case "input_number":
				var searchValue = util.InterfaceToInt(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "contains" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' like '%" + strconv.Itoa(searchValue) + "%' " + searchFields.Condition + " "
				}
			case "date_input":
				var searchValue = util.InterfaceToString(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > '" + searchValue + "' " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < '" + searchValue + "' " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != '" + searchValue + "' " + searchFields.Condition + " "
				}
			case "datetime_input":
				var searchValueStr = util.InterfaceToString(filterOption.Value)
				var searchValue = util.ConvertSingaporeTimeToUTC(searchValueStr)

				if searchValue != "" {
					if filterOption.EvaluationType == "equal" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' " + searchFields.Condition + " "
					} else if filterOption.EvaluationType == "greater_than" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > '" + searchValue + "' " + searchFields.Condition + " "
					} else if filterOption.EvaluationType == "lesser_than" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < '" + searchValue + "' " + searchFields.Condition + " "
					} else if filterOption.EvaluationType == "not_equal" {
						componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != '" + searchValue + "' " + searchFields.Condition + " "
					}
				}

			}
		}

	}
	searchQuery := "select " + selectQuery + " from " + componentName + util.TrimSuffix(componentWhereClause, "or (") + " and " + componentName + ".object_info ->> '$.objectStatus' = 'Active' and " + additionalClause + " limit " + limitValue

	fmt.Println("=============================================")
	fmt.Println("searchQuery: ", searchQuery)
	var selectedIds []map[string]interface{}
	dbConnection.Raw(searchQuery).Scan(&selectedIds)

	idList := "("
	for _, idObjects := range selectedIds {
		id := util.InterfaceToString(idObjects["id"])
		idList += id + ","
	}
	fmt.Println("idList: ", idList)
	idList = util.TrimSuffix(idList, ",") + ")"
	return idList

}

func (cm *ComponentManager) GetFilterCount(dbConnection *gorm.DB, componentName string, searchFields FilterCriteria, additionalClause string) int64 {
	var filterCriteria = searchFields.Filters

	selectQuery := componentName + ".id "
	var componentWhereClause = " where "

	for index, filterOption := range filterCriteria {
		if len(filterCriteria)-1 == index {
			switch filterOption.InterfaceType {
			case "input_text":
				if filterOption.EvaluationType == "equal" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' "
				}
			case "input_number":
				var searchValue = util.InterfaceToInt(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < " + strconv.Itoa(searchValue) + " "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != " + strconv.Itoa(searchValue) + " "
				}
			}
		} else {
			switch filterOption.InterfaceType {
			case "input_text":
				if filterOption.EvaluationType == "equal" {
					var searchValue = util.InterfaceToString(filterOption.Value)
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = '" + searchValue + "' " + searchFields.Condition + " "
				}
			case "input_number":
				var searchValue = util.InterfaceToInt(filterOption.Value)
				if filterOption.EvaluationType == "equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' = " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "greater_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' > " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "lesser_than" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' < " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				} else if filterOption.EvaluationType == "not_equal" {
					componentWhereClause = componentWhereClause + componentName + ".object_info ->> '$." + filterOption.Property + "' != " + strconv.Itoa(searchValue) + " " + searchFields.Condition + " "
				}
			}
		}

	}

	var searchQuery string
	if additionalClause == "" {
		searchQuery = "select " + selectQuery + " from " + componentName + util.TrimSuffix(componentWhereClause, "or (") + " and " + componentName + ".object_info ->> '$.objectStatus' = 'Active'"
	} else {
		searchQuery = "select " + selectQuery + " from " + componentName + util.TrimSuffix(componentWhereClause, "or (") + " and " + componentName + ".object_info ->> '$.objectStatus' = 'Active' and " + additionalClause
	}

	fmt.Println("=============================================")
	fmt.Println("searchQuery: ", searchQuery)
	var selectedIds = make([]map[string]interface{}, 0)
	dbConnection.Raw(searchQuery).Scan(&selectedIds)

	return int64(len(selectedIds))

}
