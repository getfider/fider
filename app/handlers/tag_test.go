package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

func TestCreateTagHandler_ValidRequests(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "00FF00", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusOK)

	tag, err := services.Tags.GetBySlug("feature-request")
	Expect(err).IsNil()
	Expect(tag.Name).Equals("Feature Request")
	Expect(tag.Slug).Equals("feature-request")
	Expect(tag.Color).Equals("00FF00")
	Expect(tag.IsPublic).IsTrue()
}

func TestCreateTagHandler_InvalidRequests(t *testing.T) {
	RegisterT(t)

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
			ExecutePostAsJSON(handlers.CreateEditTag(), testCase)

		Expect(status).Equals(http.StatusBadRequest)
	}
}

func TestCreateTagHandler_AlreadyInUse(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, _ := server.
		AsUser(mock.JonSnow).
		ExecutePostAsJSON(
			handlers.CreateEditTag(),
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
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusForbidden)
}

func TestEditInvalidTagHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusNotFound)
}

func TestEditExistingTagHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).Equals(http.StatusOK)
	tag, err := services.Tags.GetBySlug("bug")
	Expect(tag).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)

	tag, err = services.Tags.GetBySlug("feature-request")
	Expect(tag).IsNotNil()
	Expect(err).IsNil()
}

func TestDeleteInvalidTagHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		Execute(handlers.DeleteTag())

	Expect(status).Equals(http.StatusNotFound)
}

func TestDeleteExistingTagHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		Execute(handlers.DeleteTag())

	tag, err := services.Tags.GetBySlug("bug")
	Expect(status).Equals(http.StatusOK)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(tag).IsNil()
}

func TestDeleteExistingTagHandler_Collaborator(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, _ := server.
		AsUser(mock.AryaStark).
		AddParam("slug", "bug").
		Execute(handlers.DeleteTag())

	tag, err := services.Tags.GetBySlug("bug")
	Expect(status).Equals(http.StatusForbidden)
	Expect(tag).IsNotNil()
	Expect(err).IsNil()

}

func TestAssignTagHandler_Success(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	tag, _ := services.Tags.Add("Bug", "0000FF", true)
	post, _ := services.Posts.Add("Post Title", "Post Description")

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", tag.Slug).
		AddParam("number", post.Number).
		Execute(handlers.AssignTag())

	tags, err := services.Tags.GetAssigned(post)
	Expect(status).Equals(http.StatusOK)
	Expect(err).IsNil()
	Expect(tags[0]).Equals(tag)
}

func TestAssignTagHandler_UnknownTag(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		AddParam("number", 1).
		Execute(handlers.AssignTag())

	Expect(status).Equals(http.StatusNotFound)
}

func TestAssignOrUnassignTagHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	var testCases = []web.HandlerFunc{
		handlers.AssignTag(),
		handlers.UnassignTag(),
	}

	for _, handler := range testCases {
		server, services := mock.NewServer()
		services.SetCurrentTenant(mock.DemoTenant)
		services.SetCurrentUser(mock.JonSnow)
		tag, _ := services.Tags.Add("Bug", "0000FF", true)
		post, _ := services.Posts.Add("Post Title", "Post Description")

		status, _ := server.
			AsUser(mock.AryaStark).
			AddParam("slug", tag.Slug).
			AddParam("number", post.Number).
			Execute(handler)

		Expect(status).Equals(http.StatusForbidden)
	}
}

func TestUnassignTagHandler_Success(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	tag, _ := services.Tags.Add("Bug", "0000FF", true)
	post, _ := services.Posts.Add("Post Title", "Post Description")
	services.Tags.AssignTag(tag, post)

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", tag.Slug).
		AddParam("number", post.Number).
		Execute(handlers.UnassignTag())

	tags, err := services.Tags.GetAssigned(post)
	Expect(status).Equals(http.StatusOK)
	Expect(err).IsNil()
	Expect(tags).HasLen(0)
}
