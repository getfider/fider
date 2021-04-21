import "./List.scss"

import React from "react"
import { classSet } from "@fider/services"

interface ListProps {
  className?: string
  divided?: boolean
  hover?: boolean
}

interface ListItemProps {
  className?: string
  onClick?: () => void
}

export const List: React.FunctionComponent<ListProps> = (props) => {
  const className = classSet({
    "c-list": true,
    [props.className || ""]: true,
    "m-divided": props.divided,
    "m-hover": props.hover,
  })

  return <div className={className}>{props.children}</div>
}

export const ListItem: React.FunctionComponent<ListItemProps> = (props) => {
  const className = classSet({
    "c-list-item": true,
    [props.className || ""]: true,
    "m-selectable": props.onClick,
  })

  if (props.onClick) {
    return (
      <div className={className} onClick={props.onClick}>
        {props.children}
      </div>
    )
  }
  return <div className={className}>{props.children}</div>
}
