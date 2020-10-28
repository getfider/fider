import "./VotesPanel.scss";

import React, { useState } from "react";
import { Post, Vote } from "@fider/models";
import { Avatar } from "@fider/components";
import { Fider, classSet } from "@fider/services";
import { useFider } from "@fider/hooks";
import { useTranslation } from "react-i18next";
import { VotesModal } from "./VotesModal";

interface VotesPanelProps {
  post: Post;
  votes: Vote[];
}

export const VotesPanel = (props: VotesPanelProps) => {
  const fider = useFider();
  const { t } = useTranslation();
  const [isVotesModalOpen, setIsVotesModalOpen] = useState(false);

  const openModal = () => {
    if (canShowAll()) {
      setIsVotesModalOpen(true);
    }
  };

  const closeModal = () => setIsVotesModalOpen(false);
  const canShowAll = () => fider.session.isAuthenticated && Fider.session.user.isCollaborator;

  const extraVotesCount = props.post.votesCount - props.votes.length;
  const moreVotesClassName = classSet({
    "l-votes-more": true,
    clickable: canShowAll()
  });

  return (
    <>
      <VotesModal post={props.post} isOpen={isVotesModalOpen} onClose={closeModal} />
      <span className="subtitle">{t("showPost.votesPanel.subtitle")}</span>
      <div className="l-votes-list">
        {props.votes.map(x => (
          <Avatar key={x.user.id} user={x.user} />
        ))}
        {extraVotesCount > 0 && (
          <span onClick={openModal} className={moreVotesClassName}>
            +{extraVotesCount} {t("showPost.votesPanel.more")}
          </span>
        )}
        {props.votes.length > 0 && extraVotesCount === 0 && canShowAll() && (
          <span onClick={openModal} className={moreVotesClassName}>
            {t("showPost.votesPanel.seeDetails")}
          </span>
        )}
        {props.votes.length === 0 && <span className="info">{t("showPost.votesPanel.noneYet")}</span>}
      </div>
    </>
  );
};
