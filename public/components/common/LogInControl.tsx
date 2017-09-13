import * as React from 'react';
import { SocialLogInButton, Form } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { AuthSettings } from '@fider/models';
import { TenantService } from '@fider/services';

interface LogInControlState {
  email: string;
  sent: boolean;
}

export class LogInControl extends React.Component<{}, LogInControlState> {
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

  private async login() {
    const result = await this.service.login(this.state.email);
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
    const settings = this.session.get<AuthSettings>('auth');

    if (this.state.sent) {
      return <div>
                <p>We sent an e-mail to <b>{ this.state.email }</b> with a login link.</p>
             </div>;
    }

    const google = settings.providers.google &&
                    <div className="column">
                        <SocialLogInButton provider="google" />
                    </div>;
    const facebook = settings.providers.facebook &&
                    <div className="column">
                        <SocialLogInButton provider="facebook" />
                    </div>;
    const github = settings.providers.github &&
                    <div className="column">
                        <SocialLogInButton provider="github" />
                    </div>;

    return  <div className="login-options">
                <div className="ui stackable three column centered grid">
                  { facebook }
                  { google }
                  { github }
                </div>
                <p className="info">We'll never post to any of your accounts</p>
                <div className="ui horizontal divider">OR</div>
                <p>Enter your e-mail address to log in</p>
                <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.login() }>
                  <div className="ui small action fluid input">
                      <input onChange={(e) => this.setState({ email: e.currentTarget.value }) } type="text" placeholder="yourname@example.com" className="small" />
                      <button onClick={ () => this.form.submit() } className={`ui small positive button ${this.state.email === '' && 'disabled'}`}>
                        Log in
                      </button>
                  </div>
                </Form>
            </div>;
  }
}
