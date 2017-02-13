package handler_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/env"
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

	path := "views/*.html"
	if env.IsTest() {
		path = os.Getenv("GOPATH") + "/src/github.com/WeCanHearYou/wchy/" + path
	}

	e.Renderer = &router.HTMLRenderer{
		Templates: template.Must(template.ParseGlob(path)),
	}

	req, _ := http.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("Tenant", &model.Tenant{ID: 2, Name: "Any Tenant"})

	handler.Index(ctx)(c)

	Expect(rec.Code).To(Equal(200))
}
