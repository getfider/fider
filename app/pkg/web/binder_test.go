package web_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

var binder = web.NewDefaultBinder()

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
	ctx := newPostContext(params, `{ "name": "Jon Snow" }`, "application/json")
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

	params := make(web.StringMap, 0)
	body := `{ "name": " Jon Snow ", "email": " JON.SNOW@got.com ", "color": "ff00ad" }`
	ctx := newPostContext(params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).To(BeNil())
	Expect(u.Name).To(Equal("Jon Snow"))
	Expect(u.Email).To(Equal("jon.snow@got.com"))
	Expect(u.Color).To(Equal("FF00AD"))
	Expect(u.Other).To(Equal(""))
}

func TestDefaultBinder_POST_WithoutBody(t *testing.T) {
	RegisterTestingT(t)

	type action struct {
		Number int    `route:"number"`
		Slug   string `route:"slug"`
	}

	params := make(web.StringMap, 0)
	params["number"] = "2"
	params["slug"] = "jon-snow"
	ctx := newPostContext(params, "", web.UTF8JSONContentType)
	a := new(action)
	err := binder.Bind(a, ctx)
	Expect(err).To(BeNil())
	Expect(a.Number).To(Equal(2))
	Expect(a.Slug).To(Equal("jon-snow"))
}

func TestDefaultBinder_POST_NonJSON(t *testing.T) {
	RegisterTestingT(t)

	type user struct {
		Name string `json:"name"`
	}

	params := make(web.StringMap, 0)
	ctx := newPostContext(params, `name=JonSnow`, web.UTF8HTMLContentType)

	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).To(Equal(web.ErrContentTypeNotAllowed))
}
