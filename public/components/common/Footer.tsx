import "./Footer.scss"

import React from "react"

export const Footer = () => {
  return (
    <div id="c-footer">
      <div className="container">
        <a className="l-powered" target="_blank" rel="noopener" href="https://getfider.com/">
          <img width="20" height="20" src="https://getfider.com/images/logo-100x100.png" alt="Fider" />
          <span>Powered by Fider</span>
        </a>
      </div>
    </div>
  )
}
