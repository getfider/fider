import "./Comments.scss";

import * as React from "react";
import { CurrentUser, Comment, Post } from "@fider/models";
import { ShowComment, CommentInput } from "../";
import { actions } from "@fider/services";

interface DiscussionPanelProps {
  user?: CurrentUser;
  idea: Post;
  comments: Comment[];
}

export class DiscussionPanel extends React.Component<DiscussionPanelProps, {}> {
  constructor(props: DiscussionPanelProps) {
    super(props);
  }

  public render() {
    return (
      <div className="comments-col">
        <div className="c-comment-list">
          <span className="subtitle">Discussion</span>
          {this.props.comments.map(c => <ShowComment key={c.id} idea={this.props.idea} comment={c} />)}
          <CommentInput idea={this.props.idea} />
        </div>
      </div>
    );
  }
}
