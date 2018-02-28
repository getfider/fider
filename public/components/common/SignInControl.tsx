import * as React from 'react';
import { SocialSignInButton, Form, Button } from '@fider/components/common';
import { AuthSettings } from '@fider/models';
import { page, actions } from '@fider/services';

interface SignInControlState {
  email: string;
  sent: boolean;
}

interface SignInControlProps {
  signInByEmail: boolean;
  auth: AuthSettings;
}

export class SignInControl extends React.Component<SignInControlProps, SignInControlState> {
  private form!: Form;

  constructor(props: SignInControlProps) {
    super(props);

    this.state = {
      email: '',
      sent: false
    };
  }

  private async signIn() {
    const result = await actions.signIn(this.state.email);
    if (result.ok) {
      this.form.clearFailure();
      this.setState({ sent: true });
      setTimeout(() => {
        this.setState({ sent: false });
      }, 5000);
    } else if (result.error) {
      this.form.setFailure(result.error);
    }
  }

  public render() {
    if (this.state.sent) {
      return (
        <div>
          <p>We have just sent a confirmation link to <b>{this.state.email}</b>. <br /> Click the link and you’ll be signed in.</p>
          <p><a href="#" onClick={() => page.hideSignIn()}>OK</a></p>
        </div>
      );
    }

    const google = this.props.auth.providers.google && (
                    <div className="column">
                      <SocialSignInButton
                        oauthEndpoint={this.props.auth.endpoint}
                        provider="google"
                      />
                    </div>
                  );
    const facebook = this.props.auth.providers.facebook && (
                      <div className="column">
                        <SocialSignInButton
                          oauthEndpoint={this.props.auth.endpoint}
                          provider="facebook"
                        />
                      </div>
                    );
    const github = this.props.auth.providers.github && (
                      <div className="column">
                        <SocialSignInButton
                          oauthEndpoint={this.props.auth.endpoint}
                          provider="github"
                        />
                      </div>
                    );
    const hasOAuth = !!(google || facebook || github);

    return (
      <div className="signin-options">
          {
            hasOAuth && (
              <div>
                <div className="ui stackable three column centered grid">
                  {facebook}
                  {google}
                  {github}
                </div>
                <p className="info">We'll never post to any of your accounts</p>
                <div className="ui horizontal divider">OR</div>
              </div>
            )
          }

          { this.props.signInByEmail && <div>
            <p>Enter your email address to sign in</p>
            <Form ref={(f) => { this.form = f!; }}>
              <div id="email-signin" className="ui small action fluid input">
                  <input onChange={(e) => this.setState({ email: e.currentTarget.value })} type="text" placeholder="yourname@example.com" className="small" />
                  <Button onClick={() => this.signIn()} className={`positive ${this.state.email === '' && 'disabled'}`}>
                    Sign in
                  </Button>
              </div>
            </Form>
          </div> }
      </div>
    );
  }
}
