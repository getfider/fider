package postgres_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
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

	trx.Execute("INSERT INTO posts (title, slug, number, description, created_at, tenant_id, user_id, status) VALUES ('add twitter integration', 'add-twitter-integration', 1, 'Would be great to see it integrated with twitter', $1, 1, 1, 1)", now)
	trx.Execute("INSERT INTO posts (title, slug, number, description, created_at, tenant_id, user_id, status) VALUES ('this is my post', 'this-is-my-post', 2, 'no description', $1, 1, 2, 2)", now)

	posts.SetCurrentTenant(demoTenant)

	dbPosts, err := posts.GetAll()
	Expect(err).IsNil()
	Expect(dbPosts).HasLen(2)

	Expect(dbPosts[0].Title).Equals("this is my post")
	Expect(dbPosts[0].Slug).Equals("this-is-my-post")
	Expect(dbPosts[0].Number).Equals(2)
	Expect(dbPosts[0].Description).Equals("no description")
	Expect(dbPosts[0].User.Name).Equals("Arya Stark")
	Expect(dbPosts[0].VotesCount).Equals(0)
	Expect(dbPosts[0].Status).Equals(models.PostCompleted)

	Expect(dbPosts[1].Title).Equals("add twitter integration")
	Expect(dbPosts[1].Slug).Equals("add-twitter-integration")
	Expect(dbPosts[1].Number).Equals(1)
	Expect(dbPosts[1].Description).Equals("Would be great to see it integrated with twitter")
	Expect(dbPosts[1].User.Name).Equals("Jon Snow")
	Expect(dbPosts[1].VotesCount).Equals(0)
	Expect(dbPosts[1].Status).Equals(models.PostStarted)

	dbPosts, err = posts.Search("twitter", "trending", "10", []string{})
	Expect(err).IsNil()
	Expect(dbPosts).HasLen(1)
	Expect(dbPosts[0].Slug).Equals("add-twitter-integration")

	dbPosts, err = posts.Search("twitter", "trending", "0", []string{})
	Expect(err).IsNil()
	Expect(dbPosts).HasLen(0)
}

func TestPostStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	dbPostByID, err := posts.GetByID(post.ID)
	Expect(err).IsNil()
	Expect(dbPostByID.ID).Equals(post.ID)
	Expect(dbPostByID.Number).Equals(1)
	Expect(dbPostByID.HasVoted).IsFalse()
	Expect(dbPostByID.VotesCount).Equals(0)
	Expect(dbPostByID.Status).Equals(models.PostOpen)
	Expect(dbPostByID.Title).Equals("My new post")
	Expect(dbPostByID.Description).Equals("with this description")
	Expect(dbPostByID.User.ID).Equals(1)
	Expect(dbPostByID.User.Name).Equals("Jon Snow")
	Expect(dbPostByID.User.Email).Equals("jon.snow@got.com")

	dbPostBySlug, err := posts.GetBySlug("my-new-post")

	Expect(err).IsNil()
	Expect(dbPostBySlug.ID).Equals(post.ID)
	Expect(dbPostBySlug.Number).Equals(1)
	Expect(dbPostBySlug.HasVoted).IsFalse()
	Expect(dbPostBySlug.VotesCount).Equals(0)
	Expect(dbPostBySlug.Status).Equals(models.PostOpen)
	Expect(dbPostBySlug.Title).Equals("My new post")
	Expect(dbPostBySlug.Description).Equals("with this description")
	Expect(dbPostBySlug.User.ID).Equals(1)
	Expect(dbPostBySlug.User.Name).Equals("Jon Snow")
	Expect(dbPostBySlug.User.Email).Equals("jon.snow@got.com")
}

func TestPostStorage_GetInvalid(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)

	dbPost, err := posts.GetByID(1)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(dbPost).IsNil()
}

func TestPostStorage_AddAndReturnComments(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	posts.SetCurrentUser(jonSnow)
	posts.AddComment(post, "Comment #1")
	posts.SetCurrentUser(aryaStark)
	posts.AddComment(post, "Comment #2")

	comments, err := posts.GetCommentsByPost(post)
	Expect(err).IsNil()
	Expect(comments).HasLen(2)

	Expect(comments[0].Content).Equals("Comment #1")
	Expect(comments[0].User.Name).Equals("Jon Snow")
	Expect(comments[1].Content).Equals("Comment #2")
	Expect(comments[1].User.Name).Equals("Arya Stark")
}

func TestPostStorage_AddGetUpdateComment(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	commentID, err := posts.AddComment(post, "Comment #1")
	Expect(err).IsNil()

	comment, err := posts.GetCommentByID(commentID)
	Expect(err).IsNil()
	Expect(comment.ID).Equals(commentID)
	Expect(comment.Content).Equals("Comment #1")
	Expect(comment.User.ID).Equals(jonSnow.ID)
	Expect(comment.EditedAt).IsNil()
	Expect(comment.EditedBy).IsNil()

	posts.SetCurrentUser(aryaStark)
	err = posts.UpdateComment(commentID, "Comment #1 with edit")
	Expect(err).IsNil()

	comment, err = posts.GetCommentByID(commentID)
	Expect(err).IsNil()
	Expect(comment.ID).Equals(commentID)
	Expect(comment.Content).Equals("Comment #1 with edit")
	Expect(comment.User.ID).Equals(jonSnow.ID)
	Expect(comment.EditedAt).IsNotNil()
	Expect(comment.EditedBy.ID).Equals(aryaStark.ID)
}

func TestPostStorage_AddDeleteComment(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	commentID, err := posts.AddComment(post, "Comment #1")
	Expect(err).IsNil()

	err = posts.DeleteComment(commentID)
	Expect(err).IsNil()

	comment, err := posts.GetCommentByID(commentID)
	Expect(err).Equals(app.ErrNotFound)
	Expect(comment).IsNil()
}

func TestPostStorage_AddAndGet_DifferentTenants(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	demoPost, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	posts.SetCurrentTenant(avengersTenant)
	posts.SetCurrentUser(tonyStark)
	avengersPost, err := posts.Add("My other post", "with other description")
	Expect(err).IsNil()

	posts.SetCurrentTenant(demoTenant)
	dbPost, err := posts.GetByNumber(1)

	Expect(err).IsNil()
	Expect(dbPost.ID).Equals(demoPost.ID)
	Expect(dbPost.Number).Equals(1)
	Expect(dbPost.Title).Equals("My new post")
	Expect(dbPost.Slug).Equals("my-new-post")

	posts.SetCurrentTenant(avengersTenant)
	dbPost, err = posts.GetByNumber(1)

	Expect(err).IsNil()
	Expect(dbPost.ID).Equals(avengersPost.ID)
	Expect(dbPost.Number).Equals(1)
	Expect(dbPost.Title).Equals("My other post")
	Expect(dbPost.Slug).Equals("my-other-post")
}

func TestPostStorage_Update(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	post, err = posts.Update(post, "The new comment", "With the new description")
	Expect(err).IsNil()
	Expect(post.Title).Equals("The new comment")
	Expect(post.Description).Equals("With the new description")
	Expect(post.Slug).Equals("the-new-comment")
}

func TestPostStorage_AddVote(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	users.SetCurrentTenant(demoTenant)
	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: aryaStark})
	Expect(err).IsNil()

	dbPost, err := posts.GetByNumber(1)
	Expect(dbPost.HasVoted).IsFalse()
	Expect(dbPost.VotesCount).Equals(1)

	err = bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	dbPost, err = posts.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbPost.HasVoted).IsTrue()
	Expect(dbPost.VotesCount).Equals(2)
}

func TestPostStorage_AddVote_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")

	err := bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	dbPost, err := posts.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbPost.VotesCount).Equals(1)
}

func TestPostStorage_RemoveVote(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")

	err := bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.RemoveVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	dbPost, err := posts.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbPost.VotesCount).Equals(0)
}

func TestPostStorage_RemoveVote_Twice(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")

	err := bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.RemoveVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.RemoveVote{Post: post, User: jonSnow})
	Expect(err).IsNil()

	dbPost, err := posts.GetByNumber(1)
	Expect(err).IsNil()
	Expect(dbPost.VotesCount).Equals(0)
}

func TestPostStorage_SetResponse(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")
	err := bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "We liked this post", Status: models.PostStarted})
	Expect(err).IsNil()

	post, _ = posts.GetByID(post.ID)
	Expect(post.Response.Text).Equals("We liked this post")
	Expect(post.Status).Equals(models.PostStarted)
	Expect(post.Response.User.ID).Equals(1)
}

func TestPostStorage_SetResponse_KeepOpen(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")
	err := bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "We liked this post", Status: models.PostOpen})
	Expect(err).IsNil()
}

func TestPostStorage_SetResponse_ChangeText(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post, _ := posts.Add("My new post", "with this description")
	bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "We liked this post", Status: models.PostStarted})
	post, _ = posts.GetByID(post.ID)
	respondedAt := post.Response.RespondedAt

	bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "We liked this post and we'll work on it", Status: models.PostStarted})
	post, _ = posts.GetByID(post.ID)
	Expect(post.Response.RespondedAt).Equals(respondedAt)

	bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "We finished it", Status: models.PostCompleted})
	post, _ = posts.GetByID(post.ID)
	Expect(post.Response.RespondedAt).TemporarilySimilar(respondedAt, time.Second)
}

func TestPostStorage_SetResponse_AsDuplicate(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post1, _ := posts.Add("My new post", "with this description")
	err := bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post1, User: jonSnow})
	Expect(err).IsNil()

	posts.SetCurrentUser(aryaStark)
	post2, _ := posts.Add("My other post", "with similar description")
	err = bus.Dispatch(aryaStarkCtx, &cmd.AddVote{Post: post2, User: aryaStark})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.MarkPostAsDuplicate{Post: post2, Original: post1})
	Expect(err).IsNil()

	posts.SetCurrentUser(jonSnow)
	post1, _ = posts.GetByID(post1.ID)

	Expect(post1.VotesCount).Equals(2)
	Expect(post1.Status).Equals(models.PostOpen)
	Expect(post1.Response).IsNil()

	post2, _ = posts.GetByID(post2.ID)

	Expect(post2.Response.Text).Equals("")
	Expect(post2.VotesCount).Equals(1)
	Expect(post2.Status).Equals(models.PostDuplicate)
	Expect(post2.Response.User.ID).Equals(1)
	Expect(post2.Response.Original.Number).Equals(post1.Number)
	Expect(post2.Response.Original.Title).Equals(post1.Title)
	Expect(post2.Response.Original.Slug).Equals(post1.Slug)
	Expect(post2.Response.Original.Status).Equals(post1.Status)
}

func TestPostStorage_SetResponse_AsDeleted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, err := posts.Add("My new post", "with this description")
	Expect(err).IsNil()

	bus.Dispatch(jonSnowCtx, &cmd.SetPostResponse{Post: post, Text: "Spam!", Status: models.PostDeleted})

	post1, err := posts.GetByNumber(post.Number)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(post1).IsNil()

	post2, err := posts.GetByID(post.ID)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(post2).IsNil()
}

func TestPostStorage_AddVote_ClosedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, _ := posts.Add("My new post", "with this description")

	err := bus.Dispatch(jonSnowCtx,
		&cmd.SetPostResponse{Post: post, Text: "We liked this post", Status: models.PostCompleted},
		&cmd.AddVote{Post: post, User: jonSnow},
	)
	Expect(err).IsNil()

	dbPost, err := posts.GetByNumber(post.Number)
	Expect(err).IsNil()
	Expect(dbPost.VotesCount).Equals(0)
}

func TestPostStorage_RemoveVote_ClosedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, _ := posts.Add("My new post", "with this description")

	bus.Dispatch(
		jonSnowCtx,
		&cmd.AddVote{Post: post, User: jonSnow},
		&cmd.SetPostResponse{Post: post, Text: "We liked this post", Status: models.PostCompleted},
		&cmd.RemoveVote{Post: post, User: jonSnow},
	)

	dbPost, err := posts.GetByNumber(post.Number)
	Expect(err).IsNil()
	Expect(dbPost.VotesCount).Equals(1)
}

func TestPostStorage_WithTags(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)

	post, _ := posts.Add("My new post", "with this description")

	addBug := &cmd.AddNewTag{Name: "Bug", Color: "FF0000", IsPublic: true}
	bus.Dispatch(jonSnowCtx, addBug)

	addFeatureRequest := &cmd.AddNewTag{Name: "Feature Request", Color: "00FF00", IsPublic: false}
	bus.Dispatch(jonSnowCtx, addFeatureRequest)

	bus.Dispatch(jonSnowCtx, &cmd.AssignTag{Tag: addBug.Result, Post: post})
	bus.Dispatch(jonSnowCtx, &cmd.AssignTag{Tag: addFeatureRequest.Result, Post: post})

	post, _ = posts.GetByID(post.ID)
	Expect(post.Tags).HasLen(1)
	Expect(post.Tags[0]).Equals(addBug.Result.Slug)

	posts.SetCurrentUser(jonSnow)
	post, _ = posts.GetByID(post.ID)
	Expect(post.Tags).HasLen(2)
	Expect(post.Tags[0]).Equals(addBug.Result.Slug)
	Expect(post.Tags[1]).Equals(addFeatureRequest.Result.Slug)
}

func TestPostStorage_IsReferenced(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post1, _ := posts.Add("My first post", "with this description")
	post2, _ := posts.Add("My second post", "with this description")
	post3, _ := posts.Add("My third post", "with this description")

	bus.Dispatch(jonSnowCtx, &cmd.MarkPostAsDuplicate{Post: post2, Original: post3})
	bus.Dispatch(jonSnowCtx, &cmd.MarkPostAsDuplicate{Post: post3, Original: post1})

	referenced, err := posts.IsReferenced(post1)
	Expect(referenced).IsTrue()
	Expect(err).IsNil()

	referenced, err = posts.IsReferenced(post2)
	Expect(referenced).IsFalse()
	Expect(err).IsNil()

	referenced, err = posts.IsReferenced(post3)
	Expect(referenced).IsTrue()
	Expect(err).IsNil()
}

func TestPostStorage_ListVotesOfPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post1, _ := posts.Add("My new post", "with this description")

	bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post1, User: jonSnow})
	bus.Dispatch(jonSnowCtx, &cmd.AddVote{Post: post1, User: aryaStark})

	listVotes := &query.ListPostVotes{PostID: post1.ID}
	err := bus.Dispatch(jonSnowCtx, listVotes)
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

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post1, _ := posts.Add("My new post", "with this description")
	post2, _ := posts.Add("My other post", "with another description")

	attachments, err := posts.GetAttachments(post1, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(0)

	bytes, err := ioutil.ReadFile(env.Path("favicon.png"))
	Expect(err).IsNil()

	err = posts.SetAttachments(post1, nil, []*models.ImageUpload{
		&models.ImageUpload{
			BlobKey: "12345-test.png",
			Upload: &models.ImageUploadData{
				FileName:    "test.png",
				ContentType: "image/png",
				Content:     bytes,
			},
		},
	})
	Expect(err).IsNil()

	attachments, err = posts.GetAttachments(post1, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(1)
	Expect(attachments[0]).Equals("12345-test.png")

	attachments, err = posts.GetAttachments(post2, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(0)

	err = posts.SetAttachments(post2, nil, []*models.ImageUpload{
		&models.ImageUpload{
			BlobKey: "12345-test.png",
			Remove:  true,
		},
		&models.ImageUpload{
			BlobKey: "67890-test2.png",
			Upload: &models.ImageUploadData{
				FileName:    "test2.png",
				ContentType: "image/png",
				Content:     bytes,
			},
		},
		&models.ImageUpload{
			BlobKey: "67890-test6.png",
			Upload: &models.ImageUploadData{
				FileName:    "test6.png",
				ContentType: "image/png",
				Content:     bytes,
			},
		},
	})
	Expect(err).IsNil()

	attachments, err = posts.GetAttachments(post1, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(1)
	Expect(attachments[0]).Equals("12345-test.png")

	attachments, err = posts.GetAttachments(post2, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(2)
	Expect(attachments[0]).Equals("67890-test2.png")
	Expect(attachments[1]).Equals("67890-test6.png")

	err = posts.SetAttachments(post1, nil, []*models.ImageUpload{
		&models.ImageUpload{
			BlobKey: "12345-test.png",
			Remove:  true,
		},
	})
	Expect(err).IsNil()

	attachments, err = posts.GetAttachments(post1, nil)
	Expect(err).IsNil()
	Expect(attachments).HasLen(0)
}
