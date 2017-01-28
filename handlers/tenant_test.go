package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GivenUnknownDomain_TenantByDomain_ShouldReturn404(t *testing.T) {
	status, _ := makeRequest("GET", "/tenants/unknown")
	assert.Equal(t, 404, status)
}

func Test_GivenKnownDomain_TenantByDomain_ShouldReturnTenantData(t *testing.T) {
	status, query := makeRequest("GET", "/tenants/trishop")

	id, _ := query.Int("id")
	name, _ := query.String("name")
	domain, _ := query.String("domain")
	assert.Equal(t, 2, id)
	assert.Equal(t, "The Triathlon Shop", name)
	assert.Equal(t, "trishop", domain)
	assert.Equal(t, 200, status)
}
