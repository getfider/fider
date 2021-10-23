import React from "react"
import { useFider } from "@fider/hooks"
import { Message } from "./common"

export const ReadOnlyNotice = () => {
  const fider = useFider()
  if (!fider.isReadOnly) {
    return null
  }

  if (fider.session.isAuthenticated && fider.session.user.isAdministrator) {
    return (
      <Message alignment="center" type="warning">
        This website is currently in read-only mode because there is no active subscription. Visit{" "}
        <a className="text-link" href="/billing">
          Billing
        </a>{" "}
        to subscribe.
      </Message>
    )
  }

  return (
    <Message alignment="center" type="warning">
      This website is currently in read-only mode.
    </Message>
  )
}
