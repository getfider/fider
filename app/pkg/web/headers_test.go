package web_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestIsJSONContentType(t *testing.T) {
	testCases := []struct {
		contentType string
		expected    bool
	}{
		{"application/json", true},
		{"application/json; charset=utf-8", true},
		{"application/json;charset=UTF-8", true},
		{"Application/JSON", true},
		{"", false},
		{"text/plain", false},
		{"text/plain; x=application/json", false},
		{"application/x-www-form-urlencoded", false},
		{"multipart/form-data; boundary=x", false},
		{"application/jsonx", false},
		{"application/json;", true},
		{"application/json; charset", true},
		{"application/json; charset=utf-8; charset=utf-8", true},
		{"application/json;;;=", true},
		{"text/plain; charset", false},
		{"text/plain; charset=utf-8; charset=utf-8", false},
		{"text/plain;;;=", false},
		{"application/json/x", false},
		{"/json", false},
	}

	for _, tc := range testCases {
		t.Run(tc.contentType, func(t *testing.T) {
			RegisterT(t)
			Expect(web.IsJSONContentType(tc.contentType)).Equals(tc.expected)
		})
	}
}

func TestParseBearerToken(t *testing.T) {
	testCases := []struct {
		header        string
		expectedToken string
		expectedOk    bool
	}{
		{"Bearer MY-KEY", "MY-KEY", true},
		{"Bearer   MY-KEY  ", "MY-KEY", true},
		{"Bearer\tMY-KEY", "MY-KEY", true},
		{"BearerMY-KEY", "MY-KEY", true},
		{"bearer MY-KEY", "MY-KEY", true},
		{"BEARER MY-KEY", "MY-KEY", true},
		{"Bearer MY-KEY\r\n", "MY-KEY", true},
		{"Bearer MY-KEY \t ", "MY-KEY", true},
		{"", "", false},
		{"Bear", "", false},
		{"Bearer", "", false},
		{"Bearer ", "", false},
		{"Bearer \t\r\n", "", false},
		{"Basic dXNlcjpwYXNz", "", false},
		{"Basic Bearer MY-KEY", "", false},
		{"Token Bearer", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.header, func(t *testing.T) {
			RegisterT(t)
			token, ok := web.ParseBearerToken(tc.header)
			Expect(ok).Equals(tc.expectedOk)
			Expect(token).Equals(tc.expectedToken)
		})
	}
}
