import React from "react";
import { Post } from "@fider/models";
import { Button, List, ListItem } from "@fider/components";
import { actions, Fider } from "@fider/services";
import { FaVolumeUp, FaVolumeMute } from "react-icons/fa";

interface NotificationsPanelProps {
  post: Post;
  subscribed: boolean;
}

interface NotificationsPanelState {
  subscribed: boolean;
}

export class NotificationsPanel extends React.Component<NotificationsPanelProps, NotificationsPanelState> {
  constructor(props: NotificationsPanelProps) {
    super(props);
    this.state = {
      subscribed: this.props.subscribed
    };
  }

  private subscribeOrUnsubscribe = async () => {
    const action = this.state.subscribed ? actions.unsubscribe : actions.subscribe;

    const response = await action(this.props.post.number);
    if (response.ok) {
      this.setState(state => ({
        subscribed: !state.subscribed
      }));
    }
  };

  public render() {
    if (!Fider.session.isAuthenticated) {
      return null;
    }

    const button = this.state.subscribed ? (
      <Button fluid={true} onClick={this.subscribeOrUnsubscribe}>
        <FaVolumeMute /> Unsubscribe
      </Button>
    ) : (
      <Button fluid={true} onClick={this.subscribeOrUnsubscribe}>
        <FaVolumeUp /> Subscribe
      </Button>
    );

    const text = this.state.subscribed ? (
      <span className="info">You’re receiving notifications about activity on this post.</span>
    ) : (
      <span className="info">You'll not receive any notification about this post.</span>
    );

    return (
      <>
        <span className="subtitle">Notifications</span>
        <List>
          <ListItem>
            {button}
            {text}
          </ListItem>
        </List>
      </>
    );
  }
}
