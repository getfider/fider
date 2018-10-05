package postgres_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEventStorage_AddAndGet(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	events.SetCurrentTenant(demoTenant)

	clientIP := "127.0.0.1"
	name := "posts.create"
	event, err := events.Add(clientIP, name)
	Expect(err).IsNil()

	event, err = events.GetByID(event.ID)
	Expect(err).IsNil()
	Expect(event.TenantID).Equals(demoTenant.ID)
	Expect(event.ClientIP).Equals(clientIP)
	Expect(event.Name).Equals(name)
}
