import "./Start.page.scss";

import React from "react";
import {
  SignInControl,
  Modal,
  Button,
  DisplayError,
  Form,
  Input,
  Message,
  LegalAgreement,
  Password
} from "@fider/components";
import { jwt, actions, Failure, querystring, Fider } from "@fider/services";
import { withTranslation, WithTranslation } from "react-i18next";

interface OAuthUser {
  token: string;
  name: string;
  email: string;
}

interface SignUpPageState {
  submitted: boolean;
  tenantName: string;
  legalAgreement: boolean;
  error?: Failure;
  name?: string;
  email?: string;
  password?: string;
  subdomain: {
    available: boolean;
    message?: string;
    value?: string;
  };
}

class SignUpPage extends React.Component<WithTranslation, SignUpPageState> {
  private user?: OAuthUser;

  constructor(props: WithTranslation) {
    super(props);
    this.state = {
      submitted: false,
      legalAgreement: false,
      tenantName: "",
      subdomain: { available: false }
    };

    const token = querystring.get("token");
    if (token) {
      const data = jwt.decode(token);
      if (data) {
        this.user = {
          token,
          name: data["oauth/name"],
          email: data["oauth/email"]
        };
      }
    }
  }

  private confirm = async () => {
    const result = await actions.createTenant({
      token: this.user && this.user.token,
      legalAgreement: this.state.legalAgreement,
      tenantName: this.state.tenantName,
      subdomain: this.state.subdomain.value,
      name: this.state.name,
      email: this.state.email,
      password: this.state.password
    });

    if (result.ok) {
      if (this.user) {
        if (Fider.isSingleHostMode()) {
          location.reload();
        } else {
          let baseURL = `${location.protocol}//${this.state.subdomain.value}${Fider.settings.domain}`;
          if (location.port) {
            baseURL = `${baseURL}:${location.port}`;
          }

          location.href = baseURL;
        }
      } else {
        this.setState({ submitted: true });
      }
    } else if (result.error) {
      this.setState({ error: result.error, submitted: false });
    }
  };

  private timer?: number;
  private checkAvailability = (subdomain: string) => {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(() => {
      actions.checkAvailability(subdomain).then(result => {
        this.setState({
          subdomain: {
            value: subdomain,
            available: !result.data.message,
            message: result.data.message
          }
        });
      });
    }, 500);
  };

  private setSubdomain = async (subdomain: string) => {
    this.setState(
      {
        subdomain: {
          value: subdomain,
          available: false
        }
      },
      this.checkAvailability.bind(this, subdomain)
    );
  };

  private onAgree = (agreed: boolean): void => {
    this.setState({ legalAgreement: agreed });
  };

  private setName = (name: string): void => {
    this.setState({ name });
  };

  private setEmail = (email: string): void => {
    this.setState({ email });
  };

  private setTenantName = (tenantName: string): void => {
    this.setState({ tenantName });
  };

  private setPassword = (password: string): void => {
    this.setState({ password });
  };

  private noop = () => {
    location.reload();
  };

  public render() {
    const { t } = this.props;
    const modal = (
      <Modal.Window canClose={true} isOpen={this.state.submitted} onClose={this.noop}>
        <Modal.Header>{t("start.submitted")}</Modal.Header>
        <Modal.Content>
          <p dangerouslySetInnerHTML={{ __html: t("start.message", { email: this.state.email }) }} />
        </Modal.Content>
      </Modal.Window>
    );

    return (
      <div id="p-signup" className="page container">
        {modal}
        <img className="logo" src="https://getfider.com/images/logo-100x100.png" />

        <h3>{t("signUp.step1Title")}</h3>
        <DisplayError fields={["token"]} error={this.state.error} />

        {this.user ? (
          <p>
            Hello, <b>{this.user.name}</b> {this.user.email && `(${this.user.email})`}
          </p>
        ) : (
          <>
            <p>{t("signUp.step1Message")}</p>
            <SignInControl useEmail={false} />
            <Form error={this.state.error}>
              <Input field="name" maxLength={100} onChange={this.setName} placeholder="Name" />
              <Input field="email" maxLength={200} onChange={this.setEmail} placeholder="Email" />
              <Password field="password" maxLength={128} onChange={this.setPassword} placeholder="Password" />
            </Form>
          </>
        )}

        <h3>{t("signUp.step2Title")}</h3>

        <Form error={this.state.error}>
          <Input
            field="tenantName"
            maxLength={60}
            onChange={this.setTenantName}
            placeholder="your company or product name"
          />
          {!Fider.isSingleHostMode() && (
            <Input
              field="subdomain"
              maxLength={40}
              onChange={this.setSubdomain}
              placeholder="subdomain"
              suffix={Fider.settings.domain}
            >
              {this.state.subdomain.available && (
                <Message type="success" showIcon={true}>
                  This subdomain is available!
                </Message>
              )}
              {this.state.subdomain.message && (
                <Message type="error" showIcon={true}>
                  {this.state.subdomain.message}
                </Message>
              )}
            </Input>
          )}
        </Form>

        <h3>{t("signUp.step3Title")}</h3>

        <p>{t("signUp.step3Message")}</p>

        <Form error={this.state.error}>
          <LegalAgreement onChange={this.onAgree} />
        </Form>

        <Button color="positive" size="large" onClick={this.confirm}>
          {t("common.button.confirm")}
        </Button>
      </div>
    );
  }
}

export default withTranslation()(SignUpPage);
