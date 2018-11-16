package postgres_test

import (
	"strconv"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSubscription_NoSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)
	post1, _ := posts.Add("Post #1", "Description #1")

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(nil)
	subscribed, err := users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(jonSnow)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_RemoveSubscriber(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)
	post1, _ := posts.Add("Post #1", "Description #1")
	err := posts.RemoveSubscriber(post1, aryaStark)
	Expect(err).IsNil()

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminSubmitted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post1, _ := posts.Add("Post #1", "Description #1")

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminUnsubscribed(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)

	post1, _ := posts.Add("Post #1", "Description #1")
	posts.RemoveSubscriber(post1, aryaStark)
	posts.RemoveSubscriber(post1, jonSnow)

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_DisabledEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)
	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err := users.UpdateSettings(map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	post1, _ := posts.Add("Post #1", "Description #1")

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(post1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_VisitorEnabledNewPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err := users.UpdateSettings(map[string]string{
		models.NotificationEventNewPost.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelEmail | models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	post1, _ := posts.Add("Post #1", "Description #1")

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelEmail, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_DisabledEverything(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	users.SetCurrentTenant(demoTenant)
	disableAll := map[string]string{
		models.NotificationEventNewPost.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   "0",
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}
	users.SetCurrentUser(jonSnow)
	Expect(users.UpdateSettings(disableAll)).IsNil()
	users.SetCurrentUser(aryaStark)
	Expect(users.UpdateSettings(disableAll)).IsNil()

	post1, _ := posts.Add("Post #1", "Description #1")

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelEmail, models.NotificationEventNewPost)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = posts.GetActiveSubscribers(post1.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)
}

func TestSubscription_DeletedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(aryaStark)

	post1, _ := posts.Add("Post #1", "Description #1")
	posts.SetResponse(post1, "Invalid Post!", models.PostDeleted)

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_SubscribedToDifferentPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)

	post1, _ := posts.Add("Post #1", "Description #1")
	post2, _ := posts.Add("Post #2", "Description #2")
	err := posts.AddSubscriber(post2, aryaStark)
	Expect(err).IsNil()

	subscribers, err := posts.GetActiveSubscribers(post1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
}
