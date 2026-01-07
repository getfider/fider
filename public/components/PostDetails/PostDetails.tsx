import "./PostDetails.scss"

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
  Avatar,
  PoweredByFider,
  RSSModal,
  ResponseLozenge,
} from "@fider/components"
import { CommentInput } from "@fider/pages/ShowPost/components/CommentInput"
import { ShowComment } from "@fider/pages/ShowPost/components/ShowComment"
import { VoteSection } from "@fider/pages/ShowPost/components/VoteSection"
import CommentEditor from "@fider/components/common/form/CommentEditor"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import IconTrash from "@fider/assets/images/heroicons-trash.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { DeletePostModal } from "@fider/pages/ShowPost/components/DeletePostModal"
import { ResponseModal } from "@fider/pages/ShowPost/components/ResponseModal"
import { VotesPanel } from "@fider/pages/ShowPost/components/VotesPanel"
import { TagsPanel } from "@fider/pages/ShowPost/components/TagsPanel"
import { ActionButton } from "@fider/pages/ShowPost/components/ActionButton"
import { t } from "@lingui/macro"
import { useFider } from "@fider/hooks"
import { useAttachments } from "@fider/hooks/useAttachments"
import { FollowButton } from "@fider/pages/ShowPost/components/FollowButton"

interface PostDetailsProps {
  postNumber: number
  // Optional initial data for SSR
  initialPost?: Post
  initialSubscribed?: boolean
  initialComments?: Comment[]
  initialTags?: Tag[]
  initialVotes?: Vote[]
  initialAttachments?: string[]
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

export const PostDetails: React.FC<PostDetailsProps> = (props) => {
  // If we have initial data, use it; otherwise we'll fetch
  const [post, setPost] = useState<Post | null>(props.initialPost || null)
  const [subscribed, setSubscribed] = useState(props.initialSubscribed || false)
  const [comments, setComments] = useState<Comment[]>(props.initialComments || [])
  const [tags, setTags] = useState<Tag[]>(props.initialTags || [])
  const [votes, setVotes] = useState<Vote[]>(props.initialVotes || [])
  const [loading, setLoading] = useState(!props.initialPost)

  const [editMode, setEditMode] = useState(false)
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [isRSSModalOpen, setIsRSSModalOpen] = useState(false)
  const [showResponseModal, setShowResponseModal] = useState(false)
  const [newTitle, setNewTitle] = useState(post?.title || "")
  const [newDescription, setNewDescription] = useState(post?.description || "")
  const { attachments, handleImageUploaded, getImageSrc } = useAttachments({
    maxAttachments: 3,
  })
  const [highlightedComment, setHighlightedComment] = useState<number | undefined>(undefined)
  const [error, setError] = useState<Failure | undefined>(undefined)
  const fider = useFider()

  // Fetch data if not provided initially
  useEffect(() => {
    if (props.initialPost) {
      return // Already have data from SSR
    }

    const fetchData = async () => {
      setLoading(true)
      try {
        const [postResponse, commentsResult, tagsResult] = await Promise.all([
          fetch(`/api/v1/posts/${props.postNumber}`).then((r) => r.json()),
          fetch(`/api/v1/posts/${props.postNumber}/comments`).then((r) => r.json()),
          fetch(`/api/v1/tags`).then((r) => r.json()),
        ])

        if (postResponse) {
          setPost(postResponse)
          setNewTitle(postResponse.title)
          setNewDescription(postResponse.description)
        }

        setComments(commentsResult || [])
        setTags(tagsResult || [])

        // Fetch votes and subscription status if authenticated
        if (Fider.session.isAuthenticated) {
          try {
            const votesResult = await actions.listVotes(props.postNumber)
            // Limit votes to 24 to match SSR behavior
            setVotes(votesResult.ok ? votesResult.data.slice(0, 24) : [])

            const subResponse = await fetch(`/api/v1/posts/${props.postNumber}/subscription`)
            if (subResponse.ok) {
              const subData = await subResponse.json()
              setSubscribed(subData.subscribed || false)
            }
          } catch (e) {
            // Ignore subscription and vote errors
          }
        }
      } catch (err) {
        console.error("Failed to fetch post data:", err)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [props.postNumber, props.initialPost])

  const handleHashChange = useCallback(
    (e?: Event) => {
      const hash = window.location.hash
      const result = /#comment-([0-9]+)/.exec(hash)

      let newHighlightedComment
      if (result === null) {
        newHighlightedComment = undefined
      } else {
        const id = parseInt(result[1])
        if (comments.map((comment) => comment.id).includes(id)) {
          newHighlightedComment = id
        } else {
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
    [comments]
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
    if (!post) return
    const result = await actions.updatePost(post.number, newTitle, newDescription, attachments)
    if (result.ok) {
      setEditMode(false)
      // Update local state with new values
      setPost({
        ...post,
        title: newTitle,
        description: newDescription,
      })
      notify.success(<Trans id="showpost.save.success">Post updated successfully</Trans>)
    } else {
      setError(result.error)
    }
  }

  const canDeletePost = () => {
    if (!post) return false
    const status = PostStatus.Get(post.status)
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
    if (!post) return
    const result = await actions.approvePost(post.id)
    if (result.ok) {
      notify.success(<Trans id="showpost.moderation.approved">Post approved successfully</Trans>)
      setTimeout(() => location.reload(), 1500)
    } else {
      notify.error(<Trans id="showpost.moderation.approveerror">Failed to approve post</Trans>)
    }
  }

  const handleDeclinePost = async () => {
    if (!post) return
    const result = await actions.declinePost(post.id)
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

  if (loading || !post) {
    return (
      <div className="p-show-post">
        <div className="p-show-post__loading">
          <Trans id="showpost.loading">Loading...</Trans>
        </div>
      </div>
    )
  }

  return (
    <div className="p-show-post">
      {/* Left Sidebar - hidden on mobile, shown on desktop */}
      <div className="p-show-post__action-col p-show-post__action-col--desktop">
        <VotesPanel post={post} votes={votes} />

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
              <h1 className="p-show-post__title">{post.title}</h1>
            )}

            {/* Posted by info with status */}
            {!editMode && (
              <div className="p-show-post__meta">
                <PostMetaInfo post={post} locale={fider.currentLocale} />
              </div>
            )}
          </VStack>

          {/* Moderation status banner for unapproved posts */}
          {!editMode && !post.isApproved && (
            <div>
              {fider.session.isAuthenticated && (
                <div className="text-muted text-sm p-3 bg-yellow-50 rounded-md mt-2 border-yellow-500">
                  <Trans id="showpost.moderation.awaiting">Awaiting moderation.</Trans>
                </div>
              )}

              {/* Admin moderation buttons */}
              {fider.session.isAuthenticated && fider.session.showModerationControls && fider.session.user.isCollaborator && (
                <div className="p-3 bg-blue-50 rounded border-l-4 border-blue-500 mt-4">
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

          {/* Description - Full width */}
          {!editMode ? (
            <div className="p-show-post__description-section">
              {post.description && <Markdown className="p-show-post__description" text={post.description} style="full" />}
              {!post.description && (
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

          {tags.length >= 1 && (
            <div className="pt-7">
              <TagsPanel post={post} tags={tags} />
            </div>
          )}

          {/* Vote Section */}
          {!editMode && (
            <div className="p-show-post__vote-section">
              <VoteSection post={post} votes={post.votesCount} />
            </div>
          )}

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

                {Fider.session.isAuthenticated && canEditPost(Fider.session.user, post) && (
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
                  <ActionButton icon={IconTrash} onClick={onActionSelected("delete")} variant="danger">
                    <Trans id="action.delete">Delete</Trans>
                  </ActionButton>
                )}

                <div className="flex-grow" />

                <FollowButton post={post} subscribed={subscribed} />
              </HStack>
            </div>
          )}
        </div>

        {/* Mobile Sidebar - shown after post card on mobile */}
        <div className="p-show-post__action-col p-show-post__action-col--mobile">
          <VotesPanel post={post} votes={votes} />
        </div>

        {/* Discussion Section */}
        <div className="p-show-post__discussion-section">
          {/* Discussion Header */}
          <HStack className="p-show-post__discussion-header mt-7" align="center" spacing={4}>
            <h2 className="text-xl text-bold text-gray-900">
              <span className="text-bold">
                <Trans id="label.discussion">Discussion</Trans>
              </span>
            </h2>
            <div className="p-show-post__discussion-count">{comments.length}</div>
          </HStack>

          {/* Comment Input at top */}
          <CommentInput post={post} />

          {/* Response Details - First discussion item */}
          {post.response && <ResponseDetails status={post.status} response={post.response} />}

          {/* Comments List */}
          {comments.length > 0 && (
            <VStack spacing={4}>
              {comments.map((c) => (
                <ShowComment key={c.id} post={post} comment={c} highlighted={highlightedComment === c.id} />
              ))}
            </VStack>
          )}
        </div>

        {/* Powered by Fider - bottom of page on mobile only */}
        <div className="p-show-post__powered-by-mobile">
          <PoweredByFider slot="show-post" />
        </div>
      </div>

      {/* Modals */}
      <RSSModal isOpen={isRSSModalOpen} onClose={hideRSSModal} url={`${fider.settings.baseURL}/feed/posts/${post.number}.atom`} />
      <DeletePostModal onModalClose={() => setShowDeleteModal(false)} showModal={showDeleteModal} post={post} />
      {Fider.session.isAuthenticated && Fider.session.user.isCollaborator && (
        <ResponseModal onCloseModal={() => setShowResponseModal(false)} showModal={showResponseModal} post={post} />
      )}
    </div>
  )
}
