import "./Loader.scss"

import React, { useState } from "react"
import { useTimeout } from "@fider/hooks"

export function Loader() {
  const [show, setShow] = useState(false)

  useTimeout(() => {
    setShow(true)
  }, 500)

  return show ? <div className="c-loader" /> : null
}
