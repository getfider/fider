import * as React from 'react';
import { SystemSettings } from '@fider/models';

interface EnvironmentInfoProps {
  system: SystemSettings;
}

export class EnvironmentInfo extends React.Component<EnvironmentInfoProps, {}> {
  public render() {
    if (this.props.system.environment === 'production') {
      return null;
    }

    return (
      <div
        id="environment-info"
        className=" ui mini negative message no-border no-margin"
      >
        Env: {this.props.system.environment} |
        Compiler: {this.props.system.compiler} |
        Version: {this.props.system.version} |
        BuildTime: {this.props.system.buildTime}
      </div>
    );
  }
}
