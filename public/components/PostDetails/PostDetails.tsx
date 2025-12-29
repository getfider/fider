import "./PostDetails.scss"

import React, { useState, useEffect, useCallback } from "react"

import { Comment, Post, Tag, Vote, CurrentUser, PostStatus } from "@fider/models"
import { actions, cache, clearUrlHash, Failure, Fider, notify, timeAgo } from "@fider/services"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import IconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import { i18n } from "@lingui/core"
import IconRSS from "@fider/assets/images/heroicons-rss.svg"
import IconPencil from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconChat from "@fider/assets/images/heroicons-chat-alt-2.svg"

import { ResponseDetails, Button, UserName, Moment, Markdown, Input, Form, Icon, Avatar, Dropdown, RSSModal } from "@fider/components"
import { DiscussionPanel } from "@fider/pages/ShowPost/components/DiscussionPanel"
import CommentEditor from "@fider/components/common/form/CommentEditor"

import IconX from "@fider/assets/images/heroicons-x.svg"
import IconThumbsUp from "@fider/assets/images/heroicons-thumbsup.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"
import { FollowButton } from "@fider/pages/ShowPost/components/FollowButton"
import { VoteSection } from "@fider/pages/ShowPost/components/VoteSection"
import { DeletePostModal } from "@fider/pages/ShowPost/components/DeletePostModal"
import { ResponseModal } from "@fider/pages/ShowPost/components/ResponseModal"
import { VotesPanel } from "@fider/pages/ShowPost/components/VotesPanel"
import { TagsPanel } from "@fider/pages/ShowPost/components/TagsPanel"
import { t } from "@lingui/macro"
import { useFider } from "@fider/hooks"
import { useAttachments } from "@fider/hooks/useAttachments"

interface PostDetailsProps {
  postNumber: number
  onClose?: () => void
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
        const [postResponse, commentsResult, tagsResult, votesResult] = await Promise.all([
          fetch(`/api/v1/posts/${props.postNumber}`).then((r) => r.json()),
          fetch(`/api/v1/posts/${props.postNumber}/comments`).then((r) => r.json()),
          fetch(`/api/v1/tags`).then((r) => r.json()),
          actions.listVotes(props.postNumber),
        ])

        if (postResponse) {
          setPost(postResponse)
          setNewTitle(postResponse.title)
          setNewDescription(postResponse.description)
        }

        setComments(commentsResult || [])
        setTags(tagsResult || [])
        setVotes(votesResult.ok ? votesResult.data : [])

        // Fetch subscription status if authenticated
        if (Fider.session.isAuthenticated) {
          try {
            const subResponse = await fetch(`/api/v1/posts/${props.postNumber}/subscription`)
            if (subResponse.ok) {
              const subData = await subResponse.json()
              setSubscribed(subData.subscribed || false)
            }
          } catch (e) {
            // Ignore subscription errors
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
    if (showSuccess) {
      cache.session.remove("POST_CREATED_SUCCESS")
      notify.success(t({ id: "mysettings.notification.event.newpostcreated", message: "Your idea has been added 👍" }))
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
      <div className="post-details">
        <div className="post-details__loading">Loading...</div>
      </div>
    )
  }

  return (
    <div className="post-details">
      <div className="post-details__main-col">
        <div className="post-details__header-col">
          <VStack spacing={8}>
            <HStack justify="between">
              <VStack align="start">
                {!editMode && (
                  <HStack>
                    <Avatar user={post.user} />
                    <VStack spacing={1}>
                      <UserName user={post.user} />
                      <Moment className="text-muted" locale={Fider.currentLocale} date={post.createdAt} />
                    </VStack>
                  </HStack>
                )}
              </VStack>
              <RSSModal isOpen={isRSSModalOpen} onClose={hideRSSModal} url={`${fider.settings.baseURL}/feed/posts/${post.number}.atom`} />

              {!editMode && (
                <HStack>
                  {props.onClose && (
                    <Button onClick={props.onClose}>
                      <Icon sprite={IconX} />
                    </Button>
                  )}
                  <Dropdown position="left" renderHandle={<Icon sprite={IconDotsHorizontal} width="24" height="24" />}>
                    <Dropdown.ListItem onClick={onActionSelected("copy")} icon={IconDuplicate}>
                      <Trans id="action.copylink">Copy link</Trans>
                    </Dropdown.ListItem>
                    {Fider.session.tenant.isFeedEnabled && (
                      <Dropdown.ListItem onClick={onActionSelected("feed")} icon={IconRSS}>
                        <Trans id="action.commentsfeed">Comment Feed</Trans>
                      </Dropdown.ListItem>
                    )}
                    {Fider.session.isAuthenticated && canEditPost(Fider.session.user, post) && (
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
                  <h1 className="text-large text-break">{post.title}</h1>
                </>
              )}
            </div>

            <DeletePostModal onModalClose={() => setShowDeleteModal(false)} showModal={showDeleteModal} post={post} />
            {Fider.session.isAuthenticated && Fider.session.user.isCollaborator && (
              <ResponseModal onCloseModal={() => setShowResponseModal(false)} showModal={showResponseModal} post={post} />
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
                  {post.description && <Markdown className="description" text={post.description} style="full" />}
                  {!post.description && (
                    <em className="text-muted">
                      <Trans id="showpost.message.nodescription">No description provided.</Trans>
                    </em>
                  )}
                </>
              )}
            </VStack>
            <div className="mt-2">
              <TagsPanel post={post} tags={tags} />
            </div>

            <VStack spacing={4}>
              {!editMode ? (
                <HStack justify="between" align="start">
                  <VoteSection post={post} votes={post.votesCount} />
                  <FollowButton post={post} subscribed={subscribed} />
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

            <ResponseDetails status={post.status} response={post.response} />
          </VStack>
        </div>

        <div className="post-details__discussion_col">
          <DiscussionPanel post={post} comments={comments} highlightedComment={highlightedComment} />
        </div>
      </div>
      <div className="post-details__action-col">
        <VotesPanel post={post} votes={votes} />
      </div>
    </div>
  )
}
