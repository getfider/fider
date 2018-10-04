package postgres

import (
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
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
		TenantID:  e.tenant.ID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	err := e.trx.Get(&event.ID, `
		INSERT INTO events (tenant_id, name, created_at) 
		VALUES ($1, $2, $3)
		RETURNING id
	`, event.TenantID, event.Name, event.CreatedAt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert event")
	}

	return event, nil
}

// GetByID returns the event with the specified id
func (e *EventStorage) GetByID(id int) (*models.Event, error) {
	event := &models.Event{}
	err := e.trx.Get(event, `
		SELECT id, tenant_id, name, created_at
		FROM events
		WHERE id = $1 AND tenant_id = $2
	`, id, e.tenant.ID)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find event with id %d", id)
	}

	return event, nil
}
