package entity

import (
	"regexp"
)

type CommentString string

const mentionRegex = `@\[(.*?)\]`

func (commentString CommentString) ParseMentions() []string {
	r, _ := regexp.Compile(mentionRegex)
	matches := r.FindAllStringSubmatch(string(commentString), -1)

	mentions := []string{}

	for _, match := range matches {
		if len(match) >= 2 && match[1] != "" {
			mentions = append(mentions, match[1])
		}
	}

	return mentions
}

func (commentString CommentString) SanitizeMentions() string {
	r, _ := regexp.Compile(mentionRegex)
	return r.ReplaceAllString(string(commentString), "@$1")
}
