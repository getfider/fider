import "./Header.scss";

import * as React from "react";
import { SystemSettings, CurrentUser, Tenant } from "@fider/models";
import { SignInModal, SignInControl, EnvironmentInfo, Gravatar } from "@fider/components";
import { page, actions } from "@fider/services";

interface HeaderProps {
  user?: CurrentUser;
  system: SystemSettings;
  tenant: Tenant;
}

interface HeaderState {
  showSignIn: boolean;
  unreadNotifications: number;
}

export class Header extends React.Component<HeaderProps, HeaderState> {
  constructor(props: HeaderProps) {
    super(props);
    this.state = {
      showSignIn: false,
      unreadNotifications: 0
    };
  }

  public componentDidMount(): void {
    if (this.props.user) {
      actions.getTotalUnreadNotifications().then(result => {
        if (result.ok && result.data > 0) {
          this.setState({ unreadNotifications: result.data });
        }
      });
    }
  }

  private showModal() {
    if (!this.props.user) {
      this.setState({ showSignIn: true });
    }
  }

  public render() {
    const items = this.props.user && (
      <div className="ui vertical menu">
        <div className="name header">
          <i className="user icon" />
          {this.props.user.name}
        </div>
        <a href="/settings" className="item">
          Settings
        </a>
        <a href="/notifications" className="item">
          Notifications
          {this.state.unreadNotifications > 0 && (
            <div className="ui mini circular red label">{this.state.unreadNotifications}</div>
          )}
        </a>
        <div className="divider" />
        {this.props.user.isCollaborator && [
          <div key={1} className="header">
            <i className="setting icon" />
            Administration
          </div>,
          <a key={2} href="/admin" className="item">
            Site Settings
          </a>,
          <div key={5} className="divider" />
        ]}
        <a href="/signout?redirect=/" className="item signout">
          Sign out
        </a>
      </div>
    );

    const showRightMenu = this.props.user || !this.props.tenant.isPrivate;

    return (
      <div className="c-header">
        <EnvironmentInfo system={this.props.system} />
        <SignInModal isOpen={this.state.showSignIn} />
        <div id="menu" className="ui small borderless menu">
          <div className="ui container">
            <a href="/" className="header item">
              {this.props.tenant.name}
            </a>
            {showRightMenu && (
              <div
                onClick={() => this.showModal()}
                className={`ui right simple dropdown item signin ${!this.props.user && "subtitle"}`}
              >
                {this.props.user && <Gravatar user={this.props.user} />}
                {this.state.unreadNotifications > 0 && <div className="unread-dot" />}
                {!this.props.user && "Sign in"} {this.props.user && <i className="dropdown icon" />}
                {items}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }
}
