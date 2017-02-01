package env_test

import (
	"testing"

	"os"

	"github.com/WeCanHearYou/wchy/env"
	"github.com/stretchr/testify/assert"
)

var envs = []struct {
	go_env string
	domain string
	env    string
}{
	{"test", "test.canhearyou.com", "test"},
	{"staging", "staging.canhearyou.com", "staging"},
	{"development", "dev.canhearyou.com", "development"},
	{"production", "canhearyou.com", "production"},
	{"anything", "canhearyou.com", "development"},
}

func TestGetEnvOrDefault_ShouldReturnNullWhenKeyIsUnknown(t *testing.T) {
	value := env.GetEnvOrDefault("UNKNOWN_KEY", "some value")
	assert.Equal(t, "some value", value)
}

func TestGetEnvOrDefault_ShouldReturnEnvironmentValue(t *testing.T) {
	value := env.GetEnvOrDefault("PATH", "default path")
	assert.NotEqual(t, "default path", value)
}

func TestGetCurrentDomain_ShouldReturnCorrectDomain(t *testing.T) {
	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := env.GetCurrentDomain()
		assert.Equal(t, testCase.domain, actual)
	}
}

func TestEnvironment_ShouldReturnCorrectEnvironment(t *testing.T) {
	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := env.Current()
		assert.Equal(t, testCase.env, actual)
	}
}
