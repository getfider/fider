package cmd

import "github.com/getfider/fider/app/models"

type SetAttachments struct {
	Post        *models.Post
	Comment     *models.Comment
	Attachments []*models.ImageUpload
}
