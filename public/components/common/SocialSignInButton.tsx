import * as React from 'react';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { Button } from '@fider/components/common';
import { AuthSettings } from '@fider/models';

interface SocialSignInButtonProps {
    provider: 'google' | 'facebook' | 'github';
    size: 'small' | 'normal';
}

const providers = {
    google: {
        name: 'Google',
        class: 'google plus',
        icon: 'google plus'
    },
    facebook: {
        name: 'Facebook',
        class: 'facebook',
        icon: 'facebook'
    },
    github: {
        name: 'GitHub',
        class: 'github black',
        icon: 'github'
    }
};

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {

    @inject(injectables.Session)
    public session: Session;

    public render() {
        const auth = this.session.get<AuthSettings>('auth');

        const href = `${auth.endpoint}/oauth/${this.props.provider}?redirect=${location.href}`;
        const classes = `${providers[this.props.provider].class} ${this.props.size === 'small' ? 'circular icon' : 'fluid'}`;

        return <Button href={href} className={classes}>
                    <i className={'icon ' + providers[this.props.provider].icon}></i>
                    { this.props.size === 'normal' && `Log in with ${providers[this.props.provider].name}` }
                </Button>;
    }
}
