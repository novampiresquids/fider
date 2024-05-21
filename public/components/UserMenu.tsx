import React from "react"
import { useFider } from "@fider/hooks"
import { Avatar, Dropdown } from "./common"
import { Trans } from "@lingui/macro"

export const UserMenu = () => {
  const fider = useFider()

  return (
    <div className="c-menu-user">
      <Dropdown position="left" renderHandle={<Avatar user={fider.session.user} />}>
        <div className="p-2 text-medium uppercase">{fider.session.user.name}</div>
        <Dropdown.ListItem href="/settings">
          <Trans id="menu.mysettings">My Settings</Trans>
        </Dropdown.ListItem>
        <Dropdown.ListItem href="/">
          <Trans id="menu.myboards">My Boards</Trans>
        </Dropdown.ListItem>
        <Dropdown.Divider />

        {fider.session.user.isCollaborator && fider.session.tenant && (
          <>
            <div className="p-2 text-medium uppercase">
              <Trans id="menu.administration">Administration</Trans>
            </div>
            <Dropdown.ListItem href={`/board/${fider.session.tenant.id}/admin`}>
              <Trans id="menu.sitesettings">Board Settings</Trans>
            </Dropdown.ListItem>
            <Dropdown.Divider />
          </>
        )}
        <Dropdown.ListItem href="/signout?redirect=/">
          <Trans id="menu.signout">Sign out</Trans>
        </Dropdown.ListItem>
      </Dropdown>
    </div>
  )
}
