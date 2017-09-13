import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { SignInModal } from '@fider/components/SignInModal';
import { AdminHomePage } from '@fider/pages/admin/AdminHomePage';
import { SiteHomePage } from '@fider/pages/site/SiteHomePage';
import { ShowIdeaPage } from '@fider/pages/site/ShowIdeaPage';
import { SignUpPage } from '@fider/pages/signup/SignUpPage';
import { container, injectables } from '@fider/di';
import {
  Session,
  BrowserSession,
  IdeaService,
  HttpIdeaService,
  TenantService,
  HttpTenantService
} from '@fider/services';

import '@fider/assets/styles/main.scss';

container.bind<Session>(injectables.Session).toConstantValue(new BrowserSession(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: JSX.Element;
}

const pathRegex = [
  { regex: new RegExp('^\/$'), component: <SiteHomePage />, id: 'fdr-home-page' },
  { regex: new RegExp('^\/admin$'), component: <AdminHomePage />, id: 'fdr-admin-page' },
  { regex: new RegExp('^\/signup$'), component: <SignUpPage />, id: 'fdr-signup-page' },
  { regex: new RegExp('^\/ideas\/\\d+.*$'), component: <ShowIdeaPage />, id: 'fdr-show-idea-page' },
];

const resolveRootComponent = (path: string): PageConfiguration | undefined => {
  for (const entry of pathRegex) {
    if (entry.regex.test(path)) {
      return entry;
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const config = resolveRootComponent(location.pathname)!;
  ReactDOM.render(
    <div id={ config.id }>
      <SignInModal />
      { config.component }
    </div>, document.getElementById('root')
  );
});
