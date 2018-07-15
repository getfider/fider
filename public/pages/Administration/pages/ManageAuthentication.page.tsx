import "./ManageAuthentication.page.scss";

import * as React from "react";

import { AdminBasePage } from "../components";
import { Segment, List, ListItem, Button, Heading } from "@fider/components";
import { OAuthConfig, OAuthProviderOption } from "@fider/models";
import { OAuthForm } from "../components/OAuthForm";
import { actions, notify } from "@fider/services";

interface ManageAuthenticationPageProps {
  providers: OAuthProviderOption[];
}

interface ManageAuthenticationPageState {
  isAdding: boolean;
  editing?: OAuthConfig;
}

export class ManageAuthenticationPage extends AdminBasePage<
  ManageAuthenticationPageProps,
  ManageAuthenticationPageState
> {
  public id = "p-admin-authentication";
  public name = "authentication";
  public icon = "sign in alternate";
  public title = "Authentication";
  public subtitle = "Manage your site authentication";

  constructor(props: ManageAuthenticationPageProps) {
    super(props);
    this.state = {
      isAdding: false
    };
  }

  private addNew = async () => {
    this.setState({ isAdding: true, editing: undefined });
  };

  private edit = async (provider: string) => {
    const result = await actions.getOAuthConfig(provider);
    if (result.ok) {
      this.setState({ editing: result.data, isAdding: false });
    } else {
      notify.error("Failed to retrieve OAuth configuration. Try again later");
    }
  };

  private cancel = async () => {
    this.setState({ isAdding: false, editing: undefined });
  };

  public content() {
    if (this.state.isAdding) {
      return <OAuthForm onCancel={this.cancel} />;
    }

    if (this.state.editing) {
      return <OAuthForm config={this.state.editing} onCancel={this.cancel} />;
    }

    return (
      <>
        <Heading
          title="OAuth Providers"
          subtitle="You can use these section to add any authentication provider as long as it implements the OAuth2 protocol."
          size="small"
        />
        <Segment>
          <List divided={true}>
            {this.props.providers.map(o => (
              <ListItem key={o.provider}>
                <strong>{o.displayName}</strong>
                <p className="info">
                  <strong>Client ID:</strong> {o.clientId} <br />
                  <strong>Callback URL:</strong> {o.callbackUrl}
                </p>
                {o.isCustomProvider && (
                  <Button key={1} onClick={this.edit.bind(this, o.provider)} className="right">
                    <i className="edit icon" />Edit
                  </Button>
                )}
              </ListItem>
            ))}
          </List>
        </Segment>
        <Button color="positive" onClick={this.addNew}>
          Add new
        </Button>
      </>
    );
  }
}
