package services_test

import (
	"testing"

	"database/sql"

	"github.com/WeCanHearYou/wchy-api/services"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTenantByDomain_WhenDomainIsUnknown_ShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnError(sql.ErrNoRows)

	svc := &services.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	assert.Nil(t, tenant)
	assert.Error(t, err)
}

func TestTenantByDomain_WhenDomainIsKnown_ShouldReturnTenantInformation(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := &services.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	assert.Equal(t, tenant.ID, 234)
	assert.Equal(t, tenant.Name, "My Domain Inc.")
	assert.Equal(t, tenant.Domain, "mydomain")
	assert.Nil(t, err)
}
