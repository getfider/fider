package postgres_test

import (
	"strconv"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSubscription_NoSettings(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(1)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(2)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(changeStatusSubscribers.Result[1].ID).Equals(aryaStark.ID)

	subscribed := &query.UserSubscribedTo{PostID: newPost.Result.ID}

	err = bus.Dispatch(ctx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()

	err = bus.Dispatch(jonSnowCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(aryaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(sansaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()
}

func TestSubscription_RemoveSubscriber(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: aryaStark})
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(1)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(1)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)

	subscribed := &query.UserSubscribedTo{PostID: newPost.Result.ID}

	err = bus.Dispatch(jonSnowCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(aryaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()

	err = bus.Dispatch(sansaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()
}

func TestSubscription_AdminSubmitted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(jonSnowCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(1)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(1)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)

	subscribed := &query.UserSubscribedTo{PostID: newPost.Result.ID}

	err = bus.Dispatch(jonSnowCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(aryaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()

	err = bus.Dispatch(sansaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()
}

func TestSubscription_AdminUnsubscribed(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	bus.MustDispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: aryaStark})
	bus.MustDispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: jonSnow})

	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	err = bus.Dispatch(aryaStarkCtx, newCommentSubscribers)
	Expect(err).IsNil()
	Expect(newCommentSubscribers.Result).HasLen(0)

	subscribed := &query.UserSubscribedTo{PostID: newPost.Result.ID}

	err = bus.Dispatch(jonSnowCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()

	err = bus.Dispatch(aryaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()

	err = bus.Dispatch(sansaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()
}

func TestSubscription_DisabledEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.UpdateCurrentUserSettings{
		Settings: map[string]string{
			enum.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(enum.NotificationChannelWeb)),
		},
	})
	Expect(err).IsNil()

	newCommentWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	newCommentEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newCommentWebSubscribers, newCommentEmailSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newCommentWebSubscribers.Result).HasLen(2)
	Expect(newCommentWebSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newCommentWebSubscribers.Result[1].ID).Equals(aryaStark.ID)

	Expect(newCommentEmailSubscribers.Result).HasLen(1)
	Expect(newCommentEmailSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(2)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(changeStatusSubscribers.Result[1].ID).Equals(aryaStark.ID)

	subscribed := &query.UserSubscribedTo{PostID: newPost.Result.ID}

	err = bus.Dispatch(jonSnowCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(aryaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsTrue()

	err = bus.Dispatch(sansaStarkCtx, subscribed)
	Expect(err).IsNil()
	Expect(subscribed.Result).IsFalse()
}

func TestSubscription_VisitorEnabledNewPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.UpdateCurrentUserSettings{
		Settings: map[string]string{
			enum.NotificationEventNewPost.UserSettingsKeyName: strconv.Itoa(int(enum.NotificationChannelEmail | enum.NotificationChannelWeb)),
		},
	})
	Expect(err).IsNil()

	newPostWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewPost}
	newPostEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewPost}
	err = bus.Dispatch(aryaStarkCtx, newPostWebSubscribers, newPostEmailSubscribers)
	Expect(err).IsNil()

	Expect(newPostWebSubscribers.Result).HasLen(2)
	Expect(newPostWebSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newPostWebSubscribers.Result[1].ID).Equals(aryaStark.ID)

	Expect(newPostEmailSubscribers.Result).HasLen(2)
	Expect(newPostEmailSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newPostEmailSubscribers.Result[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_DisabledEverything(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	disableAll := map[string]string{
		enum.NotificationEventNewPost.UserSettingsKeyName:      "0",
		enum.NotificationEventNewComment.UserSettingsKeyName:   "0",
		enum.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}

	err = bus.Dispatch(aryaStarkCtx, &cmd.UpdateCurrentUserSettings{Settings: disableAll})
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.UpdateCurrentUserSettings{Settings: disableAll})
	Expect(err).IsNil()

	newPostWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewPost}
	newPostEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewPost}
	newCommentWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	newCommentEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewComment}
	changeStatusWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	changeStatusEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostWebSubscribers, newPostEmailSubscribers, newCommentWebSubscribers, newCommentEmailSubscribers, changeStatusWebSubscribers, changeStatusEmailSubscribers)
	Expect(err).IsNil()

	Expect(newPostWebSubscribers.Result).HasLen(0)
	Expect(newPostEmailSubscribers.Result).HasLen(0)
	Expect(newCommentWebSubscribers.Result).HasLen(0)
	Expect(newCommentEmailSubscribers.Result).HasLen(0)
	Expect(changeStatusWebSubscribers.Result).HasLen(0)
	Expect(changeStatusEmailSubscribers.Result).HasLen(0)
}

func TestSubscription_DeletedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "Invalid Post!", Status: enum.PostDeleted})
	Expect(err).IsNil()

	q := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(2)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
	Expect(q.Result[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_SubscribedToDifferentPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost1 := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	newPost2 := &cmd.AddNewPost{Title: "Post #2", Description: "Description #2"}
	err := bus.Dispatch(jonSnowCtx, newPost1, newPost2)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddSubscriber{Post: newPost2.Result, User: aryaStark})
	Expect(err).IsNil()

	q := &query.GetActiveSubscribers{Number: newPost1.Result.Number, Channel: enum.NotificationChannelWeb, Event: enum.NotificationEventNewComment}
	err = bus.Dispatch(jonSnowCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(1)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
}

func TestSubscription_EmailSupressed(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// Enable email notifications for new comments
	err := bus.Dispatch(aryaStarkCtx, &cmd.UpdateCurrentUserSettings{
		Settings: map[string]string{
			enum.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(enum.NotificationChannelEmail)),
		},
	})
	Expect(err).IsNil()

	newPost1 := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err = bus.Dispatch(aryaStarkCtx, newPost1)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.AddSubscriber{Post: newPost1.Result, User: aryaStark})
	Expect(err).IsNil()

	q := &query.GetActiveSubscribers{Number: newPost1.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewComment}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(2)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
	Expect(q.Result[1].ID).Equals(aryaStark.ID)

	//Supress the email and verify that AryaStark is not an active subscriber anymore
	err = bus.Dispatch(aryaStarkCtx, &cmd.SupressEmail{EmailAddresses: []string{aryaStark.Email}})
	Expect(err).IsNil()

	q = &query.GetActiveSubscribers{Number: newPost1.Result.Number, Channel: enum.NotificationChannelEmail, Event: enum.NotificationEventNewComment}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(1)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
}
