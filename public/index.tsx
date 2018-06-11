import * as React from "react";
import * as ReactDOM from "react-dom";
import { resolveRootComponent } from "@fider/router";
import { Header, Footer } from "@fider/components/common";
import { analytics, classSet } from "@fider/services";
import { ToastContainer } from "react-toastify";
import { FiderPage } from "@fider/models";

import "semantic-ui-css/components/reset.min.css";
import "semantic-ui-css/components/icon.min.css";
import "semantic-ui-css/components/dropdown.min.css";
import "semantic-ui-css/components/transition.min.css";

import "react-toastify/dist/ReactToastify.css";
import "@fider/assets/styles/main.scss";

const page = new FiderPage();
(global as any).page = page;

window.addEventListener("error", (evt: ErrorEvent) => {
  analytics.error(evt.error);
});

document.addEventListener("DOMContentLoaded", () => {
  const root = document.getElementById("root");
  if (root) {
    const config = resolveRootComponent(location.pathname);
    document.body.className = classSet({
      "is-authenticated": page.user,
      "is-staff": page.user && page.user.isCollaborator
    });
    ReactDOM.render(
      <>
        <ToastContainer position="top-right" toastClassName="c-toast" />
        {config.showHeader && <Header />}
        {React.createElement(config.component, page.props)}
        {config.showHeader && <Footer />}
      </>,
      root
    );
  }
});
