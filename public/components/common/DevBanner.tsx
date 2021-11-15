import React from "react"
import { useFider } from "@fider/hooks"

import "./DevBanner.scss"

export const DevBanner = () => {
  const fider = useFider()

  if (fider.isProduction()) {
    return null
  }

  return <div className="c-dev-banner">DEV</div>
}
