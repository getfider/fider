import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { resolveRootComponent } from '@fider/router';
import { SignInModal } from '@fider/components/SignInModal';
import { Header, Footer } from '@fider/components/common';
import { analytics } from '@fider/services';
import * as $ from 'jquery';

const w = window as any;
w.$ = w.jQuery = $;

import 'semantic-ui-css/components/site.min.js';
import 'semantic-ui-css/components/modal.min.js';
import 'semantic-ui-css/components/transition.min.js';
import 'semantic-ui-css/components/dropdown.min.js';
import 'semantic-ui-css/components/dimmer.min.js';

import 'semantic-ui-css/components/reset.min.css';
import 'semantic-ui-css/components/site.min.css';
import 'semantic-ui-css/components/segment.min.css';
import 'semantic-ui-css/components/table.min.css';
import 'semantic-ui-css/components/checkbox.min.css';
import 'semantic-ui-css/components/transition.min.css';
import 'semantic-ui-css/components/header.min.css';
import 'semantic-ui-css/components/button.min.css';
import 'semantic-ui-css/components/form.min.css';
import 'semantic-ui-css/components/container.min.css';
import 'semantic-ui-css/components/modal.min.css';
import 'semantic-ui-css/components/dropdown.min.css';
import 'semantic-ui-css/components/input.min.css';
import 'semantic-ui-css/components/label.min.css';
import 'semantic-ui-css/components/list.min.css';
import 'semantic-ui-css/components/divider.min.css';
import 'semantic-ui-css/components/item.min.css';
import 'semantic-ui-css/components/grid.min.css';
import 'semantic-ui-css/components/icon.min.css';
import 'semantic-ui-css/components/message.min.css';
import 'semantic-ui-css/components/menu.min.css';
import 'semantic-ui-css/components/dimmer.min.css';
import 'semantic-ui-css/components/comment.min.css';
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
