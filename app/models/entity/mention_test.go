package entity_test

import (
	"testing"

	"github.com/getfider/fider/app/models/entity"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestComment_ParseMentions(t *testing.T) {
	RegisterT(t)
	tests := []struct {
		name     string
		content  string
		expected []entity.Mention
	}{
		{
			name:     "no mentions",
			content:  "This is a regular comment \\\" just here",
			expected: []entity.Mention{},
		},
		{
			name:     "Simple mention",
			content:  `Hello there @{"id":1,"name":"John Doe","isNew":false} how are you`,
			expected: []entity.Mention{{ID: 1, Name: "John Doe", IsNew: false}},
		},
		{
			name:     "Simple mention 2",
			content:  `Hello there @{"id":2,"name":"John Doe Smith","isNew":true} how are you`,
			expected: []entity.Mention{{ID: 2, Name: "John Doe Smith", IsNew: true}},
		},
		{
			name:     "Multiple mentions",
			content:  `Hello there @{"id":2,"name":"John Doe Smith","isNew":true} and @{"id":1,"name":"John Doe","isNew":false} how are you`,
			expected: []entity.Mention{{ID: 2, Name: "John Doe Smith", IsNew: true}, {ID: 1, Name: "John Doe", IsNew: false}},
		},
		{
			name:     "Some odd JSON",
			content:  `Hello there @{"id":2,name:"John Doe Smith","isNew":true} how are you`,
			expected: []entity.Mention{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := &entity.Comment{Content: tt.content}
			comment.ParseMentions()
			Expect(comment.Mentions).Equals(tt.expected)
		})
	}
}

func TestStripMentionMetaData(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		`@{\"id\":1,\"name\":\"John Doe quoted\"}`:                                 "@John Doe quoted",
		`@{"id":1,"name":"John Doe"}`:                                              "@John Doe",
		`@{"id":1,"name":"JohnDoe"}`:                                               "@JohnDoe",
		`@{\"id\":1,\"name\":\"JohnDoe quoted\"}`:                                  "@JohnDoe quoted",
		`@{"id":1,"name":"John Smith Doe"}`:                                        "@John Smith Doe",
		`@{\"id\":1,\"name\":\"John Smith Doe quoted\"}`:                           "@John Smith Doe quoted",
		"Hello there how are you":                                                  "Hello there how are you",
		`Hello there @{"id":1,"name":"John Doe"}`:                                  "Hello there @John Doe",
		`Hello there @{"id":1,"name":"John Doe quoted"}`:                           "Hello there @John Doe quoted",
		`Hello both @{"id":1,"name":"John Doe"} and @{"id":2,"name":"John Smith"}`: "Hello both @John Doe and @John Smith",
	} {
		output := entity.CommentString(input).FormatMentionJson(func(mention entity.Mention) string {
			return mention.Name
		})
		Expect(output).Equals(expected)
	}
}

func TestStripMentionMetaDataDoesntBreakUserInput(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		`There is nothing here`:                                         "There is nothing here",
		`There is nothing here {ok}`:                                    "There is nothing here {ok}",
		`This is a message for {{matt}}`:                                "This is a message for {{matt}}",
		`This is a message for {{id:1,wiggles:true}}`:                   "This is a message for {{id:1,wiggles:true}}",
		`Although uncommon, someone could enter @{something} like this`: "Although uncommon, someone could enter @{something} like this",
		`Or @{"id":100,"wiggles":"yes"} something like this`:            `Or @{"id":100,"wiggles":"yes"} something like this`,
	} {
		output := entity.CommentString(input).FormatMentionJson(func(mention entity.Mention) string {
			return mention.Name
		})
		Expect(output).Equals(expected)
	}
}
