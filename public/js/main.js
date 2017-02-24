document.addEventListener("DOMContentLoaded", function() {
  var input = document.querySelector("#new-idea input");
  var submit = document.querySelector("#new-idea button");

  submit.style.display = 'none';

  input.addEventListener("keyup", function() {
    if (this.value) {
      submit.style.display = 'block';
    } else {
      submit.style.display = 'none';
    }
  });
});