import "./Avatar.scss"

import React from "react"
import { classSet } from "@fider/services"
import { isCollaborator, UserRole } from "@fider/models"

interface AvatarProps {
  user: {
    role?: UserRole
    avatarURL: string
    name: string
  }
  size?: "small" | "normal" | "large"
}

export const Avatar = (props: AvatarProps) => {
  const size = props.size || "normal"

  const className = classSet({
    "c-avatar": true,
    [`m-${size}`]: true,
    "m-staff": props.user.role && isCollaborator(props.user.role),
  })

  return <img className={className} alt={props.user.name} title={props.user.name} src={`${props.user.avatarURL}?size=50`} />
}
