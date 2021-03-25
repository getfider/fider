import React, { useState } from "react"
import { Tag, Post } from "@fider/models"
import { actions } from "@fider/services"
import { ShowTag, List, ListItem } from "@fider/components"
import { TagListItem } from "./TagListItem"
import { FaCheckCircle, FaCog } from "react-icons/fa"
import { useFider } from "@fider/hooks"

interface TagsPanelProps {
  post: Post
  tags: Tag[]
}

export const TagsPanel = (props: TagsPanelProps) => {
  const fider = useFider()
  const canEdit = fider.session.isAuthenticated && fider.session.user.isCollaborator && props.tags.length > 0

  const [isEditing, setIsEditing] = useState(false)
  const [assignedTags, setAssignedTags] = useState(props.tags.filter((t) => props.post.tags.indexOf(t.slug) >= 0))

  const assignOrUnassignTag = async (tag: Tag) => {
    const idx = assignedTags.indexOf(tag)
    let nextAssignedTags: Tag[] = []

    if (idx >= 0) {
      const response = await actions.unassignTag(tag.slug, props.post.number)
      if (response.ok) {
        nextAssignedTags = [...assignedTags]
        nextAssignedTags.splice(idx, 1)
      }
    } else {
      const response = await actions.assignTag(tag.slug, props.post.number)
      if (response.ok) {
        nextAssignedTags = [...assignedTags, tag]
      }
    }

    setAssignedTags(nextAssignedTags)
  }

  const onSubtitleClick = () => {
    if (canEdit) {
      setIsEditing(!isEditing)
    }
  }

  if (!canEdit && assignedTags.length === 0) {
    return null
  }

  const tagsList =
    assignedTags.length > 0 ? (
      <List className="c-tag-list">
        {assignedTags.map((tag) => (
          <ListItem key={tag.id}>
            <ShowTag tag={tag} />
          </ListItem>
        ))}
      </List>
    ) : (
      <span className="info">None yet</span>
    )

  const editTagsList = props.tags.length > 0 && (
    <List className="c-tag-list">
      {props.tags.map((tag) => (
        <TagListItem key={tag.id} tag={tag} assigned={assignedTags.indexOf(tag) >= 0} onClick={assignOrUnassignTag} />
      ))}
    </List>
  )

  const subtitleClasses = `subtitle ${canEdit && "active"}`
  const icon = canEdit && (isEditing ? <FaCheckCircle /> : <FaCog />)

  return (
    <>
      <span className={subtitleClasses} onClick={onSubtitleClick}>
        Tags {icon}
      </span>

      {!isEditing && tagsList}
      {isEditing && editTagsList}
    </>
  )
}
