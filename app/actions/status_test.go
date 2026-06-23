package actions_test

import (
	"context"
	"strings"
	"testing"

	"github.com/getfider/fider/app/actions"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestCreateStatus_Validate_Success(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateStatus{
		Slug:       "Triage",
		Label:      " Triage ",
		Kind:       "open",
		Color:      "blue",
		Icon:       "lightbulb",
		ShowOnHome: true,
		Filterable: true,
		SortOrder:  15,
	}

	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsTrue()
	Expect(action.Slug).Equals("triage")
	Expect(action.Label).Equals("Triage")
}

func TestCreateStatus_Validate_RejectsBadSlug(t *testing.T) {
	RegisterT(t)

	cases := []struct {
		name string
		slug string
	}{
		{"empty", ""},
		{"uppercase", "Triage"}, // normalized first, so this becomes "triage" — actually OK
		{"underscores", "in_review"},
		{"spaces", "in review"},
		{"leading-hyphen", "-triage"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			RegisterT(t)
			action := &actions.CreateStatus{Slug: tc.slug, Label: "Triage", Kind: "open", Color: "blue"}
			result := action.Validate(context.Background(), nil)
			if tc.name == "uppercase" {
				// Uppercase gets normalized in Validate, so it actually succeeds.
				Expect(result.Ok).IsTrue()
				return
			}
			Expect(result.Ok).IsFalse()
		})
	}
}

func TestCreateStatus_Validate_RejectsLongSlug(t *testing.T) {
	RegisterT(t)
	action := &actions.CreateStatus{Slug: strings.Repeat("a", 51), Label: "Triage", Kind: "open", Color: "blue"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsFalse()
}

func TestCreateStatus_Validate_RejectsUnknownKind(t *testing.T) {
	RegisterT(t)
	action := &actions.CreateStatus{Slug: "triage", Label: "Triage", Kind: "pending", Color: "blue"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsFalse()
}

func TestCreateStatus_Validate_RejectsUnknownColor(t *testing.T) {
	RegisterT(t)
	action := &actions.CreateStatus{Slug: "triage", Label: "Triage", Kind: "open", Color: "fuchsia"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsFalse()
}

func TestCreateStatus_Validate_RejectsLongLabel(t *testing.T) {
	RegisterT(t)
	action := &actions.CreateStatus{Slug: "triage", Label: strings.Repeat("L", 51), Kind: "open", Color: "blue"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsFalse()
}

func TestUpdateStatus_Validate_Success(t *testing.T) {
	RegisterT(t)
	action := &actions.UpdateStatus{Label: "Renamed", Color: "green", Icon: "check"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsTrue()
}

func TestUpdateStatus_Validate_DefaultsAppliedForEmptyColorIcon(t *testing.T) {
	RegisterT(t)
	action := &actions.UpdateStatus{Label: "Renamed"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsTrue()
	Expect(action.Color).Equals("blue")
	Expect(action.Icon).Equals("lightbulb")
}

func TestUpdateStatus_Validate_RejectsBadColor(t *testing.T) {
	RegisterT(t)
	action := &actions.UpdateStatus{Label: "Renamed", Color: "fuchsia"}
	result := action.Validate(context.Background(), nil)
	Expect(result.Ok).IsFalse()
}
