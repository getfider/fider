package query

import "github.com/getfider/fider/app/models"

type CountUnreadNotifications struct {
	Result int
}

type GetNotificationByID struct {
	ID     int
	Result *models.Notification
}

type GetActiveNotifications struct {
	Result []*models.Notification
}
