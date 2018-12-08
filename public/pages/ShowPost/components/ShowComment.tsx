import React from "react";
import { Comment, Post } from "@fider/models";
import {
  Gravatar,
  UserName,
  Moment,
  Form,
  TextArea,
  Button,
  MultiLineText,
  DropDown,
  DropDownItem,
  Modal
} from "@fider/components";
import { formatDate, Failure, actions, Fider } from "@fider/services";
import { FaEllipsisH } from "react-icons/fa";

interface ShowCommentProps {
  post: Post;
  comment: Comment;
}

interface ShowCommentState {
  comment: Comment;
  isEditing: boolean;
  newContent: string;
  error?: Failure;
  showDeleteConfirmation: boolean;
}

export class ShowComment extends React.Component<ShowCommentProps, ShowCommentState> {
  constructor(props: ShowCommentProps) {
    super(props);
    this.state = {
      comment: props.comment,
      isEditing: false,
      newContent: "",
      showDeleteConfirmation: false
    };
  }

  private canEditComment(comment: Comment): boolean {
    if (Fider.session.isAuthenticated) {
      return Fider.session.user.isCollaborator || comment.user.id === Fider.session.user.id;
    }
    return false;
  }

  private cancelEdit = async () => {
    this.setState({
      isEditing: false,
      newContent: "",
      error: undefined
    });
  };

  private saveEdit = async () => {
    const response = await actions.updateComment(this.props.post.number, this.state.comment.id, this.state.newContent);
    if (response.ok) {
      this.state.comment.content = this.state.newContent;
      this.state.comment.editedAt = new Date().toISOString();
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

  private renderEllipsis = () => {
    return <FaEllipsisH />;
  };

  private closeModal = async () => {
    this.setState({ showDeleteConfirmation: false });
  };

  private deleteComment = async () => {
    const response = await actions.deleteComment(this.props.post.number, this.props.comment.id);
    if (response.ok) {
      location.reload();
    }
  };

  private onActionSelected = (item: DropDownItem) => {
    if (item.value === "edit") {
      this.setState({ isEditing: true, newContent: this.state.comment.content, error: undefined });
    } else if (item.value === "delete") {
      this.setState({ showDeleteConfirmation: true });
    }
  };

  private modal() {
    return (
      <Modal.Window isOpen={this.state.showDeleteConfirmation} center={false} size="small">
        <Modal.Header>Delete Comment</Modal.Header>
        <Modal.Content>
          <p>
            This process is irreversible. <strong>Are you sure?</strong>
          </p>
        </Modal.Content>

        <Modal.Footer>
          <Button color="danger" onClick={this.deleteComment}>
            Delete
          </Button>
          <Button color="cancel" onClick={this.closeModal}>
            Cancel
          </Button>
        </Modal.Footer>
      </Modal.Window>
    );
  }

  public render() {
    const c = this.state.comment;

    const editedMetadata = !!c.editedAt && !!c.editedBy && (
      <div className="c-comment-metadata">
        <span title={`This comment has been edited by ${c.editedBy!.name} on ${formatDate(c.editedAt)}`}>· edited</span>
      </div>
    );

    return (
      <div className="c-comment">
        {this.modal()}
        <Gravatar user={c.user} />
        <div className="c-comment-content">
          <UserName user={c.user} />
          <div className="c-comment-metadata">
            · <Moment date={c.createdAt} />
          </div>
          {editedMetadata}
          {!this.state.isEditing && this.canEditComment(c) && (
            <DropDown
              className="l-more-actions"
              direction="left"
              inline={true}
              style="simple"
              highlightSelected={false}
              items={[
                { label: "Edit", value: "edit" },
                { label: "Delete", value: "delete", render: <span style={{ color: "red" }}>Delete</span> }
              ]}
              onChange={this.onActionSelected}
              renderControl={this.renderEllipsis}
            />
          )}
          <div className="c-comment-text">
            {this.state.isEditing ? (
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
                <Button color="cancel" size="tiny" onClick={this.cancelEdit}>
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
