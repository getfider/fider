import "./SignIn.page.scss";

import React from "react";
import { SignInControl, TenantLogo, LegalNotice } from "@fider/components";
import { notify, Fider } from "@fider/services";

const messages = {
  locked: () => (
    <>
      <p className="welcome">
        <strong>{Fider.session.tenant.name}</strong> is currently locked.
      </p>
      <p>To reactivate this site, sign in with an administrator account and update the required settings.</p>
    </>
  ),
  private: () => (
    <>
      <p className="welcome">
        <strong>{Fider.session.tenant.name}</strong> is a private space and requires an invitation to join it.
      </p>
      <p>If you have an account or an invitation, you may use following options to sign in.</p>
    </>
  )
};

export default class SignInPage extends React.Component<any, {}> {
  private onEmailSent = (email: string) => {
    notify.success(
      <span>
        We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
      </span>
    );
  };

  public render() {
    return (
      <div id="p-signin" className="page container">
        <div className="message">
          <TenantLogo size={100} />
          {Fider.session.tenant.isPrivate ? messages.private() : messages.locked()}
        </div>
        <SignInControl onEmailSent={this.onEmailSent} useEmail={true} redirectTo={Fider.settings.baseURL} />
        <LegalNotice />
      </div>
    );
  }
}
