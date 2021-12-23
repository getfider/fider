import React from "react"
import { Trans } from "@lingui/macro"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const Error401 = () => {
  return (
    <ErrorPageWrapper id="p-error401" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.unauthorized.title">Unauthorized</Trans>
      </h1>
      <p>
        <Trans id="error.unauthorized.text">You need to sign in before accessing this page.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default Error401
