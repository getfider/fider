import { Trans } from "@lingui/macro"
import React from "react"
import { ErrorPageWrapper } from "./components/ErrorPageWrapper"

const Error404 = () => {
  return (
    <ErrorPageWrapper id="p-error404" showHomeLink={true}>
      <h1 className="text-display uppercase">
        <Trans id="error.pagenotfound.title">Page not found</Trans>
      </h1>
      <p>
        <Trans id="error.pagenotfound.text">The link you clicked may be broken or the page may have been removed.</Trans>
      </p>
    </ErrorPageWrapper>
  )
}

export default Error404
