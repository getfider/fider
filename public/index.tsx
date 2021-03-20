import React, { Suspense } from "react";
import ReactDOM from "react-dom";
import { resolveRootComponent, route } from "@fider/router";
import { Header, Footer, Loader } from "@fider/components/common";
import { ErrorBoundary } from "@fider/components";
import { classSet, Fider, FiderContext, actions } from "@fider/services";
import { IconContext } from "react-icons";
import * as Pages from "@fider/AsyncPages";
import "@fider/assets/styles/index.scss";

const Loading = () => (
  <div className="page">
    <Loader />
  </div>
);

const logProductionError = (err: Error) => {
  if (Fider.isProduction()) {
    console.error(err); // tslint:disable-line
    actions.logError(`react.ErrorBoundary: ${err.message}`, err);
  }
};

window.addEventListener("unhandledrejection", (evt: PromiseRejectionEvent) => {
  if (evt.reason instanceof Error) {
    actions.logError(`window.unhandledrejection: ${evt.reason.message}`, evt.reason);
  } else if (evt.reason) {
    actions.logError(`window.unhandledrejection: ${evt.reason.toString()}`);
  }
});

window.addEventListener("error", (evt: ErrorEvent) => {
  if (evt.error && evt.colno > 0 && evt.lineno > 0) {
    actions.logError(`window.error: ${evt.message}`, evt.error);
  }
});

const routes = [
  route("", Pages.AsyncHomePage),
  route("/posts/:number*", Pages.AsyncShowPostPage),
  route("/admin/members", Pages.AsyncManageMembersPage),
  route("/admin/tags", Pages.AsyncManageTagsPage),
  route("/admin/privacy", Pages.AsyncPrivacySettingsPage),
  route("/admin/export", Pages.AsyncExportPage),
  route("/admin/invitations", Pages.AsyncInvitationsPage),
  route("/admin/authentication", Pages.AsyncManageAuthenticationPage),
  route("/admin/advanced", Pages.AsyncAdvancedSettingsPage),
  route("/admin", Pages.AsyncGeneralSettingsPage),
  route("/signin", Pages.AsyncSignInPage, false),
  route("/signup", Pages.AsyncSignUpPage, false),
  route("/signin/verify", Pages.AsyncCompleteSignInProfilePage),
  route("/invite/verify", Pages.AsyncCompleteSignInProfilePage),
  route("/notifications", Pages.AsyncMyNotificationsPage),
  route("/settings", Pages.AsyncMySettingsPage),
  route("/oauth/:string/echo", Pages.AsyncOAuthEchoPage, false),
  route("/-/ui", Pages.AsyncUIToolkitPage),
];

(() => {
  let fider;

  fider = Fider.initialize();

  __webpack_nonce__ = fider.session.contextID;
  __webpack_public_path__ = `${fider.settings.globalAssetsURL}/assets/`;

  const config = resolveRootComponent(location.pathname, routes);
  document.body.className = classSet({
    "is-authenticated": fider.session.isAuthenticated,
    "is-staff": fider.session.isAuthenticated && fider.session.user.isCollaborator,
  });
  ReactDOM.render(
    <React.StrictMode>
      <ErrorBoundary onError={logProductionError}>
        <FiderContext.Provider value={fider}>
          <IconContext.Provider value={{ className: "icon" }}>
            {config.showHeader && <Header />}
            <Suspense fallback={<Loading />}>{React.createElement(config.component, fider.session.props)}</Suspense>
            {config.showHeader && <Footer />}
          </IconContext.Provider>
        </FiderContext.Provider>
      </ErrorBoundary>
    </React.StrictMode>,
    document.getElementById("root")
  );
})();
