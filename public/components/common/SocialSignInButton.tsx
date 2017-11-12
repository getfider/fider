import * as React from 'react';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services/Session';
import { Button } from '@fider/components/common';
import { AuthSettings } from '@fider/models';

interface SocialSignInButtonProps {
  provider: 'google' | 'facebook' | 'github';
}

const providers = {
  google: {
    name: 'Google',
    class: 'google plus',
    icon: 'google'
  },
  facebook: {
    name: 'Facebook',
    class: 'facebook',
    icon: 'facebook f'
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

    return (
      <Button size="small" href={href} className={`${providers[this.props.provider].class} fluid`}>
        <i className={`icon ${providers[this.props.provider].icon}`} />
        <span>{providers[this.props.provider].name}</span>
      </Button>
    );
  }
}
