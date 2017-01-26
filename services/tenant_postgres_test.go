package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTenantByDomain_WhenDomainIsUnknown_ShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery("SELECT id, name, domain FROM tenants WHERE domain = $1").WithArgs("mydomain").WillReturnRows(rows)

	svc := NewPostgresTenantService(db)
	tenant := svc.GetByDomain("mydomain")

	assert.Nil(t, tenant)
}

func TestTenantByDomain_WhenDomainIsKnown_ShouldReturnTenantInformation(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "domain"}).AddRow(234, "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, domain FROM tenants WHERE domain = \\$1").WithArgs("mydomain").WillReturnRows(rows)

	svc := NewPostgresTenantService(db)
	tenant := svc.GetByDomain("mydomain")

	assert.Equal(t, tenant.ID, 234)
	assert.Equal(t, tenant.Name, "My Domain Inc.")
	assert.Equal(t, tenant.Domain, "mydomain")
}
