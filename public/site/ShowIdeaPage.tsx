import * as moment from 'moment';
import * as React from 'react';
import { getCurrentUser, get, getArray, isStaff } from '../storage';

import { User, Comment, Idea } from '../models';
import { CommentInput } from '../shared/CommentInput';
import { Gravatar, MultiLineText, ShowIdeaResponse } from '../shared/Common';
import { SocialSignInButton } from '../shared/SocialSignInButton';
import { SupportCounter } from '../shared/SupportCounter';
import { ResponseForm } from '../shared/ResponseForm';

import { Footer } from '../shared/Footer';
import { Header } from '../shared/Header';
import { IdeaInput } from './IdeaInput';

export class ShowIdeaPage extends React.Component<{}, {}> {
    private user: User;
    private idea: Idea;
    private comments: Comment[];

    constructor(props: {}) {
      super(props);

      this.user = getCurrentUser();
      this.idea = get<Idea>('idea');
      this.comments = getArray<Comment>('comments');
      document.title = `${this.idea.title} · Idea #${this.idea.number} · ${document.title}`;
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

        const respond = <div className="ui blue labeled submit icon button false">
                          <i className="icon announcement"></i>Respond
                        </div>;

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
                    <ShowIdeaResponse status={ this.idea.status } response={ this.idea.response } />

                    {
                      isStaff() &&
                      <div>
                        <div className="ui hidden divider"></div>
                        <ResponseForm idea={ this.idea } />
                      </div>
                    }

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
