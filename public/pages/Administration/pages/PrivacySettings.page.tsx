import "./PrivacySettings.page.scss";

import * as React from "react";

import { SystemSettings, CurrentUser, Tenant } from "@fider/models";
import { Button, ButtonClickEvent, Textarea, DisplayError, Toggle, Form2 } from "@fider/components/common";
import { actions, notify, Failure } from "@fider/services";
import { AdminBasePage } from "../components";

interface PrivacySettingsPageProps {
  user: CurrentUser;
  tenant: Tenant;
}

interface PrivacySettingsPageState {
  isPrivate: boolean;
}

export class PrivacySettingsPage extends AdminBasePage<PrivacySettingsPageProps, PrivacySettingsPageState> {
  public id = "p-admin-privacy";
  public name = "privacy";
  public icon = "key";
  public title = "Privacy";
  public subtitle = "Manage your site privacy";

  constructor(props: PrivacySettingsPageProps) {
    super(props);

    this.state = {
      isPrivate: this.props.tenant.isPrivate
    };
  }

  private toggle = async (active: boolean) => {
    this.setState(
      state => ({
        isPrivate: active
      }),
      async () => {
        const response = await actions.updateTenantPrivacy(this.state.isPrivate);
        if (response.ok) {
          notify.success("Your privacy settings have been saved.");
        }
      }
    );
  };

  public content() {
    return (
      <Form2>
        <div className="c-form-field">
          <label htmlFor="private">Private site</label>
          <Toggle disabled={!this.props.user.isAdministrator} active={this.state.isPrivate} onToggle={this.toggle} />
          <p className="info">
            A private site prevents unauthenticated users from viewing or interacting with its content. <br /> If
            enabled, only already registered and invited users will be able to sign in to this site.
          </p>
        </div>
      </Form2>
    );
  }
}
