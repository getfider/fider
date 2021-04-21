import "./Hint.scss"

import React, { useState } from "react"
import { FaTimes } from "react-icons/fa"

import { cache } from "@fider/services"

interface HintProps {
  permanentCloseKey?: string
  condition?: boolean
}

export const Hint: React.FC<HintProps> = (props) => {
  const cacheKey: string | undefined = props.permanentCloseKey ? `Hint-Closed-${props.permanentCloseKey}` : undefined
  const [isClosed, setIsClosed] = useState<boolean>(cacheKey ? cache.local.has(cacheKey) : false)

  const close = () => {
    if (cacheKey) {
      cache.local.set(cacheKey, "true")
    }
    setIsClosed(true)
  }

  if (props.condition === false || isClosed) {
    return null
  }
  return (
    <p className="c-hint">
      <strong>HINT:</strong> {props.children}
      {cacheKey && <FaTimes onClick={close} className="close" />}
    </p>
  )
}
