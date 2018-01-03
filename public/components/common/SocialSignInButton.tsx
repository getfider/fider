import * as React from 'react';
import { Button } from '@fider/components/common';

interface SocialSignInButtonProps {
  oauthEndpoint: string;
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
  public render() {
    const href = `${this.props.oauthEndpoint}/oauth/${this.props.provider}?redirect=${location.href}`;

    return (
      <Button size="small" href={href} className={`${providers[this.props.provider].class} fluid`}>
        <i className={`icon ${providers[this.props.provider].icon}`} />
        <span>{providers[this.props.provider].name}</span>
      </Button>
    );
  }
}
