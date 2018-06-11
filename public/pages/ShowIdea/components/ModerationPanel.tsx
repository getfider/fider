import * as React from "react";
import { IdeaStatus, CurrentUser, Idea } from "@fider/models";
import { actions, Failure } from "@fider/services";
import { Form, DisplayError, Modal, Button, List, ListItem, TextArea } from "@fider/components";

interface ModerationPanelProps {
  idea: Idea;
}

interface ModerationPanelState {
  showConfirmation: boolean;
  text: string;
  error?: Failure;
}

export class ModerationPanel extends React.Component<ModerationPanelProps, ModerationPanelState> {
  private form!: Form;

  constructor(props: ModerationPanelProps) {
    super(props);
    this.state = {
      text: "",
      showConfirmation: false
    };
  }

  private delete = async () => {
    const response = await actions.deleteIdea(this.props.idea.number, this.state.text);
    if (response.ok) {
      await this.closeModal();
      actions.goHome();
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
    const status = IdeaStatus.Get(this.props.idea.status);
    if (!page.user || !page.user.isAdministrator || status.closed) {
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
              placeholder="Why are you deleting this idea? (optional)"
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
          <Button onClick={this.closeModal}>Cancel</Button>
        </Modal.Footer>
      </Modal.Window>
    );

    return (
      <div>
        {modal}
        <span className="subtitle">Moderation</span>
        <List>
          <ListItem>
            <Button color="danger" size="tiny" fluid={true} onClick={this.showModal}>
              <i className="delete icon" /> Delete
            </Button>
          </ListItem>
        </List>
      </div>
    );
  }
}
