package csv

import (
	"bytes"
	gocsv "encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models"
)

//FromIdeas return a byte array of CSV file containing all ideas
func FromIdeas(ideas []*models.Idea) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := gocsv.NewWriter(buffer)

	header := []string{
		"number",
		"title",
		"description",
		"created_on",
		"created_by",
		"total_supporters",
		"total_comments",
		"status",
		"responded_by",
		"responded_on",
		"response",
		"original_number",
		"original_title",
		"tags",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, idea := range ideas {
		var (
			originalNumber string
			originalTitle  string
			respondedBy    string
			respondedOn    string
			response       string
		)

		if idea.Response != nil {
			respondedBy = idea.Response.User.Name
			respondedOn = idea.Response.RespondedOn.Format(time.RFC3339)
			response = idea.Response.Text
			if idea.Response.Original != nil {
				originalNumber = strconv.Itoa(idea.Response.Original.Number)
				originalTitle = idea.Response.Original.Title
			}
		}

		record := []string{
			strconv.Itoa(idea.Number),
			idea.Title,
			idea.Description,
			idea.CreatedOn.Format(time.RFC3339),
			idea.User.Name,
			strconv.Itoa(idea.TotalSupporters),
			strconv.Itoa(idea.TotalComments),
			models.GetIdeaStatusName(idea.Status),
			respondedBy,
			respondedOn,
			response,
			originalNumber,
			originalTitle,
			strings.Join(idea.Tags, ", "),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
