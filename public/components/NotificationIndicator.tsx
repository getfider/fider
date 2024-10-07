import "./NotificationIndicator.scss"
import NoDataIllustration from "@fider/assets/images/undraw-empty.svg"

import React, { useEffect, useState } from "react"
import IconBell from "@fider/assets/images/heroicons-bell.svg"
import { useFider } from "@fider/hooks"
import { actions, Fider } from "@fider/services"
import { Avatar, Icon, Markdown, Moment } from "./common"
import { Dropdown } from "./common/Dropdown"
import { Notification } from "@fider/models"
import { HStack, VStack } from "./layout"

import { Trans } from "@lingui/macro"

export const NotificationItem = ({ notification }: { notification: Notification }) => {
  const openNotification = () => {
    console.log(notification.link)
    window.location.href = `/notifications/${notification.id}`
  }

  return (
    <HStack spacing={4} className="px-3 pr-5 clickable hover py-2" onClick={openNotification}>
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
                <p className="text-subtitle px-4 mt-2 mb-0">
                  <Trans id="modal.notifications.unread">Unread notifications</Trans>
                </p>
                <VStack spacing={2} className="py-3 mb-2 no-lastchild-paddingzero" divide={true}>
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
                  {recent?.length === 0 && (
                    <Icon sprite={NoDataIllustration} height="120" className="mt-6 mb-2" />
                  )}
                </div>
            )}
            {recent !== undefined && recent?.length > 0 && (
              <>
                <p className="text-subtitle px-4 mb-0">
                  <Trans id="modal.notifications.previous">Previous notifications</Trans>
                </p>
                <VStack spacing={2} className="py-2 no-lastchild-paddingzero" divide={true}>
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
