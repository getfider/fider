package dbx

import "reflect"
import "fmt"

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
		mapping := typeMapper.Mapping[c]
		field := vd.FieldByName(mapping.FieldName)
		if !field.CanAddr() {
			panic(fmt.Sprintf("Field not found for column %s", c))
		}
		pointers[i] = field.Addr().Interface()
	}

	return scanner(pointers...)
}

type TypeMapper struct {
	Type    reflect.Type
	Mapping map[string]MappingInfo
}

func NewTypeMapper(t reflect.Type) TypeMapper {
	all := make(map[string]MappingInfo, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		if columnName != "" {
			all[columnName] = MappingInfo{
				FieldName:  field.Name,
				ColumnName: columnName,
				Type:       field.Type,
				Kind:       field.Type.Kind(),
			}
		}
	}
	return TypeMapper{
		Type:    t,
		Mapping: all,
	}
}

type MappingInfo struct {
	FieldName  string
	ColumnName string
	Type       reflect.Type
	Kind       reflect.Kind
}
