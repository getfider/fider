import "./TagForm.scss";

import * as React from "react";
import { OAuthConfig } from "@fider/models";
import { Failure, Fider, actions, navigator } from "@fider/services";
import { Form, Button, Input, Heading } from "@fider/components";

interface OAuthFormProps {
  config?: OAuthConfig;
  onCancel: () => void;
}

export interface OAuthFormState {
  provider: string;
  displayName: string;
  clientId: string;
  clientSecret: string;
  clientSecretEnabled: boolean;
  authorizeUrl: string;
  tokenUrl: string;
  profileUrl: string;
  scope: string;
  jsonUserIdPath: string;
  jsonUserNamePath: string;
  jsonUserEmailPath: string;
  error?: Failure;
}

export class OAuthForm extends React.Component<OAuthFormProps, OAuthFormState> {
  constructor(props: OAuthFormProps) {
    super(props);
    this.state = {
      provider: this.props.config ? this.props.config.provider : "",
      displayName: this.props.config ? this.props.config.displayName : "",
      clientId: this.props.config ? this.props.config.clientId : "",
      clientSecret: this.props.config ? this.props.config.clientSecret : "",
      clientSecretEnabled: !this.props.config,
      authorizeUrl: this.props.config ? this.props.config.authorizeUrl : "",
      tokenUrl: this.props.config ? this.props.config.tokenUrl : "",
      profileUrl: this.props.config ? this.props.config.profileUrl : "",
      scope: this.props.config ? this.props.config.scope : "",
      jsonUserIdPath: this.props.config ? this.props.config.jsonUserIdPath : "",
      jsonUserNamePath: this.props.config ? this.props.config.jsonUserNamePath : "",
      jsonUserEmailPath: this.props.config ? this.props.config.jsonUserEmailPath : ""
    };
  }

  private handleSave = async () => {
    const result = await actions.saveOAuthConfig({
      provider: this.state.provider,
      displayName: this.state.displayName,
      clientId: this.state.clientId,
      clientSecret: this.state.clientSecretEnabled ? this.state.clientSecret : "",
      authorizeUrl: this.state.authorizeUrl,
      tokenUrl: this.state.tokenUrl,
      profileUrl: this.state.profileUrl,
      scope: this.state.scope,
      jsonUserIdPath: this.state.jsonUserIdPath,
      jsonUserNamePath: this.state.jsonUserNamePath,
      jsonUserEmailPath: this.state.jsonUserEmailPath,
    });
    if (result.ok) {
      navigator.goTo("/admin/authentication");
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

  private setClientId = (clientId: string) => {
    this.setState({ clientId });
  };

  private setClientSecret = (clientSecret: string) => {
    this.setState({ clientSecret });
  };

  private setAuthorizeUrl = (authorizeUrl: string) => {
    this.setState({ authorizeUrl });
  };

  private setTokenUrl = (tokenUrl: string) => {
    this.setState({ tokenUrl });
  };

  private setProfileUrl = (profileUrl: string) => {
    this.setState({ profileUrl });
  };

  private setScope = (scope: string) => {
    this.setState({ scope });
  };

  private setJSONUserIdPath = (jsonUserIdPath: string) => {
    this.setState({ jsonUserIdPath });
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
          <Input
            field="displayName"
            label="Display Name"
            maxLength={50}
            value={this.state.displayName}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setDisplayName}
          />
          <Input
            field="clientId"
            label="Client ID"
            maxLength={100}
            value={this.state.clientId}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setClientId}
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
            field="authorizeUrl"
            label="Authorize URL"
            maxLength={300}
            value={this.state.authorizeUrl}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setAuthorizeUrl}
          />
          <Input
            field="tokenUrl"
            label="Token URL"
            maxLength={300}
            value={this.state.tokenUrl}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setTokenUrl}
          />

          <h3>User Profile</h3>
          <p className="info">
            This section is used to configure how Fider will fetch user information like Id, Name and Email after the
            authentication OAuth process.
          </p>

          <Input
            field="profileUrl"
            label="Profile API URL"
            maxLength={300}
            value={this.state.profileUrl}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setProfileUrl}
          >
            <p className="info">
              This URL is used to fetch the authenticated user details. It must return a JSON and not require any
              QueryString parameter. E.g: Google Profile URL is https://www.googleapis.com/plus/v1/people/me
            </p>
          </Input>

          <Input
            field="scope"
            label="Scope"
            maxLength={100}
            value={this.state.scope}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setScope}
          >
            <p className="info">
              It is recommended to only request the minimum scopes we need to fecth the user <strong>id</strong>,
              <strong>name</strong> and <strong>email</strong>. Multiple scopes must be separated by space.
            </p>
          </Input>

          <h4>JSON Path</h4>

          <div className="row">
            <Input
              field="jsonUserIdPath"
              label="ID"
              className="col-sm-4"
              maxLength={100}
              value={this.state.jsonUserIdPath}
              disabled={!Fider.session.user.isAdministrator}
              onChange={this.setJSONUserIdPath}
            >
              <p className="info">
                Path to extract User ID from JSON. This ID <strong>must</strong> be unique within the provider or
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
                Path to extract user Display Name from JSON. This is optional, but <strong>highly</strong> recommended.
                For the example below, the path would be <strong>profile.name</strong>.
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
                Path to extract user Email from JSON. This is optional, but <strong>highly</strong> recommended. For the
                example below, the path would be <strong>profile.emails[0]</strong>.
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

          <div className="c-form-field">
            <Button color="positive" onClick={this.handleSave}>
              Save
            </Button>
            <Button onClick={this.handleCancel}>Cancel</Button>
          </div>
        </Form>
      </>
    );
  }
}
