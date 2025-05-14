import React from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { ListPosts } from "./ListPosts"
import { actions } from "@fider/services"

import { i18n } from "@lingui/core"

import "./SimilarPosts.scss"

interface SimilarPostsProps {
  title: string
  tags: Tag[]
  user?: CurrentUser
}

interface SimilarPostsState {
  title: string
  posts: Post[]
  loading: boolean
  visible: boolean
}

export class SimilarPosts extends React.Component<SimilarPostsProps, SimilarPostsState> {
  constructor(props: SimilarPostsProps) {
    super(props)
    this.state = {
      title: props.title,
      loading: !!props.title,
      posts: [],
      visible: false,
    }
  }

  public static getDerivedStateFromProps(nextProps: SimilarPostsProps, prevState: SimilarPostsState) {
    if (nextProps.title !== prevState.title) {
      return {
        loading: true,
        title: nextProps.title,
        visible: prevState.posts.length > 0 ? prevState.visible : false,
      }
    }
    return null
  }
  public componentDidMount() {
    console.log("componentDidMount", this.state.visible)
    this.loadSimilarPosts()
  }

  private timer?: number
  public componentDidUpdate() {
    window.clearTimeout(this.timer)
    console.log("componentDidUpdate", this.state.visible)
    this.timer = window.setTimeout(this.loadSimilarPosts, 1000)
  }

  private loadSimilarPosts = () => {
    if (this.state.loading) {
      actions.searchPosts({ query: this.state.title, limit: 5 }).then((x) => {
        if (x.ok) {
          this.setState({ loading: false, posts: x.data }, () => {
            if (this.state.posts.length > 0) {
              setTimeout(() => {
                this.setState({ visible: true })
              }, 50)
            }
          })
        }
      })
    }
  }

  public render() {
    const title = i18n._("home.similar.title", { message: "Similar posts" })
    const subtitle = i18n._("home.similar.subtitle", { message: "Consider voting on existing posts instead." })

    const animationClass = this.state.visible ? "similar-posts-visible" : "similar-posts-hidden"

    return (
      <>
        {this.state.posts.length > 0 && (
          <div className={`mb-4 similar-posts-container ${animationClass}`}>
            <div className="mb-4 text-gray-700">
              {title} - {subtitle}
            </div>
            <ListPosts posts={this.state.posts} tags={this.props.tags} emptyText={`No similar posts matched '${this.props.title}'.`} minimalView={true} />
          </div>
        )}
      </>
    )
  }
}
