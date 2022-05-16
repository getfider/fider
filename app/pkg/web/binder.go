package web

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	//ErrContentTypeNotAllowed is used when POSTing a body that is not json
	ErrContentTypeNotAllowed = stdErrors.New("Only Content-Type application/json is allowed")
)

//DefaultBinder is the default HTTP binder
type DefaultBinder struct {
}

//NewDefaultBinder creates a new default binder
func NewDefaultBinder() *DefaultBinder {
	return &DefaultBinder{}
}

func methodHasBody(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodDelete ||
		method == http.MethodPut
}

//Bind request data to object i
func (b *DefaultBinder) Bind(target any, c *Context) error {
	if methodHasBody(c.Request.Method) && c.Request.ContentLength > 0 {
		contentType := strings.Split(c.Request.GetHeader("Content-Type"), ";")
		if len(contentType) == 0 || contentType[0] != JSONContentType {
			return ErrContentTypeNotAllowed
		}

		if err := json.Unmarshal([]byte(c.Request.Body), target); err != nil {
			return err
		}
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()
	for i := 0; i < targetValue.NumField(); i++ {
		if err := b.bindRoute(i, targetValue, targetType, c.params); err != nil {
			return err
		}
		b.format(i, targetValue, targetType)
	}
	return nil
}

func (b *DefaultBinder) bindRoute(idx int, target reflect.Value, targetType reflect.Type, params StringMap) error {
	name := targetType.Field(idx).Tag.Get("route")
	if name != "" {
		field := target.Field(idx)
		fieldTypeKind := field.Type().Kind()

		if isInt(fieldTypeKind) {
			value, err := strconv.ParseInt(params[name], 10, 64)
			if err != nil {
				obj := reflect.New(field.Type())
				m := obj.MethodByName("UnmarshalText")
				if m.IsValid() {
					b := []byte(params[name])
					r := m.Call([]reflect.Value{reflect.ValueOf(b)})
					if r[0].IsNil() {
						field.Set(obj.Elem())
					}
				}
			} else {
				field.SetInt(value)
			}
		} else if isString(fieldTypeKind) {
			field.SetString(params[name])
		}
	}

	return nil
}

func (b *DefaultBinder) format(idx int, target reflect.Value, targetType reflect.Type) {
	field := target.Field(idx)
	fieldType := field.Type()
	fieldTypeKind := fieldType.Kind()
	format := targetType.Field(idx).Tag.Get("format")

	if isString(fieldTypeKind) {
		value := field.Interface().(string)
		field.SetString(applyFormat(format, value))
	} else if fieldTypeKind == reflect.Slice && isString(fieldType.Elem().Kind()) {
		values := field.Interface().([]string)
		for i, value := range values {
			field.Index(i).SetString(applyFormat(format, value))
		}
	}
}

func isInt(k reflect.Kind) bool {
	return k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64
}

func isString(k reflect.Kind) bool {
	return k == reflect.String
}

func applyFormat(format string, value string) string {
	value = strings.TrimSpace(value)
	if format == "lower" {
		value = strings.ToLower(value)
	} else if format == "upper" {
		value = strings.ToUpper(value)
	}
	return value
}
