import React, { useState } from "react"
import { Post, Tag } from "@fider/models"
import { actions } from "@fider/services"
import { useFider } from "@fider/hooks"
import { TagsSelect } from "@fider/components/common/TagsSelect"

export interface TagsPanelProps {
  onDataChanged?: () => void
  post: Post
  tags: Tag[]
}

export const TagsPanel = (props: TagsPanelProps) => {
  const fider = useFider()
  const canEdit = fider.session.isAuthenticated && fider.session.user.isCollaborator && props.tags.length > 0

  const [assignedTags, setAssignedTags] = useState(props.tags.filter((t) => props.post.tags.indexOf(t.slug) >= 0))

  const assignOrUnassignTag = async (tags: Tag[]) => {
    await Promise.all([
      ...tags.filter((t) => !assignedTags.includes(t)).map((t) => actions.assignTag(t.slug, props.post.number)),
      ...assignedTags.filter((t) => !tags.includes(t)).map((t) => actions.unassignTag(t.slug, props.post.number)),
    ])

    setAssignedTags(tags)
    props.onDataChanged?.()
  }

  return <TagsSelect tags={props.tags} selected={assignedTags} canEdit={canEdit} selectionChanged={assignOrUnassignTag} asLinks />
}
