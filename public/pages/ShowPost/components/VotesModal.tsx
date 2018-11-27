import "./VotesModal.scss";

import React from "react";
import { Post, Vote } from "@fider/models";
import { Modal, Button, Loader, List, ListItem } from "@fider/components";
import { actions } from "@fider/services";

interface VotesModalProps {
  isOpen: boolean;
  post: Post;
  onClose?: () => void;
}

interface VotesModalState {
  searchText: string;
  votes: Vote[];
  isLoading: boolean;
}

export class VotesModal extends React.Component<VotesModalProps, VotesModalState> {
  constructor(props: VotesModalProps) {
    super(props);
    this.state = {
      searchText: "",
      votes: [],
      isLoading: true
    };
  }

  public componentDidMount() {
    actions.listVotes(this.props.post.number).then(response => {
      if (response.ok) {
        this.setState({ votes: response.data, isLoading: false });
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
      <Modal.Window className="c-votes-modal" isOpen={this.props.isOpen} center={false} onClose={this.props.onClose}>
        <Modal.Content>
          {this.state.isLoading && <Loader />}
          <div>
            {!this.state.isLoading && (
              <List>
                {this.state.votes.map(x => (
                  <ListItem key={x.user.id}>{x.user.name}</ListItem>
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
