import "./SignInControl.scss";

import * as React from "react";
import { SocialSignInButton, Form, Button, Input, Message } from "@fider/components";
import { device, actions, Failure, Fider, isCookieEnabled } from "@fider/services";

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
    const oauthProviders = Fider.settings.oauth.map(o => (
      <div key={o.provider} className="col-md-3 col-lg-3 col-sm-3 col-xs-6 col-6 l-social-col">
        <SocialSignInButton option={o} redirectTo={this.props.redirectTo} />
      </div>
    ));

    if (!isCookieEnabled()) {
      return (
        <Message type="error">
          <h3>Cookies Required</h3>
          <p>Cookies are not enabled on your browser. Please enable cookies in your browser preferences to continue.</p>
        </Message>
      );
    }

    return (
      <div className="c-signin-control">
        {oauthProviders.length > 0 && (
          <div>
            <div className="row">{oauthProviders}</div>
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
