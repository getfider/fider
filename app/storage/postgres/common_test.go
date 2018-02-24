package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/storage/postgres"
	. "github.com/onsi/gomega"
)

func TestToTSQuery(t *testing.T) {
	RegisterTestingT(t)

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
		Expect(output).To(Equal(testcase.expected))
	}
}
