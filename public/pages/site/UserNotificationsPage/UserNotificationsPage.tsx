import * as React from 'react';

import { CurrentUser, Notification } from '@fider/models';
import { MultiLineText, Moment } from '@fider/components';
import { page, actions } from '@fider/services';

import './UserNotificationsPage.scss';

interface UserNotificationsPageProps {
  notifications: Notification[];
}

interface UserNotificationsPageState {
  unread: Notification[];
  recent: Notification[];
}

export class UserNotificationsPage extends React.Component<UserNotificationsPageProps, UserNotificationsPageState> {

  constructor(props: UserNotificationsPageProps) {
    super(props);

    const [ unread, recent ] = (this.props.notifications || []).reduce((result, item) => {
      result[item.read ? 1 : 0].push(item);
      return result;
    }, [ [] as Notification[], [] as Notification[] ]);

    this.state = {
      unread,
      recent,
    };
  }

  private items(notifications: Notification[]): JSX.Element[] {
    return notifications.map((n) => {
      return (
        <span key={n.id} className="item">
          <a href={`/notifications/${n.id}`}>
            <MultiLineText text={n.title} style="simple" />
            <span className="info">
                <Moment date={n.createdOn} />
            </span>
          </a>
        </span>
      );
    });
  }

  private async markAllAsRead(): Promise<void> {
    const response = await actions.markAllAsRead();
    if (response.ok) {
      page.refresh();
    }
  }

  public render() {
    return (
      <div className="page ui container">
        <h2 className="ui header">
          <i className="circular id bell icon" />
          <div className="content">
            Notifications
            <div className="sub header">Stay up to date with what's happening</div>
          </div>
        </h2>

        <div className="ui grid">
          <div className="ten wide computer sixteen wide mobile column">
            <h4>
              Unread
              {
                this.state.unread.length > 0 &&
                <span
                  className="mark-as-read"
                  onClick={() => this.markAllAsRead()}
                >
                  Mark All as Read
                </span>
              }
            </h4>
            <div className="ui list">
              {this.state.unread.length > 0 && this.items(this.state.unread)}
              {this.state.unread.length === 0 && <span className="info">No unread notifications.</span>}
            </div>
            {
              this.state.recent.length > 0 &&
              <>
                <h4>Read on last 30 days</h4>
                <div className="ui list">
                  {this.items(this.state.recent)}
                </div>
              </>
            }
          </div>
        </div>

      </div>
    );
  }
}
