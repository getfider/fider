package postgres_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestNotificationStorage_TotalCount(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)

	notifications.SetCurrentUser(jonSnow)
	total, err := notifications.TotalUnread()
	Expect(err).IsNil()
	Expect(total).Equals(0)

	notifications.SetCurrentUser(nil)
	total, err = notifications.TotalUnread()
	Expect(err).IsNil()
	Expect(total).Equals(0)
}

func TestNotificationStorage_Insert_Read_Count(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(jonSnow)
	post, _ := posts.Add("Title", "Description")

	not1, err := notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br", post.ID)
	Expect(err).IsNil()
	not2, err := notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br", post.ID)
	Expect(err).IsNil()

	notifications.SetCurrentUser(aryaStark)
	total, err := notifications.TotalUnread()
	Expect(err).IsNil()
	Expect(total).Equals(2)

	Expect(notifications.MarkAsRead(not1.ID)).IsNil()
	total, err = notifications.TotalUnread()
	Expect(err).IsNil()
	Expect(total).Equals(1)

	Expect(notifications.MarkAsRead(not2.ID)).IsNil()
	total, err = notifications.TotalUnread()
	Expect(err).IsNil()
	Expect(total).Equals(0)
}

func TestNotificationStorage_GetActiveNotifications(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(jonSnow)
	post, _ := posts.Add("Title", "Description")

	notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br", post.ID)
	notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br", post.ID)

	notifications.SetCurrentUser(aryaStark)

	allNotifications, err := notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(allNotifications).HasLen(2)

	notifications.MarkAsRead(allNotifications[0].ID)
	notifications.MarkAsRead(allNotifications[1].ID)
	trx.Execute("UPDATE notifications SET updated_at = $1 WHERE id = $2", time.Now().AddDate(0, 0, -31), allNotifications[0].ID)

	allNotifications, err = notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(allNotifications).HasLen(1)
	Expect(allNotifications[0].Read).IsTrue()
}

func TestNotificationStorage_ReadAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(jonSnow)

	post, _ := posts.Add("Title", "Description")

	notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br", post.ID)
	notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br", post.ID)

	notifications.SetCurrentUser(aryaStark)

	allNotifications, err := notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(allNotifications).HasLen(2)

	notifications.MarkAllAsRead()

	allNotifications, err = notifications.GetActiveNotifications()
	Expect(err).IsNil()
	Expect(allNotifications).HasLen(2)
	Expect(allNotifications[0].Read).IsTrue()
	Expect(allNotifications[1].Read).IsTrue()
}

func TestNotificationStorage_GetNotificationByID(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(jonSnow)
	post, _ := posts.Add("Title", "Description")

	not1, err := notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br", post.ID)
	Expect(err).IsNil()

	notifications.SetCurrentUser(aryaStark)
	not1, err = notifications.GetNotification(not1.ID)
	Expect(err).IsNil()
	Expect(not1.Title).Equals("Hello World")
	Expect(not1.Link).Equals("http://www.google.com.br")
	Expect(not1.Read).IsFalse()
}

func TestNotificationStorage_GetNotificationByID_OtherUser(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	posts.SetCurrentTenant(demoTenant)
	posts.SetCurrentUser(jonSnow)
	post, _ := posts.Add("Title", "Description")

	notifications.SetCurrentTenant(demoTenant)
	notifications.SetCurrentUser(aryaStark)
	not1, err := notifications.Insert(jonSnow, "Hello World", "http://www.google.com.br", post.ID)
	Expect(err).IsNil()

	not1, err = notifications.GetNotification(not1.ID)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(not1).IsNil()
}
