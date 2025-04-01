import React from "react"
import { OAuthProviderLogo } from "@fider/components"
import { Trans } from "@lingui/react/macro"

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
    <a rel="nofollow" className="c-signin-social-button" href={href}>
      {props.option.logoURL ? <img alt={props.option.displayName} src={props.option.logoURL} /> : <OAuthProviderLogo option={props.option} />}
      <span>
        <Trans id="signin.message.socialbutton.intro">Log in with</Trans>
        &nbsp;
        {props.option.displayName}
      </span>
    </a>
  )
}
