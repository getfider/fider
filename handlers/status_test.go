package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StatusHandler_ShouldReturnStatusBasedOnContext(t *testing.T) {
	status, query := makeRequest("GET", "/api/status")

	build, _ := query.String("build")
	isHealthy, _ := query.Bool("healthy", "database")
	assert.Equal(t, "today", build)
	assert.Equal(t, false, isHealthy)
	assert.Equal(t, 200, status)
}
