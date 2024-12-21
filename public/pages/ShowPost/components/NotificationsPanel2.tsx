import React, { useState } from "react"
import { Button, Icon } from "@fider/components"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconVolumeOn from "@fider/assets/images/heroicons-volume-on.svg"
import IconVolumeOff from "@fider/assets/images/heroicons-volume-off.svg"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"
import { NotificationsPanelProps } from "./NotificationsPanel"

export const NotificationsPanel2 = (props: NotificationsPanelProps) => {
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
    <Button variant="quaternary" className="w-full text-primary-base" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconVolumeOff} />{" "}
      <span>
        <Trans id="label.unfollow">Unfollow</Trans>
      </span>
    </Button>
  ) : (
    <Button variant="quaternary" className="w-full text-primary-base" onClick={subscribeOrUnsubscribe} disabled={fider.isReadOnly}>
      <Icon sprite={IconVolumeOn} />
      <span>
        <Trans id="label.follow">Follow</Trans>
      </span>
    </Button>
  )

  return <VStack>{button}</VStack>
}
