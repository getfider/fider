import React, { useState } from "react"
import { Post, PostStatus, Vote } from "@fider/models"
import { actions } from "@fider/services"
import { Button, Icon, SignInModal } from "@fider/components"
import { useFider } from "@fider/hooks"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { Trans } from "@lingui/macro"
import { HStack, VStack } from "@fider/components/layout"
// import { VotesPanel } from "./VotesPanel"

interface VoteSectionProps {
  post: Post
  votes: Vote[]
}

export const VoteSection = (props: VoteSectionProps) => {
  const fider = useFider()
  const [hasVoted, setHasVoted] = useState(props.post.hasVoted)
  const [votes, setVotes] = useState(props.votes)
  const [post, setPost] = useState(props.post)
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)

  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
      return
    }

    const action = hasVoted ? actions.removeVote : actions.addVote

    const response = await action(props.post.number)
    if (response.ok) {
      const newVotes: Vote[] = hasVoted
        ? votes.filter((vote) => vote.user.id !== fider.session.user.id)
        : [...votes, { user: fider.session.user, createdAt: new Date() }]

      setPost({ ...post, votesCount: newVotes.length })
      setVotes(newVotes)
      setHasVoted(!hasVoted)
    }
  }

  const hideModal = () => setIsSignInModalOpen(false)

  const status = PostStatus.Get(props.post.status)
  const isDisabled = status.closed || fider.isReadOnly

  const buttonText = hasVoted ? <Trans id="action.voted">Voted!</Trans> : <Trans id="action.vote">Vote for this idea</Trans>
  const icon = hasVoted ? IconCheck : IconThumbsUp

  return (
    <VStack spacing={4}>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <div className="align-self-start">
        <Button variant="primary" onClick={voteOrUndo} disabled={isDisabled} style={{ minWidth: "180px" }}>
          <HStack spacing={2} justify="center" className="w-full">
            <Icon sprite={icon} /> <span>{buttonText}</span>
          </HStack>
        </Button>
      </div>
      <HStack>
        <span className="text-display">
          {votes.length} {votes.length === 1 ? "Vote" : "Votes"}
        </span>
        {/* <div className="pl-4">
          <VotesPanel post={post} votes={votes} hideTitle={true} />
        </div> */}
      </HStack>
    </VStack>
  )
}
