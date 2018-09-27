package postgres

import (
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
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
func (s *NotificationStorage) Insert(user *models.User, title, link string, postID int) (*models.Notification, error) {
	if user.ID == s.user.ID {
		return nil, nil
	}

	now := time.Now()
	notification := &models.Notification{
		Title:     title,
		Link:      link,
		CreatedAt: now,
		Read:      false,
	}
	err := s.trx.Get(&notification.ID, `
		INSERT INTO notifications (tenant_id, user_id, title, link, read, post_id, author_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		RETURNING id
	`, s.tenant.ID, user.ID, title, link, false, postID, s.user.ID, now)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert notification")
	}
	return notification, nil
}

// TotalUnread returns the number of unread notifications for current user
func (s *NotificationStorage) TotalUnread() (int, error) {
	total := 0
	if s.user != nil {
		err := s.trx.Scalar(&total, "SELECT COUNT(*) FROM notifications WHERE tenant_id = $1 AND user_id = $2 AND read = false", s.tenant.ID, s.user.ID)
		if err != nil {
			return 0, errors.Wrap(err, "failed count total unread notifications")
		}
	}
	return total, nil
}

// MarkAsRead given id of current user
func (s *NotificationStorage) MarkAsRead(id int) error {
	if s.user == nil {
		return nil
	}
	_, err := s.trx.Execute(`
		UPDATE notifications SET read = true, updated_at = $1
		WHERE id = $2 AND tenant_id = $3 AND user_id = $4 AND read = false
	`, time.Now(), id, s.tenant.ID, s.user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to mark notification as read")
	}
	return nil
}

// MarkAllAsRead of current user
func (s *NotificationStorage) MarkAllAsRead() error {
	if s.user == nil {
		return nil
	}
	_, err := s.trx.Execute(`
		UPDATE notifications SET read = true, updated_at = $1
		WHERE tenant_id = $2 AND user_id = $3 AND read = false
	`, time.Now(), s.tenant.ID, s.user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to mark all notifications as read")
	}
	return nil
}

// GetActiveNotifications returns all unread notifications and last 30 days of read notifications
func (s *NotificationStorage) GetActiveNotifications() ([]*models.Notification, error) {
	notifications := []*models.Notification{}
	err := s.trx.Select(&notifications, `
		SELECT id, title, link, read, created_at 
		FROM notifications 
		WHERE tenant_id = $1 AND user_id = $2
		AND (read = false OR updated_at > CURRENT_DATE - INTERVAL '30 days')
	`, s.tenant.ID, s.user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get active notifications")
	}
	return notifications, nil
}

// GetNotification returns notification by id
func (s *NotificationStorage) GetNotification(id int) (*models.Notification, error) {
	notification := &models.Notification{}
	err := s.trx.Get(notification, `
		SELECT id, title, link, read, created_at 
		FROM notifications
		WHERE id = $1 AND tenant_id = $2 AND user_id = $3
	`, id, s.tenant.ID, s.user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get notifications with id '%d'", id)
	}
	return notification, nil
}
