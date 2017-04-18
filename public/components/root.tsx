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
                          <div className="ui small statistics">
                            <div className="statistic">
                              <div className="value">
                                { x.totalSupporters }
                              </div>
                              <div className="ui mini primary animated button">
                                <div className="visible content">Want</div>
                                <div className="hidden content">
                                  <i className="heart icon"></i>
                                </div>
                              </div>
                            </div>
                          </div>
                          <div className="content">
                            <a href={`/ideas/${x.number}`} className="header">
                              <i className="idea icon"></i> { x.totalSupporters } { x.title }
                            </a>
                            <div className="description">
                              <MultiLineText text={ x.description } />
                            </div>
                            <div className="extra">
                              #{ x.number } shared by <Gravatar email={x.user.email}/> <u>{x.user.name}</u>
                              <span title={x.createdOn}>{ moment(x.createdOn).fromNow() }</span>
                            </div>
                          </div>
                        </div>);

        const displayIdeas = (ideas.length > 0) ? <div>
                      <h3>Top Ideas</h3>
                      <div className="ui divided unstackable items">
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
