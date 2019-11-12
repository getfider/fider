import "./VoteCounter.scss";

import React, { useState } from "react";
import { Post, PostStatus } from "@fider/models";
import { actions, device, classSet } from "@fider/services";
import { SignInModal } from "@fider/components";
import { FaCaretUp } from "react-icons/fa";
import { useFider } from "@fider/hooks";

interface VoteCounterProps {
  post: Post;
}

export const VoteCounter = (props: VoteCounterProps) => {
  const fider = useFider();
  const [hasVoted, setHasVoted] = useState(props.post.hasVoted);
  const [votesCount, setVotesCount] = useState(props.post.votesCount);
  const [showSignIn, setShowSignIn] = useState(false);

  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setShowSignIn(true);
      return;
    }

    const action = hasVoted ? actions.removeVote : actions.addVote;

    const response = await action(props.post.number);
    if (response.ok) {
      setVotesCount(votesCount + (hasVoted ? -1 : 1));
      setHasVoted(!hasVoted);
    }
  };

  const hideModal = () => {
    setShowSignIn(false);
  };

  const status = PostStatus.Get(props.post.status);

  const className = classSet({
    "m-voted": !status.closed && hasVoted,
    "m-disabled": status.closed,
    "no-touch": !device.isTouch()
  });

  const vote = (
    <button className={className} onClick={voteOrUndo}>
      <FaCaretUp />
      {votesCount}
    </button>
  );

  const disabled = (
    <button className={className}>
      <FaCaretUp />
      {votesCount}
    </button>
  );

  return (
    <>
      <SignInModal isOpen={showSignIn} onClose={hideModal} />
      <div className="c-vote-counter">{status.closed ? disabled : vote}</div>
    </>
  );
};
