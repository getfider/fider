import { Trans } from "@lingui/macro"
import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const NotInvited = () => {
  return (
    <ErrorPageWrapper id="p-notinvited" showHomeLink={true}>
      <h1 className="text-display">
        <Trans id="error.unauthorized.title">Not invited</Trans>
      </h1>
      <p>
        <Trans id="error.unauthorized.text">We could not find an account for your email address.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default NotInvited
