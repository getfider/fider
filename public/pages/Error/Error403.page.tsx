import React from "react"
import { Trans } from "@lingui/macro"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const Error403 = () => {
  return (
    <ErrorPageWrapper id="p-error403" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.forbidden.title">Forbidden</Trans>
      </h1>
      <p>
        <Trans id="error.forbidden.text">You are not authorized to view this page.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default Error403
