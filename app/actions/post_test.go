package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestCreateNewPost_InvalidPostTitles(t *testing.T) {
	RegisterT(t)

	services.SetCurrentUser(&models.User{ID: 1})
	services.Posts.Add("My great post", "With a great description")

	for _, title := range []string{
		"me",
		"",
		"  ",
		"signup",
		"My great great great great great great great great great great great great great great great great great post.",
		"my GREAT post",
	} {
		action := &actions.CreateNewPost{Model: &models.NewPost{Title: title}}
		result := action.Validate(nil, services)
		ExpectFailed(result, "title")
	}
}

func TestCreateNewPost_ValidPostTitles(t *testing.T) {
	RegisterT(t)

	for _, title := range []string{
		"this is my new post",
		"this post is very descriptive",
	} {
		action := &actions.CreateNewPost{Model: &models.NewPost{Title: title}}
		result := action.Validate(nil, services)
		ExpectSuccess(result)
	}
}

func TestSetResponse_InvalidStatus(t *testing.T) {
	RegisterT(t)

	action := &actions.SetResponse{Model: &models.SetResponse{
		Status: models.PostDeleted,
		Text:   "Spam!",
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "status")
}

func TestDeletePost_WhenIsBeingReferenced(t *testing.T) {
	RegisterT(t)

	services.SetCurrentUser(&models.User{ID: 1})
	post1, _ := services.Posts.Add("Post #1", "")
	post2, _ := services.Posts.Add("Post #2", "")
	services.Posts.MarkAsDuplicate(post2, post1)

	model := &models.DeletePost{
		Number: post2.Number,
		Text:   "Spam!",
	}
	action := &actions.DeletePost{Model: model}
	ExpectSuccess(action.Validate(nil, services))

	model.Number = post1.Number
	ExpectFailed(action.Validate(nil, services))
}

func TestDeleteComment(t *testing.T) {
	RegisterT(t)

	author := &models.User{ID: 1, Role: models.RoleVisitor}
	notAuthor := &models.User{ID: 2, Role: models.RoleVisitor}
	administrator := &models.User{ID: 3, Role: models.RoleAdministrator}

	services.SetCurrentUser(author)
	post1, _ := services.Posts.Add("Post #1", "")
	commentID, _ := services.Posts.AddComment(post1, "Comment #1")

	action := &actions.DeleteComment{
		Model: &models.DeleteComment{
			CommentID: commentID,
		},
	}

	authorized := action.IsAuthorized(notAuthor, services)
	Expect(authorized).IsFalse()

	authorized = action.IsAuthorized(administrator, services)
	Expect(authorized).IsTrue()
}
