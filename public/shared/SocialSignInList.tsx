import * as React from 'react';
import { SocialSignInButton } from './SocialSignInButton';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';

interface AuthSettings {
    endpoint: string;
    providers: {
        google: boolean,
        facebook: boolean
    };
}

interface SocialSignInListProps {
    size: 'small' | 'normal';
    orientation: 'horizontal' | 'vertical';
}

export class SocialSignInList extends React.Component<SocialSignInListProps, {}> {

    @inject(injectables.Session)
    public session: Session;

    constructor(props: SocialSignInListProps) {
        super(props);
    }

    public render() {
        const settings = this.session.get<AuthSettings>('auth');
        const cssClasses = this.props.orientation === 'horizontal' ? 'horizontal divided' : '';

        const google = settings.providers.google &&
                        <div className="item">
                            <SocialSignInButton provider="google" size={this.props.size} />
                        </div>;
        const facebook = settings.providers.facebook &&
                        <div className="item">
                            <SocialSignInButton provider="facebook" size={this.props.size} />
                        </div>;

        const noAuth = !facebook && !google &&
                        <div className="item">
                            <div className="ui tertiary inverted red segment">
                                There are no authentication methods enabled.
                            </div>
                        </div>;

        return <div className={`ui list ${cssClasses}`}>
                    { facebook }
                    { google }
                    { noAuth }
                </div>;
    }
}