import "./Hint.scss"

import React, { useState } from "react"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { HStack } from "@fider/components/layout"
import { cache } from "@fider/services"
import { Icon } from "./Icon"

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
    <HStack className="c-hint" justify="between" spacing={2}>
      <span>{props.children}</span>
      {cacheKey && <Icon sprite={IconX} onClick={close} className="c-hint__close h-5" />}
    </HStack>
  )
}
