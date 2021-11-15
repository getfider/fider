import React from "react"
import { TenantLogo } from "@fider/components"
import { Trans } from "@lingui/macro"

const PendingActivation = () => {
  return (
    <div id="p-notinvited" className="container page">
      <div className="w-max-7xl mx-auto text-center mt-8">
        <div className="h-20 mb-4">
          <TenantLogo size={100} useFiderIfEmpty={true} />
        </div>
        <h1 className="text-display uppercase">
          <Trans id="page.pendingactivation.title">Your account is pending activation</Trans>
        </h1>
        <p>
          <Trans id="page.pendingactivation.text">We sent you a confirmation email with a link to activate your site.</Trans>
        </p>
        <p>
          <Trans id="page.pendingactivation.text2">Please check your inbox to activate it.</Trans>
        </p>
      </div>
    </div>
  )
}

export default PendingActivation
