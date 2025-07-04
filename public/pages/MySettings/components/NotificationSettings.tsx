import React, { useState } from "react"

import { UserSettings } from "@fider/models"
import { Toggle, Field } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"

interface NotificationSettingsProps {
  userSettings: UserSettings
  settingsChanged: (settings: UserSettings) => void
}

type Channel = number
const WebChannel: Channel = 1
const EmailChannel: Channel = 2

export const NotificationSettings = (props: NotificationSettingsProps) => {
  const [userSettings, setUserSettings] = useState(props.userSettings)

  const isEnabled = (settingsKey: string, channel: Channel): boolean => {
    if (settingsKey in userSettings) {
      return (parseInt(userSettings[settingsKey], 10) & channel) > 0
    }
    return false
  }

  const toggle = async (settingsKey: string, channel: Channel) => {
    const nextSettings = {
      ...userSettings,
      [settingsKey]: (parseInt(userSettings[settingsKey], 10) ^ channel).toString(),
    }
    setUserSettings(nextSettings)
    props.settingsChanged(nextSettings)
  }

  const labelWeb = i18n._({ id: "mysettings.notification.channelweb", message: "Web" })
  const labelEmail = i18n._({ id: "mysettings.notification.channelemail", message: "Email" })

  const icon = (settingsKey: string, channel: Channel) => {
    const active = isEnabled(settingsKey, channel)
    const label = channel === WebChannel ? labelWeb : labelEmail
    const onToggle = () => toggle(settingsKey, channel)
    return <Toggle key={`${settingsKey}_${channel}`} active={active} label={label} onToggle={onToggle} />
  }

  return (
    <>
      <Field label={i18n._({ id: "label.notifications", message: "Notifications" })}>
        <p className="text-muted mb-6">
          <Trans id="mysettings.notification.title">Choose the events to recieve a notification for.</Trans>
        </p>

        <div className="notifications-settings mt-4">
          <VStack spacing={4} divide={true} className="rounded">
            <div>
              <HStack spacing={6} justify="between">
                <span>
                  <Trans id="mysettings.notification.event.newpost">New Post</Trans>
                </span>
                <HStack spacing={6}>
                  {icon("event_notification_new_post", WebChannel)}
                  {icon("event_notification_new_post", EmailChannel)}
                </HStack>
              </HStack>
            </div>
            <div>
              <HStack spacing={6} justify="between">
                <span className="mb-1">
                  <Trans id="mysettings.notification.event.discussion">New Comments</Trans>
                </span>
                <HStack spacing={6}>
                  {icon("event_notification_new_comment", WebChannel)}
                  {icon("event_notification_new_comment", EmailChannel)}
                </HStack>
              </HStack>
            </div>
            <div>
              <HStack spacing={6} justify="between">
                <div className="mb-1">
                  <Trans id="mysettings.notification.event.mention">Mentions</Trans>
                </div>
                <HStack spacing={6}>
                  {icon("event_notification_mention", WebChannel)}
                  {icon("event_notification_mention", EmailChannel)}
                </HStack>
              </HStack>
            </div>
            <div>
              <HStack spacing={6} justify="between">
                <span className="mb-1">
                  <Trans id="mysettings.notification.event.statuschanged">Status Changed</Trans>
                </span>
                <HStack spacing={6}>
                  {icon("event_notification_change_status", WebChannel)}
                  {icon("event_notification_change_status", EmailChannel)}
                </HStack>
              </HStack>
            </div>
          </VStack>
        </div>
      </Field>
    </>
  )
}
