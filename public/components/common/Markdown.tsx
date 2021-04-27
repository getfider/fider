import React from "react"
import { markdown, truncate } from "@fider/services"

import "./Markdown.scss"

interface MarkdownProps {
  className?: string
  text?: string
  maxLength?: number
  style: "full" | "plainText"
}

export const Markdown = (props: MarkdownProps) => {
  if (!props.text) {
    return null
  }

  const html = markdown[props.style](props.text)
  const className = `c-markdown ${props.className || ""}`
  const tagName = props.style === "plainText" ? "p" : "div"

  return React.createElement(tagName, {
    className,
    dangerouslySetInnerHTML: { __html: props.maxLength ? truncate(html, props.maxLength) : html },
  })
}
