// Author: huaxr
// Time:   2021/12/27 上午11:53
// Git:    huaxr

package valuate

import (
	"errors"
	"reflect"

	"github.com/Knetic/govaluate"
	"github.com/spf13/cast"
)

var functions = map[string]govaluate.ExpressionFunction{
	"equals": equals,
	"len":    length,
	"gte":    gte,
	"lte":    lte,
	"gt":     gt,
	"lt":     lt,
}

func equals(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("must have 2 args")
	}
	return cast.ToString(args[0]) == cast.ToString(args[1]), nil
}

func length(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("must have 1 args")
	}

	val := args[0]
	switch reflect.TypeOf(val).Kind() {
	case reflect.Array, reflect.Slice:
		return len(val.([]interface{})), nil
	case reflect.Map:
		return len(val.(map[string]interface{})), nil
	case reflect.String:
		return len(val.(string)), nil
	default:
		return nil, errors.New("not support parameter, only allowed map,slice,string")
	}
}

func gte(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("must have 2 args")
	}
	v1 := cast.ToFloat64(cast.ToString(args[0]))
	v2 := cast.ToFloat64(cast.ToString(args[1]))
	return v1 >= v2, nil
}

func gt(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("must have 2 args")
	}
	v1 := cast.ToFloat64(cast.ToString(args[0]))
	v2 := cast.ToFloat64(cast.ToString(args[1]))
	return v1 > v2, nil
}

func lte(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("must have 2 args")
	}
	v1 := cast.ToFloat64(cast.ToString(args[0]))
	v2 := cast.ToFloat64(cast.ToString(args[1]))
	return v1 <= v2, nil
}

func lt(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("must have 2 args")
	}
	v1 := cast.ToFloat64(cast.ToString(args[0]))
	v2 := cast.ToFloat64(cast.ToString(args[1]))
	return v1 < v2, nil
}
