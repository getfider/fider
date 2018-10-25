import * as React from "react";
import * as Loadable from "react-loadable";

import { Loader } from "@fider/components/common/Loader";

export const AsyncHomePage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Home/Home.page");
    return module.HomePage;
  },
  loading: () => <Loader />
});

export const AsyncShowPostPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/ShowPost/ShowPost.page");
    return module.ShowPostPage;
  },
  loading: () => <Loader />
});

export const AsyncManageMembersPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageMembers.page");
    return module.ManageMembersPage;
  },
  loading: () => <Loader />
});

export const AsyncManageTagsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageTags.page");
    return module.ManageTagsPage;
  },
  loading: () => <Loader />
});

export const AsyncPrivacySettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/PrivacySettings.page");
    return module.PrivacySettingsPage;
  },
  loading: () => <Loader />
});

export const AsyncExportPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/Export.page");
    return module.ExportPage;
  },
  loading: () => <Loader />
});

export const AsyncInvitationsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/Invitations.page");
    return module.InvitationsPage;
  },
  loading: () => <Loader />
});

export const AsyncManageAuthenticationPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageAuthentication.page");
    return module.ManageAuthenticationPage;
  },
  loading: () => <Loader />
});

export const AsyncAdvancedSettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/AdvancedSettings.page");
    return module.AdvancedSettingsPage;
  },
  loading: () => <Loader />
});

export const AsyncGeneralSettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/GeneralSettings.page");
    return module.GeneralSettingsPage;
  },
  loading: () => <Loader />
});

export const AsyncSignInPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/SignIn/SignIn.page");
    return module.SignInPage;
  },
  loading: () => <Loader />
});

export const AsyncSignUpPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/SignUp/SignUp.page");
    return module.SignUpPage;
  },
  loading: () => <Loader />
});

export const AsyncCompleteSignInProfilePage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page");
    return module.CompleteSignInProfilePage;
  },
  loading: () => <Loader />
});

export const AsyncMyNotificationsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/MyNotifications/MyNotifications.page");
    return module.MyNotificationsPage;
  },
  loading: () => <Loader />
});

export const AsyncMySettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/MySettings/MySettings.page");
    return module.MySettingsPage;
  },
  loading: () => <Loader />
});

export const AsyncOAuthEchoPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/OAuthEcho/OAuthEcho.page");
    return module.OAuthEchoPage;
  },
  loading: () => <Loader />
});

export const AsyncUIToolkitPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/UI/UIToolkit.page");
    return module.UIToolkitPage;
  },
  loading: () => <Loader />
});
