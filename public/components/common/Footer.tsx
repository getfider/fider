import "./Footer.scss";

import React from "react";
import { PrivacyPolicy, TermsOfService } from "@fider/components";
import { useFider } from "@fider/hooks";

export const Footer = () => {
  const fider = useFider();

  return (
    <div id="c-footer">
      <div className="container">
        {fider.settings.hasLegal && (
          <div className="l-links">
            <PrivacyPolicy />
            &middot;
            <TermsOfService />
          </div>
        )}
        <a className="l-powered" target="_blank" rel="noopener" href="https://getfider.com/">
          <img width="20" height="20" src="https://getfider.com/images/logo-100x100.png" alt="Fider" />
          <span>Powered by Fider</span>
        </a>
      </div>
    </div>
  );
};
