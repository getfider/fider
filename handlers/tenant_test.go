package handlers_test

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTenantHandler_404(t *testing.T) {
	RegisterTestingT(t)

	status, _ := makeRequest("GET", "/api/tenants/unknown")

	Expect(status).To(Equal(404))
}

func TestTenantHandler_200(t *testing.T) {
	RegisterTestingT(t)

	status, query := makeRequest("GET", "/api/tenants/trishop")

	Expect(query.Int("id")).To(Equal(2))
	Expect(query.String("name")).To(Equal("The Triathlon Shop"))
	Expect(query.String("domain")).To(Equal("trishop"))
	Expect(status).To(Equal(200))
}
