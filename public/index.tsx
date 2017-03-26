import * as React from "react";
import * as ReactDOM from "react-dom";

import { Root } from "./components/root";
import { ShowIdeaRoot } from "./components/show_idea_root";
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
