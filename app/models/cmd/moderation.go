package cmd

type ApprovePost struct {
	PostID int
}

type DeclinePost struct {
	PostID int
}

type ApproveComment struct {
	CommentID int
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

type DeclinePostAndBlock struct {
	PostID int
}

type DeclineCommentAndBlock struct {
	CommentID int
}

type ApprovePostAndVerify struct {
	PostID int
}

type ApproveCommentAndVerify struct {
	CommentID int
}

type VerifyUser struct {
	UserID int
}