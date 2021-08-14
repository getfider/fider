package webhook

import (
	"context"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/webhook"
	"time"
)

var dummyPost = &entity.Post{
	ID:     36,
	Number: 36,
	Title:  "Example dummy post title",
	Slug:   "example-dummy-post-title",
	Description: `*Those data are fake.*

This is an example of a post, but don't worry, nothing was created on your Fider instance:
It was just a test of a webhook trigger.

You can use ` + "`" + `{{ markdown .post_description }}` + "`" + ` to parse __Markdown__ to **HTML** in post!`,
	CreatedAt: time.Date(2021, time.May, 7, 18, 42, 27, 0, time.UTC),
	User: &entity.User{
		ID:    7,
		Name:  "Fider",
		Email: "contact@fider.io",
		Role:  1,
	},
	HasVoted:      true,
	VotesCount:    7,
	CommentsCount: 3,
	Status:        enum.PostStarted,
	Response: &entity.PostResponse{
		Text:        "This is a response, still in *Markdown*.",
		RespondedAt: time.Date(2021, time.July, 9, 15, 29, 57, 0, time.UTC),
		User:        nil,
	},
	Tags: []string{"tag1", "tag2"},
}

func dummyTriggerProps(c context.Context, webhookType enum.WebhookType) webhook.Props {
	props := webhook.Props{}
	author := c.Value(app.UserCtxKey).(*entity.User)
	tenant := c.Value(app.TenantCtxKey).(*entity.Tenant)
	baseURL, logoURL := web.BaseURL(c), web.LogoURL(c)
	dummyPost.User.AvatarURL = logoURL
	dummyPost.Response.User = author
	props.SetUser(author, "author")
	props.SetTenant(tenant, "tenant", baseURL, logoURL)
	switch webhookType {
	case enum.WebhookNewPost:
		props.SetPost(dummyPost, "post", baseURL, false, false)
	case enum.WebhookNewComment:
		props.SetPost(dummyPost, "post", baseURL, true, true)
		props["comment"] = "An example **comment** on a post."
	case enum.WebhookChangeStatus:
		props.SetPost(dummyPost, "post", baseURL, true, true)
		props["post_old_status"] = enum.PostOpen.Name()
	case enum.WebhookDeletePost:
		props.SetPost(dummyPost, "post", baseURL, true, true)
		props["post_status"] = enum.PostDeleted.Name()
		props["post_response_text"] = "The reason _why_ this post was deleted."
	}
	return props
}
