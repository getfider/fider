import * as React from "react";

/*
<!--
<div class="ui stackable inverted divided grid">
<div class="six wide column">
    <h4 class="ui inverted header">Trending ideas</h4>
    <div class="ui inverted link list">
    <a href="#" class="item">Link #1</a>
    <a href="#" class="item">Link #2</a>
    <a href="#" class="item">Link #3</a>
    <a href="#" class="item">Link $3</a>
    </div>
</div>
<div class="six wide column">
    <h4 class="ui inverted header">Latest activity</h4>
    <div class="ui inverted link list">
    <a href="#" class="item">Someone did something</a>
    <a href="#" class="item">New idea was posted</a>
    <a href="#" class="item">Someone commented on idea #1</a>
    <a href="#" class="item">Idea #2 changed to Shipped</a>
    </div>
</div>
<div class="four wide column">
    <h4 class="ui inverted header">Footer Header</h4>
    <p>Extra space for a call to action inside the footer that could help re-engage users.</p>
</div>
</div>
-->
*/

export class Footer extends React.Component<{}, {}> {
    public render() {
        return  <div id="footer" className="ui inverted vertical footer segment">
                    <div className="ui center aligned container">
                        <div className="ui inverted section divider"></div>
                        <div className="ui horizontal inverted small divided link list">
                            <span>Powered by</span>
                            <a className="item" target="_blank" href="http://getfider.com/">Fider</a>
                        </div>
                    </div>
                </div>;
    }
}
