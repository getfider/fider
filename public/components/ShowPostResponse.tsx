import React from "react"
import { PostResponse, PostStatus } from "@fider/models"
import { Icon, Markdown, UserName, Moment } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import HeroIconCheck from "@fider/assets/images/heroicons-check-circle.svg"
import HeroIconSparkles from "@fider/assets/images/heroicons-sparkles-outline.svg"
import HeroIconLightBulb from "@fider/assets/images/heroicons-lightbulb.svg"
import HeroIconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import HeroIconThumbsDown from "@fider/assets/images/heroicons-thumbsdown.svg"
import { HStack } from "./layout"
import { Trans } from "@lingui/react/macro"
import { useFider } from "@fider/hooks"

type Size = "micro" | "small" | "xsmall" | "normal"

interface PostResponseProps {
  status: string
  response: PostResponse | null
  size?: Size
}

export const ResponseDetails = (props: PostResponseProps): JSX.Element | null => {
  const fider = useFider()
  const status = PostStatus.Get(props.status)

  if (!props.response) {
    return null
  }

  const { bg, border } = getLozengeProps(status)

  return (
    <HStack spacing={4} align="start" className="c-response-details">
      <div className={`c-response-details__card ${bg} border ${border}`}>
        <div className="c-response-details__header">
          <HStack spacing={2} align="center">
            <UserName user={props.response.user} />
            <span className="text-xs text-gray-600">•</span>
            <Moment className="text-xs text-gray-600" locale={fider.currentLocale} date={props.response.respondedAt} />
          </HStack>
        </div>

        {props.response?.text && status !== PostStatus.Duplicate && (
          <div className="c-response-details__content">
            <Markdown text={props.response.text} style="full" />
          </div>
        )}

        {status === PostStatus.Duplicate && props.response.original && (
          <div className="c-response-details__content">
            <a className="text-link" href={`/posts/${props.response.original.number}/${props.response.original.slug}`}>
              {props.response.original.title}
            </a>
          </div>
        )}
      </div>
    </HStack>
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
    case PostStatus.Open:
      return { icon: HeroIconLightBulb, bg: "bg-blue-100", color: "text-blue-700", border: "border-blue-400" }
    default:
      return { icon: HeroIconSparkles, bg: "bg-green-100", color: "text-green-700", border: "border-green-400" }
  }
}

const getStatusTranslation = (status: PostStatus): JSX.Element => {
  switch (status) {
    case PostStatus.Open:
      return <Trans id="enum.poststatus.open">Open</Trans>
    case PostStatus.Planned:
      return <Trans id="enum.poststatus.planned">Planned</Trans>
    case PostStatus.Started:
      return <Trans id="enum.poststatus.started">Started</Trans>
    case PostStatus.Completed:
      return <Trans id="enum.poststatus.completed">Completed</Trans>
    case PostStatus.Declined:
      return <Trans id="enum.poststatus.declined">Declined</Trans>
    case PostStatus.Duplicate:
      return <Trans id="enum.poststatus.duplicate">Duplicate</Trans>
    case PostStatus.Deleted:
      return <Trans id="enum.poststatus.deleted">Deleted</Trans>
    default:
      return <>{status.title}</>
  }
}

export const ResponseLozenge = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)
  const { icon, bg, color, border } = getLozengeProps(status)
  const translatedStatus = getStatusTranslation(status)

  if (props.size == "micro") {
    return <span className={`${color} text-sm`}>{translatedStatus}</span>
  }

  if (props.size === "xsmall") {
    return (
      <div>
        <HStack align="center" className={`${color} ${bg} rounded-full p-0 px-3`}>
          <Icon sprite={icon} className={`h-4 c-status-col--${status.value}`} />
          <span className={`c-status-col--${status.value} text-xs uppercase`}>{translatedStatus}</span>
        </HStack>
      </div>
    )
  }

  return (
    <div>
      <HStack align="start" className={`${color} ${bg} border ${border} rounded-full p-1 px-3`}>
        {!props.size && <Icon sprite={icon} className={`h-5 c-status-col--${status.value}`} />}
        <span className={`c-status-col--${status.value} ${props.size === "small" ? "text-sm" : "text-semibold"}`}>{translatedStatus}</span>
      </HStack>
    </div>
  )
}
