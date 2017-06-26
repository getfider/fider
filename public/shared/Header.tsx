import * as React from 'react';
import { User, Tenant } from '../models';
import { get, getCurrentUser, isStaff } from '../storage';
import { EnvironmentInfo, Gravatar } from './Common';
import { SocialSignInList } from './SocialSignInList';

export class Header extends React.Component<{}, {}> {
    private dropdown: HTMLElement;
    private list: HTMLElement;
    private user: User;
    private tenant: Tenant;

    constructor(props: {}) {
        super(props);

        this.user = getCurrentUser();
        this.tenant = get<Tenant>('tenant');
    }

    public componentDidMount() {
        $(this.dropdown).popup({
            inline: true,
            hoverable: true,
            popup: this.list,
            position : 'bottom right'
        });
    }

    public render() {
        const items = this.user ?
                        <div className="ui divided list">
                            <div className="item">
                                <b>{ this.user.email }</b>
                            </div>
                            {
                                isStaff() &&
                                <div className="item">
                                    <a href="/admin">
                                        Administration
                                    </a>
                                </div>
                            }
                            <div className="item">
                                <a href="/logout?redirect=/">
                                    Sign out
                                </a>
                            </div>
                        </div> :
                        <SocialSignInList orientation="vertical" size="normal" />;

        return <div>
                    <EnvironmentInfo />
                    <div id="menu" className="ui menu no-border">
                        <div className="ui container">
                            <a href="/" className="header item">
                                { this.tenant.name }
                            </a>
                            <a ref={(e) => { this.dropdown = e!; } } className="item right signin">
                                <Gravatar email={this.user.email} />
                                { this.user.name || 'Sign in' }
                                <i className="dropdown icon"></i>
                            </a>
                        </div>
                    </div>
                    <div ref={(e) => { this.list = e!; } } className="fdr-profile-popup ui popup transition hidden">
                        { items }
                    </div>
               </div>;
    }
}
