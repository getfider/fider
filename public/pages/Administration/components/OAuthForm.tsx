import React from "react";
import { OAuthConfig, OAuthConfigStatus } from "@fider/models";
import { Failure, Fider, actions } from "@fider/services";
import {
  Form,
  Button,
  Input,
  Heading,
  SocialSignInButton,
  Field,
  ImageUploadState,
  ImageUploader,
  Toggle,
  OAuthProviderLogoURL
} from "@fider/components";

interface OAuthFormProps {
  config?: OAuthConfig;
  onCancel: () => void;
}

export interface OAuthFormState {
  provider: string;
  displayName: string;
  enabled: boolean;
  clientID: string;
  clientSecret: string;
  clientSecretEnabled: boolean;
  authorizeURL: string;
  tokenURL: string;
  profileURL: string;
  scope: string;
  jsonUserIDPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
  logoURL?: string;
  logo?: ImageUploadState;
  error?: Failure;
}

export class OAuthForm extends React.Component<OAuthFormProps, OAuthFormState> {
  constructor(props: OAuthFormProps) {
    super(props);
    this.state = {
      provider: this.props.config ? this.props.config.provider : "",
      enabled: this.props.config ? this.props.config.status === OAuthConfigStatus.Enabled : false,
      displayName: this.props.config ? this.props.config.displayName : "",
      clientID: this.props.config ? this.props.config.clientID : "",
      clientSecret: this.props.config ? this.props.config.clientSecret : "",
      clientSecretEnabled: !this.props.config,
      authorizeURL: this.props.config ? this.props.config.authorizeURL : "",
      tokenURL: this.props.config ? this.props.config.tokenURL : "",
      profileURL: this.props.config ? this.props.config.profileURL : "",
      scope: this.props.config ? this.props.config.scope : "",
      jsonUserIDPath: this.props.config ? this.props.config.jsonUserIDPath : "",
      jsonUserNamePath: this.props.config ? this.props.config.jsonUserNamePath : "",
      jsonUserEmailPath: this.props.config ? this.props.config.jsonUserEmailPath : "",
      logoURL: this.props.config ? OAuthProviderLogoURL(this.props.config.logoID) : ""
    };
  }

  private handleSave = async () => {
    const result = await actions.saveOAuthConfig({
      provider: this.state.provider,
      status: this.state.enabled ? OAuthConfigStatus.Enabled : OAuthConfigStatus.Disabled,
      displayName: this.state.displayName,
      clientID: this.state.clientID,
      clientSecret: this.state.clientSecretEnabled ? this.state.clientSecret : "",
      authorizeURL: this.state.authorizeURL,
      tokenURL: this.state.tokenURL,
      profileURL: this.state.profileURL,
      scope: this.state.scope,
      jsonUserIDPath: this.state.jsonUserIDPath,
      jsonUserNamePath: this.state.jsonUserNamePath,
      jsonUserEmailPath: this.state.jsonUserEmailPath,
      logo: this.state.logo
    });
    if (result.ok) {
      location.reload();
    } else {
      this.setState({ error: result.error });
    }
  };

  private handleCancel = async () => {
    this.props.onCancel();
  };

  private setDisplayName = (displayName: string) => {
    this.setState({ displayName });
  };

  private setLogo = (logo: ImageUploadState, previewURL: string) => {
    this.setState({ logo, logoURL: previewURL });
  };

  private setStatus = async (enabled: boolean) => {
    this.setState({ enabled });
  };

  private setClientID = (clientID: string) => {
    this.setState({ clientID });
  };

  private setClientSecret = (clientSecret: string) => {
    this.setState({ clientSecret });
  };

  private setAuthorizeURL = (authorizeURL: string) => {
    this.setState({ authorizeURL });
  };

  private setTokenURL = (tokenURL: string) => {
    this.setState({ tokenURL });
  };

  private setProfileURL = (profileURL: string) => {
    this.setState({ profileURL });
  };

  private setScope = (scope: string) => {
    this.setState({ scope });
  };

  private setJSONUserIDPath = (jsonUserIDPath: string) => {
    this.setState({ jsonUserIDPath });
  };

  private setJSONUserNamePath = (jsonUserNamePath: string) => {
    this.setState({ jsonUserNamePath });
  };

  private setJSONUserEmailPath = (jsonUserEmailPath: string) => {
    this.setState({ jsonUserEmailPath });
  };

  private enableClientSecret = () => {
    this.setState({ clientSecretEnabled: true, clientSecret: "" });
  };

  public render() {
    const title = this.props.config ? `OAuth Provider: ${this.props.config.displayName}` : "New OAuth Provider";
    return (
      <>
        <Heading title={title} size="small" />
        <Form error={this.state.error}>
          <div className="row">
            <div className="col-sm-9">
              <Input
                field="displayName"
                label="Display Name"
                maxLength={50}
                value={this.state.displayName}
                disabled={!Fider.session.user.isAdministrator}
                onChange={this.setDisplayName}
              />

              <ImageUploader
                label="Logo"
                field="logo"
                defaultImageURL={this.state.logoURL}
                previewMaxWidth={80}
                disabled={!Fider.session.user.isAdministrator}
                onChange={this.setLogo}
              >
                <p className="info">
                  We accept JPG, GIF and PNG images, smaller than 50KB and with an aspect ratio of 1:1 with minimum
                  dimensions of 24x24 pixels.
                </p>
              </ImageUploader>
            </div>
            <div className="col-sm-3">
              <Field label="Button Preview">
                <SocialSignInButton option={this.state} />
              </Field>
            </div>
          </div>

          <Input
            field="clientID"
            label="Client ID"
            maxLength={100}
            value={this.state.clientID}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setClientID}
          />

          <Input
            field="clientSecret"
            label="Client Secret"
            maxLength={500}
            value={this.state.clientSecret}
            disabled={!this.state.clientSecretEnabled}
            onChange={this.setClientSecret}
            afterLabel={
              !this.state.clientSecretEnabled ? (
                <>
                  <span className="info">omitted for security reasons.</span>
                  <span className="info clickable" onClick={this.enableClientSecret}>
                    change
                  </span>
                </>
              ) : (
                undefined
              )
            }
          />
          <Input
            field="authorizeURL"
            label="Authorize URL"
            maxLength={300}
            value={this.state.authorizeURL}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setAuthorizeURL}
          />
          <Input
            field="tokenURL"
            label="Token URL"
            maxLength={300}
            value={this.state.tokenURL}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setTokenURL}
          />

          <Input
            field="scope"
            label="Scope"
            maxLength={100}
            value={this.state.scope}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setScope}
          >
            <p className="info">
              It is recommended to only request the minimum scopes we need to fetch the user <strong>id</strong>,{" "}
              <strong>name</strong> and <strong>email</strong>. Multiple scopes must be separated by space.
            </p>
          </Input>

          <h3>User Profile</h3>
          <p className="info">
            This section is used to configure how Fider will fetch user after the authentication process.
          </p>

          <Input
            field="profileURL"
            label="Profile API URL"
            maxLength={300}
            value={this.state.profileURL}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setProfileURL}
          >
            <p className="info">
              The URL to fetch the authenticated user info. If empty, Fider will try to parse the user info from the
              Access Token.
            </p>
          </Input>

          <h4>JSON Path</h4>

          <div className="row">
            <Input
              field="jsonUserIDPath"
              label="ID"
              className="col-sm-4"
              maxLength={100}
              value={this.state.jsonUserIDPath}
              disabled={!Fider.session.user.isAdministrator}
              onChange={this.setJSONUserIDPath}
            >
              <p className="info">
                Path to extract User ID from the JSON. This ID <strong>must</strong> be unique within the provider or
                unexpected side effects might happen. For example below, the path would be <strong>id</strong>.
              </p>
            </Input>
            <Input
              field="jsonUserNamePath"
              label="Name"
              className="col-sm-4"
              maxLength={100}
              value={this.state.jsonUserNamePath}
              disabled={!Fider.session.user.isAdministrator}
              onChange={this.setJSONUserNamePath}
            >
              <p className="info">
                Path to extract user Display Name from the JSON. This is optional, but <strong>highly</strong>{" "}
                recommended. For the example below, the path would be <strong>profile.name</strong>.
              </p>
            </Input>
            <Input
              field="jsonUserEmailPath"
              label="Email"
              className="col-sm-4"
              maxLength={100}
              value={this.state.jsonUserEmailPath}
              disabled={!Fider.session.user.isAdministrator}
              onChange={this.setJSONUserEmailPath}
            >
              <p className="info">
                Path to extract user Email from the JSON. This is optional, but <strong>highly</strong> recommended. For
                the example below, the path would be <strong>profile.emails[0]</strong>.
              </p>
            </Input>
          </div>
          <pre>
            <h5>Example Response</h5>
            {`
{ 
  id: "35235"
  title: "Sr. Account Manager",
  profile: {
    dob: "01/05/2018",
    name: "John Doe"
    emails: [
      "john.doe@company.com"
    ]
  }
}
            `}
          </pre>

          <div className="row">
            <div className="col-sm-4">
              <Field label="Status">
                <Toggle active={this.state.enabled} onToggle={this.setStatus} />
                <span>{this.state.enabled ? "Enabled" : "Disabled"}</span>
                {this.state.enabled && (
                  <p className="info">
                    This provider will be available for everyone to use during the sign in process. It is recommended
                    that you keep it disable and test it before enabling it. The Test button is available after saving
                    this configuration.
                  </p>
                )}
                {!this.state.enabled && <p className="info">Users won't be able to sign in with this Provider.</p>}
              </Field>
            </div>
          </div>

          <div className="c-form-field">
            <Button color="positive" onClick={this.handleSave}>
              Save
            </Button>
            <Button color="cancel" onClick={this.handleCancel}>
              Cancel
            </Button>
          </div>
        </Form>
      </>
    );
  }
}
