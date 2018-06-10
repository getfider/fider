import * as React from "react";
import { SystemSettings } from "@fider/models";

export class EnvironmentInfo extends React.Component<{}, {}> {
  public render() {
    if (window.props.settings.environment === "production") {
      return null;
    }

    return (
      <div className="c-env-info">
        Env: {window.props.settings.environment} | Compiler: {window.props.settings.compiler} | Version:{" "}
        {window.props.settings.version} | BuildTime: {window.props.settings.buildTime}
      </div>
    );
  }
}
