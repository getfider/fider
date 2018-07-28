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
  ManageAuthenticationPage,
  ManageTagsPage,
  OAuthEchoPage,
  ShowPostPage,
  MySettingsPage,
  MyNotificationsPage,
  UIToolkitPage
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
    .replace(":string", ".+")
    .replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, showHeader };
};

const pathRegex = [
  route("", HomePage),
  route("/posts/:number*", ShowPostPage),
  route("/admin/members", ManageMembersPage),
  route("/admin/tags", ManageTagsPage),
  route("/admin/privacy", PrivacySettingsPage),
  route("/admin/export", ExportPage),
  route("/admin/invitations", InvitationsPage),
  route("/admin/authentication", ManageAuthenticationPage),
  route("/admin/advanced", AdvancedSettingsPage),
  route("/admin", GeneralSettingsPage),
  route("/signin", SignInPage, false),
  route("/signup", SignUpPage, false),
  route("/signin/verify", CompleteSignInProfilePage),
  route("/invite/verify", CompleteSignInProfilePage),
  route("/notifications", MyNotificationsPage),
  route("/settings", MySettingsPage),
  route("/oauth/:string/echo", OAuthEchoPage, false),
  route("/-/ui", UIToolkitPage)
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
