import "./ShowPost.page.scss"

import React, { useState, useEffect, useCallback } from "react"

import { Comment, Post, Tag, Vote, ImageUpload, CurrentUser, PostStatus } from "@fider/models"
import { actions, cache, clearUrlHash, Failure, Fider, notify, timeAgo } from "@fider/services"
import { getPostDetails } from "@fider/services/actions/post-details"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import IconDuplicate from "@fider/assets/images/heroicons-duplicate.svg"
import { i18n } from "@lingui/core"
import IconRSS from "@fider/assets/images/heroicons-rss.svg"
import IconPencil from "@fider/assets/images/heroicons-pencil-alt.svg"
import IconChat from "@fider/assets/images/heroicons-chat-alt-2.svg"

import { ResponseDetails, Button, UserName, Moment, Markdown, Input, Form, Icon, PoweredByFider, Avatar, Dropdown, Loader } from "@fider/components"
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

// Props for ShowPostPage component
interface ShowPostPageProps {
  // Or provide a post number to load the data
  postNumber: number
}

const oneHour = 3600
const canEditPost = (user: CurrentUser, post: Post) => {
  if (user.isCollaborator) {
    return true
  }

  return user.id === post.user.id && timeAgo(post.createdAt) <= oneHour
}

export default function ShowPostPage(props: ShowPostPageProps) {
  // State for loading post details
  const [loading, setLoading] = useState(props.postNumber !== undefined)
  const [loadError, setLoadError] = useState<string | null>(null)

  // State for post data
  const [post, setPost] = useState<Post | undefined>(undefined)
  const [comments, setComments] = useState<Comment[] | undefined>(undefined)
  const [tags, setTags] = useState<Tag[] | undefined>(undefined)
  const [votes, setVotes] = useState<Vote[] | undefined>(undefined)
  const [subscribed, setSubscribed] = useState<boolean | undefined>(undefined)
  const [attachments, setAttachments] = useState<string[] | undefined>(undefined)

  // Component state - initialize with defaults
  const [editMode, setEditMode] = useState(false)
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [showResponseModal, setShowResponseModal] = useState(false)
  const [newTitle, setNewTitle] = useState("")
  const [newDescription, setNewDescription] = useState("")
  const [attachmentsList, setAttachmentsList] = useState<ImageUpload[]>([])
  const [highlightedComment, setHighlightedComment] = useState<number | undefined>(undefined)
  const [error, setError] = useState<Failure | undefined>(undefined)

  // Load post data if postNumber is provided
  useEffect(() => {
    if (props.postNumber && !post) {
      loadPost()
    }
  }, [props.postNumber, post])

  const loadPost = () => {
    setLoading(true)
    setLoadError(null)

    getPostDetails(props.postNumber).then((response) => {
      setLoading(false)
      if (response.ok) {
        const data = response.data
        setPost(data.post)
        setComments(data.comments)
        setTags(data.tags)
        setVotes(data.votes)
        setSubscribed(data.subscribed)
        setAttachments(data.attachments)
        setNewTitle(data.post.title)
        setNewDescription(data.post.description)
      } else {
        setLoadError("Failed to load post")
      }
    })
  }

  const handleHashChange = useCallback(
    (e?: Event) => {
      if (!comments) return

      const hash = window.location.hash
      const result = /#comment-([0-9]+)/.exec(hash)

      let newHighlightedComment
      if (result === null) {
        // No match
        newHighlightedComment = undefined
      } else {
        // Match, extract numeric ID
        const id = parseInt(result[1])
        if (comments.map((comment) => comment.id).includes(id)) {
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
    [comments]
  )

  useEffect(() => {
    if (comments) {
      handleHashChange()
      window.addEventListener("hashchange", handleHashChange)
      return () => {
        window.removeEventListener("hashchange", handleHashChange)
      }
    }
  }, [handleHashChange, comments])

  useEffect(() => {
    const showSuccess = cache.session.get("POST_CREATED_SUCCESS")
    if (showSuccess) {
      cache.session.remove("POST_CREATED_SUCCESS")
      // Show success message/toast
      notify.success(t({ id: "mysettings.notification.event.newpostcreated", message: "Your idea has been added 👍" }))
    }
  }, [])

  // Initialize attachments from props when entering edit mode
  useEffect(() => {
    if (editMode && attachments) {
      // Convert attachment bkeys to ImageUpload objects
      const initialAttachments = attachments.map((bkey) => ({
        bkey,
        remove: false,
      })) as ImageUpload[]

      setAttachmentsList(initialAttachments)
    }
  }, [editMode, attachments])

  // Show loading state
  if (loading) {
    return (
      <div className="p-6 text-center">
        <Loader />
      </div>
    )
  }

  // Show error state
  if (loadError) {
    return <div className="p-6 text-center text-red-500">{loadError}</div>
  }

  // If no data is available, return null
  if (!post || !comments || !tags || !votes || subscribed === undefined || !attachments) {
    return null
  }

  const saveChanges = async () => {
    const result = await actions.updatePost(post.number, newTitle, newDescription, attachmentsList)
    if (result.ok) {
      loadPost()
      setEditMode(false)
    } else {
      setError(result.error)
    }
  }

  const canDeletePost = () => {
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

  const handleImageUploaded = (image: ImageUpload) => {
    setAttachmentsList((prev) => {
      // If this is a removal request, find and mark the attachment for removal
      if (image.remove && image.bkey) {
        return prev.map((att) => (att.bkey === image.bkey ? { ...att, remove: true } : att))
      }
      // Otherwise add the new upload
      return [...prev, image]
    })
  }

  const onActionSelected = (action: "copy" | "delete" | "status" | "edit") => () => {
    if (action === "copy") {
      navigator.clipboard.writeText(window.location.href)
      notify.success(<Trans id="showpost.copylink.success">Link copied to clipboard</Trans>)
    } else if (action === "delete") {
      setShowDeleteModal(true)
    } else if (action === "status") {
      setShowResponseModal(true)
    } else if (action === "edit") {
      startEdit()
    }
  }

  return (
    <>
      <div id="p-show-post" className="page container">
        <div className="p-show-post">
          <div className="p-show-post__main-col">
            <div className="p-show-post__header-col">
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

                  {!editMode && (
                    <HStack>
                      <Dropdown position="left" renderHandle={<Icon sprite={IconDotsHorizontal} width="24" height="24" />}>
                        <Dropdown.ListItem onClick={onActionSelected("copy")} icon={IconDuplicate}>
                          <Trans id="action.copylink">Copy link</Trans>
                        </Dropdown.ListItem>
                        {Fider.session.tenant.isFeedEnabled && (
                          <Dropdown.ListItem type="application/atom+xml" href={`/feed/posts/${post.number}.atom`} icon={IconRSS}>
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
                      <h1 className="text-large">{post.title}</h1>
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

            <div className="p-show-post__discussion_col">
              <DiscussionPanel post={post} comments={comments} highlightedComment={highlightedComment} />
            </div>
          </div>
          <div className="p-show-post__action-col">
            <VotesPanel post={post} votes={votes} />
            <PoweredByFider slot="show-post" className="mt-3" />
          </div>
        </div>
      </div>
    </>
  )
}
