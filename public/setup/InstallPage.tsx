import * as React from 'react';

import { Footer } from '../shared/Footer';
import { EnvironmentInfo } from '../shared/Common';
import { AppSettings } from '../models';
import { isSingleHostMode } from '../storage';
import { SocialSignInList } from '../shared/SocialSignInList';
const logo = require('../imgs/logo.png');

import './install.scss';

export class InstallPage extends React.Component<{}, {}> {
    constructor(props: {}) {
        super(props);
    }

    public render() {
      return <div>
                <EnvironmentInfo />
                <div id="fdr-install-page" className="ui container">
                    <img className="logo" src={logo} />

                    <h3 className="ui header">1. Who are you?</h3>
                    <p>We need to identify you in order to setup your new Fider instance.</p>
                    <SocialSignInList size="normal" orientation="horizontal" />
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
                                <input type="text" placeholder="orgname" />
                                <div className="ui label">.dev.fider.io</div>
                                <div className="ui left pointing red basic label">
                                That subdomain is not available!
                                </div>
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
