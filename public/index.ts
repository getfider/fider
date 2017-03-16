import * as $ from "jquery";

import init from "./init";

(<any> window).$ = (<any> window).jQuery = $;
import "semantic-ui-css/semantic.min.css";
import "semantic-ui-css/semantic.min.js";
import "./css/main.css";

$(document).ready(() => {
    init();

    $(".ui.dropdown").dropdown();

    $(".menu .signin").popup({
        inline: true,
        hoverable: true,
        popup: "#user-popup",
        position : "bottom right",
        delay: {
            show: 300
        }
    });
});
