import "./SignInControl.scss";

import * as React from "react";
import { SocialSignInButton, Form, Button, Input } from "@fider/components";
import { device, actions, Failure, Fider } from "@fider/services";

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
  constructor(props: SignInControlProps) {
    super(props);
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

  private setEmail = (email: string) => {
    this.setState({ email });
  };

  public render() {
    const google = Fider.settings.auth.providers.google && (
      <div className="col-sm l-social-col">
        <SocialSignInButton
          oauthEndpoint={Fider.settings.auth.endpoint}
          provider="google"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const facebook = Fider.settings.auth.providers.facebook && (
      <div className="col-sm l-social-col">
        <SocialSignInButton
          oauthEndpoint={Fider.settings.auth.endpoint}
          provider="facebook"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const github = Fider.settings.auth.providers.github && (
      <div className="col-sm l-social-col">
        <SocialSignInButton
          oauthEndpoint={Fider.settings.auth.endpoint}
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
            <p className="info">We will never post to these accounts on your behalf.</p>
            <div className="c-divider">OR</div>
          </div>
        )}

        {this.props.useEmail && (
          <div>
            <p>Enter your email address to sign in</p>
            <Form error={this.state.error}>
              <Input
                field="email"
                value={this.state.email}
                autoFocus={!device.isTouch()}
                onChange={this.setEmail}
                onSubmit={this.signIn}
                placeholder="yourname@example.com"
                suffix={
                  <Button color="positive" disabled={this.state.email === ""} onClick={this.signIn}>
                    Sign in
                  </Button>
                }
              />
            </Form>
          </div>
        )}
      </div>
    );
  }
}
