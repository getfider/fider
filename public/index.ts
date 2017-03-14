import * as $ from "jquery";

import init from "./init"

(<any>window).$ = (<any>window).jQuery = $;
import "semantic-ui-css/semantic.min.css";
import "./css/main.css";
import "semantic-ui-css/semantic.min.js";

$(document).ready(() => {
    init();

    $(".ui.dropdown").dropdown();

    $('.menu .signin').popup({
        inline: true,
        hoverable: true,
        popup: '#user-popup',
        position : 'bottom right',
        delay: {
        show: 300
        }
    });

    $('#user-popup a').click(function() {
        $(this).addClass("loading");
    });
    
    var input = <HTMLInputElement>document.querySelector("#new-idea-input");
    var submit = <HTMLElement>document.querySelector("#new-idea-submit");

    submit.style.display = 'none';

    input.addEventListener("keyup", function() {
        if (this.value) {
            submit.style.display = 'block';
        } else {
            submit.style.display = 'none';
        }
    });
});

