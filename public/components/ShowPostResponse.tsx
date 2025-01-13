import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Icon, Markdown } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import HeroIconSmile from "@fider/assets/images/heroicons-smile.svg"
import HeroIconStar from "@fider/assets/images/heroicons-star.svg"
import HeroIconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import HeroIconInbox from "@fider/assets/images/heroicons-inbox.svg"
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
      return HeroIconStar
    case PostStatus.Planned:
      return HeroIconThumbsUp
    case PostStatus.Started:
      return HeroIconSmile
  }
  return HeroIconInbox
}

export const ResponseStatusLabel = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)

  return (
    <>
      <HStack>
        <Icon sprite={getStatusIcon(status)} className={`h-6 c-status-col--${status.value}`} />
        <span className={`text-semibold text-lg c-status-col--${status.value}`}>{status.title}</span>
      </HStack>
    </>
  )
}
