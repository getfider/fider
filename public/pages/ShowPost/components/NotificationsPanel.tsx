import React, { useState } from "react";
import { Post } from "@fider/models";
import { Button, List, ListItem } from "@fider/components";
import { actions } from "@fider/services";
import { FaVolumeUp, FaVolumeMute } from "react-icons/fa";
import { useFider } from "@fider/hooks";
import { useTranslation } from "react-i18next";

interface NotificationsPanelProps {
  post: Post;
  subscribed: boolean;
}

export const NotificationsPanel = (props: NotificationsPanelProps) => {
  const fider = useFider();
  const { t } = useTranslation();
  const [subscribed, setSubscribed] = useState(props.subscribed);

  const subscribeOrUnsubscribe = async () => {
    const action = subscribed ? actions.unsubscribe : actions.subscribe;

    const response = await action(props.post.number);
    if (response.ok) {
      setSubscribed(!subscribed);
    }
  };

  if (!fider.session.isAuthenticated) {
    return null;
  }

  const button = subscribed ? (
    <Button fluid={true} onClick={subscribeOrUnsubscribe}>
      <FaVolumeMute /> {t("common.button.unsubscribe")}
    </Button>
  ) : (
    <Button fluid={true} onClick={subscribeOrUnsubscribe}>
      <FaVolumeUp /> {t("common.button.subscribe")}
    </Button>
  );

  const text = subscribed ? (
    <span className="info">{t("showPost.notificationPanel.receivingNotifications")}</span>
  ) : (
    <span className="info">{t("showPost.notificationPanel.notReceiveNotification")}</span>
  );

  return (
    <>
      <span className="subtitle">{t("showPost.notificationPanel.subtitle")}</span>
      <List>
        <ListItem>
          {button}
          {text}
        </ListItem>
      </List>
    </>
  );
};
