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
      setTitle(`${this.idea.title} · ${document.title}`);
    }

    public render() {

        const commentsList = this.comments.length ? this.comments.map((c) =>
          <div className="comment">
            <Gravatar name={c.user.name} hash={c.user.gravatar} />
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
        ) : <p>There are no comments yet.</p>;

        return <div>
                  <Header />
                  <div className="page ui container">
                    <div className="ui items unstackable">
                      <div className="item">

                        <SupportCounter user={this.user} idea={this.idea} />

                        <div className="idea-header">
                          <h1 className="ui header">{ this.idea.title }</h1>

                          <p>
                            <Gravatar name={this.idea.user.name} hash={this.idea.user.gravatar}/> <b>{this.idea.user.name}</b> &nbsp;
                            <span className="info">
                              <Moment date={this.idea.createdOn} />
                            </span>
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
                      <h3 className="ui dividing header">Discussion</h3>
                      { commentsList }
                      <CommentInput idea={this.idea} />
                    </div>
                  </div>
                  <Footer />
               </div>;
    }
}
