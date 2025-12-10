import React, { useState } from "react"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconPlus from "@fider/assets/images/heroicons-plus.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { Trans } from "@lingui/react/macro"
import { Post } from "@fider/models"
import { ActionButton } from "./ActionButton"

export interface NotificationsPanelProps {
  post: Post
  subscribed: boolean
}

export const FollowButton = (props: NotificationsPanelProps) => {
  const fider = useFider()
  const [subscribed, setSubscribed] = useState(props.subscribed)

  const subscribeOrUnsubscribe = async () => {
    const action = subscribed ? actions.unsubscribe : actions.subscribe

    const response = await action(props.post.number)
    if (response.ok) {
      setSubscribed(!subscribed)
    }
  }

  if (!fider.session.isAuthenticated) {
    return null
  }

  return subscribed ? (
    <ActionButton icon={IconCheck} onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Trans id="label.following">Following</Trans>
    </ActionButton>
  ) : (
    <ActionButton icon={IconPlus} onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Trans id="label.follow">Follow</Trans>
    </ActionButton>
  )
}
