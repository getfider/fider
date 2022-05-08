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
    <div id="p-error" className="container page">
      <div className="w-max-7xl mx-auto">
        <div className="h-20 text-center mb-4">
          <TenantLogo size={100} useFiderIfEmpty={true} />
        </div>
        <h1 className="text-display">Shoot! Well, this is unexpected…</h1>
        <p>An error has occurred and we&apos;re working to fix the problem!</p>
        {fider.settings && (
          <span>
            Take me back to{" "}
            <a className="text-link" href={fider.settings.baseURL}>
              {fider.settings.baseURL}
            </a>{" "}
            home page.
          </span>
        )}
      </div>
      {props.showDetails && (
        <pre>
          {props.error.toString()}
          {props.errorInfo.componentStack}
        </pre>
      )}
    </div>
  )
}
