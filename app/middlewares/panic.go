package middlewares

import (
	"io"

	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

func CatchPanic() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err := c.Failure(errors.Panicked(r))
					log.Error(c, err)
					c.Rollback()

					if f, ok := c.Response.(io.Closer); ok {
						f.Close()
					}
				}
			}()

			return next(c)
		}
	}
}
