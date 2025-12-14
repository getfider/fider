import "./VoteCounter.scss"

import React, { useState } from "react"
import { Post, PostStatus } from "@fider/models"
import { actions, classSet } from "@fider/services"
import { Icon, SignInModal } from "@fider/components"
import { useFider } from "@fider/hooks"
import ChevronUp from "@fider/assets/images/chevron-up.svg"

export interface VoteCounterProps {
  post: Post
  size?: "default" | "large"
}

export const VoteCounter = (props: VoteCounterProps) => {
  const fider = useFider()
  const { size = "default" } = props
  const [hasVoted, setHasVoted] = useState(props.post.hasVoted)
  const [votesCount, setVotesCount] = useState(props.post.votesCount)
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)

  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
      return
    }

    const action = hasVoted ? actions.removeVote : actions.addVote

    const response = await action(props.post.number)
    if (response.ok) {
      setVotesCount(votesCount + (hasVoted ? -1 : 1))
      setHasVoted(!hasVoted)
    }
  }

  const hideModal = () => setIsSignInModalOpen(false)

  const status = PostStatus.Get(props.post.status)
  const isDisabled = status.closed || fider.isReadOnly

  const className = classSet({
    "c-vote-counter__button": true,
    "c-vote-counter__button--voted": !status.closed && hasVoted,
    "c-vote-counter__button--disabled": isDisabled,
    "c-vote-counter__button--large": size === "large",
  })

  const vote = (
    <button className={className} onClick={voteOrUndo}>
      {!hasVoted && <Icon sprite={ChevronUp} className="c-vote-counter__icon" />}
      <span className="c-vote-counter__count">{votesCount}</span>
    </button>
  )

  const disabled = (
    <button className={className}>
      {!hasVoted && <Icon sprite={ChevronUp} className="c-vote-counter__icon" />}
      <span className="c-vote-counter__count">{votesCount}</span>
    </button>
  )

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <div className="c-vote-counter">{isDisabled ? disabled : vote}</div>
    </>
  )
}
