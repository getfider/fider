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
  onSubmit?: () => void
  onEmailSent?: (email: string) => void
  signInButtonText?: string
  onCodeVerified?: (result: { showProfileCompletion?: boolean; code?: string }) => void
}

export const SignInControl: React.FunctionComponent<SignInControlProps> = (props) => {
  const fider = useFider()
  const [showEmailForm, setShowEmailForm] = useState(fider.session.tenant ? fider.session.tenant.isEmailAuthAllowed : true)
  const [showCodeEntry, setShowCodeEntry] = useState(false)
  const [email, setEmail] = useState("")
  const [code, setCode] = useState("")
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [resendMessage, setResendMessage] = useState("")

  const signInText = props.signInButtonText || i18n._({ id: "action.signin", message: "Sign in" })

  const forceShowEmailForm = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault()
    setShowEmailForm(true)
  }

  const doPreSigninAction = () => {
    if (props.onSubmit) {
      props.onSubmit()
    }
  }

  const onSocialSignin = () => {
    doPreSigninAction()
  }

  const editEmail = () => {
    setShowCodeEntry(false)
    setCode("")
    setError(undefined)
    setResendMessage("")
  }

  const signIn = async () => {
    await doPreSigninAction()
    const result = await actions.signIn(email)
    if (result.ok) {
      setError(undefined)
      setShowCodeEntry(true)
      // Don't call onEmailSent - we're showing code entry inline now
    } else if (result.error) {
      setError(result.error)
    }
  }

  const verifyCode = async () => {
    const result = await actions.verifySignInCode(email, code)
    if (result.ok) {
      const data = result.data as { showProfileCompletion?: boolean } | undefined
      if (props.onCodeVerified) {
        // Let the parent component decide what to do, pass the code along
        props.onCodeVerified({ ...data, code })
      } else {
        // Default behavior: reload the page
        location.reload()
      }
    } else {
      // Handle validation errors - convert data object to Failure format
      const data = result.data as Record<string, string> | undefined
      if (data && typeof data === "object") {
        const errors = Object.entries(data).map(([field, message]) => ({
          field,
          message,
        }))
        setError({ errors })
      } else if (result.error) {
        // Display the error from the server
        setError(result.error)
      }
    }
  }

  const resendCode = async () => {
    setResendMessage("")
    const result = await actions.resendSignInCode(email)
    if (result.ok) {
      setError(undefined)
      setCode("")
      setResendMessage(i18n._({ id: "signin.code.sent", message: "A new code has been sent to your email." }))
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
            {!showCodeEntry ? (
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
                {!fider.session.tenant.isEmailAuthAllowed && (
                  <p className="text-red-700 mt-1">
                    <Trans id="signin.message.onlyadmins">Currently only allowed to sign in to an administrator account</Trans>
                  </p>
                )}
              </Form>
            ) : (
              <div>
                <p className="text-muted mb-2">
                  <Trans id="signin.code.instruction">
                    Please type in the code we just sent to <strong>{email}</strong>
                  </Trans>{" "}
                  <a
                    href="#"
                    className="text-link"
                    onClick={(e) => {
                      e.preventDefault()
                      editEmail()
                    }}
                  >
                    <Trans id="signin.code.edit">Edit</Trans>
                  </a>
                </p>
                <Form error={error}>
                  <Input
                    className="text-left"
                    field="code"
                    value={code}
                    autoFocus={!device.isTouch()}
                    onChange={setCode}
                    placeholder={i18n._({ id: "signin.code.placeholder", message: "Type in the code here" })}
                    maxLength={6}
                  />
                  <Button className="w-full justify-center" type="submit" variant="primary" disabled={code.length !== 6} onClick={verifyCode}>
                    <Trans id="signin.code.submit">Submit</Trans>
                  </Button>
                </Form>
                {resendMessage && <p className="text-green-700 mt-2">{resendMessage}</p>}
                <p className="text-center mt-2">
                  <a
                    href="#"
                    className="text-link"
                    onClick={(e) => {
                      e.preventDefault()
                      resendCode()
                    }}
                  >
                    <Trans id="signin.code.getnew">Get a new code</Trans>
                  </a>
                </p>
              </div>
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
