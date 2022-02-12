import React from "react"

import "./PoweredByFider.scss"

interface PoweredByFiderProps {
  className: string
}

export const PoweredByFider = (props: PoweredByFiderProps) => {
  return (
    <div className={`c-powered ${props.className}`}>
      <a rel="noopener" href="https://fider.io" target="_blank">
        Powered by Fider
      </a>
    </div>
  )
}
