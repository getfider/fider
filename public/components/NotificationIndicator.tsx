import "./NotificationIndicator.scss"
import NoDataIllustration from "@fider/assets/images/undraw-empty.svg"

import React, { useEffect, useState } from "react"
import IconBell from "@fider/assets/images/heroicons-bell.svg"
import { useFider } from "@fider/hooks"
import { actions, Fider } from "@fider/services"
import { Avatar, Icon, Markdown, Moment } from "./common"
import { Dropdown } from "./common/Dropdown"
import { Notification } from "@fider/models"
import { VStack } from "./layout"

import { Trans } from "@lingui/react/macro"

export const NotificationItem = ({ notification }: { notification: Notification }) => {
  return (
    <a href={`/notifications/${notification.id}`} className="px-3 pr-5 hover py-4 flex flex-x flex--spacing-4 flex-items-center">
      <Avatar user={{ name: notification.authorName, avatarURL: notification.avatarURL }} />
      <div>
        <Markdown className="c-notification-indicator-text" text={notification.title} style="full" />
        <span className="text-muted">
          <Moment locale={Fider.currentLocale} date={notification.createdAt} />
        </span>
      </div>
    </a>
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
  const [recent, setRecent] = useState<Notification[] | undefined>()
  const [unread, setUnread] = useState<Notification[] | undefined>()

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
          const [unread, recent] = (result.data || []).reduce(
            (result, item) => {
              result[item.read ? 1 : 0].push(item)
              return result
            },
            [[] as Notification[], [] as Notification[]]
          )
          setRecent(recent)
          setUnread(unread)
          setUnreadNotifications(unread.length)
        }
      })
    }
  }, [showingNotifications])

  const markAllAsRead = async (e: React.MouseEvent) => {
    e.preventDefault()
    const response = await actions.markAllAsRead()
    if (response.ok) {
      location.reload()
    }
  }

  return (
    <Dropdown
      wide={true}
      position="left"
      fullsceenSm={true}
      onToggled={(isOpen: boolean) => setShowingNotifications(isOpen)}
      renderHandle={<NotificationIcon unreadNotifications={unreadNotifications} />}
    >
      <div className="c-notifications-container">
        {showingNotifications && (unread !== undefined || recent !== undefined) && (
          <>
            {unread !== undefined && unread?.length > 0 ? (
              <>
                <p className="text-subtitle px-4 mt-4 mb-0">
                  <Trans id="modal.notifications.unread">Unread notifications</Trans>
                  {unread.length > 1 && (
                    <a href="#" className="text-link text-xs pl-6" onClick={markAllAsRead}>
                      <Trans id="action.markallasread">Mark All as Read</Trans>
                    </a>
                  )}
                </p>
                <VStack spacing={0} className="py-2" divide={false}>
                  {unread.map((n) => (
                    <NotificationItem key={n.id} notification={n} />
                  ))}
                </VStack>
              </>
            ) : (
              <div className="text-center pb-6">
                <p className="text-display text-center mt-6 px-4">
                  <Trans id="modal.notifications.nonew">No new notifications</Trans>
                </p>
                {recent?.length === 0 && <Icon sprite={NoDataIllustration} height="120" className="mt-6 mb-2" />}
              </div>
            )}
            {recent !== undefined && recent?.length > 0 && (
              <>
                <p className="text-subtitle px-4 mb-0 pt-4 bg-gray-50 border-gray-200 border-t">
                  <Trans id="modal.notifications.previous">Previous notifications</Trans>
                </p>
                <VStack spacing={0} className="py-2 bg-gray-50" divide={false}>
                  {recent.map((n) => (
                    <NotificationItem key={n.id} notification={n} />
                  ))}
                </VStack>
              </>
            )}
          </>
        )}
      </div>
    </Dropdown>
  )
}
