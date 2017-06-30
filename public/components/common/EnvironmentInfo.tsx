import * as React from 'react';
import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

export class EnvironmentInfo extends React.Component<{}, {}> {

  @inject(injectables.Session)
  public session: Session;

  public render() {
    if (!this.session.isProduction()) {
      const settings = this.session.getAppSettings();
      return <div className="ui mini negative message no-border no-margin">
                  Env: { settings.environment } |
                  Compiler: { settings.compiler } |
                  Version: { settings.version } |
                  BuildTime: { settings.buildTime }
              </div>;
    }
    return <div/>;
  }
}