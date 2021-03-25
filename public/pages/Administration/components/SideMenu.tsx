import "./SideMenu.scss"

import React from "react"
import { classSet } from "@fider/services"
import { FiderVersion } from "@fider/components"
import { useFider } from "@fider/hooks"

interface SiteMenuProps {
  activeItem: string
  visible: boolean
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
    "c-side-menu-item": true,
    "m-active": props.isActive,
  })

  if (props.isActive) {
    return (
      <span key={props.name} className={className}>
        {props.title}
      </span>
    )
  }

  return (
    <a key={props.name} className={className} href={props.href}>
      {props.title}
    </a>
  )
}

export const SideMenu = (props: SiteMenuProps) => {
  const fider = useFider()
  const activeItem = props.activeItem || "general"
  const style = { display: props.visible ? "" : "none" }

  return (
    <div className={props.className}>
      <div className="c-side-menu" style={style}>
        <SideMenuItem name="general" title="General" href="/admin" isActive={activeItem === "general"} />
        <SideMenuItem name="privacy" title="Privacy" href="/admin/privacy" isActive={activeItem === "privacy"} />
        <SideMenuItem name="members" title="Members" href="/admin/members" isActive={activeItem === "members"} />
        <SideMenuItem name="tags" title="Tags" href="/admin/tags" isActive={activeItem === "tags"} />
        <SideMenuItem name="invitations" title="Invitations" href="/admin/invitations" isActive={activeItem === "invitations"} />
        <SideMenuItem name="authentication" title="Authentication" href="/admin/authentication" isActive={activeItem === "authentication"} />
        <SideMenuItem name="advanced" title="Advanced" href="/admin/advanced" isActive={activeItem === "advanced"} />
        {fider.session.user.isAdministrator && (
          <>
            <SideMenuItem name="export" title="Export" href="/admin/export" isActive={activeItem === "export"} />
          </>
        )}
      </div>
      <FiderVersion />
    </div>
  )
}

interface SideMenuTogglerProps {
  onToggle: (active: boolean) => void
}

interface SideMenuTogglerState {
  active: boolean
}

export class SideMenuToggler extends React.Component<SideMenuTogglerProps, SideMenuTogglerState> {
  constructor(props: SideMenuTogglerProps) {
    super(props)
    this.state = {
      active: false,
    }
  }
  private toggle = () => {
    this.setState(
      (state) => ({ active: !state.active }),
      () => {
        this.props.onToggle(this.state.active)
      }
    )
  }

  public render() {
    const className = classSet({
      "c-side-menu-toggler": true,
      active: this.state.active,
    })
    return (
      <div className={className} onClick={this.toggle}>
        <div className="bar1" />
        <div className="bar2" />
        <div className="bar3" />
      </div>
    )
  }
}
