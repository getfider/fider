import React from "react"
import { CurrentUser, Comment, Post } from "@fider/models"
import { ShowComment } from "./ShowComment"
import { CommentInput } from "./CommentInput"
import { VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

interface DiscussionPanelProps {
  user?: CurrentUser
  post: Post
  comments: Comment[]
  highlightedComment?: number
}

export const DiscussionPanel = (props: DiscussionPanelProps) => {
  return (
    <>
      <VStack spacing={2} className="c-comment-list mt-8">
        <span className="text-category">
          <Trans id="label.discussion">Discussion</Trans>
        </span>
        <VStack spacing={4} className="c-comment-list">
          {props.comments.map((c) => (
            <ShowComment key={c.id} post={props.post} comment={c} highlighted={props.highlightedComment === c.id} />
          ))}
          <CommentInput post={props.post} />
        </VStack>
      </VStack>
    </>
  )
}
