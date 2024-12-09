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
