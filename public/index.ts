import * as $ from "jquery";
import init from "./init";
import * as storage from "./storage";

import "./css/main.css";

storage.setup();

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
