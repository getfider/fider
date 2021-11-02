import React, { useEffect, useRef, useState } from "react"
import { Comment, Post, ImageUpload } from "@fider/models"
import { Avatar, UserName, Moment, Form, TextArea, Button, Markdown, Modal, ImageViewer, MultiImageUploader, Dropdown, Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"
import { formatDate, Failure, actions, notify, copyToClipboard, classSet, clearUrlHash } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import { t, Trans } from "@lingui/macro"

interface ShowCommentProps {
  post: Post
  comment: Comment
  highlighted?: boolean
}

export const ShowComment = (props: ShowCommentProps) => {
  const fider = useFider()
  const node = useRef<HTMLDivElement | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [newContent, setNewContent] = useState("")
  const [isDeleteConfirmationModalOpen, setIsDeleteConfirmationModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [error, setError] = useState<Failure>()

  const handleClick = (e: MouseEvent) => {
    if (node.current == null || !node.current.contains(e.target as Node)) {
      clearUrlHash()
    }
  }

  useEffect(() => {
    if (props.highlighted) {
      document.addEventListener("mousedown", handleClick)
      return () => {
        document.removeEventListener("mousedown", handleClick)
      }
    }
  }, [props.highlighted])

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

  const closeModal = async () => {
    setIsDeleteConfirmationModalOpen(false)
  }

  const deleteComment = async () => {
    const response = await actions.deleteComment(props.post.number, props.comment.id)
    if (response.ok) {
      location.reload()
    }
  }

  const onActionSelected = (action: string) => () => {
    if (action === "copylink") {
      window.location.hash = `#comment-${props.comment.id}`
      copyToClipboard(window.location.href).then(
        () => notify.success(t({ id: "showpost.comment.copylink.success", message: "Successfully copied comment link to clipboard" })),
        () => notify.error(t({ id: "showpost.comment.copylink.error", message: "Could not copy comment link, please copy page URL" }))
      )
    } else if (action === "edit") {
      setIsEditing(true)
      setNewContent(props.comment.content)
      clearError()
    } else if (action === "delete") {
      setIsDeleteConfirmationModalOpen(true)
    }
  }

  const modal = () => {
    return (
      <Modal.Window isOpen={isDeleteConfirmationModalOpen} onClose={closeModal} center={false} size="small">
        <Modal.Header>
          <Trans id="modal.deletecomment.header">Delete Comment</Trans>
        </Modal.Header>
        <Modal.Content>
          <p>
            <Trans id="modal.deletecomment.text">
              This process is irreversible. <strong>Are you sure?</strong>
            </Trans>
          </p>
        </Modal.Content>

        <Modal.Footer>
          <Button variant="danger" onClick={deleteComment}>
            <Trans id="action.delete">Delete</Trans>
          </Button>
          <Button variant="tertiary" onClick={closeModal}>
            <Trans id="action.cancel">Cancel</Trans>
          </Button>
        </Modal.Footer>
      </Modal.Window>
    )
  }

  const comment = props.comment

  const editedMetadata = !!comment.editedAt && !!comment.editedBy && (
    <span data-tooltip={`This comment has been edited by ${comment.editedBy.name} on ${formatDate(fider.currentLocale, comment.editedAt)}`}>· edited</span>
  )

  const classList = classSet({
    "flex-grow rounded-md p-2": true,
    "bg-gray-50": !props.highlighted,
    "bg-gray-200": props.highlighted,
  })

  return (
    <div id={`comment-${comment.id}`}>
      <HStack spacing={2} center={false} className="c-comment flex-items-baseline">
        {modal()}
        <div className="pt-4">
          <Avatar user={comment.user} />
        </div>
        <div ref={node} className={classList}>
          <div className="mb-1">
            <HStack justify="between">
              <HStack>
                <UserName user={comment.user} />{" "}
                <div className="text-xs">
                  · <Moment locale={fider.currentLocale} date={comment.createdAt} /> {editedMetadata}
                </div>
              </HStack>
              {!isEditing && (
                <Dropdown position="left" renderHandle={<Icon sprite={IconDotsHorizontal} width="16" height="16" />}>
                  <Dropdown.ListItem onClick={onActionSelected("copylink")}>
                    <Trans id="action.copylink">Copy link</Trans>
                  </Dropdown.ListItem>
                  {canEditComment() && (
                    <>
                      <Dropdown.Divider />
                      <Dropdown.ListItem onClick={onActionSelected("edit")}>
                        <Trans id="action.edit">Edit</Trans>
                      </Dropdown.ListItem>
                      <Dropdown.ListItem onClick={onActionSelected("delete")} className="text-red-700">
                        <Trans id="action.delete">Delete</Trans>
                      </Dropdown.ListItem>
                    </>
                  )}
                </Dropdown>
              )}
            </HStack>
          </div>
          <div>
            {isEditing ? (
              <Form error={error}>
                <TextArea field="content" minRows={1} value={newContent} placeholder={comment.content} onChange={setNewContent} />
                <MultiImageUploader field="attachments" bkeys={comment.attachments} maxUploads={2} onChange={setAttachments} />
                <Button size="small" onClick={saveEdit} variant="primary">
                  <Trans id="action.save">Save</Trans>
                </Button>
                <Button variant="tertiary" size="small" onClick={cancelEdit}>
                  <Trans id="action.cancel">Cancel</Trans>
                </Button>
              </Form>
            ) : (
              <>
                <Markdown text={comment.content} style="full" />
                {comment.attachments && comment.attachments.map((x) => <ImageViewer key={x} bkey={x} />)}
              </>
            )}
          </div>
        </div>
      </HStack>
    </div>
  )
}
