import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Icon, Markdown } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import HeroIconCheck from "@fider/assets/images/heroicons-check-circle.svg"
import HeroIconSparkles from "@fider/assets/images/heroicons-sparkles-outline.svg"
import HeroIconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import HeroIconThumbsDown from "@fider/assets/images/heroicons-thumbsdown.svg"
import { HStack, VStack } from "./layout"
import { timeSince } from "@fider/services"

interface PostResponseProps {
  status: string
  response: PostResponse | null
  small?: boolean
}

export const ResponseDetails = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)

  if (!props.response || status === PostStatus.Open) {
    return null
  }

  return (
    <VStack align="start" spacing={4} className="bg-blue-50 p-3 border border-blue-200 rounded">
      <ResponseLozenge response={props.response} status={props.status} />
      <div className="text-semibold text-lg">{timeSince("en", new Date(), props.response.respondedAt, "date")}</div>
      {props.response?.text && status !== PostStatus.Duplicate && (
        <div className="content">
          <Markdown text={props.response.text} style="full" />
        </div>
      )}

      {status === PostStatus.Duplicate && props.response.original && (
        <div className="content">
          <a className="text-link" href={`/posts/${props.response.original.number}/${props.response.original.slug}`}>
            {props.response.original.title}
          </a>
        </div>
      )}
    </VStack>
  )
}

const getLozengeProps = (status: PostStatus): { icon: SpriteSymbol; bg: string; color: string; border: string } => {
  switch (status) {
    case PostStatus.Declined:
      return { icon: HeroIconThumbsDown, bg: "bg-red-100", color: "text-red-800", border: "border-red-300" }
    case PostStatus.Duplicate:
      return { icon: HeroIconDuplicate, bg: "bg-yellow-100", color: "text-yellow-800", border: "border-yellow-400" }
    case PostStatus.Completed:
      return { icon: HeroIconCheck, bg: "bg-green-300", color: "text-green-800", border: "border-green-500" }
    case PostStatus.Planned:
      return { icon: HeroIconThumbsUp, bg: "bg-blue-100", color: "text-blue-700", border: "border-blue-400" }
    default:
      return { icon: HeroIconSparkles, bg: "bg-green-100", color: "text-green-700", border: "border-green-400" }
  }
}

export const ResponseLozenge = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)
  const { icon, bg, color, border } = getLozengeProps(status)

  if (status === PostStatus.Open) {
    return <div />
  }

  return (
    <div>
      <HStack align="start" className={`${color} ${bg} border ${border} rounded-full p-1 px-3`}>
        {!props.small && <Icon sprite={icon} className={`h-5 c-status-col--${status.value}`} />}
        <span className={`c-status-col--${status.value} ${props.small ? "text-sm" : "text-semibold"}`}>{status.title}</span>
      </HStack>
    </div>
  )
}
