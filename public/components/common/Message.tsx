import "./Message.scss"

import React from "react"
import { classSet } from "@fider/services"
import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"
import IconExclamationCircle from "@fider/assets/images/heroicons-exclamation-circle.svg"
import IconExclamation from "@fider/assets/images/heroicons-exclamation.svg"
import { HStack } from "@fider/components/layout"
import { Icon } from "./Icon"

interface MessageProps {
  type: "success" | "warning" | "error"
  className?: string
  showIcon?: boolean
}

export const Message: React.FunctionComponent<MessageProps> = (props) => {
  const className = classSet({
    "c-message": true,
    [`c-message--${props.type}`]: true,
    [`${props.className}`]: props.className,
  })

  const icon = props.type === "error" ? IconExclamation : props.type === "warning" ? IconExclamationCircle : IconCheckCircle

  return (
    <HStack className={className} spacing={2}>
      {props.showIcon === true && <Icon className="h-5" sprite={icon} />}
      <span>{props.children}</span>
    </HStack>
  )
}
