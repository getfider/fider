package postgres_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestEventStorage_Add(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	// TODO: Add Real Test
	Expect(true).IsTrue()
}
