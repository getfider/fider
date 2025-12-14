import React, { useState } from "react"
import { SignInModal, RSSModal, TenantLogo, NotificationIndicator, UserMenu, ThemeSwitcher, Icon, Button, ModerationIndicator } from "@fider/components"
import { useFider } from "@fider/hooks"
import { HStack } from "./layout"
import { Trans } from "@lingui/react/macro"
import { i18n } from "@lingui/core"
import IconRss from "@fider/assets/images/heroicons-rss.svg"
import "./Header.scss"

interface HeaderProps {
  hasInert?: boolean
}

export const Header = (props: HeaderProps) => {
  const fider = useFider()
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [isRSSModalOpen, setIsRSSModalOpen] = useState(false)

  const handleSignInClick = () => {
    setIsSignInModalOpen(true)
  }

  const showRSSModal = (e: React.MouseEvent) => {
    e.preventDefault()
    setIsRSSModalOpen(true)
  }

  const atomFeedTitle = i18n._({ id: "action.postsfeed", message: "Posts Feed" })
  const hideSignInModal = () => setIsSignInModalOpen(false)
  const hideRSSModal = () => setIsRSSModalOpen(false)

  return (
    <div id="c-header" className="bg-white" style={{ borderBottom: "1px solid var(--colors-gray-200)" }} {...(props.hasInert && { inert: "true" })}>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideSignInModal} />
      <RSSModal isOpen={isRSSModalOpen} onClose={hideRSSModal} url={`${fider.settings.baseURL}/feed/global.atom`} />
      <HStack className="c-menu p-4 w-full">
        <div className="container c-header__container">
          <div className="flex flex-wrap flex-items-center gap-2">
            <div className="flex flex-x flex-items-center justify-between w-full">
              <a href="/" className="flex flex-x flex-items-center flex--spacing-2 h-8">
                <TenantLogo size={100} />
                <h1 className="text-header">{fider.session.tenant.name}</h1>
              </a>
              {fider.session.isAuthenticated && (
                <HStack spacing={2}>
                  {fider.session.tenant.isFeedEnabled && (
                    <button title={atomFeedTitle} className="c-themeswitcher" onClick={showRSSModal}>
                      <Icon sprite={IconRss} className="h-6 text-gray-500" />
                    </button>
                  )}
                  <ThemeSwitcher />
                  <NotificationIndicator />
                  <UserMenu />
                </HStack>
              )}
              {!fider.session.isAuthenticated && (
                <HStack spacing={2}>
                  {fider.session.tenant.isFeedEnabled && (
                    <button title={atomFeedTitle} className="c-themeswitcher" onClick={showRSSModal}>
                      <Icon sprite={IconRss} className="h-6 text-gray-500" />
                    </button>
                  )}
                  <ThemeSwitcher />
                  <Button variant="primary" size="default" onClick={handleSignInClick}>
                    <HStack spacing={1} className="flex-items-center">
                      <span>
                        <Trans id="action.signin">Sign in</Trans>
                      </span>
                    </HStack>
                  </Button>
                </HStack>
              )}
            </div>
            {fider.session.isAuthenticated && (
              <div className="c-header__moderation">
                <ModerationIndicator />
              </div>
            )}
          </div>
        </div>
      </HStack>
    </div>
  )
}
