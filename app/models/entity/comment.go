package entity

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

type ReactionCounts struct {
	Emoji      string `json:"emoji"`
	Count      int    `json:"count"`
	IncludesMe bool   `json:"includesMe"`
}

type Mention struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	IsNew bool   `json:"isNew"`
}

// Comment represents an user comment on an post
type Comment struct {
	ID             int              `json:"id"`
	Content        string           `json:"content"`
	CreatedAt      time.Time        `json:"createdAt"`
	User           *User            `json:"user"`
	Attachments    []string         `json:"attachments,omitempty"`
	EditedAt       *time.Time       `json:"editedAt,omitempty"`
	EditedBy       *User            `json:"editedBy,omitempty"`
	ReactionCounts []ReactionCounts `json:"reactionCounts,omitempty"`
	Mentions       []Mention        `json:"_"`
}

func (c *Comment) ParseMentions() {
	r, _ := regexp.Compile("@{([^}]+)}")

	// Remove escaped quotes from the input string
	input := strings.ReplaceAll(c.Content, `\"`, `"`)

	matches := r.FindAllString(input, -1)

	c.Mentions = []Mention{}

	for _, match := range matches {

		jsonMention := match[1:]

		var dat map[string]interface{}

		err := json.Unmarshal([]byte(jsonMention), &dat)
		if err == nil {
			name, nameExists := dat["name"]
			id, idExists := dat["id"]
			isNew, isNewExists := dat["isNew"]
			if !isNewExists {
				isNew = false
			}

			if nameExists && idExists {
				mention := Mention{
					ID:    int(id.(float64)),
					Name:  name.(string),
					IsNew: isNew.(bool),
				}
				c.Mentions = append(c.Mentions, mention)
			}
		}
	}
}
