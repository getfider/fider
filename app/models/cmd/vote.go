package cmd

import "github.com/getfider/fider/app/models"

type AddVote struct {
	Post *models.Post
	User *models.User
}

type RemoveVote struct {
	Post *models.Post
	User *models.User
}

type MarkPostAsDuplicate struct {
	Post     *models.Post
	Original *models.Post
}
