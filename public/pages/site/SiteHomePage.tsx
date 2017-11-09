import * as React from 'react';
import { Idea, IdeaStatus, CurrentUser, Tenant } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { IdeaInput } from '@fider/components/IdeaInput';
import { IdeaFilter, IdeaFilterFunction } from '@fider/components/IdeaFilter';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

import './SiteHomePage.scss';

const defaultShowCount = 20;

interface SiteHomePageState {
  ideas: Idea[];
  showCount: number;
}

const ListIdeaItem = (props: { idea: Idea, user?: CurrentUser }) => {
  return <div className="item">
            <SupportCounter user={props.user} idea={props.idea} />
            <div className="content">
              { props.idea.totalComments > 0 && <div className="info right">
                { props.idea.totalComments } <i className="comments outline icon"/>
              </div> }
              <a className="title" href={`/ideas/${props.idea.number}/${props.idea.slug}`}>
                { props.idea.title }
              </a>
              <MultiLineText className="description" text={ props.idea.description } style="simple" />
              <ShowIdeaResponse status={ props.idea.status } response={ props.idea.response } />
            </div>
          </div>;
};

export class SiteHomePage extends React.Component<{}, SiteHomePageState> {
    private user?: CurrentUser;
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
          ideas: IdeaFilter.getFilter(this.activeFilter)(this.allIdeas),
          showCount: defaultShowCount
        };
    }

    private filterChanged(name: string, filter: IdeaFilterFunction) {
      window.location.hash = `#${name}`;
      this.setState({
        ideas: filter(this.allIdeas),
        showCount: defaultShowCount
      });
    }

    private showMore(event: React.SyntheticEvent<TouchEvent>): void {
      event.preventDefault();
      this.setState({
        showCount: this.state.showCount + defaultShowCount
      });
    }

    public render() {
        const ideasToList = this.state.ideas.slice(0, this.state.showCount);

        const displayIdeas = (ideasToList.length > 0) ?
          <div className="ui divided unstackable items fdr-idea-list">
              { ideasToList.map((x) => <ListIdeaItem key={x.id} user={this.user} idea={x} />) }
              {
                this.state.ideas.length > this.state.showCount &&
                <h5 className="ui blue header show-more"
                    onTouchEnd={ this.showMore.bind(this) }
                    onClick={ this.showMore.bind(this) }>
                  View { this.state.ideas.length - this.state.showCount } more ideas
                </h5>
              }
          </div>
          : <p>No ideas found for given filter.</p>;

        const welcomeMessage = this.tenant.welcomeMessage ||
        `## Welcome to our feedback forum!

We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas.`;

        return <div>
                  <Header />
                  <div className="page ui container">

                    <div className="ui grid stackable">
                      <div className="six wide column">
                        <MultiLineText className="welcome-message" text={ welcomeMessage } style="full" />
                        <IdeaInput placeholder={this.tenant.invitation || 'I suggest you...'} />
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
