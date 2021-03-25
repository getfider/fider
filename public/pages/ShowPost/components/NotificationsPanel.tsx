import React, { useState } from "react"
import { Post } from "@fider/models"
import { Button, List, ListItem } from "@fider/components"
import { actions } from "@fider/services"
import { FaVolumeUp, FaVolumeMute } from "react-icons/fa"
import { useFider } from "@fider/hooks"

interface NotificationsPanelProps {
  post: Post
  subscribed: boolean
}

export const NotificationsPanel = (props: NotificationsPanelProps) => {
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
    <Button fluid={true} onClick={subscribeOrUnsubscribe}>
      <FaVolumeMute /> Unsubscribe
    </Button>
  ) : (
    <Button fluid={true} onClick={subscribeOrUnsubscribe}>
      <FaVolumeUp /> Subscribe
    </Button>
  )

  const text = subscribed ? (
    <span className="info">You’re receiving notifications about activity on this post.</span>
  ) : (
    <span className="info">You&apos;ll not receive any notification about this post.</span>
  )

  return (
    <>
      <span className="subtitle">Notifications</span>
      <List>
        <ListItem>
          {button}
          {text}
        </ListItem>
      </List>
    </>
  )
}
