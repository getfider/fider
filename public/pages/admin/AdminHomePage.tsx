import * as React from 'react';

import { Header, Footer, Button } from '@fider/components/common';

import { Tenant } from '@fider/models';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

interface AdminHomePageState {
    title: string;
    welcome: string;
    invitation: string;
}

export class AdminHomePage extends React.Component<{}, AdminHomePageState> {

    @inject(injectables.Session)
    public session: Session;

    private tenant: Tenant;

    constructor(props: {}) {
        super(props);
        this.tenant = this.session.getCurrentTenant();

        this.state = {
           title: this.tenant.name,
           welcome: '',
           invitation: ''
        };
    }

    public render() {
      return <div>
                <Header />
                    <div className="ui container">
                        <h1 className="ui header">Administration</h1>
                        <h4 className="ui header">General Settings</h4>

                        <div className="ui form">
                            <div className="six wide field">
                                <label htmlFor="title">Title</label>
                                <input id="title"
                                       type="text"
                                       placeholder="Title"
                                       maxLength={60}
                                       value={this.state.title}
                                       onChange={(e) => this.setState({ title: e.currentTarget.value })}/>
                                <p className="info">
                                    <p>Use this field to change the title that is shown on the top of your page.</p>
                                </p>
                            </div>
                            <div className="six wide field">
                                <label htmlFor="welcome-message">Welcome Message</label>
                                <textarea id="welcome-message"
                                          placeholder="Welcome Message"
                                          onChange={(e) => this.setState({ welcome: e.currentTarget.value })}/>
                                <p className="info">
                                    <p>Use this space to change message of your initial page.</p>
                                    <p>Common use cases for this is to explain what is your Company/Product, why you created this space</p>
                                    <p>This field is powered by Commonmark, which means you can style and add links to your message. Learn more <a target="_blank" href="http://commonmark.org/help/">http://commonmark.org/help/</a>.</p>
                                </p>
                            </div>
                            <div className="six wide field">
                                <label htmlFor="invitation">Invitation</label>
                                <input id="invitation"
                                       type="text"
                                       placeholder="Invitation"
                                       onChange={(e) => this.setState({ invitation: e.currentTarget.value })}/>
                                <p className="info">
                                    <p>This is your customized message invition your users to share their ideas and suggestions.</p>
                                </p>
                            </div>
                            <div className="six wide field">
                                <Button className="positive" size="tiny">Confirm</Button>
                            </div>
                        </div>
                    </div>
                <Footer />
            </div>;
    }
}
