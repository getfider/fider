package cmd

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
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
