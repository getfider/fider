import * as React from 'react';

import { AdminHomePage } from '@fider/pages/admin/AdminHomePage';
import { MembersPage } from '@fider/pages/admin/MembersPage';
import { ManageTagsPage } from '@fider/pages/admin/ManageTagsPage';
import { HomePage } from '@fider/pages/site/HomePage';
import { ShowIdeaPage } from '@fider/pages/site/ShowIdeaPage';
import { UserSettingsPage } from '@fider/pages/site/UserSettingsPage';
import { UserNotificationsPage } from '@fider/pages/site/UserNotificationsPage';
import { CompleteSignInProfilePage } from '@fider/pages/site/CompleteSignInProfilePage';
import { SignUpPage } from '@fider/pages/signup/SignUpPage';

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, id: string, showHeader: boolean): PageConfiguration => {
  path = path.replace('/', '\/')
             .replace(':number', '\\d+')
             .replace('*', '\/?.*');

  const regex = new RegExp(`^${path}$`);
  return { regex, component, id, showHeader };
};

const pathRegex = [
  route('', HomePage, 'fdr-home-page', true),
  route('/ideas/:number*', ShowIdeaPage, 'fdr-show-idea-page', true),
  route('/admin/members', MembersPage, 'fdr-admin-members-page', true),
  route('/admin/tags', ManageTagsPage, 'fdr-admin-tags-page', true),
  route('/admin', AdminHomePage, 'fdr-admin-page', true),
  route('/signup', SignUpPage, 'fdr-signup-page', false),
  route('/signin/verify', CompleteSignInProfilePage, 'fdr-complete-signin-profile', true),
  route('/notifications', UserNotificationsPage, 'fdr-user-notifications', true),
  route('/settings', UserSettingsPage, 'fdr-user-settings', true),
];

export const resolveRootComponent = (path: string): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === '/') {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of pathRegex) {
    if (entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};
