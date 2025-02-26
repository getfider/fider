// import "./ShowPostStatus.scss"

import React from "react"
import { PostStatus } from "@fider/models"
import { i18n } from "@lingui/core"

interface ShowPostStatusProps {
  status: PostStatus
}

export const ShowPostStatus = (props: ShowPostStatusProps) => {
  const id = `enum.poststatus.${props.status.value}`
  const title = i18n._(id, { message: props.status.title })

  return <span className={`c-status-label c-status-label--${props.status.value}`}>{title}</span>
}
