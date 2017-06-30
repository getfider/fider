import * as React from 'react';

import { Footer, Button, Form, EnvironmentInfo, Gravatar, DisplayError } from '@fider/components/common';
import { AppSettings } from '@fider/models';
import { SocialSignInList } from '@fider/components/SocialSignInList';
import { setTitle, getQueryString } from '@fider/utils/page';
import { decode } from '@fider/utils/jwt';
const td = require('throttle-debounce');
const logo = require('@fider/images/logo.png');

import { inject, injectables } from '@fider/di';
import { Session, TenantService, Failure } from '@fider/services';

import './signup.scss';

interface OAuthUser {
    token: string;
    name: string;
    email: string;
}

interface SignUpPageState {
    name?: string;
    error?: Failure;
    subdomain: {
        available: boolean,
        message?: string,
        value?: string
    };
}

export class SignUpPage extends React.Component<{}, SignUpPageState> {
    private settings: AppSettings;
    private user: OAuthUser;
    private form: Form;

    @inject(injectables.Session)
    public session: Session;

    @inject(injectables.TenantService)
    public service: TenantService;

    constructor(props: {}) {
        super(props);
        this.state = {
            subdomain: { available: false }
        };

        this.settings = this.session.getAppSettings();
        this.checkAvailability = td.debounce(300, this.checkAvailability);

        setTitle(this.session.isSingleHostMode() ? 'Installation · Fider' : 'New tenant sign up · Fider');

        const token = getQueryString('jwt');
        if (token) {
            const data = decode(token);
            this.user = {
                token,
                name: data['oauth/name'],
                email: data['oauth/email']
            };
        }
    }

    private async confirm() {
        this.form.clearFailure();

        const result = await this.service.create(
            this.user && this.user.token,
            this.state.name,
            this.state.subdomain.value
        );
        if (result.ok) {
            if (this.session.isSingleHostMode()) {
                location.reload();
            } else {
                location.href = `${location.protocol}//${this.state.subdomain.value}${this.settings.domain}`;
            }
        } else if (result.error) {
            this.form.setFailure(result.error);
        }
    }

    private async checkAvailability(subdomain: string) {
        const result = await this.service.checkAvailability(subdomain);

        this.setState({
            subdomain: {
                value: subdomain,
                available: !result.data.message,
                message: result.data.message,
            }
        });
    }

    public render() {
        return <div>
                <EnvironmentInfo />
                <div id="fdr-signup-page" className="ui container">
                    <img className="logo" src={logo} />

                    <h3 className="ui header">1. Who are you?</h3>

                    {
                        this.user ?
                        <div>
                            <p>
                                Hello, &nbsp;
                                <Gravatar email={this.user.email} />
                                <b>{this.user.name}</b> ({this.user.email})
                            </p>
                        </div> :
                        <div>
                            <p>We need to identify you in order to setup your new Fider instance.</p>
                            <SocialSignInList size="normal" orientation="horizontal" />
                        </div>
                    }

                    <div className="ui section divider"></div>
                    <h3 className="ui header">2. Organization details</h3>
                    <Form ref={(f) => { this.form = f!; } } onSubmit={() => this.confirm()}>
                        <div className="inline field">
                            <label>Name</label>
                            <input id="name" type="text"
                                placeholder="Your organization name"
                                maxLength={60}
                                onChange={(e) => this.setState({ name: e.currentTarget.value })}/>
                        </div>
                        { !this.session.isSingleHostMode() && <div className="inline field">
                            <label>Identity</label>
                            <div className="ui right labeled input">
                                <div className="ui label">https://</div>
                                <input id="subdomain" type="text"
                                    maxLength={40}
                                    placeholder="orgname"
                                    onChange={(e) => this.checkAvailability(e.currentTarget.value)} />
                                <div className="ui label">{ this.settings.domain }</div>
                                {
                                    this.state.subdomain.available &&
                                    <div className="ui left pointing green basic label">
                                        Great!
                                    </div>
                                }
                                {
                                    this.state.subdomain.message &&
                                    <div className="ui left pointing red basic label">
                                        { this.state.subdomain.message }
                                    </div>
                                }
                            </div>
                        </div> }
                    </Form>
                    <div className="ui section divider"></div>

                    <h3 className="ui header">3. Review and continue</h3>

                    <p>Please make sure information provided above is correct before proceeding.</p>

                    <Button classes="positive" onClick={() => this.form.submit()}>Confirm</Button>
                </div>
                <Footer />
            </div>;
    }
}
