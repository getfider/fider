import "./Comments.scss";

import * as React from "react";
import { CurrentUser, Comment, Idea } from "@fider/models";
import { CommentList, CommentInput } from "../";
import { actions } from "@fider/services";

interface DiscussionPanelProps {
  user?: CurrentUser;
  idea: Idea;
  comments: Comment[];
}

interface DiscussionPanelState {
  isEditing: boolean;
}

export class DiscussionPanel extends React.Component<DiscussionPanelProps, DiscussionPanelState> {
  constructor(props: DiscussionPanelProps) {
    super(props);
    this.state = {
      isEditing: false
    };
  }

  public render() {
    return (
      <div className="comments-col">
        <div className="c-comment-list">
          <span className="subtitle">Discussion</span>
          <CommentList
            idea={this.props.idea}
            user={this.props.user}
            comments={this.props.comments}
            onStartEdit={() => this.setState({ isEditing: true })}
            onStopEdit={() => this.setState({ isEditing: false })}
          />
          {!this.state.isEditing && <CommentInput user={this.props.user} idea={this.props.idea} />}
        </div>
      </div>
    );
  }
}
