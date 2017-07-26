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
              <div className="extra">
                shared <Moment date={x.createdOn} />
                <div style={{float: 'right'}}>{ x.totalComments } <i className="comments icon"/></div>
              </div>
              <ShowIdeaResponse status={ x.status } response={ x.response } />
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

                    <div className="ui grid stackable">
                      <div className="six wide column">
                        <h2 className="ui header">Welcome to our feedback forum!</h2>
                        <p>We'd love to hear what you're thinking about. This is the place for you to submit your feedback.</p>
                        <IdeaInput />
                      </div>
                      <div className="ten wide column">
                        {
                          this.allIdeas.length === 0
                          ? <div className="center">
                              <p><i className="icon lightbulb" aria-hidden="true"></i></p>
                              <p>It's lonely out here. Start by sharing an idea!</p>
                            </div>
                          : <div>
                              <IdeaFilter filterChanged={(filter) => this.setState({ ideas: filter(this.allIdeas) }) } />
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
