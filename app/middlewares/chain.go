package middlewares

import (
	"github.com/getfider/fider/app/pkg/web"
)

// Chain combines multiple middlewares into one
func Chain(mws ...web.MiddlewareFunc) web.MiddlewareFunc {
	return func(handler web.HandlerFunc) web.HandlerFunc {
		next := handler
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			if mw != nil {
				next = mw(next)
			}
		}
		return next
	}
}
