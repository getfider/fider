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

func getViewData(view string) (string, []models.PostStatus, string) {
	var (
		condition string
		sort      string
	)
	statuses := []models.PostStatus{
		models.PostOpen,
		models.PostStarted,
		models.PostPlanned,
	}
	switch view {
	case "recent":
		sort = "id"
	case "my-votes":
		condition = "AND has_voted = true"
		sort = "id"
	case "most-wanted":
		sort = "votes_count"
	case "most-discussed":
		sort = "comments_count"
	case "planned":
		sort = "response_date"
		statuses = []models.PostStatus{models.PostPlanned}
	case "started":
		sort = "response_date"
		statuses = []models.PostStatus{models.PostStarted}
	case "completed":
		sort = "response_date"
		statuses = []models.PostStatus{models.PostCompleted}
	case "declined":
		sort = "response_date"
		statuses = []models.PostStatus{models.PostDeclined}
	case "all":
		sort = "id"
		statuses = []models.PostStatus{
			models.PostOpen,
			models.PostStarted,
			models.PostPlanned,
			models.PostCompleted,
			models.PostDeclined,
		}
	case "trending":
		fallthrough
	default:
		sort = "((COALESCE(recent_votes_count, 0)*5 + COALESCE(recent_comments_count, 0) *3)-1) / pow((EXTRACT(EPOCH FROM current_timestamp - created_at)/3600) + 2, 1.4)"
	}
	return condition, statuses, sort
}
