package services_test

import (
	"testing"

	"database/sql"

	"github.com/WeCanHearYou/wchy/services"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTenantByDomain_WhenSubdomainIsUnknown_ShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnError(sql.ErrNoRows)

	svc := &services.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	assert.Nil(t, tenant)
	assert.Error(t, err)
}

func TestTenantByDomain_WhenSubdomainIsKnown_ShouldReturnTenantInformation(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := &services.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	assert.Equal(t, 234, tenant.ID)
	assert.Equal(t, "My Domain Inc.", tenant.Name)
	assert.Equal(t, "mydomain.test.canhearyou.com", tenant.Domain)
	assert.Nil(t, err)
}

func TestTenantByDomain_WhenCompleteDomainIsUsed_ShouldReturnTenantInformation(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := &services.PostgresTenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain.anydomain.com")

	assert.Equal(t, 234, tenant.ID)
	assert.Equal(t, "My Domain Inc.", tenant.Name)
	assert.Equal(t, "mydomain.test.canhearyou.com", tenant.Domain)
	assert.Nil(t, err)
}
