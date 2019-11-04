import "./VoteCounter.scss";

import React from "react";
import { Post, PostStatus } from "@fider/models";
import { actions, device, classSet, Fider } from "@fider/services";
import { SignInModal } from "@fider/components";
import { FaCaretUp } from "react-icons/fa";

interface VoteCounterProps {
  post: Post;
}

interface VoteCounterState {
  voted: boolean;
  count: number;
  showSignIn: boolean;
}

export class VoteCounter extends React.Component<VoteCounterProps, VoteCounterState> {
  constructor(props: VoteCounterProps) {
    super(props);
    this.state = {
      voted: props.post.hasVoted,
      count: props.post.votesCount,
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
        count: state.count + (state.voted ? -1 : 1)
      }));
    }
  };

  private hideModal = () => {
    this.setState({ showSignIn: false });
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
        <FaCaretUp />
        {this.state.count}
      </button>
    );

    const disabled = (
      <button className={className}>
        <FaCaretUp />
        {this.state.count}
      </button>
    );

    return (
      <>
        <SignInModal isOpen={this.state.showSignIn} onClose={this.hideModal} />
        <div className="c-vote-counter">{status.closed ? disabled : vote}</div>
      </>
    );
  }
}
