package query

import (
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/entity"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/enum"
)

type CountUnreadNotifications struct {
	Result int
}

type GetNotificationByID struct {
	ID     int
	Result *entity.Notification
}

type GetActiveNotifications struct {
	Result []*entity.Notification
}

type GetActiveSubscribers struct {
	Number  int
	Channel enum.NotificationChannel
	Event   enum.NotificationEvent

	Result []*entity.User
}
