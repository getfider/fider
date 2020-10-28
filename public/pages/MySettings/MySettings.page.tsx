import "./MySettings.page.scss";

import React from "react";

import { Modal, Form, Button, Heading, Input, Select, SelectOption, ImageUploader } from "@fider/components";

import { UserSettings, UserAvatarType, ImageUpload } from "@fider/models";
import { Failure, actions, Fider } from "@fider/services";
import { FaRegAddressCard } from "react-icons/fa";
import { withTranslation, WithTranslation } from "react-i18next";
import { NotificationSettings } from "./components/NotificationSettings";
import { APIKeyForm } from "./components/APIKeyForm";
import { DangerZone } from "./components/DangerZone";

interface MySettingsPageState {
  showModal: boolean;
  name: string;
  newEmail: string;
  avatar?: ImageUpload;
  avatarType: UserAvatarType;
  changingEmail: boolean;
  error?: Failure;
  userSettings: UserSettings;
}

interface MySettingsPageProps extends WithTranslation {
  userSettings: UserSettings;
}

class MySettingsPage extends React.Component<MySettingsPageProps, MySettingsPageState> {
  constructor(props: MySettingsPageProps) {
    super(props);
    this.state = {
      showModal: false,
      changingEmail: false,
      avatarType: Fider.session.user.avatarType,
      newEmail: "",
      name: Fider.session.user.name,
      userSettings: this.props.userSettings
    };
  }

  private confirm = async () => {
    const result = await actions.updateUserSettings({
      name: this.state.name,
      avatarType: this.state.avatarType,
      avatar: this.state.avatar,
      settings: this.state.userSettings
    });
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

  private avatarTypeChanged = (opt?: SelectOption) => {
    if (opt) {
      this.setState({ avatarType: opt.value as UserAvatarType });
    }
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

  private setAvatar = (avatar: ImageUpload): void => {
    this.setState({ avatar });
  };

  public render() {
    const { t } = this.props;
    const changeEmail = (
      <span className="ui info clickable" onClick={this.startChangeEmail}>
        {t("mySettings.change")}
      </span>
    );

    return (
      <div id="p-my-settings" className="page container">
        <Modal.Window isOpen={this.state.showModal} onClose={this.closeModal}>
          <Modal.Header> {t("mySettings.confirmEmail")} </Modal.Header>
          <Modal.Content>
            <div>
              <p dangerouslySetInnerHTML={{ __html: t("mySettings.confirmMessage", { email: this.state.newEmail }) }} />
              <p>
                <a href="#" onClick={this.closeModal}>
                  OK
                </a>
              </p>
            </div>
          </Modal.Content>
        </Modal.Window>

        <Heading title={t("mySettings.title")} subtitle={t("mySettings.subtitle")} icon={FaRegAddressCard} />

        <div className="row">
          <div className="col-lg-7">
            <Form error={this.state.error}>
              <Input
                label={t("email")}
                field="email"
                value={this.state.changingEmail ? this.state.newEmail : Fider.session.user.email}
                maxLength={200}
                disabled={!this.state.changingEmail}
                afterLabel={this.state.changingEmail ? undefined : changeEmail}
                onChange={this.setNewEmail}
              >
                <p className="info">
                  {Fider.session.user.email || this.state.changingEmail
                    ? t("mySettings.privateMessage")
                    : t("mySettings.dontHaveEmail")}
                </p>
                {this.state.changingEmail && (
                  <>
                    <Button color="positive" size="mini" onClick={this.submitNewEmail}>
                      {t("common.button.confirm")}
                    </Button>
                    <Button color="cancel" size="mini" onClick={this.cancelChangeEmail}>
                      {t("common.button.cancel")}
                    </Button>
                  </>
                )}
              </Input>

              <Input label="Name" field="name" value={this.state.name} maxLength={100} onChange={this.setName} />

              <Select
                label={t("avatar")}
                field="avatarType"
                defaultValue={this.state.avatarType}
                options={[
                  { label: "Letter", value: UserAvatarType.Letter },
                  { label: "Gravatar", value: UserAvatarType.Gravatar },
                  { label: "Custom", value: UserAvatarType.Custom }
                ]}
                onChange={this.avatarTypeChanged}
              >
                {this.state.avatarType === UserAvatarType.Gravatar && (
                  <p className="info">
                    A{" "}
                    <a href="https://en.gravatar.com" target="_blank">
                      Gravatar
                    </a>{" "}
                    {t("mySettings.gravatarMessage")}
                  </p>
                )}
                {this.state.avatarType === UserAvatarType.Letter && (
                  <p className="info">{t("mySettings.letterMessage")}</p>
                )}
                {this.state.avatarType === UserAvatarType.Custom && (
                  <ImageUploader
                    field="avatar"
                    previewMaxWidth={80}
                    onChange={this.setAvatar}
                    bkey={Fider.session.user.avatarBlobKey}
                  >
                    <p className="info">{t("mySettings.imageFormatMessage")}</p>
                  </ImageUploader>
                )}
              </Select>

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

export default withTranslation()(MySettingsPage);
