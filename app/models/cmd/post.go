package cmd

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/rand"
)

type AddNewPost struct {
	Title       string
	Description string

	Result *entity.Post
}

type AddNewDraftPost struct {
	Title       string
	Description string
	Code        string

	Result *entity.DraftPost
}

func GenerateNewCode() string {
	return rand.String(12)
}

type UpdatePost struct {
	Post        *entity.Post
	Title       string
	Description string

	Result *entity.Post
}

type SetPostResponse struct {
	Post   *entity.Post
	Text   string
	Status enum.PostStatus
}
