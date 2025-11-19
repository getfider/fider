import "./ShowPost.page.scss"

import React, { useState, useEffect, useCallback } from "react"

import { Comment, Post, Tag, Vote, CurrentUser, PostStatus } from "@fider/models"
import { actions, cache, clearUrlHash, Failure, Fider, notify, timeAgo } from "@fider/services"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import IconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import { i18n } from "@lingui/core"
import IconRSS from "@fider/assets/images/heroicons-rss.svg"
import IconPencil from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconChat from "@fider/assets/images/heroicons-chat-alt-2.svg"

import { ResponseDetails, Button, UserName, Moment, Markdown, Input, Form, Icon, Header, PoweredByFider, Avatar, Dropdown, RSSModal } from "@fider/components"
import { DiscussionPanel } from "./components/DiscussionPanel"
import CommentEditor from "@fider/components/common/form/CommentEditor"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { FollowButton } from "./components/FollowButton"
import { VoteSection } from "./components/VoteSection"
import { DeletePostModal } from "./components/DeletePostModal"
import { ResponseModal } from "./components/ResponseModal"
import { VotesPanel } from "./components/VotesPanel"
import { TagsPanel } from "@fider/pages/ShowPost/components/TagsPanel"
import { t } from "@lingui/macro"
import { useFider } from "@fider/hooks"
import { useAttachments } from "@fider/hooks/useAttachments"

interface ShowPostPageProps {
  post: Post
  subscribed: boolean
  comments: Comment[]
  tags: Tag[]
  votes: Vote[]
  attachments: string[]
}

const oneHour = 3600
const canEditPost = (user: CurrentUser, post: Post) => {
  if (user.isCollaborator) {
    return true
  }

  return user.id === post.user.id && timeAgo(post.createdAt) <= oneHour
}

export default function ShowPostPage(props: ShowPostPageProps) {
  const [editMode, setEditMode] = useState(false)
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [isRSSModalOpen, setIsRSSModalOpen] = useState(false)
  const [showResponseModal, setShowResponseModal] = useState(false)
  const [newTitle, setNewTitle] = useState(props.post.title)
  const [newDescription, setNewDescription] = useState(props.post.description)
  const { attachments, handleImageUploaded, getImageSrc } = useAttachments({
    maxAttachments: 3,
  })
  const [highlightedComment, setHighlightedComment] = useState<number | undefined>(undefined)
  const [error, setError] = useState<Failure | undefined>(undefined)
  const fider = useFider()

  const handleHashChange = useCallback(
    (e?: Event) => {
      const hash = window.location.hash
      const result = /#comment-([0-9]+)/.exec(hash)

      let newHighlightedComment
      if (result === null) {
        // No match
        newHighlightedComment = undefined
      } else {
        // Match, extract numeric ID
        const id = parseInt(result[1])
        if (props.comments.map((comment) => comment.id).includes(id)) {
          newHighlightedComment = id
        } else {
          // Unknown comment
          if (e?.cancelable) {
            e.preventDefault()
          } else {
            clearUrlHash(true)
          }
          notify.error(<Trans id="showpost.comment.unknownhighlighted">Unknown comment ID #{id}</Trans>)
          newHighlightedComment = undefined
        }
      }
      setHighlightedComment(newHighlightedComment)
    },
    [props.comments]
  )

  useEffect(() => {
    handleHashChange()
    window.addEventListener("hashchange", handleHashChange)
    return () => {
      window.removeEventListener("hashchange", handleHashChange)
    }
  }, [handleHashChange])

  useEffect(() => {
    const showSuccess = cache.session.get("POST_CREATED_SUCCESS")
    const showModeration = cache.session.get("POST_CREATED_MODERATION")
    const showCommentModeration = cache.session.get("COMMENT_CREATED_MODERATION")

    if (showSuccess) {
      cache.session.remove("POST_CREATED_SUCCESS")
      notify.success(t({ id: "mysettings.notification.event.newpostcreated", message: "Your idea has been added 👍" }))
    }

    if (showModeration) {
      cache.session.remove("POST_CREATED_MODERATION")
      notify.success(t({ id: "showpost.moderation.postsuccess", message: "Your idea has been submitted and is awaiting moderation 📝" }))
    }

    if (showCommentModeration) {
      cache.session.remove("COMMENT_CREATED_MODERATION")
      notify.success(t({ id: "showpost.moderation.commentsuccess", message: "Your comment has been submitted and is awaiting moderation 📝" }))
    }
  }, [])

  const saveChanges = async () => {
    const result = await actions.updatePost(props.post.number, newTitle, newDescription, attachments)
    if (result.ok) {
      location.reload()
    } else {
      setError(result.error)
    }
  }

  const canDeletePost = () => {
    const status = PostStatus.Get(props.post.status)
    if (!Fider.session.isAuthenticated || !Fider.session.user.isAdministrator || status.closed) {
      return false
    }
    return true
  }

  const cancelEdit = () => {
    setError(undefined)
    setEditMode(false)
  }

  const startEdit = () => {
    setEditMode(true)
  }

  const handleDescriptionChange = (value: string) => {
    setNewDescription(value)
  }

  const handleApprovePost = async () => {
    const result = await actions.approvePost(props.post.id)
    if (result.ok) {
      notify.success(<Trans id="showpost.moderation.approved">Post approved successfully</Trans>)
      setTimeout(() => location.reload(), 1500)
    } else {
      notify.error(<Trans id="showpost.moderation.approveerror">Failed to approve post</Trans>)
    }
  }

  const handleDeclinePost = async () => {
    const result = await actions.declinePost(props.post.id)
    if (result.ok) {
      notify.success(<Trans id="showpost.moderation.declined">Post declined successfully</Trans>)
      setTimeout(() => location.reload(), 1500)
    } else {
      notify.error(<Trans id="showpost.moderation.declineerror">Failed to decline post</Trans>)
    }
  }

  const onActionSelected = (action: "copy" | "delete" | "status" | "feed" | "edit") => () => {
    if (action === "copy") {
      navigator.clipboard.writeText(window.location.href)
      notify.success(<Trans id="showpost.copylink.success">Link copied to clipboard</Trans>)
    } else if (action === "delete") {
      setShowDeleteModal(true)
    } else if (action === "status") {
      setShowResponseModal(true)
    } else if (action === "edit") {
      startEdit()
    } else if (action == "feed") {
      setIsRSSModalOpen(true)
    }
  }

  const hideRSSModal = () => setIsRSSModalOpen(false)

  return (
    <>
      <Header />
      <div id="p-show-post" className="page container">
        <div className="p-show-post">
          <div className="p-show-post__main-col">
            <div className="p-show-post__header-col">
              <VStack spacing={8}>
                <HStack justify="between">
                  <VStack align="start">
                    {!editMode && (
                      <HStack>
                        <Avatar user={props.post.user} />
                        <VStack spacing={1}>
                          <UserName user={props.post.user} />
                          <Moment className="text-muted" locale={Fider.currentLocale} date={props.post.createdAt} />
                        </VStack>
                      </HStack>
                    )}
                  </VStack>
                  <RSSModal isOpen={isRSSModalOpen} onClose={hideRSSModal} url={`${fider.settings.baseURL}/feed/posts/${props.post.number}.atom`} />

                  {!editMode && (
                    <HStack>
                      <Dropdown position="left" renderHandle={<Icon sprite={IconDotsHorizontal} width="24" height="24" />}>
                        <Dropdown.ListItem onClick={onActionSelected("copy")} icon={IconDuplicate}>
                          <Trans id="action.copylink">Copy link</Trans>
                        </Dropdown.ListItem>
                        {Fider.session.tenant.isFeedEnabled && (
                          <Dropdown.ListItem onClick={onActionSelected("feed")} icon={IconRSS}>
                            <Trans id="action.commentsfeed">Comment Feed</Trans>
                          </Dropdown.ListItem>
                        )}
                        {Fider.session.isAuthenticated && canEditPost(Fider.session.user, props.post) && (
                          <>
                            <Dropdown.ListItem onClick={onActionSelected("edit")} icon={IconPencil}>
                              <Trans id="action.edit">Edit</Trans>
                            </Dropdown.ListItem>
                            {Fider.session.user.isCollaborator && (
                              <Dropdown.ListItem onClick={onActionSelected("status")} icon={IconChat}>
                                <Trans id="action.respond">Respond</Trans>
                              </Dropdown.ListItem>
                            )}
                          </>
                        )}
                        {canDeletePost() && (
                          <Dropdown.ListItem onClick={onActionSelected("delete")} className="text-red-700" icon={IconX}>
                            <Trans id="action.delete">Delete</Trans>
                          </Dropdown.ListItem>
                        )}
                      </Dropdown>
                    </HStack>
                  )}
                </HStack>

                <div className="flex-grow">
                  {editMode ? (
                    <Form error={error}>
                      <Input field="title" maxLength={100} value={newTitle} onChange={setNewTitle} />
                    </Form>
                  ) : (
                    <>
                      <h1 className="text-large text-break">{props.post.title}</h1>

                      {/* Moderation status banner for unapproved posts */}
                      {fider.session.showModerationControls && !props.post.isApproved && (
                        <div className="mt-4">
                          {fider.session.isAuthenticated && fider.session.user.id === props.post.user.id && (
                            <div className="text-muted text-sm p-3 bg-yellow-50 rounded border-l-4 border-yellow-500">
                              <Trans id="showpost.moderation.awaiting">Awaiting moderation.</Trans>
                            </div>
                          )}

                          {/* Admin moderation buttons */}
                          {fider.session.isAuthenticated && fider.session.user.isCollaborator && (
                            <div className="p-3 bg-blue-50 rounded border-l-4 border-blue-500">
                              <div className="mb-2 text-sm font-medium text-blue-800">
                                <Trans id="showpost.moderation.admin.title">Moderation</Trans>
                              </div>
                              <div className="text-sm text-blue-700 mb-3">
                                <Trans id="showpost.moderation.admin.description">This idea needs your approval before being published</Trans>
                              </div>
                              <HStack spacing={2}>
                                <Button variant="primary" size="small" onClick={handleApprovePost}>
                                  <Trans id="action.publish">Publish</Trans>
                                </Button>
                                <Button variant="danger" size="small" onClick={handleDeclinePost}>
                                  <Trans id="action.delete">Delete</Trans>
                                </Button>
                              </HStack>
                            </div>
                          )}
                        </div>
                      )}
                    </>
                  )}
                </div>

                <DeletePostModal onModalClose={() => setShowDeleteModal(false)} showModal={showDeleteModal} post={props.post} />
                {Fider.session.isAuthenticated && Fider.session.user.isCollaborator && (
                  <ResponseModal onCloseModal={() => setShowResponseModal(false)} showModal={showResponseModal} post={props.post} />
                )}
                <VStack>
                  {editMode ? (
                    <Form error={error}>
                      <CommentEditor
                        field="description"
                        onChange={handleDescriptionChange}
                        initialValue={newDescription}
                        disabled={false}
                        maxAttachments={3}
                        maxImageSizeKB={5 * 1024}
                        placeholder={i18n._({
                          id: "newpost.modal.description.placeholder",
                          message: "Tell us about it. Explain it fully, don't hold back, the more information the better.",
                        })}
                        onImageUploaded={handleImageUploaded}
                        onGetImageSrc={getImageSrc}
                      />
                    </Form>
                  ) : (
                    <>
                      {props.post.description && <Markdown className="description" text={props.post.description} style="full" />}
                      {!props.post.description && (
                        <em className="text-muted">
                          <Trans id="showpost.message.nodescription">No description provided.</Trans>
                        </em>
                      )}
                    </>
                  )}
                </VStack>
                <div className="mt-2">
                  <TagsPanel post={props.post} tags={props.tags} />
                </div>

                <VStack spacing={4}>
                  {!editMode ? (
                    <HStack justify="between" align="start">
                      <VoteSection post={props.post} votes={props.post.votesCount} />
                      <FollowButton post={props.post} subscribed={props.subscribed} />
                    </HStack>
                  ) : (
                    <HStack>
                      <Button variant="primary" onClick={saveChanges} disabled={Fider.isReadOnly}>
                        <Icon sprite={IconThumbsUp} />{" "}
                        <span>
                          <Trans id="action.save">Save</Trans>
                        </span>
                      </Button>
                      <Button onClick={cancelEdit} disabled={Fider.isReadOnly}>
                        <Icon sprite={IconX} />
                        <span>
                          <Trans id="action.cancel">Cancel</Trans>
                        </span>
                      </Button>
                    </HStack>
                  )}
                  <div className="border-4 border-blue-500" />
                </VStack>

                <ResponseDetails status={props.post.status} response={props.post.response} />
              </VStack>
            </div>

            <div className="p-show-post__discussion_col">
              <DiscussionPanel post={props.post} comments={props.comments} highlightedComment={highlightedComment} />
            </div>
          </div>
          <div className="p-show-post__action-col">
            <VotesPanel post={props.post} votes={props.votes} />
            <PoweredByFider slot="show-post" className="mt-3" />
          </div>
        </div>
      </div>
    </>
  )
}
