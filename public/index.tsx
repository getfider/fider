import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { AdminHomePage } from '@fider/pages/AdminHomePage';
import { SiteHomePage } from '@fider/pages/SiteHomePage';
import { ShowIdeaPage } from '@fider/pages/ShowIdeaPage';
import { SignUpPage } from '@fider/pages/SignUpPage';
import { container, injectables } from '@fider/di';
import {
    Session,
    BrowserSession,
    IdeaService,
    HttpIdeaService,
    TenantService,
    HttpTenantService
} from '@fider/services';

container.bind<Session>(injectables.Session).toConstantValue(new BrowserSession(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);

import '@fider/main.scss';

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
