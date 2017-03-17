import * as moment from "moment";
import * as React from "react";

import { Footer } from "./Footer";
import { Header } from "./Header";
import { IdeaInput } from "./IdeaInput";

const claims = (window as any)._claims;
const ideas: any[] = (window as any)._ideas || [];

export class Root extends React.Component<{}, {}> {
    public render() {

        const ideasList = ideas.map((x: any) =>
                        <div className="item" key={x.id}>
                          <div className="content">
                            <a className="header">{ x.title }</a>
                            <div className="description">
                              <p>{ x.title }</p>
                            </div>
                            <div className="extra">
                              <i className="calendar icon"></i>
                              shared { moment(x.createdOn).fromNow() }
                            </div>
                          </div>
                        </div>);

        const displayIdeas = (ideas.length > 0) ? <div>
                      <h3>Top Ideas</h3>
                      <div className="ui divided items">
                        { ideasList }
                      </div> 
                    </div>
                    : <p>No ideas have been shared yet.</p>;

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
