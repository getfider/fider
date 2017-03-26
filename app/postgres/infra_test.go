package postgres_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

func TestHealthCheckService_IsDatabaseOnline_Success(t *testing.T) {
	RegisterTestingT(t)
	db := setup()
	defer teardown(db)

	svc := &postgres.HealthCheckService{DB: db}

	Expect(svc.IsDatabaseOnline()).To(BeTrue())
}
