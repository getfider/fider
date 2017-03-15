import * as React from "react";

const authEndpoint = (window as any)._authEndpoint;

interface SocialSignInButtonProps {
    provider: string
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, {}> {
    render() {
        const providerClassName = this.props.provider == "google" ? "google plus" : "facebook";
        const providerDisplayName = this.props.provider == "google" ? "Google" : "Facebook";
        
        return  <a href={ authEndpoint + "/oauth/" + this.props.provider + "?redirect=" + location.href } className={providerClassName + " ui fluid button"}>
                  <i className={providerClassName + " icon"}></i>
                  Sign in with { providerDisplayName } 
                </a>;
    }
}