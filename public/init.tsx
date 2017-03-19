// tslint:disable-next-line:no-unused-variable
import * as React from "react";
import * as ReactDOM from "react-dom";

import { Root } from "./components/Root";

export default function init() {
    ReactDOM.render(
        location.pathname === "/" ? <Root /> : <div />,
        document.getElementById("root")
    );
}
