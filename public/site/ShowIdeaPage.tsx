import * as moment from "moment";
import * as React from "react";
import * as storage from "../storage";

import { User, Comment, Idea } from "../models";
import { CommentInput } from "../shared/CommentInput";
import { Gravatar, MultiLineText, IdeaResponse } from "../shared/Common";
import { SocialSignInButton } from "../shared/SocialSignInButton";
import { SupportCounter } from "../shared/SupportCounter";

import { Footer } from "../shared/Footer";
import { Header } from "../shared/Header";
import { IdeaInput } from "./IdeaInput";

export class ShowIdeaPage extends React.Component<{}, {}> {
    private user: User;
    private idea: Idea;
    private comments: Comment[];

    constructor(props: {}) {
      super(props);

      this.user = storage.getCurrentUser();
      this.idea = storage.get<Idea>("idea");
      this.comments = storage.get<Comment[]>("comments") || [];
    }

    public render() {

        const commentsList = this.comments.length ? this.comments.map((c) =>
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
                <MultiLineText text={ c.content } />
              </div>
            </div>
          </div>
        ) : <p>No comments yet.</p>;

        return <div>
                  <Header />
                  <div className="ui container">
                    <div className="ui items unstackable">
                      <div className="item">

                        <SupportCounter user={this.user} idea={this.idea} />

                        <div className="idea-header">
                          <h1 className="ui header">{ this.idea.title }</h1>

                          <p>
                            #{ this.idea.number } shared by &nbsp;
                            <Gravatar email={this.idea.user.email}/> <u>{this.idea.user.name}</u> &nbsp;
                            <span title={this.idea.createdOn}>{ moment(this.idea.createdOn).fromNow() }</span>
                          </p>
                        </div>
                      </div>
                    </div>

                    <MultiLineText text={ this.idea.description } />

                    <IdeaResponse status={this.idea.status}
                                  response={this.idea.response ? this.idea.response.text : undefined }
                                  createdOn={this.idea.response ? this.idea.response.createdOn : undefined }
                                  user={this.idea.user} />

                    <div className="ui comments">
                      <h3 className="ui dividing header">Comments</h3>
                      { commentsList }
                    </div>
                    <CommentInput idea={this.idea} />
                  </div>
                  <Footer />
               </div>;
    }
}
