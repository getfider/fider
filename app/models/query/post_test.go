package query_test

import (
	"testing"

	"github.com/getfider/fider/app/models/query"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSearchPosts_SetLimitFromString(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"abc", ""},
		{"0", ""},
		{"-5", ""},
		{"1", "1"},
		{"30", "30"},
		{"1000", "1000"},
		{"1001", "1000"},
		{"2000000000", "1000"},
		{"99999999999999999999", ""},
		{"all", "1000"},
	}

	for _, testCase := range testCases {
		q := &query.SearchPosts{}
		q.SetLimitFromString(testCase.input)
		Expect(q.Limit).Equals(testCase.expected)
	}
}
