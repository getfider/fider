import "./Comments.scss";

import React from "react";
import { CurrentUser, Comment, Post } from "@fider/models";
import { ShowComment } from "./ShowComment";
import { CommentInput } from "./CommentInput";
import { useTranslation } from "react-i18next";

interface DiscussionPanelProps {
  user?: CurrentUser;
  post: Post;
  comments: Comment[];
}

export const DiscussionPanel = (props: DiscussionPanelProps) => {
  const { t } = useTranslation();
  return (
    <div className="comments-col">
      <div className="c-comment-list">
        <span className="subtitle">{t("showPost.discussionPanel.discussion")}</span>
        {props.comments.map(c => (
          <ShowComment key={c.id} post={props.post} comment={c} />
        ))}
        <CommentInput post={props.post} />
      </div>
    </div>
  );
};
