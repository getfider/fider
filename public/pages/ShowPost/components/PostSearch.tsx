import * as React from "react";
import { Post, PostStatus } from "@fider/models";
import { actions } from "@fider/services";
import { FiderDropDown, FiderDropDownItem } from "@fider/components";

interface PostSearchProps {
  exclude?: number[];
  onChanged(postNumber: number): void;
}

interface PostSearchState {
  posts: Post[];
}

export class PostSearch extends React.Component<PostSearchProps, PostSearchState> {
  private timer?: number;

  constructor(props: PostSearchProps) {
    super(props);
    this.state = {
      posts: []
    };
    this.search("");
  }

  // private onSearchChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownOnSearchChangeData) => {
  //   this.search(data.searchQuery);
  // };

  private onChange = (item: FiderDropDownItem) => {
    this.props.onChanged(item.value as number);
  };

  private search = (searchQuery: string) => {
    window.clearTimeout(this.timer);
    this.timer = window.setTimeout(() => {
      actions.searchPosts({ query: searchQuery }).then(res => {
        const posts =
          this.props.exclude && this.props.exclude.length > 0
            ? res.data.filter(i => this.props.exclude!.indexOf(i.number) === -1)
            : res.data;
        this.setState({ posts });
      });
    }, 200);
  };

  private returnAll = (options: FiderDropDownItem[], value: string) => options;

  public renderItem = (item: FiderDropDownItem) => {
    const post = this.state.posts.filter(p => p.number === item.value)[0];
    const status = PostStatus.Get(post.status);
    return (
      <>
        <span className="votes">
          <i className="caret up icon" />
          {post.votesCount}
        </span>
        <span className={`status-label status-${status.value}`}>{status.title}</span>
        {post.title}
      </>
    );
  };

  public render() {
    const items = this.state.posts.map(i => {
      return {
        label: i.title,
        value: i.number
      };
    });

    return (
      <FiderDropDown
        className="c-post-search"
        // search={this.returnAll}
        items={items}
        placeholder="Search original post"
        onChange={this.onChange}
        renderItem={this.renderItem}
        // onSearchChange={this.onSearchChange}
      />
    );
  }
}
