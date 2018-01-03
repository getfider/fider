package dbx

import (
	"fmt"
	"reflect"

	"github.com/lib/pq"
)

type RowMapper struct {
	cache map[reflect.Type]TypeMapper
}

func NewRowMapper() RowMapper {
	return RowMapper{
		cache: make(map[reflect.Type]TypeMapper, 0),
	}
}

func (m *RowMapper) Map(dest interface{}, columns []string, scanner func(dest ...interface{}) error) error {
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var (
		typeMapper TypeMapper
		ok         bool
	)

	if typeMapper, ok = m.cache[t]; !ok {
		typeMapper = NewTypeMapper(t)
		m.cache[t] = typeMapper
	}

	vd := reflect.ValueOf(dest).Elem()
	pointers := make([]interface{}, len(columns))
	for i, c := range columns {
		mapping := typeMapper.Fields[c]
		field := vd

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

		if field.Kind() == reflect.Slice {
			obj := reflect.New(reflect.MakeSlice(field.Type(), 0, 0).Type()).Elem()
			field.Set(obj)
			pointers[i] = pq.Array(field.Addr().Interface())
		} else {
			pointers[i] = field.Addr().Interface()
		}
	}

	return scanner(pointers...)
}

type TypeMapper struct {
	Type   reflect.Type
	Fields map[string]FieldInfo
}

func NewTypeMapper(t reflect.Type) TypeMapper {
	all := make(map[string]FieldInfo, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		if columnName != "" {
			fieldType := field.Type
			fieldKind := fieldType.Kind()

			if fieldKind == reflect.Ptr {
				fieldType = field.Type.Elem()
				fieldKind = fieldType.Kind()
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

type FieldInfo struct {
	FieldName  []string
	ColumnName string
}
