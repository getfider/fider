import React from "react"
import { markdown, truncate } from "@fider/services"

interface MultiLineTextProps {
  className?: string
  text?: string
  maxLength?: number
  style: "full" | "plainText"
}

export const MultiLineText = (props: MultiLineTextProps) => {
  if (!props.text) {
    return null
  }

  const html = markdown[props.style](props.text)
  const className = `markdown-body ${props.className || ""}`
  const tagName = props.style === "plainText" ? "p" : "div"

  return React.createElement(tagName, {
    className,
    dangerouslySetInnerHTML: { __html: props.maxLength ? truncate(html, props.maxLength) : html },
  })
}
