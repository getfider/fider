package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"

	"github.com/getfider/fider/app"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestNotificationStorage_TotalCount(t *testing.T) {
	ctx := SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	q := &query.CountUnreadNotifications{}

	err := bus.Dispatch(ctx, q)
	Expect(err).IsNil()
	Expect(q.Result).Equals(0)

	err = bus.Dispatch(demoTenantCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).Equals(0)
}

func TestNotificationStorage_Insert_Read_Count(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification1 := &cmd.AddNewNotification{User: aryaStark, Title: "Hello World", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	addNotification2 := &cmd.AddNewNotification{User: aryaStark, Title: "Hello World", Link: "", PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, addNotification1, addNotification2)
	Expect(err).IsNil()

	q := &query.CountUnreadNotifications{}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).Equals(2)

	Expect(bus.Dispatch(aryaStarkCtx, &cmd.MarkNotificationAsRead{ID: addNotification1.Result.ID})).IsNil()
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).Equals(1)

	Expect(bus.Dispatch(aryaStarkCtx, &cmd.MarkNotificationAsRead{ID: addNotification2.Result.ID})).IsNil()
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).Equals(0)
}

func TestNotificationStorage_GetActiveNotifications(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification1 := &cmd.AddNewNotification{User: aryaStark, Title: "Hello World", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	addNotification2 := &cmd.AddNewNotification{User: aryaStark, Title: "Another thing happened", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, addNotification1, addNotification2)
	Expect(err).IsNil()

	activeNotifications := &query.GetActiveNotifications{}
	err = bus.Dispatch(aryaStarkCtx, activeNotifications)
	Expect(err).IsNil()
	Expect(activeNotifications.Result).HasLen(2)

	bus.MustDispatch(aryaStarkCtx, &cmd.MarkNotificationAsRead{ID: activeNotifications.Result[0].ID})
	bus.MustDispatch(aryaStarkCtx, &cmd.MarkNotificationAsRead{ID: activeNotifications.Result[1].ID})

	_, err = trx.Execute("UPDATE notifications SET updated_at = $1 WHERE id = $2", time.Now().AddDate(0, 0, -31), activeNotifications.Result[0].ID)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, activeNotifications)
	Expect(err).IsNil()
	Expect(activeNotifications.Result).HasLen(1)
	Expect(activeNotifications.Result[0].Read).IsTrue()
}

func TestNotificationStorage_ReadAll(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification1 := &cmd.AddNewNotification{User: aryaStark, Title: "Hello World", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	addNotification2 := &cmd.AddNewNotification{User: aryaStark, Title: "Another thing happened", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, addNotification1, addNotification2)
	Expect(err).IsNil()

	activeNotifications := &query.GetActiveNotifications{}
	err = bus.Dispatch(aryaStarkCtx, activeNotifications)
	Expect(err).IsNil()
	Expect(activeNotifications.Result).HasLen(2)

	err = bus.Dispatch(aryaStarkCtx, &cmd.MarkAllNotificationsAsRead{})
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, activeNotifications)
	Expect(err).IsNil()
	Expect(activeNotifications.Result).HasLen(2)
	Expect(activeNotifications.Result[0].Read).IsTrue()
	Expect(activeNotifications.Result[1].Read).IsTrue()
}

func TestNotificationStorage_GetNotificationByID(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification := &cmd.AddNewNotification{User: aryaStark, Title: "Hello World", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	err = bus.Dispatch(jonSnowCtx, addNotification)
	Expect(err).IsNil()

	q := &query.GetNotificationByID{ID: addNotification.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result.Title).Equals("Hello World")
	Expect(q.Result.Link).Equals("http://www.google.com.br")
	Expect(q.Result.Read).IsFalse()
}

func TestNotificationStorage_GetNotificationByID_OtherUser(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification := &cmd.AddNewNotification{User: jonSnow, Title: "Hello World", Link: "http://www.google.com.br", PostID: newPost.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, addNotification)
	Expect(err).IsNil()

	q := &query.GetNotificationByID{ID: addNotification.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(q.Result).IsNil()
}

func TestNotificationStorage_PurgeExpiredNotifications(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()
	defer ResetDatabase()

	newPost := &cmd.AddNewPost{Title: "Title", Description: "Description"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	addNotification1 := &cmd.AddNewNotification{User: jonSnow, Title: "Hello World 1", Link: "http://www.google.com", PostID: newPost.Result.ID}
	addNotification2 := &cmd.AddNewNotification{User: jonSnow, Title: "Hello World 2", Link: "http://www.google.com", PostID: newPost.Result.ID}
	addNotification3 := &cmd.AddNewNotification{User: jonSnow, Title: "Hello World 3", Link: "http://www.microsoft.com", PostID: newPost.Result.ID}
	addNotification4 := &cmd.AddNewNotification{User: jonSnow, Title: "Hello World 4", Link: "http://www.microsoft.com", PostID: newPost.Result.ID}
	err = bus.Dispatch(aryaStarkCtx, addNotification1, addNotification2, addNotification3, addNotification4)
	Expect(err).IsNil()

	rows, err := trx.Execute("UPDATE notifications SET created_at = NOW() - INTERVAL '2 years' WHERE link = 'http://www.microsoft.com'")
	Expect(err).IsNil()
	Expect(rows).Equals(int64(2))

	trx.MustCommit()

	purgeCommand := &cmd.PurgeExpiredNotifications{}
	err = bus.Dispatch(context.Background(), purgeCommand)
	Expect(err).IsNil()
	Expect(purgeCommand.NumOfDeletedNotifications).Equals(2)
}
