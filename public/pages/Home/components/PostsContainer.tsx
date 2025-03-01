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
import { PostsSort } from "./PostsSort"

export interface FilterState {
  tags: string[]
  statuses: string[]
  myVotes: boolean
}

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
  query: string // Search query
  offset: number // Offset for pagination
  limit: number
  hasMore: boolean
}

const untaggedTag: Tag = {
  id: -1,
  slug: "untagged",
  name: "untagged",
  color: "#cccccc",
  isPublic: false,
}

export class PostsContainer extends React.Component<PostsContainerProps, PostsContainerState> {
  private timer?: number
  private loadMoreRef = React.createRef<HTMLDivElement>()
  private observer?: IntersectionObserver

  constructor(props: PostsContainerProps) {
    super(props)
    const view = querystring.get("view") || "trending"
    this.state = {
      posts: this.props.posts,
      loading: false,
      view,
      query: querystring.get("query") || "",
      filterState: {
        tags: querystring.getArray("tags"),
        statuses: querystring.getArray("statuses"),
        myVotes: querystring.get("myvotes") === "true",
      },
      limit: querystring.getNumber("limit") || 15,
      offset: 0,
      hasMore: true,
    }
  }

  componentDidMount() {
    this.observer = new IntersectionObserver(this.handleObserver, {
      root: null,
      rootMargin: "0px",
      threshold: 1.0,
    })
    if (this.loadMoreRef.current) {
      this.observer.observe(this.loadMoreRef.current)
    }
  }

  componentDidUpdate(prevProps: PostsContainerProps, prevState: PostsContainerState) {
    if (this.loadMoreRef.current && this.observer) {
      this.observer.observe(this.loadMoreRef.current)
    }
  }

  componentWillUnmount() {
    if (this.observer) {
      this.observer.disconnect()
    }
    if (this.timer) clearTimeout(this.timer)
  }

  handleObserver = (entries: IntersectionObserverEntry[]) => {
    const entry = entries[0]
    if (entry.isIntersecting && !this.state.loading && this.state.hasMore) {
      this.loadMore()
    }
  }

  loadMore = () => {
    const newOffset = (this.state.offset || 0) + this.state.limit
    this.setState({ loading: true })
    actions
      .searchPosts({
        query: this.state.query,
        view: this.state.view,
        limit: this.state.limit,
        offset: newOffset,
        tags: this.state.filterState.tags,
        statuses: this.state.filterState.statuses,
        myVotes: this.state.filterState.myVotes,
      })
      .then((response) => {
        if (response.ok) {
          const newPosts: Post[] = response.data || []
          const hasMore = newPosts.length === this.state.limit
          this.setState((prevState) => ({
            posts: prevState.posts ? [...prevState.posts, ...newPosts] : newPosts,
            offset: newOffset,
            loading: false,
            hasMore,
          }))
        } else {
          this.setState({ loading: false })
        }
      })
  }

  private changeFilterCriteria<K extends keyof PostsContainerState>(
    obj: Pick<PostsContainerState, K>,
    reset: boolean
  ): void {
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
      this.setState({ offset: 0, hasMore: true }, () => {
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
    })
  }

  private async searchPosts(
    query: string,
    view: string,
    limit: number,
    tags: string[],
    statuses: string[],
    myVotes: boolean,
    reset: boolean
  ) {
    window.clearTimeout(this.timer)
    this.setState({ posts: reset ? undefined : this.state.posts, loading: true, offset: reset ? 0 : this.state.offset })
    this.timer = window.setTimeout(() => {
      actions
        .searchPosts({ query, view, limit, tags, statuses, myVotes, offset: this.state.offset })
        .then((response) => {
          if (response.ok && this.state.loading) {
            const posts: Post[] = response.data || []
            const hasMore = posts.length === limit
            this.setState({ loading: false, posts: reset ? posts : this.state.posts, hasMore })
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

  public render() {
    return (
      <div className="c-posts-container">
        <div className="c-posts-container__header mb-5">
          {!this.state.query && (
            <div className="c-posts-container__filter-col">
              <PostFilter
                tags={[untaggedTag, ...this.props.tags]}
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
              placeholder={i18n._("home.postscontainer.query.placeholder", { message: "Search" })}
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            />
          </div>
        </div>
        <ListPosts
          posts={this.state.posts}
          tags={this.props.tags}
          emptyText={i18n._("home.postscontainer.label.noresults", { message: "No results matched your search, try something different." })}
        />
        {this.state.loading && <Loader />}
        {this.state.hasMore && <div ref={this.loadMoreRef} style={{ height: "1px" }}></div>}
      </div>
    )
  }
}
