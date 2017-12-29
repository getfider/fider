import * as React from 'react';
import { inject, injectables } from '@fider/di';
import { AppSettings } from '@fider/models';

interface EnvironmentInfoProps {
  settings: AppSettings;
}

export class EnvironmentInfo extends React.Component<EnvironmentInfoProps, {}> {
  public render() {
    if (this.props.settings.environment !== 'production') {
      return (
        <div
          id="environment-info"
          className=" ui mini negative message no-border no-margin"
        >
          Env: {this.props.settings.environment} |
          Compiler: {this.props.settings.compiler} |
          Version: {this.props.settings.version} |
          BuildTime: {this.props.settings.buildTime}
        </div>
      );
    }
  }
}
