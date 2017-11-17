package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/pkg/mock"
	. "github.com/onsi/gomega"
)

func TestCreateTagHandler_ValidRequests(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "00FF00", "isPublic": true }`,
		)

	Expect(status).To(Equal(http.StatusOK))

	tag, err := services.Tags.GetBySlug("feature-request")
	Expect(err).To(BeNil())
	Expect(tag.Name).To(Equal("Feature Request"))
	Expect(tag.Slug).To(Equal("feature-request"))
	Expect(tag.Color).To(Equal("00FF00"))
	Expect(tag.IsPublic).To(BeTrue())
}

func TestCreateTagHandler_InvalidRequests(t *testing.T) {
	RegisterTestingT(t)

	var testCases = []struct {
		input    string
		failures []string
	}{
		{`{ }`, []string{"failures.name", "failures.color"}},
		{`{ "name": "" }`, []string{"failures.name", "failures.color"}},
		{`{ "name": "Bug" }`, []string{"failures.color"}},
		{`{ "name": "Bug", "color": "ABC" }`, []string{"failures.color"}},
		{`{ "name": "Bug", "color": "00000X" }`, []string{"failures.color"}},
		{`{ "name": "123456789012345678901234567890A", "color": "000000" }`, []string{"failures.name"}},
	}

	for _, testCase := range testCases {
		server, _ := mock.NewServer()
		status, query := server.
			AsUser(mock.JonSnow).
			ExecutePostAsJSON(handlers.CreateEditTag(), testCase.input)

		Expect(status).To(Equal(http.StatusBadRequest))
		for _, failure := range testCase.failures {
			Expect(query.Contains(failure)).To(BeTrue())
		}
	}

}

func TestCreateTagHandler_AlreadyInUse(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, query := server.
		AsUser(mock.JonSnow).
		ExecutePostAsJSON(
			handlers.CreateEditTag(),
			`{ "name": "Bug", "color": "0000FF", "isPublic": true }`,
		)

	Expect(status).To(Equal(http.StatusBadRequest))
	Expect(query.Contains("failures.name")).To(BeTrue())
}

func TestCreateTagHandler_Collaborator(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.AryaStark).
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).To(Equal(http.StatusForbidden))
}

func TestEditInvalidTagHandler(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).To(Equal(http.StatusNotFound))
}

func TestEditExistingTagHandler(t *testing.T) {
	RegisterTestingT(t)

	server, services := mock.NewServer()
	services.Tags.Add("Bug", "0000FF", true)

	status, _ := server.
		AsUser(mock.JonSnow).
		AddParam("slug", "bug").
		ExecutePost(
			handlers.CreateEditTag(),
			`{ "name": "Feature Request", "color": "000000", "isPublic": true }`,
		)

	Expect(status).To(Equal(http.StatusOK))
	tag, err := services.Tags.GetBySlug("bug")
	Expect(tag).To(BeNil())
	Expect(err).To(Equal(app.ErrNotFound))

	tag, err = services.Tags.GetBySlug("feature-request")
	Expect(tag).ToNot(BeNil())
	Expect(err).To(BeNil())
}
