import "./ShowPost.page.scss"

import React from "react"

import { Comment, Post, Tag, Vote } from "@fider/models"
import { Header, PoweredByFider } from "@fider/components"
import { PostDetails } from "@fider/components/PostDetails"

interface ShowPostPageProps {
  post: Post
  subscribed: boolean
  comments: Comment[]
  tags: Tag[]
  votes: Vote[]
  attachments: string[]
}

export default function ShowPostPage(props: ShowPostPageProps) {
  return (
    <>
      <Header />
      <div id="p-show-post" className="page container">
        <PostDetails
          postNumber={props.post.number}
          initialPost={props.post}
          initialSubscribed={props.subscribed}
          initialComments={props.comments}
          initialTags={props.tags}
          initialVotes={props.votes}
          initialAttachments={props.attachments}
        />
        <PoweredByFider slot="show-post" className="mt-3" />
      </div>
    </>
  )
}
