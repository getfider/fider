package query

import (
	"github.com/getfider/fider/app/models/entities"
)

type GetAttachments struct {
	Post    *entities.Post
	Comment *entities.Comment

	Result []string
}
