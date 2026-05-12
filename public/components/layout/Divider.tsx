import "./Divider.scss"

import React from "react"
import { Trans } from "@lingui/react/macro"

export const Divider = () => {
  return (
    <div className="c-divider text-gray-600">
      <Trans id="label.or">OR</Trans>
    </div>
  )
}
