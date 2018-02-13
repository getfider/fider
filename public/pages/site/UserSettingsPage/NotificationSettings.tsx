import * as React from 'react';

import { CurrentUser, UserSettings } from '@fider/models';

interface NotificationSettingsProps {
  user: CurrentUser;
  settings: UserSettings;
}

interface NotificationSettingsState {
  newIdeaWeb: boolean;
  newIdeaEmail: boolean;
  newCommentWeb: boolean;
  newCommentEmail: boolean;
  changeStatusWeb: boolean;
  changeStatusEmail: boolean;
}

const WebChannel = 1;
const EmailChannel = 1;

export class NotificationSettings extends React.Component<NotificationSettingsProps, NotificationSettingsState> {
  constructor(props: NotificationSettingsProps) {
    super(props);

    const isEnabled = (key: string, flag: number): boolean => {
      if (key in this.props.settings) {
        return (parseInt(this.props.settings[key], 10) & flag) > 0;
      }
      return false;
    };

    this.state = {
      newIdeaWeb: isEnabled('event_notification_new_idea', WebChannel),
      newIdeaEmail: isEnabled('event_notification_new_idea', EmailChannel),
      newCommentWeb: isEnabled('event_notification_new_comment', WebChannel),
      newCommentEmail: isEnabled('event_notification_new_comment', EmailChannel),
      changeStatusWeb: isEnabled('event_notification_change_status', WebChannel),
      changeStatusEmail: isEnabled('event_notification_change_status', EmailChannel),
    };
  }

  private icon(key: keyof NotificationSettingsState) {
    return this.state[key]
    ? <i className="green large toggle on icon" onClick={() => this.setState({ [key]: false } as any)} />
    : <i className="large toggle off icon" onClick={() => this.setState({ [key]: true } as any)} />;
  }

  private info(web: keyof NotificationSettingsState, email: keyof NotificationSettingsState, aboutForVisitors: string, aboutForCollaborators: string) {
    const about = this.props.user.isCollaborator ? aboutForCollaborators : aboutForVisitors;
    if (!this.state[web] && !this.state[email]) {
      return <span className="info">You'll <strong>NOT</strong> receive any notification about this event.</span>;
    }
    if (this.state[web] && !this.state[email]) {
      return <span className="info">You'll receive <strong>web</strong> notifications about {about}.</span>;
    }
    if (!this.state[web] && this.state[email]) {
      return <span className="info">You'll receive <strong>e-mail</strong> notifications about {about}.</span>;
    }
    if (this.state[web] && this.state[email]) {
      return <span className="info">You'll receive <strong>web</strong> and <strong>e-mail</strong> notifications about {about}.</span>;
    }
    return null;
  }

  public render() {
    return (
      <>
        <div className="field">
          <label htmlFor="notifications">Notifications</label>
        </div>
        <table className="ui very basic table notifications-settings">
          <thead>
            <tr>
              <th className="five wide">Event</th>
              <th className="two wide">Web <span className="info">*</span></th>
              <th className="two wide">E-mail</th>
              <th className="seven wide"><span className="info">* web notifications are comming soon!</span></th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>New Ideas</td>
              <td>{this.icon('newIdeaWeb')}</td>
              <td>{this.icon('newIdeaEmail')}</td>
              <td>{this.info('newIdeaWeb', 'newIdeaEmail', 'new ideas posted on this site', 'new ideas posted on this site')}</td>
            </tr>
            <tr>
              <td>Discussion</td>
              <td>{this.icon('newCommentWeb')}</td>
              <td>{this.icon('newCommentEmail')}</td>
              <td>{this.info('newCommentWeb', 'newCommentEmail', 'comments on ideas you\'ve subscribed to', 'comments on all ideas unless individually unsubscribed')}</td>
            </tr>
            <tr>
              <td>Status Changed</td>
              <td>{this.icon('changeStatusWeb')}</td>
              <td>{this.icon('changeStatusEmail')}</td>
              <td>{this.info('changeStatusWeb', 'changeStatusEmail', 'status change on ideas you\'ve subscribed to', 'status change on all ideas unless individually unsubscribed')}</td>
            </tr>
          </tbody>
        </table>
      </>
    );
  }
}
