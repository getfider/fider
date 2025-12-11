import React, { useState } from "react"
import { Post, PostStatus } from "@fider/models"
import { actions } from "@fider/services"
import { Button, Icon, SignInModal } from "@fider/components"
import { useFider } from "@fider/hooks"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import { Trans } from "@lingui/react/macro"
import { HStack, VStack } from "@fider/components/layout"

interface VoteSectionProps {
  post: Post
  votes: number
}

export const VoteSection = (props: VoteSectionProps) => {
  const fider = useFider()
  const [votes, setVotes] = useState(props.votes)
  const [hasVoted, setHasVoted] = useState(props.post.hasVoted)
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)

  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
      return
    }

    const action = hasVoted ? actions.removeVote : actions.addVote
    const response = await action(props.post.number)
    if (response.ok) {
      setVotes(hasVoted ? votes - 1 : votes + 1)
      setHasVoted(!hasVoted)
    }
  }

  const hideModal = () => setIsSignInModalOpen(false)

  const status = PostStatus.Get(props.post.status)
  const isDisabled = status.closed || fider.isReadOnly

  const buttonText = hasVoted ? <Trans id="action.voted">Voted!</Trans> : <Trans id="action.vote">Vote for this idea</Trans>
  const icon = hasVoted ? IconCheck : IconThumbsUp

  return (
    <>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <VStack spacing={4}>
        <div className="align-self-start">
          <Button variant="primary" onClick={voteOrUndo} disabled={isDisabled} style={{ minWidth: "180px" }}>
            <HStack spacing={2} justify="center" className="w-full">
              <Icon sprite={icon} /> <span>{buttonText}</span>
            </HStack>
          </Button>
        </div>
        <HStack align="center" spacing={2}>
          <span className="text-semibold text-2xl" style={{ fontSize: "32px", minHeight: "48px" }}>
            {votes}
          </span>
          <span className="text-semibold text-lg">{votes === 1 ? <Trans id="label.vote">Vote</Trans> : <Trans id="label.votes">Votes</Trans>}</span>
        </HStack>
      </VStack>
    </>
  )
}
