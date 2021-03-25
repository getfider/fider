import "./VotesPanel.scss"

import React, { useState } from "react"
import { Post, Vote } from "@fider/models"
import { Avatar } from "@fider/components"
import { Fider, classSet } from "@fider/services"
import { useFider } from "@fider/hooks"
import { VotesModal } from "./VotesModal"

interface VotesPanelProps {
  post: Post
  votes: Vote[]
}

export const VotesPanel = (props: VotesPanelProps) => {
  const fider = useFider()
  const [isVotesModalOpen, setIsVotesModalOpen] = useState(false)

  const openModal = () => {
    if (canShowAll()) {
      setIsVotesModalOpen(true)
    }
  }

  const closeModal = () => setIsVotesModalOpen(false)
  const canShowAll = () => fider.session.isAuthenticated && Fider.session.user.isCollaborator

  const extraVotesCount = props.post.votesCount - props.votes.length
  const moreVotesClassName = classSet({
    "l-votes-more": true,
    clickable: canShowAll(),
  })

  return (
    <>
      <VotesModal post={props.post} isOpen={isVotesModalOpen} onClose={closeModal} />
      <span className="subtitle">Voters</span>
      <div className="l-votes-list">
        {props.votes.map((x) => (
          <Avatar key={x.user.id} user={x.user} />
        ))}
        {extraVotesCount > 0 && (
          <span onClick={openModal} className={moreVotesClassName}>
            +{extraVotesCount} more
          </span>
        )}
        {props.votes.length > 0 && extraVotesCount === 0 && canShowAll() && (
          <span onClick={openModal} className={moreVotesClassName}>
            see details
          </span>
        )}
        {props.votes.length === 0 && <span className="info">None yet</span>}
      </div>
    </>
  )
}
