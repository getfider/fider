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

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(aryaStark)
	idea1, _ := ideas.Add("Idea #1", "Description #1")

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(nil)
	subscribed, err := users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(jonSnow)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

}

func TestSubscription_RemoveSubscriber(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(aryaStark)
	idea1, _ := ideas.Add("Idea #1", "Description #1")
	err := ideas.RemoveSubscriber(idea1, aryaStark)
	Expect(err).IsNil()

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminSubmitted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)

	idea1, _ := ideas.Add("Idea #1", "Description #1")

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminUnsubscribed(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(aryaStark)

	idea1, _ := ideas.Add("Idea #1", "Description #1")
	ideas.RemoveSubscriber(idea1, aryaStark)
	ideas.RemoveSubscriber(idea1, jonSnow)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_DisabledEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(aryaStark)
	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err := users.UpdateSettings(map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	idea1, _ := ideas.Add("Idea #1", "Description #1")

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(1)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(idea1.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_VisitorEnabledNewIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err := users.UpdateSettings(map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelEmail | models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	idea1, _ := ideas.Add("Idea #1", "Description #1")

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(2)
	Expect(subscribers[0].ID).Equals(jonSnow.ID)
	Expect(subscribers[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_DisabledEverything(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	ideas.SetCurrentUser(jonSnow)
	users.SetCurrentTenant(demoTenant)
	disableAll := map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   "0",
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}
	users.SetCurrentUser(jonSnow)
	Expect(users.UpdateSettings(disableAll)).IsNil()
	users.SetCurrentUser(aryaStark)
	Expect(users.UpdateSettings(disableAll)).IsNil()

	idea1, _ := ideas.Add("Idea #1", "Description #1")

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
	Expect(err).IsNil()
	Expect(subscribers).HasLen(0)
}
