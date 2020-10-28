import React from "react";
import { Post, Tag, CurrentUser } from "@fider/models";
import { Heading, Loader } from "@fider/components";
import { ListPosts } from "./ListPosts";
import { actions } from "@fider/services";
import { withTranslation, WithTranslation } from "react-i18next";
import { FaRegLightbulb } from "react-icons/fa";

interface SimilarPostsProps extends WithTranslation {
  title: string;
  tags: Tag[];
  user?: CurrentUser;
}

interface SimilarPostsState {
  title: string;
  posts: Post[];
  loading: boolean;
}

class InternalSimilarPosts extends React.Component<SimilarPostsProps, SimilarPostsState> {
  constructor(props: SimilarPostsProps) {
    super(props);
    this.state = {
      title: props.title,
      loading: !!props.title,
      posts: []
    };
  }

  public static getDerivedStateFromProps(nextProps: SimilarPostsProps, prevState: SimilarPostsState) {
    if (nextProps.title !== prevState.title) {
      return {
        loading: true,
        title: nextProps.title
      };
    }
    return null;
  }
  public componentDidMount() {
    this.loadSimilarPosts();
  }

  private timer?: number;
  public componentDidUpdate() {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(this.loadSimilarPosts, 500);
  }

  private loadSimilarPosts = () => {
    if (this.state.loading) {
      actions.searchPosts({ query: this.state.title }).then(x => {
        if (x.ok) {
          this.setState({ loading: false, posts: x.data });
        }
      });
    }
  };

  public render() {
    const { t } = this.props;
    return (
      <>
        <Heading
          title={t("home.similarPosts.title")}
          subtitle={t("home.similarPosts.subtitle")}
          icon={FaRegLightbulb}
          size="small"
          dividing={true}
        />
        {this.state.loading ? (
          <Loader />
        ) : (
          <ListPosts
            posts={this.state.posts}
            tags={this.props.tags}
            emptyText={t("home.similarPosts.emptyText", { title: this.props.title })}
          />
        )}
      </>
    );
  }
}

export const SimilarPosts = withTranslation()(InternalSimilarPosts);
