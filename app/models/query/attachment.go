package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetAttachments struct {
	Post    *entity.Post
	Comment *entity.Comment

	Result []string
}
