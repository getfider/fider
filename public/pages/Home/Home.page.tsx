import "./Home.page.scss"
import NoDataIllustration from "@fider/assets/images/undraw-no-data.svg"

import React, { useEffect, useState } from "react"
import { Post, Tag, PostStatus } from "@fider/models"
import { Markdown, Hint, PoweredByFider, Icon, Header, Button } from "@fider/components"
import { PostsContainer } from "./components/PostsContainer"
import { useFider } from "@fider/hooks"
import { VStack } from "@fider/components/layout"
import { ShareFeedback } from "./components/ShareFeedback"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { isPostPending, setPostPending } from "./components/PostCache"
import { PostDetails, PostDetailsOverlay } from "@fider/components/PostDetails"

export interface HomePageProps {
  posts: Post[]
  tags: Tag[]
  searchNoiseWords: string[]
  countPerStatus: { [key: string]: number }
}

export interface HomePageState {
  title: string
}

const Lonely = () => {
  const fider = useFider()

  return (
    <div className="text-center">
      <Hint permanentCloseKey="at-least-3-posts" condition={fider.session.isAuthenticated && fider.session.user.isAdministrator}>
        <p>
          <Trans id="home.lonely.suggestion">
            It&apos;s recommended that you create <strong>at least 3</strong> suggestions here before sharing this site. The initial content is important to
            start engaging your audience.
          </Trans>
        </p>
      </Hint>
      <Icon sprite={NoDataIllustration} height="120" className="mt-6 mb-2" />
      <p className="text-muted">
        <Trans id="home.lonely.text">No posts have been created yet.</Trans>
      </p>
    </div>
  )
}

const HomePage = (props: HomePageProps) => {
  const fider = useFider()
  const [isShareFeedbackOpen, setIsShareFeedbackOpen] = useState(isPostPending())
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null)
  const [savedScrollPosition, setSavedScrollPosition] = useState<number>(0)

  useEffect(() => {
    // If we're showing the share feedback, make sure we clear the show pending flag (for draft posts)
    if (isShareFeedbackOpen) {
      if (isPostPending()) {
        setPostPending(false)
      }
    }
  })

  // Handle post clicks from ListPosts
  const handlePostClick = (postNumber: number, slug: string) => {
    // Save current scroll position
    setSavedScrollPosition(window.scrollY)
    setSelectedPostId(postNumber)
    window.history.pushState({ selectedPostId: postNumber }, "", `/posts/${postNumber}/${slug}`)
  }

  // Handle closing the overlay
  const handleCloseOverlay = () => {
    setSelectedPostId(null)
    window.history.pushState({}, "", "/")
    // Restore scroll position after state updates
    setTimeout(() => {
      window.scrollTo(0, savedScrollPosition)
    }, 0)
  }

  // Handle browser back/forward buttons
  useEffect(() => {
    const handlePopState = () => {
      const path = window.location.pathname
      if (path === "/" || path === "") {
        setSelectedPostId(null)
        // Restore scroll position when going back to home
        setTimeout(() => {
          window.scrollTo(0, savedScrollPosition)
        }, 0)
      } else if (path.startsWith("/posts/")) {
        // Save scroll position before opening post
        setSavedScrollPosition(window.scrollY)
        // Extract post number from URL
        const match = path.match(/\/posts\/(\d+)/)
        if (match) {
          const postNumber = parseInt(match[1], 10)
          setSelectedPostId(postNumber)
        }
      }
    }

    window.addEventListener("popstate", handlePopState)
    return () => window.removeEventListener("popstate", handlePopState)
  }, [savedScrollPosition])

  const defaultWelcomeMessage = i18n._({
    id: "home.form.defaultwelcomemessage",
    message: `We'd love to hear what you're thinking about.

What can we do better? This is the place for you to vote, discuss and share ideas.`,
  })

  const defaultInvitation = i18n._({ id: "home.form.defaultinvitation", message: "Enter your suggestion here..." })

  const isLonely = () => {
    const len = Object.keys(props.countPerStatus).length
    if (len === 0) {
      return true
    }

    if (len === 1 && PostStatus.Deleted.value in props.countPerStatus) {
      return true
    }

    return false
  }

  const handleNewPost = () => {
    setIsShareFeedbackOpen(true)
  }

  return (
    <>
      <ShareFeedback
        tags={props.tags}
        placeholder={fider.session.tenant.invitation || defaultInvitation}
        isOpen={isShareFeedbackOpen && !fider.isReadOnly}
        onClose={() => setIsShareFeedbackOpen(false)}
      />
      <div style={{ display: selectedPostId ? "none" : "block" }}>
        <Header hasInert={isShareFeedbackOpen && !fider.isReadOnly} />
        <div id="p-home" className="page container" {...(isShareFeedbackOpen && !fider.isReadOnly && { inert: "true" })}>
          <div className="p-home__welcome-col">
            <VStack spacing={2} className="p-4">
              <Markdown text={fider.session.tenant.welcomeMessage || defaultWelcomeMessage} style="full" />
              <Button className="c-input" type="submit" variant="secondary" onClick={handleNewPost}>
                {fider.session.tenant.invitation || defaultInvitation}
              </Button>
            </VStack>
            <div onClick={() => setIsShareFeedbackOpen(true)}>
              <PoweredByFider slot="home-input" className="sm:hidden md:hidden lg:block mt-3" />
            </div>
          </div>
          <div className="p-home__posts-col p-4">
            {isLonely() ? <Lonely /> : <PostsContainer posts={props.posts} tags={props.tags} countPerStatus={props.countPerStatus} onPostClick={handlePostClick} />}
            <PoweredByFider slot="home-footer" className="lg:hidden xl:hidden mt-8" />
          </div>
        </div>
      </div>

      {selectedPostId && (
        <PostDetailsOverlay onClose={handleCloseOverlay}>
          <Header />
          <div className="page container">
            <PostDetails postNumber={selectedPostId} onClose={handleCloseOverlay} />
          </div>
        </PostDetailsOverlay>
      )}
    </>
  )
}

export default HomePage
