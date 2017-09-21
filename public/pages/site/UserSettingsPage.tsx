import * as React from 'react';

import { inject, injectables } from '@fider/di';
import { Footer, Header, Form, DisplayError, Button, Gravatar } from '@fider/components/common';

import { CurrentUser } from '@fider/models';
import { Session, Failure } from '@fider/services';

interface UserSettingsPageState {
  name: string;
  error?: Failure;
}

export class UserSettingsPage extends React.Component<{}, UserSettingsPageState> {
  private user: CurrentUser;

  @inject(injectables.Session)
  private session: Session;

  constructor(props: {}) {
      super(props);
      this.user = this.session.getCurrentUser()!;
      this.state = {
        name: this.user.name
      };
  }

  private async confirm() {
    console.log('Done!');
  }

  public render() {
    return <div>
              <Header />
                <div className="page ui container">
                  <h1 className="ui header">Settings</h1>

                  <div className="ui grid">
                    <div className="eight wide computer sixteen wide mobile column">
                      <div className="ui form">
                        <div className="field">
                            <label htmlFor="email">Avatar</label>
                            <p><Gravatar hash={ this.user.gravatar } name={ this.state.name } /></p>
                            <p className="info">
                                <p>
                                  We use <a href="https://en.gravatar.com/" target="blank">Gravatar</a> to display profile avatars. <br/>
                                  A letter avatar based on your name is generated for profiles without e-mail.
                                </p>
                            </p>
                        </div>
                        <div className="field">
                            <label htmlFor="email">E-mail</label>
                            <p><b>{ this.user.email }</b></p>
                            <p className="info">
                                {
                                  this.user.email ? <p>Your e-mail is private and will never be displayed to anyone.</p>
                                                  : <p>Your account doesn't have an e-mail.</p>
                                }
                            </p>
                        </div>
                        <DisplayError fields={['name']} error={this.state.error} />
                        <div className="field">
                            <label htmlFor="name">Name</label>
                            <input id="name"
                                   type="text"
                                   disabled
                                   maxLength={100}
                                   value={ this.state.name }
                                   onChange={(e) => this.setState({ name: e.currentTarget.value })}/>
                        </div>
                      </div>
                    </div>
                  </div>

                </div>
              <Footer />
           </div>;
  }
}
