import * as moment from 'moment';
import * as React from 'react';
import { Idea, User } from '../models';
import { Gravatar, MultiLineText, ShowIdeaResponse } from '../shared/Common';
import { SupportCounter } from '../shared/SupportCounter';

import { Footer } from '../shared/Footer';
import { Header } from '../shared/Header';
import { IdeaInput } from './IdeaInput';

import { inject, injectables } from '../di';
import { Session } from '../services/Session';

export class HomePage extends React.Component<{}, {}> {
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
                              <i className="idea icon"></i> { x.title }
                            </a>
                            <div className="description">
                              <MultiLineText text={ x.description } />
                            </div>
                            <div className="extra">
                              #{ x.number } shared by <Gravatar email={x.user.email}/> <u>{x.user.name}</u>
                              <span title={x.createdOn}>{ moment(x.createdOn).fromNow() }</span>
                            </div>
                            <ShowIdeaResponse status={ x.status } response={ x.response } />
                          </div>
                        </div>);

        const displayIdeas = (this.ideas.length > 0) ? <div>
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

                    <div className="ui container fdr-idea-list">
                      { displayIdeas }
                    </div>
                  </div>
                  <Footer />
               </div>;
    }
}
