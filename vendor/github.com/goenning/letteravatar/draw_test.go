package letteravatar_test

import (
	"testing"

	"github.com/goenning/letteravatar"
)

func TestExtract(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Jon Snow",
			expected: "JS",
		},
		{
			input:    "Jon von Snow",
			expected: "JS",
		},
		{
			input:    "jon.snow",
			expected: "JS",
		},
		{
			input:    "Jon Snow Start",
			expected: "JS",
		},
		{
			input:    "Arya stark",
			expected: "AS",
		},
		{
			input:    " Arya stark ",
			expected: "AS",
		},
		{
			input:    "Arya",
			expected: "A",
		},
		{
			input:    "Arya99",
			expected: "A9",
		},
		{
			input:    "J",
			expected: "J",
		},
		{
			input:    "",
			expected: "?",
		},
		{
			input:    "AryaStark",
			expected: "AS",
		},
		{
			input:    "Jon (Snow)",
			expected: "J",
		},
		{
			input:    "Jon Wallet (JS)",
			expected: "JW",
		},
	}

	for _, testCase := range testCases {
		output := letteravatar.Extract(testCase.input)
		if output != testCase.expected {
			t.Errorf("[%s]: %s != %s ", testCase.input, output, testCase.expected)
		}
	}
}
