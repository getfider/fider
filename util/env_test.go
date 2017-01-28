package util_test

import (
	"testing"

	"os"

	"github.com/WeCanHearYou/wchy-api/util"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault_ShouldReturnNullWhenKeyIsUnknown(t *testing.T) {
	value := util.GetEnvOrDefault("UNKNOWN_KEY", "some value")
	assert.Equal(t, "some value", value)
}

func TestGetEnvOrDefault_ShouldReturnEnvironmentValue(t *testing.T) {
	value := util.GetEnvOrDefault("PATH", "default path")
	assert.NotEqual(t, "default path", value)
}

func TestGetCurrentDomain_ShouldReturnCorrectDomain(t *testing.T) {
	var tests = []struct {
		env      string
		expected string
	}{
		{"test", "test.canhearyou.com"},
		{"staging", "staging.canhearyou.com"},
		{"development", "dev.canhearyou.com"},
		{"production", "canhearyou.com"},
		{"anything", "canhearyou.com"},
	}

	for _, testCase := range tests {
		os.Setenv("GO_ENV", testCase.env)
		actual := util.GetCurrentDomain()
		assert.Equal(t, testCase.expected, actual)
	}
}
