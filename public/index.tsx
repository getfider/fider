import * as React from "react";
import * as ReactDOM from "react-dom";
import { resolveRootComponent } from "@fider/router";
import { Header, Footer } from "@fider/components/common";
import { analytics, classSet, Fider } from "@fider/services";
import { ToastContainer, ToastPosition } from "react-toastify";

import "semantic-ui-css/components/reset.min.css";
import "semantic-ui-css/components/icon.min.css";
import "semantic-ui-css/components/dropdown.min.css";
import "semantic-ui-css/components/transition.min.css";

import "react-toastify/dist/ReactToastify.css";
import "@fider/assets/styles/main.scss";

window.addEventListener("error", (evt: ErrorEvent) => {
  analytics.error(evt.error);
});

document.addEventListener("DOMContentLoaded", () => {
  const fider = Fider.initialize();

  const root = document.getElementById("root");
  if (root) {
    const config = resolveRootComponent(location.pathname);
    document.body.className = classSet({
      "is-authenticated": fider.session.isAuthenticated,
      "is-staff": fider.session.isAuthenticated && fider.session.user.isCollaborator
    });
    ReactDOM.render(
      <>
        <ToastContainer position={ToastPosition.TOP_RIGHT} toastClassName="c-toast" />
        {config.showHeader && <Header />}
        {React.createElement(config.component, fider.session.props)}
        {config.showHeader && <Footer />}
      </>,
      root
    );
  }
});
