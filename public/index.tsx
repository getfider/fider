import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { HomePage as AdminHomePage } from './admin/HomePage';
import { HomePage as SiteHomePage } from './site/HomePage';
import { ShowIdeaPage } from './site/ShowIdeaPage';
import { SignUpPage } from './setup/SignUpPage';
import { container, injectables } from './di';
import {
    Session,
    BrowserSession,
    IdeaService,
    HttpIdeaService,
    TenantService,
    HttpTenantService
} from './services';

container.bind<Session>(injectables.Session).toConstantValue(new BrowserSession(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);

import './main.scss';

const pathRegex = [
    { regex: new RegExp('^\/$'), component: <SiteHomePage /> },
    { regex: new RegExp('^\/admin$'), component: <AdminHomePage /> },
    { regex: new RegExp('^\/signup$'), component: <SignUpPage /> },
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

document.addEventListener('DOMContentLoaded', () => {
  ReactDOM.render(
      resolveRootComponent(location.pathname),
      document.getElementById('root')
  );
});
