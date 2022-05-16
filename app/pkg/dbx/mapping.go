package dbx

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lib/pq"
)

//RowMapper is responsible for mapping a sql.Rows into a Struct (model)
type RowMapper struct {
	cache map[reflect.Type]TypeMapper
	sync.RWMutex
}

//NewRowMapper creates a new instance of RowMapper
func NewRowMapper() *RowMapper {
	return &RowMapper{
		cache: make(map[reflect.Type]TypeMapper),
	}
}

//Map values from scanner (usually sql.Rows.Scan) into dest based on columns
func (m *RowMapper) Map(dest any, columns []string, scanner func(dest ...any) error) error {
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var (
		typeMapper TypeMapper
		ok         bool
	)

	m.RLock()
	typeMapper, ok = m.cache[t]
	m.RUnlock()

	if !ok {
		typeMapper = NewTypeMapper(t)
		m.Lock()
		m.cache[t] = typeMapper
		m.Unlock()
	}

	pointers := make([]any, len(columns))
	for i, c := range columns {
		mapping := typeMapper.Fields[c]
		field := reflect.ValueOf(dest).Elem()

		for _, f := range mapping.FieldName {
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				field = field.Elem()
			}
			field = field.FieldByName(f)
		}

		if !field.CanAddr() {
			panic(fmt.Sprintf("Field not found for column %s", c))
		}

		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() != reflect.Uint8 {
			obj := reflect.New(reflect.MakeSlice(field.Type(), 0, 0).Type()).Elem()
			field.Set(obj)
			pointers[i] = pq.Array(field.Addr().Interface())
		} else {
			pointers[i] = field.Addr().Interface()
		}
	}

	return scanner(pointers...)
}

//TypeMapper holds information about how to map SQL ResultSet to a Struct
type TypeMapper struct {
	Type   reflect.Type
	Fields map[string]FieldInfo
}

//NewTypeMapper creates a new instance of TypeMapper for given reflect.Type
func NewTypeMapper(t reflect.Type) TypeMapper {
	all := make(map[string]FieldInfo)

	if t.Kind() != reflect.Struct {
		return TypeMapper{
			Type:   t,
			Fields: all,
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		if columnName != "" {
			fieldType := field.Type
			fieldKind := fieldType.Kind()

			if fieldKind == reflect.Ptr {
				fieldType = field.Type.Elem()
				mapper := NewTypeMapper(fieldType)
				for _, f := range mapper.Fields {
					all[columnName+"_"+f.ColumnName] = FieldInfo{
						ColumnName: columnName + "_" + f.ColumnName,
						FieldName:  append([]string{field.Name}, f.FieldName...),
					}
				}
			} else {
				all[columnName] = FieldInfo{
					FieldName:  []string{field.Name},
					ColumnName: columnName,
				}
			}
		}
	}
	return TypeMapper{
		Type:   t,
		Fields: all,
	}
}

//FieldInfo is a simple struct to map Column -> Field
type FieldInfo struct {
	FieldName  []string
	ColumnName string
}
