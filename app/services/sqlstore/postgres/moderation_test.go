package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestModerationStorage_ApproveOnlyReportsPendingItems(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	bus.MustDispatch(aryaStarkCtx, newPost)
	newComment := &cmd.AddNewComment{Post: newPost.Result, Content: "Comment #1"}
	bus.MustDispatch(aryaStarkCtx, newComment)

	_, err := trx.Execute("UPDATE posts SET is_approved = false WHERE id = $1", newPost.Result.ID)
	Expect(err).IsNil()
	_, err = trx.Execute("UPDATE comments SET is_approved = false WHERE id = $1", newComment.Result.ID)
	Expect(err).IsNil()

	approvePost := &cmd.ApprovePost{PostID: newPost.Result.ID}
	bus.MustDispatch(jonSnowCtx, approvePost)
	Expect(approvePost.Result).IsTrue()

	approvePost = &cmd.ApprovePost{PostID: newPost.Result.ID}
	bus.MustDispatch(jonSnowCtx, approvePost)
	Expect(approvePost.Result).IsFalse()

	approveComment := &cmd.ApproveComment{CommentID: newComment.Result.ID}
	bus.MustDispatch(jonSnowCtx, approveComment)
	Expect(approveComment.Result).Equals(newPost.Result.ID)

	approveComment = &cmd.ApproveComment{CommentID: newComment.Result.ID}
	bus.MustDispatch(jonSnowCtx, approveComment)
	Expect(approveComment.Result).Equals(0)
}
