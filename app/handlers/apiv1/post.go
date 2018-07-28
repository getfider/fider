package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// SearchPosts return existing posts based on search criteria
func SearchPosts() web.HandlerFunc {
	return func(c web.Context) error {
		posts, err := c.Services().Posts.Search(
			c.QueryParam("q"),
			c.QueryParam("f"),
			c.QueryParam("l"),
			c.QueryParamAsArray("t"),
		)
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(posts)
	}
}

// CreatePost creates a new post on current tenant
func CreatePost() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateNewPost)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		posts := c.Services().Posts
		post, err := posts.Add(input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		if err := posts.AddSupporter(post, c.User()); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewPost(post))

		return c.Ok(post)
	}
}
