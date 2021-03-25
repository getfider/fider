import React from "react"
import { useFider } from "@fider/hooks"

export const FiderVersion = () => {
  const fider = useFider()

  return (
    <p className="info center hidden-sm hidden-md">
      Support our{" "}
      <a rel="noopener" target="_blank" href="http://opencollective.com/fider">
        OpenCollective
      </a>
      <br />
      Fider v{fider.settings.version}
    </p>
  )
}
