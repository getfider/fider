package cmd

type ApprovePost struct {
	PostID int

	// Result is true if the post was pending and is now approved
	Result bool
}

type DeclinePost struct {
	PostID int
}

type ApproveComment struct {
	CommentID int

	// Result is the comment's post ID if the comment was pending and is now approved, otherwise 0
	Result int
}

type DeclineComment struct {
	CommentID int
}

type BulkApproveItems struct {
	PostIDs    []int
	CommentIDs []int
}

type BulkDeclineItems struct {
	PostIDs    []int
	CommentIDs []int
}

type TrustUser struct {
	UserID int
}
