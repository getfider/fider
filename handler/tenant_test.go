package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/service"
	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

type mockTenantService struct{}

func (svc mockTenantService) GetByDomain(domain string) (*model.Tenant, error) {
	if domain == "trishop" {
		return &model.Tenant{ID: 2, Name: "The Triathlon Shop", Domain: "trishop"}, nil
	}
	return nil, service.ErrNotFound
}

var ctx *context.WchyContext = &context.WchyContext{
	Tenant: &mockTenantService{},
}

func TestTenantHandler_404(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("domain")
	c.SetParamValues("unknown")
	handler.TenantByDomain(ctx)(c)

	Expect(rec.Code).To(Equal(404))
}

func TestTenantHandler_200(t *testing.T) {
	RegisterTestingT(t)

	e := echo.New()
	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("domain")
	c.SetParamValues("trishop")
	handler.TenantByDomain(ctx)(c)

	var data interface{}
	decoder := json.NewDecoder(rec.Body)
	decoder.Decode(&data)
	query := jsonq.NewQuery(data)

	Expect(query.Int("id")).To(Equal(2))
	Expect(query.String("name")).To(Equal("The Triathlon Shop"))
	Expect(query.String("domain")).To(Equal("trishop"))
	Expect(rec.Code).To(Equal(200))
}
