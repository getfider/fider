import * as React from "react";

const authEndpoint = (window as any)._authEndpoint;

interface SocialSignInButtonProps {
    provider: string;
}
interface SocialSignInButtonState {
    clicked: boolean;
}

export class SocialSignInButton extends React.Component<SocialSignInButtonProps, SocialSignInButtonState> {
    constructor() {
        super();
        this.state = {
            clicked: false
        };
    }

    render() {
        const providerClassName = this.props.provider == "google" ? "google plus" : "facebook";
        const providerDisplayName = this.props.provider == "google" ? "Google" : "Facebook";
        const oauthUrl = `${authEndpoint}/oauth/${this.props.provider}?redirect=${location.href}`;
        const cssClasses = `${providerClassName} ui fluid button ${this.state.clicked ? "loading" : ""}`;

        return  <a href={oauthUrl} className={cssClasses} onClick={() => this.setState({clicked: true})}>
                  <i className={providerClassName + " icon"}></i>
                  Sign in with { providerDisplayName } 
                </a>;
    }
}