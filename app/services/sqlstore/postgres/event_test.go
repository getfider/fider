package postgres_test

import (
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEventStorage_Add(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	clientIP := "127.0.0.1"
	eventName := "posts.create"
	err := bus.Dispatch(demoTenantCtx, &cmd.StoreEvent{
		ClientIP:  clientIP,
		EventName: eventName,
	})
	Expect(err).IsNil()

	count, err := trx.Count("SELECT * FROM events WHERE name = 'posts.create' AND client_ip = '127.0.0.1' AND tenant_id = 1")
	Expect(err).IsNil()
	Expect(count).Equals(1)
}

func TestEventStorage_AddWithNullClientIP(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	eventName := "posts.delete"
	err := bus.Dispatch(avengersTenantCtx, &cmd.StoreEvent{
		EventName: eventName,
	})
	Expect(err).IsNil()

	count, err := trx.Count("SELECT * FROM events WHERE name = 'posts.delete' AND client_ip IS NULL AND tenant_id = 2")
	Expect(err).IsNil()
	Expect(count).Equals(1)
}
