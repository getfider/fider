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
