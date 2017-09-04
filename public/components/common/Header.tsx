import * as React from 'react';
import { User, Tenant } from '@fider/models';
import { SocialSignInList, EnvironmentInfo, Gravatar } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';

export class Header extends React.Component<{}, {}> {
    private dropdown: HTMLElement;
    private list: HTMLElement;
    private user: User;
    private tenant: Tenant;

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);

        this.user = this.session.getCurrentUser();
        this.tenant = this.session.get<Tenant>('tenant');
    }

    public componentDidMount() {
        if (this.user) {
            $(this.dropdown).popup({
                inline: true,
                hoverable: true,
                popup: this.list,
                position : 'bottom right'
            });
        }
    }

    private showModal() {
        if (!this.user) {
            $('#signin-modal').modal({
                blurring: true
            }).modal('show');
        }
    }

    public render() {
        const items = <div className="ui divided list">
                            {
                                this.session.isStaff() &&
                                <div className="item">
                                    <a href="/admin">
                                        Administration
                                    </a>
                                </div>
                            }
                            <div className="item">
                                <a className="signout" href="/logout?redirect=/">
                                    Log out
                                </a>
                            </div>
                        </div>;

        return <div>
                    <EnvironmentInfo />
                    <div id="menu" className="ui small menu no-border">
                        <div className="ui container">
                            <a href="/" className="header item">
                                { this.tenant.name }
                            </a>
                            <a ref={(e) => { this.dropdown = e!; } } onClick={ () => this.showModal() } className={`item right signin ${!this.user.name && 'subtitle' }`}>
                                <Gravatar name={this.user.name} hash={this.user.gravatar} />
                                { this.user.name || 'Log in' }
                                { this.user.name && <i className="dropdown icon"></i> }
                            </a>
                        </div>
                    </div>
                    <div ref={(e) => { this.list = e!; } } className="fdr-profile-popup ui popup transition hidden">
                        { items }
                    </div>
               </div>;
    }
}
