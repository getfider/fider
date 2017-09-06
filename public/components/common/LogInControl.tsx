import * as React from 'react';
import { SocialLogInButton } from '@fider/components/common';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { AuthSettings } from '@fider/models';

interface LogInControlProps {
    size: 'small' | 'normal';
}

export class LogInControl extends React.Component<LogInControlProps, {}> {

    @inject(injectables.Session)
    public session: Session;

    constructor(props: LogInControlProps) {
        super(props);
    }

    public render() {
        const settings = this.session.get<AuthSettings>('auth');

        const google = settings.providers.google &&
                        <div className="column">
                            <SocialLogInButton provider="google" size={this.props.size} />
                        </div>;
        const facebook = settings.providers.facebook &&
                        <div className="column">
                            <SocialLogInButton provider="facebook" size={this.props.size} />
                        </div>;
        const github = settings.providers.github &&
                        <div className="column">
                            <SocialLogInButton provider="github" size={this.props.size} />
                        </div>;

        return  <div className="login-options">
                    <p className="info">We'll never post to any of your accounts.</p>
                    <div className="ui stackable three column centered grid">
                        { facebook }
                        { google }
                        { github }
                    </div>
                </div>;
    }
}

/*
<div className="ui horizontal divider">OR</div>
<div className="ui small action fluid input">
    <input type="text" placeholder="Log in with e-mail" className="small" />
    <button className="ui small disabled primary button">Log in</button>
</div>
*/
