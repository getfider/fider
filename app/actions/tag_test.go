package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/rand"
)

func TestCreateEditTag_InvalidName(t *testing.T) {
	RegisterT(t)

	services.Tags.Add("Feature Request", "000000", true)

	for _, name := range []string{
		"",
		"Feature Request",
		rand.String(31),
	} {
		action := &actions.CreateEditTag{Model: &models.CreateEditTag{Name: name, Color: "FFFFFF"}}
		result := action.Validate(nil, services)
		ExpectFailed(result, "name")
	}
}

func TestCreateEditTag_InvalidColor(t *testing.T) {
	RegisterT(t)

	for _, color := range []string{
		"",
		"ABC",
		"PPPOOO",
		"FFF",
		"000000X",
	} {
		action := &actions.CreateEditTag{Model: &models.CreateEditTag{Name: "Bug", Color: color}}
		result := action.Validate(nil, services)
		ExpectFailed(result, "color")
	}
}

func TestCreateEditTag_ValidInput(t *testing.T) {
	RegisterT(t)

	tag, _ := services.Tags.Add("To Discuss", "000000", true)

	action := &actions.CreateEditTag{Model: &models.CreateEditTag{Name: "Bug", Color: "FF0000"}}
	result := action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(action.Tag).IsNil()

	action = &actions.CreateEditTag{Model: &models.CreateEditTag{Name: "New Name", Slug: "to-discuss", Color: "FF0000"}}
	result = action.Validate(nil, services)
	ExpectSuccess(result)
	Expect(action.Tag).Equals(tag)
}
