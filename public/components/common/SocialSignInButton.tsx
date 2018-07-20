import * as React from "react";
import { Button, OAuthProviderLogo } from "@fider/components/common";
import { classSet } from "@fider/services";

interface SocialSignInButtonProps {
  option: {
    displayName: string;
    url?: string;
    logoId?: number;
    logoUrl?: string;
  };
  redirectTo?: string;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = this.props.option.url ? `${this.props.option.url}?redirect=${redirectTo}` : undefined;

    return (
      <Button href={href} fluid={true} className="m-social">
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
