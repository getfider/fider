import * as React from "react";
import { Comment, CurrentUser, Idea } from "@fider/models";
import { Gravatar, UserName, Moment, Form, TextArea, Button, MultiLineText } from "@fider/components";
import { formatDate, Failure, actions } from "@fider/services";

interface ShowCommentProps {
  idea: Idea;
  comment: Comment;
}

interface ShowCommentState {
  comment: Comment;
  isEditting: boolean;
  newContent: string;
  error?: Failure;
}

export class ShowComment extends React.Component<ShowCommentProps, ShowCommentState> {
  constructor(props: ShowCommentProps) {
    super(props);
    this.state = {
      comment: props.comment,
      isEditting: false,
      newContent: ""
    };
  }

  private canEditComment(comment: Comment): boolean {
    if (Fider.session.isAuthenticated) {
      return Fider.session.user.isCollaborator || comment.user.id === Fider.session.user.id;
    }
    return false;
  }

  private startEdit = () => {
    this.setState({ isEditting: true, newContent: this.state.comment.content, error: undefined });
  };

  private cancelEdit = async () => {
    this.setState({
      isEditting: false,
      newContent: "",
      error: undefined
    });
  };

  private saveEdit = async () => {
    const response = await actions.updateComment(this.props.idea.number, this.state.comment.id, this.state.newContent);
    if (response.ok) {
      this.state.comment.content = this.state.newContent;
      this.state.comment.editedOn = new Date().toISOString();
      this.state.comment.editedBy = Fider.session.user;
      this.setState({
        comment: this.state.comment
      });
      this.cancelEdit();
    } else {
      this.setState({ error: response.error });
    }
  };

  private setNewContent = (newContent: string) => {
    this.setState({ newContent });
  };

  public render() {
    const c = this.state.comment;

    const edittedMetadata = !!c.editedOn &&
      !!c.editedBy && (
        <div className="c-comment-metadata">
          ·{" "}
          <span title={`This comment has been edited by ${c.editedBy!.name} on ${formatDate(c.editedOn)}`}>edited</span>
        </div>
      );

    return (
      <div className="c-comment">
        <Gravatar user={c.user} />
        <div className="c-comment-content">
          <UserName user={c.user} />
          <div className="c-comment-metadata">
            · <Moment date={c.createdOn} />
          </div>
          {edittedMetadata}
          {!this.state.isEditting &&
            this.canEditComment(c) && (
              <div className="c-comment-metadata">
                ·{" "}
                <span className="clickable" onClick={this.startEdit}>
                  edit
                </span>
              </div>
            )}
          <div className="c-comment-text">
            {this.state.isEditting ? (
              <Form error={this.state.error}>
                <TextArea
                  field="content"
                  minRows={1}
                  value={this.state.newContent}
                  placeholder={c.content}
                  onChange={this.setNewContent}
                />
                <Button size="tiny" onClick={this.saveEdit} color="positive">
                  Save
                </Button>
                <Button size="tiny" onClick={this.cancelEdit}>
                  Cancel
                </Button>
              </Form>
            ) : (
              <MultiLineText text={c.content} style="simple" />
            )}
          </div>
        </div>
      </div>
    );
  }
}
