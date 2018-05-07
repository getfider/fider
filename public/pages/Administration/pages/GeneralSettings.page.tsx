import * as React from "react";

import { SystemSettings, CurrentUser, Tenant } from "@fider/models";
import { Button, ButtonClickEvent, Textarea, DisplayError } from "@fider/components/common";
import { actions, page, Failure, fileToBase64 } from "@fider/services";
import { AdminBasePage } from "../components";

interface GeneralSettingsPageProps {
  user: CurrentUser;
  tenant: Tenant;
  system: SystemSettings;
  publicIP: string;
}

interface GeneralSettingsPageState {
  logo: string;
  title: string;
  welcomeMessage: string;
  invitation: string;
  cname: string;
  error?: Failure;
}

export class GeneralSettingsPage extends AdminBasePage<GeneralSettingsPageProps, GeneralSettingsPageState> {
  public id = "p-admin-general";
  public name = "general";
  public icon = "settings";
  public title = "General";
  public subtitle = "Manage your site settings";

  constructor(props: GeneralSettingsPageProps) {
    super(props);

    this.state = {
      logo: "",
      title: this.props.tenant.name,
      cname: this.props.tenant.cname,
      welcomeMessage: this.props.tenant.welcomeMessage,
      invitation: this.props.tenant.invitation
    };

    page.setTitle(`General · Site Settings · ${document.title}`);
  }

  private async save(e: ButtonClickEvent) {
    const result = await actions.updateTenantSettings(
      this.state.logo,
      this.state.title,
      this.state.invitation,
      this.state.welcomeMessage,
      this.state.cname
    );
    if (result.ok) {
      e.preventEnable();
      location.href = `/`;
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  public fileChanged = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const base64 = await fileToBase64(e.target.files[0]);
      this.setState({
        logo: base64
      });
    }
  };

  public dnsInstructions(): JSX.Element {
    const isApex = this.state.cname.split(".").length === 2;
    const recordType = isApex ? "A" : "CNAME";
    const publicIP = this.props.publicIP || "<error>";
    const targetRecord = isApex ? publicIP : `${this.props.tenant.subdomain}${this.props.system.domain}`;
    return (
      <>
        <strong>{this.state.cname}</strong> {recordType} <strong>{targetRecord}</strong>
      </>
    );
  }

  public content() {
    return (
      <div className="ui form">

        <input type="file" onChange={this.fileChanged} />

        <DisplayError fields={["title"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="title">Title</label>
          <input
            id="title"
            type="text"
            maxLength={60}
            disabled={!this.props.user.isAdministrator}
            value={this.state.title}
            onChange={e => this.setState({ title: e.currentTarget.value })}
          />
          <div className="info">
            <p>Use this field to change the title that is shown on the top of your page.</p>
          </div>
        </div>
        <DisplayError fields={["welcomeMessage"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="welcome-message">Welcome Message</label>
          <Textarea
            id="welcome-message"
            disabled={!this.props.user.isAdministrator}
            onChange={e => this.setState({ welcomeMessage: e.currentTarget.value })}
            value={this.state.welcomeMessage}
          />
          <div className="info">
            <p>Use this space to change message of your initial page.</p>
            <p>
              Common use case for this area is a brief description of what is your company/product, why you created this
              space and how the visitors can collaborate.
            </p>
            <p>
              This field is powered by CommonMark. You can style and add links to your message. Learn more at{" "}
              <a target="_blank" href="http://commonmark.org/help/">
                {"http://commonmark.org/help/"}
              </a>.
            </p>
          </div>
        </div>
        <DisplayError fields={["invitation"]} error={this.state.error} />
        <div className="field">
          <label htmlFor="invitation">Invitation</label>
          <input
            id="invitation"
            type="text"
            maxLength={60}
            disabled={!this.props.user.isAdministrator}
            value={this.state.invitation}
            onChange={e => this.setState({ invitation: e.currentTarget.value })}
          />
          <div className="info">
            <p>This is your customized message to invite visitors to share their ideas and suggestions.</p>
          </div>
        </div>
        {!page.isSingleHostMode() && [
          <DisplayError key={1} fields={["cname"]} error={this.state.error} />,
          <div key={2} className="field">
            <label htmlFor="cname">Custom Domain</label>
            <input
              id="cname"
              type="text"
              placeholder="feedback.yourcompany.com"
              maxLength={100}
              disabled={!this.props.user.isAdministrator}
              value={this.state.cname}
              onChange={e => this.setState({ cname: e.currentTarget.value })}
            />
            <div className="info">
              {this.state.cname ? (
                [
                  <p key={0}>Enter the following record into your DNS zone records:</p>,
                  <p key={1}>{this.dnsInstructions()}</p>,
                  <p key={2}>
                    Please note that it may take up to 72 hours for the change to take effect worldwide due to DNS
                    propagation.
                  </p>
                ]
              ) : (
                <p>
                  Custom domains allow you to access your app via your own domain name (for example,{" "}
                  <code>feedback.yourcomany.com</code>).
                </p>
              )}
            </div>
          </div>
        ]}
        {this.props.user.isAdministrator && (
          <div className="field">
            <Button color="green" onClick={async e => await this.save(e)}>
              Save
            </Button>
          </div>
        )}
      </div>
    );
  }
}
