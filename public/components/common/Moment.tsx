import React from "react"
import { formatDate, timeSince } from "@fider/services"

interface MomentText {
  locale: string
  date: Date | string
  format?: "relative" | "full" | "short" | "date"
}

export const Moment = (props: MomentText) => {
  if (!props.date) {
    return <span />
  }

  const format = props.format || "relative"

  const now = new Date()
  const date = props.date instanceof Date ? props.date : new Date(props.date)
  const diff = (now.getTime() - date.getTime()) / (60 * 60 * 24 * 1000)
  const display =
    diff >= 365 && format === "relative"
      ? formatDate(props.locale, props.date, "short")
      : format === "relative"
      ? timeSince(props.locale, now, date)
      : format === "date"
      ? formatDate(props.locale, props.date, "date")
      : formatDate(props.locale, props.date, format)

  const tooltip = props.format === "short" ? formatDate(props.locale, props.date, "full") : undefined

  return (
    <span className="date" data-tooltip={tooltip}>
      {display}
    </span>
  )
}
