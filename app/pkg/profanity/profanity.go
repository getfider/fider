package profanity

import (
	"context"
	"regexp"
	"strings"

	"github.com/Spicy-Bush/fider-tarkov-community/app"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/errors"
)

// ContainsProfanity checks if the provided text contains any banned words.
// It reads the banned words from the current tenant record
// and uses a regex pattern for each banned word.
func ContainsProfanity(ctx context.Context, text string) ([]string, error) {
	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if !ok {
		return nil, errors.New("tenant not found in context")
	}

	words := strings.Split(tenant.ProfanityWords, ",")
	lowerText := strings.ToLower(text)
	found := []string{}

	for _, word := range words {
		trimmed := strings.TrimSpace(word)
		if trimmed == "" {
			continue
		}

		// Construct a regex pattern that matches:
		// 1. The banned word as a standalone word (using word boundaries \b)
		// 2. OR two or more consecutive repetitions of that word.
		// 3. We can improve tis if needed.
		//
		// Explanation:
		//   (?i)           => case-insensitive mode.
		//   \b<word>\b     => matches the word with boundaries.
		//   (?:<word>){2,} => matches two or more consecutive occurrences.
		pattern := `(?i)(\b` + regexp.QuoteMeta(trimmed) + `\b|(?:` + regexp.QuoteMeta(trimmed) + `){2,})`

		// Compile the regex.
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		// If the pattern matches anywhere in the text, record the banned word.
		if re.FindString(lowerText) != "" {
			found = append(found, trimmed)
		}
	}
	// Return the word if for whatever reason it is needed to be used.
	return found, nil
}
