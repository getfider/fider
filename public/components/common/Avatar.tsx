import "./Avatar.scss"

import React from "react"
import { UserRole } from "@fider/models"

interface AvatarProps {
  user: {
    role?: UserRole
    avatarURL: string
    name: string
  }
  size?: "small" | "normal"
}

export const Avatar = (props: AvatarProps) => {
  const size = props.size === "small" ? "h-6 w-6" : "h-8 w-8"
  return <img className={`c-avatar ${size}`} alt={props.user.name} src={`${props.user.avatarURL}?size=50`} />
}
