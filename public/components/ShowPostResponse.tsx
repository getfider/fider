import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Markdown, UserName, ShowPostStatus } from "@fider/components"
import { HStack } from "./layout"

const DuplicateDetails = (props: PostResponseProps): JSX.Element | null => {
  if (!props.response) {
    return null
  }

  const original = props.response.original
  if (!original) {
    return null
  }

  return (
    <div className="content">
      <span>&#8618;</span>{" "}
      <a className="text-link" href={`/posts/${original.number}/${original.slug}`}>
        {original.title}
      </a>
    </div>
  )
}

interface PostResponseProps {
  status: string
  response: PostResponse | null
}

const StatusDetails = (props: PostResponseProps): JSX.Element | null => {
  if (!props.response || !props.response.text) {
    return null
  }

  return (
    <div className="content">
      <Markdown text={props.response.text} style="full" />
    </div>
  )
}

export const ShowPostResponse = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)

  if (props.response && (status.show || props.response.text)) {
    return (
      <div className="p-2 shadow rounded">
        <HStack>
          {status.show && <ShowPostStatus status={status} />}
          <span className="text-xs">
            &middot; <UserName user={props.response.user} />
          </span>
        </HStack>
        {status === PostStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}
      </div>
    )
  }

  return null
}
