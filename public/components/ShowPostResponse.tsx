import React from "react"
import { PostResponse, PostStatus, Status, resolveStatus } from "@fider/models"
import { Fider } from "@fider/services"
import { Icon, Markdown, UserName, Moment, Avatar } from "@fider/components"
import HeroIconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import HeroIconCheck from "@fider/assets/images/heroicons-check-circle.svg"
import HeroIconSparkles from "@fider/assets/images/heroicons-sparkles-outline.svg"
import HeroIconLightBulb from "@fider/assets/images/heroicons-lightbulb.svg"
import HeroIconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import HeroIconThumbsDown from "@fider/assets/images/heroicons-thumbsdown.svg"
import { HStack, VStack } from "./layout"
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

  return (
    <HStack spacing={4} align="start" className="c-response-details">
      <Avatar user={props.response.user} size="large" />
      <div className="c-response-details__card">
        <div className="c-response-details__inner">
          <VStack spacing={2}>
            <HStack spacing={2} align="center">
              <UserName user={props.response.user} />
              <span className="text-xs text-gray-600">•</span>
              <Moment className="text-xs text-gray-600" locale={fider.currentLocale} date={props.response.respondedAt} />
              <ResponseLozenge status={props.status} response={props.response} size="xsmall" />
            </HStack>

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
          </VStack>
        </div>
      </div>
    </HStack>
  )
}

// Tailwind class palette per tenant-configurable color. Keep keys in sync with
// the COLOR_OPTIONS list in ManageStatuses.page.tsx.
const colorPalette: Record<string, { bg: string; color: string; border: string }> = {
  red: { bg: "bg-red-100", color: "text-red-800", border: "border-red-300" },
  yellow: { bg: "bg-yellow-100", color: "text-yellow-800", border: "border-yellow-400" },
  green: { bg: "bg-green-100", color: "text-green-800", border: "border-green-400" },
  blue: { bg: "bg-blue-100", color: "text-blue-700", border: "border-blue-400" },
  gray: { bg: "bg-gray-100", color: "text-gray-700", border: "border-gray-300" },
}

// Icon defaults per semantic kind, used when a tenant-defined status doesn't
// match the legacy hardcoded slugs.
const iconForKind = (kind: string): SpriteSymbol => {
  switch (kind) {
    case "closed-declined":
      return HeroIconThumbsDown
    case "closed-completed":
      return HeroIconCheck
    case "duplicate":
      return HeroIconDuplicate
    case "active":
      return HeroIconSparkles
    case "open":
    default:
      return HeroIconLightBulb
  }
}

const getLozengeProps = (
  status: PostStatus,
  tenantStatus: Status | null
): { icon: SpriteSymbol; bg: string; color: string; border: string } => {
  // Tenant catalogue takes precedence — that's where admin-chosen color/kind
  // live for custom statuses (feedback.fider.io/111).
  if (tenantStatus) {
    const palette = colorPalette[tenantStatus.color] || colorPalette.blue
    return { icon: iconForKind(tenantStatus.kind), ...palette }
  }
  switch (status) {
    case PostStatus.Declined:
      return { icon: HeroIconThumbsDown, ...colorPalette.red }
    case PostStatus.Duplicate:
      return { icon: HeroIconDuplicate, ...colorPalette.yellow }
    case PostStatus.Completed:
      return { icon: HeroIconCheck, ...colorPalette.green }
    case PostStatus.Planned:
      return { icon: HeroIconThumbsUp, ...colorPalette.blue }
    case PostStatus.Started:
      return { icon: HeroIconSparkles, ...colorPalette.blue }
    case PostStatus.Open:
      return { icon: HeroIconLightBulb, ...colorPalette.blue }
    default:
      return { icon: HeroIconSparkles, ...colorPalette.blue }
  }
}

const getStatusTranslation = (status: PostStatus, tenantStatus: Status | null): JSX.Element => {
  // Tenant catalogue label wins for custom slugs — i18n has no catalog entry
  // for admin-named statuses, so falling back to <Trans> would render the raw
  // message-id (e.g. "enum.poststatus.parked").
  if (tenantStatus && !["open", "planned", "started", "completed", "declined", "duplicate", "deleted"].includes(tenantStatus.slug)) {
    return <>{tenantStatus.label}</>
  }
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
      return <>{tenantStatus?.label || status.title}</>
  }
}

// Extracts the first non-empty line of the response markdown — that's where
// the Plane webhook writes the substage label (e.g. "In Beta Testing").
const extractSubstage = (text?: string): string | null => {
  if (!text) return null
  const firstLine = text.split("\n").map((l) => l.trim()).find((l) => l.length > 0)
  if (!firstLine) return null
  // Strip simple markdown bold/italic markers so the bubble reads cleanly.
  const cleaned = firstLine.replace(/[*_`]/g, "").trim()
  return cleaned.length > 60 ? cleaned.slice(0, 57) + "..." : cleaned
}

export const ResponseLozenge = (props: PostResponseProps): JSX.Element | null => {
  const status = PostStatus.Get(props.status)
  const tenantStatus = resolveStatus(Fider.session.tenant, props.status)
  const { icon, bg, color, border } = getLozengeProps(status, tenantStatus)
  const translatedStatus = getStatusTranslation(status, tenantStatus)
  const substage =
    props.size === "small" || props.size === "xsmall" || !props.size ? extractSubstage(props.response?.text || undefined) : null

  if (props.size == "micro") {
    return <span className={`${color} text-sm`}>{translatedStatus}</span>
  }

  if (props.size === "xsmall") {
    return (
      <div>
        <HStack align="center" className={`${color} ${bg} rounded-full p-0 px-3`}>
          <Icon sprite={icon} className={`h-4 c-status-col--${status.value}`} />
          <span className={`c-status-col--${status.value} text-xs uppercase`}>
            {translatedStatus}
            {substage && <span className="c-status-substage"> · {substage}</span>}
          </span>
        </HStack>
      </div>
    )
  }

  return (
    <div>
      <HStack align="start" className={`${color} ${bg} border ${border} rounded-full p-1 px-3`}>
        {!props.size && <Icon sprite={icon} className={`h-5 c-status-col--${status.value}`} />}
        <span className={`c-status-col--${status.value} ${props.size === "small" ? "text-sm" : "text-semibold"}`}>
          {translatedStatus}
          {substage && <span className="c-status-substage"> · {substage}</span>}
        </span>
      </HStack>
    </div>
  )
}
