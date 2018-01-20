package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestInvalidUserNames(t *testing.T) {
	RegisterTestingT(t)

	for _, name := range []string{
		"",
		"123456789012345678901234567890123456789012345678901", // 51 chars
	} {
		action := &actions.UpdateUserSettings{Model: &models.UpdateUserSettings{Name: name}}
		result := action.Validate(nil, services)
		ExpectFailed(result, "name")
	}
}

func TestValidUserNames(t *testing.T) {
	RegisterTestingT(t)

	for _, name := range []string{
		"Jon Snow",
		"Arya",
	} {
		action := &actions.UpdateUserSettings{Model: &models.UpdateUserSettings{Name: name}}
		result := action.Validate(nil, services)
		ExpectSuccess(result)
	}
}
