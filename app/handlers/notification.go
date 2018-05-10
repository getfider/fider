package handlers

import (
	"github.com/getfider/fider/app/pkg/web"
)

// TotalUnreadNotifications returns the total number of unread notifications
func TotalUnreadNotifications() web.HandlerFunc {
	return func(c web.Context) error {
		total, err := c.Services().Notifications.TotalUnread()
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"total": total,
		})
	}
}

// Notifications is the home for unread and recent notifications
func Notifications() web.HandlerFunc {
	return func(c web.Context) error {
		notifications, err := c.Services().Notifications.GetActiveNotifications()
		if err != nil {
			return c.Failure(err)
		}

		return c.Page(web.Props{
			Title: "Notifications",
			Data: web.Map{
				"notifications": notifications,
			},
		})
	}
}

// ReadNotification marks it as read and redirect to its content
func ReadNotification() web.HandlerFunc {
	return func(c web.Context) error {
		id, err := c.ParamAsInt("id")
		if err != nil {
			return c.Failure(err)
		}

		notification, err := c.Services().Notifications.GetNotification(id)
		if err != nil {
			return c.Failure(err)
		}

		if err = c.Services().Notifications.MarkAsRead(notification.ID); err != nil {
			return c.Failure(err)
		}

		return c.Redirect(c.BaseURL() + notification.Link)
	}
}

// ReadAllNotifications marks all unread notifications as read
func ReadAllNotifications() web.HandlerFunc {
	return func(c web.Context) error {
		if err := c.Services().Notifications.MarkAllAsRead(); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
