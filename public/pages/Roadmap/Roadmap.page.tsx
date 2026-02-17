import React, { useState } from "react"
import { Post } from "@fider/models"
import { Header, Icon } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"
import { i18n } from "@lingui/core"
import IconChatAlt2 from "@fider/assets/images/heroicons-chat-alt-2.svg"
import "./Roadmap.page.scss"

const POSTS_PER_PAGE = 10

interface RoadmapPageProps {
  planned: Post[]
  started: Post[]
  completed: Post[]
}

interface ColumnProps {
  title: string
  posts: Post[]
  colorClass: string
  icon: string
}

const RoadmapColumn = ({ title, posts, colorClass, icon }: ColumnProps) => {
  const [visibleCount, setVisibleCount] = useState(POSTS_PER_PAGE)
  const visiblePosts = posts.slice(0, visibleCount)
  const hasMore = posts.length > visibleCount
  const remaining = posts.length - visibleCount

  const handleLoadMore = () => {
    setVisibleCount((prev) => prev + POSTS_PER_PAGE)
  }

  return (
    <div className={`c-roadmap-column c-roadmap-column--${colorClass}`}>
      <div className="c-roadmap-column__header">
        <HStack justify="between" align="center">
          <span className="c-roadmap-column__title">
            <span className="c-roadmap-column__icon">{icon}</span>
            {title}
          </span>
          <span className="c-roadmap-column__count">{posts.length}</span>
        </HStack>
      </div>
      <div className="c-roadmap-column__posts">
        {visiblePosts.map((post) => (
          <a key={post.id} href={`/posts/${post.number}/${post.slug}`} className="c-roadmap-card">
            <VStack spacing={2}>
              <div className="c-roadmap-card__title">{post.title}</div>
              <HStack spacing={4} className="c-roadmap-card__meta">
                <span className="text-muted">{post.votesCount} votes</span>
                {post.commentsCount > 0 && (
                  <HStack spacing={1} className="text-muted">
                    <Icon sprite={IconChatAlt2} className="h-4" />
                    <span>{post.commentsCount}</span>
                  </HStack>
                )}
              </HStack>
            </VStack>
          </a>
        ))}
        {posts.length === 0 && <div className="c-roadmap-column__empty">{i18n._({ id: "roadmap.empty", message: "No posts" })}</div>}
        {hasMore && (
          <button className="c-roadmap-column__load-more" onClick={handleLoadMore}>
            {i18n._({ id: "roadmap.loadmore", message: "View {0} more" }, { 0: remaining })}
          </button>
        )}
      </div>
    </div>
  )
}

const RoadmapPage = (props: RoadmapPageProps) => {
  return (
    <>
      <Header />
      <div id="p-roadmap" className="page container">
        <div className="c-roadmap-board">
          <RoadmapColumn
            title={i18n._({ id: "roadmap.column.planned", message: "Planned" })}
            posts={props.planned || []}
            colorClass="planned"
            icon="○"
          />
          <RoadmapColumn
            title={i18n._({ id: "roadmap.column.started", message: "In Progress" })}
            posts={props.started || []}
            colorClass="started"
            icon="◐"
          />
          <RoadmapColumn
            title={i18n._({ id: "roadmap.column.completed", message: "Completed" })}
            posts={props.completed || []}
            colorClass="completed"
            icon="●"
          />
        </div>
      </div>
    </>
  )
}

export default RoadmapPage
