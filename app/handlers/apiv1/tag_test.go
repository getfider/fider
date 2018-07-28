package apiv1_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers/apiv1"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
)

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
		Execute(apiv1.AssignTag())

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
		Execute(apiv1.AssignTag())

	Expect(status).Equals(http.StatusNotFound)
}

func TestAssignOrUnassignTagHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	var testCases = []web.HandlerFunc{
		apiv1.AssignTag(),
		apiv1.UnassignTag(),
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
		Execute(apiv1.UnassignTag())

	tags, err := services.Tags.GetAssigned(post)
	Expect(status).Equals(http.StatusOK)
	Expect(err).IsNil()
	Expect(tags).HasLen(0)
}
