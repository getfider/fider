import { lazy } from "react";

const retry = (fn: () => Promise<{ default: any }>, retriesLeft = 5, interval = 500): Promise<{ default: any }> => {
  return new Promise((resolve, reject) => {
    fn()
      .then(resolve)
      .catch(err => {
        setTimeout(() => {
          if (retriesLeft === 1) {
            reject(err);
            return;
          }
          retry(fn, interval, retriesLeft - 1).then(resolve, reject);
        }, interval);
      });
  });
};

export const AsyncHomePage = lazy(() => retry(() => import("@fider/pages/Home/Home.page")));

export const AsyncShowPostPage = lazy(() => retry(() => import("@fider/pages/ShowPost/ShowPost.page")));

export const AsyncManageMembersPage = lazy(() =>
  retry(() => import("@fider/pages/Administration/pages/ManageMembers.page"))
);

export const AsyncManageTagsPage = lazy(() => retry(() => import("@fider/pages/Administration/pages/ManageTags.page")));

export const AsyncPrivacySettingsPage = lazy(() =>
  retry(() => import("@fider/pages/Administration/pages/PrivacySettings.page"))
);

export const AsyncExportPage = lazy(() => retry(() => import("@fider/pages/Administration/pages/Export.page")));

export const AsyncInvitationsPage = lazy(() =>
  retry(() => import("@fider/pages/Administration/pages/Invitations.page"))
);

export const AsyncManageAuthenticationPage = lazy(() =>
  import("@fider/pages/Administration/pages/ManageAuthentication.page")
);

export const AsyncAdvancedSettingsPage = lazy(() =>
  retry(() => import("@fider/pages/Administration/pages/AdvancedSettings.page"))
);

export const AsyncGeneralSettingsPage = lazy(() =>
  retry(() => import("@fider/pages/Administration/pages/GeneralSettings.page"))
);

export const AsyncSignInPage = lazy(() => retry(() => import("@fider/pages/SignIn/SignIn.page")));

export const AsyncSignUpPage = lazy(() => retry(() => import("@fider/pages/SignUp/SignUp.page")));

export const AsyncCompleteSignInProfilePage = lazy(() =>
  import("@fider/pages/CompleteSignInProfile/CompleteSignInProfile.page")
);

export const AsyncMyNotificationsPage = lazy(() =>
  retry(() => import("@fider/pages/MyNotifications/MyNotifications.page"))
);

export const AsyncMySettingsPage = lazy(() => retry(() => import("@fider/pages/MySettings/MySettings.page")));

export const AsyncOAuthEchoPage = lazy(() => retry(() => import("@fider/pages/OAuthEcho/OAuthEcho.page")));

export const AsyncUIToolkitPage = lazy(() => retry(() => import("@fider/pages/UI/UIToolkit.page")));
