import "./AvatarStack.scss"

import React from "react"
import { UserRole } from "@fider/models"
import { Avatar } from "./Avatar"

interface AvatarStackProps {
  users: Array<{
    role?: UserRole
    avatarURL: string
    name: string
  }>
}

export const AvatarStack = (props: AvatarStackProps) => {
  return (
    <div className="c-avatar-stack">
      {props.users.map((x, i) => (
        <Avatar key={i} user={x} />
      ))}
    </div>
  )
}
