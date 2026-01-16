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

type TrustUser struct {
	UserID int
}
