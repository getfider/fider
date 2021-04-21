import "./Header.scss"

import React, { useState, useEffect } from "react"
import { SignInModal, EnvironmentInfo, Avatar, TenantLogo } from "@fider/components"
import { actions } from "@fider/services"
import { FaUser, FaCog, FaCaretDown } from "react-icons/fa"
import { useFider } from "@fider/hooks"

export const Header = () => {
  const fider = useFider()
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [unreadNotifications, setUnreadNotifications] = useState(0)

  useEffect(() => {
    if (fider.session.isAuthenticated) {
      actions.getTotalUnreadNotifications().then((result) => {
        if (result.ok && result.data > 0) {
          setUnreadNotifications(result.data)
        }
      })
    }
  }, [fider.session.isAuthenticated])

  const showModal = () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
    }
  }

  const hideModal = () => setIsSignInModalOpen(false)

  const items = fider.session.isAuthenticated && (
    <div className="c-menu-user">
      <div className="c-menu-user-heading">
        <FaUser /> <span>{fider.session.user.name}</span>
      </div>
      <a href="/settings" className="c-menu-user-item">
        Settings
      </a>
      <a href="/notifications" className="c-menu-user-item">
        Notifications
        {unreadNotifications > 0 && <div className="c-unread-count">{unreadNotifications}</div>}
      </a>
      <div className="c-menu-user-divider" />
      {fider.session.user.isCollaborator && [
        <div key={1} className="c-menu-user-heading">
          <FaCog /> <span>Administration</span>
        </div>,
        <a key={2} href="/admin" className="c-menu-user-item">
          Site Settings
        </a>,
        <div key={5} className="c-menu-user-divider" />,
      ]}
      <a href="/signout?redirect=/" className="c-menu-user-item signout">
        Sign out
      </a>
    </div>
  )

  const showRightMenu = fider.session.isAuthenticated || !fider.session.tenant.isPrivate
  return (
    <div id="c-header">
      <EnvironmentInfo />
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <div className="c-menu">
        <div className="container">
          <a href="/" className="c-menu-item-title">
            <TenantLogo size={100} />
            <h1>{fider.session.tenant.name}</h1>
          </a>
          {showRightMenu && (
            <div onClick={showModal} className="c-menu-item-signin">
              {fider.session.isAuthenticated && <Avatar user={fider.session.user} />}
              {unreadNotifications > 0 && <div className="c-unread-dot" />}
              {!fider.session.isAuthenticated && <span>Sign in</span>}
              {fider.session.isAuthenticated && <FaCaretDown />}
              {items}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
