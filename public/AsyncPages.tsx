import { lazy } from "react";

export const AsyncHomePage = lazy(() => import("@fider/pages/Home/Home.page"));

export const AsyncShowPostPage = lazy(() => import("@fider/pages/ShowPost/ShowPost.page"));

export const AsyncManageMembersPage = lazy(() => import("@fider/pages/Administration/pages/ManageMembers.page"));

export const AsyncManageTagsPage = lazy(() => import("@fider/pages/Administration/pages/ManageTags.page"));

export const AsyncPrivacySettingsPage = lazy(() => import("@fider/pages/Administration/pages/PrivacySettings.page"));

export const AsyncExportPage = lazy(() => import("@fider/pages/Administration/pages/Export.page"));

export const AsyncInvitationsPage = lazy(() => import("@fider/pages/Administration/pages/Invitations.page"));

export const AsyncManageAuthenticationPage = lazy(() =>
  import("@fider/pages/Administration/pages/ManageAuthentication.page")
);

export const AsyncAdvancedSettingsPage = lazy(() => import("@fider/pages/Administration/pages/AdvancedSettings.page"));

export const AsyncGeneralSettingsPage = lazy(() => import("@fider/pages/Administration/pages/GeneralSettings.page"));

export const AsyncSignInPage = lazy(() => import("@fider/pages/SignIn/SignIn.page"));

export const AsyncSignUpPage = lazy(() => import("@fider/pages/SignUp/SignUp.page"));

export const AsyncCompleteSignInProfilePage = lazy(() =>
  import("@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page")
);

export const AsyncMyNotificationsPage = lazy(() => import("@fider/pages/MyNotifications/MyNotifications.page"));

export const AsyncMySettingsPage = lazy(() => import("@fider/pages/MySettings/MySettings.page"));

export const AsyncOAuthEchoPage = lazy(() => import("@fider/pages/OAuthEcho/OAuthEcho.page"));

export const AsyncUIToolkitPage = lazy(() => import("@fider/pages/UI/UIToolkit.page"));
