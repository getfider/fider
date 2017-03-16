import * as React from "react";
import { IdeaInput } from "./IdeaInput";
import { Header } from "./Header";
import { Footer } from "./Footer";

const claims = (window as any)._claims;
const ideas:any[] = (window as any)._ideas || [];

export class Root extends React.Component<{}, {}> {
    render() {

        const ideasList = ideas.map(x => 
                        <div className="item" key={x.id}>
                          <div className="content">
                            <a className="header">{ x.title }</a>
                            <div className="meta">
                              <span className="cinema">Union Square 14</span>
                            </div>
                            <div className="description">
                              <p>{ x.title }</p>
                            </div>
                            <div className="extra">
                              <div className="ui label">IMAX</div>
                              <div className="ui label">Additional Languages</div>
                            </div>
                          </div>
                        </div>);

        const displayIdeas = (ideas.length > 0) ? <div>
                      <h3>Top Ideas</h3>
                      <div className="ui divided items">
                        { ideasList }
                      </div> 
                    </div> 
                    : <p>No ideas have been shared yet.</p>

        return <div>
                  <Header />
                  <div className="ui container">
                    <h1 className="ui header">Welcome to our feedback forum!</h1>

                    <h3>I want ...</h3>

                    <IdeaInput />

                    <div className="ui container ideas-list">
                      { displayIdeas }
                    </div>
                  </div>
                  <Footer />
               </div>;
    }
}