package inmemory

import "github.com/getfider/fider/app/models"

// EventStorage contains read and write operations for Audit Events
type EventStorage struct {
	lastID int
	user   *models.User
	tenant *models.Tenant
	events map[int]*models.Event
}

// NewEventStorage creates a new inmemory EventStorage
func NewEventStorage() *EventStorage {
	return &EventStorage{}
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
	e.lastID++
	event := &models.Event{
		ID:       e.lastID,
		TenantID: e.tenant.ID,
		Name:     name,
	}
	e.events[event.ID] = event

	return event, nil
}

// GetByID returns the event with the specified id
func (e *EventStorage) GetByID(id int) (*models.Event, error) {
	event := e.events[id]
	return event, nil
}
