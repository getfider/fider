import * as React from 'react';
import Textarea from 'react-textarea-autosize';

import { Tenant } from '@fider/models';
import { setTitle } from '@fider/utils/page';

import { Header, Footer, Button, DisplayError } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session, TenantService, Failure } from '@fider/services';

interface AdminHomePageState {
    title: string;
    welcomeMessage: string;
    invitation: string;
    error?: Failure;
}

export class AdminHomePage extends React.Component<{}, AdminHomePageState> {

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.TenantService)
    public service: TenantService;

    private tenant: Tenant;

    constructor(props: {}) {
        super(props);
        this.tenant = this.session.getCurrentTenant();

        this.state = {
           title: this.tenant.name,
           welcomeMessage: this.tenant.welcomeMessage,
           invitation: this.tenant.invitation
        };

        setTitle(`Administration · ${document.title}`);
    }

    private async confirm() {
        const result = await this.service.updateSettings(this.state.title, this.state.invitation, this.state.welcomeMessage);
        if (result.ok) {
            location.href = `/`;
        } else if (result.error) {
            this.setState({ error: result.error });
        }
    }

    public render() {
      return <div>
                <Header />
                    <div className="page ui container">
                        <h1 className="ui header">Administration</h1>
                        <h4 className="ui header">General Settings</h4>

                        <div className="ui grid">
                            <div className="eight wide computer sixteen wide mobile column">
                                <div className="ui form">
                                    <DisplayError fields={['title']} error={this.state.error} />
                                    <div className="field">
                                        <label htmlFor="title">Title</label>
                                        <input id="title"
                                            type="text"
                                            placeholder="Title"
                                            maxLength={60}
                                            disabled={ !this.session.isAdmin() }
                                            value={ this.state.title }
                                            onChange={(e) => this.setState({ title: e.currentTarget.value })}/>
                                        <p className="info">
                                            <p>Use this field to change the title that is shown on the top of your page.</p>
                                        </p>
                                    </div>
                                    <DisplayError fields={['welcomeMessage']} error={this.state.error} />
                                    <div className="field">
                                        <label htmlFor="welcome-message">Welcome Message</label>
                                        <Textarea id="welcome-message"
                                                placeholder="Welcome Message"
                                                disabled={ !this.session.isAdmin() }
                                                onChange={(e) => this.setState({ welcomeMessage: e.currentTarget.value })}
                                                value={ this.state.welcomeMessage } />
                                        <p className="info">
                                            <p>Use this space to change message of your initial page.</p>
                                            <p>Common use case for this area is a brief description of what is your company/product, why you created this space and how the visitors can collaborate.</p>
                                            <p>This field is powered by Commonmark. You can style and add links to your message. Learn more at <a target="_blank" href="http://commonmark.org/help/">http://commonmark.org/help/</a>.</p>
                                        </p>
                                    </div>
                                    <DisplayError fields={['invitation']} error={this.state.error} />
                                    <div className="field">
                                        <label htmlFor="invitation">Invitation</label>
                                        <input id="invitation"
                                            type="text"
                                            placeholder="Invitation"
                                            maxLength={60}
                                            disabled={ !this.session.isAdmin() }
                                            value={ this.state.invitation }
                                            onChange={(e) => this.setState({ invitation: e.currentTarget.value })}/>
                                        <p className="info">
                                            <p>This is your customized message to invite visitors to share their ideas and suggestions.</p>
                                        </p>
                                    </div>
                                    {
                                        this.session.isAdmin() &&
                                        <div className="field">
                                            <Button className="positive" size="tiny" onClick={async () => await this.confirm()}>Confirm</Button>
                                        </div>
                                    }
                                </div>
                            </div>
                        </div>
                    </div>
                <Footer />
            </div>;
    }
}
