import React from "react"
import { SignInControl, TenantLogo, LegalNotice } from "@fider/components"
import { notify, Fider } from "@fider/services"

const Locked = (): JSX.Element => (
  <>
    <p className="text-title">
      <strong>{Fider.session.tenant.name}</strong> is currently locked.
    </p>
    <p>To reactivate this site, sign in with an administrator account and update the required settings.</p>
  </>
)

const Private = (): JSX.Element => (
  <>
    <p className="text-title">
      <strong>{Fider.session.tenant.name}</strong> is a private space and requires an invitation to join it.
    </p>
    <p>If you have an account or an invitation, you may use following options to sign in.</p>
  </>
)

export default class SignInPage extends React.Component<any, any> {
  private onEmailSent = (email: string) => {
    notify.success(
      <span>
        We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
      </span>
    )
  }

  public render() {
    return (
      <div id="p-signin" className="page container w-max-6xl">
        <div className="h-20 text-center mb-4">
          <TenantLogo size={100} />
        </div>
        <div className="text-center w-max-4xl mx-auto mb-4">{Fider.session.tenant.isPrivate ? <Private /> : <Locked />}</div>

        <SignInControl onEmailSent={this.onEmailSent} useEmail={true} redirectTo={Fider.settings.baseURL} />
        <LegalNotice />
      </div>
    )
  }
}
