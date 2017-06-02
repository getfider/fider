import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { HomePage as AdminHomePage } from './admin/HomePage';
import { HomePage as SiteHomePage } from './site/HomePage';
import { ShowIdeaPage } from './site/ShowIdeaPage';
import { setup } from './storage';

import './style/main.scss';

const pathRegex = [
    { regex: new RegExp('^\/$'), component: <SiteHomePage /> },
    { regex: new RegExp('^\/admin$'), component: <AdminHomePage /> },
    { regex: new RegExp('^\/ideas\/\\d+.*$'), component: <ShowIdeaPage /> },

];

const resolveRootComponent = (path: string): JSX.Element => {
    for (const entry of pathRegex) {
        if (entry.regex.test(path)) {
            return entry.component;
        }
    }

    return <div />;
};

setup();

document.addEventListener('DOMContentLoaded', () => {
  ReactDOM.render(
      resolveRootComponent(location.pathname),
      document.getElementById('root')
  );
});
