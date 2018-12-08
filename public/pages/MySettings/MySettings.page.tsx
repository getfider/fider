import "./MySettings.page.scss";

import React from "react";

import { Modal, Form, Button, Gravatar, Heading, Field, Input } from "@fider/components";

import { UserSettings } from "@fider/models";
import { Failure, actions, Fider } from "@fider/services";
import { FaRegAddressCard } from "react-icons/fa";
import { NotificationSettings } from "./components/NotificationSettings";
import { APIKeyForm } from "./components/APIKeyForm";
import { DangerZone } from "./components/DangerZone";

interface MySettingsPageState {
  showModal: boolean;
  name: string;
  newEmail: string;
  changingEmail: boolean;
  error?: Failure;
  userSettings: UserSettings;
}

interface MySettingsPageProps {
  userSettings: UserSettings;
}

export default class MySettingsPage extends React.Component<MySettingsPageProps, MySettingsPageState> {
  constructor(props: MySettingsPageProps) {
    super(props);
    this.state = {
      showModal: false,
      changingEmail: false,
      newEmail: "",
      name: Fider.session.user.name,
      userSettings: this.props.userSettings
    };
  }

  private confirm = async () => {
    const result = await actions.updateUserSettings(this.state.name, this.state.userSettings);
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

  private startChangeEmail = () => {
    this.setState({ changingEmail: true });
  };

  private cancelChangeEmail = async () => {
    this.setState({
      changingEmail: false,
      newEmail: "",
      error: undefined
    });
  };

  private setName = (name: string) => {
    this.setState({ name });
  };

  private setNotificationSettings = (userSettings: UserSettings) => {
    this.setState({ userSettings });
  };

  private closeModal = () => {
    this.setState({ showModal: false });
  };

  private setNewEmail = (newEmail: string) => {
    this.setState({ newEmail });
  };

  public render() {
    const changeEmail = (
      <span className="ui info clickable" onClick={this.startChangeEmail}>
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
                <a href="#" onClick={this.closeModal}>
                  OK
                </a>
              </p>
            </div>
          </Modal.Content>
        </Modal.Window>

        <Heading title="Settings" subtitle="Manage your profile settings" icon={FaRegAddressCard} />

        <div className="row">
          <div className="col-lg-7">
            <Form error={this.state.error}>
              <Field label="Avatar">
                <p>
                  <Gravatar user={Fider.session.user} />
                </p>
                <div className="info">
                  <p>
                    This site uses{" "}
                    <a href="https://en.gravatar.com/" target="blank">
                      Gravatar
                    </a>{" "}
                    to display profile avatars. <br />A letter avatar based on your name is generated for profiles
                    without a Gravatar.
                  </p>
                </div>
              </Field>

              <Input
                label="Email"
                field="email"
                value={this.state.changingEmail ? this.state.newEmail : Fider.session.user.email}
                maxLength={200}
                disabled={!this.state.changingEmail}
                afterLabel={this.state.changingEmail ? undefined : changeEmail}
                onChange={this.setNewEmail}
              >
                <p className="info">
                  {Fider.session.user.email || this.state.changingEmail
                    ? "Your email is private and will never be publicly displayed"
                    : "Your account doesn't have an email."}
                </p>
                {this.state.changingEmail && (
                  <>
                    <Button color="positive" size="mini" onClick={this.submitNewEmail}>
                      Confirm
                    </Button>
                    <Button color="cancel" size="mini" onClick={this.cancelChangeEmail}>
                      Cancel
                    </Button>
                  </>
                )}
              </Input>

              <Input label="Name" field="name" value={this.state.name} maxLength={100} onChange={this.setName} />

              <NotificationSettings
                userSettings={this.props.userSettings}
                settingsChanged={this.setNotificationSettings}
              />

              <Button color="positive" onClick={this.confirm}>
                Save
              </Button>
            </Form>
          </div>
        </div>

        {Fider.session.user.isCollaborator && (
          <div className="row">
            <div className="col-lg-7">
              <APIKeyForm />
            </div>
          </div>
        )}

        <div className="row">
          <div className="col-lg-7">
            <DangerZone />
          </div>
        </div>
      </div>
    );
  }
}
