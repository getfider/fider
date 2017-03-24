// tslint:disable-next-line:no-unused-variable
import * as React from "react";
import * as ReactDOM from "react-dom";

import { Root } from "./components/Root";
import { ShowIdeaRoot } from "./components/ShowIdeaRoot";

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

export default function init() {
    ReactDOM.render(
        resolveRootComponent(location.pathname),
        document.getElementById("root")
    );
}
