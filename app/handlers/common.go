package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

//Status returns some useful information
func Status(settings *models.SystemSettings) web.HandlerFunc {
	return func(c web.Context) error {
		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)

		return c.Ok(web.Map{
			"build":       settings.BuildTime,
			"version":     settings.Version,
			"env":         settings.Environment,
			"compiler":    settings.Compiler,
			"now":         time.Now().Format("2006.01.02.150405"),
			"goroutines":  runtime.NumGoroutine(),
			"workerQueue": c.Engine().Worker().Length(),
			"heapInMB":    memStats.HeapAlloc / 1048576,
			"stackInMB":   memStats.StackInuse / 1048576,
		})
	}
}

//Page returns a page without properties
func Page() web.HandlerFunc {
	return func(c web.Context) error {
		return c.Page(web.Map{})
	}
}

func validateKey(kind models.EmailVerificationKind, c web.Context) (*models.EmailVerification, error) {
	key := c.QueryParam("k")

	result, err := c.Services().Tenants.FindVerificationByKey(kind, key)
	if err != nil {
		if err == app.ErrNotFound {
			return nil, c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
		}
		return nil, c.Failure(err)
	}

	//If key has been used, just go back to home page
	if result.VerifiedOn != nil {
		return nil, c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
	}

	//If key expired, go back to home page
	if time.Now().After(result.ExpiresOn) {
		err = c.Services().Tenants.SetKeyAsVerified(key)
		if err != nil {
			return nil, c.Failure(err)
		}
		return nil, c.Redirect(http.StatusTemporaryRedirect, c.BaseURL())
	}

	return result, nil
}
