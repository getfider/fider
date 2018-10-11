package web_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

var binder = web.NewDefaultBinder()

func TestDefaultBinder_FromParams(t *testing.T) {
	type updateUser struct {
		Number int    `route:"number"`
		Slug   string `route:"slug"`
		Name   string `json:"name"`
	}

	RegisterT(t)
	params := make(web.StringMap, 0)
	params["number"] = "2"
	params["slug"] = "jon-snow"

	for _, method := range []string{"POST", "PUT", "DELETE"} {
		ctx := newBodyContext(method, params, `{ "name": "Jon Snow" }`, "application/json")
		u := new(updateUser)
		err := binder.Bind(u, ctx)
		Expect(err).IsNil()
		Expect(u.Number).Equals(2)
		Expect(u.Slug).Equals("jon-snow")
		Expect(u.Name).Equals("Jon Snow")
	}
}

func TestDefaultBinder_TrimSpaces(t *testing.T) {
	RegisterT(t)

	type user struct {
		Name  string `json:"name"`
		Email string `json:"email" format:"lower"`
		Color string `json:"color" format:"upper"`
		Other string
	}

	params := make(web.StringMap, 0)
	body := `{ "name": " Jon Snow ", "email": " JON.SNOW@got.com ", "color": "ff00ad" }`
	ctx := newBodyContext("POST", params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(u.Name).Equals("Jon Snow")
	Expect(err).IsNil()
	Expect(u.Email).Equals("jon.snow@got.com")
	Expect(u.Color).Equals("FF00AD")
	Expect(u.Other).Equals("")
}

func TestDefaultBinder_Base64ToBytes(t *testing.T) {
	RegisterT(t)

	resume, _ := ioutil.ReadFile(env.Path("./README.md"))
	resumeBase64 := base64.StdEncoding.EncodeToString(resume)

	type user struct {
		Name   string `json:"name"`
		Resume []byte `json:"resume"`
	}

	params := make(web.StringMap, 0)
	body := fmt.Sprintf(`{ "name": "Jon Snow", "resume": "%s" }`, resumeBase64)
	ctx := newBodyContext("POST", params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.Name).Equals("Jon Snow")
	Expect(u.Resume).Equals(resume)
}

func TestDefaultBinder_NestedJson(t *testing.T) {
	RegisterT(t)

	resume, _ := ioutil.ReadFile(env.Path("./README.md"))
	resumeBase64 := base64.StdEncoding.EncodeToString(resume)

	type user struct {
		Name   string `json:"name"`
		Resume struct {
			FileName string `json:"fileName"`
			File     []byte `json:"file"`
		}
	}

	params := make(web.StringMap, 0)
	body := fmt.Sprintf(`{ "name": "Jon Snow", "resume": { "fileName": "README.md", "file": "%s" } }`, resumeBase64)
	ctx := newBodyContext("POST", params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.Name).Equals("Jon Snow")
	Expect(u.Resume.FileName).Equals("README.md")
	Expect(u.Resume.File).Equals(resume)
}

func TestDefaultBinder_Array_TrimSpaces(t *testing.T) {
	RegisterT(t)

	type user struct {
		Providers []string `json:"providers" format:"lower"`
	}

	params := make(web.StringMap, 0)
	body := `{ "providers": [ " Google", " FACEBOOK ", "   MicroSoft    " ] }`
	ctx := newBodyContext("POST", params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.Providers).Equals([]string{
		"google",
		"facebook",
		"microsoft",
	})
}

func TestDefaultBinder_DELETE(t *testing.T) {
	RegisterT(t)

	type user struct {
		Name  string `json:"name"`
		Other string
	}

	params := make(web.StringMap, 0)
	body := `{ "name": " Jon Snow " }`
	ctx := newBodyContext("DELETE", params, body, "application/json")
	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.Name).Equals("Jon Snow")
	Expect(u.Other).Equals("")
}

func TestDefaultBinder_POST_WithoutBody(t *testing.T) {
	RegisterT(t)

	type action struct {
		Number int    `route:"number"`
		Slug   string `route:"slug"`
	}

	params := make(web.StringMap, 0)
	params["number"] = "2"
	params["slug"] = "jon-snow"
	ctx := newBodyContext("POST", params, "", web.UTF8JSONContentType)
	a := new(action)
	err := binder.Bind(a, ctx)
	Expect(err).IsNil()
	Expect(a.Number).Equals(2)
	Expect(a.Slug).Equals("jon-snow")
}

func TestDefaultBinder_POST_NonJSON(t *testing.T) {
	RegisterT(t)

	type user struct {
		Name string `json:"name"`
	}

	params := make(web.StringMap, 0)
	ctx := newBodyContext("POST", params, `name=JonSnow`, web.UTF8HTMLContentType)

	u := new(user)
	err := binder.Bind(u, ctx)
	Expect(err).Equals(web.ErrContentTypeNotAllowed)
}

type Size int
type Size32 int32
type Size64 int64

const (
	Small   Size   = 1
	Large   Size   = 2
	Small32 Size32 = 3
	Large32 Size32 = 4
	Small64 Size64 = 5
	Large64 Size64 = 6
)

type pizza struct {
	TheSize   Size   `route:"size"`
	TheSize32 Size32 `route:"size32"`
	TheSize64 Size64 `route:"size64"`
}

func (p *Size) UnmarshalText(s []byte) error {
	if string(s) == "large" {
		*p = Large
	} else {
		*p = Small
	}
	return nil
}

func TestDefaultBinder_POST_CustomType(t *testing.T) {
	RegisterT(t)

	params := make(web.StringMap, 0)
	params["size"] = "1"
	params["size32"] = "4"
	params["size64"] = "5"
	ctx := newBodyContext("POST", params, `{ }`, "application/json")

	u := new(pizza)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.TheSize).Equals(Small)
	Expect(u.TheSize32).Equals(Large32)
	Expect(u.TheSize64).Equals(Small64)
}

func TestDefaultBinder_POST_CustomType_Unmarshal(t *testing.T) {
	RegisterT(t)

	type pizza struct {
		TheSize Size `route:"size"`
	}

	params := make(web.StringMap, 0)
	params["size"] = "large"
	ctx := newBodyContext("POST", params, `{ }`, "application/json")

	u := new(pizza)
	err := binder.Bind(u, ctx)
	Expect(err).IsNil()
	Expect(u.TheSize).Equals(Large)
}
