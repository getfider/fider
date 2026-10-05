package apiv1_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/middlewares"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestCreatePostHandler(t *testing.T) {
	RegisterT(t)

	var newPost *cmd.AddNewPost
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewPost) error {
		newPost = c
		c.Result = &entity.Post{
			ID:          1,
			Title:       c.Title,
			Description: c.Description,
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newPost.Title).Equals("My newest post :)")
	Expect(newPost.Description).Equals("")
}

func TestCreatePostHandler_AppendsUnreferencedAttachments(t *testing.T) {
	RegisterT(t)

	var newPost *cmd.AddNewPost
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewPost) error {
		newPost = c
		c.Result = &entity.Post{ID: 1, Title: c.Title, Description: c.Description}
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{
			"title": "Post with attachments",
			"description": "Already referenced: ![](fider-image:attachments/referenced.png)",
			"attachments": [
				{ "bkey": "attachments/referenced.png" },
				{ "bkey": "attachments/standalone.png" }
			]
		}`)

	Expect(code).Equals(http.StatusOK)
	// The already-referenced attachment is left as-is (not duplicated), and the
	// unreferenced one is appended at the end as a fider-image markdown reference.
	Expect(newPost.Description).Equals("Already referenced: ![](fider-image:attachments/referenced.png)\n\n![](fider-image:attachments/standalone.png)")
}

func executeSearchPosts(limit string, configure func(server *mock.Server) *mock.Server) (int, *query.SearchPosts, bool) {
	var searchPosts *query.SearchPosts
	bus.AddHandler(func(ctx context.Context, q *query.SearchPosts) error {
		searchPosts = q
		return nil
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetUserByAPIKey) error {
		if q.APIKey == "1234567890" {
			q.Result = mock.JonSnow
			return nil
		}
		return app.ErrNotFound
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		q.Result = mock.JonSnow
		return nil
	})

	isAuthenticated := false
	server := mock.NewServer().
		OnTenant(mock.DemoTenant).
		WithURL("http://demo.test.fider.io/api/v1/posts?limit=" + limit).
		Use(middlewares.User()).
		Use(func(next web.HandlerFunc) web.HandlerFunc {
			return func(c *web.Context) error {
				isAuthenticated = c.IsAuthenticated()
				return next(c)
			}
		})
	status, _ := configure(server).Execute(apiv1.SearchPosts())
	return status, searchPosts, isAuthenticated
}

func TestSearchPostsHandler_ClampsLimit(t *testing.T) {
	RegisterT(t)

	sessionToken, _ := jwt.Encode(jwt.FiderClaims{UserID: mock.JonSnow.ID, UserName: mock.JonSnow.Name})

	anonymous := func(server *mock.Server) *mock.Server { return server }
	withSession := func(server *mock.Server) *mock.Server {
		return server.AddCookie(web.CookieAuthName, sessionToken)
	}

	for _, limit := range []string{"all", "2000000000"} {
		status, searchPosts, isAuthenticated := executeSearchPosts(limit, anonymous)
		Expect(status).Equals(http.StatusOK)
		Expect(isAuthenticated).IsFalse()
		Expect(searchPosts.Limit).Equals("1000")

		status, searchPosts, isAuthenticated = executeSearchPosts(limit, withSession)
		Expect(status).Equals(http.StatusOK)
		Expect(isAuthenticated).IsTrue()
		Expect(searchPosts.Limit).Equals("1000")
	}
}

func TestSearchPostsHandler_APIKeyIsNotClamped(t *testing.T) {
	RegisterT(t)

	withAPIKey := func(server *mock.Server) *mock.Server {
		return server.AddHeader("Authorization", "Bearer 1234567890")
	}

	status, searchPosts, isAuthenticated := executeSearchPosts("all", withAPIKey)
	Expect(status).Equals(http.StatusOK)
	Expect(isAuthenticated).IsTrue()
	Expect(searchPosts.Limit).Equals("all")

	status, searchPosts, _ = executeSearchPosts("5000", withAPIKey)
	Expect(status).Equals(http.StatusOK)
	Expect(searchPosts.Limit).Equals("5000")

	status, searchPosts, _ = executeSearchPosts("-5", withAPIKey)
	Expect(status).Equals(http.StatusOK)
	Expect(searchPosts.Limit).Equals("")
}

func TestCreatePostHandler_WithoutTitle(t *testing.T) {
	RegisterT(t)

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		ExecutePost(apiv1.CreatePost(), `{ "title": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCreatePostHandler_WithNonExistentTag(t *testing.T) {
	if env.Config.PostCreationWithTagsEnabled {
		RegisterT(t)

		bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
			return app.ErrNotFound
		})
		bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
			return app.ErrNotFound
		})

		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)", "tags": ["inexistent_tag"]}`)

		Expect(code).Equals(http.StatusBadRequest)
	}
}

func TestCreatePostHandler_WithPrivateTagAsVisitor(t *testing.T) {
	if env.Config.PostCreationWithTagsEnabled {
		RegisterT(t)

		privateTag := &entity.Tag{
			ID:       1,
			Name:     "private_tag",
			Slug:     "private_tag",
			Color:    "blue",
			IsPublic: false,
		}
		bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
			if q.Slug == "private_tag" {
				q.Result = privateTag
				return nil
			}
			return app.ErrNotFound
		})

		bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
			return app.ErrNotFound
		})

		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.AryaStark).
			ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)", "tags": ["private_tag"]}`)

		Expect(code).Equals(http.StatusForbidden)
	}
}

func TestCreatePostHandler_WithPublicTagAsVisitor(t *testing.T) {
	if env.Config.PostCreationWithTagsEnabled {
		RegisterT(t)

		var newPost *cmd.AddNewPost
		bus.AddHandler(func(ctx context.Context, c *cmd.AddNewPost) error {
			newPost = c
			c.Result = &entity.Post{
				ID:          1,
				Title:       c.Title,
				Description: c.Description,
			}
			return nil
		})

		publicTag := &entity.Tag{
			ID:       1,
			Name:     "public_tag",
			Slug:     "public_tag",
			Color:    "red",
			IsPublic: true,
		}
		bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
			if q.Slug == "public_tag" {
				q.Result = publicTag
				return nil
			}
			return app.ErrNotFound
		})

		var tagAssignment *cmd.AssignTag
		bus.AddHandler(func(ctx context.Context, c *cmd.AssignTag) error {
			tagAssignment = c
			return nil
		})

		bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
			return app.ErrNotFound
		})

		bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
		bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error { return nil })
		bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.AryaStark).
			ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)", "tags": ["public_tag"]}`)

		Expect(code).Equals(http.StatusOK)
		Expect(tagAssignment.Tag).Equals(publicTag)
		Expect(tagAssignment.Post).Equals(newPost.Result)
	}
}

func TestCreatePostHandler_WithPublicTagAndPrivateTagAsCollaborator(t *testing.T) {
	if env.Config.PostCreationWithTagsEnabled {
		RegisterT(t)

		var newPost *cmd.AddNewPost
		bus.AddHandler(func(ctx context.Context, c *cmd.AddNewPost) error {
			newPost = c
			c.Result = &entity.Post{
				ID:          1,
				Title:       c.Title,
				Description: c.Description,
			}
			return nil
		})

		publicTag := &entity.Tag{
			ID:       1,
			Name:     "public_tag",
			Slug:     "public_tag",
			Color:    "red",
			IsPublic: true,
		}
		privateTag := &entity.Tag{
			ID:       1,
			Name:     "private_tag",
			Slug:     "private_tag",
			Color:    "blue",
			IsPublic: false,
		}
		bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
			if q.Slug == "public_tag" {
				q.Result = publicTag
				return nil
			}
			if q.Slug == "private_tag" {
				q.Result = privateTag
				return nil
			}
			return app.ErrNotFound
		})

		tagAssignments := make([]*cmd.AssignTag, 2)
		bus.AddHandler(func(ctx context.Context, c *cmd.AssignTag) error {
			switch c.Tag.Slug {
			case "public_tag":
				tagAssignments[0] = c
			case "private_tag":
				tagAssignments[1] = c
			}
			return nil
		})

		bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
			return app.ErrNotFound
		})

		bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
		bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error { return nil })
		bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			ExecutePost(apiv1.CreatePost(), `{ "title": "My newest post :)", "tags": ["public_tag", "private_tag"]}`)

		Expect(code).Equals(http.StatusOK)
		Expect(tagAssignments[0].Tag).Equals(publicTag)
		Expect(tagAssignments[1].Tag).Equals(privateTag)
		Expect(tagAssignments[0].Post).Equals(newPost.Result)
		Expect(tagAssignments[1].Post).Equals(newPost.Result)
	}
}

func TestGetPostHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecuteAsJSON(apiv1.GetPost())

	Expect(code).Equals(http.StatusOK)
	Expect(query.String("title")).Equals(post.Title)
	Expect(query.String("description")).Equals(post.Description)
}

func TestUpdatePostHandler_TenantStaff(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "With a description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	var updatePost *cmd.UpdatePost
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdatePost) error {
		updatePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(updatePost.Post).Equals(post)
	Expect(updatePost.Title).Equals("the new title")
	Expect(updatePost.Description).Equals("new description")
}

func TestUpdatePostHandler_AppendsUnreferencedAttachments(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "With a description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	var updatePost *cmd.UpdatePost
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdatePost) error {
		updatePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.UpdatePost(), `{
			"title": "the new title",
			"description": "Already referenced: ![](fider-image:attachments/referenced.png)",
			"attachments": [
				{ "bkey": "attachments/referenced.png" },
				{ "bkey": "attachments/standalone.png" }
			]
		}`)

	Expect(code).Equals(http.StatusOK)
	// The already-referenced attachment is left as-is (not duplicated), and the
	// unreferenced one is appended at the end as a fider-image markdown reference.
	Expect(updatePost.Description).Equals("Already referenced: ![](fider-image:attachments/referenced.png)\n\n![](fider-image:attachments/standalone.png)")
}

func TestUpdatePostHandler_NonAuthorized(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:          5,
		Number:      5,
		Title:       "My First Post",
		Description: "Such an amazing description",
		User:        mock.JonSnow,
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestUpdatePostHandler_IsOwner_AfterGracePeriod(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:          5,
		Number:      5,
		Title:       "My First Post",
		Description: "Such an amazing description",
		User:        mock.AryaStark,
		CreatedAt:   time.Now().UTC().Add(-2 * time.Hour),
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestUpdatePostHandler_IsOwner_WithinGracePeriod(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{
		ID:          5,
		Number:      5,
		Title:       "My First Post",
		Description: "Such an amazing description",
		User:        mock.AryaStark,
		CreatedAt:   time.Now().UTC(),
	}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})
	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.UploadImages) error { return nil })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, cmd *cmd.UpdatePost) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", "5").
		ExecutePost(apiv1.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestUpdatePostHandler_InvalidTitle(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "", "description": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdatePostHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 5, Number: 5, Title: "My First Post", Description: "Such an amazing description"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error { return app.ErrNotFound })
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error { return app.ErrNotFound })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 999).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "This is a good title!", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestUpdatePostHandler_DuplicateTitle(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{ID: 1, Number: 1, Title: "My First Post", Slug: "my-first-post"}
	post2 := &entity.Post{ID: 2, Number: 2, Title: "My Second Post", Slug: "my-second-post"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post1.Number {
			q.Result = post1
			return nil
		}
		if q.Number == post2.Number {
			q.Result = post2
			return nil
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostBySlug) error {
		if q.Slug == post1.Slug {
			q.Result = post1
			return nil
		}
		if q.Slug == post2.Slug {
			q.Result = post2
			return nil
		}
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.Number).
		ExecutePost(apiv1.UpdatePost(), `{ "title": "My Second Post", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "My First Post", Slug: "my-first-post"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	var setResponse *cmd.SetPostResponse
	bus.AddHandler(func(ctx context.Context, c *cmd.SetPostResponse) error {
		setResponse = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.SetResponse(), fmt.Sprintf(`{ "status": "%s", "text": "Done!" }`, enum.PostCompleted.Name()))

	Expect(code).Equals(http.StatusOK)
	Expect(setResponse.Post).Equals(post)
	Expect(setResponse.Status).Equals(enum.PostCompleted)
	Expect(setResponse.Text).Equals("Done!")
}

func TestSetResponseHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 5).
		ExecutePost(apiv1.SetResponse(), fmt.Sprintf(`{ "status": "%s", "text": "Done!" }`, enum.PostCompleted.Name()))

	Expect(code).Equals(http.StatusForbidden)
}

func TestSetResponseHandler_Duplicate(t *testing.T) {
	RegisterT(t)

	var markAsDuplicate *cmd.MarkPostAsDuplicate
	bus.AddHandler(func(ctx context.Context, c *cmd.MarkPostAsDuplicate) error {
		markAsDuplicate = c
		return nil
	})

	post1 := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	post2 := &entity.Post{ID: 2, Number: 2, Title: "The Post #2", Description: "The Description #2"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post1.Number {
			q.Result = post1
			return nil
		}
		if q.Number == post2.Number {
			q.Result = post2
			return nil
		}
		return app.ErrNotFound
	})

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": %d }`, enum.PostDuplicate.Name(), post2.Number)
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusOK)
	Expect(markAsDuplicate.Post).Equals(post1)
	Expect(markAsDuplicate.Original).Equals(post2)
}

func TestSetResponseHandler_Duplicate_NotFound(t *testing.T) {
	RegisterT(t)

	post1 := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post1.Number {
			q.Result = post1
			return nil
		}
		return app.ErrNotFound
	})

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": 9999 }`, enum.PostDuplicate.Name())
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler_Duplicate_Itself(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		if q.Number == post.Number {
			q.Result = post
			return nil
		}
		return app.ErrNotFound
	})

	body := fmt.Sprintf(`{ "status": "%s", "originalNumber": %d }`, enum.PostDuplicate.Name(), post.Number)
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.ID).
		ExecutePost(apiv1.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddVoteHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var addVote *cmd.AddVote
	bus.AddHandler(func(ctx context.Context, c *cmd.AddVote) error {
		addVote = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		Execute(apiv1.AddVote())

	Expect(code).Equals(http.StatusOK)
	Expect(addVote.Post).Equals(post)
	Expect(addVote.User).Equals(mock.AryaStark)
}

func TestAddVoteHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		return app.ErrNotFound
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 999).
		Execute(apiv1.AddVote())

	Expect(code).Equals(http.StatusNotFound)
}

func TestRemoveVoteHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var removeVote *cmd.RemoveVote
	bus.AddHandler(func(ctx context.Context, c *cmd.RemoveVote) error {
		removeVote = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.ID).
		Execute(apiv1.RemoveVote())

	Expect(code).Equals(http.StatusOK)
	Expect(removeVote.Post).Equals(post)
	Expect(removeVote.User).Equals(mock.AryaStark)
}

func TestDeletePostHandler_Authorized(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.PostIsReferenced) error {
		q.Result = false
		return nil
	})

	var deletePost *cmd.SetPostResponse
	bus.AddHandler(func(ctx context.Context, c *cmd.SetPostResponse) error {
		deletePost = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.DeletePost(), `{ }`)

	Expect(code).Equals(http.StatusOK)
	Expect(deletePost.Post).Equals(post)
	Expect(deletePost.Status).Equals(enum.PostDeleted)
	Expect(deletePost.Text).Equals("")
}

func TestPostCommentHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var newComment *cmd.AddNewComment
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewComment) error {
		newComment = c
		c.Result = &entity.Comment{ID: 1, Content: c.Content}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newComment.Post).Equals(post)
	Expect(newComment.Content).Equals("This is a comment!")
}

func TestPostCommentHandlerMentions(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	var newComment *cmd.AddNewComment
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewComment) error {
		newComment = c
		c.Result = &entity.Comment{ID: 1, Content: c.Content}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.PostComment(), `{ "content": "Hello @[Jon Snow]!" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(newComment.Post).Equals(post)
	Expect(newComment.Content).Equals("Hello @[Jon Snow]!")
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(apiv1.PostComment(), `{ "content": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.AryaStark}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByID) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.SetAttachments) error { return nil })
	bus.AddHandler(func(ctx context.Context, c *cmd.UploadImages) error { return nil })

	var updateComment *cmd.UpdateComment
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateComment) error {
		updateComment = c
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		ExecutePost(apiv1.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusOK)
	Expect(updateComment.Content).Equals("My first comment has been edited")
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.JonSnow}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		ExecutePost(apiv1.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestListCommentHandler(t *testing.T) {
	RegisterT(t)

	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1", Description: "The Description #1"}
	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = post
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentsByPost) error {
		q.Result = []*entity.Comment{
			{ID: 1, Content: "First Comment"},
			{ID: 2, Content: "First Comment"},
		}
		return nil
	})

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecuteAsJSON(apiv1.ListComments())

	Expect(code).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(2)
}

func TestCommentReactionToggleHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = &entity.Post{ID: 1, Number: 1}
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.AryaStark}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	testCases := []struct {
		name     string
		user     *entity.User
		reaction string
	}{
		{"JonSnow reacts with like", mock.JonSnow, "👍"},
		{"AryaStark reacts with smile", mock.AryaStark, "👍"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var toggleReaction *cmd.ToggleCommentReaction
			bus.AddHandler(func(ctx context.Context, c *cmd.ToggleCommentReaction) error {
				toggleReaction = c
				return nil
			})

			code, _ := mock.NewServer().
				OnTenant(mock.DemoTenant).
				AsUser(tc.user).
				AddParam("number", 1).
				AddParam("id", comment.ID).
				AddParam("reaction", tc.reaction).
				ExecutePost(apiv1.ToggleReaction(), ``)

			Expect(code).Equals(http.StatusOK)
			Expect(toggleReaction.Emoji).Equals(tc.reaction)
			Expect(toggleReaction.Comment).Equals(comment)
			Expect(toggleReaction.User).Equals(tc.user)
		})
	}
}

func TestCommentReactionToggleHandler_InvalidEmoji(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = &entity.Post{ID: 1, Number: 1}
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.AryaStark}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ToggleCommentReaction) error {
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 1).
		AddParam("id", comment.ID).
		AddParam("reaction", "like").
		ExecutePost(apiv1.ToggleReaction(), ``)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestCommentReactionToggleHandler_UnAuthorised(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = &entity.Post{ID: 1, Number: 1}
		return nil
	})

	comment := &entity.Comment{ID: 5, Content: "Old comment text", User: mock.AryaStark}
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		q.Result = comment
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ToggleCommentReaction) error {
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AddParam("number", 1).
		AddParam("id", comment.ID).
		AddParam("reaction", "👍").
		ExecutePost(apiv1.ToggleReaction(), ``)

	Expect(code).Equals(http.StatusForbidden)
}

func TestCommentReactionToggleHandler_MismatchingTenantAndComment(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		q.Result = &entity.Post{ID: 1, Number: 1}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.ToggleCommentReaction) error {
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 1).
		AddParam("id", 1).
		AddParam("reaction", "👍").
		ExecutePost(apiv1.ToggleReaction(), ``)

	Expect(code).Equals(http.StatusNotFound)
}

// mockCommentOnPost registers a visible post #1 holding comment #5. Any other post number is
// treated as not visible to the caller, as GetPostByNumber does for deleted or unapproved posts.
func mockCommentOnPost() (*entity.Post, *entity.Comment) {
	post := &entity.Post{ID: 1, Number: 1, Title: "The Post #1"}
	other := &entity.Post{ID: 2, Number: 2, Title: "The Post #2"}
	comment := &entity.Comment{ID: 5, Content: "A comment", User: mock.AryaStark}

	bus.AddHandler(func(ctx context.Context, q *query.GetPostByNumber) error {
		switch q.Number {
		case post.Number:
			q.Result = post
		case other.Number:
			q.Result = other
		default:
			return app.ErrNotFound
		}
		return nil
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentByID) error {
		if q.CommentID != comment.ID || (q.PostID != 0 && q.PostID != post.ID) {
			return app.ErrNotFound
		}
		q.Result = comment
		return nil
	})

	return post, comment
}

func TestGetCommentHandler(t *testing.T) {
	RegisterT(t)
	post, comment := mockCommentOnPost()

	code, query := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		ExecuteAsJSON(apiv1.GetComment())

	Expect(code).Equals(http.StatusOK)
	Expect(query.Int32("id")).Equals(comment.ID)
	Expect(query.String("content")).Equals(comment.Content)
}

func TestGetCommentHandler_HiddenPost(t *testing.T) {
	RegisterT(t)
	_, comment := mockCommentOnPost()

	// Post #3 stands in for a deleted, declined or unapproved post
	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AddParam("number", 3).
		AddParam("id", comment.ID).
		ExecuteAsJSON(apiv1.GetComment())

	Expect(code).Equals(http.StatusNotFound)
}

func TestGetCommentHandler_WrongPost(t *testing.T) {
	RegisterT(t)
	_, comment := mockCommentOnPost()

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AddParam("number", 2).
		AddParam("id", comment.ID).
		ExecuteAsJSON(apiv1.GetComment())

	Expect(code).Equals(http.StatusNotFound)
}

func TestCommentReactionToggleHandler_HiddenOrWrongPost(t *testing.T) {
	RegisterT(t)
	_, comment := mockCommentOnPost()

	toggled := false
	bus.AddHandler(func(ctx context.Context, c *cmd.ToggleCommentReaction) error {
		toggled = true
		return nil
	})

	for _, number := range []int{2, 3} {
		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.JonSnow).
			AddParam("number", number).
			AddParam("id", comment.ID).
			AddParam("reaction", "👍").
			ExecutePost(apiv1.ToggleReaction(), ``)

		Expect(code).Equals(http.StatusNotFound)
	}
	Expect(toggled).IsFalse()
}

func TestUpdateCommentHandler_HiddenOrWrongPost(t *testing.T) {
	RegisterT(t)
	_, comment := mockCommentOnPost()

	updated := false
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateComment) error {
		updated = true
		return nil
	})

	for _, number := range []int{2, 3} {
		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.AryaStark).
			AddParam("number", number).
			AddParam("id", comment.ID).
			ExecutePost(apiv1.UpdateComment(), `{ "content": "Edited" }`)

		Expect(code).Equals(http.StatusForbidden)
	}
	Expect(updated).IsFalse()
}

func TestDeleteCommentHandler_HiddenOrWrongPost(t *testing.T) {
	RegisterT(t)
	_, comment := mockCommentOnPost()

	deleted := false
	bus.AddHandler(func(ctx context.Context, c *cmd.DeleteComment) error {
		deleted = true
		return nil
	})

	for _, number := range []int{2, 3} {
		code, _ := mock.NewServer().
			OnTenant(mock.DemoTenant).
			AsUser(mock.AryaStark).
			AddParam("number", number).
			AddParam("id", comment.ID).
			Execute(apiv1.DeleteComment())

		Expect(code).Equals(http.StatusForbidden)
	}
	Expect(deleted).IsFalse()
}

func TestDeleteCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)
	post, comment := mockCommentOnPost()

	var deleteComment *cmd.DeleteComment
	bus.AddHandler(func(ctx context.Context, c *cmd.DeleteComment) error {
		deleteComment = c
		return nil
	})

	code, _ := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", comment.ID).
		Execute(apiv1.DeleteComment())

	Expect(code).Equals(http.StatusOK)
	Expect(deleteComment.CommentID).Equals(comment.ID)
}
