import "./SignIn.page.scss";

import * as React from "react";
import { SignInControl } from "@fider/components/common";
import { page } from "@fider/services";
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

  public render() {
    return (
      <div className="page ui container">
        <div className="message">
          <p className="welcome">
            <strong>{this.props.tenant.name}</strong> is a private space and requires an invitation to join it.
          </p>
          <p>If you have an account or an invitation, you may use following options to sign in.</p>
        </div>
        <SignInControl useEmail={true} redirectTo={page.getBaseUrl()} />
      </div>
    );
  }
}
