package postgres

import (
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

// EventStorage contains read and write operations for Audit Events
type EventStorage struct {
	trx    *dbx.Trx
	user   *models.User
	tenant *models.Tenant
}

// NewEventStorage creates a new inmemory EventStorage
func NewEventStorage(trx *dbx.Trx) *EventStorage {
	return &EventStorage{
		trx: trx,
	}
}

// SetCurrentTenant to current context
func (e *EventStorage) SetCurrentTenant(tenant *models.Tenant) {
	e.tenant = tenant
}

// SetCurrentUser to current context
func (e *EventStorage) SetCurrentUser(user *models.User) {
	e.user = user
}

// Add stores a new event
func (e *EventStorage) Add(name string) (*models.Event, error) {
	// TODO: Add Query Logic, Review struct tags
	event := &models.Event{
		ID:       0,
		TenantID: e.tenant.ID,
		Name:     name,
	}

	return event, nil
}
