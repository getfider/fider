import React from "react"
import { Tag } from "@fider/models"
import { ListItem, ShowTag } from "@fider/components"
import { FaCheck } from "react-icons/fa"

interface TagListItemProps {
  tag: Tag
  assigned: boolean
  onClick: (tag: Tag) => void
}

export const TagListItem = (props: TagListItemProps) => {
  const onClick = () => {
    props.onClick(props.tag)
  }

  return (
    <ListItem onClick={onClick}>
      {props.assigned ? <FaCheck /> : <svg className="icon" />}
      <ShowTag tag={props.tag} size="mini" circular={true} />
      <span>{props.tag.name}</span>
    </ListItem>
  )
}
