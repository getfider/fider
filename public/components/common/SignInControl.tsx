import "./SignInControl.scss";

import React from "react";
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
    const providersLen = Fider.settings.oauth.length;

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
        {providersLen > 0 && (
          <div className="l-signin-social">
            <div className="row">
              {Fider.settings.oauth.map((o, i) => (
                <React.Fragment key={o.provider}>
                  {i % 4 === 0 && <div className="col-lf" />}
                  <div className={`col-sm l-social-col ${providersLen === 1 ? "l-social-col-100" : ""}`}>
                    <SocialSignInButton option={o} redirectTo={this.props.redirectTo} />
                  </div>
                </React.Fragment>
              ))}
            </div>
            <p className="info">We will never post to these accounts on your behalf.</p>
          </div>
        )}

        {providersLen > 0 && <div className="c-divider">OR</div>}

        {this.props.useEmail && (
          <div className="l-signin-email">
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
