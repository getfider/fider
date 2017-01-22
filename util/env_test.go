package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault_ShouldReturnNullWhenKeyIsUnknown(t *testing.T) {
	value := GetEnvOrDefault("UNKNOWN_KEY", "some value")
	assert.Equal(t, "some value", value)
}

func TestGetEnvOrDefault_ShouldReturnEnvironmentValue(t *testing.T) {
	value := GetEnvOrDefault("PATH", "default path")
	assert.NotEqual(t, "default path", value)
}
