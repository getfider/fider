import * as React from "react";
const logo = require("@fider/assets/images/logo-small.png");

export const Footer = () => {
  return (
    <div id="footer">
      <div id="powered-by" className="ui container">
        <a target="_blank" href="https://getfider.com/">
          <img src={logo} alt="Fider" />
          <span>Powered by Fider</span>
        </a>
      </div>
    </div>
  );
};
