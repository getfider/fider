import * as React from "react";
import * as ReactDOM from "react-dom";

import { IdeaInput } from "./components/IdeaInput";

export default function init() {
    ReactDOM.render(
        <IdeaInput />,
        document.getElementById("root")
    );
}