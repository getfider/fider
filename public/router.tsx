import * as React from "react";

import {
  HomePage,
  SignInPage,
  SignUpPage,
  ManageMembersPage,
  CompleteSignInProfilePage,
  PrivacySettingsPage,
  InvitationsPage,
  ExportPage,
  GeneralSettingsPage,
  AdvancedSettingsPage,
  ManageTagsPage,
  ShowIdeaPage,
  MySettingsPage,
  MyNotificationsPage
} from "@fider/pages";

interface PageConfiguration {
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, showHeader: boolean = true): PageConfiguration => {
  path = path
    .replace("/", "/")
    .replace(":number", "\\d+")
    .replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, showHeader };
};

const pathRegex = [
  route("", HomePage),
  route("/ideas/:number*", ShowIdeaPage),
  route("/admin/members", ManageMembersPage),
  route("/admin/tags", ManageTagsPage),
  route("/admin/privacy", PrivacySettingsPage),
  route("/admin/export", ExportPage),
  route("/admin/invitations", InvitationsPage),
  route("/admin/advanced", AdvancedSettingsPage),
  route("/admin", GeneralSettingsPage),
  route("/signin", SignInPage, false),
  route("/signup", SignUpPage, false),
  route("/signin/verify", CompleteSignInProfilePage),
  route("/invite/verify", CompleteSignInProfilePage),
  route("/notifications", MyNotificationsPage),
  route("/settings", MySettingsPage)
];

export const resolveRootComponent = (path: string): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of pathRegex) {
    if (entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};
