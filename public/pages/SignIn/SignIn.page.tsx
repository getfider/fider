import React from "react"
import { SignInControl, TenantLogo, LegalNotice } from "@fider/components"
import { notify } from "@fider/services"
import { Trans } from "@lingui/macro"
import { useFider } from "@fider/hooks"

const Locked = (): JSX.Element => {
  const fider = useFider()
  return (
    <>
      <p className="text-title">
        <Trans id="signin.message.locked.title">
          <strong>{fider.session.tenant.name}</strong> is currently locked.
        </Trans>
      </p>
      <Trans id="signin.message.locked.text">To reactivate this site, sign in with an administrator account and update the required settings.</Trans>
    </>
  )
}

const Private = (): JSX.Element => {
  const fider = useFider()
  return (
    <>
      <p className="text-title">
        <Trans id="signin.message.private.title">
          <strong>{fider.session.tenant.name}</strong> is a private space and requires an invitation to join it.
        </Trans>
      </p>
      <Trans id="signin.message.private.text">If you have an account or an invitation, you may use following options to sign in.</Trans>
    </>
  )
}

export const SignInPage = () => {
  const fider = useFider()

  const onEmailSent = (email: string) => {
    notify.success(
      <span>
        <Trans id="signin.message.emailsent">
          We have just sent a confirmation link to <b>{email}</b>. Click the link and you’ll be signed in.
        </Trans>
      </span>
    )
  }

  return (
    <div id="p-signin" className="page container w-max-6xl">
      <div className="h-20 text-center mb-4">
        <TenantLogo size={100} />
      </div>
      <div className="text-center w-max-4xl mx-auto mb-4">{fider.session.tenant.isPrivate ? <Private /> : <Locked />}</div>

      <SignInControl onEmailSent={onEmailSent} useEmail={true} redirectTo={fider.settings.baseURL} />
      <LegalNotice />
    </div>
  )
}

export default SignInPage
