import React from "react";
import Loadable from "react-loadable";

import { Loader } from "@fider/components/common/Loader";

const Loading = () => (
  <div className="page">
    <Loader />
  </div>
);

export const AsyncHomePage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "Home.page" */ "@fider/pages/Home/Home.page");
    return module.HomePage;
  },
  loading: () => <Loading />
});

export const AsyncShowPostPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "ShowPost.page" */ "@fider/pages/ShowPost/ShowPost.page");
    return module.ShowPostPage;
  },
  loading: () => <Loading />
});

export const AsyncManageMembersPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "ManageMembers.page" */ "@fider/pages/Administration/pages/ManageMembers.page");
    return module.ManageMembersPage;
  },
  loading: () => <Loading />
});

export const AsyncManageTagsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "ManageTags.page" */ "@fider/pages/Administration/pages/ManageTags.page");
    return module.ManageTagsPage;
  },
  loading: () => <Loading />
});

export const AsyncPrivacySettingsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "PrivacySettings.page" */ "@fider/pages/Administration/pages/PrivacySettings.page");
    return module.PrivacySettingsPage;
  },
  loading: () => <Loading />
});

export const AsyncExportPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "Export.page" */ "@fider/pages/Administration/pages/Export.page");
    return module.ExportPage;
  },
  loading: () => <Loading />
});

export const AsyncInvitationsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "Invitations.page" */ "@fider/pages/Administration/pages/Invitations.page");
    return module.InvitationsPage;
  },
  loading: () => <Loading />
});

export const AsyncManageAuthenticationPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "ManageAuthentication.page" */ "@fider/pages/Administration/pages/ManageAuthentication.page");
    return module.ManageAuthenticationPage;
  },
  loading: () => <Loading />
});

export const AsyncAdvancedSettingsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "AdvancedSettings.page" */ "@fider/pages/Administration/pages/AdvancedSettings.page");
    return module.AdvancedSettingsPage;
  },
  loading: () => <Loading />
});

export const AsyncGeneralSettingsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "GeneralSettings.page" */ "@fider/pages/Administration/pages/GeneralSettings.page");
    return module.GeneralSettingsPage;
  },
  loading: () => <Loading />
});

export const AsyncSignInPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "SignIn.page" */ "@fider/pages/SignIn/SignIn.page");
    return module.SignInPage;
  },
  loading: () => <Loading />
});

export const AsyncSignUpPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "SignUp.page" */ "@fider/pages/SignUp/SignUp.page");
    return module.SignUpPage;
  },
  loading: () => <Loading />
});

export const AsyncCompleteSignInProfilePage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "CompleteSignInProfile.page" */ "@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page");
    return module.CompleteSignInProfilePage;
  },
  loading: () => <Loading />
});

export const AsyncMyNotificationsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "MyNotifications.page" */ "@fider/pages/MyNotifications/MyNotifications.page");
    return module.MyNotificationsPage;
  },
  loading: () => <Loading />
});

export const AsyncMySettingsPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "MySettings.page" */ "@fider/pages/MySettings/MySettings.page");
    return module.MySettingsPage;
  },
  loading: () => <Loading />
});

export const AsyncOAuthEchoPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "OAuthEcho.page" */ "@fider/pages/OAuthEcho/OAuthEcho.page");
    return module.OAuthEchoPage;
  },
  loading: () => <Loading />
});

export const AsyncUIToolkitPage = Loadable({
  loader: async () => {
    const module = await import(/* webpackChunkName: "UIToolkit.page" */ "@fider/pages/UI/UIToolkit.page");
    return module.UIToolkitPage;
  },
  loading: () => <Loading />
});
