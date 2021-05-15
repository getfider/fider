package cmd

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/entities"
)

type AddVote struct {
	Post *models.Post
	User *entities.User
}

type RemoveVote struct {
	Post *models.Post
	User *entities.User
}

type MarkPostAsDuplicate struct {
	Post     *models.Post
	Original *models.Post
}
