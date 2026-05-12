import React from "react"
import { useFider } from "@fider/hooks"
import { Avatar, Dropdown } from "./common"
import { Trans } from "@lingui/react/macro"
import IconCog from "@fider/assets/images/heroicons-cog.svg"
import IconWrench from "@fider/assets/images/heroicons-wrenchscrewdriver.svg"
import IconLeft from "@fider/assets/images/heroicons-arrowleft-rectangle.svg"

export const UserMenu = () => {
  const fider = useFider()

  return (
    <div className="c-menu-user">
      <Dropdown position="left" renderHandle={<Avatar user={fider.session.user} />}>
        <div className="p-2 text-medium uppercase">{fider.session.user.name}</div>
        <Dropdown.ListItem href="/settings" icon={IconCog}>
          <Trans id="menu.mysettings">My Settings</Trans>
        </Dropdown.ListItem>
        <Dropdown.Divider />

        {fider.session.user.isCollaborator && (
          <>
            <div className="p-2 text-medium uppercase">
              <Trans id="menu.administration">Administration</Trans>
            </div>
            <Dropdown.ListItem href="/admin" icon={IconWrench}>
              <Trans id="menu.sitesettings">Site Settings</Trans>
            </Dropdown.ListItem>
            <Dropdown.Divider />
          </>
        )}
        <Dropdown.ListItem href="/signout" icon={IconLeft}>
          <Trans id="menu.signout">Sign out</Trans>
        </Dropdown.ListItem>
      </Dropdown>
    </div>
  )
}
