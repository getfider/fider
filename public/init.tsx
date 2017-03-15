import * as React from "react";
import * as ReactDOM from "react-dom";

import { Root } from "./components/Root";

export default function init() {
    ReactDOM.render(
        <Root />,
        document.getElementById("root")
    );
}