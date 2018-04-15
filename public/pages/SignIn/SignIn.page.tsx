import "./SignIn.page.scss";

import * as React from "react";
import { SignInControl } from "@fider/components/common";
import { page, notify } from "@fider/services";
import { Tenant } from "@fider/models";

interface SignInPageProps {
  tenant: Tenant;
}

export class SignInPage extends React.Component<SignInPageProps, {}> {
  constructor(props: SignInPageProps) {
    super(props);
    this.state = {};

    page.setTitle(`Sign in · ${document.title}`);
  }

  private onEmailSent = (email: string) => {
    notify.success(
      <span>
        We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
      </span>
    );
  };

  public render() {
    return (
      <div className="page ui container">
        <div className="message">
          <p className="welcome">
            <strong>{this.props.tenant.name}</strong> is a private space and requires an invitation to join it.
          </p>
          <p>If you have an account or an invitation, you may use following options to sign in.</p>
        </div>
        <SignInControl onEmailSent={this.onEmailSent} useEmail={true} redirectTo={page.getBaseUrl()} />
      </div>
    );
  }
}
