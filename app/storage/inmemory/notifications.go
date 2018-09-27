package inmemory

import (
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// NotificationStorage contains read and write operations for notifications
type NotificationStorage struct {
	lastID        int
	tenant        *models.Tenant
	user          *models.User
	notifications map[*models.User][]*models.Notification
}

// NewNotificationStorage creates a new NotificationStorage
func NewNotificationStorage() *NotificationStorage {
	return &NotificationStorage{
		lastID:        0,
		notifications: make(map[*models.User][]*models.Notification, 0),
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
	s.lastID = s.lastID + 1
	newNotification := &models.Notification{
		ID:        s.lastID,
		Title:     title,
		Link:      link,
		Read:      false,
		CreatedAt: time.Now(),
	}
	all, ok := s.notifications[user]
	if !ok {
		all = make([]*models.Notification, 0)
	}
	s.notifications[user] = append(all, newNotification)
	return newNotification, nil
}

// TotalUnread returns the number of unread notifications for current user
func (s *NotificationStorage) TotalUnread() (int, error) {
	all, ok := s.notifications[s.user]
	if !ok {
		return 0, nil
	}
	return len(all), nil
}

// MarkAsRead given id of current user
func (s *NotificationStorage) MarkAsRead(id int) error {
	all, ok := s.notifications[s.user]
	if !ok {
		return nil
	}

	for _, notification := range all {
		if notification.ID == id {
			notification.Read = true
		}
	}
	return nil
}

// MarkAllAsRead of current user
func (s *NotificationStorage) MarkAllAsRead() error {
	all, ok := s.notifications[s.user]
	if !ok {
		return nil
	}

	for _, notification := range all {
		notification.Read = true
	}

	return nil
}

// GetActiveNotifications returns all unread notifications and last 30 days of read notifications
func (s *NotificationStorage) GetActiveNotifications() ([]*models.Notification, error) {
	all, ok := s.notifications[s.user]
	if !ok {
		return make([]*models.Notification, 0), nil
	}
	return all, nil
}

// GetNotification returns notification by id
func (s *NotificationStorage) GetNotification(id int) (*models.Notification, error) {
	all, ok := s.notifications[s.user]
	if !ok {
		return nil, app.ErrNotFound
	}

	for _, notification := range all {
		if notification.ID == id {
			return notification, nil
		}
	}
	return nil, app.ErrNotFound
}
