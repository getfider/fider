import "./Error.page.scss"

import React from "react"
import { TenantLogo } from "@fider/components"
import { useFider } from "@fider/hooks"

interface ErrorPageProps {
  error: Error
  errorInfo: React.ErrorInfo
  showDetails?: boolean
}

export const ErrorPage = (props: ErrorPageProps) => {
  const fider = useFider()

  return (
    <div id="p-error" className="container failure-page">
      <TenantLogo size={100} useFiderIfEmpty={true} />
      <h1>Shoot! Well, this is unexpected…</h1>
      <p>An error has occurred and we&apos;re working to fix the problem!</p>
      {fider.settings && (
        <span>
          Take me back to <a href={fider.settings.baseURL}>{fider.settings.baseURL}</a> home page.
        </span>
      )}
      {props.showDetails && (
        <pre className="error">
          {props.error.toString()}
          {props.errorInfo.componentStack}
        </pre>
      )}
    </div>
  )
}
