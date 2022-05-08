import { Trans } from "@lingui/macro"
import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const Error410 = () => {
  return (
    <ErrorPageWrapper id="p-error410" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.expired.title">Expired</Trans>
      </h1>
      <p>
        <Trans id="error.expired.text">The link you clicked has expired.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default Error410
