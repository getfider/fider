package web_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

var binder = web.NewDefaultBinder()

func newGetContext(params web.StringMap) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func newPostContext(params web.StringMap, body string) *web.Context {
	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	ctx := e.NewContext(res, req, params)
	return &ctx
}

func TestDefaultBinder_FromParams(t *testing.T) {
	type updateUser struct {
		Number int    `route:"number"`
		Slug   string `route:"slug"`
		Name   string `json:"name"`
	}

	RegisterTestingT(t)
	params := make(web.StringMap, 0)
	params["number"] = "2"
	params["slug"] = "jon-snow"
	ctx := newPostContext(params, `{ "name": "Jon Snow" }`)
	u := new(updateUser)
	err := binder.Bind(u, ctx)
	Expect(err).To(BeNil())
	Expect(u.Number).To(Equal(2))
	Expect(u.Slug).To(Equal("jon-snow"))
	Expect(u.Name).To(Equal("Jon Snow"))
}

func TestDefaultBinder_TrimSpaces(t *testing.T) {
	RegisterTestingT(t)

	type user struct {
		Name  string `json:"name"`
		Email string `json:"email" format:"lower"`
		Color string `json:"color" format:"upper"`
		Other string
	}

	ctx := newPostContext(make(web.StringMap, 0), `{ "name": " Jon Snow ", "email": " JON.SNOW@got.com ", "color": "ff00ad" }`)
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).To(BeNil())
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.Email).To(Equal("jon.snow@got.com"))
	Expect(u.Color).To(Equal("FF00AD"))
	Expect(u.Other).To(Equal(""))
}
