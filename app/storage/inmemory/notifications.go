package inmemory

import (
	"github.com/getfider/fider/app/models"
)

// NotificationStorage contains read and write operations for notifications
type NotificationStorage struct {
	tenant *models.Tenant
	user   *models.User
}

// NewNotificationStorage creates a new NotificationStorage
func NewNotificationStorage() *NotificationStorage {
	return &NotificationStorage{}
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
func (s *NotificationStorage) Insert(user *models.User, title, link string, ideaID, authorID int) (*models.Notification, error) {
	return nil, nil
}

// TotalUnread returns the number of unread notifications for current user
func (s *NotificationStorage) TotalUnread() (int, error) {
	return 0, nil
}

// MarkAsRead given id of current user
func (s *NotificationStorage) MarkAsRead(id int) error {
	return nil
}

// MarkAllAsRead of current user
func (s *NotificationStorage) MarkAllAsRead() error {
	return nil
}

// GetActiveNotifications returns all unread notifications and last 30 days of read notifications
func (s *NotificationStorage) GetActiveNotifications() ([]*models.Notification, error) {
	return make([]*models.Notification, 0), nil
}

// GetNotification returns notification by id
func (s *NotificationStorage) GetNotification(id int) (*models.Notification, error) {
	return nil, nil
}
