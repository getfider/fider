package handlers

import "github.com/getfider/fider/app/pkg/web"

func GetMentions() web.HandlerFunc {
	return func(c web.Context) error {
		query := c.QueryParam("query")
		users, err := c.Services().Users.FindLike(query)

		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(users)
	}

}
