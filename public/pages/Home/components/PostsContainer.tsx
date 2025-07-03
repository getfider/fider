import "./PostsContainer.scss"

import React from "react"

import { Post, Tag, CurrentUser } from "@fider/models"
import { Loader, Input } from "@fider/components"
import { actions, navigator, querystring } from "@fider/services"
import IconSearch from "@fider/assets/images/heroicons-search.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import { PostFilter } from "./PostFilter"
import { ListPosts } from "./ListPosts"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { PostsSort } from "./PostsSort"

interface PostsContainerProps {
  user?: CurrentUser
  posts: Post[]
  tags: Tag[]
  countPerStatus: { [key: string]: number }
}

interface PostsContainerState {
  loading: boolean
  posts?: Post[] // All posts
  view: string
  filterState: FilterState // Filter state
  query: string // Seach query
  limit?: number // Limit
}

export interface FilterState {
  tags: string[]
  statuses: string[]
  myVotes: boolean
}

export class PostsContainer extends React.Component<PostsContainerProps, PostsContainerState> {
  constructor(props: PostsContainerProps) {
    super(props)

    const view = querystring.get("view")

    this.state = {
      posts: this.props.posts,
      loading: false,
      view,
      query: querystring.get("query"),
      filterState: { tags: querystring.getArray("tags"), statuses: querystring.getArray("statuses"), myVotes: querystring.get("myvotes") === "true" },
      limit: querystring.getNumber("limit"),
    }
  }

  private changeFilterCriteria<K extends keyof PostsContainerState>(obj: Pick<PostsContainerState, K>, reset: boolean): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase()
      navigator.replaceState(
        querystring.stringify({
          statuses: this.state.filterState.statuses,
          tags: this.state.filterState.tags,
          myvotes: this.state.filterState.myVotes ? "true" : undefined,
          query,
          view: this.state.view,
          limit: this.state.limit,
        })
      )

      this.searchPosts(
        query,
        this.state.view || "trending",
        this.state.limit,
        this.state.filterState.tags,
        this.state.filterState.statuses,
        this.state.filterState.myVotes,
        reset
      )
    })
  }

  private timer?: number
  private async searchPosts(query: string, view: string, limit: number | undefined, tags: string[], statuses: string[], myVotes: boolean, reset: boolean) {
    window.clearTimeout(this.timer)
    this.setState({ posts: reset ? undefined : this.state.posts, loading: true })
    this.timer = window.setTimeout(() => {
      actions.searchPosts({ query, view: view, limit, tags, statuses, myVotes }).then((response) => {
        if (response.ok && this.state.loading) {
          this.setState({ loading: false, posts: response.data })
        }
      })
    }, 500)
  }

  private handleFilterChanged = (filterState: FilterState) => {
    this.changeFilterCriteria({ filterState }, true)
  }

  private handleSearchFilterChanged = (query: string) => {
    this.changeFilterCriteria({ query }, true)
  }

  private handleSortChanged = (view: string) => {
    this.changeFilterCriteria({ view }, true)
  }

  private clearSearch = () => {
    this.changeFilterCriteria({ query: "" }, true)
  }

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault()
    this.changeFilterCriteria({ limit: (this.state.limit || 30) + 10 }, false)
  }

  private getShowMoreLink = (): string | undefined => {
    if (this.state.posts && this.state.posts.length >= (this.state.limit || 30)) {
      return querystring.set("limit", (this.state.limit || 30) + 10)
    }
  }

  public render() {
    const showMoreLink = this.getShowMoreLink()

    return (
      <div className="c-posts-container">
        <div className="c-posts-container__header mb-5">
          {!this.state.query && (
            <div className="c-posts-container__filter-col">
              <PostFilter
                tags={this.props.tags}
                activeFilter={this.state.filterState}
                filtersChanged={this.handleFilterChanged}
                countPerStatus={this.props.countPerStatus}
              />
              <PostsSort onChange={this.handleSortChanged} value={this.state.view} />
            </div>
          )}
          <div className="c-posts-container__search-col">
            <Input
              field="query"
              icon={this.state.query ? IconX : IconSearch}
              onIconClick={this.state.query ? this.clearSearch : undefined}
              placeholder={i18n._({ id: "home.postscontainer.query.placeholder", message: "Search" })}
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            />
          </div>
        </div>
        <ListPosts
          posts={this.state.posts}
          tags={this.props.tags}
          emptyText={i18n._({ id: "home.postscontainer.label.noresults", message: "No results matched your search, try something different." })}
        />
        {this.state.loading && <Loader />}
        {showMoreLink && (
          <div className="my-4 ml-4">
            <a href={showMoreLink} className="text-primary-base text-medium hover:underline" onClick={this.showMore}>
              <Trans id="home.postscontainer.label.viewmore">View more posts</Trans>
            </a>
          </div>
        )}
      </div>
    )
  }
}
