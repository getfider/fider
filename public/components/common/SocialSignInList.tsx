import * as React from 'react';
import { SocialSignInButton } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { AuthSettings } from '@fider/models';

interface SocialSignInListProps {
    size: 'small' | 'normal';
}

export class SocialSignInList extends React.Component<SocialSignInListProps, {}> {

    @inject(injectables.Session)
    public session: Session;

    constructor(props: SocialSignInListProps) {
        super(props);
    }

    public render() {
        const settings = this.session.get<AuthSettings>('auth');

        const google = settings.providers.google &&
                        <div className="column">
                            <SocialSignInButton provider="google" size={this.props.size} />
                        </div>;
        const facebook = settings.providers.facebook &&
                        <div className="column">
                            <SocialSignInButton provider="facebook" size={this.props.size} />
                        </div>;
        const github = settings.providers.github &&
                        <div className="column">
                            <SocialSignInButton provider="github" size={this.props.size} />
                        </div>;

        return  <div className="signin-options">
                    <p className="info">We'll never post to any of your accounts.</p>
                    <div className="ui stackable three column centered grid">
                        { facebook }
                        { google }
                        { github }
                    </div>
                    <div className="ui horizontal divider">OR</div>
                    <div id="signin-email" className="ui small action fluid input">
                        <input type="text" placeholder="Log in with e-mail" className="small" />
                        <button className="ui small disabled primary button">Log in</button>
                    </div>
                </div>;
    }
}