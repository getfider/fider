package web

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

//DefaultBinder is the default HTTP binder
type DefaultBinder struct {
}

//NewDefaultBinder creates a new default binder
func NewDefaultBinder() *DefaultBinder {
	return &DefaultBinder{}
}

//Bind request data to object i
func (b *DefaultBinder) Bind(i interface{}, c *Context) error {
	req := c.Request()
	if req.Method == http.MethodPost {
		if err := json.NewDecoder(req.Body).Decode(i); err != nil {
			return err
		}
		b.format(i)
	}
	return nil
}

func (b *DefaultBinder) format(i interface{}) {
	value := reflect.ValueOf(i).Elem()
	stringType := reflect.TypeOf("")

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		// Ignore fields that don't have the same type as a string
		if field.Type() != stringType {
			continue
		}

		format := value.Type().Field(i).Tag.Get("format")
		str := field.Interface().(string)
		str = strings.TrimSpace(str)
		if format == "lower" {
			str = strings.ToLower(str)
		}
		field.SetString(str)
	}
}
