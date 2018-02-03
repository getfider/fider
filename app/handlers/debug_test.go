package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestRuntimeStatsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	status, query := server.ExecuteAsJSON(handlers.RuntimeStats())

	Expect(query.Contains("goroutines")).To(BeTrue())
	Expect(query.Contains("heapInMB")).To(BeTrue())
	Expect(query.Contains("stackInMB")).To(BeTrue())
	Expect(status).To(Equal(http.StatusOK))
}
