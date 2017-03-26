import * as moment from "moment";
import * as React from "react";
import { Idea } from "../models";
import * as storage from "../storage";
import { Gravatar, MultiLineText } from "./common";

import { Footer } from "./footer";
import { Header } from "./header";
import { IdeaInput } from "./idea_input";

export class Root extends React.Component<{}, {}> {
    public render() {
        const ideas = storage.get<Idea[]>("ideas") || [];

        const ideasList = ideas.map((x) =>
                        <div className="item" key={x.id}>
                          <div className="content">
                            <a href={`/ideas/${x.id}`} className="header">
                              <i className="idea icon"></i>{ x.title }
                            </a>
                            <div className="description">
                              <MultiLineText text={ x.description } />
                            </div>
                            <div className="extra">
                              <Gravatar email={x.user.email}/> <u>{x.user.name}</u>
                              shared <u title={x.createdOn}>{ moment(x.createdOn).fromNow() }</u>
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
