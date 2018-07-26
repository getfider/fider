package csv_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/csv"
)

func TestExportPostsToCSV_Empty(t *testing.T) {
	RegisterT(t)

	posts := []*models.Idea{}
	expected, err := ioutil.ReadFile("./testdata/empty.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

func TestExportPostsToCSV_OneIdea(t *testing.T) {
	RegisterT(t)

	posts := []*models.Idea{
		declinedIdea,
	}

	expected, err := ioutil.ReadFile("./testdata/one-idea.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

func TestExportPostsToCSV_MorePosts(t *testing.T) {
	RegisterT(t)

	posts := []*models.Idea{
		declinedIdea,
		openIdea,
		duplicateIdea,
	}

	expected, err := ioutil.ReadFile("./testdata/more-posts.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

var declinedIdea = &models.Idea{
	Number:      10,
	Title:       "Go is fast",
	Description: "Very tiny description",
	CreatedOn:   time.Date(2018, 3, 23, 19, 33, 22, 0, time.UTC),
	User: &models.User{
		Name: "Faceless",
	},
	TotalSupporters: 4,
	TotalComments:   2,
	Status:          models.IdeaDeclined,
	Response: &models.IdeaResponse{
		Text:        "Nothing we need to do",
		RespondedOn: time.Date(2018, 4, 4, 19, 48, 10, 0, time.UTC),
		User: &models.User{
			Name: "John Snow",
		},
	},
	Tags: []string{"easy", "ignored"},
}

var openIdea = &models.Idea{
	Number:      15,
	Title:       "Go is great",
	Description: "",
	CreatedOn:   time.Date(2018, 2, 21, 15, 51, 35, 0, time.UTC),
	User: &models.User{
		Name: "Someone else",
	},
	TotalSupporters: 4,
	TotalComments:   2,
	Status:          models.IdeaOpen,
}

var duplicateIdea = &models.Idea{
	Number:      20,
	Title:       "Go is easy",
	Description: "",
	CreatedOn:   time.Date(2018, 1, 12, 1, 46, 59, 0, time.UTC),
	User: &models.User{
		Name: "Faceless",
	},
	TotalSupporters: 4,
	TotalComments:   2,
	Status:          models.IdeaDuplicate,
	Response: &models.IdeaResponse{
		Text:        "This has already been suggested",
		RespondedOn: time.Date(2018, 3, 17, 10, 15, 42, 0, time.UTC),
		User: &models.User{
			Name: "Arya Stark",
		},
		Original: &models.OriginalIdea{
			Number: 99,
			Title:  "Go is very easy",
		},
	},
	Tags: []string{"this-tag-has,comma"},
}
