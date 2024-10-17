package postgres_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/sqlstore/postgres"
)

func TestToTSQuery(t *testing.T) {
	RegisterT(t)

	var testcases = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"123 hello", "123|hello"},
		{" hello  ", "hello"},
		{" hello$ world$ ", "hello|world"},
		{" yes, please ", "yes|please"},
		{" yes / please ", "yes|please"},
		{" hello 'world' ", "hello|world"},
		{"hello|world", "hello|world"},
		{"hello | world", "hello|world"},
		{"hello & world", "hello|world"},
	}

	for _, testcase := range testcases {
		output := postgres.ToTSQuery(testcase.input)
		Expect(output).Equals(testcase.expected)
	}
}

func withTenant(ctx context.Context, tenant *entity.Tenant) context.Context {
	return context.WithValue(ctx, app.TenantCtxKey, tenant)
}

func withUser(ctx context.Context, user *entity.User) context.Context {
	ctx = context.WithValue(ctx, app.TenantCtxKey, user.Tenant)
	ctx = context.WithValue(ctx, app.UserCtxKey, user)
	return ctx
}

func TestToggleReaction_Add(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	newComment := &cmd.AddNewComment{Post: newPost.Result, Content: "This is my comment"}
	err = bus.Dispatch(jonSnowCtx, newComment)
	Expect(err).IsNil()

	// Now add a reaction
	reaction := &cmd.ToggleCommentReaction{Comment: newComment.Result, Emoji: "👍", User: jonSnow}
	err = bus.Dispatch(jonSnowCtx, reaction)
	Expect(err).IsNil()
	Expect(reaction.Result).IsTrue()

	// Get the comment, and check that the reaction was added
	commentByID := &query.GetCommentsByPost{Post: &entity.Post{ID: newPost.Result.ID}}
	err = bus.Dispatch(jonSnowCtx, commentByID)
	Expect(err).IsNil()

	Expect(commentByID.Result).IsNotNil()
	Expect(commentByID.Result[0].ReactionCounts).IsNotNil()
	Expect(len(commentByID.Result[0].ReactionCounts)).Equals(1)
	Expect(commentByID.Result[0].ReactionCounts[0].Emoji).Equals("👍")
	Expect(commentByID.Result[0].ReactionCounts[0].Count).Equals(1)
	Expect(commentByID.Result[0].ReactionCounts[0].IncludesMe).IsTrue()

	// Now remove the reaction
	reaction = &cmd.ToggleCommentReaction{Comment: newComment.Result, Emoji: "👍", User: jonSnow}
	err = bus.Dispatch(jonSnowCtx, reaction)
	Expect(err).IsNil()
	Expect(reaction.Result).IsFalse()

	// Get the comment, and check that the reaction was removed
	commentByID = &query.GetCommentsByPost{Post: &entity.Post{ID: newPost.Result.ID}}
	err = bus.Dispatch(jonSnowCtx, commentByID)
	Expect(err).IsNil()

	Expect(commentByID.Result).IsNotNil()
	Expect(commentByID.Result[0].ReactionCounts).IsNil()

}
