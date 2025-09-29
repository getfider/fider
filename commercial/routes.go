package commercial

import (
	commercialHandlers "github.com/getfider/fider/commercial/handlers"
	commercialApiv1 "github.com/getfider/fider/commercial/handlers/apiv1"
	"github.com/getfider/fider/app/pkg/web"
)

// Override system for commercial routes

var (
	// Commercial handler overrides
	CommercialModerationPage     = commercialHandlers.ModerationPage
	CommercialGetModerationItems = commercialHandlers.GetModerationItems
	CommercialGetModerationCount = commercialHandlers.GetModerationCount

	// Commercial API v1 handler overrides
	CommercialApprovePost             = commercialApiv1.ApprovePost
	CommercialDeclinePost             = commercialApiv1.DeclinePost
	CommercialApproveComment          = commercialApiv1.ApproveComment
	CommercialDeclineComment          = commercialApiv1.DeclineComment
	CommercialDeclinePostAndBlock     = commercialApiv1.DeclinePostAndBlock
	CommercialDeclineCommentAndBlock  = commercialApiv1.DeclineCommentAndBlock
	CommercialApprovePostAndVerify    = commercialApiv1.ApprovePostAndVerify
	CommercialApproveCommentAndVerify = commercialApiv1.ApproveCommentAndVerify
)

// GetModerationHandler returns the appropriate moderation handler based on license
func GetModerationHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialModerationPage
}

// GetModerationItemsHandler returns the appropriate moderation items handler based on license
func GetModerationItemsHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialGetModerationItems
}

// GetModerationCountHandler returns the appropriate moderation count handler based on license
func GetModerationCountHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialGetModerationCount
}

// GetApprovePostHandler returns the appropriate approve post handler based on license
func GetApprovePostHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialApprovePost
}

// GetDeclinePostHandler returns the appropriate decline post handler based on license
func GetDeclinePostHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialDeclinePost
}

// GetApproveCommentHandler returns the appropriate approve comment handler based on license
func GetApproveCommentHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialApproveComment
}

// GetDeclineCommentHandler returns the appropriate decline comment handler based on license
func GetDeclineCommentHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialDeclineComment
}

// GetDeclinePostAndBlockHandler returns the appropriate decline post and block handler based on license
func GetDeclinePostAndBlockHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialDeclinePostAndBlock
}

// GetDeclineCommentAndBlockHandler returns the appropriate decline comment and block handler based on license
func GetDeclineCommentAndBlockHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialDeclineCommentAndBlock
}

// GetApprovePostAndVerifyHandler returns the appropriate approve post and verify handler based on license
func GetApprovePostAndVerifyHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialApprovePostAndVerify
}

// GetApproveCommentAndVerifyHandler returns the appropriate approve comment and verify handler based on license
func GetApproveCommentAndVerifyHandler(fallback func() web.HandlerFunc) func() web.HandlerFunc {
	return CommercialApproveCommentAndVerify
}