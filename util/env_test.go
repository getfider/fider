package util_test

import (
	"testing"

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
