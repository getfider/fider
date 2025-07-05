package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetAttachments struct {
	Post    *entity.Post
	Comment *entity.Comment

	Result []string
}

// IsAttachmentReferenced checks if a blob key is referenced in the attachments table
type IsAttachmentReferenced struct {
	BlobKey string

	Result bool
}
