import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconPlus from "@fider/assets/images/heroicons-plus.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"
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
    <Button className="w-full text-gray-800 bg-white border border-gray-800 no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconCheck} />{" "}
      <span>
        <Trans id="label.following">Following</Trans>
      </span>
    </Button>
  ) : (
    <Button className="w-full text-blue-600 bg-white border border-blue-600 no-focus" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconPlus} />
      <span>
        <Trans id="label.follow">Follow</Trans>
      </span>
    </Button>
  )

  return (
    <VStack spacing={2}>
      <h3 className="text-base font-semibold text-gray-900">
        <Trans id="showpost.stayupdated.title">Stay Updated</Trans>
      </h3>
      <p className="text-sm text-gray-600">
        <Trans id="showpost.stayupdated.description">Get notified when this post receives updates</Trans>
      </p>
      {button}
    </VStack>
  )
}
