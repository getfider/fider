import './SignUp.page.scss';

import * as React from 'react';
import { SignInControl, Modal, Button, DisplayError } from '@fider/components/common';
import { SystemSettings } from '@fider/models';
import { jwt, page, actions, Failure } from '@fider/services';

const logo = require('@fider/assets/images/logo-small.png');

interface OAuthUser {
  token: string;
  name: string;
  email: string;
}

interface SignUpPageProps {
  system: SystemSettings;
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

export class SignUpPage extends React.Component<SignUpPageProps, SignUpPageState> {
    private user?: OAuthUser;

    constructor(props: SignUpPageProps) {
      super(props);
      this.state = {
        submitted: false,
        tenantName: '',
        subdomain: { available: false }
      };

      page.setTitle(page.isSingleHostMode() ? 'Installation · Fider' : 'Sign up · Fider');

      const token = page.getQueryString('token');
      if (token) {
        const data = jwt.decode(token);
        if (data) {
          this.user = {
            token,
            name: data['oauth/name'],
            email: data['oauth/email']
          };
        }
      }
    }

    private async confirm() {
      const result = await actions.createTenant({
        token: this.user && this.user.token,
        tenantName: this.state.tenantName,
        subdomain: this.state.subdomain.value,
        name: this.state.name,
        email: this.state.email,
      });

      if (result.ok) {
        if (result.data.token) {
          if (page.isSingleHostMode()) {
            location.reload();
          } else {
            let baseUrl = `${location.protocol}//${this.state.subdomain.value}${this.props.system.domain}`;
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
      const result = await actions.checkAvailability(subdomain);

      this.setState({
        subdomain: {
          value: subdomain,
          available: !result.data.message,
          message: result.data.message,
        }
      });
    }

    public render() {
      const modal = (
        <Modal.Window canClose={false} isOpen={this.state.submitted}>
          <Modal.Header>Thank you for registering!</Modal.Header>
          <Modal.Content>
            <p>We have just sent a confirmation link to <b>{this.state.email}</b>. <br /> Click the link to finish your registration.</p>
          </Modal.Content>
        </Modal.Window>
      );

      return (
        <div className="page ui container">
          {modal}
          <img className="logo" src={logo} />

          <h3 className="ui header">1. Who are you?</h3>
          <DisplayError fields={['token']} error={this.state.error} />

          {
            this.user ?
            <p>
              Hello, &nbsp;
              <b>{this.user.name}</b> {this.user.email && `(${this.user.email})`}
            </p> :
            <>
              <p>We need to identify you to setup your new Fider account.</p>
              <SignInControl useEmail={false} />
              <div className="ui form">
                <DisplayError fields={['name', 'email']} error={this.state.error} />
                <div className="fluid field">
                  <input id="name" onChange={(e) => this.setState({ name: e.currentTarget.value })} type="text" placeholder="your name" className="small" />
                </div>
                <div className="fluid field">
                  <input id="email" onChange={(e) => this.setState({ email: e.currentTarget.value })} type="text" placeholder="yourname@example.com" className="small" />
                </div>
              </div>
            </>
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
              !page.isSingleHostMode() && <div className="fluid field">
                <div className="ui right labeled input">
                  <input
                    id="subdomain"
                    type="text"
                    maxLength={40}
                    placeholder="subdomain"
                    onChange={(e) => this.checkAvailability(e.currentTarget.value)}
                  />
                  <div className="ui label">{this.props.system.domain}</div>
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

          <Button color="green" size="large" onClick={() => this.confirm()}>Confirm</Button>
        </div>
      );
    }
}
