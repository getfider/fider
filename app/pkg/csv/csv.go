package csv

import (
	"bytes"
	gocsv "encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/entity"
)

//FromPosts return a byte array of CSV file containing all posts
func FromPosts(posts []*entity.Post) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := gocsv.NewWriter(buffer)

	header := []string{
		"number",
		"title",
		"description",
		"created_at",
		"created_by",
		"votes_count",
		"comments_count",
		"status",
		"responded_by",
		"responded_at",
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
			respondedAt    string
			response       string
		)

		if post.Response != nil {
			respondedBy = post.Response.User.Name
			respondedAt = post.Response.RespondedAt.Format(time.RFC3339)
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
			post.CreatedAt.Format(time.RFC3339),
			post.User.Name,
			strconv.Itoa(post.VotesCount),
			strconv.Itoa(post.CommentsCount),
			post.Status.Name(),
			respondedBy,
			respondedAt,
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
