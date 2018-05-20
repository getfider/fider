import "./SignInControl.scss";

import * as React from "react";
import { SocialSignInButton, Form, Button, Form2, Input } from "@fider/components";
import { AuthSettings } from "@fider/models";
import { page, device, actions, Failure } from "@fider/services";

interface SignInControlState {
  email: string;
  error?: Failure;
}

interface SignInControlProps {
  useEmail: boolean;
  redirectTo?: string;
  onEmailSent?: (email: string) => void;
}

export class SignInControl extends React.Component<SignInControlProps, SignInControlState> {
  private settings: AuthSettings;

  constructor(props: SignInControlProps) {
    super(props);
    this.settings = page.authSettings();

    this.state = {
      email: ""
    };
  }

  private onEmailKeyDown = (event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.keyCode === 13) {
      // ENTER
      this.signIn();
      event.preventDefault();
    }
  };

  private signIn = async () => {
    const result = await actions.signIn(this.state.email);
    if (result.ok) {
      if (this.props.onEmailSent) {
        this.props.onEmailSent(this.state.email);
      }
      this.setState({ email: "", error: undefined });
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  };

  public render() {
    const google = this.settings.providers.google && (
      <div className="col-md-4 l-social-col">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="google"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const facebook = this.settings.providers.facebook && (
      <div className="col-md-4 l-social-col">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="facebook"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const github = this.settings.providers.github && (
      <div className="col-md-4 l-social-col">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="github"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const hasOAuth = !!(google || facebook || github);

    return (
      <div className="c-signin-control">
        {hasOAuth && (
          <div>
            <div className="row">
              {facebook}
              {google}
              {github}
            </div>
            <p className="info">We'll never post to any of your accounts</p>
            <div className="c-divider">OR</div>
          </div>
        )}

        {this.props.useEmail && (
          <div>
            <p>Enter your email address to sign in</p>
            <Form2 error={this.state.error}>
              <Input
                field="email"
                value={this.state.email}
                autoFocus={!device.isTouch()}
                onChange={email => this.setState({ email })}
                onSubmit={this.signIn}
                placeholder="yourname@example.com"
                suffix={
                  <Button color="positive" disabled={this.state.email === ""} onClick={this.signIn}>
                    Sign in
                  </Button>
                }
              />
            </Form2>
          </div>
        )}
      </div>
    );
  }
}
