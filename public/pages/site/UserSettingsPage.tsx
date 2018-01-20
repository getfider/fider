import * as React from 'react';

import { Form, DisplayError, Button, Gravatar } from '@fider/components/common';

import { CurrentUser } from '@fider/models';
import { Failure, actions } from '@fider/services';

import './UserSettingsPage.scss';

interface UserSettingsPageState {
  name: string;
  newEmail: string;
  changingEmail: boolean;
  error?: Failure;
}

interface UserSettingsPageProps {
  user: CurrentUser;
}

export class UserSettingsPage extends React.Component<UserSettingsPageProps, UserSettingsPageState> {

  constructor(props: UserSettingsPageProps) {
    super(props);
    this.state = {
      changingEmail: false,
      newEmail: '',
      name: this.props.user.name
    };
  }

  private async confirm() {
    const result = await actions.updateUserSettings(this.state.name);
    if (result.ok) {
      location.reload();
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  private async submitNewEmail() {
    const result = await actions.changeUserEmail(this.state.newEmail);
    if (result.ok) {
      this.setState({ error: undefined, changingEmail: false }, () => {
        $('#confirmation-modal').modal('show');
      });
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  public render() {
    return (
      <div className="page ui container">
        <div id="confirmation-modal" className="ui mini modal">
          <div className="header">Confirm your new e-mail</div>
          <div className="content">
            <div>
              <p>We have just sent a confirmation link to <b>{this.state.newEmail}</b>. <br /> Click the link to update your e-mail.</p>
              <p><a href="#" onClick={() => $('#confirmation-modal').modal('hide')}>OK</a></p>
            </div>
          </div>
        </div>

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
                  <p><Gravatar user={this.props.user} /></p>
                  <div className="info">
                      <p>
                        This site uses <a href="https://en.gravatar.com/" target="blank">Gravatar</a> to display profile avatars. <br/>
                        A letter avatar based on your name is generated for profiles without a Gravatar.
                      </p>
                  </div>
              </div>
              <DisplayError fields={['email']} error={this.state.error} />
              <div className="field">
                  <label htmlFor="email">E-mail <span className="info">Your e-mail is private and will never be displayed to anyone.</span></label>
                  {
                    this.state.changingEmail ?
                    <>
                      <p>
                        <input
                          id="new-email"
                          type="text"
                          style={{'max-width': '200px', 'margin-right': '10px'}}
                          maxLength={200}
                          placeholder={this.props.user.email}
                          value={this.state.newEmail}
                          onChange={(e) => this.setState({ newEmail: e.currentTarget.value })}
                        />
                        <Button className="positive" size="mini" onClick={async () => await this.submitNewEmail()}>Confirm</Button>
                        <Button size="mini" onClick={async () => this.setState({ changingEmail: false, error: undefined })}>Cancel</Button>
                      </p>
                    </>
                    :
                    <p>
                      {
                        this.props.user.email
                        ? <b>{this.props.user.email}</b>
                        : <span className="info">Your account doesn't have an e-mail.</span>
                      }
                      <span className="ui info clickable" onClick={() => this.setState({ changingEmail: true })}>change</span>
                    </p>
                  }
                  
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
    );
  }
}
