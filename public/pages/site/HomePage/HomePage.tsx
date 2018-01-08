import * as React from 'react';
import { Idea, Tag, IdeaStatus, CurrentUser, Tenant } from '@fider/models';
import { ShowTag, ShowIdeaResponse, SupportCounter, Gravatar, MultiLineText, Moment } from '@fider/components';
import { IdeaInput, TagsFilter, IdeaFilter, IdeaFilterFunction } from './';
import { page } from '@fider/services';

import './HomePage.scss';

const defaultShowCount = 20;

export interface HomePageProps {
  user?: CurrentUser;
  tenant: Tenant;
  ideas: Idea[];
  tags: Tag[];
}

interface HomePageState {
  ideas: Idea[];
  search?: string;
  showCount: number;
  searching: boolean;
  activeFilter: string;
  tags: string[];
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

export class HomePage extends React.Component<HomePageProps, HomePageState> {
  private filter: HTMLDivElement;

  constructor(props: HomePageProps) {
    super(props);

    const search = page.getQueryString('q');
    const tags = page.getQueryStringArray('t');
    const activeFilter = window.location.hash.substring(1);
    this.state = {
      ideas: this.filterIdeas(activeFilter, search, tags),
      showCount: defaultShowCount,
      activeFilter,
      searching: !!search,
      search,
      tags
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

  private selectedTagsChanged(tags: string[]): void {
    const ideas = this.filterIdeas(this.state.activeFilter, this.state.search, tags);
    this.setState({
      ideas,
      tags,
    });
  }

  private filterIdeas(activeFilter: string, search: string | undefined, tags: string[]): Idea[] {
    let path = '';
    let ideas = [];

    if (search) {
      const s = search.trim().toLowerCase();
      path += `?q=${encodeURIComponent(s).replace(/%20/g, '+')}`;
      ideas = this.props.ideas.filter((idea) => {
        const terms = s.split(' ').filter((x) => x.length >= 2);
        return (
          this.containsAll(idea.title.toLowerCase(), terms) ||
          this.containsAll(idea.description.toLowerCase(), terms) ||
          (idea.response && this.containsAll(idea.response.text.toLowerCase(), terms))
        );
      });
    } else {
      path += activeFilter ? `#${activeFilter}` : '';
      ideas = IdeaFilter.getFilter(activeFilter)(this.props.ideas);
    }

    if (tags.length > 0) {
      const prefix = (!path) ? '?' : '&';
      path += `${prefix}t=${tags.join(',')}`;
    }

    if (history.replaceState) {
      const newUrl = page.getBaseUrl() + path;
      window.history.replaceState({ path: newUrl }, '', newUrl);
    }

    if (tags.length > 0) {
      const tagsToFilter = this.props.tags.filter((x) => tags.indexOf(x.slug) >= 0).map((x) => x.id);
      ideas = ideas.filter(
        (i) => i.tags.filter(
          (t) => tagsToFilter.indexOf(t) >= 0
        ).length === tagsToFilter.length
      );
    }

    return ideas;
  }

  private filterChanged(name: string) {
    this.setState({
      ideas: this.filterIdeas(name, this.state.search, this.state.tags),
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
      ideas: this.filterIdeas(this.state.activeFilter, '', this.state.tags),
    });
  }

  private searchIdea(input: string): void {
    this.setState({
      search: input,
      ideas: this.filterIdeas(this.state.activeFilter, input, this.state.tags)
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
                user={this.props.user}
                idea={idea}
                tags={this.props.tags.filter((tag) => idea.tags.indexOf(tag.id) >= 0)}
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

    const welcomeMessage = this.props.tenant.welcomeMessage ||
    `## Welcome to our feedback forum!

We'd love to hear what you're thinking about. What can we do better? This is the place for you to vote, discuss and share ideas.`;

    return (
      <div className="page ui container">
        <div className="ui grid stackable">
          <div className="six wide column">
            <MultiLineText className="welcome-message" text={welcomeMessage} style="full" />
            <IdeaInput
              user={this.props.user}
              placeholder={this.props.tenant.invitation || 'I suggest you...'}
            />
          </div>
          <div className="ten wide column">
            {
              this.props.ideas.length === 0
              ? <EmptyList />
              : <div>
                  <div className="ui grid">
                    {
                      !this.state.searching && <div className="ten wide mobile ten wide tablet twelve wide computer column filter-column">
                      <IdeaFilter
                        ideas={this.props.ideas}
                        activeFilter={this.state.activeFilter}
                        filterChanged={(name) => this.filterChanged(name)}
                      />
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
                  <TagsFilter
                    tags={this.props.tags}
                    selectionChanged={(selected) => this.selectedTagsChanged(selected)}
                    defaultSelection={this.state.tags}
                  />
                  {displayIdeas}
                </div>
            }
          </div>
        </div>
      </div>
    );
  }
}
