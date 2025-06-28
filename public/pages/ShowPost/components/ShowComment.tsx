import React, { useEffect, useRef, useState } from "react"
import { Comment, Post, ImageUpload } from "@fider/models"
import { Reactions, Avatar, UserName, Moment, Form, Button, Markdown, Modal, Dropdown, Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"
import { formatDate, Failure, actions, notify, copyToClipboard, classSet, clearUrlHash } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import { t } from "@lingui/core/macro"
import { Trans } from "@lingui/react/macro"
import CommentEditor from "@fider/components/common/form/CommentEditor"
import { extractImageBkeys } from "@fider/services/bkey"

interface ShowCommentProps {
  post: Post
  comment: Comment
  highlighted?: boolean
  onToggleReaction?: () => void
}

export const ShowComment = (props: ShowCommentProps) => {
  const fider = useFider()
  const node = useRef<HTMLDivElement | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [newContent, setNewContent] = useState<string>(props.comment.content)
  const [isDeleteConfirmationModalOpen, setIsDeleteConfirmationModalOpen] = useState(false)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [localReactionCounts, setLocalReactionCounts] = useState(props.comment.reactionCounts)
  const emojiSelectorRef = useRef<HTMLDivElement>(null)

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

  useEffect(() => {
    if (isEditing) {
      const bkeys = extractImageBkeys(props.comment.content)
      if (bkeys.length > 0) {
        // Create ImageUpload objects for each bkey found in the comment
        const existingAttachments = bkeys.map(
          (bkey) =>
            ({
              bkey,
              remove: false,
            } as ImageUpload)
        )

        // Initialize attachments state with existing images
        setAttachments(existingAttachments)
      }
    }
  }, [isEditing])

  const canEditComment = (): boolean => {
    if (fider.session.isAuthenticated) {
      return fider.session.user.isCollaborator || props.comment.user.id === fider.session.user.id
    }
    return false
  }

  const clearError = () => setError(undefined)

  const cancelEdit = async () => {
    setIsEditing(false)
    setNewContent(props.comment.content)
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

  const toggleReaction = async (emoji: string) => {
    const response = await actions.toggleCommentReaction(props.post.number, comment.id, emoji)
    if (response.ok) {
      const added = response.data.added

      setLocalReactionCounts((prevCounts) => {
        const newCounts = [...(prevCounts ?? [])]
        const reactionIndex = newCounts.findIndex((r) => r.emoji === emoji)
        if (reactionIndex !== -1) {
          const newCount = added ? newCounts[reactionIndex].count + 1 : newCounts[reactionIndex].count - 1
          if (newCount === 0) {
            newCounts.splice(reactionIndex, 1)
          } else {
            newCounts[reactionIndex] = {
              ...newCounts[reactionIndex],
              count: newCount,
              includesMe: added,
            }
          }
        } else if (added) {
          newCounts.push({ emoji, count: 1, includesMe: true })
        }
        return newCounts
      })
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
    "bg-gray-100": !props.highlighted,
    "bg-gray-200": props.highlighted,
  })

  return (
    <div id={`comment-${comment.id}`}>
      <HStack spacing={2} className="c-comment flex-items-baseline">
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
                <CommentEditor
                  field="content"
                  disabled={!fider.session.isAuthenticated}
                  initialValue={newContent}
                  onChange={setNewContent}
                  placeholder={comment.content}
                  onImageUploaded={(upload) => {
                    // Handle image uploads and removals
                    setAttachments((prev) => {
                      // If this is a removal request, find and mark the attachment for removal
                      if (upload.remove && upload.bkey) {
                        return prev.map((att) => (att.bkey === upload.bkey ? { ...att, remove: true } : att))
                      }
                      // Otherwise add the new upload
                      return [...prev, upload]
                    })
                  }}
                />
                <div className="mt-2">
                  <Button size="small" onClick={saveEdit} variant="primary">
                    <Trans id="action.save">Save</Trans>
                  </Button>
                  <Button variant="tertiary" size="small" onClick={cancelEdit}>
                    <Trans id="action.cancel">Cancel</Trans>
                  </Button>
                </div>
              </Form>
            ) : (
              <>
                <Markdown text={comment.content} style="full" />
                <Reactions reactions={localReactionCounts} emojiSelectorRef={emojiSelectorRef} toggleReaction={toggleReaction} />
              </>
            )}
          </div>
        </div>
      </HStack>
    </div>
  )
}
