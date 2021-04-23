import "./NotificationIndicator.scss"

import React, { useEffect, useState } from "react"
import IconBell from "@fider/assets/images/heroicons-bell.svg"
import { useFider } from "@fider/hooks"
import { actions } from "@fider/services"
import { Icon } from "./common"

export const NotificationIndicator = () => {
  const fider = useFider()
  const [unreadNotifications, setUnreadNotifications] = useState(0)

  useEffect(() => {
    if (fider.session.isAuthenticated) {
      actions.getTotalUnreadNotifications().then((result) => {
        if (result.ok && result.data > 0) {
          setUnreadNotifications(result.data)
        }
      })
    }
  }, [fider.session.isAuthenticated])

  return (
    <a href="/notifications" className="c-notification-indicator">
      <Icon sprite={IconBell} className="h-6 text-gray-500" />
      {unreadNotifications > 0 && <div className="c-notification-indicator-unread-counter" />}
    </a>
  )
}
