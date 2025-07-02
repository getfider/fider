import React, { useEffect, useRef, useState } from "react"
import { Comment, Post } from "@fider/models"
import { Reactions, Avatar, UserName, Moment, Form, Button, Markdown, Modal, Dropdown, Icon } from "@fider/components"
import { HStack } from "@fider/components/layout"
import { formatDate, Failure, actions, notify, copyToClipboard, classSet, clearUrlHash } from "@fider/services"
import { useFider } from "@fider/hooks"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import { t } from "@lingui/core/macro"
import { Trans } from "@lingui/react/macro"
import CommentEditor from "@fider/components/common/form/CommentEditor"
import { useAttachments } from "@fider/hooks/useAttachments"

import "./ShowComment.scss"

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
  const { attachments, handleImageUploaded, getImageSrc } = useAttachments({
    maxAttachments: 2,
  })
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

  const handleApproveComment = async () => {
    const result = await actions.approveComment(props.comment.id)
    if (result.ok) {
      notify.success(<Trans id="showpost.moderation.comment.approved">Comment approved successfully</Trans>)
      setTimeout(() => location.reload(), 1500)
    } else {
      notify.error(<Trans id="showpost.moderation.comment.approveerror">Failed to approve comment</Trans>)
    }
  }

  const handleDeclineComment = async () => {
    const result = await actions.declineComment(props.comment.id)
    if (result.ok) {
      notify.success(<Trans id="showpost.moderation.comment.declined">Comment declined successfully</Trans>)
      setTimeout(() => location.reload(), 1500)
    } else {
      notify.error(<Trans id="showpost.moderation.comment.declineerror">Failed to decline comment</Trans>)
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
    "flex-grow rounded-md p-2 comment-area": true,
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
                  maxAttachments={2}
                  maxImageSizeKB={5 * 1024}
                  onGetImageSrc={getImageSrc}
                  onImageUploaded={handleImageUploaded}
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
                
                {/* Moderation status banner for unapproved comments */}
                {fider.session.tenant.isModerationEnabled && !comment.isApproved && (
                  <div className="mt-3">
                    {fider.session.isAuthenticated && fider.session.user.id === comment.user.id && (
                      <div className="text-muted text-xs p-2 bg-yellow-50 rounded border-l-4 border-yellow-500">
                        <Trans id="showpost.moderation.comment.awaiting">
                          This comment is awaiting moderation by an administrator before being visible to other users.
                        </Trans>
                      </div>
                    )}
                    
                    {/* Admin moderation buttons */}
                    {fider.session.isAuthenticated && fider.session.user.isCollaborator && (
                      <div className="p-2 bg-blue-50 rounded border-l-4 border-blue-500">
                        <div className="mb-1 text-xs font-medium text-blue-800">
                          <Trans id="showpost.moderation.comment.admin.title">Comment Moderation</Trans>
                        </div>
                        <div className="text-xs text-blue-700 mb-2">
                          <Trans id="showpost.moderation.comment.admin.description">
                            This comment is awaiting your approval to be visible to all users.
                          </Trans>
                        </div>
                        <HStack spacing={1}>
                          <Button variant="primary" size="small" onClick={handleApproveComment}>
                            <Trans id="action.approve">Approve</Trans>
                          </Button>
                          <Button variant="danger" size="small" onClick={handleDeclineComment}>
                            <Trans id="action.decline">Decline</Trans>
                          </Button>
                        </HStack>
                      </div>
                    )}
                  </div>
                )}
                
                <Reactions reactions={localReactionCounts} emojiSelectorRef={emojiSelectorRef} toggleReaction={toggleReaction} />
              </>
            )}
          </div>
        </div>
      </HStack>
    </div>
  )
}
