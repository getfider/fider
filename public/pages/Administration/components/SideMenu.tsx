import "./SideMenu.scss"

import React, { useState } from "react"
import { classSet, basePath } from "@fider/services"
import { Icon } from "@fider/components"
import { useFider } from "@fider/hooks"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconMenu from "@fider/assets/images/heroicons-menu.svg"
import { VStack } from "@fider/components/layout"

interface SiteMenuProps {
  activeItem: string
  className?: string
}

interface SideMenuItemProps {
  name: string
  title: string
  isActive: boolean
  href: string
}

const SideMenuItem = (props: SideMenuItemProps) => {
  const className = classSet({
    "c-side-menu__item": true,
    "c-side-menu__item--active": props.isActive,
  })

  return (
    <a key={props.name} className={className} href={props.href}>
      {props.title}
    </a>
  )
}

export const SideMenu = (props: SiteMenuProps) => {
  const fider = useFider()
  const activeItem = props.activeItem || "general"

  return (
    <div className="js-admin-menu sm:hidden md:hidden lg:block">
      <VStack spacing={0} className="c-side-menu rounded-md shadow bg-white">
        <SideMenuItem name="general" title="General" href={`${basePath()}/admin`} isActive={activeItem === "general"} />
        <SideMenuItem name="privacy" title="Privacy" href={`${basePath()}/admin/privacy`} isActive={activeItem === "privacy"} />
        <SideMenuItem name="users" title="Users" href={`${basePath()}/admin/users`} isActive={activeItem === "users"} />
        <SideMenuItem name="tags" title="Tags" href={`${basePath()}/admin/tags`} isActive={activeItem === "tags"} />
        <SideMenuItem name="invitations" title="Invitations" href={`${basePath()}/admin/invitations`} isActive={activeItem === "invitations"} />
        <SideMenuItem name="authentication" title="Authentication" href={`${basePath()}/admin/authentication`} isActive={activeItem === "authentication"} />
        <SideMenuItem name="advanced" title="Advanced" href={`${basePath()}/admin/advanced`} isActive={activeItem === "advanced"} />
        {fider.session.user.isAdministrator && (
          <>
            {fider.settings.isBillingEnabled && <SideMenuItem name="billing" title="Billing" href={`${basePath()}/admin/billing`} isActive={activeItem === "billing"} />}
            <SideMenuItem name="webhooks" title="Webhooks" href={`${basePath()}/admin/webhooks`} isActive={activeItem === "webhooks"} />
            <SideMenuItem name="export" title="Export" href={`${basePath()}/admin/export`} isActive={activeItem === "export"} />
          </>
        )}
      </VStack>
    </div>
  )
}

export const SideMenuToggler = () => {
  const [isActive, setIsActive] = useState(false)

  const toggle = () => {
    const classes = ["sm:hidden", "md:hidden"]
    const el = document.querySelector(".js-admin-menu") as HTMLElement
    if (el && !isActive) {
      el.classList.remove(...classes)
    } else if (el && isActive) {
      el.classList.add(...classes)
    }
    setIsActive(!isActive)
  }

  return (
    <div className="h-8 w-8 lg:hidden xl:hidden" onClick={toggle}>
      {isActive ? <Icon sprite={IconX} /> : <Icon sprite={IconMenu} />}
    </div>
  )
}
