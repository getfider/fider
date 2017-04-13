package postgres_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/postgres"
	. "github.com/onsi/gomega"
)

func TestHealthCheckService_IsDatabaseOnline_Success(t *testing.T) {
	RegisterTestingT(t)
	db, _ := dbx.New()
	defer db.Close()

	svc := &postgres.HealthCheckService{DB: db}

	Expect(svc.IsDatabaseOnline()).To(BeTrue())
}
