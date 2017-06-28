package middlewares

import (
	"fmt"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage/postgres"
	"github.com/labstack/gommon/log"
)

//Setup current context with some services
func Setup(db *dbx.Database) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			logger := c.Logger().(*log.Logger)
			logger.Debugf("HTTP Request %s", logger.Color().Bold(logger.Color().RedBg(c.Request().Method+" "+c.Request().URL.String())))

			trx, err := db.Begin()
			if err != nil {
				return err
			}

			c.SetActiveTransaction(trx)
			c.SetServices(&app.Services{
				Tenants: postgres.NewTenantStorage(trx),
			})

			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					c.Error(err)
				}
			}()

			err = next(c)

			if err != nil {
				trx.Rollback()
				return err
			}
			return trx.Commit()
		}
	}
}

//AddServices adds services to current context
func AddServices() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			trx := c.ActiveTransaction()
			tenant := c.Tenant()
			services := c.Services()
			services.OAuth = &oauth.HTTPService{}
			services.Ideas = postgres.NewIdeaStorage(tenant, trx)
			services.Users = postgres.NewUserStorage(trx)
			return next(c)
		}
	}
}
