import * as React from 'react';
const logo = require('@fider/images/logo.png');

export const Footer = () => {
        return  <div id="footer" className="ui vertical footer segment">
                    <div className="ui center aligned container">
                        <div className="ui inverted section divider"></div>
                        <div className="ui horizontal small divided link list">
                            <div id="powered-by">
                                <p>
                                    <span>Powered by </span>
                                    <b>
                                        <a className="item" target="_blank" href="http://getfider.com/">Fider</a>
                                    </b>
                                </p>
                                <p><img src={logo} /></p>
                            </div>
                        </div>
                    </div>
                </div>;
};
