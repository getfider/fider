import "./ShowPost.page.scss"

import React, { useState, useEffect, useCallback } from "react"

import { Comment, Post, Tag, Vote, ImageUpload, CurrentUser, PostStatus } from "@fider/models"
import { actions, clearUrlHash, Failure, Fider, notify, timeAgo } from "@fider/services"
import IconDotsHorizontal from "@fider/assets/images/heroicons-dots-horizontal.svg"
import IconRss from "@fider/assets/images/heroicons-rss.svg"

import {
  ResponseDetails,
  Button,
  UserName,
  Moment,
  Markdown,
  Input,
  Form,
  TextArea,
  MultiImageUploader,
  ImageViewer,
  Icon,
  Header,
  PoweredByFider,
  Avatar,
  Dropdown,
} from "@fider/components"
import { DiscussionPanel } from "./components/DiscussionPanel"

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
  const [showResponseModal, setShowResponseModal] = useState(false)
  const [newTitle, setNewTitle] = useState(props.post.title)
  const [newDescription, setNewDescription] = useState(props.post.description)
  const [attachments, setAttachments] = useState<ImageUpload[]>([])
  const [highlightedComment, setHighlightedComment] = useState<number | undefined>(undefined)
  const [error, setError] = useState<Failure | undefined>(undefined)

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

                  {!editMode && (
                    <HStack>
                      {Fider.session.tenant.isFeedEnabled && (
                        <a title="ATOM Feed (Comments)" type="application/atom+xml" href={`/feed/posts/${props.post.id}.atom`}>
                          <Icon sprite={IconRss} width="24" height="24" />
                        </a>
                      )}
                      <Dropdown position="left" renderHandle={<Icon sprite={IconDotsHorizontal} width="24" height="24" />}>
                        <Dropdown.ListItem onClick={onActionSelected("copy")}>
                          <Trans id="action.copylink">Copy link</Trans>
                        </Dropdown.ListItem>
                        {Fider.session.isAuthenticated && canEditPost(Fider.session.user, props.post) && (
                          <>
                            <Dropdown.ListItem onClick={onActionSelected("edit")}>
                              <Trans id="action.edit">Edit</Trans>
                            </Dropdown.ListItem>
                            {Fider.session.user.isCollaborator && (
                              <Dropdown.ListItem onClick={onActionSelected("status")}>
                                <Trans id="action.respond">Respond</Trans>
                              </Dropdown.ListItem>
                            )}
                          </>
                        )}
                        {canDeletePost() && (
                          <Dropdown.ListItem onClick={onActionSelected("delete")} className="text-red-700">
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
                      <h1 className="text-large">{props.post.title}</h1>
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
                      <TextArea field="description" value={newDescription} onChange={setNewDescription} />
                      <MultiImageUploader field="attachments" bkeys={props.attachments} maxUploads={3} onChange={setAttachments} />
                    </Form>
                  ) : (
                    <>
                      {props.post.description && <Markdown className="description" text={props.post.description} style="full" />}
                      {!props.post.description && (
                        <em className="text-muted">
                          <Trans id="showpost.message.nodescription">No description provided.</Trans>
                        </em>
                      )}
                      {props.attachments.map((x) => (
                        <ImageViewer key={x} bkey={x} />
                      ))}
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
