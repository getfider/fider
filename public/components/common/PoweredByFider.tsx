import React from "react"
import { classSet, FiderContext } from "@fider/services"

import "./PoweredByFider.scss"

interface PoweredByFiderProps {
  slot: string
  className?: string
}

export const PoweredByFider = (props: PoweredByFiderProps) => {
  const fider = React.useContext(FiderContext)
  const source = encodeURIComponent(window?.location?.host || "")
  const medium = "powered-by"
  const campaign = props.slot
  const version = fider.settings?.version
  const versionString = fider.isSingleHostMode() && version && version !== "dev" ? `v${version}` : ""

  const className = classSet({
    "c-powered": true,
    [props.className || ""]: props.className,
  })

  return (
    <div className={className}>
      <a rel="noopener" className="text-2xs" href={`https://fider.io?utm_source=${source}&utm_medium=${medium}&utm_campaign=${campaign}`} target="_blank">
        Powered by Fider ⚡
      </a>
      {versionString && <span className="text-2xs block">{versionString}</span>}
    </div>
  )
}
