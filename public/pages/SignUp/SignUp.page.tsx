import "./SignUp.page.scss";

import React from "react";
import { Modal, Button, Form, Input, LegalAgreement, Password } from "@fider/components";
import { actions, Failure } from "@fider/services";
import { withTranslation, WithTranslation } from "react-i18next";

interface SignUpPageState {
  submitted: boolean;
  legalAgreement: boolean;
  error?: Failure;
  name?: string;
  email?: string;
  password?: string;
}

class SignUpPage extends React.Component<WithTranslation, SignUpPageState> {
  constructor(props: WithTranslation) {
    super(props);
    this.state = {
      submitted: false,
      legalAgreement: false
    };
  }

  private confirm = async () => {
    const result = await actions.signUp({
      legalAgreement: this.state.legalAgreement,
      email: this.state.email,
      password: this.state.password,
      name: this.state.name
    });
    if (result.ok) {
      const baseURL = `${location.protocol}//${location.host}`;
      location.href = baseURL + "/signin";
    } else if (result.error) {
      this.setState({ error: result.error, submitted: false });
    }
    // const result = await actions.createTenant({
    //   token: this.user && this.user.token,
    //   legalAgreement: this.state.legalAgreement,
    //   tenantName: this.state.tenantName,
    //   subdomain: this.state.subdomain.value,
    //   name: this.state.name,
    //   email: this.state.email
    // });
    // if (result.ok) {
    //   if (this.user) {
    //     if (Fider.isSingleHostMode()) {
    //       location.reload();
    //     } else {
    //       let baseURL = `${location.protocol}//${this.state.subdomain.value}${Fider.settings.domain}`;
    //       if (location.port) {
    //         baseURL = `${baseURL}:${location.port}`;
    //       }
    //       location.href = baseURL;
    //     }
    //   } else {
    //     this.setState({ submitted: true });
    //   }
    // } else if (result.error) {
    //   this.setState({ error: result.error, submitted: false });
    // }
  };

  private onAgree = (agreed: boolean): void => {
    this.setState({ legalAgreement: agreed });
  };

  private setEmail = (email: string): void => {
    this.setState({ email });
  };

  private setPassword = (password: string): void => {
    this.setState({ password });
  };

  private setName = (name: string): void => {
    this.setState({ name });
  };

  private noop = () => {
    // do nothing
  };

  public render() {
    const { t } = this.props;
    const modal = (
      <Modal.Window canClose={false} isOpen={this.state.submitted} onClose={this.noop}>
        <Modal.Header>{t("signUp.thankYou")}</Modal.Header>
        <Modal.Content>
          <p dangerouslySetInnerHTML={{ __html: t("signUp.sentLink", { email: this.state.email }) }} />
        </Modal.Content>
      </Modal.Window>
    );

    return (
      <div id="p-signup" className="page container">
        {modal}
        <img className="logo" src="https://getfider.com/images/logo-100x100.png" />

        <p>{t("signUp.message")}</p>
        <Form error={this.state.error}>
          <Input field="name" maxLength={200} onChange={this.setName} placeholder="Name" />
          <Input field="email" maxLength={200} onChange={this.setEmail} placeholder="Email" />
          <Password field="password" maxLength={128} onChange={this.setPassword} placeholder="Password" />
        </Form>

        <Form error={this.state.error}>
          <LegalAgreement onChange={this.onAgree} />
        </Form>

        <Button className="c-button m-positive" color="positive" size="large" onClick={this.confirm}>
          {t("common.button.joinNow")}
        </Button>
        <div className="c-divider">OR</div>
        <Button color="default" size="normal" onClick={this.confirm}>
          {t("signUp.signIn")}
        </Button>
      </div>
    );
  }
}

export default withTranslation()(SignUpPage);
