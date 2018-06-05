import "./Comments.scss";

import * as React from "react";
import { CurrentUser, Comment, Idea } from "@fider/models";
import { ShowComment, CommentInput } from "../";
import { actions } from "@fider/services";

interface DiscussionPanelProps {
  user?: CurrentUser;
  idea: Idea;
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
          {this.props.comments.map(c => (
            <ShowComment key={c.id} idea={this.props.idea} user={this.props.user} comment={c} />
          ))}
          <CommentInput user={this.props.user} idea={this.props.idea} />
        </div>
      </div>
    );
  }
}
