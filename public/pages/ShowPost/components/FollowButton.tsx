import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconPlus from "@fider/assets/images/heroicons-plus.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"
import { Post } from "@fider/models"

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

  const button = subscribed ? (
    <Button className="w-full text-gray-800 bg-white border border-gray-800 no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconCheck} />{" "}
      <span>
        <Trans id="label.following">Following</Trans>
      </span>
    </Button>
  ) : (
    <Button className="w-full text-blue-500 bg-white border border-blue-500 no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconPlus} />
      <span>
        <Trans id="label.follow">Follow</Trans>
      </span>
    </Button>
  )

  return <VStack>{button}</VStack>
}
