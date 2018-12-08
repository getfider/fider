import React from "react";
import { PostStatus, Post } from "@fider/models";
import { actions, navigator, Failure, Fider } from "@fider/services";
import { Form, Modal, Button, List, ListItem, TextArea } from "@fider/components";

interface ModerationPanelProps {
  post: Post;
}

interface ModerationPanelState {
  showConfirmation: boolean;
  text: string;
  error?: Failure;
}

export class ModerationPanel extends React.Component<ModerationPanelProps, ModerationPanelState> {
  constructor(props: ModerationPanelProps) {
    super(props);
    this.state = {
      text: "",
      showConfirmation: false
    };
  }

  private delete = async () => {
    const response = await actions.deletePost(this.props.post.number, this.state.text);
    if (response.ok) {
      await this.closeModal();
      navigator.goHome();
    } else if (response.error) {
      this.setState({ error: this.state.error });
    }
  };

  private closeModal = async () => {
    this.setState({ showConfirmation: false });
  };

  private showModal = async () => {
    this.setState({ showConfirmation: true });
  };

  private setText = (text: string) => {
    this.setState({ text });
  };

  public render() {
    const status = PostStatus.Get(this.props.post.status);
    if (!Fider.session.isAuthenticated || !Fider.session.user.isAdministrator || status.closed) {
      return null;
    }

    const modal = (
      <Modal.Window isOpen={this.state.showConfirmation} center={false} size="large">
        <Modal.Content>
          <Form error={this.state.error}>
            <TextArea
              field="text"
              onChange={this.setText}
              value={this.state.text}
              placeholder="Why are you deleting this post? (optional)"
            >
              <span className="info">
                This operation <strong>cannot</strong> be undone.
              </span>
            </TextArea>
          </Form>
        </Modal.Content>

        <Modal.Footer>
          <Button color="danger" onClick={this.delete}>
            Delete
          </Button>
          <Button color="cancel" onClick={this.closeModal}>
            Cancel
          </Button>
        </Modal.Footer>
      </Modal.Window>
    );

    return (
      <>
        {modal}
        <span className="subtitle">Moderation</span>
        <List>
          <ListItem>
            <Button color="danger" size="tiny" fluid={true} onClick={this.showModal}>
              Delete
            </Button>
          </ListItem>
        </List>
      </>
    );
  }
}
