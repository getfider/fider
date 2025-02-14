import React, { useState } from "react"
import { Post, Tag } from "@fider/models"
import { actions } from "@fider/services"
import { ShowTag, Button } from "@fider/components"
import { TagListItem } from "./TagListItem"
import { useFider } from "@fider/hooks"

import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

export interface TagsPanelProps {
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

  const tagsList = (
    <HStack spacing={2} align="center">
      {assignedTags.length > 0 && assignedTags.map((tag) => <ShowTag key={tag.id} tag={tag} link />)}
      <HStack spacing={1} align="center" className="clickable" onClick={onSubtitleClick}>
        <span>
          <Trans id="label.edittags">Edit tags</Trans>
        </span>
      </HStack>
    </HStack>
  )

  const editTagsList = props.tags.length > 0 && (
    <VStack justify="between" className="flex-items-start">
      {props.tags.map((tag) => (
        <TagListItem key={tag.id} tag={tag} assigned={assignedTags.indexOf(tag) >= 0} onClick={assignOrUnassignTag} />
      ))}
      <Button variant="secondary" size="small" onClick={onSubtitleClick}>
        <Trans id="action.close">Close</Trans>
      </Button>
    </VStack>
  )

  if (fider.isReadOnly) {
    return (
      <VStack>
        <HStack spacing={2} className="text-category">
          <Trans id="label.tags">Tags</Trans>
        </HStack>
        {tagsList}
      </VStack>
    )
  }

  return (
    <VStack>
      <HStack spacing={2} align="center" className="text-primary-base text-xs">
        {!isEditing && tagsList}
        {isEditing && editTagsList}
      </HStack>
    </VStack>
  )
}
