import * as React from 'react';
import { Idea, Tag, IdeaStatus, CurrentUser, Tenant } from '@fider/models';
import { Gravatar, MultiLineText, Moment, Header, Footer } from '@fider/components/common';
import { ShowIdeaResponse } from '@fider/components/ShowIdeaResponse';
import { SupportCounter } from '@fider/components/SupportCounter';
import { ShowTag } from '@fider/components/ShowTag';
import { IdeaInput } from './IdeaInput';
import { IdeaFilter, IdeaFilterFunction } from './IdeaFilter';

import { inject, injectables } from '@fider/di';
import { Session } from '@fider/services';
import { getBaseUrl, getQueryString } from '@fider/utils/page';

import { TagsFilter } from './TagsFilter';

import './HomePage.scss';

const defaultShowCount = 20;

interface HomePageState {
  ideas: Idea[];
  showCount: number;
  searching: boolean;
  search?: string;
  activeFilter: string;
}

const EmptyList = () => {
  return (
    <div className="center">
      <p><i className="icon lightbulb" aria-hidden="true" /></p>
      <p>It's lonely out here. Start by sharing an idea!</p>
    </div>
  );
};

const ListIdeaItem = (props: { idea: Idea, user?: CurrentUser, tags: Tag[] }) => {
  return (
    <div className="item">
      <SupportCounter user={props.user} idea={props.idea} />
      <div className="content">
        {
          props.idea.totalComments > 0 &&
          <div className="info right">
            {props.idea.totalComments} <i className="comments outline icon"/>
          </div>
        }
        <a className="title" href={`/ideas/${props.idea.number}/${props.idea.slug}`}>
          {props.idea.title}
        </a>
        <MultiLineText className="description" text={props.idea.description} style="simple" />
        <ShowIdeaResponse status={props.idea.status} response={props.idea.response} />
        {
          props.tags.map((tag) => (
            <ShowTag key={tag.id} size="tiny" tag={tag} />
          ))
        }
      </div>
    </div>
  );
};

export class HomePage extends React.Component<{}, HomePageState> {
    private user?: CurrentUser;
    private tenant: Tenant;
    private allIdeas: Idea[];
    private allTags: Tag[];
    private filter: HTMLDivElement;

    @inject(injectables.Session)
    public session: Session;

    constructor(props: {}) {
        super(props);
        this.user = this.session.getCurrentUser();
        this.tenant = this.session.getCurrentTenant();
        this.allIdeas = this.session.getArray<Idea>('ideas');
        this.allTags = this.session.getArray<Tag>('tags');

        const search = getQueryString('q');
        const activeFilter = window.location.hash.substring(1);
        this.state = {
          ideas: this.filterIdeas(activeFilter, search),
          showCount: defaultShowCount,
          activeFilter,
          searching: !!search,
          search
        };
    }

    private containsAll(str: string, substrings: string[]): boolean {
        for (let i = 0; i !== substrings.length; i++) {
           if (str.indexOf(substrings[i]) === - 1) {
             return false;
           }
        }
        return true;
    }

    private selectedTagsChanged(tags: number[]): void {
      const ideas = this.filterIdeas(this.state.activeFilter, this.state.search);
      this.setState({
        ideas: ideas.filter(
          (i) => i.tags.filter(
            (t) => tags.indexOf(t) >= 0
          ).length === tags.length
        ),
      });
    }

    private filterIdeas(activeFilter: string, search: string | undefined): Idea[] {
      let newUrl = getBaseUrl();
      let ideas = [];

      if (search) {
        const s = search.trim().toLowerCase();
        newUrl += `?q=${encodeURIComponent(s).replace(/%20/g, '+')}`;
        ideas = this.allIdeas.filter((idea) => {
          const terms = s.split(' ').filter((x) => x.length >= 2);
          return (
            this.containsAll(idea.title.toLowerCase(), terms) ||
            this.containsAll(idea.description.toLowerCase(), terms) ||
            (idea.response && this.containsAll(idea.response.text.toLowerCase(), terms))
          );
        });
      } else {
        newUrl += activeFilter ? `#${activeFilter}` : '';
        ideas = IdeaFilter.getFilter(activeFilter)(this.allIdeas);
      }

      if (history.replaceState) {
        window.history.replaceState({ path: newUrl }, '', newUrl);
      }

      return ideas;
    }

    private filterChanged(name: string) {
      this.setState({
        ideas: this.filterIdeas(name, this.state.search),
        showCount: defaultShowCount,
        activeFilter: name,
      });
    }

    private showMore(event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void {
      event.preventDefault();
      this.setState({
        showCount: this.state.showCount + defaultShowCount
      });
    }

    private resetSearch() {
      this.setState({
        search: '',
        searching: false,
        ideas: this.filterIdeas(this.state.activeFilter, ''),
      });
    }

    private searchIdea(input: string): void {
      this.setState({
        search: input,
        ideas: this.filterIdeas(this.state.activeFilter, input)
      });
    }

  public render() {
    const ideasToList = this.state.ideas.slice(0, this.state.showCount);

    const displayIdeas = (ideasToList.length > 0)
      ? (
        <div className="ui divided unstackable items fdr-idea-list">
            {
              ideasToList.map((idea) =>
              <ListIdeaItem
                key={idea.id}
                user={this.user}
                idea={idea}
                tags={this.allTags.filter((tag) => idea.tags.indexOf(tag.id) >= 0)}
              />)}
            {
              this.state.ideas.length > this.state.showCount &&
              <h5
                className="ui blue header show-more"
                onTouchEnd={(e) => this.showMore(e)}
                onClick={(e) => this.showMore(e)}
              >
                View {this.state.ideas.length - this.state.showCount} more ideas
              </h5>
            }
        </div>
      )
      : <p className="no-ideas-found">No ideas found for given filter.</p>;

    const welcomeMessage = this.tenant.welcomeMessage ||
    `## Welcome to our feedback forum!

We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas.`;

    return (
      <div>
        <Header />
        <div className="page ui container">

          <div className="ui grid stackable">
            <div className="six wide column">
              <MultiLineText className="welcome-message" text={welcomeMessage} style="full" />
              <IdeaInput placeholder={this.tenant.invitation || 'I suggest you...'} />
            </div>
            <div className="ten wide column">
              {
                this.allIdeas.length === 0
                ? <EmptyList />
                : <div>
                    <div className="ui grid">
                      {
                        !this.state.searching && <div className="ten wide mobile ten wide tablet twelve wide computer column filter-column">
                        <IdeaFilter activeFilter={this.state.activeFilter} filterChanged={(name) => this.filterChanged(name)} />
                      </div>
                      }
                      <div className={!this.state.searching ? `six wide mobile six wide tablet four wide computer column` : 'column'}>
                        <div className="ui search">
                          <div className="ui icon fluid input">
                            <input
                              onFocus={() => this.setState({ searching: true })}
                              onBlur={(x) => this.setState({ searching: !!this.state.search })}
                              onChange={(x) => this.searchIdea(x.currentTarget.value)}
                              value={this.state.search}
                              type="text"
                              placeholder="Search..."
                            />
                            {
                              this.state.searching
                              ? <i onClick={() => this.resetSearch()} className="cancel link icon" />
                              : <i className="search icon" />
                            }
                          </div>
                        </div>
                      </div>
                    </div>
                    <TagsFilter tags={this.allTags} selectionChanged={(selected) => this.selectedTagsChanged(selected)} />
                    {displayIdeas}
                  </div>
              }
            </div>
          </div>
        </div>
        <Footer />
      </div>
    );
  }
}
