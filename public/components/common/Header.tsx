import "./Header.scss";

import * as React from "react";
import { SystemSettings, CurrentUser, Tenant } from "@fider/models";
import { SignInModal, SignInControl, EnvironmentInfo, Gravatar, Logo } from "@fider/components";
import { actions, classSet } from "@fider/services";

interface HeaderState {
  showSignIn: boolean;
  unreadNotifications: number;
}

export class Header extends React.Component<{}, HeaderState> {
  constructor(props: {}) {
    super(props);
    this.state = {
      showSignIn: false,
      unreadNotifications: 0
    };
  }

  public componentDidMount(): void {
    if (page.user) {
      actions.getTotalUnreadNotifications().then(result => {
        if (result.ok && result.data > 0) {
          this.setState({ unreadNotifications: result.data });
        }
      });
    }
  }

  private showModal = () => {
    if (!page.user) {
      this.setState({ showSignIn: true });
    }
  };

  public render() {
    const items = page.user && (
      <div className="c-menu-user">
        <div className="c-menu-user-heading">
          <i className="user icon" />
          {page.user.name}
        </div>
        <a href="/settings" className="c-menu-user-item">
          Settings
        </a>
        <a href="/notifications" className="c-menu-user-item">
          Notifications
          {this.state.unreadNotifications > 0 && <div className="c-unread-count">{this.state.unreadNotifications}</div>}
        </a>
        <div className="c-menu-user-divider" />
        {page.user.isCollaborator && [
          <div key={1} className="c-menu-user-heading">
            <i className="setting icon" />
            Administration
          </div>,
          <a key={2} href="/admin" className="c-menu-user-item">
            Site Settings
          </a>,
          <div key={5} className="c-menu-user-divider" />
        ]}
        <a href="/signout?redirect=/" className="c-menu-user-item signout">
          Sign out
        </a>
      </div>
    );

    const showRightMenu = page.user || !page.tenant.isPrivate;
    return (
      <div id="c-header">
        <EnvironmentInfo />
        <SignInModal isOpen={this.state.showSignIn} />
        <div className="c-menu">
          <div className="container">
            <a href="/" className="c-menu-item-title">
              <Logo size={100} />
              <span>{page.tenant.name}</span>
            </a>
            {showRightMenu && (
              <div onClick={this.showModal} className="c-menu-item-signin">
                {page.user && <Gravatar user={page.user} />}
                {this.state.unreadNotifications > 0 && <div className="c-unread-dot" />}
                {!page.user && <span>Sign in</span>} {page.user && <i className="dropdown icon" />}
                {items}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }
}
