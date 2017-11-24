import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { SignInModal } from '@fider/components/SignInModal';
import { AdminHomePage } from '@fider/pages/admin/AdminHomePage';
import { MembersPage } from '@fider/pages/admin/MembersPage';
import { ManageTagsPage } from '@fider/pages/admin/ManageTagsPage';
import { HomePage } from '@fider/pages/site/HomePage';
import { ShowIdeaPage } from '@fider/pages/site/ShowIdeaPage';
import { UserSettingsPage } from '@fider/pages/site/UserSettingsPage';
import { CompleteSignInProfilePage } from '@fider/pages/site/CompleteSignInProfilePage';
import { SignUpPage } from '@fider/pages/signup/SignUpPage';
import { container, injectables } from '@fider/di';

import {
  Session,
  BrowserSession,
  IdeaService,
  HttpIdeaService,
  TenantService,
  HttpTenantService,
  UserService,
  HttpUserService,
  TagService,
  HttpTagService
} from '@fider/services';

import '@fider/assets/styles/main.scss';

container.bind<Session>(injectables.Session).toConstantValue(new BrowserSession(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);
container.bind<UserService>(injectables.UserService).to(HttpUserService);
container.bind<TagService>(injectables.TagService).to(HttpTagService);

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: any;
}

const pathRegex = [
  { regex: new RegExp('^\/$'), component: HomePage, id: 'fdr-home-page' },
  { regex: new RegExp('^\/admin\/members$'), component: MembersPage, id: 'fdr-admin-members-page' },
  { regex: new RegExp('^\/admin\/tags$'), component: ManageTagsPage, id: 'fdr-admin-tags-page' },
  { regex: new RegExp('^\/admin$'), component: AdminHomePage, id: 'fdr-admin-page' },
  { regex: new RegExp('^\/signup$'), component: SignUpPage, id: 'fdr-signup-page' },
  { regex: new RegExp('^\/ideas\/\\d+.*$'), component: ShowIdeaPage, id: 'fdr-show-idea-page' },
  { regex: new RegExp('^\/signin\/verify'), component: CompleteSignInProfilePage, id: 'fdr-complete-signin-profile' },
  { regex: new RegExp('^\/settings$'), component: UserSettingsPage, id: 'fdr-user-settings' },
];

const resolveRootComponent = (path: string): PageConfiguration | undefined => {
  for (const entry of pathRegex) {
    if (entry.regex.test(path)) {
      return entry;
    }
  }
};

document.addEventListener('DOMContentLoaded', () => {
  const root = document.getElementById('root');
  if (root) {
    const config = resolveRootComponent(location.pathname)!;
    ReactDOM.render(
      <div>
        <SignInModal />
        <div id={config.id}>
          {React.createElement(config.component)}
        </div>
      </div>, root
    );
  }
});
