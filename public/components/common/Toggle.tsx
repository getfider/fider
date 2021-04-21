import "./Toggle.scss"

import React, { useState } from "react"
import { classSet } from "@fider/services"

interface ToggleProps {
  label?: string
  active: boolean
  disabled?: boolean
  onToggle?: (active: boolean) => void
}

export const Toggle: React.StatelessComponent<ToggleProps> = (props) => {
  const [active, setActive] = useState(props.active)

  const toggle = () => {
    if (props.disabled) {
      return
    }

    const newActive = !active
    setActive(newActive)
    if (props.onToggle) {
      props.onToggle(newActive)
    }
  }

  const className = classSet({
    "c-toggle": true,
    "m-disabled": !!props.disabled,
  })

  return (
    <span className={className} onClick={toggle}>
      <input type="checkbox" checked={active} readOnly={true} />
      <label>
        <span className="switch" />
      </label>
      <span className="text">{!!props.label && props.label}</span>
    </span>
  )
}
