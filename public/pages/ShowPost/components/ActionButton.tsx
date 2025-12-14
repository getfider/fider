import "./ActionButton.scss"

import React from "react"
import { Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"

interface ActionButtonProps {
  icon: SpriteSymbol | string
  onClick: () => void
  children: React.ReactNode
  disabled?: boolean
  variant?: "normal" | "danger"
}

export const ActionButton = (props: ActionButtonProps) => {
  const { icon, onClick, children, variant = "normal" } = props
  const className = `c-action-button ${variant === "danger" ? "c-action-button--danger" : ""}`

  return (
    <button className={className} onClick={onClick} disabled={props.disabled}>
      <HStack spacing={2} align="center">
        <Icon sprite={icon} className="c-action-button__icon" />
        <span className="c-action-button__text">{children}</span>
      </HStack>
    </button>
  )
}
