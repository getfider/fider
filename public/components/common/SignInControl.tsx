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
  onCodeVerified?: () => void
}

enum EmailSigninStep {
  EnterEmail,
  EnterName,
  EnterCode,
}

export const SignInControl: React.FunctionComponent<SignInControlProps> = (props) => {
  const fider = useFider()
  const [showEmailForm, setShowEmailForm] = useState(fider.session.tenant ? fider.session.tenant.isEmailAuthAllowed : true)
  const [email, setEmail] = useState("")
  const [emailSignInStep, setEmailSignInStep] = useState(EmailSigninStep.EnterEmail)
  const [userName, setUserName] = useState("")
  const [code, setCode] = useState("")
  const [error, setError] = useState<Failure | undefined>(undefined)
  const [resendMessage, setResendMessage] = useState("")

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
    setEmailSignInStep(EmailSigninStep.EnterEmail)
    setUserName("")
    setCode("")
    setError(undefined)
    setResendMessage("")
  }

  const signIn = async () => {
    await doPreSigninAction()
    const result = await actions.signIn(email)
    if (result.ok) {
      setError(undefined)
      const data = result.data as { userExists?: boolean } | undefined
      if (data && data.userExists === false) {
        // New user - show name field
        setEmailSignInStep(EmailSigninStep.EnterName)
      } else {
        // Existing user - show code entry
        setEmailSignInStep(EmailSigninStep.EnterCode)
      }
    } else if (result.error) {
      setError(result.error)
    }
  }

  const submitNewUser = async () => {
    doPreSigninAction()
    const result = await actions.signInNewUser(email, userName)
    if (result.ok) {
      setError(undefined)
      setEmailSignInStep(EmailSigninStep.EnterCode)
    } else if (result.error) {
      setError(result.error)
    }
  }

  const verifyCode = async () => {
    const result = await actions.verifySignInCode(email, code)
    if (result.ok) {
      if (props.onCodeVerified) {
        // Let the parent component decide what to do
        props.onCodeVerified()
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

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    // Form submission handler that routes to the correct function based on step
    if (emailSignInStep === EmailSigninStep.EnterEmail) {
      await signIn()
    } else if (emailSignInStep === EmailSigninStep.EnterName) {
      await submitNewUser()
    } else if (emailSignInStep === EmailSigninStep.EnterCode) {
      await verifyCode()
    }
  }

  const renderSigninEmailButton = () => {
    if (emailSignInStep == EmailSigninStep.EnterEmail) {
      return (
        <Button className="w-full justify-center" type="submit" variant="primary">
          <Trans id="signin.message.email">Continue with Email</Trans>
        </Button>
      )
    }
    if (emailSignInStep == EmailSigninStep.EnterName) {
      return (
        <Button className="w-full justify-center" type="submit" variant="primary">
          <Trans id="action.signup">Sign up</Trans>
        </Button>
      )
    }
    if (emailSignInStep == EmailSigninStep.EnterCode) {
      return (
        <Button className="w-full justify-center" type="submit" variant="primary" disabled={code.length !== 6}>
          <Trans id="action.submit">Submit</Trans>
        </Button>
      )
    }
  }

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
            <Form error={error} autoComplete={emailSignInStep === EmailSigninStep.EnterCode ? "on" : "off"} onSubmit={handleFormSubmit}>
              {(emailSignInStep == EmailSigninStep.EnterEmail || emailSignInStep == EmailSigninStep.EnterName) && renderEmailField()}

              {emailSignInStep == EmailSigninStep.EnterName && renderNameField()}

              {emailSignInStep == EmailSigninStep.EnterCode && renderCodeField()}

              <div className="pt-3">{renderSigninEmailButton()}</div>
            </Form>
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

  function renderNameField() {
    return (
      <Input
        className="text-left"
        field="name"
        value={userName}
        autoFocus={!device.isTouch()}
        onChange={setUserName}
        placeholder={i18n._({ id: "signin.name.placeholder", message: "Your name" })}
        maxLength={100}
      />
    )
  }

  function renderEmailField(): React.ReactNode {
    return (
      <>
        <Input
          className="text-left"
          field="email"
          value={email}
          disabled={emailSignInStep == EmailSigninStep.EnterName}
          autoFocus={!device.isTouch()}
          autoComplete="email"
          onChange={setEmail}
          placeholder={i18n._({ id: "signin.email.placeholder", message: "Email address" })}
        />
        {!fider.session.tenant.isEmailAuthAllowed && (
          <p className="text-red-700 mt-1">
            <Trans id="signin.message.onlyadmins">Currently only allowed to sign in to an administrator account</Trans>
          </p>
        )}
      </>
    )
  }

  function renderCodeField(): React.ReactNode {
    return (
      <>
        <p className="text-muted mb-2 text-left">
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
        <Input
          className="text-left"
          field="code"
          value={code}
          autoFocus={!device.isTouch()}
          autoComplete="one-time-code"
          inputMode="numeric"
          onChange={setCode}
          placeholder={i18n._({ id: "signin.code.placeholder", message: "Type in the code here" })}
          maxLength={6}
        />
        {resendMessage && <p className="text-green-700 mt-2">{resendMessage}</p>}
        <p className="text-center mt-2 text-muted text-left">
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
      </>
    )
  }
}
