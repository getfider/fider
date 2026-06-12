import "./Home.page.scss"
import NoDataIllustration from "@fider/assets/images/undraw-no-data.svg"
import IconPlusCircle from "@fider/assets/images/heroicons-pluscircle.svg"
import IconArrowLeft from "@fider/assets/images/heroicons-arrowleft.svg"

import React, { useEffect, useState, useRef } from "react"
import { Post, Tag, PostStatus } from "@fider/models"
import { Markdown, Hint, PoweredByFider, Icon, Header, Button } from "@fider/components"
import { PostsContainer } from "./components/PostsContainer"
import { useFider } from "@fider/hooks"
import { basePath } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import { ShareFeedback } from "./components/ShareFeedback"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { isPostPending, setPostPending } from "./components/PostCache"
import { PostDetails } from "@fider/components/PostDetails"

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
  const postsContainerRef = useRef<PostsContainer>(null)
  const [isShareFeedbackOpen, setIsShareFeedbackOpen] = useState(isPostPending())
  const [selectedPostId, setSelectedPostId] = useState<number | null>(null)
  const [savedScrollPosition, setSavedScrollPosition] = useState<number>(0)
  const [isPostDirty, setIsPostDirty] = useState(false)
  const [savedSearch, setSavedSearch] = useState("")

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
    // Save current scroll position and search params
    setSavedScrollPosition(window.scrollY)
    setSavedSearch(window.location.search)
    setSelectedPostId(postNumber)
    setLastOpenedPostId(postNumber) // Track which post was opened
    setIsPostDirty(false) // Reset dirty flag when opening overlay
    window.history.pushState({ selectedPostId: postNumber }, "", `${basePath()}/posts/${postNumber}/${slug}`)
  }

  // Handle closing the overlay
  const handleCloseOverlay = () => {
    setSelectedPostId(null)
    window.history.pushState({}, "", `${basePath()}/${savedSearch}`)
  }

  // Track which post was opened so we can update just that one
  const [lastOpenedPostId, setLastOpenedPostId] = useState<number | null>(null)

  // Update single post when closing overlay if data changed, and always restore scroll
  useEffect(() => {
    if (selectedPostId === null && lastOpenedPostId !== null) {
      if (isPostDirty && postsContainerRef.current) {
        postsContainerRef.current.updateSinglePost(lastOpenedPostId)
        setIsPostDirty(false)
      }
      setLastOpenedPostId(null)

      // Always restore scroll position when returning to home
      setTimeout(() => {
        window.scrollTo(0, savedScrollPosition)
      }, 0)
    }
  }, [selectedPostId, lastOpenedPostId, isPostDirty, savedScrollPosition])

  // Handle browser back/forward buttons
  useEffect(() => {
    const handlePopState = () => {
      const path = window.location.pathname
      if (path === "/" || path === "") {
        setSelectedPostId(null)
        // Scroll restoration is handled in the useEffect above
      } else if (path.startsWith("/posts/")) {
        // Save scroll position and search params before opening post
        setSavedScrollPosition(window.scrollY)
        setSavedSearch(window.location.search)
        // Extract post number from URL
        const match = path.match(/\/posts\/(\d+)/)
        if (match) {
          const postNumber = parseInt(match[1], 10)
          setSelectedPostId(postNumber)
          setLastOpenedPostId(postNumber) // Track which post was opened
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

  const parseWelcomeHeader = (text: string): JSX.Element[] => {
    const parts: JSX.Element[] = []
    let currentIndex = 0
    const regex = /_([^_]+)_/g
    let match: RegExpExecArray | null

    while ((match = regex.exec(text)) !== null) {
      // Add text before the match
      if (match.index > currentIndex) {
        parts.push(<span key={currentIndex}>{text.slice(currentIndex, match.index)}</span>)
      }
      // Add the highlighted text
      parts.push(
        <span key={match.index} className="header-emphasis">
          {match[1]}
        </span>
      )
      currentIndex = regex.lastIndex
    }

    // Add remaining text
    if (currentIndex < text.length) {
      parts.push(<span key={currentIndex}>{text.slice(currentIndex)}</span>)
    }

    return parts
  }

  return (
    <>
      <ShareFeedback
        tags={props.tags}
        placeholder={fider.session.tenant.invitation || defaultInvitation}
        isOpen={isShareFeedbackOpen && !fider.isReadOnly}
        onClose={() => setIsShareFeedbackOpen(false)}
      />
      <div>
        <Header hasInert={isShareFeedbackOpen && !fider.isReadOnly} />
        <div
          id="p-home"
          className="page container"
          style={selectedPostId !== null ? { display: "none" } : undefined}
          {...(isShareFeedbackOpen && !fider.isReadOnly && { inert: "true" })}
        >
          <div className="p-home__welcome-col">
            <VStack spacing={6}>
              <div>
                {fider.session.tenant.welcomeHeader && <h1 className="p-home__welcome-title mb-5">{parseWelcomeHeader(fider.session.tenant.welcomeHeader)}</h1>}
                <Markdown className="p-home__welcome-body" text={fider.session.tenant.welcomeMessage || defaultWelcomeMessage} style="full" />
              </div>
            </VStack>
            <div>
              <PoweredByFider slot="home-input" className="sm:hidden md:hidden lg:block mt-3" />
            </div>
          </div>
          <div className="p-home__posts-col">
            <button className="p-home__add-idea-btn" onClick={handleNewPost}>
              <HStack spacing={4} align="center">
                <Icon sprite={IconPlusCircle} className="p-home__add-idea-icon" />
                <span>{fider.session.tenant.invitation || defaultInvitation}</span>
              </HStack>
            </button>
            {isLonely() ? (
              <Lonely />
            ) : (
              <PostsContainer
                ref={postsContainerRef}
                posts={props.posts}
                tags={props.tags}
                countPerStatus={props.countPerStatus}
                onPostClick={handlePostClick}
              />
            )}
            <PoweredByFider slot="home-footer" className="lg:hidden xl:hidden mt-8" />
          </div>
        </div>
        {selectedPostId !== null && (
          <div className="page container">
            <Button onClick={handleCloseOverlay} variant="link">
              <HStack spacing={2}>
                <Icon sprite={IconArrowLeft} />
                <span className="text-body clickable text-blue-600 hover">
                  <Trans id="postdetails.backtoall">Back to all suggestions</Trans>
                </span>
              </HStack>
            </Button>
            <PostDetails postNumber={selectedPostId} onDataChanged={() => setIsPostDirty(true)} />
          </div>
        )}
      </div>
    </>
  )
}

export default HomePage
