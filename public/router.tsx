import * as React from "react";

import {
  HomePage,
  SignUpPage,
  ManageMembersPage,
  CompleteSignInProfilePage,
  AdminHomePage,
  ManageTagsPage,
  ShowIdeaPage,
  MySettingsPage,
  MyNotificationsPage
} from "@fider/pages";

interface PageConfiguration {
  id: string;
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, id: string, showHeader: boolean): PageConfiguration => {
  path = path
    .replace("/", "/")
    .replace(":number", "\\d+")
    .replace("*", "/?.*");

  const regex = new RegExp(`^${path}$`);
  return { regex, component, id, showHeader };
};

const pathRegex = [
  route("", HomePage, "p-home", true),
  route("/ideas/:number*", ShowIdeaPage, "p-show-idea", true),
  route("/admin/members", ManageMembersPage, "p-manage-members", true),
  route("/admin/tags", ManageTagsPage, "p-manage-tags", true),
  route("/admin", AdminHomePage, "p-admin-home", true),
  route("/signup", SignUpPage, "p-signup", false),
  route("/signin/verify", CompleteSignInProfilePage, "p-complete-signin-profile", true),
  route("/notifications", MyNotificationsPage, "p-my-notifications", true),
  route("/settings", MySettingsPage, "p-my-settings", true)
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
