import * as React from 'react';

import { SignInControl, Footer, Button, EnvironmentInfo, Gravatar, DisplayError } from '@fider/components/common';
import { AppSettings } from '@fider/models';
import { setTitle, getQueryString } from '@fider/utils/page';
import { decode } from '@fider/utils/jwt';
const td = require('throttle-debounce');
const logo = require('@fider/assets/images/logo-small.png');

import { inject, injectables } from '@fider/di';
import { Session, TenantService, Failure } from '@fider/services';
import { showModal } from '@fider/utils/page';

import './SignUpPage.scss';

interface OAuthUser {
    token: string;
    name: string;
    email: string;
}

interface SignUpPageState {
    submitted: boolean;
    tenantName: string;
    error?: Failure;
    name?: string;
    email?: string;
    subdomain: {
        available: boolean,
        message?: string,
        value?: string
    };
}

export class SignUpPage extends React.Component<{}, SignUpPageState> {
    private settings: AppSettings;
    private user: OAuthUser;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.TenantService)
    public service: TenantService;

    constructor(props: {}) {
      super(props);
      this.state = {
        submitted: false,
        tenantName: '',
        subdomain: { available: false }
      };

      this.settings = this.session.getAppSettings();
      this.checkAvailability = td.debounce(300, this.checkAvailability);

      setTitle(this.session.isSingleHostMode() ? 'Installation · Fider' : 'Sign up · Fider');

      const token = getQueryString('token');
      if (token) {
        const data = decode(token);
        if (data) {
          this.user = {
            token,
            name: data['oauth/name'],
            email: data['oauth/email']
          };
        }
      }
    }

    public componentDidUpdate() {
      if (this.state.submitted) {
        showModal('#submitted-modal', { closable: false });
      }
    }

    private async confirm() {
      const result = await this.service.create({
        token: this.user && this.user.token,
        tenantName: this.state.tenantName,
        subdomain: this.state.subdomain.value,
        name: this.state.name,
        email: this.state.email,
      });

      if (result.ok) {
        if (result.data.token) {
          if (this.session.isSingleHostMode()) {
            location.reload();
          } else {
            let baseUrl = `${location.protocol}//${this.state.subdomain.value}${this.settings.domain}`;
            if (location.port) {
              baseUrl = `${baseUrl}:${location.port}`;
            }

            location.href = `${baseUrl}?token=${result.data.token}`;
          }
        } else {
          this.setState({ submitted: true });
        }
      } else if (result.error) {
        this.setState({ error: result.error, submitted: false });
      }
    }

    private async checkAvailability(subdomain: string) {
      const result = await this.service.checkAvailability(subdomain);

      this.setState({
        subdomain: {
          value: subdomain,
          available: !result.data.message,
          message: result.data.message,
        }
      });
    }

    public render() {
      const modal = this.state.submitted && (
        <div id="submitted-modal" className="ui modal small">
          <div className="header">Thank you for registering!</div>
          <div className="content">
            <p>Please check your inbox for a confirmation link to finish your registration.</p>
          </div>
        </div>
      );

      return (
        <div>
          <EnvironmentInfo />
          <div className="page ui container">
            {modal}
            <img className="logo" src={logo} />

            <h3 className="ui header">1. Who are you?</h3>
            <DisplayError fields={['token']} error={this.state.error} />

            {
              this.user ?
              <div>
                <p>
                  Hello, &nbsp;
                  <Gravatar name={this.user.name} email={this.user.email} />
                  <b>{this.user.name}</b> ({this.user.email})
                </p>
              </div> :
              <div>
                <p>We need to identify you to setup your new Fider account.</p>
                <SignInControl signInByEmail={false} />
                  <div className="ui form">
                  <DisplayError fields={['name', 'email']} error={this.state.error} />
                  <div className="fluid field">
                    <input id="name" onChange={(e) => this.setState({ name: e.currentTarget.value })} type="text" placeholder="your name" className="small" />
                  </div>
                  <div className="fluid field">
                    <input id="email" onChange={(e) => this.setState({ email: e.currentTarget.value })} type="text" placeholder="yourname@example.com" className="small" />
                  </div>
                </div>
              </div>
            }

            <div className="ui section divider" />

            <h3 className="ui header">2. What is this Feedback Forum for?</h3>

            <DisplayError fields={['tenantName', 'subdomain']} error={this.state.error} />
            <div className="ui form">
              <div className="fluid field">
                <input
                  id="tenantName"
                  type="text"
                  placeholder="your company or product name"
                  maxLength={60}
                  onChange={(e) => this.setState({ tenantName: e.currentTarget.value })}
                />
              </div>
              {
                !this.session.isSingleHostMode() && <div className="fluid field">
                  <div className="ui right labeled input">
                    <input
                      id="subdomain"
                      type="text"
                      maxLength={40}
                      placeholder="subdomain"
                      onChange={(e) => this.checkAvailability(e.currentTarget.value)}
                    />
                    <div className="ui label">{this.settings.domain}</div>
                    {
                      this.state.subdomain.available &&
                      <div className="ui left pointing green basic label">
                        Great!
                      </div>
                    }
                    {
                      this.state.subdomain.message &&
                      <div className="ui left pointing red basic label">
                        {this.state.subdomain.message}
                      </div>
                    }
                  </div>
                </div>
              }
            </div>
            <div className="ui section divider" />

            <h3 className="ui header">3. Review and continue</h3>

            <p>Make sure information provided above is correct before proceeding.</p>

            <Button className="positive" size="large" onClick={() => this.confirm()}>Confirm</Button>
          </div>
          <Footer />
        </div>
      );
    }
}
