import * as React from 'react';
import { AppSettings, CurrentUser, Tenant } from '@fider/models';
import { SignInControl, EnvironmentInfo, Gravatar } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { showSignIn } from '@fider/utils/page';

interface HeaderProps {
  user?: CurrentUser;
  settings: AppSettings;
  tenant: Tenant;
}

export class Header extends React.Component<HeaderProps, {}> {

  constructor(props: HeaderProps) {
    super(props);
  }

  private showModal() {
    if (!this.props.user) {
      showSignIn();
    }
  }

  public render() {
    const items = this.props.user && (
      <div className="menu">
          <div className="name header">
            <i className="user icon" />
            {this.props.user.name}
          </div>
          <a href="/settings" className="item">Settings</a>
          <div className="divider" />
          {
            this.props.user.isCollaborator && [
              <div key={1} className="header">
                <i className="setting icon" />
                Administration
              </div>,
              <a key={2} href="/admin" className="item">General Settings</a>,
              <a key={3} href="/admin/members" className="item">Members</a>,
              <a key={4} href="/admin/tags" className="item">Tags</a>,
              <div key={5} className="divider" />
            ]
          }
          <a href="/signout?redirect=/" className="item signout">Sign out</a>
      </div>
    );

    return (
      <div>
        <EnvironmentInfo settings={this.props.settings}/>
        <div id="menu" className="ui small borderless menu">
          <div className="ui container">
            <a href="/" className="header item">
              {this.props.tenant.name}
            </a>
            <div onClick={() => this.showModal()} className={`ui right simple dropdown item signin ${!this.props.user && 'subtitle'}`}>
              {this.props.user && <Gravatar user={this.props.user} />}
              {!this.props.user && 'Sign in'} {this.props.user && <i className="dropdown icon" />}
              {items}
            </div>
          </div>
        </div>
      </div>
    );

    }
}
