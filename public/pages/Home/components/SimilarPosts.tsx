import React from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { PageTitle, Loader } from "@fider/components"
import { ListPosts } from "./ListPosts"
import { actions } from "@fider/services"

import { t } from "@lingui/macro"

interface SimilarPostsProps {
  title: string
  tags: Tag[]
  user?: CurrentUser
}

interface SimilarPostsState {
  title: string
  posts: Post[]
  loading: boolean
}

export class SimilarPosts extends React.Component<SimilarPostsProps, SimilarPostsState> {
  constructor(props: SimilarPostsProps) {
    super(props)
    this.state = {
      title: props.title,
      loading: !!props.title,
      posts: [],
    }
  }

  public static getDerivedStateFromProps(nextProps: SimilarPostsProps, prevState: SimilarPostsState) {
    if (nextProps.title !== prevState.title) {
      return {
        loading: true,
        title: nextProps.title,
      }
    }
    return null
  }
  public componentDidMount() {
    this.loadSimilarPosts()
  }

  private timer?: number
  public componentDidUpdate() {
    window.clearTimeout(this.timer)
    this.timer = window.setTimeout(this.loadSimilarPosts, 500)
  }

  private loadSimilarPosts = () => {
    if (this.state.loading) {
      actions.searchPosts({ query: this.state.title }).then((x) => {
        if (x.ok) {
          this.setState({ loading: false, posts: x.data })
        }
      })
    }
  }

  public render() {
    const title = t({ id: "home.similar.title", message: "Similar posts" })
    const subtitle = t({ id: "home.similar.subtitle", message: "Consider voting on existing posts instead." })

    return (
      <>
        <PageTitle title={title} subtitle={subtitle} />
        {this.state.loading ? (
          <Loader />
        ) : (
          <ListPosts posts={this.state.posts} tags={this.props.tags} emptyText={`No similar posts matched '${this.props.title}'.`} />
        )}
      </>
    )
  }
}
