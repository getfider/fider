import * as React from 'react';
import { SocialSignInButton, Form } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { AuthSettings } from '@fider/models';
import { TenantService } from '@fider/services';
import { hideSignIn } from '@fider/utils/page';

interface SignInControlState {
  email: string;
  sent: boolean;
}

export class SignInControl extends React.Component<{}, SignInControlState> {
  private form: Form;

  @inject(injectables.Session)
  public session: Session;

  @inject(injectables.TenantService)
  public service: TenantService;

  constructor() {
    super();
    this.state = {
      email: '',
      sent: false
    };
  }

  private async signIn() {
    const result = await this.service.signIn(this.state.email);
    if (result.ok) {
      this.form.clearFailure();
      this.setState({ sent: true });
      setTimeout(() => {
        this.setState({ sent: false });
      }, 10000);
    } else if (result.error) {
      this.form.setFailure(result.error);
    }
  }

  public render() {
    const settings = this.session.get<AuthSettings>('auth');

    if (this.state.sent) {
      return <div>
                <p>We sent a sign in link to <b>{ this.state.email }</b>. <br /> Please check your inbox.</p>
                <p><a href="#" onClick={() => hideSignIn()}>OK</a></p>
             </div>;
    }

    const google = settings.providers.google &&
                    <div className="column">
                        <SocialSignInButton provider="google" />
                    </div>;
    const facebook = settings.providers.facebook &&
                    <div className="column">
                        <SocialSignInButton provider="facebook" />
                    </div>;
    const github = settings.providers.github &&
                    <div className="column">
                        <SocialSignInButton provider="github" />
                    </div>;

    return  <div className="signin-options">
                <div className="ui stackable three column centered grid">
                  { facebook }
                  { google }
                  { github }
                </div>
                <p className="info">We'll never post to any of your accounts</p>
                <div className="ui horizontal divider">OR</div>
                <p>Enter your e-mail address to sign in</p>
                <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.signIn() }>
                  <div className="ui small action fluid input">
                      <input onChange={(e) => this.setState({ email: e.currentTarget.value }) } type="text" placeholder="yourname@example.com" className="small" />
                      <button onClick={ () => this.form.submit() } className={`ui small positive button ${this.state.email === '' && 'disabled'}`}>
                        Sign in
                      </button>
                  </div>
                </Form>
            </div>;
  }
}
