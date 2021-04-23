import React from "react"

import { Button, OAuthProviderLogo, Icon } from "@fider/components"
import { OAuthConfig, OAuthProviderOption } from "@fider/models"
import { OAuthForm } from "../components/OAuthForm"
import { actions, notify, Fider } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

import IconPlay from "@fider/assets/images/heroicons-play.svg"
import IconPencilAlt from "@fider/assets/images/heroicons-pencil-alt.svg"

import { HStack, VStack } from "@fider/components/layout"

interface ManageAuthenticationPageProps {
  providers: OAuthProviderOption[]
}

interface ManageAuthenticationPageState {
  isAdding: boolean
  editing?: OAuthConfig
}

export default class ManageAuthenticationPage extends AdminBasePage<ManageAuthenticationPageProps, ManageAuthenticationPageState> {
  public id = "p-admin-authentication"
  public name = "authentication"
  public title = "Authentication"
  public subtitle = "Manage your site authentication"

  constructor(props: ManageAuthenticationPageProps) {
    super(props)
    this.state = {
      isAdding: false,
    }
  }

  private addNew = async () => {
    this.setState({ isAdding: true, editing: undefined })
  }

  private edit = async (provider: string) => {
    const result = await actions.getOAuthConfig(provider)
    if (result.ok) {
      this.setState({ editing: result.data, isAdding: false })
    } else {
      notify.error("Failed to retrieve OAuth configuration. Try again later")
    }
  }

  private startTest = async (provider: string) => {
    const redirect = `${Fider.settings.baseURL}/oauth/${provider}/echo`
    window.open(`/oauth/${provider}?redirect=${redirect}`, "oauth-test", "width=1100,height=600,status=no,menubar=no")
  }

  private cancel = async () => {
    this.setState({ isAdding: false, editing: undefined })
  }

  public content() {
    if (this.state.isAdding) {
      return <OAuthForm onCancel={this.cancel} />
    }

    if (this.state.editing) {
      return <OAuthForm config={this.state.editing} onCancel={this.cancel} />
    }

    const enabled = <span className="text-green-700">Enabled</span>
    const disabled = <span className="text-red-700">Disabled</span>

    return (
      <>
        <h2 className="text-display">OAuth Providers</h2>
        <p>
          You can use these section to add any authentication provider thats supports the OAuth2 protocol. Additional information is available in our{" "}
          <a rel="noopener" className="text-link" target="_blank" href="https://getfider.com/docs/configuring-oauth/">
            OAuth Documentation
          </a>
          .
        </p>
        <VStack spacing={6}>
          {this.props.providers.map((o) => (
            <div key={o.provider}>
              <HStack justify="between">
                <HStack className="h-6">
                  <OAuthProviderLogo option={o} />
                  <strong>{o.displayName}</strong>
                </HStack>
                {o.isCustomProvider && (
                  <HStack>
                    {Fider.session.user.isAdministrator && (
                      <Button onClick={this.edit.bind(this, o.provider)} size="small">
                        <Icon sprite={IconPencilAlt} />
                        <span>Edit</span>
                      </Button>
                    )}
                    <Button onClick={this.startTest.bind(this, o.provider)} size="small">
                      <Icon sprite={IconPlay} />
                      <span>Test</span>
                    </Button>
                  </HStack>
                )}
              </HStack>
              <div className="text-xs block my-1">{o.isEnabled ? enabled : disabled}</div>
              {o.isCustomProvider && (
                <span className="text-muted">
                  <strong>Client ID:</strong> {o.clientID} <br />
                  <strong>Callback URL:</strong> {o.callbackURL}
                </span>
              )}
            </div>
          ))}
          <div>
            {Fider.session.user.isAdministrator && (
              <Button variant="secondary" onClick={this.addNew}>
                Add new
              </Button>
            )}
          </div>
        </VStack>
      </>
    )
  }
}
