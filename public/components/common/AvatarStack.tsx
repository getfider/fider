import "./AvatarStack.scss"

import React from "react"
import { UserRole } from "@fider/models"
import { Avatar } from "./Avatar"
import { classSet } from "@fider/services"

interface AvatarStackProps {
  overlap?: boolean
  users: Array<{
    role?: UserRole
    avatarURL: string
    name: string
  }>
}

export const AvatarStack = (props: AvatarStackProps) => {
  const classes = classSet({
    "c-avatar-stack": true,
    "c-avatar-stack--overlap": props.overlap ?? true,
  })

  return (
    <div className={classes}>
      {props.users.map((x, i) => (
        <Avatar key={i} user={x} />
      ))}
    </div>
  )
}
