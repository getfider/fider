package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestCreateNewIdea_InvalidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	services.Ideas.Add("My great idea", "With a great description", 1)

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
		result := action.Validate(services)
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
		result := action.Validate(services)
		ExpectSuccess(result)
	}
}
