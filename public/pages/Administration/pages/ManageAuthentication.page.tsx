import React from "react"

import { Button, OAuthProviderLogo, Icon, Toggle, Form } from "@fider/components"
import { OAuthConfig, OAuthProviderOption } from "@fider/models"
import { OAuthForm } from "../components/OAuthForm"
import { actions, notify, Fider, Failure } from "@fider/services"
import { AdminBasePage } from "../components/AdminBasePage"

import IconPlay from "@fider/assets/images/heroicons-play.svg"
import IconPencilAlt from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconAt from "@fider/assets/images/heroicons-at.svg"

import { HStack, VStack } from "@fider/components/layout"

interface ManageAuthenticationPageProps {
  providers: OAuthProviderOption[]
}

interface ManageAuthenticationPageState {
  isAdding: boolean
  isEmailAuthAllowed: boolean
  editing?: OAuthConfig
  error?: Failure
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
      isEmailAuthAllowed: Fider.session.tenant.isEmailAuthAllowed,
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

  private toggleEmailAuth = async (active: boolean) => {
    this.setState(
      () => ({
        isEmailAuthAllowed: active,
      }),
      async () => {
        const response = await actions.updateTenantEmailAuthAllowed(this.state.isEmailAuthAllowed)
        if (!response.ok) {
          this.setState(
            () => ({
              isEmailAuthAllowed: !active,
              error: response.error,
            }),
            () => notify.error("Unable to save this setting.")
          )
        }
      }
    )
  }

  private toggleSystemProvider = async (provider: OAuthProviderOption, active: boolean) => {
    const response = await actions.setSystemProviderStatus(provider.provider, active)
    if (response.ok) {
      provider.isEnabled = active
      this.forceUpdate()
    }
  }

  public content() {
    const isAdmin = Fider.session.user.isAdministrator
    const enabledProvidersCount = this.props.providers.filter((o) => o.isEnabled).length
    const canDisableEmailAuth = enabledProvidersCount > 0
    const cantDisable = !this.state.isEmailAuthAllowed && enabledProvidersCount == 1

    if (this.state.isAdding) {
      return <OAuthForm cantDisable={cantDisable} onCancel={this.cancel} />
    }

    if (this.state.editing) {
      return <OAuthForm cantDisable={cantDisable} config={this.state.editing} onCancel={this.cancel} />
    }

    const enabled = <span className="text-green-700">Enabled</span>
    const disabled = <span className="text-red-700">Disabled</span>

    const systemProviders = this.props.providers.filter((o) => !o.isCustomProvider)
    const customProviders = this.props.providers.filter((o) => o.isCustomProvider)

    return (
      <VStack spacing={8}>
        <div>
          <h2 className="text-display">General Authentication</h2>
          <p className="text-muted mb-4">Choose how people can sign in. Administrators can always sign in with their email.</p>
          <Form error={this.state.error}>
            <VStack spacing={6}>
              <div>
                <HStack justify="between">
                  <HStack className="h-6">
                    <Icon sprite={IconAt} width="24" height="24" />
                    <strong>Email address</strong>
                  </HStack>
                  <Toggle
                    field="isEmailAuthAllowed"
                    ariaLabel="Email address"
                    disabled={!isAdmin || !canDisableEmailAuth}
                    active={this.state.isEmailAuthAllowed}
                    onToggle={this.toggleEmailAuth}
                  />
                </HStack>
                {!canDisableEmailAuth && <p className="text-muted mt-1">Enable another sign in method before disabling email.</p>}
              </div>
              {systemProviders.map((o) => (
                <HStack key={o.provider} justify="between">
                  <HStack className="h-6">
                    <OAuthProviderLogo option={o} />
                    <strong>{o.displayName}</strong>
                  </HStack>
                  {o.clientID ? (
                    <Toggle
                      field={`provider-${o.provider}`}
                      ariaLabel={o.displayName}
                      disabled={!isAdmin || (cantDisable && o.isEnabled)}
                      active={o.isEnabled}
                      onToggle={this.toggleSystemProvider.bind(this, o)}
                    />
                  ) : (
                    <span className="text-muted">Not configured</span>
                  )}
                </HStack>
              ))}
            </VStack>
          </Form>
        </div>
        <div>
          <h2 className="text-display">Custom OAuth Providers</h2>
          <p className="text-muted mb-4">
            Add any authentication provider that supports the OAuth2 protocol. Additional information is available in our{" "}
            <a rel="noopener" className="text-link" target="_blank" href="https://docs.fider.io/configuring-oauth">
              OAuth Documentation
            </a>
            .
          </p>
          <VStack spacing={6}>
            {customProviders.map((o) => (
              <div key={o.provider}>
                <HStack justify="between">
                  <HStack className="h-6">
                    <OAuthProviderLogo option={o} />
                    <strong>{o.displayName}</strong>
                  </HStack>
                  {isAdmin && (
                    <HStack>
                      <Button onClick={this.edit.bind(this, o.provider)} size="small">
                        <Icon sprite={IconPencilAlt} />
                        <span>Edit</span>
                      </Button>
                      <Button onClick={this.startTest.bind(this, o.provider)} size="small">
                        <Icon sprite={IconPlay} />
                        <span>Test</span>
                      </Button>
                    </HStack>
                  )}
                </HStack>
                <div className="text-xs block my-1">{o.isEnabled ? enabled : disabled}</div>
                <span className="text-muted">
                  <strong>Client ID:</strong> {o.clientID} <br />
                  <strong>Callback URL:</strong> {o.callbackURL}
                </span>
              </div>
            ))}
            {isAdmin && (
              <div>
                <Button variant="secondary" onClick={this.addNew}>
                  Add new
                </Button>
              </div>
            )}
          </VStack>
        </div>
      </VStack>
    )
  }
}
