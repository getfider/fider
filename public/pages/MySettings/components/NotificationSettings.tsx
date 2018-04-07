import * as React from "react";

import { CurrentUser, UserSettings } from "@fider/models";
import { Toggle } from "@fider/components";

interface NotificationSettingsProps {
  user: CurrentUser;
  settings: UserSettings;
  settingsChanged: (settings: UserSettings) => void;
}

interface NotificationSettingsState {
  settings: UserSettings;
}

type Channel = number;
const WebChannel: Channel = 1;
const EmailChannel: Channel = 2;

export class NotificationSettings extends React.Component<NotificationSettingsProps, NotificationSettingsState> {
  constructor(props: NotificationSettingsProps) {
    super(props);

    this.state = {
      settings: this.props.settings
    };
  }

  private isEnabled(settingsKey: string, channel: Channel): boolean {
    if (settingsKey in this.state.settings) {
      return (parseInt(this.state.settings[settingsKey], 10) & channel) > 0;
    }
    return false;
  }

  private toggle(settingsKey: string, channel: Channel) {
    const settings = { ...this.state.settings };
    settings[settingsKey] = (parseInt(this.state.settings[settingsKey], 10) ^ channel).toString();

    this.setState({ settings });
    this.props.settingsChanged(settings);
  }

  private icon(settingsKey: string, channel: Channel) {
    const active = this.isEnabled(settingsKey, channel);
    const label = channel === WebChannel ? "Web" : "Email"; 
    return (
      <Toggle
        key={`${settingsKey}_${channel}`}
        active={active}
        label={label}
        onToggle={this.toggle.bind(this, settingsKey, channel)}
      />
    );
  }

  private info(settingsKey: string, aboutForVisitors: string, aboutForCollaborators: string) {
    const about = this.props.user.isCollaborator ? aboutForCollaborators : aboutForVisitors;
    const webEnabled = this.isEnabled(settingsKey, WebChannel);
    const emailEnabled = this.isEnabled(settingsKey, EmailChannel);

    if (!webEnabled && !emailEnabled) {
      return (
        <p className="info">
          You'll <strong>NOT</strong> receive any notification about this event.
        </p>
      );
    } else if (webEnabled && !emailEnabled) {
      return (
        <p className="info">
          You'll receive <strong>web</strong> notifications about {about}.
        </p>
      );
    } else if (!webEnabled && emailEnabled) {
      return (
        <p className="info">
          You'll receive <strong>email</strong> notifications about {about}.
        </p>
      );
    } else if (webEnabled && emailEnabled) {
      return (
        <p className="info">
          You'll receive <strong>web</strong> and <strong>email</strong> notifications about {about}.
        </p>
      );
    }
    return null;
  }

  public render() {
    return (
      <>
        <div className="field">
          <label htmlFor="notifications">Notifications</label>
          <p className="info">Use following panel to choose which events you'd like to receive notification</p>
        </div>

        <div className="ui segments notifications-settings">
          <div className="ui segment">
            <span className="event-title">New Idea</span>
            {this.info("event_notification_new_idea", "new ideas posted on this site", "new ideas posted on this site")}
            <p>
              {this.icon("event_notification_new_idea", WebChannel)}
              {this.icon("event_notification_new_idea", EmailChannel)}
            </p>
          </div>
          <div className="ui segment">
            <span className="event-title">Discussion</span>
            {this.info(
              "event_notification_new_comment",
              "comments on ideas you've subscribed to",
              "comments on all ideas unless individually unsubscribed"
            )}
            <p>
              {this.icon("event_notification_new_comment", WebChannel)}
              {this.icon("event_notification_new_comment", EmailChannel)}
            </p>
          </div>
          <div className="ui segment">
            <span className="event-title">Status Changed</span>
            {this.info(
              "event_notification_change_status",
              "status change on ideas you've subscribed to",
              "status change on all ideas unless individually unsubscribed"
            )}
            <p>
              {this.icon("event_notification_change_status", WebChannel)}
              {this.icon("event_notification_change_status", EmailChannel)}
            </p>
          </div>
        </div>
      </>
    );
  }
}
