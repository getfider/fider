import React, { useState } from "react"
import { Comment, Post, ImageUpload } from "@fider/models"
import {
  Avatar,
  UserName,
  Moment,
  Form,
  TextArea,
  Button,
  MultiLineText,
  DropDown,
  DropDownItem,
  Modal,
  ImageViewer,
  MultiImageUploader,
} from "@fider/components"
import { formatDate, Failure, actions } from "@fider/services"
import { FaEllipsisH } from "react-icons/fa"
import { useFider } from "@fider/hooks"

interface ShowCommentProps {
  post: Post
  comment: Comment
}

export const ShowComment = (props: ShowCommentProps) => {
  const fider = useFider()
  const [isEditing, setIsEditing] = useState(false)
  const [newContent, setNewContent] = useState("")
  const [isDeleteConfirmationModalOpen, setIsDeleteConfirmationModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure>()

  const canEditComment = (): boolean => {
    if (fider.session.isAuthenticated) {
      return fider.session.user.isCollaborator || props.comment.user.id === fider.session.user.id
    }
    return false
  }

  const clearError = () => setError(undefined)

  const cancelEdit = async () => {
    setIsEditing(false)
    setNewContent("")
    clearError()
  }

  const saveEdit = async () => {
    const response = await actions.updateComment(props.post.number, props.comment.id, newContent, attachments)
    if (response.ok) {
      location.reload()
    } else {
      setError(response.error)
    }
  }

  const renderEllipsis = () => {
    return <FaEllipsisH />
  }

  const closeModal = async () => {
    setIsDeleteConfirmationModalOpen(false)
  }

  const deleteComment = async () => {
    const response = await actions.deleteComment(props.post.number, props.comment.id)
    if (response.ok) {
      location.reload()
    }
  }

  const onActionSelected = (item: DropDownItem) => {
    if (item.value === "edit") {
      setIsEditing(true)
      setNewContent(props.comment.content)
      clearError()
    } else if (item.value === "delete") {
      setIsDeleteConfirmationModalOpen(true)
    }
  }

  const modal = () => {
    return (
      <Modal.Window isOpen={isDeleteConfirmationModalOpen} onClose={closeModal} center={false} size="small">
        <Modal.Header>Delete Comment</Modal.Header>
        <Modal.Content>
          <p>
            This process is irreversible. <strong>Are you sure?</strong>
          </p>
        </Modal.Content>

        <Modal.Footer>
          <Button color="danger" onClick={deleteComment}>
            Delete
          </Button>
          <Button color="cancel" onClick={closeModal}>
            Cancel
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }

  const comment = props.comment

  const editedMetadata = !!comment.editedAt && !!comment.editedBy && (
    <div className="c-comment-metadata">
      <span title={`This comment has been edited by ${comment.editedBy.name} on ${formatDate(comment.editedAt)}`}>· edited</span>
    </div>
  )

  return (
    <div className="c-comment">
      {modal()}
      <Avatar user={comment.user} />
      <div className="c-comment-content">
        <UserName user={comment.user} />
        <div className="c-comment-metadata">
          · <Moment date={comment.createdAt} />
        </div>
        {editedMetadata}
        {!isEditing && canEditComment() && (
          <DropDown
            className="l-more-actions"
            direction="left"
            inline={true}
            style="simple"
            highlightSelected={false}
            items={[
              { label: "Edit", value: "edit" },
              { label: "Delete", value: "delete", render: <span style={{ color: "red" }}>Delete</span> },
            ]}
            onChange={onActionSelected}
            renderControl={renderEllipsis}
          />
        )}
        <div className="c-comment-text">
          {isEditing ? (
            <Form error={error}>
              <TextArea field="content" minRows={1} value={newContent} placeholder={comment.content} onChange={setNewContent} />
              <MultiImageUploader field="attachments" bkeys={comment.attachments} maxUploads={2} previewMaxWidth={100} onChange={setAttachments} />
              <Button size="tiny" onClick={saveEdit} color="positive">
                Save
              </Button>
              <Button color="cancel" size="tiny" onClick={cancelEdit}>
                Cancel
              </Button>
            </Form>
          ) : (
            <>
              <MultiLineText text={comment.content} style="full" />
              {comment.attachments && comment.attachments.map((x) => <ImageViewer key={x} bkey={x} />)}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
