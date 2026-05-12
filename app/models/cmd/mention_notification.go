package cmd

type AddMentionNotification struct {
	UserID    int `db:"user_id"`
	CommentID int `db:"comment_id"`
	PostID    int `db:"post_id"`
}
