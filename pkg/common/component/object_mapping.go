package component

type FunctionArguments struct {
	Field string  `json:"field"`
	Value *string `json:"value"`
}
type Function struct {
	Name      string              `json:"name"`
	Arguments []FunctionArguments `json:"arguments"`
}
type Service struct {
	Name         string `json:"name"`
	Call         string `json:"call"`
	ServiceParam []struct {
		Type  string `json:"type"`
		Field string `json:"field"`
	} `json:"serviceParam"`
	CommonRouteLink string `json:"commonRouteLink"`
}
type Catalyst struct {
	Name         string `json:"name"`
	Call         string `json:"call"`
	ServiceParam []struct {
		Type  string `json:"type"`
		Field string `json:"field"`
	} `json:"serviceParam"`
}
type ObjectMapping struct {
	Service    *Service          `json:"service"`
	Call       *string           `json:"call"`
	Catalyst   []Catalyst        `json:"catalyst"`
	Name       string            `json:"name"`
	Query      *QueryObject      `json:"query"`
	Function   *Function         `json:"function"`
	Predefined *PredefinedObject `json:"predefined"`
	Builder    struct {
		SingleDropdown           *SingleDropdownBuilder           `json:"singleDropdown"`
		SingleValueDropdown      *SingleValueDropdownBuilder      `json:"singleValueDropdown"`
		MultiValueDropdown       *MultiValueDropdownBuilder       `json:"multiValueDropdown"`
		FieldPropertyAssignment  *FieldPropertyAssignmentBuilder  `json:"fieldPropertyAssignment"`
		ObjectFieldAssignment    *ObjectFieldAssignmentBuilder    `json:"ObjectFieldAssignment"`
		MultiDropdownObject      *MultiDropdownObjectBuilder      `json:"multiDropdownObject"`
		SingleDropdownObject     *SingleDropdownObjectBuilder     `json:"singleDropdownObject"`
		GroupBy                  *GroupByBuilder                  `json:"groupBy"`
		SingleObjectValueArray   *SingleObjectValueArrayBuilder   `json:"singleObjectValueArray"`
		ObjectArray              *ObjectArrayBuilder              `json:"objectArray"`
		KeyValue                 *KeyValueBuilder                 `json:"keyValue"`
		Table                    *TableBuilder                    `json:"table"`
		SingleValue              *SingleValueBuilder              `json:"singleValue"`
		SingleValueFromId        *SingleValueFromIdBuilder        `json:"singleValueFromId"`
		IdValue                  *IdValueBuilder                  `json:"idValue"`
		Email                    *EmailBuilder                    `json:"email"`
		Date                     *DateBuilder                     `json:"date"`
		SingleValueCondition     *SingleValueConditionBuilder     `json:"singleValueCondition"`
		TableFieldsToObject      *TableFieldsToObjectBuilder      `json:"tableFieldsToObject"`
		TableFieldsToObjectArray *TableFieldsToObjectArrayBuilder `json:"tableFieldsToObjectArray"`
		DBExecution              *DbExecutionBuilder              `json:"dbExecution"`
		EmptyValueValidator      *EmptyValueValidator             `json:"emptyValueValidator"`
	} `json:"builder"`
}

type ReplacementFields struct {
	Type          string         `json:"type"`
	Field         string         `json:"field"`
	Format        string         `json:"format"`
	Property      string         `json:"property"`
	ObjectMapping *ObjectMapping `json:"objectMapping"`
	Duration      int            `json:"duration"`
	Resolution    string         `json:"resolution"`
}

type QueryOutFieldsMapping struct {
	Field     string `json:"field"`
	Formatter string `json:"formatter"`
}
type QueryObject struct {
	Query             string                  `json:"query"`
	ReplacementFields []ReplacementFields     `json:"replacementFields"`
	OutFieldsMapping  []QueryOutFieldsMapping `json:"outFieldsMapping"`
}

type PredefinedObject struct {
	Data []map[string]interface{} `json:"data"`
}

type SingleDropdownBuilder struct {
	IsReverseMapping bool   `json:""`
	Index            string `json:"index"`
	Value            string `json:"value"`
	ReferenceField   string `json:"referenceField"`
}

type SingleValueDropdownBuilder struct {
	Field string `json:"field"`
}

type MultiDropdownObjectBuilder struct {
	Index  string   `json:"index"`
	Values []string `json:"values"`
}
type SingleDropdownObjectBuilder struct {
	Index  string   `json:"index"`
	Values []string `json:"values"`
}

type MultiValueDropdownBuilder struct {
	Index string `json:"index"`
	Value string `json:"value"`
}
type GroupByBuilder struct {
	Id                string `json:"id"`
	Value             string `json:"value"`
	GroupByColumnName string `json:"groupByColumnName"`
}

type SingleValueConditionBuilder struct {
	Value  int `json:"value"`
	Action struct {
		Type              string `json:"type"`
		Query             string `json:"query"`
		IdColumn          string `json:"idColumn"`
		ReplacementFields []struct {
			Type     string `json:"type"`
			Field    string `json:"field"`
			Format   string `json:"format"`
			Property string `json:"property"`
		} `json:"replacementFields"`
	} `json:"action"`
	Condition string `json:"condition"`
	Field     string `json:"field"`
}

type EmptyValueValidator struct {
}

type ObjectArrayBuilder struct {
	Values []string `json:"values"`
}

type SingleObjectValueArrayBuilder struct {
	Field string `json:"field"`
}

type DbExecutionBuilder struct {
}

type TableFieldsToObjectBuilder struct {
	Fields []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"fields"`
}

type TableFieldsToObjectArrayBuilder struct {
	Fields []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"fields"`
}

type KeyValueBuilder struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TableBuilder struct {
	Schema          []TableSchema `json:"schema"`
	CommonRouteLink string        `json:"commonRouteLink"`
}

type SingleValueBuilder struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type SingleValueFromIdBuilder struct {
	Field   string `json:"field"`
	Type    string `json:"type"`
	IdField string `json:"idField"`
}

type FieldPropertyAssignmentBuilder struct {
	Field    string `json:"field"`
	Property string `json:"property"`
}
type ObjectFieldAssignmentBuilder struct {
	Property string `json:"property"`
}
type IdValueBuilder struct {
	Id    string `json:"id"`
	Field string `json:"field"`
}

type EmailBuilder struct {
	ReplacementFields []ReplacementFields `json:"replacementFields"`
	Subject           string              `json:"subject"`
}

type DateBuilder struct {
	Function string `json:"function"`
}
