package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TenantHandler", func() {
	It("should return 404 when domain is unknown", func() {
		status, _ := makeRequest("GET", "/api/tenants/unknown")

		Expect(status).To(Equal(404))
	})
	It("should return tenant data when domain is known", func() {
		status, query := makeRequest("GET", "/api/tenants/trishop")

		Expect(query.Int("id")).To(Equal(2))
		Expect(query.String("name")).To(Equal("The Triathlon Shop"))
		Expect(query.String("domain")).To(Equal("trishop"))
		Expect(status).To(Equal(200))
	})
})
