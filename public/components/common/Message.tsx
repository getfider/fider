import "./Message.scss"

import React from "react"
import { classSet } from "@fider/services"
import { FaBan, FaRegCheckCircle, FaExclamationTriangle } from "react-icons/fa"

interface MessageProps {
  type: "success" | "warning" | "error"
  showIcon?: boolean
}

export const Message: React.FunctionComponent<MessageProps> = (props) => {
  const className = classSet({
    "c-message": true,
    [`m-${props.type}`]: true,
  })

  const icon = props.type === "error" ? <FaBan /> : props.type === "warning" ? <FaExclamationTriangle /> : <FaRegCheckCircle />

  return (
    <p className={className}>
      {props.showIcon === true && icon}
      <span>{props.children}</span>
    </p>
  )
}
