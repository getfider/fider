import "./Roadmap.page.scss"
import "@fider/pages/Home/components/PostsContainer.scss"
import IconArrowLeft from "@fider/assets/images/heroicons-arrowleft.svg"

import React, { useState, useCallback } from "react"
import { Post, Tag } from "@fider/models"
import { Header, Button, Icon, ResponseLozenge } from "@fider/components"
import { ListPosts } from "@fider/pages/Home/components/ListPosts"
import { VStack, HStack } from "@fider/components/layout"
import { useFider, usePostOverlay } from "@fider/hooks"
import { actions } from "@fider/services"
import { PostDetails } from "@fider/components/PostDetails"
import { Trans } from "@lingui/react/macro"

interface RoadmapPageProps {
  plannedPosts?: Post[]
  startedPosts?: Post[]
  completedPosts?: Post[]
  tags?: Tag[]
}

interface RoadmapColumnProps {
  status: string
  posts: Post[]
  tags: Tag[]
  onPostClick?: (postNumber: number, slug: string) => void
}

const RoadmapColumn = (props: RoadmapColumnProps) => {
  return (
    <div className="c-roadmap-column">
      <div className="c-roadmap-column__header">
        <ResponseLozenge status={props.status} response={null} />
      </div>
      <div className="c-roadmap-column__body">
        <ListPosts posts={props.posts} tags={props.tags} showStatus={false} emptyText="" maxVisible={20} onPostClick={props.onPostClick} />
      </div>
    </div>
  )
}

const RoadmapBoard = (props: RoadmapPageProps) => {
  const [plannedPosts, setPlannedPosts] = useState<Post[]>(props.plannedPosts || [])
  const [startedPosts, setStartedPosts] = useState<Post[]>(props.startedPosts || [])
  const [completedPosts, setCompletedPosts] = useState<Post[]>(props.completedPosts || [])
  const tags = props.tags || []

  const reloadPosts = useCallback(async () => {
    const [planned, started, completed] = await Promise.all([
      actions.searchPosts({ view: "planned" }),
      actions.searchPosts({ view: "started" }),
      actions.searchPosts({ view: "completed" }),
    ])
    if (planned.ok) setPlannedPosts(planned.data)
    if (started.ok) setStartedPosts(started.data)
    if (completed.ok) setCompletedPosts(completed.data)
  }, [])

  const { selectedPostId, handlePostClick, handleCloseOverlay, setIsPostDirty } = usePostOverlay({
    basePath: "/roadmap",
    onPostClosed: () => reloadPosts(),
  })

  const hasNoActivePosts = plannedPosts.length === 0 && startedPosts.length === 0

  if (hasNoActivePosts && selectedPostId === null) {
    return <RoadmapBlankState />
  }

  return (
    <div id="p-roadmap" className="page container">
      <div style={selectedPostId !== null ? { display: "none" } : undefined}>
        <VStack spacing={4}>
          <div className="c-roadmap-board">
            <RoadmapColumn status="planned" posts={plannedPosts} tags={tags} onPostClick={handlePostClick} />
            <RoadmapColumn status="started" posts={startedPosts} tags={tags} onPostClick={handlePostClick} />
            <RoadmapColumn status="completed" posts={completedPosts} tags={tags} onPostClick={handlePostClick} />
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
