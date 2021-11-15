import React from "react"
import { Header, TenantLogo } from "@fider/components"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/macro"

interface ErrorPageWrapperProps {
  id: string
  showHomeLink: boolean
  children?: React.ReactNode
}

export const ErrorPageWrapper = (props: ErrorPageWrapperProps) => {
  const fider = useFider()

  return (
    <>
      {fider.session.tenant && <Header />}
      <div id={props.id} className="container page">
        <div className="w-max-7xl mx-auto text-center mt-8">
          <div className="h-20 mb-4">
            <TenantLogo size={100} useFiderIfEmpty={true} />
          </div>
          {props.children}
          {props.showHomeLink && fider.session.tenant && (
            <p>
              <Trans id="page.backhome">
                Take me back to{" "}
                <a className="text-link" href={fider.settings.baseURL}>
                  {fider.settings.baseURL}
                </a>{" "}
                home page.
              </Trans>
            </p>
          )}
        </div>
      </div>
    </>
  )
}
