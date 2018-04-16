import * as React from "react";
import { Button } from "@fider/components/common";

interface SocialSignInButtonProps {
  oauthEndpoint: string;
  provider: "google" | "facebook" | "github";
  redirectTo?: string;
}

const providers = {
  google: {
    name: "Google",
    class: "social google"
  },
  facebook: {
    name: "Facebook",
    class: "social facebook"
  },
  github: {
    name: "GitHub",
    class: "social github"
  }
};

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = `${this.props.oauthEndpoint}/oauth/${this.props.provider}?redirect=${redirectTo}`;

    return (
      <Button href={href} fluid={true} className={providers[this.props.provider].class}>
        <i className="svg" />
        <span>{providers[this.props.provider].name}</span>
      </Button>
    );
  }
}
