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
  onClick?: () => void
}

export const SocialSignInButton = (props: SocialSignInButtonProps) => {
  const redirectTo = props.redirectTo || window.location.href
  const href = props.option.url ? `${props.option.url}?redirect=${redirectTo}` : ""

  const handleClick = async (e: React.MouseEvent<HTMLAnchorElement>) => {
    // If there is an onClick then let that run and check it finishes OK before doing the oauth.
    if (props.onClick) {
      e.preventDefault()
      props.onClick()
      window.location.href = href
    }
  }

  return (
    <a rel="nofollow" className="c-signin-social-button" href={href} onClick={handleClick}>
      {props.option.logoURL ? <img alt={props.option.displayName} src={props.option.logoURL} /> : <OAuthProviderLogo option={props.option} />}
      <span>
        <Trans id="signin.message.socialbutton.intro">Continue with</Trans>
        &nbsp;
        {props.option.displayName}
      </span>
    </a>
  )
}
