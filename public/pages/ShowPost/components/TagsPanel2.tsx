import React, { useState } from "react"
import { Tag } from "@fider/models"
import { actions } from "@fider/services"
import { ShowTag, Icon } from "@fider/components"
import { TagListItem } from "./TagListItem"
import { useFider } from "@fider/hooks"

import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"
import IconPlusCircle from "@fider/assets/images/heroicons-pluscircle.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/macro"
import { TagsPanelProps } from "./TagsPanel"

export const TagsPanel2 = (props: TagsPanelProps) => {
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

  const tagsList = assignedTags.length > 0 && (
    <VStack spacing={2} className="flex-items-baseline">
      {assignedTags.map((tag) => (
        <ShowTag key={tag.id} tag={tag} link />
      ))}
    </VStack>
  )
  const editTagsList = props.tags.length > 0 && (
    <VStack>
      {props.tags.map((tag) => (
        <TagListItem key={tag.id} tag={tag} assigned={assignedTags.indexOf(tag) >= 0} onClick={assignOrUnassignTag} />
      ))}
    </VStack>
  )

  const icon = canEdit && (isEditing ? <Icon sprite={IconCheckCircle} className="h-4" /> : <Icon sprite={IconPlusCircle} className="h-5 text-primary" />)

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
      <HStack spacing={2} center={true} className="text-primary-base text-xs">
        {!isEditing && tagsList}
        {isEditing && editTagsList}
        <HStack spacing={1} className="clickable" onClick={onSubtitleClick}>
          {icon}
          <span>Edit Tags</span>
        </HStack>
      </HStack>
    </VStack>
  )
}
