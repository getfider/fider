import * as React from 'react';
const logo = require('@fider/assets/images/logo-small.png');

export const Footer = () => {
  return (
    <div id="footer" className="ui vertical footer segment">
      <div className="ui center aligned container">
        <div className="ui inverted section divider"/>
        <div className="ui horizontal small divided link list">
          <div id="powered-by">
            <a className="item" target="_blank" href="http://getfider.com/">
              <span>Powered by Fider</span>
              <img src={logo} />
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};
