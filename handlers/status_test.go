package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StatusHandler", func() {
	It("should return status based on context", func() {
		status, query := makeRequest("GET", "/api/status")

		Expect(query.String("build")).To(Equal("today"))
		Expect(query.Bool("healthy", "database")).To(Equal(false))
		Expect(status).To(Equal(200))
	})
})
