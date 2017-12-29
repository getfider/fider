import * as React from 'react';

import { AppSettings, CurrentUser, Tenant } from '@fider/models';
import { isSingleHostMode, setTitle } from '@fider/utils/page';

import { Header, Footer, Button, Textarea, DisplayError } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session, TenantService, Failure } from '@fider/services';

interface AdminHomePageProps {
  user: CurrentUser;
  tenant: Tenant;
  settings: AppSettings;
}

interface AdminHomePageState {
  title: string;
  welcomeMessage: string;
  invitation: string;
  cname: string;
  error?: Failure;
}

export class AdminHomePage extends React.Component<AdminHomePageProps, AdminHomePageState> {

  @inject(injectables.TenantService)
  public service: TenantService;

  constructor(props: AdminHomePageProps) {
    super(props);

    this.state = {
      title: this.props.tenant.name,
      cname: this.props.tenant.cname,
      welcomeMessage: this.props.tenant.welcomeMessage,
      invitation: this.props.tenant.invitation
    };

    setTitle(`Administration · ${document.title}`);
  }

  private async confirm() {
    const result = await this.service.updateSettings(
      this.state.title,
      this.state.invitation,
      this.state.welcomeMessage,
      this.state.cname,
    );
    if (result.ok) {
      location.href = `/`;
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  }

  public render() {
    return (
      <div className="page ui container">
        <h2 className="ui header">
          <i className="circular settings icon" />
          <div className="content">
            General Settings
            <div className="sub header">Manage your site settings</div>
          </div>
        </h2>

        <div className="ui grid">
          <div className="eight wide computer sixteen wide mobile column">
            <div className="ui form">
              <DisplayError fields={['title']} error={this.state.error} />
              <div className="field">
                <label htmlFor="title">Title</label>
                <input
                  id="title"
                  type="text"
                  maxLength={60}
                  disabled={!this.props.user.isAdministrator}
                  value={this.state.title}
                  onChange={(e) => this.setState({ title: e.currentTarget.value })}
                />
                <div className="info">
                  <p>Use this field to change the title that is shown on the top of your page.</p>
                </div>
              </div>
              <DisplayError fields={['welcomeMessage']} error={this.state.error} />
              <div className="field">
                <label htmlFor="welcome-message">Welcome Message</label>
                <Textarea
                  id="welcome-message"
                  disabled={!this.props.user.isAdministrator}
                  onChange={(e) => this.setState({ welcomeMessage: e.currentTarget.value })}
                  value={this.state.welcomeMessage}
                />
                <div className="info">
                  <p>Use this space to change message of your initial page.</p>
                  <p>Common use case for this area is a brief description of what is your company/product, why you created this space and how the visitors can collaborate.</p>
                  <p>This field is powered by Commonmark. You can style and add links to your message. Learn more at <a target="_blank" href="http://commonmark.org/help/">http://commonmark.org/help/</a>.</p>
                </div>
              </div>
              <DisplayError fields={['invitation']} error={this.state.error} />
              <div className="field">
                <label htmlFor="invitation">Invitation</label>
                <input
                  id="invitation"
                  type="text"
                  maxLength={60}
                  disabled={!this.props.user.isAdministrator}
                  value={this.state.invitation}
                  onChange={(e) => this.setState({ invitation: e.currentTarget.value })}
                />
                <div className="info">
                  <p>This is your customized message to invite visitors to share their ideas and suggestions.</p>
                </div>
              </div>
              {
                !isSingleHostMode() && [
                  <DisplayError key={1} fields={['cname']} error={this.state.error} />,
                  <div key={2} className="field">
                    <label htmlFor="cname">Custom Domain</label>
                    <input
                        id="cname"
                        type="text"
                        placeholder="feedback.yourcompany.com"
                        maxLength={100}
                        disabled={!this.props.user.isAdministrator}
                        value={this.state.cname}
                        onChange={(e) => this.setState({ cname: e.currentTarget.value })}
                    />
                    <div className="info">
                      {
                        this.state.cname ? [
                          <p key={0}>Enter the following record into your DNS zone records:</p>,
                          <p key={1}><strong>{this.state.cname}</strong> CNAME <strong>{this.props.tenant.subdomain}{this.props.settings.domain}</strong></p>,
                          <p key={2}>Please note that it may take up to 72 hours for the change to take effect worldwide due to DNS propagation.</p>
                        ] :
                        <p>Custom domains allow you to access your app via your own domain name (for example, <code>feedback.yourcomany.com</code>).</p>
                      }
                    </div>
                  </div>
                ]
              }
              {
                this.props.user.isAdministrator &&
                <div className="field">
                  <Button className="positive" size="tiny" onClick={async () => await this.confirm()}>Confirm</Button>
                </div>
              }
            </div>
          </div>
        </div>
      </div>
    );
  }
}
