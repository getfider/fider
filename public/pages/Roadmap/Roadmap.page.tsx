import "./Roadmap.page.scss"
import "@fider/pages/Home/components/PostsContainer.scss"
import IconArrowLeft from "@fider/assets/images/heroicons-arrowleft.svg"

import React, { useState, useCallback } from "react"
import { Post, Tag } from "@fider/models"
import { Header, ResponseLozenge, Button, Icon } from "@fider/components"
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

const RoadmapUpsell = () => {
  const fider = useFider()
  const isAdmin = fider.session.isAuthenticated && fider.session.user.isAdministrator
  const showBillingLink = isAdmin && fider.settings.isBillingEnabled

  return (
    <div className="page container">
      <div className="text-center py-8">
        <h2 className="text-display mb-2">
          <Trans id="roadmap.upsell.title">Roadmap</Trans>
        </h2>
        <p className="text-muted mb-4">
          <Trans id="roadmap.upsell.description">
            The roadmap view is a PRO feature that lets you visualize your planned, in-progress, and completed posts in a kanban-style board.
          </Trans>
        </p>
        {showBillingLink && (
          <a href="/admin/billing" className="text-link">
            <Trans id="roadmap.upsell.billing">Upgrade to PRO</Trans>
          </a>
        )}
      </div>
    </div>
  )
}

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
