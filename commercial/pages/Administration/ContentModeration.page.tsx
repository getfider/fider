import "./ContentModeration.page.scss"

import React, { useState, useEffect } from "react"
import { Button, Avatar, Loader, Icon, Markdown } from "@fider/components/common"
import { Header } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { actions, basePath, chopString, http, notify } from "@fider/services"
import { User, UserStatus } from "@fider/models"
import { useFider } from "@fider/hooks"
import { Trans } from "@lingui/react/macro"
import IconCheck from "@fider/assets/images/heroicons-check.svg"
import IconX from "@fider/assets/images/heroicons-x.svg"
import IconShieldCheck from "@fider/assets/images/heroicons-shieldcheck.svg"
import IconBan from "@fider/assets/images/heroicons-x-circle.svg"
import { Moment } from "@fider/components/common"

interface ModerationItem {
  type: "post" | "comment"
  id: number
  postId?: number
  postNumber?: number
  postSlug?: string
  title?: string
  content: string
  user: User
  createdAt: string
  postTitle?: string
}

interface ContentModerationPageState {
  items: ModerationItem[]
  loading: boolean
}

const ContentModerationPage = () => {
  const [state, setState] = useState<ContentModerationPageState>({
    items: [],
    loading: true,
  })
  const fider = useFider()

  const fetchItems = async () => {
    setState((prev) => ({ ...prev, loading: true }))
    const result = await http.get<{ items: ModerationItem[] }>("/_api/admin/moderation/items")
    if (result.ok) {
      setState({ items: result.data.items, loading: false })
    } else {
      setState((prev) => ({ ...prev, loading: false }))
      notify.error(<Trans id="moderation.fetch.error">Failed to fetch moderation items</Trans>)
    }
  }

  useEffect(() => {
    fetchItems()
  }, [])

  const handleApprovePost = async (postId: number) => {
    const result = await actions.approvePost(postId)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.published">Post published successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.publish.error">Failed to publish post</Trans>)
    }
  }

  const handleDeclinePost = async (postId: number) => {
    const result = await actions.declinePost(postId)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.deleted">Post deleted successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.delete.error">Failed to delete post</Trans>)
    }
  }

  const handleApproveComment = async (commentId: number) => {
    const result = await actions.approveComment(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.published">Comment published successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.publish.error">Failed to publish comment</Trans>)
    }
  }

  const handleDeclineComment = async (commentId: number) => {
    const result = await actions.declineComment(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.deleted">Comment deleted successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.delete.error">Failed to delete comment</Trans>)
    }
  }

  const handleApprovePostAndVerify = async (postId: number) => {
    const result = await actions.approvePostAndVerify(postId)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.published.verified">Post published and user verified</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.publish.verify.error">Failed to publish post and verify user</Trans>)
    }
  }

  const handleDeclinePostAndBlock = async (postId: number) => {
    const result = await actions.declinePostAndBlock(postId)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.deleted.blocked">Post deleted and user blocked</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.delete.block.error">Failed to delete post and block user</Trans>)
    }
  }

  const handleApproveCommentAndVerify = async (commentId: number) => {
    const result = await actions.approveCommentAndVerify(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.published.verified">Comment published and user verified</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.publish.verify.error">Failed to publish comment and verify user</Trans>)
    }
  }

  const handleDeclineCommentAndBlock = async (commentId: number) => {
    const result = await actions.declineCommentAndBlock(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.deleted.blocked">Comment deleted and user blocked</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.delete.block.error">Failed to delete comment and block user</Trans>)
    }
  }

  const renderDivider = (title: string, count: number) => {
    return (
      <div className="c-moderation-page__divider">
        <div className="c-moderation-page__divider-title">
          {title} ({count})
        </div>
      </div>
    )
  }

  const handlePostClick = (link: string) => {
    window.location.href = link
  }

  const renderModerationItem = (item: ModerationItem) => {
    const title = item.type == "post" ? item.title : item.postTitle
    const link = item.type == "post" ? `${basePath()}/posts/${item.postNumber}/${item.postSlug}` : `${basePath()}/posts/${item.postNumber}/${item.postSlug}#comment-${item.id}`
    const blocked = item.user.status === UserStatus.Blocked && <span className="text-red-700">blocked</span>

    return (
      <div key={`${item.type}-${item.id}`} className="c-moderation-item" onClick={() => handlePostClick(link)}>
        <div className="c-moderation-item__content">
          <div className="c-moderation-item__header">
            <HStack spacing={4} align="start">
              <Avatar user={item.user} size="large" />
              <VStack spacing={1} className="flex-grow c-moderation-item__user-info">
                <div className="text-semibold c-moderation-item__user-name">{item.user.name}</div>
                <div className="c-moderation-item__user-email">&lt;{item.user.email}&gt;</div>
                {blocked && <span className="text-red-700">{blocked}</span>}
              </VStack>
            </HStack>
            <div className="c-moderation-item__timestamp">
              <Moment date={item.createdAt} locale={fider.currentLocale} />
            </div>
          </div>

          <VStack spacing={2} className="c-moderation-item__body">
            {item.type === "post" && <h3 className="text-semibold m-0">{title}</h3>}
            <div className="text-body text-break">
              <Markdown text={chopString(item.content, 200)} style="plainText" />
            </div>
            {item.type === "comment" && <div className="text-muted text-break">{title}</div>}

            <div className="c-moderation-item__actions invisible" onClick={(e) => e.stopPropagation()}>
              <Button size="small" variant="secondary" onClick={() => (item.type === "post" ? handleApprovePost(item.id) : handleApproveComment(item.id))}>
                <Icon sprite={IconCheck} />
                <span>
                  <Trans id="action.publish">Publish</Trans>
                </span>
              </Button>
              <Button size="small" variant="secondary" onClick={() => (item.type === "post" ? handleDeclinePost(item.id) : handleDeclineComment(item.id))}>
                <Icon sprite={IconX} />
                <span>
                  <Trans id="action.delete">Delete</Trans>
                </span>
              </Button>
              <Button
                size="small"
                variant="secondary"
                onClick={() => (item.type === "post" ? handleApprovePostAndVerify(item.id) : handleApproveCommentAndVerify(item.id))}
              >
                <Icon sprite={IconShieldCheck} />
                <span>
                  <Trans id="action.publish.verify">Publish & Trust</Trans>
                </span>
              </Button>
              <Button
                size="small"
                variant="secondary"
                onClick={() => (item.type === "post" ? handleDeclinePostAndBlock(item.id) : handleDeclineCommentAndBlock(item.id))}
              >
                <Icon sprite={IconBan} />
                <span>
                  <Trans id="action.delete.block">Delete & Block</Trans>
                </span>
              </Button>
            </div>
          </VStack>
        </div>
      </div>
    )
  }

  const posts = state.items.filter((item) => item.type === "post")
  const comments = state.items.filter((item) => item.type === "comment")

  return (
    <>
      <Header />
      <div id="p-admin-moderation" className="page container">
        <VStack spacing={2}>
          <h1 className="text-display">
            <Trans id="moderation.title">Moderation Queue</Trans>
          </h1>
          <p className="text-body">
            <Trans id="moderation.subtitle">
              These ideas and comments are from people outside of your trusted users list, you decide if they get published.
            </Trans>
          </p>
        </VStack>

        <div className="c-moderation-page">
          {state.loading ? (
            <Loader />
          ) : state.items.length === 0 ? (
            <div className="text-center p-8 rounded mt-4">
              <p>
                <Trans id="moderation.empty">All content has been moderated. You&apos;re all caught up!</Trans>
              </p>
            </div>
          ) : (
            <>
              {posts.length > 0 && (
                <>
                  {renderDivider("New ideas", posts.length)}
                  <div className="c-moderation-page__list">{posts.map(renderModerationItem)}</div>
                </>
              )}

              {comments.length > 0 && (
                <>
                  {renderDivider("New comments", comments.length)}
                  <div className="c-moderation-page__list">{comments.map(renderModerationItem)}</div>
                </>
              )}
            </>
          )}
        </div>
      </div>
    </>
  )
}

export default ContentModerationPage
