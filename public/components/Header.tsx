import React, { useState } from "react"
import { SignInModal, RSSModal, TenantLogo, NotificationIndicator, UserMenu, ThemeSwitcher, Icon, Button, ModerationIndicator } from "@fider/components"
import { useFider } from "@fider/hooks"
import { Link } from "./common"
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

  const pathname = typeof window !== "undefined" ? window.location.pathname : "/"
  const isRoadmapActive = pathname === "/roadmap"
  const isFeedbackActive = !isRoadmapActive

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
          <div className="c-header__row">
            <Link href="/" className="c-header__brand flex flex-x flex-items-center flex--spacing-2 h-8">
              <TenantLogo size={100} />
              <h1 className="text-header">{fider.session.tenant.name}</h1>
            </Link>
            <HStack spacing={4} className="c-header__nav flex-items-center">
              <Link href="/" className={`c-header__nav-link ${isFeedbackActive ? "c-header__nav-link--active" : ""}`}>
                <Trans id="header.nav.feedback">All Feedback</Trans>
              </Link>
              <Link href="/roadmap" className={`c-header__nav-link ${isRoadmapActive ? "c-header__nav-link--active" : ""}`}>
                <Trans id="header.nav.roadmap">Roadmap</Trans>
              </Link>
            </HStack>
            {fider.session.isAuthenticated && (
              <div className="c-header__moderation">
                <ModerationIndicator />
              </div>
            )}
            {fider.session.isAuthenticated && (
              <HStack spacing={2} className="c-header__actions">
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
              <HStack spacing={2} className="c-header__actions">
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
        </div>
      </HStack>
    </div>
  )
}
