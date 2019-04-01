package cmd

import "github.com/getfider/fider/app/models"

type SetPostResponse struct {
	Post   *models.Post
	Text   string
	Status models.PostStatus
}
