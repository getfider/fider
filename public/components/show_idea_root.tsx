import * as moment from "moment";
import * as React from "react";
import * as storage from "../storage";

import { Comment, Idea } from "../models";
import { CommentInput } from "./comment_input";
import { Gravatar, MultiLineText } from "./common";
import { SocialSignInButton } from "./SocialSignInButton";
import { SupportCounter } from "./SupportCounter";

import { Footer } from "./footer";
import { Header } from "./header";
import { IdeaInput } from "./idea_input";

export class ShowIdeaRoot extends React.Component<{}, {}> {
    public render() {
        const user = storage.getCurrentUser();
        const idea = storage.get<Idea>("idea");
        const comments = storage.get<Comment[]>("comments") || [];

        const commentsList = comments.length ? comments.map((c) =>
          <div className="comment">
            <a className="avatar">
              <Gravatar email={c.user.email} />
            </a>
            <div className="content">
              <span className="author">{ c.user.name }</span>
              <div className="metadata">
                <span className="date" title={c.createdOn}>{ moment(c.createdOn).fromNow() }</span>
              </div>
              <div className="text">
                { c.content }
              </div>
            </div>
          </div>
        ) : <p>No comments yet.</p>;

        return <div>
                  <Header />
                  <div className="ui container">
                    <div className="ui items unstackable">
                      <div className="item">

                        <SupportCounter user={user} idea={idea} />

                        <div className="idea-header">
                          <h1 className="ui header">{ idea.title }</h1>

                          <p>
                            #{ idea.number } shared by &nbsp;
                            <Gravatar email={idea.user.email}/> <u>{idea.user.name}</u> &nbsp;
                            <span title={idea.createdOn}>{ moment(idea.createdOn).fromNow() }</span>
                          </p>
                        </div>
                      </div>
                    </div>

                    <p>{ idea.description }</p>

                    <div className="ui comments">
                      <h3 className="ui dividing header">Comments</h3>
                      { commentsList }
                    </div>
                    <CommentInput idea={idea} />
                  </div>
                  <Footer />
               </div>;
    }
}
