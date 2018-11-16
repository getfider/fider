import { lazy, ComponentType } from "react";

type LazyImport = () => Promise<{ default: ComponentType<any> }>;

const MAX_RETRIES = 5;

const retry = (fn: LazyImport, retriesLeft = MAX_RETRIES, interval = 500): Promise<{ default: ComponentType<any> }> => {
  return new Promise((resolve, reject) => {
    fn()
      .then(resolve)
      .catch(err => {
        setTimeout(() => {
          if (retriesLeft === 1) {
            reject(new Error(`${err} after ${MAX_RETRIES} retries`));
            return;
          }
          retry(fn, interval, retriesLeft - 1).then(resolve, reject);
        }, interval);
      });
  });
};

const load = (fn: LazyImport) => lazy(() => retry(() => fn()));

export const AsyncHomePage = load(() => import("@fider/pages/Home/Home.page"));

export const AsyncShowPostPage = load(() => import("@fider/pages/ShowPost/ShowPost.page"));

export const AsyncManageMembersPage = load(() => import("@fider/pages/Administration/pages/ManageMembers.page"));

export const AsyncManageTagsPage = load(() => import("@fider/pages/Administration/pages/ManageTags.page"));

export const AsyncPrivacySettingsPage = load(() => import("@fider/pages/Administration/pages/PrivacySettings.page"));

export const AsyncExportPage = load(() => import("@fider/pages/Administration/pages/Export.page"));

export const AsyncInvitationsPage = load(() => import("@fider/pages/Administration/pages/Invitations.page"));

export const AsyncManageAuthenticationPage = load(() =>
  import("@fider/pages/Administration/pages/ManageAuthentication.page")
);

export const AsyncAdvancedSettingsPage = load(() => import("@fider/pages/Administration/pages/AdvancedSettings.page"));

export const AsyncGeneralSettingsPage = load(() => import("@fider/pages/Administration/pages/GeneralSettings.page"));

export const AsyncSignInPage = load(() => import("@fider/pages/SignIn/SignIn.page"));

export const AsyncSignUpPage = load(() => import("@fider/pages/SignUp/SignUp.page"));

export const AsyncCompleteSignInProfilePage = load(() =>
  import("@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page")
);

export const AsyncMyNotificationsPage = load(() => import("@fider/pages/MyNotifications/MyNotifications.page"));

export const AsyncMySettingsPage = load(() => import("@fider/pages/MySettings/MySettings.page"));

export const AsyncOAuthEchoPage = load(() => import("@fider/pages/OAuthEcho/OAuthEcho.page"));

export const AsyncUIToolkitPage = load(() => import("@fider/pages/UI/UIToolkit.page"));
