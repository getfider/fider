import * as React from "react";

const authEndpoint = (window as any)._authEndpoint;

interface SocialSignInButtonProps {
    provider: string;
    small?: boolean;
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

    public render() {
        const providerClassName = this.props.provider === "google" ? "google plus" : "facebook";
        const providerDisplayName = this.props.provider === "google" ? "Google" : "Facebook";
        const oauthUrl = `${authEndpoint}/oauth/${this.props.provider}?redirect=${location.href}`;
        const cssClasses = `ui button 
                            ${providerClassName} 
                            ${this.state.clicked ? "loading disabled" : ""} 
                            ${this.props.small ? "circular icon" : "fluid"}`;

        if (this.props.small) {
            return <a href={oauthUrl} className={cssClasses} onClick={() => this.setState({clicked: true})}>
                    <i className={providerClassName + " icon"}></i>
                   </a>;
        } else {
            return  <a href={oauthUrl} className={cssClasses} onClick={() => this.setState({clicked: true})}>
                        <i className={providerClassName + " icon"}></i>
                        Sign in with { providerDisplayName } 
                    </a>;
        }
    }
}
