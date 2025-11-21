package apiv1

import (
	"github.com/getfider/fider/app/pkg/web"
)

// Handler registry for commercial API v1 overrides
type HandlerRegistry struct {
	ApprovePost             func() web.HandlerFunc
	DeclinePost             func() web.HandlerFunc
	ApproveComment          func() web.HandlerFunc
	DeclineComment          func() web.HandlerFunc
	DeclinePostAndBlock     func() web.HandlerFunc
	DeclineCommentAndBlock  func() web.HandlerFunc
	ApprovePostAndVerify    func() web.HandlerFunc
	ApproveCommentAndVerify func() web.HandlerFunc
}

var registry = &HandlerRegistry{
	ApprovePost:             ApprovePost,
	DeclinePost:             DeclinePost,
	ApproveComment:          ApproveComment,
	DeclineComment:          DeclineComment,
	DeclinePostAndBlock:     DeclinePostAndBlock,
	DeclineCommentAndBlock:  DeclineCommentAndBlock,
	ApprovePostAndVerify:    ApprovePostAndVerify,
	ApproveCommentAndVerify: ApproveCommentAndVerify,
}

// RegisterModerationHandlers allows commercial package to override moderation handlers
func RegisterModerationHandlers(
	approvePost func() web.HandlerFunc,
	declinePost func() web.HandlerFunc,
	approveComment func() web.HandlerFunc,
	declineComment func() web.HandlerFunc,
	declinePostAndBlock func() web.HandlerFunc,
	declineCommentAndBlock func() web.HandlerFunc,
	approvePostAndVerify func() web.HandlerFunc,
	approveCommentAndVerify func() web.HandlerFunc,
) {
	registry.ApprovePost = approvePost
	registry.DeclinePost = declinePost
	registry.ApproveComment = approveComment
	registry.DeclineComment = declineComment
	registry.DeclinePostAndBlock = declinePostAndBlock
	registry.DeclineCommentAndBlock = declineCommentAndBlock
	registry.ApprovePostAndVerify = approvePostAndVerify
	registry.ApproveCommentAndVerify = approveCommentAndVerify
}

// Handler getters that return the registered handlers
func GetApprovePostHandler() web.HandlerFunc {
	return registry.ApprovePost()
}

func GetDeclinePostHandler() web.HandlerFunc {
	return registry.DeclinePost()
}

func GetApproveCommentHandler() web.HandlerFunc {
	return registry.ApproveComment()
}

func GetDeclineCommentHandler() web.HandlerFunc {
	return registry.DeclineComment()
}

func GetDeclinePostAndBlockHandler() web.HandlerFunc {
	return registry.DeclinePostAndBlock()
}

func GetDeclineCommentAndBlockHandler() web.HandlerFunc {
	return registry.DeclineCommentAndBlock()
}

func GetApprovePostAndVerifyHandler() web.HandlerFunc {
	return registry.ApprovePostAndVerify()
}

func GetApproveCommentAndVerifyHandler() web.HandlerFunc {
	return registry.ApproveCommentAndVerify()
}