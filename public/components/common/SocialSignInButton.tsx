import * as React from "react";
import { Button } from "@fider/components/common";

interface SocialSignInButtonProps {
  provider: "google" | "facebook" | "github";
  oauthEndpoint?: string;
  redirectTo?: string;
}

const providers = {
  google: {
    name: "Google",
    class: "m-social m-google"
  },
  facebook: {
    name: "Facebook",
    class: "m-social m-facebook"
  },
  github: {
    name: "GitHub",
    class: "m-social m-github"
  }
};

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
  public render() {
    const redirectTo = this.props.redirectTo || location.href;
    const href = this.props.oauthEndpoint
      ? `${this.props.oauthEndpoint}/oauth/${this.props.provider}?redirect=${redirectTo}`
      : undefined;

    return (
      <Button href={href} fluid={true} className={providers[this.props.provider].class}>
        <i className="svg" />
        <span>{providers[this.props.provider].name}</span>
      </Button>
    );
  }
}
