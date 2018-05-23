import * as React from "react";

import { IdeaInput, ListIdeas, TagsFilter, IdeaFilter } from "../";

import { Idea, Tag, IdeaStatus, CurrentUser } from "@fider/models";
import { Loader, MultiLineText, Heading, Field, Input } from "@fider/components";
import { page, actions } from "@fider/services";

interface IdeasContainerProps {
  user?: CurrentUser;
  ideas: Idea[];
  tags: Tag[];
  newIdeaTitle: string;
  countPerStatus: { [key: string]: number };
}

interface IdeasContainerState {
  loading: boolean;
  ideas: Idea[];
  filter: string;
  tags: string[];
  query: string;
}

export class IdeasContainer extends React.Component<IdeasContainerProps, IdeasContainerState> {
  private timer?: number;

  constructor(props: IdeasContainerProps) {
    super(props);
    const query = page.getQueryString("q");
    const tags = page.getQueryStringArray("t");
    const filter = page.getQueryString("f");

    this.state = {
      ideas: [],
      loading: false,
      filter,
      query,
      tags
    };
  }

  public componentWillReceiveProps(nextProps: IdeasContainerProps) {
    if (nextProps.newIdeaTitle) {
      this.searchIdeas(nextProps.newIdeaTitle, "", [], 200);
    } else if (this.state.query || this.state.filter || this.state.tags.length > 0) {
      this.searchIdeas(this.state.query, this.state.filter, this.state.tags);
    } else {
      this.setState({ loading: false, ideas: nextProps.ideas });
    }
  }

  private changeFilterCriteria<K extends keyof IdeasContainerState>(
    obj: Pick<IdeasContainerState, K>,
    delay: number = 0
  ): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase();
      page.replaceState(
        page.toQueryString({
          t: this.state.tags,
          q: query,
          f: this.state.filter
        })
      );

      this.searchIdeas(query, this.state.filter, this.state.tags, delay);
    });
  }

  private async searchIdeas(query: string, filter: string, tags: string[], delay: number = 0) {
    window.clearTimeout(this.timer);
    this.setState({ loading: true });
    this.timer = window.setTimeout(() => {
      actions.searchIdeas(query, filter, tags).then(response => {
        if (this.state.loading) {
          this.setState({ loading: false, ideas: response.data });
        }
      });
    }, delay);
  }

  public render() {
    if (this.props.newIdeaTitle) {
      return (
        <>
          <Heading
            title="Similar ideas"
            subtitle="Consider voting on existing ideas before posting a new one."
            icon="lightbulb"
            size="small"
            dividing={true}
          />
          {this.state.loading ? (
            <Loader />
          ) : (
            <ListIdeas
              ideas={this.state.ideas}
              tags={this.props.tags}
              user={this.props.user}
              emptyText={`No similar ideas matched '${this.props.newIdeaTitle}'.`}
            />
          )}
        </>
      );
    }

    return (
      <>
        <div className="row">
          {!this.state.query && (
            <div className="col-sm-7 col-md-8 col-lg-9">
              <Field>
                <IdeaFilter
                  activeFilter={this.state.filter}
                  filterChanged={filter => this.changeFilterCriteria({ filter })}
                  countPerStatus={this.props.countPerStatus}
                />
                <TagsFilter
                  tags={this.props.tags}
                  selectionChanged={tags => this.changeFilterCriteria({ tags })}
                  defaultSelection={this.state.tags}
                />
              </Field>
            </div>
          )}
          <div className={!this.state.query ? `col-sm-5 col-md-4 col-lg-3` : "col-sm-12"}>
            <Input
              field="query"
              icon={this.state.query ? "cancel" : "search"}
              onIconClick={this.state.query ? () => this.changeFilterCriteria({ query: "" }) : undefined}
              placeholder="Search..."
              value={this.state.query}
              onChange={query => this.changeFilterCriteria({ query }, 200)}
            />
          </div>
        </div>
        {this.state.loading ? (
          <Loader />
        ) : (
          <ListIdeas
            ideas={this.state.ideas}
            tags={this.props.tags}
            user={this.props.user}
            emptyText={"No results matched your search, try something different."}
          />
        )}
      </>
    );
  }
}
