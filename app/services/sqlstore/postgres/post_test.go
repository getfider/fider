package postgres_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestPostStorage_GetAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	now := time.Now()

	_, err := trx.Execute("INSERT INTO posts (title, slug, number, description, created_at, tenant_id, user_id, status) VALUES ('add twitter integration', 'add-twitter-integration', 1, 'Would be great to see it integrated with twitter', $1, 1, 1, 1)", now)
	Expect(err).IsNil()

	_, err = trx.Execute("INSERT INTO posts (title, slug, number, description, created_at, tenant_id, user_id, status) VALUES ('this is my post', 'this-is-my-post', 2, 'no description', $1, 1, 2, 2)", now)
	Expect(err).IsNil()

	allPosts := &query.GetAllPosts{}
	err = bus.Dispatch(demoTenantCtx, allPosts)
	Expect(err).IsNil()
	Expect(allPosts.Result).HasLen(2)

	Expect(allPosts.Result[0].Title).Equals("this is my post")
	Expect(allPosts.Result[0].Slug).Equals("this-is-my-post")
	Expect(allPosts.Result[0].Number).Equals(2)
	Expect(allPosts.Result[0].Description).Equals("no description")
	Expect(allPosts.Result[0].User.Name).Equals("Arya Stark")
	Expect(allPosts.Result[0].VotesCount).Equals(0)
	Expect(allPosts.Result[0].Status).Equals(enum.PostCompleted)

	Expect(allPosts.Result[1].Title).Equals("add twitter integration")
	Expect(allPosts.Result[1].Slug).Equals("add-twitter-integration")
	Expect(allPosts.Result[1].Number).Equals(1)
	Expect(allPosts.Result[1].Description).Equals("Would be great to see it integrated with twitter")
	Expect(allPosts.Result[1].User.Name).Equals("Jon Snow")
	Expect(allPosts.Result[1].VotesCount).Equals(0)
	Expect(allPosts.Result[1].Status).Equals(enum.PostStarted)

	search10 := &query.SearchPosts{Query: "twitter", Limit: "10"}
	search0 := &query.SearchPosts{Query: "twitter", Limit: "0"}
	err = bus.Dispatch(demoTenantCtx, search10, search0)
	Expect(err).IsNil()

	Expect(search10.Result).HasLen(1)
	Expect(search10.Result[0].Slug).Equals("add-twitter-integration")
	Expect(search0.Result).HasLen(0)
}

func TestPostStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	postByID := &query.GetPostByID{PostID: newPost.Result.ID}
	postBySlug := &query.GetPostBySlug{Slug: "my-new-post"}
	err = bus.Dispatch(jonSnowCtx, postByID, postBySlug)
	Expect(err).IsNil()

	Expect(postByID.Result.ID).Equals(newPost.Result.ID)
	Expect(postByID.Result.Number).Equals(1)
	Expect(postByID.Result.HasVoted).IsFalse()
	Expect(postByID.Result.VotesCount).Equals(0)
	Expect(postByID.Result.Status).Equals(enum.PostOpen)
	Expect(postByID.Result.Title).Equals("My new post")
	Expect(postByID.Result.Description).Equals("with this description")
	Expect(postByID.Result.User.ID).Equals(1)
	Expect(postByID.Result.User.Name).Equals("Jon Snow")
	Expect(postByID.Result.User.Email).Equals("jon.snow@got.com")

	Expect(postBySlug.Result.ID).Equals(newPost.Result.ID)
	Expect(postBySlug.Result.Number).Equals(1)
	Expect(postBySlug.Result.HasVoted).IsFalse()
	Expect(postBySlug.Result.VotesCount).Equals(0)
	Expect(postBySlug.Result.Status).Equals(enum.PostOpen)
	Expect(postBySlug.Result.Title).Equals("My new post")
	Expect(postBySlug.Result.Description).Equals("with this description")
	Expect(postBySlug.Result.User.ID).Equals(1)
	Expect(postBySlug.Result.User.Name).Equals("Jon Snow")
	Expect(postBySlug.Result.User.Email).Equals("jon.snow@got.com")
}

func TestPostStorage_GetInvalid(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	postByID := &query.GetPostByID{PostID: 1}
	err := bus.Dispatch(jonSnowCtx, postByID)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(postByID.Result).IsNil()
}

func TestPostStorage_AddAndReturnComments(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddNewComment{Post: newPost.Result, Content: "Comment #1"})
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.AddNewComment{Post: newPost.Result, Content: "Comment #2"})
	Expect(err).IsNil()

	commentsByPost := &query.GetCommentsByPost{Post: newPost.Result}
	err = bus.Dispatch(aryaStarkCtx, commentsByPost)
	Expect(err).IsNil()
	Expect(commentsByPost.Result).HasLen(2)

	Expect(commentsByPost.Result[0].Content).Equals("Comment #1")
	Expect(commentsByPost.Result[0].User.Name).Equals("Jon Snow")
	Expect(commentsByPost.Result[1].Content).Equals("Comment #2")
	Expect(commentsByPost.Result[1].User.Name).Equals("Arya Stark")
}

func TestPostStorage_AddGetUpdateComment(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNewComment := &cmd.AddNewComment{Post: newPost.Result, Content: "Comment #1"}
	err = bus.Dispatch(jonSnowCtx, addNewComment)
	Expect(err).IsNil()

	commentByID := &query.GetCommentByID{CommentID: addNewComment.Result.ID}
	err = bus.Dispatch(jonSnowCtx, commentByID)
	Expect(err).IsNil()

	Expect(commentByID.Result.ID).Equals(addNewComment.Result.ID)
	Expect(commentByID.Result.Content).Equals("Comment #1")
	Expect(commentByID.Result.User.ID).Equals(jonSnow.ID)
	Expect(commentByID.Result.EditedAt).IsNil()
	Expect(commentByID.Result.EditedBy).IsNil()

	updateComment := &cmd.UpdateComment{CommentID: addNewComment.Result.ID, Content: "Comment #1 with edit"}
	err = bus.Dispatch(aryaStarkCtx, updateComment)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, commentByID)
	Expect(err).IsNil()

	Expect(commentByID.Result.ID).Equals(addNewComment.Result.ID)
	Expect(commentByID.Result.Content).Equals("Comment #1 with edit")
	Expect(commentByID.Result.User.ID).Equals(jonSnow.ID)
	Expect(commentByID.Result.EditedAt).IsNotNil()
	Expect(commentByID.Result.EditedBy.ID).Equals(aryaStark.ID)
}

func TestPostStorage_AddDeleteComment(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNewComment := &cmd.AddNewComment{Post: newPost.Result, Content: "Comment #1"}
	err = bus.Dispatch(jonSnowCtx, addNewComment)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.DeleteComment{CommentID: addNewComment.Result.ID})
	Expect(err).IsNil()

	commentByID := &query.GetCommentByID{CommentID: addNewComment.Result.ID}
	err = bus.Dispatch(jonSnowCtx, commentByID)
	Expect(err).Equals(app.ErrNotFound)
	Expect(commentByID.Result).IsNil()
}

func TestPostStorage_AddAndGet_DifferentTenants(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	jonSnowNewPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, jonSnowNewPost)
	Expect(err).IsNil()

	tonyStarkNewPost := &cmd.AddNewPost{Title: "My other post", Description: "with other description"}
	err = bus.Dispatch(tonyStarkCtx, tonyStarkNewPost)
	Expect(err).IsNil()

	getDemoTenantPost1 := &query.GetPostByNumber{Number: 1}
	err = bus.Dispatch(demoTenantCtx, getDemoTenantPost1)
	Expect(err).IsNil()
	Expect(getDemoTenantPost1.Result.ID).Equals(jonSnowNewPost.Result.ID)
	Expect(getDemoTenantPost1.Result.Number).Equals(1)
	Expect(getDemoTenantPost1.Result.Title).Equals("My new post")
	Expect(getDemoTenantPost1.Result.Slug).Equals("my-new-post")

	getAvengersPost1 := &query.GetPostByNumber{Number: 1}
	err = bus.Dispatch(avengersTenantCtx, getAvengersPost1)
	Expect(err).IsNil()
	Expect(getAvengersPost1.Result.ID).Equals(tonyStarkNewPost.Result.ID)
	Expect(getAvengersPost1.Result.Number).Equals(1)
	Expect(getAvengersPost1.Result.Title).Equals("My other post")
	Expect(getAvengersPost1.Result.Slug).Equals("my-other-post")
}

func TestPostStorage_Update(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.UpdatePost{Post: newPost.Result, Title: "The new title", Description: "With the new description"})
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.Title).Equals("The new title")
	Expect(getPost.Result.Description).Equals("With the new description")
	Expect(getPost.Result.Slug).Equals("the-new-title")
}

func TestPostStorage_AddVote(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: newPost.Result, User: aryaStark})
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.HasVoted).IsFalse()
	Expect(getPost.Result.VotesCount).Equals(1)

	err = bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: newPost.Result, User: jonSnow})
	Expect(err).IsNil()

	getPost = &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.HasVoted).IsTrue()
	Expect(getPost.Result.VotesCount).Equals(2)
}

func TestPostStorage_AddVote_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
	)
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.HasVoted).IsTrue()
	Expect(getPost.Result.VotesCount).Equals(1)
}

func TestPostStorage_RemoveVote(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
		&cmd.RemoveVote{Post: newPost.Result, User: jonSnow},
	)
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.HasVoted).IsFalse()
	Expect(getPost.Result.VotesCount).Equals(0)
}

func TestPostStorage_RemoveVote_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
		&cmd.RemoveVote{Post: newPost.Result, User: jonSnow},
		&cmd.RemoveVote{Post: newPost.Result, User: jonSnow},
	)
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.HasVoted).IsFalse()
	Expect(getPost.Result.VotesCount).Equals(0)
}

func TestPostStorage_SetResponse(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post", Status: enum.PostStarted})
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.Response.Text).Equals("We liked this post")
	Expect(getPost.Result.Status).Equals(enum.PostStarted)
	Expect(getPost.Result.Response.User.ID).Equals(1)
}

func TestPostStorage_SetResponse_KeepOpen(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post", Status: enum.PostOpen})
	Expect(err).IsNil()
}

func TestPostStorage_SetResponse_ChangeText(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	getPost := &query.GetPostByID{PostID: newPost.Result.ID}

	bus.MustDispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post", Status: enum.PostStarted})
	bus.MustDispatch(jonSnowCtx, getPost)
	firstResponseAt := getPost.Result.Response.RespondedAt

	bus.MustDispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post and we'll work on it", Status: enum.PostStarted})
	bus.MustDispatch(jonSnowCtx, getPost)
	Expect(getPost.Result.Response.RespondedAt).Equals(firstResponseAt)

	bus.MustDispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "We finished it", Status: enum.PostCompleted})
	bus.MustDispatch(jonSnowCtx, getPost)
	Expect(getPost.Result.Response.RespondedAt).TemporarilySimilar(firstResponseAt, time.Second)
}

func TestPostStorage_SetResponse_AsDuplicate(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost1 := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost1)
	Expect(err).IsNil()

	newPost2 := &cmd.AddNewPost{Title: "My other post", Description: "with similar description"}
	err = bus.Dispatch(aryaStarkCtx, newPost2)
	Expect(err).IsNil()

	err = bus.Dispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: newPost1.Result, User: jonSnow},
		&cmd.AddVote{Post: newPost2.Result, User: aryaStark},
		&cmd.MarkPostAsDuplicate{Post: newPost2.Result, Original: newPost1.Result},
	)
	Expect(err).IsNil()

	getPost1 := &query.GetPostByID{PostID: newPost1.Result.ID}
	getPost2 := &query.GetPostByID{PostID: newPost2.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, getPost1, getPost2)
	Expect(err).IsNil()

	Expect(getPost1.Result.VotesCount).Equals(2)
	Expect(getPost1.Result.Status).Equals(enum.PostOpen)
	Expect(getPost1.Result.Response).IsNil()

	Expect(getPost2.Result.Response.Text).Equals("")
	Expect(getPost2.Result.VotesCount).Equals(1)
	Expect(getPost2.Result.Status).Equals(enum.PostDuplicate)
	Expect(getPost2.Result.Response.User.ID).Equals(1)
	Expect(getPost2.Result.Response.Original.Number).Equals(newPost1.Result.Number)
	Expect(getPost2.Result.Response.Original.Title).Equals(newPost1.Result.Title)
	Expect(getPost2.Result.Response.Original.Slug).Equals(newPost1.Result.Slug)
	Expect(getPost2.Result.Response.Original.Status).Equals(newPost1.Result.Status)
}

func TestPostStorage_SetResponse_AsDeleted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	bus.MustDispatch(jonSnowCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "Spam!", Status: enum.PostDeleted})

	postByID := &query.GetPostByID{PostID: newPost.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, postByID)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(postByID.Result).IsNil()

	postByNumber := &query.GetPostByNumber{Number: newPost.Result.Number}
	err = bus.Dispatch(aryaStarkCtx, postByNumber)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(postByNumber.Result).IsNil()
}

func TestPostStorage_AddVote_ClosedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx,
		&cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post", Status: enum.PostCompleted},
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
	)
	Expect(err).IsNil()

	getPost := &query.GetPostByNumber{Number: newPost.Result.Number}
	err = bus.Dispatch(aryaStarkCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.VotesCount).Equals(0)
}

func TestPostStorage_RemoveVote_ClosedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	bus.MustDispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: newPost.Result, User: jonSnow},
		&cmd.SetPostResponse{Post: newPost.Result, Text: "We liked this post", Status: enum.PostCompleted},
		&cmd.RemoveVote{Post: newPost.Result, User: jonSnow},
	)

	getPost := &query.GetPostByNumber{Number: newPost.Result.Number}
	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.VotesCount).Equals(1)
}

func TestPostStorage_WithTags(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	addBug := &cmd.AddNewTag{Name: "Bug", Color: "FF0000", IsPublic: true}
	addFeatureRequest := &cmd.AddNewTag{Name: "Feature Request", Color: "00FF00", IsPublic: false}
	bus.MustDispatch(jonSnowCtx, addBug, addFeatureRequest)
	bus.MustDispatch(jonSnowCtx, &cmd.AssignTag{Tag: addBug.Result, Post: newPost.Result})
	bus.MustDispatch(jonSnowCtx, &cmd.AssignTag{Tag: addFeatureRequest.Result, Post: newPost.Result})

	getPost := &query.GetPostByNumber{Number: newPost.Result.Number}
	err = bus.Dispatch(aryaStarkCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.Tags).HasLen(1)
	Expect(getPost.Result.Tags[0]).Equals(addBug.Result.Slug)

	err = bus.Dispatch(jonSnowCtx, getPost)
	Expect(err).IsNil()
	Expect(getPost.Result.Tags).HasLen(2)
	Expect(getPost.Result.Tags[0]).Equals(addBug.Result.Slug)
	Expect(getPost.Result.Tags[1]).Equals(addFeatureRequest.Result.Slug)
}

func TestPostStorage_IsReferenced(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost1 := &cmd.AddNewPost{Title: "My first post", Description: "with this description"}
	newPost2 := &cmd.AddNewPost{Title: "My second post", Description: "with this description"}
	newPost3 := &cmd.AddNewPost{Title: "My third post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost1, newPost2, newPost3)
	Expect(err).IsNil()

	bus.MustDispatch(jonSnowCtx, &cmd.MarkPostAsDuplicate{Post: newPost2.Result, Original: newPost3.Result})
	bus.MustDispatch(jonSnowCtx, &cmd.MarkPostAsDuplicate{Post: newPost3.Result, Original: newPost1.Result})

	isReferenced1 := &query.PostIsReferenced{PostID: newPost1.Result.ID}
	isReferenced2 := &query.PostIsReferenced{PostID: newPost2.Result.ID}
	isReferenced3 := &query.PostIsReferenced{PostID: newPost3.Result.ID}

	err = bus.Dispatch(jonSnowCtx, isReferenced1, isReferenced2, isReferenced3)
	Expect(err).IsNil()
	Expect(isReferenced1.Result).IsTrue()
	Expect(isReferenced2.Result).IsFalse()
	Expect(isReferenced3.Result).IsTrue()
}

func TestPostStorage_ListVotesOfPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	bus.MustDispatch(jonSnowCtx, &cmd.AddVote{Post: newPost.Result, User: jonSnow})
	bus.MustDispatch(jonSnowCtx, &cmd.AddVote{Post: newPost.Result, User: aryaStark})

	listVotes := &query.ListPostVotes{PostID: newPost.Result.ID, IncludeEmail: true}
	err = bus.Dispatch(jonSnowCtx, listVotes)
	Expect(err).IsNil()
	Expect(listVotes.Result).HasLen(2)

	Expect(listVotes.Result[0].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(listVotes.Result[0].User.Name).Equals("Jon Snow")
	Expect(listVotes.Result[0].User.Email).Equals("jon.snow@got.com")

	Expect(listVotes.Result[1].CreatedAt).TemporarilySimilar(time.Now(), 5*time.Second)
	Expect(listVotes.Result[1].User.Name).Equals("Arya Stark")
	Expect(listVotes.Result[1].User.Email).Equals("arya.stark@got.com")
}

func TestPostStorage_Attachments(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost1 := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	newPost2 := &cmd.AddNewPost{Title: "My other post", Description: "with another description"}
	err := bus.Dispatch(jonSnowCtx, newPost1, newPost2)
	Expect(err).IsNil()

	getAttachments1 := &query.GetAttachments{Post: newPost1.Result}
	getAttachments2 := &query.GetAttachments{Post: newPost2.Result}

	err = bus.Dispatch(jonSnowCtx, getAttachments1)
	Expect(err).IsNil()
	Expect(getAttachments1.Result).HasLen(0)

	bytes, err := ioutil.ReadFile(env.Path("favicon.png"))
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.SetAttachments{
		Post: newPost1.Result,
		Attachments: []*dto.ImageUpload{
			{
				BlobKey: "12345-test.png",
				Upload: &dto.ImageUploadData{
					FileName:    "test.png",
					ContentType: "image/png",
					Content:     bytes,
				},
			},
		},
	})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, getAttachments1, getAttachments2)
	Expect(err).IsNil()
	Expect(getAttachments1.Result).HasLen(1)
	Expect(getAttachments1.Result[0]).Equals("12345-test.png")
	Expect(getAttachments2.Result).HasLen(0)

	err = bus.Dispatch(jonSnowCtx, &cmd.SetAttachments{
		Post: newPost2.Result,
		Attachments: []*dto.ImageUpload{
			{
				BlobKey: "12345-test.png",
				Remove:  true,
			},
			{
				BlobKey: "67890-test2.png",
				Upload: &dto.ImageUploadData{
					FileName:    "test2.png",
					ContentType: "image/png",
					Content:     bytes,
				},
			},
			{
				BlobKey: "67890-test6.png",
				Upload: &dto.ImageUploadData{
					FileName:    "test6.png",
					ContentType: "image/png",
					Content:     bytes,
				},
			},
		},
	})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, getAttachments1, getAttachments2)
	Expect(err).IsNil()
	Expect(getAttachments1.Result).HasLen(1)
	Expect(getAttachments1.Result[0]).Equals("12345-test.png")
	Expect(getAttachments2.Result).HasLen(2)
	Expect(getAttachments2.Result[0]).Equals("67890-test2.png")
	Expect(getAttachments2.Result[1]).Equals("67890-test6.png")

	err = bus.Dispatch(jonSnowCtx, &cmd.SetAttachments{
		Post: newPost1.Result,
		Attachments: []*dto.ImageUpload{
			{
				BlobKey: "12345-test.png",
				Remove:  true,
			},
		},
	})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, getAttachments1)
	Expect(err).IsNil()
	Expect(getAttachments1.Result).HasLen(0)
}
