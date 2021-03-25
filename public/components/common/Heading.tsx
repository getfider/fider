import "./Heading.scss"

import React from "react"
import { classSet } from "@fider/services"
import { IconType } from "react-icons"

interface HeadingLogo {
  title: string
  dividing?: boolean
  size?: "normal" | "small"
  icon?: IconType
  subtitle?: string
  className?: string
}

const Header: React.FunctionComponent<{ level: number; className: string }> = (props) =>
  React.createElement(`h${props.level}`, { className: props.className }, props.children)

export const Heading = (props: HeadingLogo) => {
  const size = props.size || "normal"
  const level = size === "normal" ? 2 : 3
  const className = classSet({
    "c-heading": true,
    "m-dividing": props.dividing || false,
    [`m-${size}`]: true,
    [`${props.className}`]: props.className,
  })

  const iconClassName = classSet({
    "c-heading-icon": true,
    circular: level <= 2,
  })

  const icon = props.icon && <div className={iconClassName}>{React.createElement(props.icon)}</div>

  return (
    <Header level={level} className={className}>
      {icon}
      <div className="c-heading-content">
        {props.title}
        <div className="c-heading-subtitle">{props.subtitle}</div>
      </div>
    </Header>
  )
}
