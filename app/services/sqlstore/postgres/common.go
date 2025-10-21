package postgres

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/web"
)

// allowedTextRunes matches any character that is NOT a Unicode letter (\p{L}), number (\p{N}), space, or pipe.
// This allows both Latin and non-Latin scripts to be preserved for full-text search.
var allowedTextRunes = regexp.MustCompile(`[^\p{L}\p{N} |]+`)
var replaceOr = strings.NewReplacer("|", " ")

// ToTSQuery converts input to another string that can be safely used for ts_query
func ToTSQuery(input string) string {
	input = replaceOr.Replace(allowedTextRunes.ReplaceAllString(input, ""))
	return strings.Join(strings.Fields(input), "|")
}

// SanitizeString converts input to another string that only contains utf-8 characters and not-null characters
func SanitizeString(input string) string {
	input = strings.Replace(input, "\u0000", "", -1)
	return strings.ToValidUTF8(input, "")
}

// mapLocaleToTSConfig maps a tenant's locale short key to the corresponding PostgreSQL text search configuration.
// See the /locale folder for supported locales. Returns 'simple' if no match is found.
func MapLocaleToTSConfig(locale string) string {
	switch locale {
	case "ar":
		return "arabic"
	case "de":
		return "german"
	case "en":
		return "english"
	case "es-ES":
		return "spanish"
	case "fr":
		return "french"
	case "it":
		return "italian"
	case "nl":
		return "dutch"
	case "pt-BR":
		return "portuguese"
	case "ru":
		return "russian"
	case "sv-SE":
		return "swedish"
	case "tr":
		return "turkish"
	case "cs":
	case "el":
	case "fa":
	case "ja":
	case "ko":
	case "pl":
	case "si-LK":
	case "sk":
	case "zh-CN":
	default:
		return "simple"
	}

	return "simple"
}

func getViewData(query query.SearchPosts) (string, []enum.PostStatus, string) {
	var (
		condition string
		sort      string
	)
	statusFilters := query.Statuses
	if len(statusFilters) == 0 {
		// Use a sensible default list of status filters
		statusFilters = []enum.PostStatus{
			enum.PostOpen,
			enum.PostStarted,
			enum.PostPlanned,
		}
	}

	if query.MyVotesOnly {
		condition = "AND has_voted = true"
	}

	switch query.View {
	case "recent":
		sort = "id"
	case "most-wanted":
		sort = "votes_count"
	case "most-discussed":
		sort = "comments_count"
	case "my-votes":
		// Depracated: You can instead filter on my votes only for more flexibility than using this view.
		condition = "AND has_voted = true"
		sort = "id"
	case "planned":
		// Depracated: Use status filters instead
		sort = "response_date"
		statusFilters = []enum.PostStatus{enum.PostPlanned}
	case "started":
		// Depracated: Use status filters instead
		sort = "response_date"
		statusFilters = []enum.PostStatus{enum.PostStarted}
	case "completed":
		// Depracated: Use status filters instead
		sort = "response_date"
		statusFilters = []enum.PostStatus{enum.PostCompleted}
	case "declined":
		// Depracated: Use status filters instead
		sort = "response_date"
		statusFilters = []enum.PostStatus{enum.PostDeclined}
	case "all":
		sort = "id"
		statusFilters = []enum.PostStatus{
			enum.PostOpen,
			enum.PostStarted,
			enum.PostPlanned,
			enum.PostCompleted,
			enum.PostDeclined,
		}
	case "trending":
		fallthrough
	default:
		sort = "((COALESCE(recent_votes_count, 0)*5 + COALESCE(recent_comments_count, 0) *3)-1) / pow((EXTRACT(EPOCH FROM current_timestamp - created_at)/3600) + 2, 1.4)"
	}

	if query.NoTagsOnly {
		condition += " AND tags = '{}'"
	}

	if len(query.Tags) > 0 {
		condition += " AND tags && $3"
	}
	return condition, statusFilters, sort
}

func buildAvatarURL(ctx context.Context, avatarType enum.AvatarType, id int, name, avatarBlobKey string) string {
	if name == "" {
		name = "-"
	}

	if avatarType == enum.AvatarTypeCustom {
		return web.AssetsURL(ctx, "/static/images/%s", avatarBlobKey)
	} else {
		return web.AssetsURL(ctx, "/static/avatars/%s/%d/%s", avatarType.String(), id, url.PathEscape(name))
	}
}
