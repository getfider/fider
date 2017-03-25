import * as moment from "moment";
import * as React from "react";
import { Comment, Idea } from "../models";
import * as storage from "../storage";
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

        const addComment = user ? <form className="ui reply form">
          <div className="field">
            <textarea></textarea>
          </div>
          <div className="ui blue labeled submit icon button">
            <i className="icon edit"></i> Add Comment
          </div>
        </form> :
        <div className="ui message">
          <div className="header">
            Please log in before leaving a comment
          </div>
          <p>This only takes a second and you'll be good to go!</p>
          <SocialSignInButton provider="facebook" small={true} />
          <SocialSignInButton provider="google" small={true} />
        </div>;

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
                    { addComment }
                  </div>
                  <Footer />
               </div>;
    }
}
