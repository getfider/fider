package postgres_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/storage/postgres"
)

func TestToTSQuery(t *testing.T) {
	RegisterT(t)

	var testcases = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"123 hello", "123|hello"},
		{" hello  ", "hello"},
		{" hello$ world$ ", "hello|world"},
		{" yes, please ", "yes|please"},
		{" yes / please ", "yes|please"},
		{" hello 'world' ", "hello|world"},
		{"hello|world", "hello|world"},
		{"hello | world", "hello|world"},
		{"hello & world", "hello|world"},
	}

	for _, testcase := range testcases {
		output := postgres.ToTSQuery(testcase.input)
		Expect(output).Equals(testcase.expected)
	}
}
