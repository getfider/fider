package cmd

import (
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entities"
)

type SetAttachments struct {
	Post        *entities.Post
	Comment     *entities.Comment
	Attachments []*dto.ImageUpload
}

type UploadImage struct {
	Image  *dto.ImageUpload
	Folder string
}

type UploadImages struct {
	Images []*dto.ImageUpload
	Folder string
}
