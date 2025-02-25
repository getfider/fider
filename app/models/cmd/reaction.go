package cmd

import "github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"

type ToggleCommentReaction struct {
	Comment *entity.Comment
	Emoji   string
	User    *entity.User
	Result  bool
}
