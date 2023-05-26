// Author: huaxr
// Time:   2021/7/2 上午11:29
// Git:    huaxr

package parser

import (
	"fmt"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/oliveagle/jsonpath"
	"reflect"
	"strings"
)

const (
	Dollar            = "$"
	Point             = "."
	ContextPrefix     = Dollar + Point
	LastPrefix        = Dollar + ContextPrefix
	TriggerPrefix     = Dollar + LastPrefix
	EnvPrefix         = ContextPrefix + "_env."
	TriggerContextKey = "_trigger"
)

// parse param mapper is nodeName-context like kv pair, the trigger
func parse(prefix, express string, mapper interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			logx.L().Errorf("parse panic %v", err)
		}
	}()

	// $$.a.b[0]  =>  a.b[0]
	e := express[len(prefix):]
	if len(e) == 0 {
		return mapper
	}

	switch reflect.TypeOf(mapper).Kind() {
	case reflect.Map:
		fields := strings.Split(e, Point)
		switch len(fields) {
		case 0:
			return mapper
		case 1:
			// handler exception.
			if v, ok := mapper.(map[string]interface{})[fields[0]]; ok {
				return v
			}
		default:
			if res, ok := mapper.(map[string]interface{})[fields[0]]; ok {
				if res == nil {
					logx.L().Errorf("err when get subkey:%v, cause res is nil", fields[0])
					return nil
				}
				switch reflect.TypeOf(res).Kind() {
				case reflect.Map:
					// if only add field, it just can visit deepth of 2.
					// set debug here and observer.
					express = prefix + strings.Join(fields[1:], ".")
					childMap, ok := res.(map[string]interface{})
					if !ok {
						logx.L().Errorf("mapper key:%v is not map type", fields[0])
						return nil
					}
					return parse(prefix, express, childMap)

				default:
					return res
				}
			}
		}
	default:
		return mapper
	}

	return nil
}

func lookUp(e string, mapper interface{}) (interface{}, error) {
	// a.b[0] => $.a.b[0]
	pat, err := jsonpath.Compile(ContextPrefix + e)
	if err != nil {
		logx.L().Errorf("err compile %v", err)
		return nil, err
	}
	res, err := pat.Lookup(mapper)
	if err != nil {
		return nil, fmt.Errorf("compile parse express %v fail %v", e, err)
	}
	return res, nil
}

// GetKey has a express paramters, it will be called in
// attribute parse and fetch, edge expression calling.
// express format:
// e.g. $.a.b[?(@.price <.expensive)].c  $.a.b[:].c  $.a.b.c
func GetKey(express string, last interface{}, mapper map[string]interface{}) (interface{}, error) {
	// mapper could be empty, message slim is inevitably cuts keys from mapper.
	// so tthis logic is not considerable.
	express = strings.Trim(express, " ")

	// using jsonpath parse expression
	if strings.Contains(express, "[") {
		var lookup interface{}
		var e string
		if strings.HasPrefix(express, ContextPrefix) {
			lookup = mapper
			e = express[len(ContextPrefix):]
		} else if strings.HasPrefix(express, LastPrefix) {
			lookup = last
			e = express[len(LastPrefix):]
		} else if strings.HasPrefix(express, TriggerPrefix) {
			lookup = mapper[TriggerContextKey]
			e = express[len(TriggerPrefix):]
		} else {
			return nil, fmt.Errorf("err format:%v", express)
		}
		return lookUp(e, lookup)
	}

	if strings.HasPrefix(express, ContextPrefix) {
		if strings.HasPrefix(express, EnvPrefix) {
			return parse(EnvPrefix, express, mapper), nil
		}
		return parse(ContextPrefix, express, mapper), nil
	}

	if strings.HasPrefix(express, LastPrefix) {
		return parse(LastPrefix, express, last), nil
	}
	// GetSlaveMessage input attributes using $$$. represent that
	// get the storage from the TriggerMessage context, which initialized
	// already when GetTriggerMessage getting called.
	if strings.HasPrefix(express, TriggerPrefix) {
		return parse(TriggerPrefix, express, mapper[TriggerContextKey]), nil
	}

	return nil, fmt.Errorf("err express format:%v", express)
}
