import React from "react"
import { Header } from "@fider/components"
import { useFider } from "@fider/hooks"

interface AppShellProps {
  page: string
  children: React.ReactNode
}

// Pages that render without the site header
const noHeader = /^(SignIn|SignUp|Legal|OAuthEcho)\//

// Persistent frame around the current page: survives client-side navigation
export const AppShell = (props: AppShellProps) => {
  const fider = useFider()
  const showHeader = !!fider.session.tenant && !noHeader.test(props.page)
  return (
    <>
      {showHeader && <Header />}
      {props.children}
    </>
  )
}
