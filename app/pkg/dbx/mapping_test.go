package dbx_test

import (
	"reflect"
	"testing"

	"github.com/getfider/fider/app/pkg/dbx"
	. "github.com/onsi/gomega"
)

func createScanner(values ...interface{}) func(dest ...interface{}) error {
	return func(dest ...interface{}) error {
		for i := 0; i < len(values); i++ {
			reflect.Indirect(reflect.ValueOf(dest[i])).Set(reflect.ValueOf(values[i]))
		}
		return nil
	}
}

func TestMapping_SimpleStruct(t *testing.T) {
	RegisterTestingT(t)
	mapper := dbx.NewRowMapper()
	u := user{}
	mapper.Map(&u, []string{"id", "name"}, createScanner(1, "Hello World"))
	Expect(u.ID).To(Equal(1))
	Expect(u.Name).To(Equal("Hello World"))
}
