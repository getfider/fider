package app

// OneYearCache adds Cache-Control header for one year
func OneYearCache() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			c.Response().Header().Add("Cache-Control", "public, max-age=30672000")
			return next(c)
		}
	}
}
