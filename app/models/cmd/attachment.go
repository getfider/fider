package cmd

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/entities"
)

type SetAttachments struct {
	Post        *entities.Post
	Comment     *entities.Comment
	Attachments []*models.ImageUpload
}

type UploadImage struct {
	Image  *models.ImageUpload
	Folder string
}

type UploadImages struct {
	Images []*models.ImageUpload
	Folder string
}
