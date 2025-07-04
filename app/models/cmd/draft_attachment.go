package cmd

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
)

type SetDraftAttachments struct {
	DraftPost   *entity.DraftPost
	Attachments []*dto.ImageUpload
}
