package cmd

import "github.com/getfider/fider/app/models/entity"

type ToggleCommentReaction struct {
	Comment *entity.Comment
	Emoji   string
	User    *entity.User
	Result  bool
}
