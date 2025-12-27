import React, { useState } from "react"
import { Post, Vote } from "@fider/models"
import { AvatarStack, Button } from "@fider/components"
import { Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { VotesModal } from "./VotesModal"
import { ManageVotersModal } from "./ManageVotersModal"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

interface VotesPanelProps {
  post: Post
  hideTitle?: boolean
  votes: Vote[]
  onVotesChanged?: (delta: number) => void
}

export const VotesPanel = (props: VotesPanelProps) => {
  const fider = useFider()
  const [isVotesModalOpen, setIsVotesModalOpen] = useState(false)
  const [isManageVotersModalOpen, setIsManageVotersModalOpen] = useState(false)
  const canShowAll = fider.session.isAuthenticated && Fider.session.user.isCollaborator
  const isAdmin = fider.session.isAuthenticated && fider.session.user.isAdministrator

  const openModal = () => {
    if (isAdmin) {
      setIsManageVotersModalOpen(true)
    } else if (canShowAll) {
      setIsVotesModalOpen(true)
    }
  }

  const closeModal = () => {
    setIsVotesModalOpen(false)
    setIsManageVotersModalOpen(false)
  }

  const handleVotesChanged = (delta: number) => {
    if (props.onVotesChanged) {
      props.onVotesChanged(delta)
    }
  }

  const extraVotesCount = props.post.votesCount - props.votes.length

  return (
    <VStack spacing={4}>
      <VotesModal post={props.post} isOpen={isVotesModalOpen} onClose={closeModal} />
      <ManageVotersModal post={props.post} isOpen={isManageVotersModalOpen} onClose={closeModal} onVotesChanged={handleVotesChanged} />
      {!props.hideTitle && (
        <span className="text-category">
          <Trans id="label.voters">Voters</Trans>
        </span>
      )}
      {props.votes.length > 0 && (
        <HStack spacing={0} className="gap-2">
          <AvatarStack users={props.votes.map((x) => x.user)} overlap={false} />
        </HStack>
      )}
      {extraVotesCount > 0 && (
        <Button variant="tertiary" disabled={!canShowAll} size="small" onClick={openModal}>
          <Trans id="showpost.votespanel.more">+{extraVotesCount} more</Trans>
        </Button>
      )}
      {props.votes.length > 0 && extraVotesCount === 0 && canShowAll && (
        <Button variant="tertiary" size="small" disabled={!canShowAll} onClick={openModal}>
          <Trans id="showpost.votespanel.seedetails">see details</Trans>
        </Button>
      )}
      {props.votes.length === 0 && (
        <span className="text-muted">
          <Trans id="label.none">None</Trans>
        </span>
      )}
    </VStack>
  )
}
