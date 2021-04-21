import React, { useState } from "react"
import { Tag } from "@fider/models"
import { ListItem, ShowTag, Button } from "@fider/components"
import { TagFormState, TagForm } from "./TagForm"
import { actions, Failure } from "@fider/services"
import { FaTimes, FaEdit } from "react-icons/fa"
import { useFider } from "@fider/hooks"

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
      <>
        <div className="content">
          <b>Are you sure?</b>{" "}
          <span>
            The tag <ShowTag tag={tag} /> will be removed from all posts.
          </span>
        </div>
        <Button className="right" onClick={resetState} color="cancel">
          Cancel
        </Button>
        <Button color="danger" className="right" onClick={deleteTag}>
          Delete tag
        </Button>
      </>
    )
  }

  const renderViewMode = () => {
    const buttons = fider.session.user.isAdministrator && [
      <Button size="mini" key={0} onClick={startDelete} className="right">
        <FaTimes />
        Delete
      </Button>,
      <Button size="mini" key={1} onClick={startEdit} className="right">
        <FaEdit />
        Edit
      </Button>,
    ]

    return (
      <>
        <ShowTag tag={tag} />
        {buttons}
      </>
    )
  }

  const renderEditMode = () => {
    return <TagForm name={props.tag.name} color={props.tag.color} isPublic={props.tag.isPublic} onSave={updateTag} onCancel={resetState} />
  }

  const view = state === "delete" ? renderDeleteMode() : state === "edit" ? renderEditMode() : renderViewMode()

  return <ListItem>{view}</ListItem>
}
