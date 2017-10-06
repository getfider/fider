import * as React from 'react';
import { User, Tenant } from '@fider/models';
import { SignInControl, EnvironmentInfo, Gravatar } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { showSignIn } from '@fider/utils/page';

export class Header extends React.Component<{}, {}> {
  private user?: User;
  private tenant: Tenant;

  @inject(injectables.Session)
  public session: Session;

  constructor(props: {}) {
    super(props);

    this.user = this.session.getCurrentUser();
    this.tenant = this.session.get<Tenant>('tenant');
  }

  private showModal() {
    if (!this.user) {
      showSignIn();
    }
  }

  public render() {
    const items = this.user && <div className="menu">
                      <div className="name header">
                        <i className="user icon"></i>
                        { this.user.name }
                      </div>
                      <a href="/settings" className="item">Settings</a>
                      <div className="divider"></div>
                      {
                        this.session.isStaff() && [
                          <div key={1} className="header">
                            <i className="setting icon"></i>
                            Administration
                          </div>,
                          <a key={2} href="/admin" className="item">General Settings</a>,
                          <div key={3} className="divider"></div>
                        ]
                      }
                      <a href="/signout?redirect=/" className="item signout">Sign out</a>
                  </div>;

    return <div>
              <EnvironmentInfo />
              <div id="menu" className="ui small borderless menu">
                <div className="ui container">
                  <a href="/" className="header item">
                    { this.tenant.name }
                  </a>
                  <div onClick={ () => this.showModal() } className={`ui right simple dropdown item signin ${!this.user && 'subtitle'}`}>
                    { this.user && <Gravatar name={this.user.name} hash={this.user.gravatar} /> }
                    { !this.user && 'Sign in' } { this.user && <i className="dropdown icon"></i> }
                    { items }
                  </div>
                </div>
              </div>
            </div>;
    }
}
