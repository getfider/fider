import * as React from "react";
import { Post, PostStatus } from "@fider/models";
import { actions } from "@fider/services";
import { Dropdown, DropdownProps, DropdownItemProps, DropdownOnSearchChangeData } from "@fider/components";

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

  private onSearchChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownOnSearchChangeData) => {
    this.search(data.searchQuery);
  };

  private onChange = (e: React.SyntheticEvent<HTMLElement>, data: DropdownProps) => {
    this.props.onChanged(parseInt(data.value as string, 10));
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

  private returnAll = (options: DropdownItemProps[], value: string) => options;

  public render() {
    const options = this.state.posts.map(i => {
      const status = PostStatus.Get(i.status);
      return {
        key: i.number,
        text: i.title,
        value: i.number,
        content: (
          <>
            <span className="votes">
              <i className="caret up icon" />
              {i.totalVotes}
            </span>
            <span className={`status-label status-${status.value}`}>{status.title}</span>
            {i.title}
          </>
        )
      };
    });

    return (
      <Dropdown
        className="c-post-search"
        fluid={true}
        selectOnBlur={false}
        selection={true}
        search={this.returnAll}
        options={options}
        placeholder="Search original post"
        onChange={this.onChange}
        onSearchChange={this.onSearchChange}
      />
    );
  }
}
