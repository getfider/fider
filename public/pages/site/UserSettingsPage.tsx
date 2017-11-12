import * as React from 'react';

import { inject, injectables } from '@fider/di';
import { Footer, Header, Form, DisplayError, Button, Gravatar } from '@fider/components/common';

import { CurrentUser } from '@fider/models';
import { Session, Failure, UserService } from '@fider/services';

interface UserSettingsPageState {
  name: string;
  error?: Failure;
}

export class UserSettingsPage extends React.Component<{}, UserSettingsPageState> {
  private user: CurrentUser;

  @inject(injectables.Session)
  private session: Session;

  @inject(injectables.UserService)
  private userService: UserService;

  constructor(props: {}) {
      super(props);
      this.user = this.session.getCurrentUser()!;
      this.state = {
        name: this.user.name
      };
  }

  private async confirm() {
    const result = await this.userService.updateSettings(this.state.name);
    if (result.ok) {
        location.reload();
    } else if (result.error) {
        this.setState({ error: result.error });
    }
  }

  public render() {
    return (
      <div>
        <Header />
          <div className="page ui container">
            <h2 className="ui header">
              <i className="circular id badge icon" />
              <div className="content">
                Settings
                <div className="sub header">Manage your profile settings</div>
              </div>
            </h2>

            <div className="ui grid">
              <div className="eight wide computer sixteen wide mobile column">
                <div className="ui form">
                  <div className="field">
                      <label htmlFor="email">Avatar</label>
                      <p><Gravatar user={this.user} /></p>
                      <div className="info">
                          <p>
                            This site uses <a href="https://en.gravatar.com/" target="blank">Gravatar</a> to display profile avatars. <br/>
                            A letter avatar based on your name is generated for profiles without a Gravatar.
                          </p>
                      </div>
                  </div>
                  <div className="field">
                      <label htmlFor="email">E-mail</label>
                      <p><b>{this.user.email}</b></p>
                      <div className="info">
                        {
                          this.user.email
                          ? <p>Your e-mail is private and will never be displayed to anyone.</p>
                          : <p>Your account doesn't have an e-mail.</p>
                        }
                      </div>
                  </div>
                  <DisplayError fields={['name']} error={this.state.error} />
                  <div className="field">
                    <label htmlFor="name">Name</label>
                    <input
                      id="name"
                      type="text"
                      maxLength={100}
                      value={this.state.name}
                      onChange={(e) => this.setState({ name: e.currentTarget.value })}
                    />
                  </div>
                  <div className="field">
                    <Button className="positive" size="tiny" onClick={async () => await this.confirm()}>Confirm</Button>
                  </div>
                </div>
              </div>
            </div>

          </div>
        <Footer />
      </div>
    );
  }
}
