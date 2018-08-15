import "./ListPosts.scss";

import * as React from "react";
import { Post, Tag, PostStatus, CurrentUser } from "@fider/models";
import {
  ShowTag,
  ShowPostResponse,
  SupportCounter,
  Gravatar,
  MultiLineText,
  Moment,
  ListItem,
  List
} from "@fider/components";

interface ListPostsProps {
  posts?: Post[];
  tags: Tag[];
  emptyText: string;
}

const ListPostItem = (props: { post: Post; user?: CurrentUser; tags: Tag[] }) => {
  return (
    <ListItem>
      <SupportCounter post={props.post} />
      <div className="c-list-item-content">
        {props.post.totalComments > 0 && (
          <div className="info right">
            {props.post.totalComments} <i className="comments outline icon" />
          </div>
        )}
        <a className="c-list-item-title" href={`/posts/${props.post.number}/${props.post.slug}`}>
          {props.post.title}
        </a>
        <MultiLineText className="c-list-item-description" text={props.post.description} style="simple" />
        <ShowPostResponse status={props.post.status} response={props.post.response} />
        {props.tags.map(tag => (
          <ShowTag key={tag.id} size="tiny" tag={tag} />
        ))}
      </div>
    </ListItem>
  );
};

export class ListPosts extends React.Component<ListPostsProps, {}> {
  constructor(props: ListPostsProps) {
    super(props);
  }

  public render() {
    if (!this.props.posts) {
      return null;
    }

    if (this.props.posts.length === 0) {
      return <p className="center">{this.props.emptyText}</p>;
    }

    return (
      <List className="c-post-list" divided={true}>
        {this.props.posts.map(post => (
          <ListPostItem
            key={post.id}
            post={post}
            tags={this.props.tags.filter(tag => post.tags.indexOf(tag.slug) >= 0)}
          />
        ))}
      </List>
    );
  }
}
