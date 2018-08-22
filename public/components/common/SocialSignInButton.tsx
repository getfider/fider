import * as React from "react";
import { Button, OAuthProviderLogo } from "@fider/components/common";
import { classSet } from "@fider/services";

interface SocialSignInButtonProps {
  option: {
    displayName: string;
    provider?: string;
    url?: string;
    logoID?: number;
    logoUrl?: string;
  };
  redirectTo?: string;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = this.props.option.url ? `${this.props.option.url}?redirect=${redirectTo}` : undefined;
    const className = classSet({
      "m-social": true,
      [`m-${this.props.option.provider}`]: this.props.option.provider
    });
    return (
      <Button href={href} rel="nofollow" fluid={true} className={className}>
        {this.props.option.logoUrl ? (
          <img src={this.props.option.logoUrl} />
        ) : (
          <OAuthProviderLogo option={this.props.option} />
        )}
        <span>{this.props.option.displayName}</span>
      </Button>
    );
  }
}
