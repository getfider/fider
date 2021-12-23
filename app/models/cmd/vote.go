package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

type AddVote struct {
	Post *entity.Post
	User *entity.User
}

type RemoveVote struct {
	Post *entity.Post
	User *entity.User
}

type MarkPostAsDuplicate struct {
	Post     *entity.Post
	Original *entity.Post
}
