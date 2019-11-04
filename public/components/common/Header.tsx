import "./Header.scss";

import React from "react";
import { SignInModal, EnvironmentInfo, Avatar, TenantLogo, TenantStatusInfo } from "@fider/components";
import { actions, Fider } from "@fider/services";
import { FaUser, FaCog, FaCaretDown } from "react-icons/fa";

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
    if (Fider.session.isAuthenticated) {
      actions.getTotalUnreadNotifications().then(result => {
        if (result.ok && result.data > 0) {
          this.setState({ unreadNotifications: result.data });
        }
      });
    }
  }

  private showModal = () => {
    if (!Fider.session.isAuthenticated) {
      this.setState({ showSignIn: true });
    }
  };

  private hideModal = () => {
    this.setState({ showSignIn: false });
  };

  public render() {
    const items = Fider.session.isAuthenticated && (
      <div className="c-menu-user">
        <div className="c-menu-user-heading">
          <FaUser /> <span>{Fider.session.user.name}</span>
        </div>
        <a href="/settings" className="c-menu-user-item">
          Settings
        </a>
        <a href="/notifications" className="c-menu-user-item">
          Notifications
          {this.state.unreadNotifications > 0 && <div className="c-unread-count">{this.state.unreadNotifications}</div>}
        </a>
        <div className="c-menu-user-divider" />
        {Fider.session.user.isCollaborator && [
          <div key={1} className="c-menu-user-heading">
            <FaCog /> <span>Administration</span>
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

    const showRightMenu = Fider.session.isAuthenticated || !Fider.session.tenant.isPrivate;
    return (
      <div id="c-header">
        <EnvironmentInfo />
        <SignInModal isOpen={this.state.showSignIn} onClose={this.hideModal} />
        <div className="c-menu">
          <div className="container">
            <a href="/" className="c-menu-item-title">
              <TenantLogo size={100} />
              <span>{Fider.session.tenant.name}</span>
            </a>
            {showRightMenu && (
              <div onClick={this.showModal} className="c-menu-item-signin">
                {Fider.session.isAuthenticated && <Avatar user={Fider.session.user} />}
                {this.state.unreadNotifications > 0 && <div className="c-unread-dot" />}
                {!Fider.session.isAuthenticated && <span>Sign in</span>}
                {Fider.session.isAuthenticated && <FaCaretDown />}
                {items}
              </div>
            )}
          </div>
        </div>
        <TenantStatusInfo />
      </div>
    );
  }
}
