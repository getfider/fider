import "./VotersPanel.scss";

import React from "react";
import { User, Post } from "@fider/models";
import { Gravatar } from "@fider/components";
import { Fider, classSet } from "@fider/services";
import { VotersModal } from "./VotersModal";

interface VotersPanelProps {
  post: Post;
  voters: {
    total: number;
    list: User[];
  };
}

interface VotersPanelState {
  showModal: boolean;
}

export class VotersPanel extends React.Component<VotersPanelProps, VotersPanelState> {
  constructor(props: VotersPanelProps) {
    super(props);
    this.state = {
      showModal: false
    };
  }

  private showModal = () => {
    this.setState({ showModal: true });
  };

  private hideModal = () => {
    this.setState({ showModal: false });
  };

  public render() {
    const extraVotersCount = this.props.voters.total - this.props.voters.list.length;
    const canShowAll = Fider.session.isAuthenticated && Fider.session.user.isCollaborator;
    const moreVotersClassName = classSet({
      "l-voters-more": true,
      clickable: canShowAll
    });

    return (
      <div>
        <VotersModal post={this.props.post} isOpen={this.state.showModal} onClose={this.hideModal} />
        <span className="subtitle">Voters</span>
        <div className="l-voters-list">
          {this.props.voters.list.map(x => (
            <Gravatar key={x.id} user={x} />
          ))}
          {extraVotersCount > 0 && (
            <span onClick={this.showModal} className={moreVotersClassName}>
              +{extraVotersCount} more
            </span>
          )}
        </div>
      </div>
    );
  }
}
