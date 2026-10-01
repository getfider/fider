package query_test

import (
	"testing"

	"github.com/getfider/fider/app/models/query"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSearchPosts_SetLimitFromString(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		input     string
		unlimited bool
		expected  string
	}{
		{"", false, ""},
		{"abc", false, ""},
		{"0", false, ""},
		{"-5", false, ""},
		{"1", false, "1"},
		{"30", false, "30"},
		{"1000", false, "1000"},
		{"1001", false, "1000"},
		{"2000000000", false, "1000"},
		{"99999999999999999999", false, ""},
		{"all", false, "1000"},
		{"", true, ""},
		{"abc", true, ""},
		{"-5", true, ""},
		{"30", true, "30"},
		{"5000", true, "5000"},
		{"all", true, "all"},
	}

	for _, testCase := range testCases {
		q := &query.SearchPosts{}
		q.SetLimitFromString(testCase.input, testCase.unlimited)
		Expect(q.Limit).Equals(testCase.expected)
	}
}
