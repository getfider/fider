import React from "react"
import { resolveHref } from "@fider/services"

type AnchorProps = React.AnchorHTMLAttributes<HTMLAnchorElement>

interface LinkProps extends Omit<AnchorProps, "href"> {
  href: string
  children?: React.ReactNode
}

// Link is a thin wrapper around <a> that runs the href through resolveHref()
// so callers can write natural root-relative paths (e.g. "/posts/1") without
// worrying about sub-path hosting. Use in place of <a> for any in-app
// navigation. External URLs, fragments, and already-prefixed hrefs pass
// through unchanged.
export const Link: React.FC<LinkProps> = ({ href, children, ...rest }) => {
  return (
    <a href={resolveHref(href)} {...rest}>
      {children}
    </a>
  )
}
