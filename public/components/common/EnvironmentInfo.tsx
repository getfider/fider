import * as React from 'react';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';
import { AppSettings } from '@fider/models';

interface EnvironmentInfoProps {
  settings: AppSettings;
}

export class EnvironmentInfo extends React.Component<EnvironmentInfoProps, {}> {

  @inject(injectables.Session)
  public session: Session;

  public render() {
    if (!this.session.isProduction()) {
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
