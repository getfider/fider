package postgres

import (
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

// NotificationStorage contains read and write operations for notifications
type NotificationStorage struct {
	trx    *dbx.Trx
	tenant *models.Tenant
	user   *models.User
}

// NewNotificationStorage creates a new NotificationStorage
func NewNotificationStorage(trx *dbx.Trx) *NotificationStorage {
	return &NotificationStorage{
		trx: trx,
	}
}

// SetCurrentTenant to current context
func (s *NotificationStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *NotificationStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// Insert notification for given user
func (s *NotificationStorage) Insert(user *models.User, title, link string) error {
	now := time.Now()
	return s.trx.Execute(`
		INSERT INTO notifications (tenant_id, user_id, title, link, read, created_on, updated_on) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, s.tenant.ID, user.ID, title, link, false, now, now)
}

// TotalUnread returns the number of unread notifications for current user
func (s *NotificationStorage) TotalUnread() (int, error) {
	var total int
	err := s.trx.Scalar(&total, "SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = false", s.tenant.ID, s.user.ID)
	if err != nil {
		return 0, err
	}
	return total, nil
}
