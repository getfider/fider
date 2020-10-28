import React, { useState } from "react";

import { UserSettings } from "@fider/models";
import { Toggle, Segment, Segments, Field } from "@fider/components";
import { useFider } from "@fider/hooks";
import { useTranslation } from "react-i18next";

interface NotificationSettingsProps {
  userSettings: UserSettings;
  settingsChanged: (settings: UserSettings) => void;
}

type Channel = number;
const WebChannel: Channel = 1;
const EmailChannel: Channel = 2;

export const NotificationSettings = (props: NotificationSettingsProps) => {
  const { t } = useTranslation();
  const fider = useFider();
  const [userSettings, setUserSettings] = useState(props.userSettings);

  const isEnabled = (settingsKey: string, channel: Channel): boolean => {
    if (settingsKey in userSettings) {
      return (parseInt(userSettings[settingsKey], 10) & channel) > 0;
    }
    return false;
  };

  const toggle = async (settingsKey: string, channel: Channel) => {
    const nextSettings = {
      ...userSettings,
      [settingsKey]: (parseInt(userSettings[settingsKey], 10) ^ channel).toString()
    };
    setUserSettings(nextSettings);
    props.settingsChanged(nextSettings);
  };

  const icon = (settingsKey: string, channel: Channel) => {
    const active = isEnabled(settingsKey, channel);
    const label = channel === WebChannel ? "Web" : "Email";
    const onToggle = () => toggle(settingsKey, channel);
    return <Toggle key={`${settingsKey}_${channel}`} active={active} label={label} onToggle={onToggle} />;
  };

  const info = (settingsKey: string, aboutForVisitors: string, aboutForCollaborators: string) => {
    const about = fider.session.user.isCollaborator ? aboutForCollaborators : aboutForVisitors;
    const webEnabled = isEnabled(settingsKey, WebChannel);
    const emailEnabled = isEnabled(settingsKey, EmailChannel);

    if (!webEnabled && !emailEnabled) {
      return <p className="info">{t("mySettings.notReceiveAnyNotification")}</p>;
    } else if (webEnabled && !emailEnabled) {
      return <p className="info">{t("mySettings.youWillReceiveWeb", { about })}</p>;
    } else if (!webEnabled && emailEnabled) {
      return <p className="info">{t("mySettings.youWillReceiveEmail", { about })}</p>;
    } else if (webEnabled && emailEnabled) {
      return <p className="info">{t("mySettings.youWillReceiveEmailAndWeb", { about })}</p>;
    }
    return null;
  };

  return (
    <>
      <Field label={t("notifications")}>
        <p className="info">{t("mySettings.chooseEvents")}</p>
      </Field>

      <div className="notifications-settings">
        <Segments>
          <Segment>
            <span className="event-title">{t("mySettings.newPostTitle")}</span>
            {info("event_notification_new_post", t("mySettings.newPostEvent"), t("mySettings.newPostEvent"))}
            <p>
              {icon("event_notification_new_post", WebChannel)}
              {icon("event_notification_new_post", EmailChannel)}
            </p>
          </Segment>
          <Segment>
            <span className="event-title">{t("mySettings.discussion")}</span>
            {info("event_notification_new_comment", t("mySettings.postSubscribedTo"), t("mySettings.postUnsubscribed"))}
            <p>
              {icon("event_notification_new_comment", WebChannel)}
              {icon("event_notification_new_comment", EmailChannel)}
            </p>
          </Segment>
          <Segment>
            <span className="event-title">{t("mySettings.statusChanged")}</span>
            {info(
              "event_notification_change_status",
              t("mySettings.statusSubscribed"),
              t("mySettings.statusUnsubscribed")
            )}
            <p>
              {icon("event_notification_change_status", WebChannel)}
              {icon("event_notification_change_status", EmailChannel)}
            </p>
          </Segment>
        </Segments>
      </div>
    </>
  );
};
