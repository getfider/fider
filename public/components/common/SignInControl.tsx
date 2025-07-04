import "./SignInControl.scss"

import React, { useState } from "react"
import { SocialSignInButton, Form, Button, Input, Message } from "@fider/components"
import { Divider } from "@fider/components/layout"
import { device, actions, Failure, isCookieEnabled } from "@fider/services"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"

interface SignInControlProps {
  useEmail: boolean
  redirectTo?: string
  onSubmit?: () => Promise<SignInSubmitResponse>
  onEmailSent?: (email: string) => void
  signInButtonText?: string
}

export interface SignInSubmitResponse {
  ok: boolean
  code?: string
}

export const SignInControl: React.FunctionComponent<SignInControlProps> = (props) => {
  const fider = useFider()
  const [showEmailForm, setShowEmailForm] = useState(fider.session.tenant ? fider.session.tenant.isEmailAuthAllowed : true)
  const [email, setEmail] = useState("")
  const [error, setError] = useState<Failure | undefined>(undefined)

  const signInText = props.signInButtonText || i18n._({ id: "action.signin", message: "Sign in" })

  const forceShowEmailForm = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault()
    setShowEmailForm(true)
  }

  const doPreSigninAction = async (): Promise<SignInSubmitResponse> => {
    let signInResponse: SignInSubmitResponse = { ok: true }
    if (props.onSubmit) {
      signInResponse = await props.onSubmit()
    }
    return signInResponse
  }

  const onSocialSignin = async () => {
    return await doPreSigninAction()
  }

  const signIn = async () => {
    const signInResponse = await doPreSigninAction()
    if (!signInResponse.ok) {
      return
    }
    const result = await actions.signIn(email, signInResponse.code)
    if (result.ok) {
      setEmail("")
      setError(undefined)
      if (props.onEmailSent) {
        props.onEmailSent(email)
      }
    } else if (result.error) {
      setError(result.error)
    }
  }

  const providersLen = fider.settings.oauth.length

  if (!isCookieEnabled()) {
    return (
      <Message type="error">
        <h3 className="text-display">Cookies Required</h3>
        <p>Cookies are not enabled on your browser. Please enable cookies in your browser preferences to continue.</p>
      </Message>
    )
  }

  return (
    <div className="c-signin-control">
      {providersLen > 0 && (
        <>
          <div className="c-signin-control__oauth pb-3">
            {fider.settings.oauth.map((o) => (
              <React.Fragment key={o.provider}>
                <SocialSignInButton onClick={onSocialSignin} option={o} redirectTo={props.redirectTo} />
              </React.Fragment>
            ))}
          </div>
          {props.useEmail && <Divider />}
        </>
      )}

      {props.useEmail &&
        (showEmailForm ? (
          <div className="pt-3">
            <Form error={error}>
              <Input
                className="text-left"
                field="email"
                value={email}
                autoFocus={!device.isTouch()}
                onChange={setEmail}
                placeholder={i18n._({ id: "signin.email.placeholder", message: "Email address" })}
              />
              <Button className="w-full justify-center" type="submit" variant="primary" disabled={email === ""} onClick={signIn}>
                {signInText}
              </Button>
            </Form>
            {!fider.session.tenant.isEmailAuthAllowed && (
              <p className="text-red-700 mt-1">
                <Trans id="signin.message.onlyadmins">Currently only allowed to sign in to an administrator account</Trans>
              </p>
            )}
          </div>
        ) : (
          <div>
            <p className="text-muted">
              <Trans id="signin.message.emaildisabled">
                Email authentication has been disabled by an administrator. If you have an administrator account and need to bypass this restriction, please{" "}
                <a href="#" className="text-bold" onClick={forceShowEmailForm}>
                  click here
                </a>
                .
              </Trans>
            </p>
          </div>
        ))}
    </div>
  )
}
