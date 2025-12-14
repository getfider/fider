import "./VotesPanel.scss"

import React, { useState } from "react"
import { Post, Vote } from "@fider/models"
import { Avatar, Icon } from "@fider/components"
import { Fider } from "@fider/services"
import { useFider } from "@fider/hooks"
import { VotesModal } from "./VotesModal"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import IconPerson from "@fider/assets/images/heroicons-person.svg"

interface VotesPanelProps {
  post: Post
  hideTitle?: boolean
  votes: Vote[]
}

const MAX_AVATARS_SHOWN = 12

export const VotesPanel = (props: VotesPanelProps) => {
  const fider = useFider()
  const [isVotesModalOpen, setIsVotesModalOpen] = useState(false)
  const canShowAll = fider.session.isAuthenticated && Fider.session.user.isCollaborator

  const openModal = () => {
    if (canShowAll) {
      setIsVotesModalOpen(true)
    }
  }

  const closeModal = () => setIsVotesModalOpen(false)

  const totalVotesCount = props.post.votesCount
  const visibleVotes = props.votes.slice(0, MAX_AVATARS_SHOWN)
  const remainingCount = totalVotesCount - MAX_AVATARS_SHOWN

  return (
    <VStack spacing={4} className="c-votes-panel card">
      <VotesModal post={props.post} isOpen={isVotesModalOpen} onClose={closeModal} />

      <div className="c-votes-panel__header">
        {!props.hideTitle && (
          <span className="text-bold text-gray-900">
            <Trans id="label.voters">Voters</Trans>
          </span>
        )}
        <div className="c-votes-panel__count-badge">{totalVotesCount}</div>
      </div>

      {props.votes.length > 0 && (
        <div className="c-votes-panel__avatars">
          {visibleVotes.map((vote, i) => (
            <Avatar key={i} user={vote.user} />
          ))}
          {remainingCount > 0 ? (
            <button className="c-votes-panel__more-avatar" onClick={openModal} disabled={!canShowAll}>
              +{remainingCount}
            </button>
          ) : (
            <button className="c-votes-panel__more-avatar" onClick={openModal} disabled={!canShowAll}>
              <Icon sprite={IconPerson} />
            </button>
          )}
        </div>
      )}

      {props.votes.length === 0 && (
        <span className="text-muted">
          <Trans id="label.none">None</Trans>
        </span>
      )}
    </VStack>
  )
}
