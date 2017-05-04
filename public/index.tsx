import * as React from "react";
import * as ReactDOM from "react-dom";

import { Root } from "./components/Root";
import { ShowIdeaRoot } from "./components/ShowIdeaRoot";
import { setup } from "./storage";

import "./css/main.css";

const pathRegex = [
    { regex: new RegExp("^\/$"), component: <Root /> },
    { regex: new RegExp("^\/ideas\/\\d+$"), component: <ShowIdeaRoot /> },
];

const resolveRootComponent = (path: string): JSX.Element => {
    for (const entry of pathRegex) {
        if (entry.regex.test(path)) {
            return entry.component;
        }
    }

    return <div />;
};

setup();

document.addEventListener("DOMContentLoaded", () => {
  ReactDOM.render(
      resolveRootComponent(location.pathname),
      document.getElementById("root")
  );
});
