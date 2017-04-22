import * as React from "react";
import { Tenant } from "../models";
import { get, getCurrentUser } from "../storage";
import { Gravatar } from "./common";
import { SocialSignInButton } from "./social_signin_button";

export class Header extends React.Component<{}, {}> {
    private dropdown: HTMLElement;
    private list: HTMLElement;

    public componentDidMount() {
        $(this.dropdown).popup({
            inline: true,
            hoverable: true,
            popup: this.list,
            position : "bottom right"
        });
    }

    public render() {
        const user = getCurrentUser();
        const tenant = get<Tenant>("tenant");

        const items = user ? <div className="ui list">
                                <div className="item">
                                    { user.email }
                                </div>
                                <div className="item right">
                                    <a href="/logout?redirect=/">
                                        Sign out
                                    </a>
                                </div>
                             </div> :
                             <div className="ui list">
                                <div className="item">
                                    <SocialSignInButton provider="facebook"/>
                                </div>
                                <div className="item">
                                    <SocialSignInButton provider="google"/>
                                </div>
                             </div>;

        return <div>
                <div className="ui menu">
                    <div className="ui container">
                        <a href="/" className="header item">
                            { tenant.name }
                        </a>
                        <a ref={(e) => { this.dropdown = e; } } className="item right signin">
                            <Gravatar email={user.email} />
                            { user.name || "Sign in" }
                            <i className="dropdown icon"></i>
                        </a>
                    </div>
                </div>
                <div ref={(e) => { this.list = e; } } className="ui popup transition hidden">
                    { items }
                </div>
               </div>;
    }
}
