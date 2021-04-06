import "./ShowTag.scss"

import React from "react"
import { Tag } from "@fider/models"
import { classSet } from "@fider/services"

interface TagProps {
  tag: Tag
  circular?: boolean
}

const getRGB = (color: string) => {
  const r = color.substring(0, 2)
  const g = color.substring(2, 4)
  const b = color.substring(4, 6)

  return {
    R: parseInt(r, 16),
    G: parseInt(g, 16),
    B: parseInt(b, 16),
  }
}

const textColor = (color: string) => {
  const components = getRGB(color)
  const bgDelta = components.R * 0.299 + components.G * 0.587 + components.B * 0.114
  return bgDelta > 140 ? "#333" : "#fff"
}

export const ShowTag = (props: TagProps) => {
  const className = classSet({
    "c-tag": true,
    "m-circular": props.circular === true,
  })

  return (
    <div
      title={`${props.tag.name}${!props.tag.isPublic ? " (Private)" : ""}`}
      className={className}
      style={{
        backgroundColor: `#${props.tag.color}`,
        color: textColor(props.tag.color),
      }}
    >
      {!props.tag.isPublic && !props.circular && (
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
          <path
            fillRule="evenodd"
            d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
            clipRule="evenodd"
          />
        </svg>
      )}
      {props.circular ? "" : props.tag.name || "Tag"}
    </div>
  )
}
