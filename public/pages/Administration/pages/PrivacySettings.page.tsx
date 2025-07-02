import React from "react"
import { Toggle, Form, Field } from "@fider/components"
import { actions, notify, Fider } from "@fider/services"
import { AdminBasePage } from "@fider/pages/Administration/components/AdminBasePage"

export interface PrivacySettingsPageState {
  isPrivate: boolean
  isFeedEnabled: boolean
  isModerationEnabled: boolean
}

export default class PrivacySettingsPage extends AdminBasePage<any, PrivacySettingsPageState> {
  public id = "p-admin-privacy"
  public name = "privacy"
  public title = "Privacy"
  public subtitle = "Manage your site's privacy"

  constructor(props: any) {
    super(props)

    this.state = {
      isPrivate: Fider.session.tenant.isPrivate,
      isFeedEnabled: Fider.session.tenant.isFeedEnabled,
      isModerationEnabled: Fider.session.tenant.isModerationEnabled,
    }
  }

  private updatePrivacySettings = async (isPrivate: boolean, isFeedEnabled: boolean, isModerationEnabled?: boolean) => {
    this.setState(
      {
        isPrivate,
        isFeedEnabled: isPrivate ? false : isFeedEnabled, // Disable feed if site is private
        isModerationEnabled: isModerationEnabled !== undefined ? isModerationEnabled : this.state.isModerationEnabled,
      },
      async () => {
        const response = await actions.updateTenantPrivacy(this.state)
        if (response.ok) {
          notify.success("Your privacy settings have been saved.")
        }
      }
    )
  }

  private privacyToggle = async (active: boolean) => {
    this.updatePrivacySettings(active, this.state.isFeedEnabled)
  }

  private atomFeedToggle = async (enabled: boolean) => {
    this.updatePrivacySettings(this.state.isPrivate, enabled)
  }

  private moderationToggle = async (enabled: boolean) => {
    this.updatePrivacySettings(this.state.isPrivate, this.state.isFeedEnabled, enabled)
  }

  public content() {
    return (
      <Form>
        <Field label="Private Site">
          <Toggle disabled={!Fider.session.user.isAdministrator} active={this.state.isPrivate} onToggle={this.privacyToggle} />
          <p className="text-muted mt-1">
            A private site prevents unauthenticated users from viewing or interacting with its content. <br /> When enabled, only already registered users,
            invited users and users from trusted OAuth providers will have access to this site. Disables the feed feature.
          </p>
        </Field>
        <Field label="ATOM Feed">
          <Toggle disabled={!Fider.session.user.isAdministrator || this.state.isPrivate} active={this.state.isFeedEnabled} onToggle={this.atomFeedToggle} />
          <p className="text-muted mt-1">
            This feature lets users access this site via a feed reader. <br /> When enabled, the site makes its posts and comments available using the ATOM
            format. Links to feeds and autodiscovery metadata are shown on the site.
          </p>
        </Field>
        <Field label="Content Moderation">
          <Toggle disabled={!Fider.session.user.isAdministrator} active={this.state.isModerationEnabled} onToggle={this.moderationToggle} />
          <p className="text-muted mt-1">
            When enabled, new posts and comments will require approval from an administrator before being visible to other users. <br /> 
            Content creators can see their own unmoderated content, but it will be hidden from other users until approved.
          </p>
        </Field>
      </Form>
    )
  }
}
