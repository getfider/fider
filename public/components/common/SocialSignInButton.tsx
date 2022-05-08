import React from "react"
import { Button, OAuthProviderLogo } from "@fider/components"

interface SocialSignInButtonProps {
  option: {
    displayName: string
    provider?: string
    url?: string
    logoBlobKey?: string
    logoURL?: string
  }
  className?: string
  redirectTo?: string
}

export const SocialSignInButton = (props: SocialSignInButtonProps) => {
  const redirectTo = props.redirectTo || window.location.href
  const href = props.option.url ? `${props.option.url}?redirect=${redirectTo}` : undefined

  return (
    <Button href={href} rel="nofollow" className={props.className}>
      {props.option.logoURL ? <img alt={props.option.displayName} src={props.option.logoURL} /> : <OAuthProviderLogo option={props.option} />}
      <span>{props.option.displayName}</span>
    </Button>
  )
}
