import React from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { Fider, FiderContext } from './services/fider';
import { IconContext } from 'react-icons';
import { Footer, Header } from './components';
import HomePage from './pages/Home/Home.page';
import AdvancedSettingsPage from './pages/Administration/pages/AdvancedSettings.page';
import ExportPage from './pages/Administration/pages/Export.page';
import GeneralSettingsPage from './pages/Administration/pages/GeneralSettings.page';
import InvitationsPage from './pages/Administration/pages/Invitations.page';
import ManageAuthenticationPage from './pages/Administration/pages/ManageAuthentication.page';
import ManageMembersPage from './pages/Administration/pages/ManageMembers.page';
import ManageTagsPage from './pages/Administration/pages/ManageTags.page';
import PrivacySettingsPage from './pages/Administration/pages/PrivacySettings.page';
import CompleteSignInProfilePage from './pages/CompleteSignInProfile/CompleteSignInProfile.page';
import MyNotificationsPage from './pages/MyNotifications/MyNotifications.page';
import MySettingsPage from './pages/MySettings/MySettings.page';
import OAuthEchoPage from './pages/OAuthEcho/OAuthEcho.page';
import ShowPostPage from './pages/ShowPost/ShowPost.page';
import SignInPage from './pages/SignIn/SignIn.page';
import SignUpPage from './pages/SignUp/SignUp.page';
import UIToolkitPage from './pages/UI/UIToolkit.page';


interface PageConfiguration {
  regex: RegExp;
  component: any;
  showHeader: boolean;
}

const route = (path: string, component: any, showHeader: boolean = true): PageConfiguration => {
  path = path.replace("/", "/").replace(":number", "\\d+").replace(":string", ".+").replace("*", "/?.*");

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
  route("/-/ui", UIToolkitPage),
];

export const resolveRootComponent = (path: string): PageConfiguration => {
  if (path.length > 0 && path.charAt(path.length - 1) === "/") {
    path = path.substring(0, path.length - 1);
  }
  for (const entry of pathRegex) {
    if (entry && entry.regex.test(path)) {
      return entry;
    }
  }
  throw new Error(`Component not found for route ${path}.`);
};

function doWork(pathname: string, args: any) {
  let fider = Fider.initialize({...args });
  const config = resolveRootComponent(pathname);

  return renderToStaticMarkup(
    <FiderContext.Provider value={fider}>
      <IconContext.Provider value={{ className: "icon" }}>
        {config.showHeader && <Header />}
        {React.createElement(config.component, args.props)}
        {config.showHeader && <Footer />}
      </IconContext.Provider>
    </FiderContext.Provider>
  )
}

(globalThis as any).doWork = doWork