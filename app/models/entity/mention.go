package entity

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Mention struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	IsNew bool   `json:"isNew"`
}

type CommentString string

const mentionRegex = `@{([^{}]+)}`

func (commentString CommentString) ParseMentions() []Mention {
	r, _ := regexp.Compile(mentionRegex)

	// Remove escaped quotes from the input string
	input := strings.ReplaceAll(string(commentString), `\"`, `"`)

	matches := r.FindAllString(input, -1)

	mentions := []Mention{}

	for _, match := range matches {

		jsonMention := match[1:]

		var mention Mention
		err := json.Unmarshal([]byte(jsonMention), &mention)
		if err == nil {
			if mention.ID > 0 && mention.Name != "" {
				mentions = append(mentions, mention)
			}
		}
	}

	return mentions
}

func (mentionString CommentString) FormatMentionJson(jsonOperator func(Mention) string) string {

	r, _ := regexp.Compile(mentionRegex)

	// Remove escaped quotes from the input string
	input := strings.ReplaceAll(string(mentionString), `\"`, `"`)

	return r.ReplaceAllStringFunc(input, func(match string) string {
		jsonMention := match[1:]

		var mention Mention

		err := json.Unmarshal([]byte(jsonMention), &mention)
		if err != nil {
			return match
		}

		if mention.ID == 0 || mention.Name == "" {
			return match
		}

		return "@" + jsonOperator(mention)
	})

}
