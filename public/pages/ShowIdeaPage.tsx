import * as React from 'react';

import { User, Comment, Idea } from '@fider/models';
import { setTitle } from '@fider/utils/page';

import { CommentInput } from '@fider/components/CommentInput';
import { ResponseForm } from '@fider/components/ResponseForm';
import { SupportCounter } from '@fider/components/SupportCounter';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { Button, Gravatar, Moment, MultiLineText, Footer, Header, SocialSignInButton } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

export class ShowIdeaPage extends React.Component<{}, {}> {
    private user: User;
    private idea: Idea;
    private comments: Comment[];

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
      super(props);

      this.user = this.session.getCurrentUser();
      this.idea = this.session.get<Idea>('idea');
      this.comments = this.session.getArray<Comment>('comments');
      setTitle(`${this.idea.title} · Idea #${this.idea.number} · ${document.title}`);
    }

    public render() {

        const commentsList = this.comments.length ? this.comments.map((c) =>
          <div className="comment">
            <a className="avatar">
              <Gravatar hash={c.user.gravatar} />
            </a>
            <div className="content">
              <span className="author">{ c.user.name }</span>
              <div className="metadata">
                <Moment date={c.createdOn} />
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
                            <Gravatar hash={this.idea.user.gravatar}/> <u>{this.idea.user.name}</u> &nbsp;
                            <Moment date={this.idea.createdOn} />
                          </p>
                        </div>
                      </div>
                    </div>

                    <MultiLineText className="description" text={ this.idea.description } markdown={true} />
                    <ShowIdeaResponse status={ this.idea.status } response={ this.idea.response } />

                    {
                      this.session.isStaff() &&
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
