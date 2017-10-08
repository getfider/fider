package web

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	intType    = reflect.TypeOf(0)
	stringType = reflect.TypeOf("")
)

//DefaultBinder is the default HTTP binder
type DefaultBinder struct {
}

//NewDefaultBinder creates a new default binder
func NewDefaultBinder() *DefaultBinder {
	return &DefaultBinder{}
}

//Bind request data to object i
func (b *DefaultBinder) Bind(target interface{}, c *Context) error {
	if c.Request.Method == http.MethodPost {
		if err := json.NewDecoder(c.Request.Body).Decode(target); err != nil {
			return err
		}
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()
	for i := 0; i < targetValue.NumField(); i++ {
		b.bindRoute(i, targetValue, targetType, c.params)
		b.format(i, targetValue, targetType)
	}
	return nil
}

func (b *DefaultBinder) bindRoute(idx int, target reflect.Value, targetType reflect.Type, params StringMap) error {
	name := targetType.Field(idx).Tag.Get("route")
	if name != "" {
		field := target.Field(idx)
		fieldType := field.Type()
		if fieldType == intType {
			value, _ := strconv.ParseInt(params[name], 10, 64)
			field.SetInt(value)
		} else if fieldType == stringType {
			field.SetString(params[name])
		}
	}

	return nil
}

func (b *DefaultBinder) format(idx int, target reflect.Value, targetType reflect.Type) {
	field := target.Field(idx)

	if field.Type() != stringType {
		return
	}

	format := targetType.Field(idx).Tag.Get("format")
	str := field.Interface().(string)
	str = strings.TrimSpace(str)
	if format == "lower" {
		str = strings.ToLower(str)
	}
	field.SetString(str)
}
