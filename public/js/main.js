document.addEventListener("DOMContentLoaded", function() {
  $('.ui.dropdown').dropdown();
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
  
  var input = document.querySelector("#new-idea-input");
  var submit = document.querySelector("#new-idea-submit");

  submit.style.display = 'none';

  input.addEventListener("keyup", function() {
    if (this.value) {
      submit.style.display = 'block';
    } else {
      submit.style.display = 'none';
    }
  });
});