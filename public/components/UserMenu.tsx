import React from "react"
import { useFider } from "@fider/hooks"
import { Avatar, Dropdown } from "./common"
import { Trans } from "@lingui/react/macro"
import IconCog from "@fider/assets/images/heroicons-cog.svg"
import IconShieldCheck from "@fider/assets/images/heroicons-shieldcheck.svg"
import IconXCircle from "@fider/assets/images/heroicons-x-circle.svg"

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
            <Dropdown.ListItem href="/admin" icon={IconShieldCheck}>
              <Trans id="menu.sitesettings">Site Settings</Trans>
            </Dropdown.ListItem>
            <Dropdown.Divider />
          </>
        )}
        <Dropdown.ListItem href="/signout" icon={IconXCircle}>
          <Trans id="menu.signout">Sign out</Trans>
        </Dropdown.ListItem>
      </Dropdown>
    </div>
  )
}
