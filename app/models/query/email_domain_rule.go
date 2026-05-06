package query

import "github.com/getfider/fider/app/models/entity"

type GetEmailDomainRules struct {
	Result struct {
		Deny  []*entity.EmailDomainRule
		Allow []*entity.EmailDomainRule
	}
}

// GetDisposableUsers returns users whose email matches the current
// disposable-blocking rules (bundled list + tenant deny - tenant allow).
type GetDisposableUsers struct {
	Limit  int
	Result struct {
		Total int
		Users []*GetDisposableUsersRow
	}
}

type GetDisposableUsersRow struct {
	UserID       int    `json:"userID"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	VoteCount    int    `json:"voteCount"`
	PostCount    int    `json:"postCount"`
	CommentCount int    `json:"commentCount"`
}
