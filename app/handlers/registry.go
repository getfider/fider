package handlers

import (
	"github.com/getfider/fider/app/pkg/web"
)

// Handler registry for commercial overrides
type HandlerRegistry struct {
	ModerationPage     func() web.HandlerFunc
	GetModerationItems func() web.HandlerFunc
	GetModerationCount func() web.HandlerFunc
}

var registry = &HandlerRegistry{
	ModerationPage:     ModerationPage,
	GetModerationItems: GetModerationItems,
	GetModerationCount: GetModerationCount,
}

// RegisterModerationHandlers allows commercial package to override moderation handlers
func RegisterModerationHandlers(
	moderationPage func() web.HandlerFunc,
	getModerationItems func() web.HandlerFunc,
	getModerationCount func() web.HandlerFunc,
) {
	registry.ModerationPage = moderationPage
	registry.GetModerationItems = getModerationItems
	registry.GetModerationCount = getModerationCount
}

// Handler getters that return the registered handlers
func GetModerationPageHandler() web.HandlerFunc {
	return registry.ModerationPage()
}

func GetModerationItemsHandler() web.HandlerFunc {
	return registry.GetModerationItems()
}

func GetModerationCountHandler() web.HandlerFunc {
	return registry.GetModerationCount()
}