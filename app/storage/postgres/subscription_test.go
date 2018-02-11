package postgres_test

import (
	"strconv"
	"testing"

	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestSubscription_NoSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea1, _ := ideas.Add("Idea #1", "Description #1", aryaStark.ID)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))
}

func TestSubscription_RemoveSubscriber(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	idea1, _ := ideas.Add("Idea #1", "Description #1", aryaStark.ID)
	err := ideas.RemoveSubscriber(idea1.Number, aryaStark.ID)
	Expect(err).To(BeNil())

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
}

func TestSubscription_AdminSubmitted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)

	idea1, _ := ideas.Add("Idea #1", "Description #1", jonSnow.ID)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
}

func TestSubscription_DisabledEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	err := users.UpdateSettings(aryaStark.ID, map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelWeb)),
	})
	Expect(err).To(BeNil())

	idea1, _ := ideas.Add("Idea #1", "Description #1", aryaStark.ID)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(1))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))
}

func TestSubscription_VisitorEnabledNewIdea(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	err := users.UpdateSettings(aryaStark.ID, map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelEmail | models.NotificationChannelWeb)),
	})
	Expect(err).To(BeNil())

	idea1, _ := ideas.Add("Idea #1", "Description #1", jonSnow.ID)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(2))
	Expect(subscribers[0].ID).To(Equal(jonSnow.ID))
	Expect(subscribers[1].ID).To(Equal(aryaStark.ID))
}

func TestSubscription_DisabledEverything(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	ideas.SetCurrentTenant(demoTenant)
	disableAll := map[string]string{
		models.NotificationEventNewIdea.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   "0",
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}
	err := users.UpdateSettings(jonSnow.ID, disableAll)
	Expect(err).To(BeNil())
	err = users.UpdateSettings(aryaStark.ID, disableAll)
	Expect(err).To(BeNil())

	idea1, _ := ideas.Add("Idea #1", "Description #1", jonSnow.ID)

	subscribers, err := ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelWeb, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewIdea)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventNewComment)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))

	subscribers, err = ideas.GetActiveSubscribers(idea1.Number, models.NotificationChannelEmail, models.NotificationEventChangeStatus)
	Expect(err).To(BeNil())
	Expect(len(subscribers)).To(Equal(0))
}
