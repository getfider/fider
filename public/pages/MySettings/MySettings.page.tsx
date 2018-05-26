import "./MySettings.page.scss";

import * as React from "react";

import { Modal, Form, DisplayError, Button, Gravatar, Heading, Field, Input } from "@fider/components";
import { DangerZone, NotificationSettings } from "./";

import { CurrentUser, UserSettings } from "@fider/models";
import { Failure, actions } from "@fider/services";

interface MySettingsPageState {
  showModal: boolean;
  name: string;
  newEmail: string;
  changingEmail: boolean;
  error?: Failure;
  settings: UserSettings;
}

interface MySettingsPageProps {
  user: CurrentUser;
  settings: UserSettings;
}

export class MySettingsPage extends React.Component<MySettingsPageProps, MySettingsPageState> {
  constructor(props: MySettingsPageProps) {
    super(props);
    this.state = {
      showModal: false,
      changingEmail: false,
      newEmail: "",
      name: this.props.user.name,
      settings: this.props.settings
    };
  }

  private confirm = async () => {
    const result = await actions.updateUserSettings(this.state.name, this.state.settings);
    if (result.ok) {
      location.reload();
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  };

  private submitNewEmail = async () => {
    const result = await actions.changeUserEmail(this.state.newEmail);
    if (result.ok) {
      this.setState({
        error: undefined,
        changingEmail: false,
        showModal: true
      });
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  };

  public render() {
    const changeEmail = (
      <span className="ui info clickable" onClick={() => this.setState({ changingEmail: true })}>
        change
      </span>
    );

    return (
      <div id="p-my-settings" className="page container">
        <Modal.Window isOpen={this.state.showModal} canClose={true} center={true}>
          <Modal.Header>Confirm your new email</Modal.Header>
          <Modal.Content>
            <div>
              <p>
                We have just sent a confirmation link to <b>{this.state.newEmail}</b>. <br /> Click the link to update
                your email.
              </p>
              <p>
                <a href="#" onClick={() => this.setState({ showModal: false })}>
                  OK
                </a>
              </p>
            </div>
          </Modal.Content>
        </Modal.Window>

        <Heading title="Settings" subtitle="Manage your profile settings" icon="id badge" />

        <div className="row">
          <div className="col-lg-7">
            <Form error={this.state.error}>
              <Field label="Avatar">
                <p>
                  <Gravatar user={this.props.user} />
                </p>
                <div className="info">
                  <p>
                    This site uses{" "}
                    <a href="https://en.gravatar.com/" target="blank">
                      Gravatar
                    </a>{" "}
                    to display profile avatars. <br />
                    A letter avatar based on your name is generated for profiles without a Gravatar.
                  </p>
                </div>
              </Field>

              <Input
                label="Email"
                field="email"
                value={this.state.changingEmail ? this.state.newEmail : this.props.user.email}
                maxLength={200}
                disabled={!this.state.changingEmail}
                afterLabel={this.state.changingEmail ? undefined : changeEmail}
                onChange={newEmail => this.setState({ newEmail })}
              >
                <p className="info">
                  {this.props.user.email || this.state.changingEmail
                    ? "Your email is private and will never be displayed to anyone"
                    : "Your account doesn't have an email."}
                </p>
                {this.state.changingEmail && (
                  <>
                    <Button color="positive" size="mini" onClick={this.submitNewEmail}>
                      Confirm
                    </Button>
                    <Button
                      size="mini"
                      onClick={async () =>
                        this.setState({
                          changingEmail: false,
                          newEmail: "",
                          error: undefined
                        })
                      }
                    >
                      Cancel
                    </Button>
                  </>
                )}
              </Input>

              <Input
                label="Name"
                field="name"
                value={this.state.name}
                maxLength={100}
                onChange={name => this.setState({ name })}
              />

              <NotificationSettings
                user={this.props.user}
                settings={this.props.settings}
                settingsChanged={settings => this.setState({ settings })}
              />

              <Button color="positive" onClick={this.confirm}>
                Save changes
              </Button>
            </Form>
          </div>
        </div>

        <div className="row">
          <div className="col-lg-7">
            <DangerZone user={this.props.user} />
          </div>
        </div>
      </div>
    );
  }
}
