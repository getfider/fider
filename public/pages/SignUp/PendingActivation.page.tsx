import React, { useState } from "react"
import { TenantLogo, Button, Message } from "@fider/components"
import { Trans } from "@lingui/react/macro"
import { actions, Failure } from "@fider/services"

const PendingActivation = () => {
  const [isResending, setIsResending] = useState(false)
  const [successMessage, setSuccessMessage] = useState("")
  const [error, setError] = useState<Failure | undefined>()

  const resendEmail = async () => {
    setIsResending(true)
    setError(undefined)
    setSuccessMessage("")

    const result = await actions.resendSignUpEmail()

    setIsResending(false)

    if (result.ok) {
      setSuccessMessage("Verification email sent successfully! Please check your inbox.")
    } else if (result.error) {
      setError(result.error)
    }
  }

  return (
    <div id="p-notinvited" className="container page">
      <div className="w-max-7xl mx-auto text-center mt-8">
        <div className="h-20 mb-4">
          <TenantLogo size={100} useFiderIfEmpty={true} />
        </div>
        <h1 className="text-display uppercase">
          <Trans id="page.pendingactivation.title">Your account is pending activation</Trans>
        </h1>
        <p>
          <Trans id="page.pendingactivation.text">We sent you a confirmation email with a link to activate your site.</Trans>
        </p>
        <p>
          <Trans id="page.pendingactivation.text2">Please check your inbox to activate it.</Trans>
        </p>

        <div className="mt-6">
          {successMessage && (
            <Message className="mb-4" type="success" showIcon={true}>
              {successMessage}
            </Message>
          )}

          {error && error.errors && error.errors.length > 0 && (
            <Message className="mb-4" type="error" showIcon={true}>
              {error.errors.map((e, i) => (
                <div key={i}>{e.message}</div>
              ))}
            </Message>
          )}

          <p className="text-muted mb-2">
            <Trans id="page.pendingactivation.didntreceive">Didn&apos;t receive the email?</Trans>
          </p>

          <Button variant="primary" onClick={resendEmail} disabled={isResending}>
            {isResending ? (
              <Trans id="page.pendingactivation.resending">Resending...</Trans>
            ) : (
              <Trans id="page.pendingactivation.resend">Resend verification email</Trans>
            )}
          </Button>
        </div>
      </div>
    </div>
  )
}

export default PendingActivation
