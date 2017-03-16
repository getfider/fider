import * as React from "react";
import { SocialSignInButton } from "./SocialSignInButton";

const tenant = (window as any)._tenant;
const claims = (window as any)._claims;
const gravatar = (window as any)._gravatar;

export class Header extends React.Component<{}, {}> {
    public render() {
        const profile = claims ?
                        <a className="item right signin">
                            <img className="ui avatar image" src={ "https://www.gravatar.com/avatar/" + gravatar } />
                            { claims["user/name"] }
                            <i className="dropdown icon"></i>
                        </a> :
                        <a className="item right signin">
                            <img className="ui avatar image" src="https://www.gravatar.com/avatar/" />
                            Sign in
                            <i className="dropdown icon"></i>
                        </a>;

        const dropdown = claims ?
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
                    <div className="header item">
                        { tenant.name }
                    </div>
                    { profile }
                    </div>
                </div>
                { dropdown }
               </div>;
    }
}
