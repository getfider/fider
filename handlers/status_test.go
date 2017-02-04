package handlers_test

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)

	status, query := makeRequest("GET", "/api/status")

	Expect(query.String("build")).To(Equal("today"))
	Expect(query.Bool("healthy", "database")).To(Equal(false))
	Expect(status).To(Equal(200))

}
