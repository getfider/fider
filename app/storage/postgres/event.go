package postgres

import (
	"database/sql"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
)

type dbEvent struct {
	ID        int            `db:"id"`
	TenantID  int            `db:"tenant_id"`
	ClientIP  sql.NullString `db:"client_ip"`
	Name      string         `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
}

func (d *dbEvent) toModel() *models.Event {
	return &models.Event{
		ID:        d.ID,
		ClientIP:  d.ClientIP.String,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
	}
}

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
func (e *EventStorage) Add(clientIP, name string) (*models.Event, error) {
	event := &models.Event{
		ClientIP:  clientIP,
		Name:      name,
		CreatedAt: time.Now(),
	}
	dbClientIP := sql.NullString{
		String: clientIP,
		Valid:  len(clientIP) > 0,
	}

	err := e.trx.Get(&event.ID, `
		INSERT INTO events (tenant_id, client_ip, name, created_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, e.tenant.ID, dbClientIP, event.Name, event.CreatedAt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to insert event")
	}

	return event, nil
}

// GetByID returns the event with the specified id
func (e *EventStorage) GetByID(id int) (*models.Event, error) {
	event := &dbEvent{}
	err := e.trx.Get(event, `
		SELECT id, tenant_id, client_ip, name, created_at
		FROM events
		WHERE id = $1 AND tenant_id = $2
	`, id, e.tenant.ID)

	if err != nil {
		return nil, errors.Wrap(err, "failed to find event with id %d", id)
	}

	return event.toModel(), nil
}
