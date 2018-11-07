package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

//Health always returns OK
func Health() web.HandlerFunc {
	return func(c web.Context) error {
		return c.Ok(web.Map{})
	}
}

//LegalPage returns a legal page with content from a file
func LegalPage(title, file string) web.HandlerFunc {
	return func(c web.Context) error {
		bytes, err := ioutil.ReadFile(env.Etc(file))
		if err != nil {
			return c.NotFound()
		}

		return c.Render(http.StatusOK, "legal.html", web.Props{
			Title: title,
			Data: web.Map{
				"Content": string(bytes),
			},
		})
	}
}

//Sitemap returns the sitemap.xml of current site
func Sitemap() web.HandlerFunc {
	return func(c web.Context) error {
		if c.Tenant().IsPrivate {
			return c.NotFound()
		}

		posts, err := c.Services().Posts.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		baseURL := c.BaseURL()
		text := strings.Builder{}
		text.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
		text.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
		text.WriteString(fmt.Sprintf("<url> <loc>%s</loc> </url>", baseURL))
		for _, post := range posts {
			text.WriteString(fmt.Sprintf("<url> <loc>%s/posts/%d/%s</loc> </url>", baseURL, post.Number, post.Slug))
		}
		text.WriteString(`</urlset>`)

		c.Response.Header().Del("Content-Security-Policy")
		return c.XML(http.StatusOK, text.String())
	}
}

//RobotsTXT return content of robots.txt file
func RobotsTXT() web.HandlerFunc {
	return func(c web.Context) error {
		bytes, err := ioutil.ReadFile(env.Path("./robots.txt"))
		if err != nil {
			return c.NotFound()
		}
		sitemapURL := c.BaseURL() + "/sitemap.xml"
		content := fmt.Sprintf("%s\nSitemap: %s", bytes, sitemapURL)
		return c.String(http.StatusOK, content)
	}
}

//Page returns a page without properties
func Page(title, description string) web.HandlerFunc {
	return func(c web.Context) error {
		return c.Page(web.Props{
			Title:       title,
			Description: description,
		})
	}
}

//BrowserNotSupported returns an error page for browser that Fider dosn't support
func BrowserNotSupported() web.HandlerFunc {
	return func(c web.Context) error {
		return c.Render(http.StatusOK, "browser-not-supported.html", web.Props{
			Title:       "Browser not supported",
			Description: "We don't support this version of your browser",
		})
	}
}

//NewLogError is the input model for UI errors
type NewLogError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//LogError logs an error coming from the UI
func LogError() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(NewLogError)
		err := c.Bind(input)
		if err != nil {
			return c.Failure(err)
		}
		c.Logger().Errorf(input.Message, log.Props{
			"Data": input.Data,
		})
		return c.Ok(web.Map{})
	}
}

func validateKey(kind models.EmailVerificationKind, c web.Context) (*models.EmailVerification, error) {
	key := c.QueryParam("k")

	//If key has been used, return NotFound
	result, err := c.Services().Tenants.FindVerificationByKey(kind, key)
	if err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			return nil, c.NotFound()
		}
		return nil, c.Failure(err)
	}

	//If key has been used, return Gone
	if result.VerifiedAt != nil {
		return nil, c.Gone()
	}

	//If key expired, return Gone
	if time.Now().After(result.ExpiresAt) {
		err = c.Services().Tenants.SetKeyAsVerified(key)
		if err != nil {
			return nil, c.Failure(err)
		}
		return nil, c.Gone()
	}

	return result, nil
}
