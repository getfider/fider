import "./PortalDirectory.page.scss"

import React, { useState } from "react"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { Input } from "@fider/components"

export interface PortalSummary {
  name: string
  url: string
  host: string
  logoURL?: string
}

interface PortalDirectoryPageProps {
  portals: PortalSummary[]
}

const PortalCard = (props: { portal: PortalSummary }) => {
  const { portal } = props

  return (
    <a className="c-portal-card" href={portal.url}>
      <div className="c-portal-card__logo">
        {portal.logoURL ? <img src={portal.logoURL} alt={portal.name} /> : <div className="c-portal-card__initial">{portal.name.charAt(0).toUpperCase()}</div>}
      </div>
      <div className="c-portal-card__text">
        <span className="c-portal-card__name text-title">{portal.name}</span>
        <span className="c-portal-card__host text-muted text-xs">{portal.host}</span>
      </div>
    </a>
  )
}

export const PortalDirectoryPage = (props: PortalDirectoryPageProps) => {
  const [filter, setFilter] = useState("")

  const term = filter.trim().toLowerCase()
  const visible = term ? props.portals.filter((portal) => portal.name.toLowerCase().includes(term) || portal.host.toLowerCase().includes(term)) : props.portals

  return (
    <div id="p-portal-directory" className="page container">
      <h1 className="text-display mb-2">
        <Trans id="portaldirectory.title">Portals</Trans>
      </h1>
      <p className="text-muted mb-4">
        <Trans id="portaldirectory.subtitle">Browse the feedback portals hosted here.</Trans>
      </p>

      {props.portals.length === 0 ? (
        <p className="c-portal-directory__empty text-muted">
          <Trans id="portaldirectory.empty">There are no portals yet.</Trans>
        </p>
      ) : (
        <>
          <Input field="filter" placeholder={i18n._({ id: "portaldirectory.search.placeholder", message: "Search portals" })} onChange={setFilter} />

          {visible.length === 0 ? (
            <p className="c-portal-directory__nomatch text-muted mt-4">
              <Trans id="portaldirectory.nomatch">No portals match your search.</Trans>
            </p>
          ) : (
            <div className="c-portal-list mt-4">
              {visible.map((portal) => (
                <PortalCard key={portal.url} portal={portal} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default PortalDirectoryPage
