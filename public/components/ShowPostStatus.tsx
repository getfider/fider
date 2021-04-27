import "./ShowPostStatus.scss"

import React from "react"
import { PostStatus } from "@fider/models"

interface ShowPostStatusProps {
  status: PostStatus
}

export const ShowPostStatus = (props: ShowPostStatusProps) => {
  return <span className={`c-status-label c-status-label--${props.status.value}`}>{props.status.title}</span>
}
