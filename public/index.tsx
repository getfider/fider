import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { resolveRootComponent } from '@fider/router';
import { SignInModal } from '@fider/components/SignInModal';
import { Header, Footer } from '@fider/components/common';
import { analytics } from '@fider/services';
import * as $ from 'jquery';

const w = window as any;
w.$ = w.jQuery = $;

import 'semantic-ui/dist/semantic.min.js';
import 'semantic-ui/dist/semantic.min.css';
import '@fider/assets/styles/main.scss';

w.props = { };
w.set = (key: string, value: any): void => {
  w.props[key] = value;
};

window.addEventListener('error', (evt: ErrorEvent) => {
  analytics.error(evt.error);
});

document.addEventListener('DOMContentLoaded', () => {
  const root = document.getElementById('root');
  if (root) {
    const config = resolveRootComponent(location.pathname);
    ReactDOM.render(
      <div>
        {!w.props.user && React.createElement(SignInModal, w.props)}
        <div id={config.id}>
          {config.showHeader && React.createElement(Header, w.props)}
          {React.createElement(config.component, w.props)}
          {React.createElement(Footer, w.props)}
        </div>
      </div>,
      root
    );
  }
});
