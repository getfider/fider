import "./MyNotifications.page.scss";

import React from "react";

import { Notification } from "@fider/models";
import { MultiLineText, Moment, Heading, List, ListItem } from "@fider/components";
import { actions } from "@fider/services";
import { FaBell } from "react-icons/fa";
import { withTranslation, WithTranslation } from "react-i18next";

interface MyNotificationsPageProps extends WithTranslation {
  notifications: Notification[];
}

interface MyNotificationsPageState {
  unread: Notification[];
  recent: Notification[];
}

class MyNotificationsPage extends React.Component<MyNotificationsPageProps, MyNotificationsPageState> {
  constructor(props: MyNotificationsPageProps) {
    super(props);

    const [unread, recent] = (this.props.notifications || []).reduce(
      (result, item) => {
        result[item.read ? 1 : 0].push(item);
        return result;
      },
      [[] as Notification[], [] as Notification[]]
    );

    this.state = {
      unread,
      recent
    };
  }

  private items(notifications: Notification[]): JSX.Element[] {
    return notifications.map(n => {
      return (
        <ListItem key={n.id}>
          <a href={`/notifications/${n.id}`}>
            <MultiLineText text={n.title} style="simple" />
            <span className="info">
              <Moment date={n.createdAt} />
            </span>
          </a>
        </ListItem>
      );
    });
  }

  private markAllAsRead = async () => {
    const response = await actions.markAllAsRead();
    if (response.ok) {
      location.reload();
    }
  };

  public render() {
    const { t } = this.props;
    return (
      <div id="p-my-notifications" className="page container">
        <Heading title={t("myNotifications.title")} subtitle={t("myNotifications.subtitle")} icon={FaBell} />

        <h4>
          {t("myNotifications.unread")}
          {this.state.unread.length > 0 && (
            <span className="mark-as-read" onClick={this.markAllAsRead}>
              {t("myNotifications.markAsRead")}
            </span>
          )}
        </h4>
        <List>
          {this.state.unread.length > 0 && this.items(this.state.unread)}
          {this.state.unread.length === 0 && <span className="info">{t("myNotifications.noUnread")}</span>}
        </List>
        {this.state.recent.length > 0 && (
          <>
            <h4>{t("myNotifications.readOnDays")}</h4>
            <List>
              <ListItem>{this.items(this.state.recent)}</ListItem>
            </List>
          </>
        )}
      </div>
    );
  }
}

export default withTranslation()(MyNotificationsPage);
