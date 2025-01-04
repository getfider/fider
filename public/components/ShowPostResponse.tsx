import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Icon, Markdown } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
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
    <HStack>
      <Icon sprite={HeroIconDuplicate} className="h-6 text-gray-500" />
      <a className="text-link" href={`/posts/${original.number}/${original.slug}`}>
        {original.title}
      </a>
    </HStack>
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
      <div>
        {status.show && <div className={`p-2 text-white text-center text-semibold c-status-bg--${status.value}`}>{status.title}</div>}
        <div className="pt-2">{status === PostStatus.Duplicate ? DuplicateDetails(props) : StatusDetails(props)}</div>
      </div>
    )
  }

  return null
}
