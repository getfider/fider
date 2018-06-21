package log_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/log"
)

func TestLevel_Parse(t *testing.T) {
	RegisterT(t)

	Expect(log.ParseLevel("DEBUG")).Equals(log.DEBUG)
	Expect(log.ParseLevel("INFO")).Equals(log.INFO)
	Expect(log.ParseLevel("WARN")).Equals(log.WARN)
	Expect(log.ParseLevel("ERROR")).Equals(log.ERROR)
	Expect(log.ParseLevel("NONE")).Equals(log.INFO)
	Expect(log.ParseLevel("")).Equals(log.INFO)
}

func TestLevel_ToString(t *testing.T) {
	RegisterT(t)

	Expect(log.DEBUG.String()).Equals("DEBUG")
	Expect(log.INFO.String()).Equals("INFO")
	Expect(log.WARN.String()).Equals("WARN")
	Expect(log.ERROR.String()).Equals("ERROR")
	Expect(log.DEBUG.String()).Equals("DEBUG")
	Expect(log.Level(88).String()).Equals("???")
}

func TestPropsMerge(t *testing.T) {
	RegisterT(t)

	Expect(log.Props{
		"Name": "John",
	}.Merge(nil)).Equals(log.Props{"Name": "John"})

	Expect(log.Props{
		"Name": "John",
	}.Merge(log.Props{
		"Name": "Maria",
	})).Equals(log.Props{"Name": "Maria"})

	Expect(log.Props(nil).Merge(log.Props{
		"Name": "Maria",
	})).Equals(log.Props{"Name": "Maria"})

	Expect(log.Props{
		"Name": "John",
		"DOB":  "05/02/1998",
	}.Merge(log.Props{
		"Name": "Maria",
		"Age":  66,
	})).Equals(log.Props{
		"Age":  66,
		"DOB":  "05/02/1998",
		"Name": "Maria",
	})
}

func TestParseText(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		format   string
		props    log.Props
		colorize bool
		expected string
	}{
		{"Hello World", nil, true, "Hello World"},
		{"Hello World", log.Props{}, true, "Hello World"},
		{"Hello @{Name}", log.Props{
			"Name": "John",
		}, true, "Hello John"},
		{"My name is @{Name} and I'm @{Age} years old", log.Props{
			"Name": "John",
			"Age":  55,
		}, true, "My name is John and I'm 55 years old"},
		{"Hello @{Name}", nil, true, "Hello @{Name}"},
		{"Hello @{Name}", log.Props{"Age": 55}, true, "Hello <nil>"},
		{"Hello @{Name:blue}", log.Props{
			"Name": "John",
		}, true, "Hello \033[34mJohn\033[0m"},
		{"Hello @{Name:undefined}", log.Props{
			"Name": "John",
		}, true, "Hello John"},
		{"Hello @{Name:blue}", log.Props{
			"Name": "John",
		}, false, "Hello John"},
	}

	for _, testCase := range testCases {
		text := log.Parse(testCase.format, testCase.props, testCase.colorize)
		Expect(text).Equals(testCase.expected)
	}
}
