import "./NotificationIndicator.scss"

import React, { useEffect, useState } from "react"
import IconBell from "@fider/assets/images/heroicons-bell.svg"
import { useFider } from "@fider/hooks"
import { actions, Fider } from "@fider/services"
import { Avatar, Icon, Markdown, Moment } from "./common"
import { Dropdown } from "./common/Dropdown"
import { Notification } from "@fider/models"
import { HStack, VStack } from "./layout"

export const NotificationItem = ({ notification }: { notification: Notification }) => {
  return (
    <HStack spacing={4} className="px-3">
      <Avatar user={{ name: notification.authorName, avatarURL: notification.avatarURL }} />
      <div>
        <Markdown className="c-notification-indicator-text" text={notification.title} style="full" />
        <span className="text-muted">
          <Moment locale={Fider.currentLocale} date={notification.createdAt} />
        </span>
      </div>
    </HStack>
  )
}

const NotificationIcon = ({ unreadNotifications }: { unreadNotifications: number }) => {
  return (
    <>
      <span className="c-notification-indicator mr-3">
        <Icon sprite={IconBell} className="h-6 text-gray-500" />
        {unreadNotifications > 0 && <div className="c-notification-indicator-unread-counter" />}
      </span>
    </>
  )
}

export const NotificationIndicator = () => {
  const fider = useFider()
  const [unreadNotifications, setUnreadNotifications] = useState(0)
  const [showingNotifications, setShowingNotifications] = useState(false)
  const [notifications, setNotifications] = useState<Notification[] | undefined>()

  useEffect(() => {
    if (fider.session.isAuthenticated) {
      actions.getTotalUnreadNotifications().then((result) => {
        if (result.ok && result.data > 0) {
          setUnreadNotifications(result.data)
        }
      })
    }
  }, [fider.session.isAuthenticated])

  useEffect(() => {
    if (showingNotifications) {
      actions.getAllNotifications().then((result) => {
        if (result) {
          setNotifications(result.data)
        }
      })
    }
  }, [showingNotifications])

  return (
    <Dropdown
      wide={true}
      position="left"
      onToggled={() => setShowingNotifications(!showingNotifications)}
      renderHandle={<NotificationIcon unreadNotifications={unreadNotifications} />}
    >
      <div>
        {showingNotifications && notifications !== undefined && notifications?.length > 0 && (
          <>
            <p className="text-display text-center my-6 bg-gray-50">No new notifications</p>
            <p className="text-title px-4 bg-gray-50">Previous notifications</p>
            <VStack spacing={2} className="c-notifications-container py-2" divide={true}>
              {notifications.map((n) => (
                <NotificationItem key={n.id} notification={n} />
              ))}
            </VStack>
          </>
        )}
      </div>
    </Dropdown>
  )
}
