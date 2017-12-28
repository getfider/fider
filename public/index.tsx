import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { container, injectables } from '@fider/di';
import { Session } from '@fider/services';
import { resolveRootComponent } from '@fider/router';
import { SignInModal } from '@fider/components/SignInModal';
import { Header, Footer } from '@fider/components/common';

import * as $ from 'jquery';
(window as any).$ = (window as any).jQuery = $;
import 'semantic-ui/dist/semantic.min.js';
import 'semantic-ui/dist/semantic.min.css';
import '@fider/assets/styles/main.scss';

document.addEventListener('DOMContentLoaded', () => {
  const root = document.getElementById('root');
  if (root) {
    const session = container.get<Session>(injectables.Session);
    const config = resolveRootComponent(location.pathname);
    ReactDOM.render(
      <div>
        {React.createElement(SignInModal, session.props())}
        <div id={config.id}>
          {config.showHeader && React.createElement(Header, session.props())}
          {React.createElement(config.component, session.props())}
          {React.createElement(Footer, session.props())}
        </div>
      </div>, root
    );
  }
});
