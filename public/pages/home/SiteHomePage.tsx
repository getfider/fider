import * as React from 'react';
import { Idea, IdeaStatus, User } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { IdeaInput } from './IdeaInput';
import { IdeaFilter } from './IdeaFilter';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

import './SiteHomePage.scss';

interface SiteHomePageState {
  ideas: Idea[];
}

export class SiteHomePage extends React.Component<{}, SiteHomePageState> {
    private user: User;
    private allIdeas: Idea[];
    private filter: HTMLDivElement;

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.user = this.session.getCurrentUser();
        this.allIdeas = this.session.get<Idea[]>('ideas') || [];

        this.state = {
          ideas: IdeaFilter.defaultFilter(this.allIdeas)
        };
    }

    public render() {
        const ideasList = this.state.ideas.map((x) =>
          <div className="item" key={x.id}>
            <SupportCounter user={this.user} idea={x} />
            <div className="content">
              <a href={`/ideas/${x.number}/${x.slug}`} className="header">
                { x.title }
              </a>
              <MultiLineText className="description" text={ x.description } markdown={true} />
              <ShowIdeaResponse status={ x.status } response={ x.response } />
              <div className="extra">
                XXX comments
              </div>
            </div>
          </div>);

        const displayIdeas = (this.state.ideas.length > 0) ?
          <div className="ui divided unstackable items fdr-idea-list">
              { ideasList }
          </div>
          : <p>No ideas found for given filter.</p>;

        return <div className="SiteHomePage">
                  <Header />
                  <div className="ui container">
                    <h1 className="ui header">Welcome to our feedback forum!</h1>

                    <div className="ui grid">
                      <div className="six wide column">
                        <p>We'd love to hear what you're thinking about. This is the place for you to submit your feedback.</p>
                        <IdeaInput />
                      </div>
                      <div className="ten wide column">
                        <IdeaFilter filterChanged={(filter) => this.setState({ ideas: filter(this.allIdeas) }) } />
                        { displayIdeas }
                      </div>
                    </div>

                  </div>
                  <Footer />
               </div>;
    }
}
