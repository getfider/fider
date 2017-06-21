import * as React from 'react';

import { Footer } from '../shared/Footer';
import { EnvironmentInfo, Gravatar } from '../shared/Common';
import { AppSettings } from '../models';
import { isSingleHostMode, getAppSettings } from '../storage';
import { SocialSignInList } from '../shared/SocialSignInList';
import { setTitle, getQueryString } from '../page';
import axios from 'axios';
const td = require('throttle-debounce');

const logo = require('../imgs/logo.png');

import './signup.scss';

interface OAuthUser {
    jwt: string;
    name: string;
    email: string;
}

interface SignUpPageState {
    subdomain: { available: boolean, message?: string };
}

export class SignUpPage extends React.Component<{}, SignUpPageState> {
    private settings: AppSettings;
    private user: OAuthUser;

    constructor(props: {}) {
        super(props);
        this.state = {
            subdomain: { available: false }
        };

        this.checkAvailability = td.debounce(300, this.checkAvailability);

        this.settings = getAppSettings();
        setTitle(isSingleHostMode() ? 'Installation · Fider' : 'New tenant sign up · Fider');

        const jwt = getQueryString('jwt');
        if (jwt) {
            const segments = jwt.split('.');
            const data = JSON.parse(window.atob(segments[1]));
            this.user = {
                jwt,
                name: data['oauth/name'],
                email: data['oauth/email']
            };
        }
    }

    private async checkAvailability(subdomain: string) {
        const url = `/api/tenants/${subdomain}/availability`;
        const result = await axios.get(url);

        this.setState({
            subdomain: {
                available: !result.data.message,
                message: result.data.message
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
                    <div className="ui form">
                        <div className="inline field">
                            <label>Name</label>
                            <input type="text" placeholder="Your organization name"/>
                        </div>
                        { !isSingleHostMode() && <div className="inline field">
                            <label>Identity</label>
                            <div className="ui right labeled input">
                                <div className="ui label">https://</div>
                                <input type="text" placeholder="orgname" onChange={(e) => this.checkAvailability(e.currentTarget.value)} />
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

                    <button className="ui positive button disabled">Confirm</button>
                </div>
                <Footer />
            </div>;
    }
}
