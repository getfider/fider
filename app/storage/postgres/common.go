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

func getFilterData(filter string) ([]int, string) {
	var sort string
	statuses := []int{
		models.IdeaOpen,
		models.IdeaStarted,
		models.IdeaPlanned,
	}
	switch filter {
	case "recent":
		sort = "id"
	case "most-wanted":
		sort = "supporters"
	case "most-discussed":
		sort = "comments"
	case "planned":
		sort = "response_date"
		statuses = []int{models.IdeaPlanned}
	case "started":
		sort = "response_date"
		statuses = []int{models.IdeaStarted}
	case "completed":
		sort = "response_date"
		statuses = []int{models.IdeaCompleted}
	case "declined":
		sort = "response_date"
		statuses = []int{models.IdeaDeclined}
	case "all":
		sort = "id"
		statuses = []int{
			models.IdeaOpen,
			models.IdeaStarted,
			models.IdeaPlanned,
			models.IdeaCompleted,
			models.IdeaDeclined,
		}
	case "trending":
		fallthrough
	default:
		sort = "((COALESCE(recent_supporters, 0)*5 + COALESCE(recent_comments, 0) *3)-1) / pow((EXTRACT(EPOCH FROM current_timestamp - created_on)/3600) + 2, 1.4)"
	}
	return statuses, sort
}
