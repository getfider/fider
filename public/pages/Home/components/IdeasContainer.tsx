import * as React from "react";

import { IdeaInput, ListIdeas, TagsFilter, IdeaFilter } from "../";

import { Idea, Tag, IdeaStatus, CurrentUser } from "@fider/models";
import { Loader, MultiLineText, Heading, Field, Input } from "@fider/components";
import { page, actions } from "@fider/services";

interface IdeasContainerProps {
  user?: CurrentUser;
  ideas: Idea[];
  tags: Tag[];
  countPerStatus: { [key: string]: number };
}

interface IdeasContainerState {
  loading: boolean;
  ideas?: Idea[];
  filter: string;
  tags: string[];
  query: string;
  limit?: number;
}

export class IdeasContainer extends React.Component<IdeasContainerProps, IdeasContainerState> {
  constructor(props: IdeasContainerProps) {
    super(props);

    this.state = {
      ideas: this.props.ideas,
      loading: false,
      filter: page.getQueryString("f"),
      query: page.getQueryString("q"),
      tags: page.getQueryStringArray("t"),
      limit: page.getQueryStringAsNumber("l")
    };
  }

  private changeFilterCriteria<K extends keyof IdeasContainerState>(
    obj: Pick<IdeasContainerState, K>,
    reset: boolean
  ): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase();
      page.replaceState(
        page.toQueryString({
          t: this.state.tags,
          q: query,
          f: this.state.filter,
          l: this.state.limit
        })
      );

      this.searchIdeas(query, this.state.filter, this.state.limit, this.state.tags, reset);
    });
  }

  private timer?: number;
  private async searchIdeas(query: string, filter: string, limit: number | undefined, tags: string[], reset: boolean) {
    window.clearTimeout(this.timer);
    this.setState({ ideas: reset ? undefined : this.state.ideas, loading: true });
    this.timer = window.setTimeout(() => {
      actions.searchIdeas({ query, filter, limit, tags }).then(response => {
        if (this.state.loading) {
          this.setState({ loading: false, ideas: response.data });
        }
      });
    }, 200);
  }

  private handleFilterChanged = (filter: string) => {
    this.changeFilterCriteria({ filter }, true);
  };

  private handleTagsFilterChanged = (tags: string[]) => {
    this.changeFilterCriteria({ tags }, true);
  };

  private handleSearchFilterChanged = (query: string) => {
    this.changeFilterCriteria({ query }, true);
  };

  private handleSearchClick = (query: string) => {
    this.changeFilterCriteria({ query }, true);
  };

  private clearSearch = () => {
    this.changeFilterCriteria({ query: "" }, true);
  };

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault();
    this.changeFilterCriteria({ limit: (this.state.limit || 30) + 10 }, false);
  };

  private canShowMore = (): boolean => {
    return this.state.ideas ? this.state.ideas.length >= (this.state.limit || 30) : false;
  };

  public render() {
    return (
      <>
        <div className="row">
          {!this.state.query && (
            <div className="l-filter-col col-sm-7 col-md-8 col-lg-9 mb-2">
              <Field>
                <IdeaFilter
                  user={this.props.user}
                  activeFilter={this.state.filter}
                  filterChanged={this.handleFilterChanged}
                  countPerStatus={this.props.countPerStatus}
                />
                <TagsFilter
                  tags={this.props.tags}
                  selectionChanged={this.handleTagsFilterChanged}
                  defaultSelection={this.state.tags}
                />
              </Field>
            </div>
          )}
          <div className={!this.state.query ? `l-search-col col-sm-5 col-md-4 col-lg-3 mb-2` : "col-sm-12 mb-2"}>
            <Input
              field="query"
              icon={this.state.query ? "cancel" : "search"}
              onIconClick={this.state.query ? this.clearSearch : undefined}
              placeholder="Search..."
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            />
          </div>
        </div>
        <ListIdeas
          ideas={this.state.ideas}
          tags={this.props.tags}
          user={this.props.user}
          emptyText={"No results matched your search, try something different."}
        />
        {this.state.loading && <Loader />}
        {this.canShowMore() && (
          <h5 className="c-idea-list-show-more" onTouchEnd={this.showMore} onClick={this.showMore}>
            View more ideas
          </h5>
        )}
      </>
    );
  }
}
