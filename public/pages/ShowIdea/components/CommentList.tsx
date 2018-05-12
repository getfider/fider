import * as React from "react";
import { Idea, Comment, CurrentUser } from "@fider/models";
import { Failure, actions, formatDate } from "@fider/services";
import { DisplayError, Textarea, Button, UserName, Gravatar, Moment, MultiLineText } from "@fider/components/common";

interface CommentListProps {
  idea: Idea;
  comments: Comment[];
  user?: CurrentUser;
  onStartEdit?: () => void;
  onStopEdit?: () => void;
}

interface CommentListState {
  editingComment?: Comment;
  editCommentNewContent: string;
  error?: Failure;
}

export class CommentList extends React.Component<CommentListProps, CommentListState> {
  constructor(props: CommentListProps) {
    super(props);
    this.state = {
      editCommentNewContent: ""
    };
  }

  private async startEdit(comment: Comment): Promise<void> {
    this.setState({
      editingComment: comment,
      editCommentNewContent: comment.content,
      error: undefined
    });

    if (this.props.onStartEdit) {
      this.props.onStartEdit();
    }
  }

  private async cancelEdit(): Promise<void> {
    this.setState({
      editingComment: undefined,
      editCommentNewContent: "",
      error: undefined
    });

    if (this.props.onStopEdit) {
      this.props.onStopEdit();
    }
  }

  private async confirmEdit(): Promise<void> {
    if (this.state.editingComment) {
      const response = await actions.updateComment(
        this.props.idea.number,
        this.state.editingComment.id,
        this.state.editCommentNewContent
      );
      if (response.ok) {
        this.state.editingComment.content = this.state.editCommentNewContent;
        this.state.editingComment.editedOn = new Date().toISOString();
        this.state.editingComment.editedBy = this.props.user;
        this.cancelEdit();
      } else {
        this.setState({ error: response.error });
      }
    }
  }

  private canEditComment(comment: Comment): boolean {
    if (this.props.user) {
      return this.props.user.isCollaborator || comment.user.id === this.props.user.id;
    }
    return false;
  }

  public render() {
    return this.props.comments.map(c => {
      return (
        <div key={c.id} className="comment">
          <Gravatar user={c.user} />
          <div className="content">
            <UserName user={c.user} />
            <div className="metadata">
              · <Moment date={c.createdOn} />
            </div>
            {!!c.editedOn &&
              !!c.editedBy && (
                <div className="metadata">
                  ·{" "}
                  <span title={`This comment has been edited by ${c.editedBy!.name} on ${formatDate(c.editedOn)}`}>
                    edited
                  </span>
                </div>
              )}
            {this.canEditComment(c) && (
              <div className="metadata">
                ·{" "}
                <span className="clickable" onClick={() => this.startEdit(c)}>
                  edit
                </span>
              </div>
            )}
            <div className="text">
              {c === this.state.editingComment ? (
                <div className="ui form">
                  <DisplayError error={this.state.error} />
                  <div className="field">
                    <Textarea
                      rows={1}
                      defaultValue={c.content}
                      placeholder={c.content}
                      onChange={e =>
                        this.setState({
                          editCommentNewContent: e.currentTarget.value
                        })
                      }
                    />
                  </div>
                  <Button size="tiny" onClick={() => this.confirmEdit()} color="positive">
                    Save
                  </Button>
                  <Button size="tiny" onClick={() => this.cancelEdit()}>
                    Cancel
                  </Button>
                </div>
              ) : (
                <MultiLineText text={c.content} style="simple" />
              )}
            </div>
          </div>
        </div>
      );
    });
  }
}
