package transform

import (
	"fmt"
	"strconv"
)

// AnyToString converts anything to string
func AnyToString(input interface{}) (res string) {
	switch input := input.(type) {

	case int:
		res = strconv.Itoa(input)

	case int8:
		res = strconv.FormatInt(int64(input), 10)

	case int16:
		res = strconv.FormatInt(int64(input), 10)

	case int32:
		res = strconv.FormatInt(int64(input), 10)

	case int64:
		res = strconv.FormatInt(input, 10)

	case string:
		res = input

	case float32:
		res = strconv.FormatFloat(float64(input), 'f', -1, 32)

	case float64:
		res = strconv.FormatFloat(input, 'f', -1, 64)

	case bool:
		res = strconv.FormatBool(input)

	default:
		if input == nil {
			res = ""
		} else {
			res = fmt.Sprintf("%v", input)
		}
	}

	return
}

// AnyToSliceOfString anything to string of slices
func AnyToSliceOfString(input interface{}) (res []string) {

	switch input := input.(type) {

	case []interface{}:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []*string:
		for _, item := range input {
			val := AnyToString(*item)
			res = append(res, val)
		}

	case []string:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []int:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []int8:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []int16:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []int32:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []int64:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []float32:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []float64:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	case []bool:
		for _, item := range input {
			val := AnyToString(item)
			res = append(res, val)
		}

	default:
		//fmt.Printf("AnyToSliceOfString, unknown type: %T | %+v\n", val, input)
	}

	return
}
