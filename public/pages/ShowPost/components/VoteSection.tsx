import React, { useState } from "react"
import { Post, PostStatus } from "@fider/models"
import { actions } from "@fider/services"
import { Button, Icon, SignInModal } from "@fider/components"
import { useFider } from "@fider/hooks"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import IconPlusCircle from "@fider/assets/images/heroicons-pluscircle.svg"
import { Trans } from "@lingui/macro"
import { HStack, VStack } from "@fider/components/layout"
import { ManageVotersModal } from "./ManageVotersModal"

interface VoteSectionProps {
  post: Post
  votes: number
}

export const VoteSection = (props: VoteSectionProps) => {
  const fider = useFider()
  const [votes, setVotes] = useState(props.votes)
  const [hasVoted, setHasVoted] = useState(props.post.hasVoted)
  const [isSignInModalOpen, setIsSignInModalOpen] = useState(false)
  const [isManageVotersModalOpen, setIsManageVotersModalOpen] = useState(false)

  const voteOrUndo = async () => {
    if (!fider.session.isAuthenticated) {
      setIsSignInModalOpen(true)
      return
    }

    const response = await actions.toggleVote(props.post.number)
    if (response.ok) {
      setVotes(response.data.voted ? votes + 1 : votes - 1)
      setHasVoted(response.data.voted)
    }
  }

  const hideModal = () => setIsSignInModalOpen(false)
  const openManageVotersModal = () => setIsManageVotersModalOpen(true)
  const closeManageVotersModal = () => setIsManageVotersModalOpen(false)
  const handleVotesChanged = (delta: number) => setVotes(votes + delta)

  const status = PostStatus.Get(props.post.status)
  const isDisabled = status.closed || fider.isReadOnly
  const canManageVoters = fider.session.isAuthenticated && fider.session.user.isAdministrator && !status.closed && !fider.isReadOnly

  const buttonText = hasVoted ? <Trans id="action.voted">Voted!</Trans> : <Trans id="action.vote">Vote for this idea</Trans>
  const icon = hasVoted ? IconCheck : IconThumbsUp

  return (
    <VStack spacing={4}>
      <SignInModal isOpen={isSignInModalOpen} onClose={hideModal} />
      <ManageVotersModal isOpen={isManageVotersModalOpen} onClose={closeManageVotersModal} post={props.post} onVotesChanged={handleVotesChanged} />
      <div className="align-self-start">
        <HStack spacing={2}>
          <Button variant="primary" onClick={voteOrUndo} disabled={isDisabled} style={{ minWidth: "180px" }}>
            <HStack spacing={2} justify="center" className="w-full">
              <Icon sprite={icon} /> <span>{buttonText}</span>
            </HStack>
          </Button>
          {canManageVoters && (
            <Button variant="secondary" onClick={openManageVotersModal}>
              <Icon sprite={IconPlusCircle} />
            </Button>
          )}
        </HStack>
      </div>
      <HStack align="center">
        <span className="text-semibold text-2xl" style={{ fontSize: "32px", minHeight: "48px" }}>
          {votes}
        </span>
        <span className="text-semibold text-lg">{votes === 1 ? "Vote" : "Votes"}</span>
      </HStack>
    </VStack>
  )
}
