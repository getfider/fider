import * as React from 'react';

import { Footer } from '../shared/Footer';
import { EnvironmentInfo, Gravatar } from '../shared/Common';
import { AppSettings } from '../models';
import { SocialSignInList } from '../shared/SocialSignInList';
import { setTitle, getQueryString } from '../page';
import { DisplayError } from '../shared/Common';
import axios from 'axios';
import { decode } from '../jwt';
const td = require('throttle-debounce');
const logo = require('../imgs/logo.png');

import { inject, injectables } from '../di';
import { Session } from '../services/Session';

import './signup.scss';

interface OAuthUser {
    token: string;
    name: string;
    email: string;
}

interface SignUpPageState {
    name?: string;
    clicked: boolean;
    error?: Error;
    subdomain: {
        available: boolean,
        message?: string,
        value?: string
    };
}

export class SignUpPage extends React.Component<{}, SignUpPageState> {
    private settings: AppSettings;
    private user: OAuthUser;

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.state = {
            clicked: false,
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
        this.setState({
            clicked: true,
            error: undefined
        });

        try {
            const response = await axios.post('/api/tenants', {
                token: this.user ? this.user.token : null,
                name: this.state.name,
                subdomain: this.state.subdomain.value,
            });

            if (this.session.isSingleHostMode()) {
                location.reload();
            } else {
                location.href = `${location.protocol}//${this.state.subdomain.value}${this.settings.domain}`;
            }
        } catch (ex) {
            this.setState({
                clicked: false,
                error: ex.response.data
            });
        }
    }

    private async checkAvailability(subdomain: string) {
        const url = `/api/tenants/${subdomain}/availability`;
        const result = await axios.get(url);

        this.setState({
            subdomain: {
                value: subdomain,
                available: !result.data.message,
                message: result.data.message
            }
        });
    }

    public render() {
        const buttonClasses = `ui positive button ${this.state.clicked && 'loading disabled'}`;
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
                    <div className="ui form">
                        <div className="inline field">
                            <label>Name</label>
                            <input id="name" type="text" placeholder="Your organization name" onChange={(e) => this.setState({ name: e.currentTarget.value })}/>
                        </div>
                        { !this.session.isSingleHostMode() && <div className="inline field">
                            <label>Identity</label>
                            <div className="ui right labeled input">
                                <div className="ui label">https://</div>
                                <input id="subdomain" type="text" placeholder="orgname" onChange={(e) => this.checkAvailability(e.currentTarget.value)} />
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
                    </div>
                    <div className="ui section divider"></div>

                    <h3 className="ui header">3. Review and continue</h3>

                    <p>Please make sure information provided above is correct before proceeding.</p>
                    <DisplayError error={this.state.error} />

                    <button className={buttonClasses} onClick={() => this.confirm()}>Confirm</button>
                </div>
                <Footer />
            </div>;
    }
}
