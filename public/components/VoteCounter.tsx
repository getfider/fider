import "./VoteCounter.scss";

import * as React from "react";
import { Post, PostStatus } from "@fider/models";
import { actions, device, classSet, Fider } from "@fider/services";
import { SignInModal } from "@fider/components";

interface VoteCounterProps {
  post: Post;
}

interface VoteCounterState {
  voted: boolean;
  total: number;
  showSignIn: boolean;
}

export class VoteCounter extends React.Component<VoteCounterProps, VoteCounterState> {
  constructor(props: VoteCounterProps) {
    super(props);
    this.state = {
      voted: props.post.viewerVoted,
      total: props.post.totalVotes,
      showSignIn: false
    };
  }

  private voteOrUndo = async () => {
    if (!Fider.session.isAuthenticated) {
      this.setState({ showSignIn: true });
      return;
    }

    const action = this.state.voted ? actions.removeVote : actions.addVote;

    const response = await action(this.props.post.number);
    if (response.ok) {
      this.setState(state => ({
        voted: !state.voted,
        total: state.total + (state.voted ? -1 : 1)
      }));
    }
  };

  public render() {
    const status = PostStatus.Get(this.props.post.status);

    const className = classSet({
      "m-voted": !status.closed && this.state.voted,
      "m-disabled": status.closed,
      "no-touch": !device.isTouch()
    });

    const vote = (
      <button className={className} onClick={this.voteOrUndo}>
        <i className="caret up icon" />
        {this.state.total}
      </button>
    );

    const disabled = (
      <button className={className}>
        <i className="caret up icon" />
        {this.state.total}
      </button>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} />
        <div className="c-vote-counter">{status.closed ? disabled : vote}</div>
      </>
    );
  }
}
