import * as React from 'react';
import * as ReactDOM from 'react-dom';

import * as $ from 'jquery';
(window as any).$ = (window as any).jQuery = $;
import 'semantic-ui/dist/semantic.min.js';
import 'semantic-ui/dist/semantic.min.css';
import '@fider/assets/styles/main.scss';

import { SignInModal } from '@fider/components/SignInModal';
import { AdminHomePage } from '@fider/pages/admin/AdminHomePage';
import { MembersPage } from '@fider/pages/admin/MembersPage';
import { ManageTagsPage } from '@fider/pages/admin/ManageTagsPage';
import { HomePage } from '@fider/pages/site/HomePage';
import { ShowIdeaPage } from '@fider/pages/site/ShowIdeaPage';
import { UserSettingsPage } from '@fider/pages/site/UserSettingsPage';
import { CompleteSignInProfilePage } from '@fider/pages/site/CompleteSignInProfilePage';
import { SignUpPage } from '@fider/pages/signup/SignUpPage';

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: any;
}

const route = (path: string, component: any, id: string): PageConfiguration => {
  path = path.replace('/', '\/')
             .replace(':number', '\\d+')
             .replace('*', '.*');

  const regex = new RegExp(`^${path}$`);
  return { regex, component, id };
};

const pathRegex = [
  route('/', HomePage, 'fdr-home-page'),
  route('/ideas/:number/*', ShowIdeaPage, 'fdr-show-idea-page'),
  route('/admin/members', MembersPage, 'fdr-admin-members-page'),
  route('/admin/tags', ManageTagsPage, 'fdr-admin-tags-page'),
  route('/admin', AdminHomePage, 'fdr-admin-page'),
  route('/signup', SignUpPage, 'fdr-signup-page'),
  route('/signin/verify', CompleteSignInProfilePage, 'fdr-complete-signin-profile'),
  route('/settings', UserSettingsPage, 'fdr-user-settings'),
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
