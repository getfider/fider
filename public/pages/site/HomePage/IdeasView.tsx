import * as React from 'react';
import { Idea, Tag, IdeaStatus, CurrentUser } from '@fider/models';
import { Loader, MultiLineText } from '@fider/components';
import { IdeaInput, ListIdeas, TagsFilter, IdeaFilter } from './';
import { page, actions } from '@fider/services';

interface IdeasViewProps {
  user?: CurrentUser;
  ideas: Idea[];
  tags: Tag[];
  newIdeaTitle: string;
  countPerStatus: { [key: string]: number };
}

interface IdeasViewState {
  loading: boolean;
  ideas: Idea[];
  filter: string;
  tags: string[];
  query: string;
}

export class IdeasView extends React.Component<IdeasViewProps, IdeasViewState> {
    private timer: any;

    constructor(props: IdeasViewProps) {
      super(props);
      const query = page.getQueryString('q');
      const tags = page.getQueryStringArray('t');
      const filter = page.getQueryString('f');

      this.state = {
        ideas: [],
        loading: false,
        filter,
        query,
        tags
      };
    }

    public componentWillReceiveProps(nextProps: IdeasViewProps) {
      if (nextProps.newIdeaTitle) {
        this.searchIdeas(nextProps.newIdeaTitle, '', [], 200);
      } else if (this.state.query || this.state.filter || this.state.tags.length > 0) {
        this.searchIdeas(this.state.query, this.state.filter, this.state.tags);
      } else {
        this.setState({ ideas: nextProps.ideas });
      }
    }

    private changeFilterCriteria<K extends keyof IdeasViewState>(obj: Pick<IdeasViewState, K>, delay: number = 0): void {
      this.setState(obj, () => {
        const query = this.state.query.trim().toLowerCase();
        page.replaceState(page.toQueryString({
          t: this.state.tags,
          q: query,
          f: this.state.filter,
        }));

        this.searchIdeas(query, this.state.filter, this.state.tags, delay);
      });
    }

    private async searchIdeas(query: string, filter: string, tags: string[], delay: number = 0) {
      clearTimeout(this.timer);
      this.setState({ loading: true });
      this.timer = setTimeout(() => {
        actions.searchIdeas(query, filter, tags).then((response) => {
          this.setState({ loading: false, ideas: response.data });
        });
      }, delay);
    }

    public render() {
      if (this.props.newIdeaTitle) {
        return (
          <>
            <h3 className="ui dividing header">
              <i className="lightbulb icon blue" />
              <div className="content">
                Similar ideas
                <div className="sub header">Consider voting on existing ideas before posting a new one.</div>
              </div>
            </h3>
            {
              this.state.loading
              ? <Loader />
              : <ListIdeas
                ideas={this.state.ideas}
                tags={this.props.tags}
                user={this.props.user}
                emptyText={`No similar ideas matched '${this.props.newIdeaTitle}'.`}
              />
            }
          </>
        );
      }

      return (
        <>
          <div className="ui grid">
            {
              !this.state.query && <div className="ten wide mobile ten wide tablet twelve wide computer column filter-column">
                <IdeaFilter
                  activeFilter={this.state.filter}
                  filterChanged={(filter) => this.changeFilterCriteria({ filter })}
                  countPerStatus={this.props.countPerStatus}
                />
              </div>
            }
            <div className={!this.state.query ? `six wide mobile six wide tablet four wide computer column` : 'column'}>
              <div className="ui search">
                <div className="ui icon fluid input">
                  <input
                    onChange={(x) => this.changeFilterCriteria({ query: x.currentTarget.value }, 200)}
                    value={this.state.query}
                    type="text"
                    placeholder="Search..."
                  />
                  {
                    this.state.query
                    ? <i onClick={() => this.changeFilterCriteria({ query: '' })} className="cancel link icon" />
                    : <i className="search icon" />
                  }
                </div>
              </div>
            </div>
          </div>
          <TagsFilter
            tags={this.props.tags}
            selectionChanged={(tags) => this.changeFilterCriteria({ tags })}
            defaultSelection={this.state.tags}
          />
          {
            this.state.loading
            ? <Loader />
            : <ListIdeas
              ideas={this.state.ideas}
              tags={this.props.tags}
              user={this.props.user}
              emptyText={'No results matched your search, try something different.'}
            />
          }
        </>
      );
    }
}
