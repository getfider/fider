package handlers

import (
	"github.com/getfider/fider/app/pkg/web"
)

// ManageTags is the home page for managing tags
func ManageTags() web.HandlerFunc {
	return func(c web.Context) error {
		tags, err := c.Services().Tags.GetAll()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title: "Manage Tags · Site Settings",
			Data: web.Map{
				"tags": tags,
			},
		})
	}
}
