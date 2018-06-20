package log_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/log"
)

func TestParseText(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		format   string
		props    log.Props
		expected string
	}{
		{"Hello World", nil, "Hello World"},
		{"Hello World", log.Props{}, "Hello World"},
		{"Hello @{Name}", log.Props{
			"Name": "John",
		}, "Hello John"},
		{"My name is @{Name} and I'm @{Age} years old", log.Props{
			"Name": "John",
			"Age":  55,
		}, "My name is John and I'm 55 years old"},
		{"Hello @{Name}", nil, "Hello @{Name}"},
		{"Hello @{Name}", log.Props{"Age": 55}, "Hello <nil>"},
		{"Hello @{Name:blue}", log.Props{
			"Name": "John",
		}, "Hello \033[34mJohn\033[0m"},
		{"Hello @{Name:undefined}", log.Props{
			"Name": "John",
		}, "Hello John"},
	}

	for _, testCase := range testCases {
		text := log.Parse(testCase.format, testCase.props)
		Expect(text).Equals(testCase.expected)
	}
}
