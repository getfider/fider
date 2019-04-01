package apiv1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app/handlers/apiv1"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestCreateTagHandler_ValidRequests(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		return app.ErrNotFound
	})

	var addNewTag *cmd.AddNewTag
	bus.AddHandler(func(ctx context.Context, c *cmd.AddNewTag) error {
		addNewTag = c
		return nil
	})

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		ExecutePost(
			apiv1.CreateEditTag(),
			`{ "name": "Feature Request", "color": "00FF00", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusOK)

	Expect(addNewTag.Name).Equals("Feature Request")
	Expect(addNewTag.Color).Equals("00FF00")
	Expect(addNewTag.IsPublic).IsTrue()
}

func TestCreateTagHandler_InvalidRequests(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		return app.ErrNotFound
	})

	var testCases = []string{
		`{ }`,
		`{ "name": "" }`,
		`{ "name": "Bug" }`,
		`{ "name": "Bug", "color": "ABC" }`,
		`{ "name": "Bug", "color": "00000X" }`,
		`{ "name": "123456789012345678901234567890A", "color": "000000" }`,
	}

	for _, testCase := range testCases {
		server, _ := mock.NewServer()
		status, _ := server.
			AsUser(mock.JonSnow).
			ExecutePostAsJSON(apiv1.CreateEditTag(), testCase)

		Expect(status).Equals(http.StatusBadRequest)
	}
}

func TestCreateTagHandler_AlreadyInUse(t *testing.T) {
	RegisterT(t)

	tag := &models.Tag{Name: "Bug", Slug: "bug"}
	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		q.Result = tag
		return nil
	})

	server, _ := mock.NewServer()

	status, _ := server.
		AsUser(mock.JonSnow).
		ExecutePostAsJSON(
			apiv1.CreateEditTag(),
			`{ "name": "Bug", "color": "0000FF", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusBadRequest)
}

func TestCreateTagHandler_Collaborator(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.AryaStark).
		ExecutePost(
			apiv1.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusForbidden)
}

func TestEditInvalidTagHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		return app.ErrNotFound
	})

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			apiv1.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusNotFound)
}

func TestEditExistingTagHandler(t *testing.T) {
	RegisterT(t)

	tag := &models.Tag{ID: 5, Name: "Bug", Slug: "bug", Color: "0000FF", IsPublic: true}
	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		q.Result = tag
		return nil
	})

	var updateTag *cmd.UpdateTag
	bus.AddHandler(func(ctx context.Context, c *cmd.UpdateTag) error {
		updateTag = c
		return nil
	})

	server, _ := mock.NewServer()

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			apiv1.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusOK)
	Expect(updateTag.TagID).Equals(5)
	Expect(updateTag.Name).Equals("Feature Request")
	Expect(updateTag.Color).Equals("000000")
	Expect(updateTag.IsPublic).IsTrue()
}

func TestDeleteInvalidTagHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		return app.ErrNotFound
	})

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		Execute(apiv1.DeleteTag())

	Expect(status).Equals(http.StatusNotFound)
}

func TestDeleteExistingTagHandler(t *testing.T) {
	RegisterT(t)

	tag := &models.Tag{ID: 5, Name: "Bug", Slug: "bug", Color: "0000FF", IsPublic: true}
	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		q.Result = tag
		return nil
	})

	var deleteTag *cmd.DeleteTag
	bus.AddHandler(func(ctx context.Context, c *cmd.DeleteTag) error {
		deleteTag = c
		return nil
	})

	server, _ := mock.NewServer()

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		Execute(apiv1.DeleteTag())

	Expect(status).Equals(http.StatusOK)
	Expect(deleteTag.Tag).Equals(tag)
}

func TestDeleteExistingTagHandler_Collaborator(t *testing.T) {
	RegisterT(t)

	tag := &models.Tag{ID: 5, Name: "Bug", Slug: "bug", Color: "0000FF", IsPublic: true}
	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		q.Result = tag
		return nil
	})

	server, _ := mock.NewServer()

	status, _ := server.
		AsUser(mock.AryaStark).
		AddParam("slug", "bug").
		Execute(apiv1.DeleteTag())

	Expect(status).Equals(http.StatusForbidden)
}

// EXPERIMENTAL-BUS: re-enable when PostStorage is on SQLStore
// func TestAssignTagHandler_Success(t *testing.T) {
// 	RegisterT(t)

// 	server, services := mock.NewServer()
// 	services.SetCurrentTenant(mock.DemoTenant)
// 	services.SetCurrentUser(mock.JonSnow)
// 	tag, _ := services.Tags.Add("Bug", "0000FF", true)
// 	post, _ := services.Posts.Add("Post Title", "Post Description")

// 	status, _ := server.
// 		AsUser(mock.JonSnow).
// 		AddParam("slug", tag.Slug).
// 		AddParam("number", post.Number).
// 		Execute(apiv1.AssignTag())

// 	tags, err := services.Tags.GetAssigned(post)
// 	Expect(status).Equals(http.StatusOK)
// 	Expect(err).IsNil()
// 	Expect(tags[0]).Equals(tag)
// }

func TestAssignTagHandler_UnknownTag(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetTagBySlug) error {
		return app.ErrNotFound
	})

	server, services := mock.NewServer()
	services.Posts.SetCurrentTenant(mock.DemoTenant)
	services.Posts.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("Post Title", "Post Description")

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		AddParam("number", post.ID).
		Execute(apiv1.AssignTag())

	Expect(status).Equals(http.StatusNotFound)
}

// EXPERIMENTAL-BUS: re-enable when PostStorage is on SQLStore
// func TestAssignOrUnassignTagHandler_Unauthorized(t *testing.T) {
// 	RegisterT(t)

// 	var testCases = []web.HandlerFunc{
// 		apiv1.AssignTag(),
// 		apiv1.UnassignTag(),
// 	}

// 	for _, handler := range testCases {
// 		server, services := mock.NewServer()
// 		services.SetCurrentTenant(mock.DemoTenant)
// 		services.SetCurrentUser(mock.JonSnow)
// 		tag, _ := services.Tags.Add("Bug", "0000FF", true)
// 		post, _ := services.Posts.Add("Post Title", "Post Description")

// 		status, _ := server.
// 			AsUser(mock.AryaStark).
// 			AddParam("slug", tag.Slug).
// 			AddParam("number", post.Number).
// 			Execute(handler)

// 		Expect(status).Equals(http.StatusForbidden)
// 	}
// }

// EXPERIMENTAL-BUS: re-enable when PostStorage is on SQLStore
// func TestUnassignTagHandler_Success(t *testing.T) {
// 	RegisterT(t)

// 	server, services := mock.NewServer()
// 	services.SetCurrentTenant(mock.DemoTenant)
// 	services.SetCurrentUser(mock.JonSnow)
// 	tag, _ := services.Tags.Add("Bug", "0000FF", true)
// 	post, _ := services.Posts.Add("Post Title", "Post Description")
// 	services.Tags.AssignTag(tag, post)

// 	status, _ := server.
// 		AsUser(mock.JonSnow).
// 		AddParam("slug", tag.Slug).
// 		AddParam("number", post.Number).
// 		Execute(apiv1.UnassignTag())

// 	tags, err := services.Tags.GetAssigned(post)
// 	Expect(status).Equals(http.StatusOK)
// 	Expect(err).IsNil()
// 	Expect(tags).HasLen(0)
// }

func TestListTagsHandler(t *testing.T) {
	RegisterT(t)

	tag1 := &models.Tag{ID: 2, Name: "Bug", Slug: "bug", Color: "0000FF", IsPublic: true}
	tag2 := &models.Tag{ID: 5, Name: "Feature Request", Slug: "feature-request", Color: "0000FF", IsPublic: true}
	bus.AddHandler(func(ctx context.Context, q *query.GetAllTags) error {
		q.Result = []*models.Tag{tag1, tag2}
		return nil
	})

	server, _ := mock.NewServer()

	status, query := server.
		AsUser(mock.JonSnow).
		ExecuteAsJSON(apiv1.ListTags())

	Expect(status).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(2)
}
