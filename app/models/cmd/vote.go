package cmd

import (
	"github.com/getfider/fider/app/models/entities"
)

type AddVote struct {
	Post *entities.Post
	User *entities.User
}

type RemoveVote struct {
	Post *entities.Post
	User *entities.User
}

type MarkPostAsDuplicate struct {
	Post     *entities.Post
	Original *entities.Post
}
