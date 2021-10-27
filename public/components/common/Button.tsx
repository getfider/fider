import "./Button.scss"

import React, { useEffect, useRef, useState } from "react"
import { classSet } from "@fider/services"

interface ButtonProps {
  className?: string
  disabled?: boolean
  href?: string
  rel?: "nofollow"
  target?: "_self" | "_blank" | "_parent" | "_top"
  type?: "button" | "submit"
  variant?: "primary" | "danger" | "secondary" | "tertiary"
  size?: "small" | "default" | "large"
  onClick?: (event: ButtonClickEvent) => Promise<any> | void
}

export class ButtonClickEvent {
  private shouldEnable = true
  public preventEnable(): void {
    this.shouldEnable = false
  }
  public canEnable(): boolean {
    return this.shouldEnable
  }
}

export const Button: React.FC<ButtonProps> = (props) => {
  const [clicked, setClicked] = useState(false)
  const unmountedContainer = useRef(false)

  useEffect(() => {
    return () => {
      unmountedContainer.current = true
    }
  }, [])

  const className = classSet({
    "c-button": true,
    [`c-button--${props.size}`]: props.size,
    [`c-button--${props.variant}`]: props.variant,
    "c-button--loading": clicked,
    "c-button--disabled": clicked || props.disabled,
    [props.className || ""]: props.className,
    "shadow-sm": props.variant !== "tertiary",
  })

  const click = async (e?: React.SyntheticEvent<HTMLElement>) => {
    if (e) {
      e.preventDefault()
      e.stopPropagation()
    }

    if (clicked) {
      return
    }

    const event = new ButtonClickEvent()
    setClicked(true)

    if (props.onClick) {
      await props.onClick(event)

      if (!unmountedContainer.current && event.canEnable()) {
        setClicked(false)
      }
    }
  }

  const buttonContent = props.href ? (
    <a href={props.href} rel={props.rel} target={props.target} className={className}>
      {props.children}
    </a>
  ) : (
    <button type={props.type} className={className} onClick={click}>
      {props.children}
    </button>
  )

  return buttonContent
}

Button.defaultProps = {
  size: "default",
  variant: "secondary",
  type: "button",
}
