import React from "react";
import { Fider } from "@fider/services";

export const EnvironmentInfo = () => {
  if (Fider.isProduction()) {
    return null;
  }

  return (
    <div className="c-env-info">
      Env: {Fider.settings.environment} | Compiler: {Fider.settings.compiler} | Version: {Fider.settings.version} |
      BuildTime: {Fider.settings.buildTime || "N/A"} |
      {!Fider.isSingleHostMode() && `TenantID: ${Fider.session.tenant.id}`} |{" "}
      {Fider.session.isAuthenticated && `UserID: ${Fider.session.user.id}`}
    </div>
  );
};
