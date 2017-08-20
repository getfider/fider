import * as React from 'react';
import { Idea, IdeaStatus, User, Tenant } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { IdeaInput } from './IdeaInput';
import { IdeaFilter, IdeaFilterFunction } from './IdeaFilter';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

import './SiteHomePage.scss';

interface SiteHomePageState {
  ideas: Idea[];
}

export class SiteHomePage extends React.Component<{}, SiteHomePageState> {
    private user: User;
    private tenant: Tenant;
    private allIdeas: Idea[];
    private filter: HTMLDivElement;
    private activeFilter: string;

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.user = this.session.getCurrentUser();
        this.tenant = this.session.getCurrentTenant();
        this.allIdeas = this.session.get<Idea[]>('ideas') || [];

        this.activeFilter = window.location.hash.substring(1);
        this.state = {
          ideas: IdeaFilter.getFilter(this.activeFilter)(this.allIdeas)
        };
    }

    private filterChanged(name: string, filter: IdeaFilterFunction) {
      window.location.hash = `#${name}`;
      this.setState({
        ideas: filter(this.allIdeas)
      });
    }

    public render() {
        const ideasList = this.state.ideas.map((x) =>
          <div className="item" key={x.id}>
            <SupportCounter user={this.user} idea={x} />
            <div className="content">
              { x.totalComments > 0 && <div className="info right">
                { x.totalComments } <i className="comments outline icon"/>
              </div> }
              <a href={`/ideas/${x.number}/${x.slug}`} className="header">
                { x.title }
              </a>
              <MultiLineText className="description" text={ x.description } markdown={true} />
              <ShowIdeaResponse status={ x.status } response={ x.response } />
              <div className="info">
                <Moment date={x.createdOn} />
              </div>
            </div>
          </div>);

        const displayIdeas = (this.state.ideas.length > 0) ?
          <div className="ui divided unstackable items fdr-idea-list">
              { ideasList }
          </div>
          : <p>No ideas found for given filter.</p>;

        const welcomeMessage = this.tenant.welcomeMessage ||
        `## Welcome to our feedback forum!

We'd love to hear what you're thinking about. This is the place for you to submit your feedback.`;

        return <div className="SiteHomePage">
                  <Header />
                  <div className="page ui container">

                    <div className="ui grid stackable">
                      <div className="six wide column">
                        <MultiLineText className="welcome-message" text={ welcomeMessage } markdown={true} />
                        <IdeaInput placeholder={this.tenant.invitation || 'Tell us your ideas'} />
                      </div>
                      <div className="ten wide column">
                        {
                          this.allIdeas.length === 0
                          ? <div className="center">
                              <p><i className="icon lightbulb" aria-hidden="true"></i></p>
                              <p>It's lonely out here. Start by sharing an idea!</p>
                            </div>
                          : <div>
                              <IdeaFilter activeFilter={ this.activeFilter } filterChanged={ this.filterChanged.bind(this) } />
                              { displayIdeas }
                            </div>
                        }
                      </div>
                    </div>

                  </div>
                  <Footer />
               </div>;
    }
}
