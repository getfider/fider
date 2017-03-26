import * as React from "react";
import { Tenant } from "../models";
import { get, getCurrentUser } from "../storage";
import { Gravatar } from "./Common";
import { SocialSignInButton } from "./SocialSignInButton";

export class Header extends React.Component<{}, {}> {
    public render() {
        const user = getCurrentUser();
        const tenant = get<Tenant>("tenant");

        const profile = user ?
                        <a className="item right signin">
                            <Gravatar email={user.email} />
                            { user.name }
                            <i className="dropdown icon"></i>
                        </a> :
                        <a className="item right signin">
                            <Gravatar />
                            Sign in
                            <i className="dropdown icon"></i>
                        </a>;

        const dropdown = user ?
                        <div id="user-popup" className="ui popup top left transition hidden">
                            <div className="ui">
                                <div className="item">
                                <a href="/logout?redirect=/">
                                    Sign out
                                </a>
                                </div>
                            </div>
                        </div> :
                        <div id="user-popup" className="ui popup top left transition hidden">
                            <div className="ui list">
                                <div className="item">
                                    <SocialSignInButton provider="facebook"/>
                                </div>
                                <div className="item">
                                    <SocialSignInButton provider="google"/>
                                </div>
                            </div>
                        </div>;

        return <div>
                <div className="ui menu">
                    <div className="ui container">
                        <a href="/" className="header item">
                            { tenant.name }
                        </a>
                        { profile }
                    </div>
                </div>
                { dropdown }
               </div>;
    }
}
