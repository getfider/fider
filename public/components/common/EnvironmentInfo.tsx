import * as React from "react";
import { SystemSettings } from "@fider/models";

export class EnvironmentInfo extends React.Component<{}, {}> {
  public render() {
    if (page.isProduction()) {
      return null;
    }

    return (
      <div className="c-env-info">
        Env: {page.settings.environment} | Compiler: {page.settings.compiler} | Version: {page.settings.version} |
        BuildTime: {page.settings.buildTime}
      </div>
    );
  }
}
