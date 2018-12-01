import "./VotesPanel.scss";

import React from "react";
import { Post, Vote } from "@fider/models";
import { Gravatar } from "@fider/components";
import { Fider, classSet } from "@fider/services";
import { VotesModal } from "./VotesModal";

interface VotesPanelProps {
  post: Post;
  votes: Vote[];
}

interface VotesPanelState {
  showModal: boolean;
}

export class VotesPanel extends React.Component<VotesPanelProps, VotesPanelState> {
  constructor(props: VotesPanelProps) {
    super(props);
    this.state = {
      showModal: false
    };
  }

  private showModal = () => {
    if (this.canShowAll()) {
      this.setState({ showModal: true });
    }
  };

  private hideModal = () => {
    this.setState({ showModal: false });
  };

  private canShowAll = () => {
    return Fider.session.isAuthenticated && Fider.session.user.isCollaborator;
  };

  public render() {
    const extraVotesCount = this.props.post.votesCount - this.props.votes.length;
    const moreVotesClassName = classSet({
      "l-votes-more": true,
      clickable: this.canShowAll()
    });

    return (
      <>
        <VotesModal post={this.props.post} isOpen={this.state.showModal} onClose={this.hideModal} />
        <span className="subtitle">Voters</span>
        <div className="l-votes-list">
          {this.props.votes.map(x => (
            <Gravatar key={x.user.id} user={x.user} />
          ))}
          {extraVotesCount > 0 && (
            <span onClick={this.showModal} className={moreVotesClassName}>
              +{extraVotesCount} more
            </span>
          )}
          {this.props.votes.length > 0 && extraVotesCount === 0 && this.canShowAll() && (
            <span onClick={this.showModal} className={moreVotesClassName}>
              see details
            </span>
          )}
          {this.props.votes.length === 0 && <span className="info">None yet</span>}
        </div>
      </>
    );
  }
}
