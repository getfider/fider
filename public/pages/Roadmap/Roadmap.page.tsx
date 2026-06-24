import "./Roadmap.page.scss"
import IconArrowLeft from "@fider/assets/images/heroicons-arrowleft.svg"
import IconCheckCircle from "@fider/assets/images/heroicons-check-circle.svg"

import React, { useState, useCallback } from "react"
import { Post, Status, Tag } from "@fider/models"
import { Header, Button, Icon, ResponseLozenge, ShowTag, Moment } from "@fider/components"
import { VStack, HStack } from "@fider/components/layout"
import { useFider, usePostOverlay } from "@fider/hooks"
import { actions } from "@fider/services"
import { PostDetails } from "@fider/components/PostDetails"
import { Trans } from "@lingui/react/macro"

interface RoadmapColumnData {
  status: Status
  posts: Post[]
}

interface RoadmapPageProps {
  columns?: RoadmapColumnData[]
  tags?: Tag[]
}

interface RoadmapColumnProps {
  status: string
  kind?: string
  posts: Post[]
  tags: Tag[]
  currentLimit: number
  onShowMore: () => void
  onPostClick?: (postNumber: number, slug: string) => void
}

// Must match the Limit sent by the server-side RoadmapPage handler. The "Show
// more" link uses posts.length >= currentLimit as its heuristic, so the SSR
// payload and the client's first render must agree on the initial cap.
const ROADMAP_DEFAULT_LIMIT = 10
const ROADMAP_LIMIT_STEP = 10

const RoadmapPost = (props: { post: Post; tags: Tag[]; status: string; kind?: string; onPostClick?: (postNumber: number, slug: string) => void }) => {
  const fider = useFider()
  const isModerationEnabled = fider.session.tenant.isModerationEnabled
  const isPending = isModerationEnabled && !props.post.isApproved
  // Anything in a closed-completed kind hides the upvote affordance.
  const isCompleted = props.kind === "closed-completed" || props.status === "completed"

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    if (props.onPostClick) {
      e.preventDefault()
      props.onPostClick(props.post.number, props.post.slug)
    }
  }

  return (
    <a href={`/posts/${props.post.number}/${props.post.slug}`} className="c-roadmap-post-link" onClick={handleClick}>
      <VStack className="c-roadmap-post w-full" spacing={2}>
        <HStack spacing={2} align="start" className="w-full">
          <h3 className="c-roadmap-post__title text-break">{props.post.title}</h3>
          {isPending && (
            <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded flex-shrink-0">
              <Trans id="post.pending">pending</Trans>
            </span>
          )}
        </HStack>
        {props.tags.length >= 1 && (
          <HStack spacing={0} className="gap-x-4 flex-wrap">
            {props.tags.map((tag) => (
              <ShowTag key={tag.id} tag={tag} />
            ))}
          </HStack>
        )}
        {isCompleted && props.post.response?.respondedAt ? (
          <HStack spacing={1} className="c-roadmap-post__completed flex-items-center">
            <Icon sprite={IconCheckCircle} className="h-4 w-4 text-green-600" />
            <Moment locale={fider.currentLocale} date={props.post.response.respondedAt} />
          </HStack>
        ) : (
          <span className="c-roadmap-post__votes">
            <span className="text-semibold">{props.post.votesCount}</span>{" "}
            {props.post.votesCount === 1 ? <Trans id="label.vote">Vote</Trans> : <Trans id="label.votes">Votes</Trans>}
          </span>
        )}
      </VStack>
    </a>
  )
}

const RoadmapColumn = (props: RoadmapColumnProps) => {
  // If we received at least as many posts as we asked for there may be more on
  // the server — same heuristic the Home feed uses (PostsContainer.getShowMoreLink).
  const hasMore = props.posts.length >= props.currentLimit

  return (
    <div className="c-roadmap-column">
      <div className="c-roadmap-column__header">
        <ResponseLozenge status={props.status} response={null} />
      </div>
      <div className="c-roadmap-column__body">
        {props.posts.map((post) => (
          <RoadmapPost
            key={post.id}
            post={post}
            tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)}
            status={props.status}
            kind={props.kind}
            onPostClick={props.onPostClick}
          />
        ))}
        {hasMore && (
          <div className="my-4 text-center">
            <a
              href="#"
              className="text-primary-base text-medium hover:underline"
              onClick={(e) => {
                e.preventDefault()
                props.onShowMore()
              }}
            >
              <Trans id="roadmap.column.showmore">Show more</Trans>
            </a>
          </div>
        )}
      </div>
    </div>
  )
}

interface ColumnState {
  status: Status
  posts: Post[]
  limit: number
}

const RoadmapBoard = (props: RoadmapPageProps) => {
  const [columns, setColumns] = useState<ColumnState[]>((props.columns || []).map((c) => ({ status: c.status, posts: c.posts, limit: ROADMAP_DEFAULT_LIMIT })))
  const tags = props.tags || []

  const reloadPosts = useCallback(async () => {
    const next = await Promise.all(
      columns.map(async (col) => {
        const result = await actions.searchPosts({ statuses: [col.status.slug], limit: col.limit })
        return result.ok ? { ...col, posts: result.data as Post[] } : col
      })
    )
    setColumns(next)
  }, [columns])

  const showMore = async (slug: string) => {
    const target = columns.find((c) => c.status.slug === slug)
    if (!target) return
    const nextLimit = target.limit + ROADMAP_LIMIT_STEP
    const result = await actions.searchPosts({ statuses: [slug], limit: nextLimit })
    if (!result.ok) return
    setColumns(columns.map((c) => (c.status.slug === slug ? { ...c, posts: result.data as Post[], limit: nextLimit } : c)))
  }

  const { selectedPostId, handlePostClick, handleCloseOverlay, setIsPostDirty } = usePostOverlay({
    basePath: "/roadmap",
    onPostClosed: () => reloadPosts(),
  })

  // Any column the admin opted into the roadmap (via show_on_roadmap) counts —
  // if Completed is enabled and has posts, the board is not "waiting for its
  // first update."
  const hasAnyOpenWork = columns.some((c) => c.posts.length > 0)

  if (!hasAnyOpenWork && selectedPostId === null) {
    return <RoadmapBlankState />
  }

  return (
    <div id="p-roadmap" className="page container">
      <div style={selectedPostId !== null ? { display: "none" } : undefined}>
        <VStack spacing={4}>
          <div className="c-roadmap-board">
            {columns.map((col) => (
              <RoadmapColumn
                key={col.status.slug}
                status={col.status.slug}
                kind={col.status.kind}
                posts={col.posts}
                tags={tags}
                currentLimit={col.limit}
                onShowMore={() => showMore(col.status.slug)}
                onPostClick={handlePostClick}
              />
            ))}
          </div>
        </VStack>
      </div>
      {selectedPostId !== null && (
        <div>
          <Button onClick={handleCloseOverlay} variant="link">
            <HStack spacing={2}>
              <Icon sprite={IconArrowLeft} />
              <span className="text-body clickable text-blue-600 hover">
                <Trans id="postdetails.backtoroadmap">Back to roadmap</Trans>
              </span>
            </HStack>
          </Button>
          <PostDetails postNumber={selectedPostId} onDataChanged={() => setIsPostDirty(true)} />
        </div>
      )}
    </div>
  )
}

const SkeletonCard = () => (
  <div className="c-roadmap-upsell__card">
    <div className="c-roadmap-upsell__bar c-roadmap-upsell__bar--lg" />
    <div className="c-roadmap-upsell__bar c-roadmap-upsell__bar--sm" />
    <div className="c-roadmap-upsell__bar c-roadmap-upsell__bar--md" />
    <div className="c-roadmap-upsell__bar c-roadmap-upsell__bar--md" />
  </div>
)

const SkeletonColumn = ({ status }: { status: string }) => (
  <div className="c-roadmap-column">
    <div className="c-roadmap-column__header">
      <ResponseLozenge status={status} response={null} />
    </div>
    <div className="c-roadmap-column__body">
      <SkeletonCard />
      <SkeletonCard />
    </div>
  </div>
)

const RoadmapSkeletonBackdrop = () => (
  <div className="c-roadmap-upsell__skeleton" aria-hidden="true">
    <div className="c-roadmap-board">
      <SkeletonColumn status="planned" />
      <SkeletonColumn status="started" />
      <SkeletonColumn status="completed" />
    </div>
  </div>
)

const RoadmapUpsell = () => {
  const fider = useFider()
  const isAdmin = fider.session.isAuthenticated && fider.session.user.isAdministrator
  const showBillingCta = isAdmin && fider.settings.isBillingEnabled

  return (
    <div id="p-roadmap-upsell" className="page container">
      <RoadmapSkeletonBackdrop />
      <VStack spacing={4} className="c-roadmap-upsell flex-items-center text-center">
        <h1 className="c-roadmap-upsell__title text-display">
          <Trans id="roadmap.upsell.title">See what&apos;s happening in the Roadmap view</Trans>
        </h1>
        <p className="c-roadmap-upsell__subtitle text-muted">
          <Trans id="roadmap.upsell.description">Upgrade to Pro to unlock your Roadmap</Trans>
        </p>
        {showBillingCta && (
          <a href="/admin/billing">
            <Button variant="primary" size="large">
              <Trans id="roadmap.upsell.billing">Upgrade to PRO</Trans>
            </Button>
          </a>
        )}
      </VStack>
    </div>
  )
}

const RoadmapBlankState = () => (
  <div id="p-roadmap-blank" className="page container">
    <RoadmapSkeletonBackdrop />
    <VStack spacing={4} className="c-roadmap-upsell flex-items-center text-center">
      <h1 className="c-roadmap-upsell__title text-display">
        <Trans id="roadmap.blank.title">Your roadmap is waiting for its first update</Trans>
      </h1>
      <p className="c-roadmap-upsell__subtitle text-muted">
        <Trans id="roadmap.blank.description">Mark posts as planned or in progress and they&apos;ll show up here on the roadmap.</Trans>
      </p>
    </VStack>
  </div>
)

const RoadmapPage = (props: RoadmapPageProps) => {
  const fider = useFider()
  const hasRoadmap = fider.isSingleHostMode() || fider.session.tenant.isPro

  return (
    <>
      <Header />
      {hasRoadmap ? <RoadmapBoard {...props} /> : <RoadmapUpsell />}
    </>
  )
}

export default RoadmapPage
