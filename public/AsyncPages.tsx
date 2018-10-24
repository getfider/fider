import * as React from "react";
import * as Loadable from "react-loadable";

const Loading = () => <div>Loading...</div>;

export const AsyncHomePage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Home/Home.page");
    return module.HomePage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncShowPostPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/ShowPost/ShowPost.page");
    return module.ShowPostPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncManageMembersPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageMembers.page");
    return module.ManageMembersPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncManageTagsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageTags.page");
    return module.ManageTagsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncPrivacySettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/PrivacySettings.page");
    return module.PrivacySettingsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncExportPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/Export.page");
    return module.ExportPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncInvitationsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/Invitations.page");
    return module.InvitationsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncManageAuthenticationPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/ManageAuthentication.page");
    return module.ManageAuthenticationPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncAdvancedSettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/AdvancedSettings.page");
    return module.AdvancedSettingsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncGeneralSettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/Administration/pages/GeneralSettings.page");
    return module.GeneralSettingsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncSignInPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/SignIn/SignIn.page");
    return module.SignInPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncSignUpPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/SignUp/SignUp.page");
    return module.SignUpPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncCompleteSignInProfilePage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page");
    return module.CompleteSignInProfilePage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncMyNotificationsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/MyNotifications/MyNotifications.page");
    return module.MyNotificationsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncMySettingsPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/MySettings/MySettings.page");
    return module.MySettingsPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncOAuthEchoPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/OAuthEcho/OAuthEcho.page");
    return module.OAuthEchoPage;
  },
  loading: Loading,
  delay: 400
});

export const AsyncUIToolkitPage = Loadable({
  loader: async () => {
    const module = await import("@fider/pages/UI/UIToolkit.page");
    return module.UIToolkitPage;
  },
  loading: Loading,
  delay: 400
});
