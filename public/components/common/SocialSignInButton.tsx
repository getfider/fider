import * as React from "react";
import { Button } from "@fider/components/common";
import { OAuthProviderOption } from "@fider/models";
import { classSet } from "@fider/services";

interface OAuthProviderPreview {
  displayName: string;
  logoUrl?: string;
}

interface SocialSignInButtonProps {
  option: OAuthProviderOption | OAuthProviderPreview;
  redirectTo?: string;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = "url" in this.props.option ? `${this.props.option.url}?redirect=${redirectTo}` : undefined;
    const className = classSet({
      "m-social": true
    });

    return (
      <Button href={href} fluid={true} className={className}>
        {this.props.option.logoUrl && <img alt={this.props.option.displayName} src={this.props.option.logoUrl} />}
        <span>{this.props.option.displayName}</span>
      </Button>
    );
  }
}
