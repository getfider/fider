package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"

	. "github.com/onsi/gomega"
)

func TestNotificationStorage_TotalCount(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)

	notifications.SetCurrentUser(jonSnow)
	Expect(notifications.TotalUnread()).To(Equal(0))

	notifications.SetCurrentUser(nil)
	Expect(notifications.TotalUnread()).To(Equal(0))
}

func TestNotificationStorage_Insert_Read_Count(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)

	not1, err := notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br")
	Expect(err).To(BeNil())
	not2, err := notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br")
	Expect(err).To(BeNil())

	notifications.SetCurrentUser(aryaStark)
	Expect(notifications.TotalUnread()).To(Equal(2))

	Expect(notifications.MarkAsRead(not1.ID)).To(BeNil())
	Expect(notifications.TotalUnread()).To(Equal(1))
	Expect(notifications.MarkAsRead(not2.ID)).To(BeNil())
	Expect(notifications.TotalUnread()).To(Equal(0))
}

func TestNotificationStorage_GetActiveNotifications(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(aryaStark)

	notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br")
	notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br")

	allNotifications, err := notifications.GetActiveNotifications()
	Expect(err).To(BeNil())
	Expect(allNotifications).To(HaveLen(2))

	notifications.MarkAsRead(allNotifications[0].ID)
	notifications.MarkAsRead(allNotifications[1].ID)
	trx.Execute("UPDATE notifications SET updated_on = $1 WHERE id = $2", time.Now().AddDate(0, 0, -31), allNotifications[0].ID)

	allNotifications, err = notifications.GetActiveNotifications()
	Expect(err).To(BeNil())
	Expect(allNotifications).To(HaveLen(1))
}

func TestNotificationStorage_GetNotificationById(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(aryaStark)
	not1, err := notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br")
	Expect(err).To(BeNil())
	not1, err = notifications.GetNotification(not1.ID)
	Expect(err).To(BeNil())
	Expect(not1.Title).To(Equal("Hello World"))
	Expect(not1.Link).To(Equal("http://www.google.com.br"))
	Expect(not1.Read).To(BeFalse())
}

func TestNotificationStorage_GetNotificationById_OtherUser(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)
	not1, err := notifications.Insert(jonSnow, "Hello World", "http://www.google.com.br")
	Expect(err).To(BeNil())

	notifications.SetCurrentUser(aryaStark)
	not1, err = notifications.GetNotification(not1.ID)
	Expect(err).To(Equal(app.ErrNotFound))
	Expect(not1).To(BeNil())
}
