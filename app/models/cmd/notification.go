package cmd

import "github.com/getfider/fider/app/models"

type MarkAllNotificationsAsRead struct{}

type PurgeExpiredNotifications struct {
	NumOfDeletedNotifications int
}

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

type AddSubscriber struct {
	Post *models.Post
	User *models.User
}

type RemoveSubscriber struct {
	Post *models.Post
	User *models.User
}
