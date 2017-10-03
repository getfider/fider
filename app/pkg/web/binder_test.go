package web_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

type user struct {
	Name  string `json:"name"`
	Email string `json:"email" format:"lower"`
	Other string
}

func TestDefaultBinder_TrimSpaces(t *testing.T) {
	RegisterTestingT(t)

	e := web.New(nil)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{ "name": " Jon Snow ", "email": " JON.SNOW@got.com " }`))
	u := new(user)
	ctx := e.NewContext(res, req, nil)
	binder := web.NewDefaultBinder()
	err := binder.Bind(u, &ctx)
	Expect(err).To(BeNil())
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.Email).To(Equal("jon.snow@got.com"))
	Expect(u.Other).To(Equal(""))
}
