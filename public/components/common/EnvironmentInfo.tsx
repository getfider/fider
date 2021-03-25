import React from "react"
import { useFider } from "@fider/hooks"

export const EnvironmentInfo = () => {
  const fider = useFider()

  if (fider.isProduction()) {
    return null
  }

  return (
    <div className="c-env-info">
      Env: {fider.settings.environment} | Compiler: {fider.settings.compiler} | Version: {fider.settings.version} | BuildTime:{" "}
      {fider.settings.buildTime || "N/A"} |{!fider.isSingleHostMode() && `TenantID: ${fider.session.tenant.id}`} |{" "}
      {fider.session.isAuthenticated && `UserID: ${fider.session.user.id}`}
    </div>
  )
}
