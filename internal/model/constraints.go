package model

import (
	"fmt"
	"reflect"
	"strconv"
)

type Filter struct {
	Should []FilterParam `json:"should,omitempty"`
	Must   []FilterParam `json:"must,omitempty"`
}

type FilterParam struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
	Type  string `json:"type"`
}

type ConditionStatement struct {
	Main   string        `json:"main"`
	Values []interface{} `json:"condition.Values"`
}

func Constraints(request Request) (where ConditionStatement, or ConditionStatement, err error) {
	//err is unused for now, add if needed
	where = makeCondition(request.Metadata.Filter.Must)
	or = makeCondition(request.Metadata.Filter.Should)
	return
}

func makeCondition(filters []FilterParam) (condition ConditionStatement) {

	condition.Values = make([]interface{}, 0)
	condition.Main = ""

	for i, filter := range filters {
		condition.Main += filter.Key + " "
		switch filter.Type {
		case "eq":
			condition.Main += "= "
		case "neq":
			condition.Main += "<> "
		case "lt":
			condition.Main += "< "
		case "gt":
			condition.Main += "> "
		case "lte":
			condition.Main += "<= "
		case "gte":
			condition.Main += ">= "
		case "in":
			condition.Main += "IN "
		case "like":
			condition.Main += "LIKE "
		}

		condition.Main += "?"
		if i < len(filters)-1 {
			condition.Main += " AND "
		}

		switch filter.Value.(type) {
		case string:
			value := filter.Value.(string)
			condition.Values = append(condition.Values, value)
		case []string:
			value := filter.Value.([]string)
			condition.Values = append(condition.Values, value)
		case int:
			value := strconv.Itoa(filter.Value.(int))
			condition.Values = append(condition.Values, value)
		case float64:
			value := fmt.Sprintf("%.0f", filter.Value)
			condition.Values = append(condition.Values, value)
		default:
			fmt.Println("not found", reflect.TypeOf(filter.Value))
		}
	}

	return
}
