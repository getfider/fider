import * as React from 'react';
import { Idea, User } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { IdeaInput } from '@fider/components/IdeaInput';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

export class SiteHomePage extends React.Component<{}, {}> {
    private user: User;
    private ideas: Idea[];

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.user = this.session.getCurrentUser();
        this.ideas = this.session.get<Idea[]>('ideas') || [];
    }

    public render() {
        const ideasList = this.ideas.map((x) =>
                        <div className="item" key={x.id}>
                          <SupportCounter user={this.user} idea={x} />
                          <div className="content">
                            <a href={`/ideas/${x.number}/${x.slug}`} className="header">
                              { x.title }
                            </a>
                            <div className="description">
                              <MultiLineText text={ x.description } markdown={true} />
                            </div>
                            <div className="extra">
                              #{ x.number } shared by <Gravatar hash={x.user.gravatar}/> <u>{x.user.name}</u>
                              <Moment date={x.createdOn} />
                            </div>
                            <ShowIdeaResponse status={ x.status } response={ x.response } />
                          </div>
                        </div>);

        const displayIdeas = (this.ideas.length > 0) ? <div>
                      <h3>Recent Ideas</h3>
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

                    <div className="ui container fdr-idea-list">
                      { displayIdeas }
                    </div>
                  </div>
                  <Footer />
               </div>;
    }
}
