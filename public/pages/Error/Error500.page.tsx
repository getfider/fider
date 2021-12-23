import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"
import { Trans } from "@lingui/macro"

const Error500 = () => {
  return (
    <ErrorPageWrapper id="p-error500" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.internalerror.title">Shoot! Well, this is unexpected…</Trans>
      </h1>
      <p>
        <Trans id="error.internalerror.text">An error has occurred and we&apos;re working to fix the problem! We’ll be up and running shortly.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default Error500
