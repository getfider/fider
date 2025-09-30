package apiv1

import (
	"github.com/getfider/fider/app/pkg/web"
)

// ApprovePost approves a post (stub - requires commercial license)
func ApprovePost() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// DeclinePost declines (deletes) a post (stub - requires commercial license)
func DeclinePost() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// ApproveComment approves a comment (stub - requires commercial license)
func ApproveComment() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// DeclineComment declines (deletes) a comment (stub - requires commercial license)
func DeclineComment() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// DeclinePostAndBlock declines (deletes) a post and blocks the user (stub - requires commercial license)
func DeclinePostAndBlock() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// DeclineCommentAndBlock declines (deletes) a comment and blocks the user (stub - requires commercial license)
func DeclineCommentAndBlock() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// ApprovePostAndVerify approves a post and verifies the user (stub - requires commercial license)
func ApprovePostAndVerify() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}

// ApproveCommentAndVerify approves a comment and verifies the user (stub - requires commercial license)
func ApproveCommentAndVerify() web.HandlerFunc {
	return func(c *web.Context) error {
		return c.BadRequest(web.Map{"error": "Content moderation requires commercial license"})
	}
}
