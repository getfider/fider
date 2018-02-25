package postgres_test

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNotificationStorage_Insert_Count(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	notifications.SetCurrentTenant(demoTenant)
	Expect(notifications.Insert(jonSnow, "Hello World", "http://www.google.com.br")).To(BeNil())
	Expect(notifications.Insert(aryaStark, "Hello World", "http://www.google.com.br")).To(BeNil())
	Expect(notifications.Insert(aryaStark, "Another thing happened", "http://www.google.com.br")).To(BeNil())

	notifications.SetCurrentUser(jonSnow)
	Expect(notifications.TotalUnread()).To(Equal(1))

	notifications.SetCurrentUser(aryaStark)
	Expect(notifications.TotalUnread()).To(Equal(2))

	notifications.SetCurrentUser(sansaStark)
	Expect(notifications.TotalUnread()).To(Equal(0))
}
