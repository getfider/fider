import "./Comments.scss"

import React from "react"
import { CurrentUser, Comment, Post } from "@fider/models"
import { ShowComment } from "./ShowComment"
import { CommentInput } from "./CommentInput"
import PostIllustration from "@fider/assets/images/undraw-post.svg"

interface DiscussionPanelProps {
  user?: CurrentUser
  post: Post
  comments: Comment[]
}

export const DiscussionPanel = (props: DiscussionPanelProps) => {
  return (
    <div className="comments-col">
      <div className="c-comment-list">
        <span className="subtitle">Discussion</span>
        <CommentInput post={props.post} />
        {props.comments.map((c) => (
          <ShowComment key={c.id} post={props.post} comment={c} />
        ))}
        {props.comments.length === 0 && (
          <div className="center">
            <p>
              <img alt="No Posts" height="100" src={PostIllustration} />
            </p>
            <p>No one has commented yet.</p>
          </div>
        )}
      </div>
    </div>
  )
}
