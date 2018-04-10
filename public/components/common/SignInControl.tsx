import * as React from "react";
import { SocialSignInButton, Form, Button } from "@fider/components/common";
import { AuthSettings } from "@fider/models";
import { page, device, actions } from "@fider/services";

interface SignInControlState {
  email: string;
}

interface SignInControlProps {
  useEmail: boolean;
  redirectTo?: string;
  onEmailSent?: (email: string) => void;
}

export class SignInControl extends React.Component<SignInControlProps, SignInControlState> {
  private form!: Form;
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

  private async signIn() {
    const result = await actions.signIn(this.state.email);
    if (result.ok) {
      this.form.clearFailure();
      if (this.props.onEmailSent) {
        this.props.onEmailSent(this.state.email);
      }
    } else if (result.error) {
      this.form.setFailure(result.error);
    }
  }

  public render() {
    const google = this.settings.providers.google && (
      <div className="column">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="google"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const facebook = this.settings.providers.facebook && (
      <div className="column">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="facebook"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const github = this.settings.providers.github && (
      <div className="column">
        <SocialSignInButton
          oauthEndpoint={this.settings.endpoint}
          provider="github"
          redirectTo={this.props.redirectTo}
        />
      </div>
    );
    const hasOAuth = !!(google || facebook || github);

    return (
      <div className="signin-options">
        {hasOAuth && (
          <div>
            <div className="ui stackable three column centered grid">
              {facebook}
              {google}
              {github}
            </div>
            <p className="info">We'll never post to any of your accounts</p>
            <div className="ui horizontal divider">OR</div>
          </div>
        )}

        {this.props.useEmail && (
          <div>
            <p>Enter your email address to sign in</p>
            <Form
              ref={f => {
                this.form = f!;
              }}
            >
              <div id="email-signin" className="ui small action fluid input">
                <input
                  autoFocus={!device.isTouch()}
                  onChange={e => this.setState({ email: e.currentTarget.value })}
                  onKeyDown={this.onEmailKeyDown}
                  type="text"
                  placeholder="yourname@example.com"
                  className="small"
                />
                <Button color="green" disabled={this.state.email === ""} onClick={() => this.signIn()}>
                  Sign in
                </Button>
              </div>
            </Form>
          </div>
        )}
      </div>
    );
  }
}
