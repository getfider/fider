package postgres

import (
	"context"
	"fmt"
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
	input = replaceOr.Replace(allowedTextRunes.ReplaceAllString(input, " "))
	return strings.Join(strings.Fields(input), "|")
}

// SanitizeString converts input to another string that only contains utf-8 characters and not-null characters
func SanitizeString(input string) string {
	input = strings.ReplaceAll(input, "\u0000", "")
	return strings.ToValidUTF8(input, "")
}

// MapLocaleToTSConfig maps a tenant's locale short key to the corresponding PostgreSQL text search configuration.
// Returns 'simple' if no match is found or if PostgreSQL doesn't have native support for the language.
// All locale definitions are centralized in app/models/enum/locale.go
func MapLocaleToTSConfig(locale string) string {
	return enum.MapLocaleToTSConfig(locale)
}

// ViewSelector lets the caller decide whether the SearchPosts SQL filters by
// concrete status slugs (when the user explicitly chose them) or by semantic
// kinds (the default/all home views, which automatically pick up custom
// statuses an admin added with kind=active or kind=closed-*).
type ViewSelector struct {
	Condition string
	SlugSet   []string // non-empty when filtering by exact slugs
	KindSet   []string // non-empty when filtering by kind via JOINed statuses
	Sort      string
}

func getViewData(query query.SearchPosts, tagsPlaceholder int) ViewSelector {
	v := ViewSelector{}

	// Explicit slug list from the user — slug match wins, kind ignored.
	if len(query.Statuses) > 0 {
		v.SlugSet = query.Statuses
	}

	if query.MyVotesOnly {
		v.Condition = "AND has_voted = true"
	}

	switch query.View {
	case "recent":
		v.Sort = "id"
	case "most-wanted":
		v.Sort = "votes_count"
	case "most-discussed":
		v.Sort = "comments_count"
	case "my-votes":
		// Deprecated: You can instead filter on my votes only for more flexibility than using this view.
		v.Condition = "AND has_voted = true"
		v.Sort = "id"
	case "planned":
		// Deprecated: Use status filters instead
		v.Sort = "response_date"
		v.SlugSet = []string{"planned"}
	case "started":
		// Deprecated: Use status filters instead
		v.Sort = "response_date"
		v.SlugSet = []string{"started"}
	case "completed":
		// Deprecated: Use status filters instead
		v.Sort = "response_date"
		v.SlugSet = []string{"completed"}
	case "declined":
		// Deprecated: Use status filters instead
		v.Sort = "response_date"
		v.SlugSet = []string{"declined"}
	case "all":
		v.Sort = "id"
		if len(v.SlugSet) == 0 {
			// All non-deleted statuses regardless of admin-added kinds.
			v.KindSet = []string{"open", "active", "closed-completed", "closed-declined", "duplicate"}
		}
	case "trending":
		fallthrough
	default:
		v.Sort = "((COALESCE(recent_votes_count, 0)*5 + COALESCE(recent_comments_count, 0) *3)-1) / pow((EXTRACT(EPOCH FROM current_timestamp - created_at)/3600) + 2, 1.4)"
	}

	// Default home view picks up Open + any tenant-defined "active" kind, so
	// admin-added custom in-progress statuses surface without a code change.
	if len(v.SlugSet) == 0 && len(v.KindSet) == 0 {
		v.KindSet = []string{"open", "active"}
	}

	if query.NoTagsOnly {
		// NoTagsOnly takes precedence: combining "untagged" with specific tag
		// filters produces contradictory SQL that always returns zero rows.
		query.Tags = nil
		v.Condition += " AND tags = '{}'"
	}

	if len(query.Tags) > 0 {
		v.Condition += fmt.Sprintf(" AND tags && $%d", tagsPlaceholder)
	}
	return v
}

// buildStatusFilter converts a ViewSelector into the SQL fragment that goes
// into buildPostQuery's caller-supplied WHERE clause, plus the matching
// $2-array parameter. Slug match uses posts.status_slug directly; kind match
// uses the JOINed statuses table (alias ps in buildPostQuery).
func buildStatusFilter(v ViewSelector) (string, []string) {
	if len(v.SlugSet) > 0 {
		return "p.status_slug = ANY($2)", v.SlugSet
	}
	return "ps.kind = ANY($2)", v.KindSet
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
