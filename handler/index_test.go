package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/model"
	"github.com/WeCanHearYou/wchy/router"
	"github.com/labstack/echo"
	. "github.com/onsi/gomega"
)

type mockIdeaService struct{}

func (svc mockIdeaService) GetAll(tenantID int) ([]*model.Idea, error) {
	return make([]*model.Idea, 0), nil
}

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	ctx := &context.WchyContext{
		Idea: &mockIdeaService{},
	}

	e := echo.New()
	e.Renderer = router.NewHTMLRenderer()

	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("Tenant", &model.Tenant{ID: 2, Name: "Any Tenant"})

	handler.Index(ctx)(c)

	Expect(rec.Code).To(Equal(200))
}
