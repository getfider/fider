import "./VotersModal.scss";

import React from "react";
import { Post, User } from "@fider/models";
import { Modal, Button, Loader, List, ListItem } from "@fider/components";
import { actions } from "@fider/services";

interface VotersModalProps {
  isOpen: boolean;
  post: Post;
  onClose?: () => void;
}

interface VotersModalState {
  searchText: string;
  users: User[];
  isLoading: boolean;
}

export class VotersModal extends React.Component<VotersModalProps, VotersModalState> {
  constructor(props: VotersModalProps) {
    super(props);
    this.state = {
      searchText: "",
      users: [],
      isLoading: true
    };
  }

  public componentDidMount() {
    actions.listVoters(this.props.post.number).then(response => {
      if (response.ok) {
        this.setState({ users: response.data, isLoading: false });
      }
    });
  }

  private closeModal = async () => {
    if (this.props.onClose) {
      this.props.onClose();
    }
  };

  public render() {
    return (
      <Modal.Window className="c-voters-modal" isOpen={this.props.isOpen} center={false} onClose={this.props.onClose}>
        <Modal.Content>
          {this.state.isLoading && <Loader />}
          <div>
            {!this.state.isLoading && (
              <List>
                {this.state.users.map(x => (
                  <ListItem key={x.id}>{x.name}</ListItem>
                ))}
              </List>
            )}
          </div>
        </Modal.Content>

        <Modal.Footer>
          <Button color="cancel" onClick={this.closeModal}>
            Close
          </Button>
        </Modal.Footer>
      </Modal.Window>
    );
  }
}
