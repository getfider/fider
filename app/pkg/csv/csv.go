package csv

import (
	"bytes"
	gocsv "encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models"
)

//FromPosts return a byte array of CSV file containing all posts
func FromPosts(posts []*models.Post) ([]byte, error) {
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

	for _, post := range posts {
		var (
			originalNumber string
			originalTitle  string
			respondedBy    string
			respondedOn    string
			response       string
		)

		if post.Response != nil {
			respondedBy = post.Response.User.Name
			respondedOn = post.Response.RespondedOn.Format(time.RFC3339)
			response = post.Response.Text
			if post.Response.Original != nil {
				originalNumber = strconv.Itoa(post.Response.Original.Number)
				originalTitle = post.Response.Original.Title
			}
		}

		record := []string{
			strconv.Itoa(post.Number),
			post.Title,
			post.Description,
			post.CreatedOn.Format(time.RFC3339),
			post.User.Name,
			strconv.Itoa(post.TotalSupporters),
			strconv.Itoa(post.TotalComments),
			models.GetPostStatusName(post.Status),
			respondedBy,
			respondedOn,
			response,
			originalNumber,
			originalTitle,
			strings.Join(post.Tags, ", "),
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
