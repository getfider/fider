package cmd

import "github.com/getfider/fider/app/models"

type MarkAllNotificationsAsRead struct{}

type MarkNotificationAsRead struct {
	ID int
}

type AddNewNotification struct {
	User   *models.User
	Title  string
	Link   string
	PostID int

	Result *models.Notification
}
