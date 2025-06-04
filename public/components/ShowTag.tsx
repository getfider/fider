import "./ShowTag.scss"

import React, { MouseEvent } from "react"
import { Tag } from "@fider/models"
import { classSet } from "@fider/services"
import ShieldCheck from "@fider/assets/images/heroicons-shieldcheck.svg"
import { Icon } from "./common"

interface TagProps {
  tag: Tag
  circular?: boolean
  link?: boolean
  noBackground?: boolean
}

// const textColor = (color: string) => {
//   const components = getRGB(color)
//   const bgDelta = components.R * 0.299 + components.G * 0.587 + components.B * 0.114
//   return bgDelta > 140 ? "#333" : "#fff"
// }

export const ShowTag = (props: TagProps) => {
  const className = classSet({
    "c-tag": true,
    "c-tag--circular": props.circular === true,
    "c-tag--transparent": props.noBackground === true,
  })

  const clickHandler = (event: MouseEvent) => {
    if (!props.link) {
      event.preventDefault()
    }
  }

  return (
    <a
      // always add an href, so the tag can be selected by keyboard in the TagsSelect
      href={props.link && props.tag.slug ? `/?tags=${props.tag.slug}` : ""}
      title={`${props.tag.name}${props.tag.isPublic ? "" : " (Private)"}`}
      className={className}
      onClick={clickHandler}
    >
      <span
        style={{
          backgroundColor: `#${props.tag.color}`,
        }}
      ></span>
      {!props.tag.isPublic && !props.circular && <Icon height="14" width="14" sprite={ShieldCheck} className="mr-1" />}
      {props.circular ? "" : props.tag.name || "Tag"}
    </a>
  )
}
