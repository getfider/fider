import React, { useState } from "react"
import { Tag } from "@fider/models"
import { ShowTag, Button, Icon } from "@fider/components"
import { TagFormState, TagForm } from "./TagForm"
import { actions, Failure } from "@fider/services"
import { useFider } from "@fider/hooks"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconPencilAlt from "@fider/assets/images/heroicons-pencil-alt.svg"
import { HStack, VStack } from "@fider/components/layout"

interface TagListItemProps {
  tag: Tag
  onTagEdited: (tag: Tag) => void
  onTagDeleted: (tag: Tag) => void
}

export const TagListItem = (props: TagListItemProps) => {
  const fider = useFider()
  const [tag] = useState(props.tag)
  const [state, setState] = useState<"view" | "edit" | "delete">("view")

  const startDelete = async () => setState("delete")
  const startEdit = async () => setState("edit")
  const resetState = async () => setState("view")

  const deleteTag = async () => {
    const result = await actions.deleteTag(tag.slug)
    if (result.ok) {
      resetState()
      props.onTagDeleted(tag)
    }
  }

  const updateTag = async (data: TagFormState): Promise<Failure | undefined> => {
    const result = await actions.updateTag(tag.slug, data.name, data.color, data.isPublic)
    if (result.ok) {
      tag.name = result.data.name
      tag.slug = result.data.slug
      tag.color = result.data.color
      tag.isPublic = result.data.isPublic

      resetState()
      props.onTagEdited(tag)
    } else {
      return result.error
    }
  }

  const renderDeleteMode = () => {
    return (
      <VStack spacing={2}>
        <div>
          <b>Are you sure?</b>{" "}
          <span>
            The tag <ShowTag tag={tag} /> will be removed from all posts.
          </span>
        </div>
        <div>
          <Button variant="danger" onClick={deleteTag}>
            Delete tag
          </Button>
          <Button onClick={resetState} variant="tertiary">
            Cancel
          </Button>
        </div>
      </VStack>
    )
  }

  const renderViewMode = () => {
    const buttons = fider.session.user.isAdministrator && [
      <Button size="small" key={0} onClick={startEdit}>
        <Icon sprite={IconPencilAlt} />
        <span>Edit</span>
      </Button>,
      <Button size="small" key={1} onClick={startDelete}>
        <Icon sprite={IconX} />
        <span>Delete</span>
      </Button>,
    ]

    return (
      <HStack justify="between">
        <ShowTag tag={tag} link />
        <HStack>{buttons}</HStack>
      </HStack>
    )
  }

  const renderEditMode = () => {
    return <TagForm name={props.tag.name} color={props.tag.color} isPublic={props.tag.isPublic} onSave={updateTag} onCancel={resetState} />
  }

  return state === "delete" ? renderDeleteMode() : state === "edit" ? renderEditMode() : renderViewMode()
}
