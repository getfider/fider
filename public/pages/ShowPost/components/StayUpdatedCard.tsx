import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconPlus from "@fider/assets/images/heroicons-plus.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { Post } from "@fider/models"

export interface StayUpdatedCardProps {
  post: Post
  subscribed: boolean
}

export const StayUpdatedCard = (props: StayUpdatedCardProps) => {
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
    <Button className="p-show-post__action-col__follow-button--following no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconCheck} />{" "}
      <span>
        <Trans id="label.following">Following</Trans>
      </span>
    </Button>
  ) : (
    <Button className="p-show-post__action-col__follow-button--follow no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconPlus} />
      <span>
        <Trans id="label.follow">Follow</Trans>
      </span>
    </Button>
  )

  return (
    <VStack spacing={2} className="card">
      <h3 className="text-bold text-gray-900">
        <Trans id="showpost.stayupdated.title">Stay Updated</Trans>
      </h3>
      <p className="text-sm text-gray-600">
        <Trans id="showpost.stayupdated.description">Get notified about updates to this idea</Trans>
      </p>
      {button}
    </VStack>
  )
}
