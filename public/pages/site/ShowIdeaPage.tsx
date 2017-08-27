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

import './ShowIdeaPage.scss';

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
                · <Moment date={c.createdOn} />
              </div>
              <div className="text">
                <MultiLineText text={ c.content } />
              </div>
            </div>
          </div>
        ) : <p>No comments have been added yet.</p>;

        return <div className="ShowIdeaPage">
                  <Header />
                  <div className="page ui container">
                    <div className="ui stackable grid container">
                      <div className="twelve wide column">
                        <div className="ui items unstackable">
                          <div className="item">
                            <SupportCounter user={this.user} idea={this.idea} />

                            <div className="idea-header">
                              <h1 className="ui header">{ this.idea.title }</h1>

                              <span className="info">
                                Shared <Moment date={this.idea.createdOn} /> by <Gravatar name={ this.idea.user.name } hash={ this.idea.user.gravatar }/> <b>{ this.idea.user.name }</b>
                              </span>
                            </div>
                          </div>
                        </div>

                        <span className="subtitle">Description</span>
                        {
                          this.idea.description
                          ? <MultiLineText className="description" text={ this.idea.description } markdown={true} />
                          : <p className="description">This idea has no description.</p>
                        }

                        <ShowIdeaResponse status={ this.idea.status } response={ this.idea.response } />

                      </div>

                      {
                        this.session.isStaff() &&
                        <div className="four wide column">
                          <span className="subtitle">Actions</span>
                          <br /><br />
                          <ResponseForm idea={ this.idea } />
                        </div>
                      }

                      <div className="sixteen wide column">
                          <div className="ui comments">
                            <span className="subtitle">Discussion</span>
                            { commentsList }
                            <CommentInput idea={this.idea} />
                          </div>
                      </div>
                    </div>

                  </div>
                  <Footer />
               </div>;
    }
}
