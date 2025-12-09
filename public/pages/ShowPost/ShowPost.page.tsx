import "./ShowPost.page.scss"

import React, { useState, useEffect, useCallback } from "react"

import { Comment, Post, Tag, Vote, CurrentUser, PostStatus } from "@fider/models"
import { actions, cache, clearUrlHash, Failure, Fider, notify, timeAgo } from "@fider/services"
import IconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import { i18n } from "@lingui/core"
import IconRSS from "@fider/assets/images/heroicons-rss.svg"
import IconPencil from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconChat from "@fider/assets/images/heroicons-chat-alt-2.svg"

import {
  ResponseDetails,
  Button,
  UserName,
  Moment,
  Markdown,
  Input,
  Form,
  Icon,
  Header,
  PoweredByFider,
  Avatar,
  RSSModal,
  ResponseLozenge,
} from "@fider/components"
import { CommentInput } from "./components/CommentInput"
import { ShowComment } from "./components/ShowComment"
import { VoteSection } from "./components/VoteSection"
import CommentEditor from "@fider/components/common/form/CommentEditor"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { StayUpdatedCard } from "./components/StayUpdatedCard"
import { DeletePostModal } from "./components/DeletePostModal"
import { ResponseModal } from "./components/ResponseModal"
import { VotesPanel } from "./components/VotesPanel"
import { TagsPanel } from "@fider/pages/ShowPost/components/TagsPanel"
import { ActionButton } from "./components/ActionButton"
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

const PostMetaInfo = ({ post, locale }: { post: Post; locale: string }) => (
  <HStack spacing={2} align="center">
    <Avatar user={post.user} size="small" />
    <span className="text-sm text-gray-600">
      <Trans id="showpost.postedby">Posted by</Trans> <UserName user={post.user} />
    </span>
    <span className="text-sm text-gray-400">•</span>
    <Moment className="text-sm text-gray-600" locale={locale} date={post.createdAt} />
    <span className="text-sm text-gray-400">•</span>
    <ResponseLozenge status={post.status} response={post.response} size="xsmall" />
  </HStack>
)

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
    if (showSuccess) {
      cache.session.remove("POST_CREATED_SUCCESS")
      // Show success message/toast
      notify.success(t({ id: "mysettings.notification.event.newpostcreated", message: "Your idea has been added 👍" }))
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
          {/* Left Sidebar */}
          <div className="p-show-post__action-col">
            <StayUpdatedCard post={props.post} subscribed={props.subscribed} />

            <VotesPanel post={props.post} votes={props.votes} />

            <PoweredByFider slot="show-post" className="mt-3" />
          </div>

          <div className="p-show-post__main-col">
            {/* Post Card */}
            <div className="p-show-post__post-card">
              {/* Title and Meta */}
              <VStack spacing={4}>
                {/* Title */}
                {editMode ? (
                  <Form error={error}>
                    <Input field="title" maxLength={100} value={newTitle} onChange={setNewTitle} />
                  </Form>
                ) : (
                  <h1 className="p-show-post__title">{props.post.title}</h1>
                )}

                {/* Posted by info with status */}
                {!editMode && (
                  <>
                    <div className="p-show-post__meta">
                      <PostMetaInfo post={props.post} locale={fider.currentLocale} />
                    </div>
                    <ResponseDetails status={props.post.status} response={props.post.response} />
                  </>
                )}
              </VStack>

              {/* Description - Full width */}
              {!editMode ? (
                <div className="p-show-post__description-section">
                  {props.post.description && <Markdown className="p-show-post__description" text={props.post.description} style="full" />}
                  {!props.post.description && (
                    <em className="text-muted">
                      <Trans id="showpost.message.nodescription">No description provided.</Trans>
                    </em>
                  )}
                </div>
              ) : (
                <div className="p-show-post__description-section">
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
                </div>
              )}

              {/* Vote Section */}
              {!editMode && (
                <div className="p-show-post__vote-section">
                  <VoteSection post={props.post} votes={props.post.votesCount} />
                </div>
              )}

              {/* Tags inline */}
              <div className="pt-7">
                <TagsPanel post={props.post} tags={props.tags} />
              </div>

              {/* Edit Mode Actions */}
              {editMode && (
                <HStack className="mt-6">
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

              {/* Bottom Action Bar */}
              {!editMode && (
                <div className="p-show-post__actions">
                  <HStack spacing={0} align="center" className="flex-wrap gap-2">
                    <ActionButton icon={IconDuplicate} onClick={onActionSelected("copy")}>
                      <Trans id="action.copylink">Copy link</Trans>
                    </ActionButton>

                    {Fider.session.isAuthenticated && canEditPost(Fider.session.user, props.post) && (
                      <ActionButton icon={IconPencil} onClick={onActionSelected("edit")}>
                        <Trans id="action.edit">Edit</Trans>
                      </ActionButton>
                    )}

                    {Fider.session.isAuthenticated && Fider.session.user.isCollaborator && (
                      <ActionButton icon={IconChat} onClick={onActionSelected("status")}>
                        <Trans id="action.respond">Respond</Trans>
                      </ActionButton>
                    )}

                    {Fider.session.tenant.isFeedEnabled && (
                      <ActionButton icon={IconRSS} onClick={onActionSelected("feed")}>
                        <Trans id="action.commentsfeed">Comment Feed</Trans>
                      </ActionButton>
                    )}

                    {canDeletePost() && (
                      <ActionButton icon={IconX} onClick={onActionSelected("delete")} variant="danger">
                        <Trans id="action.delete">Delete</Trans>
                      </ActionButton>
                    )}

                    <div className="flex-grow" />
                  </HStack>
                </div>
              )}
            </div>

            {/* Discussion Section */}
            <div className="p-show-post__discussion-section">
              {/* Discussion Header */}
              <HStack className="p-show-post__discussion-header" align="center" spacing={4}>
                <h2 className="text-xl text-bold text-gray-900">
                  <span className="text-bold">
                    <Trans id="label.discussion">Discussion</Trans>
                  </span>
                </h2>
                <div className="p-show-post__discussion-count">{props.comments.length}</div>
              </HStack>

              {/* Comment Input at top */}
              <CommentInput post={props.post} />

              {/* Comments List */}
              {props.comments.length > 0 && (
                <VStack spacing={4}>
                  {props.comments.map((c) => (
                    <ShowComment key={c.id} post={props.post} comment={c} highlighted={highlightedComment === c.id} />
                  ))}
                </VStack>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Modals */}
      <RSSModal isOpen={isRSSModalOpen} onClose={hideRSSModal} url={`${fider.settings.baseURL}/feed/posts/${props.post.number}.atom`} />
      <DeletePostModal onModalClose={() => setShowDeleteModal(false)} showModal={showDeleteModal} post={props.post} />
      {Fider.session.isAuthenticated && Fider.session.user.isCollaborator && (
        <ResponseModal onCloseModal={() => setShowResponseModal(false)} showModal={showResponseModal} post={props.post} />
      )}
    </>
  )
}
