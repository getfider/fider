import { Trans } from "@lingui/react/macro"
import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const AccessDenied = () => {
  return (
    <ErrorPageWrapper id="p-access-denied" showHomeLink={true}>
      <h1 className="text-display">
        <Trans id="error.accessdenied.title">Access Denied</Trans>
      </h1>
      <p>
        <Trans id="error.accessdenied.text">You do not have the required permissions to access this site.</Trans>
      </p>
      <p className="text-muted">
        <Trans id="error.accessdenied.contact">If you believe this is an error, please contact your administrator.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default AccessDenied
