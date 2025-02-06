package query

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
)

type GetAttachments struct {
	Post    *entity.Post
	Comment *entity.Comment

	Result []string
}
