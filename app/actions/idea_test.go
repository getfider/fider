package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestCreateNewIdea_InvalidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	services.SetCurrentUser(&models.User{ID: 1})
	services.Ideas.Add("My great idea", "With a great description")

	for _, title := range []string{
		"me",
		"",
		"  ",
		"signup",
		"My great great great great great great great great great great great great great great great great great idea.",
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
		"my GREAT idea",
	} {
		action := &actions.CreateNewIdea{Model: &models.NewIdea{Title: title}}
		result := action.Validate(nil, services)
		ExpectFailed(result, "title")
	}
}

func TestCreateNewIdea_ValidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	for _, title := range []string{
		"this is my new idea",
		"this idea is very descriptive",
	} {
		action := &actions.CreateNewIdea{Model: &models.NewIdea{Title: title}}
		result := action.Validate(nil, services)
		ExpectSuccess(result)
	}
}

func TestSetResponse_InvalidStatus(t *testing.T) {
	RegisterTestingT(t)

	action := &actions.SetResponse{Model: &models.SetResponse{
		Status: models.IdeaDeleted,
		Text:   "Spam!",
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "status")
}

func TestDeleteIdea_WhenIsBeingReferenced(t *testing.T) {
	RegisterTestingT(t)

	services.SetCurrentUser(&models.User{ID: 1})
	idea1, _ := services.Ideas.Add("Idea #1", "")
	idea2, _ := services.Ideas.Add("Idea #2", "")
	services.Ideas.MarkAsDuplicate(idea2, idea1)

	model := &models.DeleteIdea{
		Number: idea2.Number,
		Text:   "Spam!",
	}
	action := &actions.DeleteIdea{Model: model}
	ExpectSuccess(action.Validate(nil, services))

	model.Number = idea1.Number
	ExpectFailed(action.Validate(nil, services))
}
