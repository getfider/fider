package handlers

import (
	"encoding/json"
	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

// ExportTagsJSON returns a JSON file with all tags
func ExportTagsJSON() web.HandlerFunc {
	return func(c *web.Context) error {
		getAllTags := &query.GetAllTags{}
		if err := bus.Dispatch(c, getAllTags); err != nil {
			return c.Failure(err)
		}

		data, err := json.MarshalIndent(getAllTags.Result, "", "  ")
		if err != nil {
			return c.Failure(errors.Wrap(err, "failed to marshal tags"))
		}

		return c.Attachment("tags.json", "application/json", data)
	}
}

// importTagInput is the shape of each tag in an import payload
type importTagInput struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	IsPublic bool   `json:"isPublic"`
}

// ImportTagsResult is returned after an import
type ImportTagsResult struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}

// ImportTagsJSON accepts a JSON array of tags and creates missing ones
func ImportTagsJSON() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Request.ContentLength == 0 {
			return c.BadRequest(web.Map{"error": "request body is empty"})
		}

		var input []importTagInput
		if err := json.Unmarshal([]byte(c.Request.Body), &input); err != nil {
			return c.BadRequest(web.Map{"error": "invalid JSON: " + err.Error()})
		}

		if len(input) == 0 {
			return c.BadRequest(web.Map{"error": "no tags provided"})
		}

		// Fetch existing tags once so we can skip duplicates
		getAllTags := &query.GetAllTags{}
		if err := bus.Dispatch(c, getAllTags); err != nil {
			return c.Failure(err)
		}

		existing := make(map[string]*entity.Tag, len(getAllTags.Result))
		for _, t := range getAllTags.Result {
			existing[strings.ToLower(t.Name)] = t
		}

		result := ImportTagsResult{}

		for _, t := range input {
			name := strings.TrimSpace(t.Name)
			color := strings.ToUpper(strings.TrimSpace(t.Color))

			if name == "" {
				result.Errors = append(result.Errors, "tag with empty name skipped")
				result.Skipped++
				continue
			}
			if len(color) != 6 {
				result.Errors = append(result.Errors, "tag '"+name+"': color must be exactly 6 hex characters")
				result.Skipped++
				continue
			}

			if _, exists := existing[strings.ToLower(name)]; exists {
				result.Skipped++
				continue
			}

			addTag := &cmd.AddNewTag{
				Name:     name,
				Color:    color,
				IsPublic: t.IsPublic,
			}
			if err := bus.Dispatch(c, addTag); err != nil {
				if errors.Cause(err) == app.ErrNotFound {
					result.Skipped++
				} else {
					result.Errors = append(result.Errors, "tag '"+name+"': "+err.Error())
					result.Skipped++
				}
				continue
			}

			existing[strings.ToLower(name)] = addTag.Result
			result.Created++
		}

		return c.Ok(result)
	}
}
