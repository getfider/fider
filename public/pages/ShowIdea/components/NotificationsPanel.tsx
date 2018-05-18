import * as React from "react";
import { CurrentUser, Idea } from "@fider/models";
import { Button, List, ListItem } from "@fider/components";
import { actions } from "@fider/services";

interface NotificationsPanelProps {
  user: CurrentUser | undefined;
  idea: Idea;
  subscribed: boolean;
}

interface NotificationsPanelState {
  subscribed: boolean;
}

export class NotificationsPanel extends React.Component<NotificationsPanelProps, NotificationsPanelState> {
  constructor(props: NotificationsPanelProps) {
    super(props);
    this.state = this.props;
  }

  private subscribeOrUnsubscribe = async () => {
    const action = this.state.subscribed ? actions.unsubscribe : actions.subscribe;

    const response = await action(this.props.idea.number);
    if (response.ok) {
      this.setState(state => ({
        subscribed: !state.subscribed
      }));
    }
  };

  public render() {
    if (!this.props.user) {
      return null;
    }

    const button = this.state.subscribed ? (
      <Button fluid={true} onClick={this.subscribeOrUnsubscribe}>
        <i className="volume off icon" /> Unsubscribe
      </Button>
    ) : (
      <Button fluid={true} onClick={this.subscribeOrUnsubscribe}>
        <i className="volume up icon" /> Subscribe
      </Button>
    );

    const text = this.state.subscribed ? (
      <span className="info">You’re receiving notifications about activity on this idea.</span>
    ) : (
      <span className="info">You'll not receive any notification about this idea.</span>
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
