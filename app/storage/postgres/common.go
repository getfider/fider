package postgres

import (
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models"
)

var onlyalphanumeric = regexp.MustCompile("[^a-zA-Z0-9 |]+")
var replaceOr = strings.NewReplacer("|", " ")

// ToTSQuery converts input to another string that can be safely used for ts_query
func ToTSQuery(input string) string {
	input = replaceOr.Replace(onlyalphanumeric.ReplaceAllString(input, ""))
	return strings.Join(strings.Fields(input), "|")
}

func getFilterData(filter string) (string, []int, string) {
	var (
		condition string
		sort      string
	)
	statuses := []int{
		models.PostOpen,
		models.PostStarted,
		models.PostPlanned,
	}
	switch filter {
	case "recent":
		sort = "id"
	case "my-votes":
		condition = "AND viewer_voted = true"
		sort = "id"
	case "most-wanted":
		sort = "votes"
	case "most-discussed":
		sort = "comments"
	case "planned":
		sort = "response_date"
		statuses = []int{models.PostPlanned}
	case "started":
		sort = "response_date"
		statuses = []int{models.PostStarted}
	case "completed":
		sort = "response_date"
		statuses = []int{models.PostCompleted}
	case "declined":
		sort = "response_date"
		statuses = []int{models.PostDeclined}
	case "all":
		sort = "id"
		statuses = []int{
			models.PostOpen,
			models.PostStarted,
			models.PostPlanned,
			models.PostCompleted,
			models.PostDeclined,
		}
	case "trending":
		fallthrough
	default:
		sort = "((COALESCE(recent_votes, 0)*5 + COALESCE(recent_comments, 0) *3)-1) / pow((EXTRACT(EPOCH FROM current_timestamp - created_on)/3600) + 2, 1.4)"
	}
	return condition, statuses, sort
}
