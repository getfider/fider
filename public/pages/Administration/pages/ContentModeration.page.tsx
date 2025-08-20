import "./ContentModeration.page.scss"

import React, { useState, useEffect } from "react"
import { Button, Avatar, Loader, Icon, Markdown } from "@fider/components/common"
import { Header } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { actions, chopString, http, notify } from "@fider/services"
import { User } from "@fider/models"
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
      notify.success(<Trans id="moderation.post.approved">Post approved successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.approve.error">Failed to approve post</Trans>)
    }
  }

  const handleDeclinePost = async (postId: number) => {
    const result = await actions.declinePost(postId)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.declined">Post declined successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.decline.error">Failed to decline post</Trans>)
    }
  }

  const handleApproveComment = async (commentId: number) => {
    const result = await actions.approveComment(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.approved">Comment approved successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.approve.error">Failed to approve comment</Trans>)
    }
  }

  const handleDeclineComment = async (commentId: number) => {
    const result = await actions.declineComment(commentId)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.declined">Comment declined successfully</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.decline.error">Failed to decline comment</Trans>)
    }
  }

  const handleApprovePostAndVerify = async (postId: number) => {
    const result = await http.post(`/api/v1/admin/moderation/posts/${postId}/approve-and-verify`)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.approved.verified">Post approved and user verified</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.approve.verify.error">Failed to approve post and verify user</Trans>)
    }
  }

  const handleDeclinePostAndBlock = async (postId: number) => {
    const result = await http.post(`/api/v1/admin/moderation/posts/${postId}/decline-and-block`)
    if (result.ok) {
      notify.success(<Trans id="moderation.post.declined.blocked">Post declined and user blocked</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "post" && item.id === postId)),
      }))
    } else {
      notify.error(<Trans id="moderation.post.decline.block.error">Failed to decline post and block user</Trans>)
    }
  }

  const handleApproveCommentAndVerify = async (commentId: number) => {
    const result = await http.post(`/api/v1/admin/moderation/comments/${commentId}/approve-and-verify`)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.approved.verified">Comment approved and user verified</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.approve.verify.error">Failed to approve comment and verify user</Trans>)
    }
  }

  const handleDeclineCommentAndBlock = async (commentId: number) => {
    const result = await http.post(`/api/v1/admin/moderation/comments/${commentId}/decline-and-block`)
    if (result.ok) {
      notify.success(<Trans id="moderation.comment.declined.blocked">Comment declined and user blocked</Trans>)
      setState((prev) => ({
        ...prev,
        items: prev.items.filter((item) => !(item.type === "comment" && item.id === commentId)),
      }))
    } else {
      notify.error(<Trans id="moderation.comment.decline.block.error">Failed to decline comment and block user</Trans>)
    }
  }

  const renderDivider = (title: string, count: number) => {
    return (
      <HStack spacing={2} className="py-4 w-full">
        <div className="text-blue-500 text-title">
          {title} ({count})
        </div>
        <hr className="border-blue-500 flex-grow"></hr>
      </HStack>
    )
  }

  const handlePostClick = (link: string) => {
    window.location.href = link
  }

  const renderModerationItem = (item: ModerationItem) => {
    const containerClasses = `c-moderation-item flex flex-y p-3 rounded-md hover clickable`

    const title = item.type == "post" ? item.title : item.postTitle
    const link = item.type == "post" ? `/posts/${item.postNumber}/${item.postSlug}` : `/posts/${item.postNumber}/${item.postSlug}#comment-${item.id}`

    return (
      <div key={`${item.type}-${item.id}`} className={containerClasses} onClick={() => handlePostClick(link)}>
        <div className="c-moderation-item__content">
          <HStack spacing={4} align="start" justify="between">
            <HStack spacing={4} align="start">
              <Avatar user={item.user} size="normal" />
              <VStack spacing={2} className="flex-grow">
                <>
                  <span className="text-medium">
                    {item.user.name} <span className="text-normal">&lt;{item.user.email}&gt;</span>
                  </span>
                  <h3 className="text-medium m-0">{title}</h3>
                  <p className="m-0 text-body text-break">
                    <Markdown text={chopString(item.content, 200)} style="plainText" />
                  </p>
                </>

                <div className="c-moderation-item__actions invisible pt-1" onClick={(e) => e.stopPropagation()}>
                  <HStack spacing={2}>
                    <Button
                      size="small"
                      variant="secondary"
                      onClick={() => (item.type === "post" ? handleApprovePost(item.id) : handleApproveComment(item.id))}
                    >
                      <Icon sprite={IconCheck} />
                      <span>
                        <Trans id="action.approve">Approve</Trans>
                      </span>
                    </Button>
                    <Button
                      size="small"
                      variant="secondary"
                      onClick={() => (item.type === "post" ? handleDeclinePost(item.id) : handleDeclineComment(item.id))}
                    >
                      <Icon sprite={IconX} />
                      <span>
                        <Trans id="action.decline">Decline</Trans>
                      </span>
                    </Button>
                    <Button
                      size="small"
                      variant="secondary"
                      onClick={() => (item.type === "post" ? handleApprovePostAndVerify(item.id) : handleApproveCommentAndVerify(item.id))}
                    >
                      <Icon sprite={IconShieldCheck} />
                      <span>
                        <Trans id="action.approve.verify">Approve & Verify</Trans>
                      </span>
                    </Button>
                    <Button
                      size="small"
                      variant="secondary"
                      onClick={() => (item.type === "post" ? handleDeclinePostAndBlock(item.id) : handleDeclineCommentAndBlock(item.id))}
                    >
                      <Icon sprite={IconBan} />
                      <span>
                        <Trans id="action.decline.block">Decline & Block</Trans>
                      </span>
                    </Button>
                  </HStack>
                </div>
              </VStack>
            </HStack>
            <Moment date={item.createdAt} locale={fider.currentLocale} />
          </HStack>
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
        <h1 className="text-large">
          <Trans id="moderation.title">Moderation Queue</Trans>
        </h1>
        <p className="text-body text-lg mt-3">
          <Trans id="moderation.subtitle">These posts and comments are from people outside of your trusted users list, you decide if they get published.</Trans>
        </p>

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
            <div className="mt-4">
              {posts.length > 0 && (
                <>
                  {renderDivider("New ideas", posts.length)}
                  <div className="mb-2">
                    <div className="flex flex-y">{posts.map(renderModerationItem)}</div>
                  </div>
                </>
              )}

              {comments.length > 0 && (
                <>
                  {renderDivider("New comments", comments.length)}
                  <div className="mb-2">
                    <div className="flex flex-y">{comments.map(renderModerationItem)}</div>
                  </div>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </>
  )
}

export default ContentModerationPage
