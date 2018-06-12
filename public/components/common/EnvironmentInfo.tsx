import * as React from "react";
import { SystemSettings } from "@fider/models";

export class EnvironmentInfo extends React.Component<{}, {}> {
  public render() {
    if (Fider.isProduction()) {
      return null;
    }

    return (
      <div className="c-env-info">
        Env: {Fider.settings.environment} | Compiler: {Fider.settings.compiler} | Version: {Fider.settings.version} |
        BuildTime: {Fider.settings.buildTime}
      </div>
    );
  }
}
