import * as React from "react";
import * as Loadable from "react-loadable";

interface PageConfiguration {
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

export const LoadablePage = (importPathSuffix: string, exportName: string) => {
  return Loadable({
    loader: async () => {
      const module = await import(`@fider/pages/${importPathSuffix}`);
      return module[exportName];
    },
    loading: () => <div>Loading...</div>,
    delay: 400
  });
};

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
  route("", LoadablePage("Home/Home.page", "HomePage")),
  route("/posts/:number*", LoadablePage("ShowPost/ShowPost.page", "ShowPostPage")),
  route("/admin/members", LoadablePage("Administration/pages/ManageMembers.page", "ManageMembersPage")),
  route("/admin/tags", LoadablePage("Administration/pages/ManageTags.page", "ManageTagsPage")),
  route("/admin/privacy", LoadablePage("Administration/pages/PrivacySettings.page", "PrivacySettingsPage")),
  route("/admin/export", LoadablePage("Administration/pages/Export.page", "ExportPage")),
  route("/admin/invitations", LoadablePage("Administration/pages/Invitations.page", "InvitationsPage")),
  route(
    "/admin/authentication",
    LoadablePage("Administration/pages/ManageAuthentication.page", "ManageAuthenticationPage")
  ),
  route("/admin/advanced", LoadablePage("Administration/pages/AdvancedSettings.page", "AdvancedSettingsPage")),
  route("/admin", LoadablePage("Administration/pages/GeneralSettings.page", "GeneralSettingsPage")),
  route("/signin", LoadablePage("SignIn/SignIn.page", "SignInPage"), false),
  route("/signup", LoadablePage("SignUp/SignUp.page", "SignUpPage"), false),
  route(
    "/signin/verify",
    LoadablePage("CompleteSignInProfile/CompleteSignInProfile.page", "CompleteSignInProfilePage")
  ),
  route(
    "/invite/verify",
    LoadablePage("CompleteSignInProfile/CompleteSignInProfile.page", "CompleteSignInProfilePage")
  ),
  route("/notifications", LoadablePage("MyNotifications/MyNotifications.page", "MyNotificationsPage")),
  route("/settings", LoadablePage("MySettings/MySettings.page", "MySettingsPage")),
  route("/oauth/:string/echo", LoadablePage("OAuthEcho/OAuthEcho.page", "OAuthEchoPage"), false),
  route("/-/ui", LoadablePage("UI/UIToolkit.page", "UIToolkitPage"))
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
