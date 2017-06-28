import * as React from 'react';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';
import { Button } from '../components/common/Button';

interface SocialSignInButtonProps {
    provider: string;
    size: 'small' | 'normal';
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {

    @inject(injectables.Session)
    public session: Session;

    public render() {
        const auth = this.session.get<any>('auth');
        const providerClassName = this.props.provider === 'google' ? 'google plus' : 'facebook';
        const providerDisplayName = this.props.provider === 'google' ? 'Google' : 'Facebook';
        const href = `${auth.endpoint}/oauth/${this.props.provider}?redirect=${location.href}`;
        const classes = `${providerClassName} ${this.props.size === 'small' ? 'circular icon' : 'fluid'}`;

        return <Button href={href} classes={classes}>
                    <i className={providerClassName + ' icon'}></i>
                    { this.props.size === 'normal' && `Sign in with ${providerDisplayName}` }
                </Button>;
    }
}
