import * as React from 'react';
import { Idea, IdeaStatus, User } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { IdeaInput } from '@fider/components/IdeaInput';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';

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
          ideas: this.allIdeas.filter(this.filterers.recent).sort(this.sorterers.recent)
        };
    }

    private filterers: {[key: string]: (idea: Idea) => boolean} = {
      'recent': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
      'most-wanted': (idea: Idea) => idea.status !== IdeaStatus.Completed.value && idea.status !== IdeaStatus.Declined.value,
      'completed': (idea: Idea) => idea.status === IdeaStatus.Completed.value,
      'declined': (idea: Idea) => idea.status === IdeaStatus.Declined.value
    };

    private sorterers: {[key: string]: (left: Idea, right: Idea) => number} = {
      'recent': (left: Idea, right: Idea) => new Date(right.createdOn).getTime() - new Date(left.createdOn).getTime(),
      'most-wanted': (left: Idea, right: Idea) => right.totalSupporters - left.totalSupporters,
      'completed': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime(),
      'declined': (left: Idea, right: Idea) => new Date(left.response.respondedOn).getTime() - new Date(left.response.respondedOn).getTime()
    };

    private applyFilter(value: string) {
      this.setState({
        ideas: this.allIdeas.filter(this.filterers[value]).sort(this.sorterers[value])
      });
    }

    public componentDidMount() {
        $(this.filter).dropdown({
          onChange: (value: string) => {
            this.applyFilter(value);
          }
        });
    }

    public render() {
        const ideasList = this.state.ideas.map((x) =>
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

        const displayIdeas = (this.state.ideas.length > 0) ?
          <div className="ui divided unstackable items">
              { ideasList }
          </div>
          : <p>No ideas found for given filter.</p>;

        return <div>
                  <Header />
                  <div className="ui container">
                    <h1 className="ui header">Welcome to our feedback forum!</h1>

                    <h3>I want ...</h3>

                    <IdeaInput />

                      <div className="ui container fdr-idea-list">
                        <h4 className="ui header">
                          <i className="filter icon"></i>
                          <div className="content">
                            Showing {' '}
                            <div className="ui inline dropdown" ref={(e) => this.filter = e!}>
                              <div className="text">recent ideas</div>
                              <i className="dropdown icon"></i>
                              <div className="menu">
                                <div className="header">What do you want to see?</div>
                                <div className="item" data-value="recent" data-text="recent ideas">Recent</div>
                                <div className="item" data-value="most-wanted" data-text="most wanted ideas">Most Wanted</div>
                                <div className="item" data-value="completed" data-text="completed ideas">Completed</div>
                                <div className="item" data-value="declined" data-text="declined ideas">Declined</div>
                              </div>
                            </div>
                          </div>
                        </h4>
                        { displayIdeas }
                      </div>
                  </div>
                  <Footer />
               </div>;
    }
}
