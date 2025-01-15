import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Icon, Markdown } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import HeroIconCheck from "@fider/assets/images/heroicons-check-circle.svg"
import HeroIconSparkles from "@fider/assets/images/heroicons-sparkles-outline.svg"
import HeroIconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import { t } from "@lingui/macro"

import { HStack } from "./layout"
import { timeSince } from "@fider/services"

// const DuplicateDetails = (props: PostResponseProps): JSX.Element | null => {
//   if (!props.response) {
//     return null
//   }

//   const original = props.response.original
//   if (!original) {
//     return null
//   }

//   return (
//     <HStack>
//       <Icon sprite={HeroIconDuplicate} className="h-6 text-gray-500" />
//       <a className="text-link" href={`/posts/${original.number}/${original.slug}`}>
//         {original.title}
//       </a>
//     </HStack>
//   )
// }

interface PostResponseProps {
  status: string
  response: PostResponse | null
}

export const ResponseDetails = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)
  const id = `enum.poststatus.${status.value}`
  const statusLabel = t({ id, message: status.title })

  if (!props.response || status === PostStatus.Open) {
    return null
  }

  const title = t({
    id: "showpost.response.date",
    message: "Status changed to {status} on {statusDate}",
    values: { status: statusLabel, statusDate: timeSince("en", new Date(), props.response.respondedAt, "date") },
  })

  return (
    <div className="bg-blue-50 p-2 border border-blue-200 rounded">
      <div>{title}</div>
      {props.response?.text && status !== PostStatus.Duplicate && (
        <div className="content pt-1">
          <Markdown text={props.response.text} style="full" />
        </div>
      )}

      {status === PostStatus.Duplicate && props.response.original && (
        <div className="content pt-1">
          <a className="text-link" href={`/posts/${props.response.original.number}/${props.response.original.slug}`}>
            {props.response.original.title}
          </a>
        </div>
      )}
    </div>
  )
}

const getStatusIcon = (status: PostStatus): SpriteSymbol => {
  switch (status) {
    case PostStatus.Duplicate:
      return HeroIconDuplicate
    case PostStatus.Completed:
      return HeroIconCheck
    case PostStatus.Planned:
      return HeroIconThumbsUp
    case PostStatus.Started:
      return HeroIconSparkles
  }
  return HeroIconThumbsUp
}

export const ResponseLozenge = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)

  if (status === PostStatus.Open) {
    return <div />
  }

  return (
    <>
      <HStack align="start" className="align-self-start bg-blue-100 border border-blue-300 rounded-full p-1 px-3 mb-4">
        <Icon sprite={getStatusIcon(status)} className={`h-5 c-status-col--${status.value}`} />
        <span className={`text-semibold c-status-col--${status.value}`}>{status.title}</span>
      </HStack>
    </>
  )
}
