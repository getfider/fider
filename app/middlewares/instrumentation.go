package middlewares

import (
	"fmt"
	"strconv"
	"time"

	"github.com/getfider/fider/app/metrics"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/julienschmidt/httprouter"
)

// Instrumentation adds Prometheus HTTP Middlewares
func Instrumentation() web.MiddlewareFunc {
	if !env.Config.Metrics.Enabled {
		return nil
	}

	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			begin := time.Now()

			err := next(c)

			routeName := c.Param(httprouter.MatchedRoutePathParam)
			operation := fmt.Sprintf("%s %s", c.Request.Method, routeName)
			if routeName == "" && c.Request.URL.Path != "/" {
				operation = "No Route Matched"
			}

			code := strconv.Itoa(c.Response.StatusCode)
			metrics.HttpRequests.WithLabelValues(code, operation).Inc()
			metrics.HttpDuration.WithLabelValues(operation).Observe(time.Since(begin).Seconds())

			return err
		}
	}
}
