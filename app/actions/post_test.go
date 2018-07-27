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
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
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
