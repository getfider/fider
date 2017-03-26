import * as moment from "moment";
import * as React from "react";
import * as storage from "../storage";

import { Comment, Idea } from "../models";
import { CommentInput } from "./CommentInput";
import { Gravatar, MultiLineText } from "./Common";
import { SocialSignInButton } from "./SocialSignInButton";

import { Footer } from "./Footer";
import { Header } from "./Header";
import { IdeaInput } from "./IdeaInput";

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
                    <h1 className="ui header">{ idea.title }</h1>

                    <p>{ idea.description }</p>

                    <p>
                      <Gravatar email={idea.user.email}/> <u>{idea.user.name}</u>
                      &nbsp;shared <u title={idea.createdOn}>{ moment(idea.createdOn).fromNow() }</u>
                    </p>

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
