import React from "react"
import { useFider } from "@fider/hooks"

export const FiderVersion = () => {
  const fider = useFider()

  return (
    <p className="text-muted mt-2 text-center sm:hidden md:hidden">
      Support our{" "}
      <a className="text-link" rel="noopener" target="_blank" href="http://opencollective.com/fider">
        OpenCollective
      </a>
      <br />
      Fider v{fider.settings.version}
    </p>
  )
}
