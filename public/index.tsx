import * as React from "react";
import * as ReactDOM from "react-dom";
import { resolveRootComponent } from "@fider/router";
import { Header, Footer } from "@fider/components/common";
import { analytics, classSet } from "@fider/services";
import { ToastContainer } from "react-toastify";

import "semantic-ui-css/components/reset.min.css";
import "semantic-ui-css/components/icon.min.css";
import "semantic-ui-css/components/dropdown.min.css";
import "semantic-ui-css/components/transition.min.css";

import "react-toastify/dist/ReactToastify.css";
import "@fider/assets/styles/main.scss";

window.props = {} as any;
window.set = (key: string, value: any): void => {
  window.props[key] = value;
};

window.addEventListener("error", (evt: ErrorEvent) => {
  analytics.error(evt.error);
});

document.addEventListener("DOMContentLoaded", () => {
  const root = document.getElementById("root");
  if (root) {
    const config = resolveRootComponent(location.pathname);
    document.body.className = classSet({
      "is-authenticated": window.props.user,
      "is-staff": window.props.user && window.props.user.isCollaborator
    });
    ReactDOM.render(
      <>
        <ToastContainer position="top-right" toastClassName="c-toast" />
        {config.showHeader && React.createElement(Header, window.props)}
        {React.createElement(config.component, window.props)}
        {config.showHeader && React.createElement(Footer, window.props)}
      </>,
      root
    );
  }
});
