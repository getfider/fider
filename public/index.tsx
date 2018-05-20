import * as React from "react";
import * as ReactDOM from "react-dom";
import { resolveRootComponent } from "@fider/router";
import { Header, Footer } from "@fider/components/common";
import { analytics } from "@fider/services";
import { ToastContainer } from "react-toastify";

import "semantic-ui-css/components/reset.min.css";
import "semantic-ui-css/components/icon.min.css";
import "semantic-ui-css/components/dropdown.min.css";
import "semantic-ui-css/components/transition.min.css";

import "semantic-ui-css/components/site.min.css";
import "semantic-ui-css/components/form.min.css";
import "semantic-ui-css/components/input.min.css";
import "semantic-ui-css/components/label.min.css";
import "semantic-ui-css/components/message.min.css";
import "semantic-ui-css/components/checkbox.min.css";

import "react-toastify/dist/ReactToastify.css";
import "@fider/assets/styles/main.scss";

const w = window as any;
w.props = {};
w.set = (key: string, value: any): void => {
  w.props[key] = value;
};

window.addEventListener("error", (evt: ErrorEvent) => {
  analytics.error(evt.error);
});

document.addEventListener("DOMContentLoaded", () => {
  const root = document.getElementById("root");
  if (root) {
    const config = resolveRootComponent(location.pathname);
    ReactDOM.render(
      <>
        <ToastContainer position="top-right" toastClassName="c-toast" />
        {config.showHeader && React.createElement(Header, w.props)}
        {React.createElement(config.component, w.props)}
        {config.showHeader && React.createElement(Footer, w.props)}
      </>,
      root
    );
  }
});
